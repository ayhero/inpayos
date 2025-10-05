# 📊 inpayos 支付网关架构设计文档

> **版本**: v1.1  
> **更新时间**: 2025年10月5日  
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

inpayos 是一个**多语言支持的支付网关服务**，采用分层多服务架构设计。系统支持**商户(Merchant)**和**收银团队(CashierTeam)**两类平级用户角色，通过统一的 OpenAPI 提供支付服务，支持银行直连和第三方支付渠道。

### 1.1 核心设计理念

- **平级角色**: Merchant 和 CashierTeam 作为两类平级的业务角色，各自拥有独立的数据和权限范围
- **统一账户**: 基于统一账户模型，通过 UserType 区分不同角色类型
- **统一接口**: 所有用户通过 OpenAPI 接入，屏蔽底层渠道差异
- **渠道灵活**: 支持多种支付渠道，易于扩展
- **管理分层**: 三层管理体系 - MerchantAdmin、CashierAdmin、GlobalAdmin
- **权限隔离**: 严格的数据和权限隔离机制
- **多语言**: 完整的国际化支持

## 2. 系统架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                          用户层 (User Layer)                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│                    🏢 商户系统 (Merchants)                       │
│                     电商/企业应用/SaaS系统                        │
│                     UserType: "merchant"                       │
│                                                                 │
│                  唯一的外部用户，通过OpenAPI接入                   │
│                                                                 │
└─────────────────────────────┬───────────────────────────────────┘
                              │
                              │ [支付请求/收银服务]
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                    OpenAPI - 商户支付网关接口                     │
├─────────────────────────────────────────────────────────────────┤
│                      🌐 OpenAPI Gateway                          │
│                   专为商户提供的统一接口                           │
│                   Port: 8080                                   │
│                   Auth: API Key (商户密钥)                      │
│                                                                 │
│  功能：渠道路由 | 收银台 | 交易管理 | 状态回调                    │
│  服务对象：商户系统 (唯一用户)                                   │
└─────┬───────────────┬───────────────────────────────────────┘
      │               │
      │ [智能路由分发]  │
      ▼               ▼
┌─────────┐    ┌─────────────────┐
│💰Cashier│    │🔗 ThirdParty   │
│Channel  │    │APIs         │
│         │    │             │
│Port:8081│    │(银行+第三方) │
│API Key  │    │APIs         │
└─────────┘    └─────────────────┘
      │
      │ [CashierAPI渠道处理]
      ▼
┌─────────────────────────────────────────────────────────────────┐
│                   CashierAPI - 收银渠道服务实现                   │
├─────────────────────────────────────────────────────────────────┤
│                   CashierAPI 调用层 (CashierTeam)                    │
│                                                                 │
│  CashierTeam_A │ CashierTeam_B │ CashierTeam_C │ CashierTeam_N│
│  (收银团队A)    │ (收银团队B)    │ (收银团队C)    │ (更多团队)    │
│                                                                 │
│  CashierAPI的具体实现层，各团队独立运营：                         │
│  • 资金池管理   • 银行卡账户     • 风控策略   • 对账流程          │
│  • 收银员配置   • 业务规则      • 费率设置   • 运营数据          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                     分层管理体系 (Admin Layer)                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  🏪 MerchantAdmin              👨‍💼 GlobalAdmin                 │
│  商户自管理后台                  全局运营管理                      │
│  Port: 8082                    Port: 8084                     │
│  JWT + Merchant权限             JWT + Admin权限                │
│                                                               │
│  管理范围：                      管理范围：                      │
│  ├ 自己的商户信息                ├ 所有商户数据                  │
│  ├ 自己的交易记录                ├ 所有渠道数据                  │
│  ├ 自己的API配置                ├ CashierTeam运营管理           │
│  ├ 自己的收银台                 ├ 全局系统配置                  │
│  └ 自己的渠道偏好                └ 平台监控告警                  │
│                                                               │
│  💡 CashierAdmin作为GlobalAdmin的子功能存在                    │
│     用于管理收银渠道和团队运营                                  │
│                                                               │
└───────┬─────────────────────────────────────┬───────────────────┘
        │                 │                 │
        └─────────────────┼─────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│                   核心共享服务层 (Core Shared Services)          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  💳 TransactionService  � AccountService   � FlowService     │
