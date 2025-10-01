# inpayos - 统一支付系统

基于MVC架构设计的支付系统，整合BankPayOS系列服务功能。

## 功能特性

- 统一账户体系（商户、收银员、银行）
- 多币种资产管理
- 完整的资金流水记录
- 支持多种交易类型（代收、代付、退款、转账）
- 基于MVC的项目架构

## 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL 13+
- Redis 6+ (可选)

### 快速启动（推荐）

使用Docker Compose一键启动所有依赖服务：

```bash
# 启动PostgreSQL和Redis
docker-compose up -d

# 等待数据库就绪
sleep 10

# 安装Go依赖
go mod tidy

# 运行应用
go run main/main.go
```

### 手动配置数据库

如果不使用Docker，请手动安装PostgreSQL：

1. 安装PostgreSQL并创建数据库：
```sql
CREATE USER inpayos WITH PASSWORD 'password';
CREATE DATABASE inpayos OWNER inpayos;
CREATE DATABASE inpayos_dev OWNER inpayos;
GRANT ALL PRIVILEGES ON DATABASE inpayos TO inpayos;
GRANT ALL PRIVILEGES ON DATABASE inpayos_dev TO inpayos;

-- 创建扩展
\c inpayos;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

\c inpayos_dev;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
```

2. 修改配置文件 `config.yaml` 中的数据库连接信息（如果需要）

### 运行应用

```bash
go run main/main.go
```

### API测试

#### 创建账户
```bash
curl -X POST http://localhost:8080/api/v1/accounts \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_001",
    "user_type": "merchant",
    "currency": "USD"
  }'
```

#### 查询余额
```bash
curl "http://localhost:8080/api/v1/accounts/balance?user_id=user_001&user_type=merchant&currency=USD"
```

#### 更新余额
```bash
curl -X POST http://localhost:8080/api/v1/accounts/balance \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_001",
    "user_type": "merchant",
    "currency": "USD",
    "operation": "add",
    "amount": "100.00",
    "business_type": "deposit",
    "description": "Initial deposit"
  }'
```

## 项目结构

```
inpayos/
├── main/                     # 应用入口
├── internal/                 # 内部包
│   ├── config/              # 配置管理
│   ├── handlers/            # HTTP处理器
│   ├── models/              # 数据模型
│   ├── protocol/            # API协议
│   ├── services/            # 业务逻辑
│   └── utils/               # 工具函数
├── test/                    # 测试文件
├── docs/                    # 文档
├── config.yaml              # 配置文件
├── dev.yaml                 # 开发环境配置
└── go.mod                   # Go模块定义
```

## 核心模型

### Account 账户模型
- 支持多用户类型（merchant、cashier、bank）
- 多币种资产管理
- 支持可用余额、冻结余额、保证金、预留余额

### Transaction 交易模型
- 支持代收、代付、退款、转账
- 完整的交易状态管理
- 支持通知回调

### FundFlow 资金流水模型
- 记录所有资金变动
- 支持余额快照
- 便于对账和审计

## 开发指南

### 添加新的API接口

1. 在 `internal/protocol/` 中定义请求/响应结构
2. 在 `internal/services/` 中实现业务逻辑
3. 在 `internal/handlers/` 中添加HTTP处理器
4. 在 `main/main.go` 中注册路由

### 数据库迁移

项目使用GORM进行数据库操作，表结构通过模型定义自动创建。

## 许可证

MIT License