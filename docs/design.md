# inpayos 项目设计文档

## 1. 项目结构

```
inpayos/
├── .github/
├── .vscode/
├── bin/
├── build/
├── docs/
│   └── design.md
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── i18n/
│   ├── locales/
│   ├── log/
│   ├── middleware/
│   ├── models/
│   ├── protocol/
│   ├── queue/
│   ├── services/
│   ├── task/
│   └── utils/
├── logs/
├── main/
├── nginx/
├── scripts/
├── test/
├── .gitignore
├── config.yaml
├── dev.yaml
├── Dockerfile
├── go.mod
├── go.sum
├── local.yaml
├── Makefile
└── prod.yaml
```

## 2. 核心模型设计

- 账户模型支持多用户类型（UserID + UserType + Currency 唯一）
- 资产模型支持多资金属性（可用、冻结、保证金、预留等）
- 资金流水模型记录所有变动
- 交易模型支持多种交易类型

## 3. 业务流程

- 账户创建、查询、余额变动
- 资金流水自动生成
- 交易处理（收款、付款、退款等）
- 结算与风控

## 4. 实施步骤

1. 初始化 inpayos 项目目录结构
2. 创建 docs/design.md 设计文档
3. 实现 internal/models 下的核心数据模型
4. 实现 internal/handlers 下的 API 处理器
5. 实现 internal/services 下的业务逻辑
6. 配置 main 入口文件和基础路由
7. 编写测试用例（test/）
8. 完善配置文件和 Dockerfile
9. 持续迭代和优化

---

详细设计与代码实现请见各模块说明。