│  交易服务               账户服务             流水服务             │
│                                                                 │
│  💰 DepositService     � WithdrawService   ⚖️ SettlementService│
│  充值服务              提现服务             结算规则服务           │
│                                                                 │
│  ⏰ TaskService        📢 MessageService                       │
│  定时任务服务           消息服务                                   │
│                                                                 │
│  📦 各系统专属业务模块在各自的Admin层实现                         │
│                                                                 │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│                 统一账户数据层 (Unified Data Layer)               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  💾 统一账户表 (t_accounts)                                     │
│  UserID + UserType + Currency 唯一索引                         │
│  ├ UserType: "merchant"      - 商户账户                        │
│  └ UserType: "admin"         - 管理员账户                      │
│                                                                 │
│  💡 说明：                                                      │
│  • Cashier作为支付渠道，不是用户角色                            │
│  • CashierTeam是渠道运营方，通过GlobalAdmin管理                │
│                                                                 │
│  🏪 商户数据表           💼 渠道运营表         🌐 全局数据表       │
│  ├ t_merchants          ├ t_cashiers         ├ t_admins        │
│  ├ t_merchant_admins    ├ t_cashier_admins   ├ system_config   │
│  ├ merchant_transactions├ (渠道交易数据)      ├ audit_logs      │
│  └ merchant_configs     └ cashier_configs    └ global_stats    │
│                                                                 │
│  🚀 缓存分区 (Cache Partitions)                                │
│  ├ merchant_sessions      ├ cashier_sessions     ├ global_cache │
│  └ merchant_configs       └ cashier_configs      └ system_status│
│                                                                 │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│                   外部系统 (External Systems)                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  🏦 银行系统              🔗 第三方支付           📡 用户回调      │
│  Banking Systems       Third-Party APIs      User Webhooks    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 3. 核心业务流程

### 3.1 商户支付处理流程

```
🏢 商户 (唯一用户)
   │
   │ [API调用]
   ▼
🌐 OpenAPI Gateway (商户专用入口)
   │
   │ [智能路由选择]
   ▼
┌─────────────────┐
│   渠道选择器     │
│ ChannelRouter   │
└─────┬───────────┘
      │
┌─────▼─────┬─────────────┐
│           │             │
▼           ▼             ▼
💰CashierAPI    🔗ThirdPartyAPI
(收银渠道)      (银行+第三方API)
│                     │
▼                     │
┌─────────────┐       │
│💼 CashierTeam│       │
│  调用层      │       │
│ (具体实现)    │       │
└─────┬───────┘       │
      │               │
      └───────────────▼
      📡 处理结果回调
            │
            ▼
      🏢 商户回调通知
```

### 3.2 管理流程

```
👨‍💼 GlobalAdmin (全局运营管控)
     │
     │ [系统级管理]
     ├─────────────────────────────────┐
     │                                 │
     ▼                                 ▼
🏪 MerchantAdmin                � CashierTeam运营管理
  (商户自管理)                    (渠道运营管理)
     │                                 │
     │ [UserType: merchant]            │ [渠道管理功能]
     ▼                                 ▼
📊 商户数据范围                   💰 渠道运营数据范围
├ t_accounts (merchant)          ├ t_cashiers (渠道数据)
├ 商户交易数据                    ├ 渠道交易统计
├ 商户配置                       ├ 渠道配置
├ API密钥                        ├ 收银员管理
└ 收银台设置                     └ 银行卡管理

架构特点：
• 商户是唯一的外部用户，通过OpenAPI接入
• CashierTeam是渠道运营方，通过GlobalAdmin管理
• 所有支付服务最终服务于商户
```

### 3.3 商户支付业务流程

```
🏢 商户支付请求 (充值/提现/支付)
        │
        │ [通过OpenAPI]
        ▼
    📦 TransactionService
    (统一交易抽象层)
             │
             ▼
       ┌─────────────┐
       │渠道路由选择  │
       │ChannelCode  │
       └─────┬───────┘
             │
    ┌──────▼──────┬──────────────┐
    │             │              │
    ▼             ▼              ▼
�CashierAPI   🏦BankAPI    🔗ThirdPartyAPI
(收银渠道)     (银行渠道)    (第三方渠道)
    │             │              │
    ▼             ▼              ▼
📋 业务表记录 (SType="merchant", Sid=商户ID)
             │
             ▼
    💳 AccountService
    (统一账户服务)
             │
             ▼
       t_accounts表
    (UserType="merchant")
             │
             ▼
    🔄 商户资金变动处理
             │
             ▼
    📡 商户回调通知
```

### 3.4 商户账户管理流程

```
商户注册请求
    │
    ▼
┌───────────────┐
│ AccountService │ ──➤ 创建商户账户
└───────┬───────┘
        │
        ▼
    t_accounts表
    (UserType="merchant")
        │
        ▼
    商户账户创建完成
        │
        ▼
    通过OpenAPI提供服务
```

### 3.5 渠道运营管理流程

