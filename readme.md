# Collectify

[![Go Report Card](https://goreportcard.com/badge/github.com/Jinvic/collectify)](https://goreportcard.com/report/github.com/Jinvic/collectify)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> [!WARNING]
> 此项目正在开发中，仅完成基础功能。

**Collectify** 是一个轻量级、易于使用的个人收藏管理系统，帮助你整理和管理你的书籍、音乐、影视等各类收藏品。无论是为了记录你的阅读进度，还是为了整理你喜爱的音乐专辑，Collectify 都能为你提供所需的一切工具。

## 功能特性

- **多类型收藏管理**：支持书籍、音乐、影视等多种类型的收藏品管理
- **自定义字段**：为不同类别的收藏品设置自定义字段
- **标签系统**：通过标签对收藏品进行分类和检索
- **收藏夹功能**：创建不同的收藏夹来组织你的收藏品
- **状态追踪**：跟踪收藏品的完成状态（待完成、进行中、已完成等）
- **搜索功能**：强大的搜索功能，支持按名称、标签、字段值等搜索
- **权限控制**：可选的身份验证功能，保护你的数据安全
- **Web UI**: 基于 React 和 Material-UI 的现代化用户界面
- **数据持久化**: 使用 SQLite 数据库存储数据
- **RESTful API**: 提供 API 接口，方便与其他系统集成
- **Docker 支持**: 提供 Dockerfile 和 docker-compose.yml，方便部署

## 快速开始

启动应用后，通过 Web 界面访问：http://localhost:8080

### 使用 Docker（推荐）

```bash
# 克隆项目
git clone https://github.com/Jinvic/collectify.git
cd collectify

# 复制并修改环境配置文件
cp .env.example .env
# 编辑 .env 文件以满足你的需求

# 启动服务
docker-compose up -d

# 访问应用
# 默认地址: http://localhost:8080
```

### 从源码构建

#### 前置要求

- Go 1.23+
- Node.js 18+
- pnpm

#### 构建步骤

```bash
# 克隆项目
git clone https://github.com/Jinvic/collectify.git
cd collectify

# 构建前端
cd web
pnpm install
pnpm run build
cd ..

# 构建后端
go build -o collectify .

# 运行应用
./collectify
```

## 配置说明

Collectify 支持通过环境变量进行配置。可以查看 `.env.example` 文件了解所有可用配置项。

所有环境变量都以 `COLLECTIFY_` 为前缀。

### 数据库配置

- `COLLECTIFY_DATABASE_TYPE`：数据库类型
  - 默认值：`sqlite`
  
- `COLLECTIFY_DATABASE_DSN`：数据库连接字符串
  - SQLite 示例：`collectify.db` （相对路径）或 `/app/data/collectify.db` （绝对路径）

### 服务器配置

- `COLLECTIFY_SERVER_PORT`：服务器监听端口
  - 默认值：`8080`
  
- `COLLECTIFY_SERVER_MODE`：服务器运行模式
  - 可选值：`release`、`debug`
  - 默认值：`release`
  - `debug` 模式提供详细日志信息
  - `release` 模式为生产环境优化

### 认证配置

- `COLLECTIFY_AUTH_ENABLE`：是否启用认证功能
  - 可选值：`true`、`false`
  - 默认值：`false`
  - 启用后需要登录才能使用管理功能

- `COLLECTIFY_AUTH_JWT_SECRET`：JWT 密钥
  - 默认值：空
  - 生产环境必须设置此值
  - 开发模式下如果留空会自动生成随机密钥
  - 可使用 `openssl rand -base64 32` 生成安全密钥

- `COLLECTIFY_AUTH_EXPIRE_DAY`：JWT Token 过期时间（天）
  - 默认值：`15`
  - 设置为 `0` 表示永不过期

## 开发指南

### 技术栈

- **后端**:
  - [Go](https://golang.org/)
  - [Gin Web Framework](https://github.com/gin-gonic/gin)
  - [GORM](https://gorm.io/)
  - [SQLite](https://www.sqlite.org/)
  - [env](https://github.com/caarlos0/env)
  - [Cobra](https://github.com/spf13/cobra)
  - [JWT](https://github.com/golang-jwt/jwt/v5)
- **前端**:
  - [React](https://reactjs.org/)
  - [React Router](https://reactrouter.com/)
  - [Material-UI (MUI)](https://mui.com/)
  - [TanStack Query (React Query)](https://tanstack.com/query/)
  - [Axios](https://axios-http.com/)

### 数据模型

Collectify 的核心数据模型包括：

- **Category（类别）**：收藏品的类别，如书籍、电影、音乐等
- **Item（收藏品）**：具体的收藏品，如某本书、某部电影
- **Field（字段）**：自定义字段，用于扩展收藏品信息
- **Tag（标签）**：标签，用于标记和分类收藏品
- **Collection（收藏夹）**：收藏夹，用于组织收藏品

### 项目结构

```bash
collectify/
├── cmd/           # 命令行接口
├── internal/      # 核心业务逻辑
│   ├── cli/       # 命令行实现
│   ├── config/    # 配置管理
│   ├── conn/      # 数据库连接
│   ├── dao/       # 数据访问对象
│   ├── handler/   # HTTP 处理器
│   ├── middleware/ # 中间件
│   ├── model/     # 数据模型
│   ├── pkg/       # 工具包
│   ├── router/    # 路由配置
│   └── service/   # 业务逻辑
├── web/           # 前端代码
├── go.mod         # Go 模块文件
└── main.go        # 程序入口
```

### 开发环境搭建

```bash
# 安装依赖
go mod tidy

# 前端开发
cd web
pnpm install
pnpm start
```

## TODO List

- [ ] 回收站功能
- [ ] 字段关联藏品
- [ ] 藏品设置私密
- [ ] 名称别名（支持搜索）
- [ ] 完善 API 接口文档
- [ ] 数据导入/导出功能
- [ ] 更多数据库支持
- [ ] 多用户支持

## 许可证

本项目遵循 MIT 许可证。详情请参见 [LICENSE](LICENSE) 文件。
