// Package utility 提供模型 Key 的 AES-256-GCM 加解密能力.
//
// KEK 来源优先级：环境变量 MODEL_KEY_ENC_KEY > 配置项 modelKey.encKey。
// 为兼容手工配置的任意长度口令：32 字节直接用；base64/hex 解码后
// 为 32 字节亦可；其他一律经 SHA-256 派生为 32 字节，避免启动即崩。
package utility

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"os"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

func modelKeyKEK(ctx context.Context) ([]byte, error) {
	raw := strings.TrimSpace(os.Getenv("MODEL_KEY_ENC_KEY"))
	if raw == "" && ctx != nil {
		raw = strings.TrimSpace(g.Cfg().MustGet(ctx, "modelKey.encKey", "").String())
	}
	if raw == "" {
		return nil, errors.New("missing MODEL_KEY_ENC_KEY (env) or modelKey.encKey (config)")
	}
	return deriveKEK(raw), nil
}

// deriveKEK 口令派生 32 字节 KEK：base64/hex 解码后为 32 字节直接用，
// 32 字节原文直接用，其他一律经 SHA-256 派生，避免启动即崩。
func deriveKEK(raw string) []byte {
	if decoded, err := base64.StdEncoding.DecodeString(raw); err == nil && len(decoded) == 32 {
		return decoded
	}
	if decoded, err := hex.DecodeString(raw); err == nil && len(decoded) == 32 {
		return decoded
	}
	if len(raw) == 32 {
		return []byte(raw)
	}
	sum := sha256.Sum256([]byte(raw))
	return sum[:]
}

// KEKProvider 由 service 层注册：从 db_system_configs 取入库 KEK。
// utility 不 import dao（防循环依赖），启动时由 sysconfig 包 init 注册。
// 未注册（单测/极简启动）则返回 ""，回退 env/config。
var KEKProvider func(ctx context.Context) string

// kekFromRaw 显式 KEK（入库值优先）：非空则直接派生，为空回退 DB，再回退 env/config。
// 调用方（user/sysconfig）从 db_system_configs 取到 KEK 后走此路径，
// utility 自身不 import dao，保持底层包无 DB 依赖。
func kekFromRaw(ctx context.Context, raw string) ([]byte, error) {
	if trimmed := strings.TrimSpace(raw); trimmed != "" {
		return deriveKEK(trimmed), nil
	}
	if KEKProvider != nil {
		if dbKEK := strings.TrimSpace(KEKProvider(ctx)); dbKEK != "" {
			return deriveKEK(dbKEK), nil
		}
	}
	return modelKeyKEK(ctx)
}

func modelKeyGCM(ctx context.Context) (cipher.AEAD, error) {
	key, err := modelKeyKEK(ctx)
	if err != nil {
		return nil, err
	}
	return gcmFromKey(key)
}

func gcmFromKey(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

// EncryptModelKey 加密明文 key，返回 base64(ciphertext)、base64(nonce)。
func EncryptModelKey(ctx context.Context, plaintext string) (cipherB64, nonceB64 string, err error) {
	if plaintext == "" {
		return "", "", errors.New("apiKey is empty")
	}
	aead, err := modelKeyGCM(ctx)
	if err != nil {
		return "", "", err
	}
	return sealWithAEAD(aead, plaintext)
}

// EncryptModelKeyWithKEK 用显式 KEK（入库值）加密；kekRaw 为空回退 env/config。
func EncryptModelKeyWithKEK(ctx context.Context, plaintext, kekRaw string) (cipherB64, nonceB64 string, err error) {
	if plaintext == "" {
		return "", "", errors.New("apiKey is empty")
	}
	key, err := kekFromRaw(ctx, kekRaw)
	if err != nil {
		return "", "", err
	}
	aead, err := gcmFromKey(key)
	if err != nil {
		return "", "", err
	}
	return sealWithAEAD(aead, plaintext)
}

func sealWithAEAD(aead cipher.AEAD, plaintext string) (cipherB64, nonceB64 string, err error) {
	nonce := make([]byte, aead.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", "", err
	}
	sealed := aead.Seal(nil, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed),
		base64.StdEncoding.EncodeToString(nonce), nil
}

// DecryptModelKey 解密 DB 中的 (cipherB64, nonceB64)，只在构造 model 瞬间调用，
// 调用方 MUST 不落盘、不打日志。
func DecryptModelKey(ctx context.Context, cipherB64, nonceB64 string) (string, error) {
	aead, err := modelKeyGCM(ctx)
	if err != nil {
		return "", err
	}
	return openWithAEAD(aead, cipherB64, nonceB64)
}

// DecryptModelKeyWithKEK 用显式 KEK（入库值）解密；kekRaw 为空回退 env/config。
func DecryptModelKeyWithKEK(ctx context.Context, cipherB64, nonceB64, kekRaw string) (string, error) {
	key, err := kekFromRaw(ctx, kekRaw)
	if err != nil {
		return "", err
	}
	aead, err := gcmFromKey(key)
	if err != nil {
		return "", err
	}
	return openWithAEAD(aead, cipherB64, nonceB64)
}

func openWithAEAD(aead cipher.AEAD, cipherB64, nonceB64 string) (string, error) {
	sealed, err := base64.StdEncoding.DecodeString(cipherB64)
	if err != nil {
		return "", errors.New("invalid cipher encoding")
	}
	nonce, err := base64.StdEncoding.DecodeString(nonceB64)
	if err != nil {
		return "", errors.New("invalid nonce encoding")
	}
	plain, err := aead.Open(nil, nonce, sealed, nil)
	if err != nil {
		return "", errors.New("decrypt model key failed: wrong MODEL_KEY_ENC_KEY or corrupted data")
	}
	return string(plain), nil
}

// MaskAPIKey 脱敏回显：sk-cp-****9O7p。
func MaskAPIKey(key string) string {
	if key == "" {
		return ""
	}
	if len(key) <= 8 {
		return "****"
	}
	return key[:5] + "****" + key[len(key)-4:]
}