```
CashierTeam管理
    │
    ▼
┌─────────────────┐
│ GlobalAdmin     │ ──➤ 渠道运营管理
│ CashierAdmin功能 │
└─────────┬───────┘
          │
          ▼
      t_cashiers表
      (渠道运营数据)
          │
          ▼
    CashierAPI渠道服务
    (ChannelCode="Cashier")
```

### 3.5 业务表统一设计流程

```
📋 业务表设计原则
├── 统一表结构 (Deposit, Withdraw)
├── Sid字段 (服务主体ID)
├── SType字段 (服务类型标识)
└── AccountID字段 (关联统一账户)

🎯 SType类型定义:
├── "merchant" - 商户业务 (唯一的用户角色)
└── "admin" - 管理员业务 (预留)

💡 说明：
• 只有商户是真正的用户，拥有账户和业务数据
• CashierTeam是渠道运营方，通过ChannelCode区分
• 所有业务最终服务于商户

🔄 数据处理流程:
1. 业务请求 → 识别SType → 路由到对应Service
2. Service处理 → 写入统一业务表 → 更新Account
3. 状态变更 → 触发回调 → 通知业务方
```

## 4. 架构分层说明

### 4.1 用户层 (User Layer)

**商户系统 (Merchant Systems)**
- **角色**: 唯一的外部用户角色
- **UserType**: "merchant"
- **接入方式**: 通过 OpenAPI 接入
- **典型用户**: 电商平台、企业应用、SaaS系统
- **主要需求**: 支付接入、交易管理、资金结算
- **服务特点**: 系统的所有功能都围绕商户需求设计

### 4.2 网关层 (Gateway Layer)

**OpenAPI - 商户支付网关接口**
- **端口**: 8080
- **认证**: API Key (商户密钥)
- **服务对象**: 仅服务于商户系统
- **核心功能**:
  - 为商户提供统一支付接口
  - 智能渠道路由（包括CashierAPI渠道）
  - 商户收银台服务
  - 商户交易状态管理
  - 商户支付回调处理

### 4.3 渠道层 (Channel Layer)

**CashierAPI - 收银渠道服务**
- **端口**: 8081
- **认证**: API Key (内部渠道调用)
- **渠道标识**: ChannelCode = "Cashier"
- **角色定位**: 被OpenAPI调用的支付渠道之一
- **调用方式**: 
  - 商户 → OpenAPI → 路由选择 → CashierAPI → CashierTeam调用层
  - 不直接面向商户，作为内部渠道存在
- **调用层架构**: CashierTeam在CashierAPI的调用层处理具体业务
- **运营管理**: 多个CashierTeam独立运营，通过GlobalAdmin管理

**第三方API渠道**
- 银行直连 API (不再单独作为渠道)
- 其他第三方支付 API
- 统一作为ThirdPartyAPI渠道处理

### 4.4 管理层 (Admin Layer)

#### MerchantAdmin - 商户自管理后台
- **端口**: 8082
- **认证**: JWT + Merchant 权限
- **数据范围**: UserType = "merchant" 的数据
- **管理范围**: 自己的商户数据
- **功能**:
  - 商户账户管理 (自己的Account记录)
  - API密钥管理  
  - 交易数据查询 (自己的交易)
  - 收银台配置
  - 渠道偏好设置
  - 商户信息维护

#### GlobalAdmin - 全局运营管理
- **端口**: 8084
- **认证**: JWT + Admin 权限
- **数据范围**: 全局数据和运营统计
- **管理范围**: 整个系统的运营管理
- **功能**:
  - 商户准入审核和管理
  - **CashierTeam渠道运营管理**:
    - 收银团队准入审核
    - 收银员管理 (Cashier表)
    - 渠道资金池管理
    - 银行卡管理
    - 渠道交易监控对账
    - 风控参数配置
  - 全局系统配置
  - 平台监控告警
  - 跨渠道数据统计

### 4.5 服务层 (Service Layer)

#### 4.5.1 核心共享服务层

