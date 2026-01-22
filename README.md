# SeedGo 后端脚手架

SeedGo 是一个基于 Go 语言的企业级现代后端开发脚手架，采用分层架构设计，集成了用户管理、角色权限控制（RBAC）、多租户支持等核心功能。它旨在帮助开发者快速搭建高性能、可扩展的
Web 应用程序。

## 🛠 技术栈

- **编程语言**: Go (1.25+)
- **Web 框架**: [Gin](https://github.com/gin-gonic/gin)
- **ORM 框架**: [GORM](https://gorm.io/) (MySQL)
- **配置管理**: [Viper](https://github.com/spf13/viper)
- **认证授权**: JWT (JSON Web Token)
- **容器化**: Docker & Docker Compose

## ✨ 核心特性

- **模块化设计**: 采用 `internal` 目录隔离私有代码，结构清晰。
- **RBAC 权限模型**: 内置完善的用户、角色、权限管理逻辑。
- **多租户架构**: 支持多租户数据隔离与管理。
- **标准化响应**: 统一的 API 响应格式与错误处理。
- **环境隔离**: 支持本地开发、测试、生产等多环境配置。

## 📸 系统截图

<p align="center">
  <img src="docs/images/login.png" alt="Login" width="48%" />
  <img src="docs/images/dashboard.png" alt="Dashboard" width="48%" />
</p>
<p align="center">
  <img src="docs/images/table.png" alt="Table" width="32%" />
  <img src="docs/images/create.png" alt="Create" width="32%" />
  <img src="docs/images/perms.png" alt="Perms" width="32%" />
</p>

## 📂 目录结构

```text
├── cmd/                # 命令行入口 (migrate, test 等)
├── config/             # 配置文件 (local.yaml 等)
├── docker/             # Docker 环境配置
│   ├── docker-compose.yml # 依赖服务编排
│   └── mysql/          # MySQL 配置与初始化脚本
├── internal/           # 私有业务代码
│   ├── app/            # 应用层 (Router 注册)
│   ├── biz/            # 业务逻辑层
│   ├── middleware/     # Gin 中间件 (Auth, CORS)
│   ├── model/          # 数据库模型定义
│   ├── shared/         # 共享逻辑 (BaseController 等)
│   └── system/         # 系统模块 (用户, 角色, 权限, 租户)
├── pkg/                # 公共工具包 (DB, JWT, Utils)
├── main.go             # 项目主入口
└── Makefile            # 构建脚本
```

## 🚀 快速开始

### 1. 环境准备

确保你的本地环境已安装以下工具：

- [Go](https://go.dev/dl/) (1.25 或更高版本)
- [Docker](https://www.docker.com/) & Docker Compose
- [Make](https://www.gnu.org/software/make/) (可选，用于简化构建命令)

### 2. 启动依赖环境 (Docker)

项目依赖 MySQL 数据库。我们使用 Docker Compose 快速启动运行环境。

```bash
cd docker
docker-compose up -d
```

该命令将启动一个 MySQL 8.0 容器，并在端口 `3306` 监听。

- **数据库名**: `smart`
- **Root 密码**: `123456`
- **数据持久化**: 数据保存在 `docker/mysql/data` 目录下。

### 3. 配置项目

检查 `config/local.yaml` 文件，确保数据库连接配置与 Docker 环境一致：

```yaml
database:
  dsn: "root:123456@tcp(127.0.0.1:3306)/smart?charset=utf8mb4&parseTime=True&loc=Local"
```

### 4. 运行项目

#### 使用 Go 直接运行

```bash
go mod tidy
go run main.go
```

#### 使用 Makefile 编译运行

```bash
# 编译二进制文件
make build

# 运行二进制文件
./main
```

服务启动后，默认监听在 `http://localhost:3000`。

## 🐳 构建应用 Docker 镜像 (可选)

如果你希望将应用程序本身也打包为 Docker 镜像，可以在项目根目录创建一个 `Dockerfile`：

```dockerfile
# 第一阶段：构建
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 设置 Go 代理以加速下载
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# 第二阶段：运行
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/config ./config

# 暴露端口
EXPOSE 3000

CMD ["./server"]
```

构建并运行镜像：

```bash
# 构建镜像
docker build -t seedgo-app .

# 运行容器 (注意：如果 MySQL 也在 Docker 中，需要让它们在同一网络下，或使用宿主机 IP)
docker run -p 3000:3000 seedgo-app
```

## 📝 常用命令

| 命令                     | 说明                 |
|------------------------|--------------------|
| `make build`           | 编译项目生成二进制文件 `main` |
| `go run main.go`       | 本地开发模式运行           |
| `docker-compose up -d` | 启动数据库服务            |
| `docker-compose down`  | 停止并移除数据库服务         |

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request 来改进本项目。

## 📄 许可证

[BUSINESS SOURCE LICENSE](LICENSE)

个人和非盈利组织以研究学习试用免费，用于商业用途目的和其余商业组织需购买授权。
