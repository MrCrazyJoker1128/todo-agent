-- =============================================================================
-- 新库种子数据：首次启动自动建。幂等（CONFLICT DO NOTHING），重复执行安全。
-- 1. 默认 KEK：固定值 todo-agent-default-kek（开发/演示用；生产请在设置页轮换）。
-- 2. 演示账号：demo@todo-agent.local / 123456（MD5: e10adc3949ba59abbe56e057f20f883e）。
--    前端登录页 MD5 后传输，后端按 password_hash 明文比对，故此处存 MD5 值。
-- 生产部署如不需要演示账号，删掉第 2 段即可。
-- =============================================================================

INSERT INTO db_system_configs (key, value) VALUES
  ('model_key_enc_key', 'todo-agent-default-kek')
ON CONFLICT (key) DO NOTHING;

INSERT INTO db_users (id, email, password_hash, nickname) VALUES
  (gen_random_uuid(), 'demo@todo-agent.local', 'e10adc3949ba59abbe56e057f20f883e', '演示用户')
ON CONFLICT DO NOTHING;