```
📦 核心业务服务层 (Core Business Services)
├── 💳 TransactionService (交易服务)
│   ├── 统一交易处理抽象层
│   ├── 交易状态管理
│   └── 交易路由分发
│
├── 👤 AccountService (账户服务)
│   ├── 统一账户管理
│   ├── 余额操作
│   └── 账户状态控制
│
├── 📊 FlowService (流水服务)
│   ├── 资金流水记录
│   ├── 流水查询统计
│   └── 流水对账处理
│
├── 💰 DepositService (充值服务)
│   ├── 充值业务处理
│   ├── 跨角色充值支持
│   └── 充值状态管理
│
├── 💸 WithdrawService (提现服务)
│   ├── 提现业务处理
│   ├── 跨角色提现支持
│   └── 提现审核管理
│
├── ⚖️ SettlementService (结算规则服务)
│   ├── 结算规则配置
│   ├── 结算周期管理
│   └── 结算费率计算
│
├── ⏰ TaskService (定时任务服务)
│   ├── 定时任务调度
│   ├── 任务状态监控
│   └── 任务执行记录
│
└── 📢 MessageService (消息服务)
    ├── 系统消息通知
    ├── 回调消息处理
    └── 消息队列管理

📦 各系统专属业务模块
├── 🏪 MerchantAdmin 专属模块
│   ├── 商户注册认证
│   ├── API密钥管理
│   ├── 收银台配置
│   ├── 渠道偏好设置
│   └── 商户报表统计
│
├── � CashierAdmin 专属模块
│   ├── 收银员管理
│   ├── 银行卡管理
│   ├── 风控参数配置
│   ├── 团队权限管理
│   └── 团队运营数据
│
└── 👨‍� GlobalAdmin 专属模块
    ├── 系统配置管理
    ├── 权限体系管理
    ├── 平台监控告警
    ├── 审计日志管理
    └── 全局数据统计
```

#### 4.5.2 核心服务设计原则

**共享服务层特点**：
- **业务无关性**: 不包含特定角色的业务逻辑
- **高度复用**: 所有角色和系统都可以调用
- **统一接口**: 提供标准化的服务接口
- **数据中性**: 通过SType等字段区分不同业务主体

**核心服务说明**：
- **TransactionService**: 统一交易抽象层，管理所有交易的生命周期
- **AccountService**: 统一账户服务，处理跨角色的账户操作
- **FlowService**: 资金流水服务，记录所有资金变动
- **DepositService**: 充值服务，支持商户和收银团队充值
- **WithdrawService**: 提现服务，支持商户和收银团队提现
- **SettlementService**: 结算规则服务，处理结算逻辑和费率计算
- **TaskService**: 定时任务服务，处理系统级定时任务
- **MessageService**: 消息服务，处理系统通知和回调

#### 4.5.3 服务调用模式

```go
// 示例：核心服务统一接口设计
type DepositService interface {
    CreateDeposit(stype, sid string, req *CreateDepositRequest) (*Deposit, error)
    QueryDeposit(stype, sid string, req *QueryDepositRequest) (*Deposit, error)
    CancelDeposit(stype, sid string, req *CancelDepositRequest) (*Deposit, error)
}

type WithdrawService interface {
    CreateWithdraw(stype, sid string, req *CreateWithdrawRequest) (*Withdraw, error)
    QueryWithdraw(stype, sid string, req *QueryWithdrawRequest) (*Withdraw, error)
    CancelWithdraw(stype, sid string, req *CancelWithdrawRequest) (*Withdraw, error)
}

type AccountService interface {
    GetAccount(userType, userID, currency string) (*Account, error)
    UpdateBalance(userType, userID, currency string, amount decimal.Decimal) error
    FreezeBalance(userType, userID, currency string, amount decimal.Decimal) error
}

// 各系统调用示例
// 商户系统调用
depositService.CreateDeposit("merchant", merchantID, req)

// 收银团队系统调用
depositService.CreateDeposit("cashier_team", teamID, req)

// 管理后台调用
accountService.GetAccount("merchant", merchantID, "USD")

// 说明：
// - Cashier是支付渠道，通过ChannelCode="Cashier"标识
// - CashierTeam在CashierAPI的调用层处理具体业务
// - 没有独立的银行渠道，银行API归入ThirdPartyAPI
```

#### 4.5.4 业务模块分离原则

**核心服务** (共享)：
- 只处理数据操作和业务规则
- 不包含UI逻辑和特定角色权限
- 通过参数区分不同业务主体

**专属业务模块** (各系统独有)：
- 包含特定角色的业务逻辑
- 处理权限验证和数据过滤
- 提供角色特定的UI和API接口
- 调用核心服务完成数据操作

### 4.6 数据层 (Data Layer)

**统一账户数据策略**:
- **统一账户表**: t_accounts (UserID + UserType + Currency 唯一索引)
  - UserType = "merchant": 商户账户数据
  - UserType = "cashier_team": 收银团队账户数据
  - UserType = "admin": 管理员账户数据

**角色相关数据**:
- **商户数据**: t_merchants, t_merchant_admins, merchant_transactions, merchant_configs
- **收银团队数据**: t_cashiers, t_cashier_admins, cashier_transactions, cashier_configs
- **全局数据**: t_admins, system_config, global_stats, audit_logs

**缓存分区**:
- **商户缓存**: merchant_sessions, merchant_configs
- **收银团队缓存**: cashier_sessions, cashier_configs
- **全局缓存**: global_cache, system_status

