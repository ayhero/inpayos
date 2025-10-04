# 📊 inpayos 支付网关架构设计文档

> **版本**: v1.0  
> **更新时间**: 2025年10月4日  
> **文档类型**: 系统架构设计

## 📋 目录

- [1. 架构概述](#1-架构概述)
- [2. 系统架构图](#2-系统架构图)
- [3. 核心业务流程](#3-核心业务流程)
- [4. 架构分层说明](#4-架构分层说明)
- [5. 数据模型设计](#5-数据模型设计)
- [6. 安全认证体系](#6-安全认证体系)
- [7. 部署架构](#7-部署架构)
- [8. 扩展性设计](#8-扩展性设计)

## 1. 架构概述

inpayos 是一个**多语言支持的支付网关服务**，采用分层多服务架构设计。系统为商户提供统一的支付接口入口，支持银行直连和第三方支付渠道，其中 Cashier 作为核心支付渠道之一，支持多账户独立运营。

### 1.1 核心设计理念

- **统一接口**: 商户只通过 OpenAPI 接入，屏蔽渠道差异
- **渠道灵活**: 支持多种支付渠道，易于扩展
- **管理分层**: 不同角色管理各自数据范围
- **权限隔离**: 严格的数据和权限隔离机制
- **多语言**: 完整的国际化支持

## 2. 系统架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                          用户层 (User Layer)                      │
├─────────────────────────────────────────────────────────────────┤
│                      🏢 商户系统                                  │
│                   Merchant Systems                              │
│                  (电商/企业应用)                                  │
└─────────────────────────────┬───────────────────────────────────┘
                              │
                              │ [支付请求]
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                    OpenAPI - 支付网关主接口                       │
├─────────────────────────────────────────────────────────────────┤
│                      🌐 OpenAPI Gateway                          │
│                    统一支付接口入口                               │
│                   Port: 8080                                   │
│                   Auth: API Key (商户密钥)                      │
│                                                                 │
│  功能：渠道路由 | 收银台 | 交易管理 | 状态回调                    │
└─────┬───────────────┬───────────────────────────────────────┘
      │               │
      │ [路由分发]      │
      ▼               ▼
┌─────────┐    ┌─────────────┐
│💰Cashier│    │🏦 Bank      │
│Channel  │    │Channel      │
│         │    │             │
│Port:8081│    │Third-Party  │
│API Key  │    │API          │
└─────────┘    └─────────────┘
      │
      │ [Cashier渠道详细架构]
      ▼
┌─────────────────────────────────────────────────────────────────┐
│                   CashierAPI - 支付渠道实现                      │
├─────────────────────────────────────────────────────────────────┤
│  🏦 多Cashier账户支持 (ChannelCode = "Cashier")                  │
│                                                                 │
│  AccountID_001 │ AccountID_002 │ AccountID_003 │ AccountID_N   │
│  (团队A)       │ (团队B)       │ (团队C)       │ (团队...)      │
│                                                                 │
│  各自独立的：                                                    │
│  • 资金池管理   • 银行通道配置   • 风控策略   • 对账流程          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                       管理层 (Admin Layer)                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  🏪 MerchantAdmin    💼 CashierAdmin      👨‍💼 GlobalAdmin     │
│  商户数据管理         Cashier渠道管理       全局数据管控           │
│  Port: 8082         Port: 8083          Port: 8084           │
│  JWT + Merchant     JWT + Cashier       JWT + Admin          │
│                                                               │
│  管理范围：          管理范围：            管理范围：             │
│  ├ 自己的商户数据     ├ 自己的Cashier账户   ├ 所有商户数据        │
│  ├ 自己的交易记录     ├ 自己的资金流水      ├ 所有Cashier数据     │
│  ├ 自己的API配置     ├ 自己的银行通道      ├ 全局系统配置        │
│  ├ 自己的收银台      ├ 自己的风控参数      ├ 平台监控告警        │
│  └ 自己的渠道偏好     └ 自己的对账数据      └ 跨渠道数据统计      │
│                                                               │
└───────┬─────────────────┬─────────────────┬───────────────────┘
        │                 │                 │
        └─────────────────┼─────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│                   共享业务服务层 (Shared Services)                │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  💳 TransactionService  🛒 CheckoutService  👤 AccountService   │
│  交易处理服务            收银台服务           账户服务             │
│                                                                 │
│  ⚙️ ConfigService      📢 NotificationService                  │
│  配置服务               通知服务                                  │
│                                                                 │
│  🔐 AuthService        📊 ReportService                        │
│  认证服务               报表服务                                  │
│                                                                 │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│                     数据隔离层 (Data Layer)                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  💾 数据库分区 (Database Partitions)                            │
│                                                                 │
│  🏪 商户数据           💼 Cashier数据         🌐 全局数据        │
│  ├ merchant_*         ├ cashier_account_*   ├ system_config    │
│  ├ merchant_txn_*     ├ cashier_txn_*       ├ global_stats     │
│  ├ merchant_config_*  ├ cashier_balance_*   ├ audit_logs       │
│  └ api_keys_*         └ bank_channels_*     └ notifications    │
│                                                                 │
│  🚀 缓存分区 (Cache Partitions)                                │
│  ├ merchant_sessions  ├ cashier_sessions    ├ global_cache     │
│  └ merchant_configs   └ cashier_configs     └ system_status    │
│                                                                 │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│                   外部系统 (External Systems)                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  🏦 银行系统              🔗 第三方支付           📡 商户回调      │
│  Banking Systems       Third-Party APIs      Merchant Webhooks│
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 3. 核心业务流程

### 3.1 支付处理流程

```
🏢 商户 ──[API调用]──➤ 🌐 OpenAPI ──[渠道路由]──➤ 💰 CashierAPI(AccountID_X)
                                 │                         │
                                 ├[其他渠道]                ▼
                                 └─➤ 🏦 银行渠道       🏦 银行系统/第三方
                                                           │
🏢 商户 ◀──[Webhook]──── 📡 状态通知 ◀──[回调]─────────────┘
```

### 3.2 数据管理流程

```
👨‍💼 GlobalAdmin ──[全局管控]──➤ 🌐 所有数据
     │
     ├─[审核开通]─➤ 🏪 MerchantAdmin ──[管理]──➤ 📊 自己商户数据
     │
     └─[审核开通]─➤ 💼 CashierAdmin ──[管理]──➤ 💰 自己Cashier数据
```

## 4. 架构分层说明

### 4.1 用户层 (User Layer)

**商户系统 (Merchant Systems)**
- **角色**: 唯一的外部用户
- **接入方式**: 通过 OpenAPI 接入
- **典型用户**: 电商平台、企业应用、SaaS系统

### 4.2 网关层 (Gateway Layer)

**OpenAPI - 支付网关主接口**
- **端口**: 8080
- **认证**: API Key (商户密钥)
- **核心功能**:
  - 统一支付接口封装
  - 智能渠道路由
  - 收银台服务
  - 交易状态管理
  - 支付回调处理

### 4.3 渠道层 (Channel Layer)

**CashierAPI - Cashier支付渠道**
- **端口**: 8081
- **认证**: API Key (渠道间调用)
- **渠道标识**: ChannelCode = "Cashier"
- **特点**: 支持多 AccountID，每个账户独立运营

**第三方渠道**
- 银行直连渠道
- 其他第三方支付

### 4.4 管理层 (Admin Layer)

#### MerchantAdmin - 商户管理后台
- **端口**: 8082
- **认证**: JWT + Merchant 权限
- **管理范围**: 自己的商户数据
- **功能**:
  - 商户账户管理
  - API密钥管理  
  - 交易数据查询
  - 收银台配置
  - 渠道偏好设置

#### CashierAdmin - Cashier渠道管理
- **端口**: 8083
- **认证**: JWT + Cashier 权限
- **管理范围**: 自己的Cashier账户数据
- **功能**:
  - Cashier账户管理
  - 资金池管理
  - 交易监控对账
  - 风控参数配置
  - 银行通道管理

#### GlobalAdmin - 全局管理
- **端口**: 8084
- **认证**: JWT + Admin 权限
- **管理范围**: 全局数据和跨渠道统计
- **功能**:
  - 商户准入审核
  - Cashier渠道准入管理
  - 全局系统配置
  - 平台监控告警
  - 跨渠道数据统计

### 4.5 服务层 (Service Layer)

所有管理层共享的业务服务：

- **TransactionService**: 交易处理服务
- **CheckoutService**: 收银台服务
- **AccountService**: 账户服务
- **ConfigService**: 配置服务
- **NotificationService**: 通知服务
- **AuthService**: 认证服务
- **ReportService**: 报表服务

### 4.6 数据层 (Data Layer)

**数据隔离策略**:
- **商户数据**: merchant_*, merchant_txn_*, merchant_config_*
- **Cashier数据**: cashier_account_*, cashier_txn_*, cashier_balance_*
- **全局数据**: system_config, global_stats, audit_logs

**缓存分区**:
- **商户缓存**: merchant_sessions, merchant_configs
- **Cashier缓存**: cashier_sessions, cashier_configs
- **全局缓存**: global_cache, system_status

## 5. 数据模型设计

### 5.1 渠道账户模型

```go
type ChannelAccount struct {
    ID          int64            `json:"id"`
    MID         string           `json:"mid"`           // 商户ID
    ChannelCode string           `json:"channel_code"`  // 渠道代码
    AccountID   string           `json:"account_id"`    // 账户ID (唯一)
    Secret      string           `json:"secret"`        // 账户密钥
    Detail      protocol.MapData `json:"detail"`        // 账户详细信息
    Pkgs        []string         `json:"pkgs"`          // 支持的功能包
    Status      int              `json:"status"`        // 账户状态
    Settings    protocol.MapData `json:"settings"`      // 账户设置
    Groups      []string         `json:"groups"`        // 所属组
    CreatedAt   int64            `json:"created_at"`
    UpdatedAt   int64            `json:"updated_at"`
}
```

### 5.2 渠道组模型

```go
type ChannelGroup struct {
    ID        int64         `json:"id"`
    Code      string        `json:"code"`      // 组代码
    Name      string        `json:"name"`      // 组名称
    Status    int           `json:"status"`    // 组状态
    Setting   *GroupSetting `json:"setting"`   // 组设置
    Members   GroupMembers  `json:"members"`   // 组成员
    CreatedAt int64         `json:"created_at"`
    UpdatedAt int64         `json:"updated_at"`
}

type GroupSetting struct {
    Strategy  string `json:"strategy"`   // 路由策略
    Weight    string `json:"weight"`     // 权重配置
    RankType  string `json:"rank_type"`  // 排序类型
    TimeIndex string `json:"time_index"` // 时间索引
    DataIndex string `json:"data_index"` // 数据索引
    Timezone  string `json:"timezone"`   // 时区
}
```

## 6. 安全认证体系

### 6.1 认证方式

#### API Key 认证
- **应用场景**: OpenAPI ↔ 商户系统，OpenAPI ↔ CashierAPI
- **特点**: 无状态认证，适合系统间调用
- **实现**: `middleware.APIKeyAuth()`

#### JWT 认证
- **应用场景**: 管理后台用户认证
- **特点**: 有状态会话认证，支持用户登录
- **JWT Claims**:
```go
type JWTClaims struct {
    UserType   string `json:"user_type"` // merchant, admin, cashier
    MerchantID string `json:"merchant_id"` 
    Role       string `json:"role"`
}
```

### 6.2 权限控制

#### 三层权限体系
- **Merchant**: 商户权限，只能访问自己的数据
- **Cashier**: Cashier权限，只能访问自己账户的数据
- **Admin**: 管理员权限，可以访问全局数据

#### 权限中间件
- `middleware.MerchantPermissionMiddleware()`
- `middleware.CashierPermissionMiddleware()`
- `middleware.AdminPermissionMiddleware()`

### 6.3 多语言错误处理

#### 错误码体系
- **1000-1999**: 系统级错误
- **2000-2999**: 请求相关错误
- **3000-3999**: 认证相关错误
- **4000-4999**: 商户相关错误
- **5000-5999**: 交易相关错误
- **6000-6999**: 渠道相关错误
- **7000-7999**: Webhook相关错误
- **8000-8999**: 配置相关错误

#### 错误响应格式
```go
type Result struct {
    Code string `json:"code"`
    Msg  string `json:"msg"`
    Data any    `json:"data,omitempty"`
}
```

## 7. 部署架构

### 7.1 端口分配

| 服务 | 端口 | 说明 |
|------|------|------|
| OpenAPI | 8080 | 支付网关主接口 |
| CashierAPI | 8081 | Cashier支付渠道 |
| MerchantAdmin | 8082 | 商户管理后台 |
| CashierAdmin | 8083 | Cashier管理后台 |
| GlobalAdmin | 8084 | 全局管理后台 |

### 7.2 部署选项

#### 单体部署
- 所有服务运行在同一进程
- 通过不同端口提供服务
- 适合中小规模部署

#### 微服务部署
- 每个服务独立部署
- 通过服务发现和负载均衡
- 适合大规模分布式部署

## 8. 扩展性设计

### 8.1 渠道扩展

新增支付渠道只需：
1. 实现标准的渠道接口
2. 在渠道配置中注册
3. 配置路由规则

### 8.2 Cashier账户扩展

新增Cashier账户：
1. 创建新的AccountID
2. 配置独立的资金池和银行通道
3. 分配给Cashier团队管理

### 8.3 管理功能扩展

- 支持新的管理角色
- 扩展权限控制粒度
- 增加新的管理功能模块

### 8.4 国际化扩展

- 支持新语言的错误消息
- 多时区支持
- 多币种支持

## 📝 总结

inpayos 采用的分层多服务架构设计具有以下优势：

1. **统一接口**: 商户只需对接 OpenAPI，简化集成
2. **渠道灵活**: 支持多种支付渠道，易于扩展
3. **权限清晰**: 不同角色管理各自数据，权限隔离
4. **扩展性强**: 可单体部署或拆分为微服务
5. **国际化**: 完整的多语言支持
6. **安全可控**: 多层认证和权限控制机制

这种架构设计既满足了支付网关的复杂业务需求，又保持了良好的可扩展性和可维护性，为构建企业级支付平台提供了坚实的基础。