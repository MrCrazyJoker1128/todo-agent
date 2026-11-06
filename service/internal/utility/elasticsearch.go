package utility

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// esClient 客户端单例
var esClient *elasticsearch.Client

// ESConfig ES连接配置
type ESConfig struct {
	Addresses []string
	Username  string
	Password  string
}

// getESConfig 从GoFrame配置读取ES配置，无配置时使用默认值
func getESConfig() ESConfig {
	config := ESConfig{
		Addresses: []string{"http://127.0.0.1:9200"},
		Username:  "elastic",
		Password:  "",
	}
	v := g.Config("default")
	ctx := gctx.New()
	if v != nil {
		if addresses, err := v.Get(ctx, "elasticsearch.addresses"); err == nil {
			config.Addresses = addresses.Strings()
		}
		if username, err := v.Get(ctx, "elasticsearch.username"); err == nil {
			config.Username = username.String()
		}
		if password, err := v.Get(ctx, "elasticsearch.password"); err == nil {
			config.Password = password.String()
		}
	}
	return config
}

// getESClient 创建ES客户端单例，已存在则直接返回
func getESClient() (*elasticsearch.Client, error) {
	if esClient != nil {
		return esClient, nil
	}

	config := getESConfig()

	esCfg := elasticsearch.Config{
		Addresses: config.Addresses,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	if config.Username != "" {
		esCfg.Username = config.Username
		esCfg.Password = config.Password
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	esClient = client
	return esClient, nil
}

// ES 获取ES客户端实例
func ES() (*elasticsearch.Client, error) {
	return getESClient()
}

// ESPing 检查ES集群连接是否正常
func ESPing() error {
	client, err := getESClient()
	if err != nil {
		return err
	}
	res, err := client.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping elasticsearch: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("elasticsearch ping error: %s", res.String())
	}
	return nil
}

// ESIndex 索引或更新文档，id为空则自动生成
func ESIndex(ctx context.Context, index, id string, document any) error {
	client, err := getESClient()
	if err != nil {
		return err
	}

	body, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	// wait_for:写入后等到可被搜索再返回;否则前端紧接着的列表查询看不到刚写的文档
	//(ES 默认 refresh_interval=1s),表现为上传成功后列表不刷新。
	res, err := client.Index(index, bytes.NewReader(body),
		client.Index.WithContext(ctx),
		client.Index.WithDocumentID(id),
		client.Index.WithRefresh("wait_for"),
	)
	if err != nil {
		return fmt.Errorf("es index error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("es index error: %s", res.String())
	}
	return nil
}

// ESSearch 执行搜索查询，返回原始响应字节
func ESSearch(ctx context.Context, index string, query any) ([]byte, error) {
	client, err := getESClient()
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(index),
		client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("es search error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("es search error: %s", res.String())
	}

	return io.ReadAll(res.Body)
}

// ESDelete 根据id删除文档，文档不存在返回成功
func ESDelete(ctx context.Context, index, id string) error {
	client, err := getESClient()
	if err != nil {
		return err
	}

	res, err := client.Delete(index, id, client.Delete.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("es delete error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("es delete error: %s", res.String())
	}
	return nil
}

// ESGet 根据id获取文档，返回原始响应字节
func ESGet(ctx context.Context, index, id string) ([]byte, error) {
	client, err := getESClient()
	if err != nil {
		return nil, err
	}

	res, err := client.Get(index, id, client.Get.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("es get error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return nil, fmt.Errorf("es get error: %s", res.String())
	}

	return io.ReadAll(res.Body)
}