## 5. 数据模型设计

### 5.1 统一账户模型

```go
// Account 统一账户表 - 支持商户和收银团队两类平级角色
type Account struct {
    ID        uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
    AccountID string `json:"account_id" gorm:"uniqueIndex"`          // 账户唯一标识
    Salt      string `json:"salt"`                                   // 加密salt
    *AccountValues
    CreatedAt int64 `json:"created_at" gorm:"autoCreateTime:milli"`
}

type AccountValues struct {
    UserID       *string `json:"user_id"`   // 用户ID（商户ID或收银团队ID）
    UserType     *string `json:"user_type"` // 用户类型：merchant, cashier_team, admin
    Ccy          *string `json:"ccy"`       // 币种
    Asset        *Asset  `json:"asset"`     // 资产信息
    Status       *int    `json:"status"`    // 账户状态
    Version      *int64  `json:"version"`   // 版本号（乐观锁）
    LastActiveAt *int64  `json:"last_active_at"` // 最后活跃时间
    UpdatedAt    int64   `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

// 索引：UserID + UserType + Currency 唯一
```

### 5.2 收银团队模型

```go
// Cashier 收银员/收银团队表
type Cashier struct {
    ID        uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
    CashierID string `json:"cashier_id" gorm:"uniqueIndex"`      // 收银员唯一标识
    AccountID string `json:"account_id" gorm:"index"`           // 关联的账户ID
    Salt      string `json:"salt"`
    *CashierValues
    CreatedAt int64 `json:"created_at" gorm:"autoCreateTime:milli"`
    UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

type CashierValues struct {
    // 基础信息
    Type        *string `json:"type"`         // private(私户), corporate(公户)
    BankCode    *string `json:"bank_code"`    // 银行代码
    BankName    *string `json:"bank_name"`    // 银行名称
    CardNumber  *string `json:"card_number"`  // 卡号
    HolderName  *string `json:"holder_name"`  // 持卡人姓名
    HolderPhone *string `json:"holder_phone"` // 持卡人手机
    HolderEmail *string `json:"holder_email"` // 持卡人邮箱
    
    // 地域信息
    Country     *string `json:"country"`      // 国家
    CountryCode *string `json:"country_code"` // 国家代码
    Province    *string `json:"province"`     // 省/州
    City        *string `json:"city"`         // 城市
    
    // 业务配置
    Ccy          *string           `json:"ccy"`           // 币种
    PayinStatus  *string           `json:"payin_status"`  // 收款状态
    PayinConfig  *protocol.MapData `json:"payin_config"`  // 收款配置
    PayoutStatus *string           `json:"payout_status"` // 付款状态
    PayoutConfig *protocol.MapData `json:"payout_config"` // 付款配置
    Status       *string           `json:"status"`        // 状态
    
    // 其他信息
    ExpireAt *int64  `json:"expire_at"` // 过期时间
    Logo     *string `json:"logo"`      // 头像/标志
    Remark   *string `json:"remark"`    // 备注
}
```

### 5.3 收银团队管理员模型

```go
// CashierAdmin 收银团队管理员表
type CashierAdmin struct {
    ID      int64  `json:"id" gorm:"primaryKey;autoIncrement"`
    AdminID string `json:"admin_id" gorm:"uniqueIndex"`        // 管理员唯一标识
    Salt    string `json:"salt"`
    *CashierAdminValues
    CreatedAt int64 `json:"created_at" gorm:"autoCreateTime:milli"`
}

type CashierAdminValues struct {
    // 基本信息
    Username *string `json:"username" gorm:"uniqueIndex"` // 用户名
    Email    *string `json:"email" gorm:"uniqueIndex"`    // 邮箱
    Phone    *string `json:"phone"`                       // 手机号
    
    // 个人信息
    FirstName *string `json:"first_name"` // 名
    LastName  *string `json:"last_name"`  // 姓
    FullName  *string `json:"full_name"`  // 全名
    Avatar    *string `json:"avatar"`     // 头像
    
    // 认证信息
    PasswordHash     *string `json:"-"`                    // 密码哈希
    PasswordSalt     *string `json:"-"`                    // 密码salt
    EmailVerified    *bool   `json:"email_verified"`       // 邮箱验证状态
    PhoneVerified    *bool   `json:"phone_verified"`       // 手机验证状态
    TwoFactorEnabled *bool   `json:"two_factor_enabled"`   // 2FA启用状态
    TwoFactorSecret  *string `json:"-"`                    // 2FA密钥
    
    // 角色和权限
    Role        *string `json:"role"`        // super_admin, admin, moderator, support, analyst
    Permissions *string `json:"permissions"` // JSON数组存储权限列表
    Department  *string `json:"department"`  // 部门
    JobTitle    *string `json:"job_title"`   // 职位
    
    // 状态管理
    Status       *string `json:"status"`        // active, inactive, suspended, locked
    ActiveStatus *string `json:"active_status"` // online, offline, busy
    
    // 登录相关
    LastLoginAt    *int64  `json:"last_login_at"`    // 最后登录时间
    LastLoginIP    *string `json:"last_login_ip"`    // 最后登录IP
    LoginCount     *int    `json:"login_count"`      // 登录次数
    FailedAttempts *int    `json:"failed_attempts"`  // 失败尝试次数
    LastFailedAt   *int64  `json:"last_failed_at"`   // 最后失败时间
    LockedUntil    *int64  `json:"locked_until"`     // 锁定截止时间
    
    // 会话管理
    CurrentSessionID      *string `json:"current_session_id"`        // 当前会话ID
    SessionCount          *int    `json:"session_count"`             // 会话数量
    MaxConcurrentSessions *int    `json:"max_concurrent_sessions"`   // 最大并发会话数
    
    // 其他字段...
    UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
```

### 5.4 全局管理员模型

```go
// Admin 全局管理员表
type Admin struct {
    ID      int64  `json:"id" gorm:"primaryKey;autoIncrement"`
    AdminID string `json:"admin_id" gorm:"uniqueIndex"`
    Salt    string `json:"salt"`
    *AdminValues
    CreatedAt int64 `json:"created_at" gorm:"autoCreateTime:milli"`
}

type AdminValues struct {
    Username     *string `json:"username" gorm:"uniqueIndex"`     // 用户名
    Email        *string `json:"email" gorm:"uniqueIndex"`        // 邮箱
    Role         *string `json:"role"`                           // 角色
    Status       *string `json:"status"`                         // 状态
    ActiveStatus *string `json:"active_status"`                  // 活跃状态
    UpdatedAt    int64   `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
```

### 5.5 统一业务交易模型

#### 5.5.1 通用交易抽象层

```go
// Transaction 通用交易记录表（作为所有业务交易的抽象层）
type Transaction struct {
    ID        uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
    Tid       string `json:"tid" gorm:"index"`                    // 交易ID
    CashierID string `json:"cashier_id" gorm:"index"`             // 收银员ID
    Mid       string `json:"mid" gorm:"index"`                    // 商户ID
    UserID    string `json:"user_id" gorm:"index"`                // 用户ID
    TrxID     string `json:"transaction_id" gorm:"uniqueIndex"`   // 交易唯一标识
    ReqID     string `json:"req_id" gorm:"index"`                 // 请求ID
    TrxType   string `json:"trx_type" gorm:"index"`               // 交易类型
    *TransactionValues
    CreatedAt int64 `json:"created_at" gorm:"autoCreateTime:milli"`
}

// 支持的交易类型：receipt, payment, refund, transfer, deposit, withdraw
```

#### 5.5.2 统一充值模型

```go
// Deposit 充值记录表 - 支持商户和收银团队统一充值
type Deposit struct {
    ID        uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
    TrxID     string `json:"trx_id" gorm:"uniqueIndex"`           // 交易唯一标识
    Sid       string `json:"sid" gorm:"index"`                    // 服务主体ID (商户ID或收银团队ID或收银员ID)
    SType     string `json:"s_type" gorm:"index"`                 // 服务类型：merchant, cashier_team, cashier
    AccountID string `json:"account_id" gorm:"index"`             // 关联账户ID
    *DepositValues
    CreatedAt int64 `json:"created_at" gorm:"autoCreateTime:milli"`
    UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

type DepositValues struct {
    Status      *string          `json:"status"`       // pending, processing, success, failed
    Ccy         *string          `json:"ccy"`          // 币种
    Amount      *decimal.Decimal `json:"amount"`       // 充值金额
    Fee         *decimal.Decimal `json:"fee"`          // 手续费
    ChannelCode *string          `json:"channel_code"` // 渠道代码
    NotifyURL   *string          `json:"notify_url"`   // 回调地址
    Country     *string          `json:"country"`      // 国家
    // 状态时间戳
    CanceledAt  *int64 `json:"canceled_at"`  // 取消时间
    CompletedAt *int64 `json:"completed_at"` // 完成时间
    ExpiredAt   *int64 `json:"expired_at"`   // 过期时间
    ConfirmedAt *int64 `json:"confirmed_at"` // 确认时间
}

// 通过 Sid + SType 区分业务主体:
// - SType="merchant", Sid=商户ID: 商户充值
// - SType="cashier_team", Sid=团队ID: 收银团队充值
```

#### 5.5.3 统一提现模型

```go
// Withdraw 提现记录表 - 支持商户和收银团队统一提现
type Withdraw struct {
    ID        uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
    Sid       string `json:"sid" gorm:"index"`                    // 服务主体ID (商户ID或收银团队ID或收银员ID)
    SType     string `json:"s_type" gorm:"index"`                 // 服务类型：merchant, cashier_team, cashier
    TrxID     string `json:"trx_id" gorm:"uniqueIndex"`           // 交易唯一标识
    AccountID string `json:"account_id" gorm:"index"`             // 关联账户ID
    *WithdrawValues
    CreatedAt int64 `json:"created_at" gorm:"autoCreateTime:milli"`
    UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

type WithdrawValues struct {
    Status      *string          `json:"status"`       // pending, processing, success, failed
    Ccy         *string          `json:"ccy"`          // 币种
    Amount      *decimal.Decimal `json:"amount"`       // 提现金额
    Fee         *decimal.Decimal `json:"fee"`          // 手续费
    ChannelCode *string          `json:"channel_code"` // 渠道代码
    NotifyURL   *string          `json:"notify_url"`   // 回调地址
    Country     *string          `json:"country"`      // 国家
    // 状态时间戳
    CanceledAt  *int64 `json:"canceled_at"`  // 取消时间
    CompletedAt *int64 `json:"completed_at"` // 完成时间
    ExpiredAt   *int64 `json:"expired_at"`   // 过期时间
    ConfirmedAt *int64 `json:"confirmed_at"` // 确认时间
}

// 通过 Sid + SType 区分业务主体:
// - SType="merchant", Sid=商户ID: 商户提现
// - SType="cashier_team", Sid=团队ID: 收银团队提现
```

#### 5.5.4 资产模型

```go
// Asset 资产模型，支持多资金属性
type Asset struct {
    Balance          decimal.Decimal `json:"balance"`           // 总余额
    AvailableBalance decimal.Decimal `json:"available_balance"` // 可用余额
    FrozenBalance    decimal.Decimal `json:"frozen_balance"`    // 冻结余额
    MarginBalance    decimal.Decimal `json:"margin_balance"`    // 保证金余额
    ReserveBalance   decimal.Decimal `json:"reserve_balance"`   // 预留余额
    Currency         string          `json:"currency"`          // 币种
    UpdatedAt        int64           `json:"updated_at"`        // 更新时间
}
```

## 6. 安全认证体系

### 6.1 认证方式

#### API Key 认证
- **应用场景**: OpenAPI ↔ 用户系统，OpenAPI ↔ CashierAPI
- **特点**: 无状态认证，适合系统间调用
- **实现**: `middleware.APIKeyAuth()`

#### JWT 认证
- **应用场景**: 管理后台用户认证
- **特点**: 有状态会话认证，支持用户登录
- **JWT Claims**:
```go
type JWTClaims struct {
    UserType   string `json:"user_type"` // merchant, cashier_team, admin
    MerchantID string `json:"merchant_id"` 
    CashierID  string `json:"cashier_id"`
    Role       string `json:"role"`
}
```

### 6.2 权限控制

#### 三层权限体系
- **Merchant**: 商户权限，只能访问自己的数据
- **CashierTeam**: 收银团队权限，只能访问自己团队的数据
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
| CashierAdmin | 8083 | 收银团队管理后台 |
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

### 8.1 业务角色扩展

**新增用户角色类型**：
1. 在统一Account表中新增UserType值
2. 在业务表中新增对应的SType值
3. 实现对应的Service接口
4. 配置相应的权限和管理后台

**示例：新增"代理商"角色**：
```go
// 1. Account表支持
UserType = "agent"

// 2. 业务表支持  
SType = "agent"
Sid = "代理商ID"

// 3. Service实现
type AgentDepositService struct {
    // 实现统一的DepositServiceInterface
    // SType = "agent"
}

// 4. 管理后台
AgentAdmin (Port: 8085)
```

### 8.2 业务类型扩展

**新增业务交易类型**：
1. 创建新的业务表（如Transfer表）
2. 使用统一的Sid + SType设计
3. 实现对应的Service接口
4. 在Transaction抽象层注册新的TrxType

**统一业务表模板**：
```go
type NewBusinessTable struct {
    ID        uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
    Sid       string `json:"sid" gorm:"index"`          // 服务主体ID
    SType     string `json:"s_type" gorm:"index"`       // 服务类型
    TrxID     string `json:"trx_id" gorm:"uniqueIndex"` // 交易唯一标识
    AccountID string `json:"account_id" gorm:"index"`   // 关联账户ID
    *NewBusinessValues
    CreatedAt int64 `json:"created_at" gorm:"autoCreateTime:milli"`
    UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
```

### 8.3 渠道扩展

**当前渠道架构**：
- **CashierAPI渠道**：ChannelCode = "Cashier"
  - CashierTeam在调用层处理具体业务
  - 多团队独立运营
- **ThirdPartyAPI渠道**：包含银行和第三方支付API
  - 不再单独设置银行渠道
  - 统一作为第三方API处理

**新增支付渠道**：
1. 实现标准的渠道接口
2. 在渠道配置中注册ChannelCode
3. 配置路由规则
4. 如需调用层，可参考CashierTeam模式

### 8.4 收银团队扩展

新增收银团队：
1. 创建新的团队Account（UserType="cashier_team"）
2. 配置独立的资金池和银行通道
3. 分配给CashierAdmin团队管理
4. 业务数据自动通过SType区分

### 8.5 管理功能扩展

- 支持新的管理角色和权限体系
- 扩展权限控制粒度
- 增加新的管理功能模块
- 基于SType的数据权限控制

### 8.6 系统集成扩展

- 支持新语言的错误消息
- 多时区支持
- 多币种支持
- 统一的API接口设计便于第三方集成

## 📝 总结

inpayos 采用的分层多服务架构设计具有以下优势：

### 🎯 核心特色

1. **商户中心化设计**: 
   - 商户(Merchant)是唯一的外部用户角色
   - 所有功能和服务都围绕商户需求设计
   - 通过OpenAPI为商户提供统一的支付服务入口

2. **渠道化架构**: 
   - CashierTeam在CashierAPI调用层处理具体业务，不是独立用户角色
   - CashierAPI作为被调用的渠道之一，通过ChannelCode="Cashier"标识
   - ThirdPartyAPI渠道统一处理银行和第三方支付API
   - 简化的双渠道架构：CashierAPI + ThirdPartyAPI

3. **统一业务表设计**: 
   - **Deposit/Withdraw表统一设计**：通过SType字段区分业务主体
   - **SType字段**：主要标识"merchant"业务
   - **Sid字段**：标识具体的商户ID
   - **渠道路由**：通过ChannelCode进行渠道选择和路由
   - **一套表结构支持多种业务角色**，简化维护和扩展

4. **核心共享服务设计**: 
   - **核心业务服务共享**：交易、账户、流水、充值、提现、结算规则、定时任务、消息服务
   - **专属业务模块分离**：各系统的特定业务逻辑在各自Admin层实现
   - **服务无关性**：核心服务不包含特定角色业务逻辑，高度复用
   - **统一接口规范**：通过SType参数区分不同业务主体

5. **三层管理体系**: 
   - MerchantAdmin：商户自管理 + 商户专属业务模块
   - CashierAdmin：收银团队自管理 + 团队专属业务模块
   - GlobalAdmin：全局系统管理 + 平台专属业务模块
   - 各层权限清晰，数据隔离，业务模块独立

### 🏗️ 架构优势

1. **统一接口**: 所有用户通过 OpenAPI 接入，简化集成
2. **渠道灵活**: 支持多种支付渠道，易于扩展
3. **权限清晰**: 不同角色管理各自数据，严格权限隔离
4. **扩展性强**: 可单体部署或拆分为微服务
5. **国际化**: 完整的多语言支持
6. **安全可控**: 多层认证和权限控制机制

### 🚀 业务价值

1. **灵活性**: 支持商户直接接入和专业收银团队服务两种业务模式
2. **可扩展**: 平级角色设计便于后续增加新的用户类型
3. **数据安全**: 统一账户体系确保数据一致性和安全性
4. **管理高效**: 分层管理体系提升运营效率
5. **开发效率**: 核心服务共享减少重复开发，一套服务支持多种角色
6. **维护简单**: 业务表结构统一，核心服务集中管理
7. **职责清晰**: 核心服务专注数据操作，业务模块专注角色逻辑

### 🔧 技术优势

1. **表结构优化**: Deposit/Withdraw等业务表通过Sid+SType统一设计，避免重复表结构
2. **服务层分离**: 核心共享服务专注数据操作，专属业务模块处理角色逻辑
3. **接口标准化**: 核心服务提供统一接口，通过参数区分不同业务主体
4. **数据一致性**: 统一的Account表确保跨角色的数据一致性
5. **权限隔离**: 基于SType的数据隔离保证业务安全性
6. **模块解耦**: 核心服务与业务模块解耦，便于独立开发和维护

这种架构设计既满足了支付网关的复杂业务需求，又保持了良好的可扩展性和可维护性。特别是**核心共享服务设计**、**统一业务表设计**和**CashierTeam作为与商户平级的角色**的创新设计，为构建更加灵活、高效和专业的支付平台提供了坚实的基础。