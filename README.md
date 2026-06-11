# Gin Todo App

一个基于 Go 语言 Gin 框架开发的轻量级待办事项 Web 应用，使用 MySQL 数据库持久化数据，并集成 Redis 缓存提升性能。

![运行截图](screenshot/demo.png)

## ✨ 功能特性

- ✅ **添加任务** - 快速添加新的待办事项
- ✅ **标记完成** - 一键标记任务为已完成状态
- ✅ **删除任务** - 删除不需要的任务
- ✅ **数据持久化** - 使用 MySQL 存储所有数据
- ✅ **Redis 缓存** - 列表查询结果缓存，提升响应速度

## 🛠️ 技术栈

| 组件 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.26+ |
| 框架 | Gin | Latest |
| 数据库 | MySQL | 9.3+ |
| 缓存 | Redis | Latest |

## 📁 项目结构

```
gin-todo/
├── main.go              # 主入口文件，包含路由和业务逻辑
├── templates/           # HTML 模板目录
│   └── list.html        # 待办事项列表页面
├── scripts/             # 数据库脚本
│   └── create-tables.sql # 初始化数据库和表结构
├── screenshot/          # 截图资源
│   └── demo.png         # 运行截图
├── compose.yaml         # Docker Compose 配置
├── Dockerfile           # 应用 Docker 配置
├── Dockerfile.mysql     # MySQL Docker 配置
├── go.mod               # Go 依赖管理
├── go.sum               # 依赖校验文件
└── README.md            # 项目说明文档
```

## 🚀 快速开始

### 方式一：使用 Docker Compose（推荐）

```bash
# 启动所有服务
docker compose up --build

# 访问应用
# 打开浏览器访问 http://localhost:8080
```

### 方式二：本地开发运行

#### 前置条件

- 安装 MySQL 和 Redis 并启动服务
- MySQL（后台运行）：`brew services start mysql`
- Redis（前台运行）：`redis-stack-server`
- 确保你有数据库的 `root` 权限，或已创建相应的数据库用户

#### 配置环境变量

```bash
export DBUSER=root
export DBPASS=你的密码
# 可选，默认为 127.0.0.1
export DBHOST=127.0.0.1  
```

#### 初始化数据库

首次运行前，执行以下命令创建数据库和表：

```bash
mysql -u root -p < scripts/create-tables.sql
```

#### 启动服务

```bash
go run .
```

打开浏览器访问 http://localhost:8080

## 🔌 API 接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/` | 获取待办事项列表（HTML 页面） |
| POST | `/add` | 添加新任务 |
| GET | `/complete/:id` | 标记任务为已完成 |
| GET | `/delete/:id` | 删除指定任务 |

## 📝 使用说明

1. 在首页的输入框中填写待办事项内容
2. 点击「添加」按钮将任务添加到列表
3. 点击任务右侧的「完成」链接标记任务为已完成
4. 点击任务右侧的「删除」链接移除任务

## 🐳 Docker 部署

使用 Docker Compose 一键启动完整的应用栈（包含应用、MySQL、Redis）：

```bash
docker compose up --build
```

启动后访问 http://localhost:8080 即可使用应用。

> **说明**：所有配置已在 `compose.yaml` 中自动完成，无需额外设置环境变量。

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！