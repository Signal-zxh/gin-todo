
CREATE DATABASE IF NOT EXISTS gin_todo;
USE gin_todo;
DROP TABLE IF EXISTS todos;
CREATE TABLE todos (
  id         INT AUTO_INCREMENT NOT NULL,
  task       VARCHAR(255) NOT NULL,
  completed  BOOLEAN DEFAULT FALSE,
  PRIMARY KEY (`id`)
);

INSERT INTO todos
  (task, completed)
VALUES
  ('尝试一下这个样例Todo，点击右边来完成它', false),
  ('在上面的文本框填写一个Todo，然后添加它', false),
  ('我会了！', false);