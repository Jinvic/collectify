# Collectify

[![Go Report Card](https://goreportcard.com/badge/github.com/Jinvic/collectify)](https://goreportcard.com/report/github.com/Jinvic/collectify)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> [!WARNING]
> 此项目正在开发中，仅完成部分基础功能。

**Collectify** 是一个轻量级、易于使用的个人收藏管理系统，帮助你整理和管理你的书籍、音乐、影视等各类收藏品。无论是为了记录你的阅读进度，还是为了整理你喜爱的音乐专辑，Collectify 都能为你提供所需的一切工具。

## 功能特性

* **多类型收藏管理**: 灵活的分类和字段系统，可自定义管理书籍、音乐、电影等各种收藏品。
* **Web UI**: 基于 React 和 Material-UI 的现代化用户界面。
* **数据持久化**: 使用 SQLite 数据库存储数据。
* **回收站**: 可选的回收站功能，防止误删数据。
* **RESTful API**: 提供 API 接口，方便与其他系统集成。
* **Docker 支持**: 提供 Dockerfile 和 docker-compose.yml，方便部署。
* **简单鉴权**: 可选的用户认证功能，保护您的收藏数据。

## 快速开始

### 使用 Docker Compose (推荐)

这是最简单的启动方式。

1. **获取docker compose 文件**:
    * 克隆项目:
  
        ```bash
        git clone https://github.com/Jinvic/collectify.git
        cd collectify
        ```

    * 或者单独下载所需文件:
  
        ```bash
        mkdir collectify
        cd collectify
        wget https://raw.githubusercontent.com/Jinvic/collectify/master/docker-compose.yml
        wget https://raw.githubusercontent.com/Jinvic/collectify/master/.env.example
        ```

2. **配置环境变量**:
    * 复制 `.env.example` 文件为 `.env`。
    * 根据需要修改 `.env` 文件中的配置。
    * 如果要启用认证功能，需要设置 `COLLECTIFY_AUTH_ENABLE=true`，并配置 `COLLECTIFY_AUTH_JWT_SECRET`。
      * 在生产环境中，必须手动设置 JWT 密钥。
      * 在开发环境中，如果不设置 JWT 密钥，系统会自动生成一个随机密钥。

3. **启动服务**:

    ```bash
    docker-compose up -d
    ```

4. **访问应用**:
    * Web UI: `http://localhost:8080`
    * API: `http://localhost:8080/api/`

### 本地开发

#### 后端

1. **安装 Go**: 确保你已安装 Go 1.23 或更高版本。
2. **获取依赖**: 在项目根目录下运行 `go mod tidy`。
3. **运行后端**:
    * 开发模式: `go run main.go` (默认监听 :8080)。
    * 或构建后运行:

        ```bash
        go build -o collectify .
        ./collectify # (Linux/macOS)
        # 或
        collectify.exe # (Windows)
        ```

#### 前端 (Web UI)

前端使用 React 构建，位于 `web` 目录。

1. **安装 Node.js**: 确保你已安装 Node.js (建议 LTS 版本)。
2. **安装 Pnpm**: 确保你已安装 Pnpm。
3. **安装前端依赖**: 在 `web` 目录下运行 `pnpm install`。
4. **开发**:
    * 启动前端开发服务器: 在 `web` 目录下运行 `pnpm start`。这将在 `http://localhost:3000` 启动一个热重载的开发服务器。它通过 `package.json` 中的 `proxy` 设置代理 API 请求到后端 (`http://localhost:8080`)。
    * 同时确保后端服务 (`go run main.go`) 正在运行。
5. **构建**:
    * 构建生产版本: 在 `web` 目录下运行 `pnpm run build`。这会将所有静态资源生成到 `web/build` 目录。
    * 构建后端时，它会将 `web/build` 目录嵌入到二进制文件中，使得应用可以作为一个整体部署。

## 项目结构

```bash
├── cmd                 # Cobra 命令行接口
├── internal            # 核心代码
│   ├── cli             # 命令行具体实现
│   ├── config          # 配置管理 (caarlos0/env)
│   ├── dao             # 数据访问对象 (泛型封装)
│   ├── conn            # 数据库初始化 (GORM)
│   ├── handler         # HTTP 处理函数 (Gin)
│   ├── middleware      # Gin 中间件
│   ├── model           # 数据模型定义
│   ├── pkg             # 公共工具包
│   ├── router          # 路由定义 (Gin)
│   ├── service         # 业务逻辑层
│   └── test            # 测试文件
├── web                 # React 前端代码
│   ├── build           # (构建产物) 前端静态文件
│   ├── public          # 公共静态资源
│   └── src             # 前端源代码
├── .env.example        # 环境变量示例文件
├── docker-compose.yml  # Docker Compose 配置
├── Dockerfile          # Docker 构建文件
├── go.mod              # Go 模块依赖
├── main.go             # 程序入口
└── ...
```

## 技术栈

* **后端**:
  * [Go](https://golang.org/)
  * [Gin Web Framework](https://github.com/gin-gonic/gin)
  * [GORM](https://gorm.io/)
  * [SQLite](https://www.sqlite.org/) (数据库)
  * [env](https://github.com/caarlos0/env) (配置)
  * [Cobra](https://github.com/spf13/cobra) (CLI)
* **前端**:
  * [React](https://reactjs.org/)
  * [React Router](https://reactrouter.com/)
  * [Material-UI (MUI)](https://mui.com/)
  * [TanStack Query (React Query)](https://tanstack.com/query/)
  * [Axios](https://axios-http.com/)

## TODO List

* [ ] 标签
* [ ] 收藏夹
* [ ] 回收站
* [ ] 简单鉴权

## ~~Q&A~~ 碎碎念

* 前端有BUG？

    不会前端，[vibe coding](https://zh.wikipedia.org/wiki/Vibe_coding)做的。

    有bug也不知道怎么改。

* 多用户？

    个人用不上，懒得写。

* 想要XX功能？

    ~~看心情~~ 欢迎贡献

* TUI？

    画饼。

* GUI?

    画饼。

* 数据导入？

    画饼。

## 许可证

本项目遵循 MIT 许可证。详情请参见 [LICENSE](LICENSE) 文件。
