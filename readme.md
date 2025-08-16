# Collectify

[![Go Report Card](https://goreportcard.com/badge/github.com/Jinvic/collectify)](https://goreportcard.com/report/github.com/Jinvic/collectify)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> [!WARNING]
> 此项目正在开发中，仍未完成基础功能，仅作同步Git仓库使用。

**Collectify** 是一个轻量级、易于使用的个人收藏管理系统，帮助你整理和管理你的书籍、音乐、影视等各类收藏品。无论是为了记录你的阅读进度，还是为了整理你喜爱的音乐专辑，Collectify 都能为你提供所需的一切工具。

## 快速开始

### 后端

1. **安装 Go**: 确保你已安装 Go 1.23 或更高版本。
2. **获取依赖**: 在项目根目录下运行 `go mod tidy`。
3. **运行后端**:
    * 开发模式: `go run main.go` (默认监听 :8080)。
    * 或构建后运行: `go build -o collectify .` 然后 `./collectify` (Linux/macOS) 或 `collectify.exe` (Windows)。

### 前端 (Web UI)

前端使用 React 构建，位于 `web` 目录。

1. **安装 Node.js**: 确保你已安装 Node.js (建议 LTS 版本)。
2. **安装 Pnpm**: 确保你已安装 Pnpm。
3. **安装前端依赖**: 在 `web` 目录下运行 `pnpm install`。
4. **开发**:
    * 启动前端开发服务器: 在 `web` 目录下运行 `pnpm start`。这将在 `http://localhost:3000` 启动一个热重载的开发服务器。它通过 `package.json` 中的 `proxy` 设置代理 API 请求到后端 (`http://localhost:8080`)。
    * 同时确保后端服务 (`go run main.go`) 正在运行。
5. **构建**:
    * 构建生产版本: 在 `web` 目录下运行 `pnpm run build`。这会将所有静态资源生成到 `web/build` 目录。
    * 构建后端时，它会自动检测 `web/build` 目录并提供这些静态文件，使得应用可以作为一个整体部署。

## 许可证

本项目遵循 MIT 许可证。详情请参见 [LICENSE](LICENSE) 文件。
