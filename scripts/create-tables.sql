CREATE DATABASE IF NOT EXISTS gin_todo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE gin_todo;
CREATE TABLE IF NOT EXISTS todos (
  id         INT AUTO_INCREMENT NOT NULL,
  task       VARCHAR(255) NOT NULL,
  completed  BOOLEAN DEFAULT FALSE,
  PRIMARY KEY (`id`)
)CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 先清空再插入（避免重复数据）
DELETE FROM todos;

INSERT INTO todos
  (task, completed)
VALUES
  ('尝试一下这个样例Todo，点击右边来完成它', false),
  ('在上面的文本框填写一个Todo，然后添加它', false),
  ('我会了！', false);