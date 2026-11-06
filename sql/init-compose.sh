#!/bin/sh
# compose PG 初始化入口（/docker-entrypoint-initdb.d 只在 volume 为空时执行一次）。
# 顺序：扩展 -> 基表 -> 全部 migration（migration 全部幂等，可重复执行）。
set -e
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<'EOSQL'
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS zhparser;
-- zhparser FTS 配置：扩展自带 parser，configuration 需手动建。
-- scws token 类型名（与 amutu/zhparser 文档一致）：n 名词、v 动词、a 形容词、
-- i 成语、e 叹词、l 习惯用语、t 时间词、其余 r* 归 simple。
-- 旧链路中文标题/作者搜索用；新链路（ES）不依赖此配置，失败不影响主流程。
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_ts_config WHERE cfgname = 'zhparser') THEN
    CREATE TEXT SEARCH CONFIGURATION zhparser (PARSER = zhparser);
    ALTER TEXT SEARCH CONFIGURATION zhparser ADD MAPPING FOR n,v,a,i,e,l,t WITH simple;
  END IF;
EXCEPTION WHEN OTHERS THEN
  RAISE NOTICE 'zhparser FTS config skipped: %', SQLERRM;
END $$;
EOSQL
for f in /docker-entrypoint-initdb.d/schema/*.sql; do
  echo "applying $f"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f "$f"
done
echo "DB init done."
