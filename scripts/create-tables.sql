USE gin_todo;
CREATE TABLE IF NOT EXISTS todos (
  id         INT AUTO_INCREMENT NOT NULL,
  task       VARCHAR(255) NOT NULL,
  completed  BOOLEAN DEFAULT FALSE,
  PRIMARY KEY (`id`)
)CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

INSERT INTO todos (task, completed)
SELECT '尝试一下这个样例Todo，点击右边来完成它', false
WHERE NOT EXISTS (SELECT 1 FROM todos WHERE task = '尝试一下这个样例Todo，点击右边来完成它');

INSERT INTO todos (task, completed)
SELECT '在上面的文本框填写一个Todo，然后添加它', false
WHERE NOT EXISTS (SELECT 1 FROM todos WHERE task = '在上面的文本框填写一个Todo，然后添加它');

INSERT INTO todos (task, completed)
SELECT '我会了！', false
WHERE NOT EXISTS (SELECT 1 FROM todos WHERE task = '我会了！');