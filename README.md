<h1 align="center">BronyaBot</h1>

---

## 📚 工学云 & 超星学习通自动化项目

> 一站式签到解决方案！
> BronyaBot 是一个基于 Go 语言开发的自动化签到项目，专为工学云平台和超星学习通设计。通过模拟登录、自动签到、考试答题（学习通部分尚在开发中），让你轻松完成每日任务！适用于个人用户和团队使用。
---

## 🤖 关于 BronyaBot 的命名

项目名称 BronyaBot 灵感来源于《崩坏3》中的角色 布洛妮娅·扎伊切克 (Bronya Zaychik)。

* ~~布洛妮娅以 高智商、高效率和冷静 的特点著称，这与本项目追求的 稳定性 和 高效性 完美契合。~~
* ~~她善于分析和解决复杂问题，正如 BronyaBot 对签到流程的模拟和优化，甚至应对 工学云最新版滑块验证 的复杂挑战。~~

~~取名 BronyaBot，不仅是对布洛妮娅智慧与坚韧的致敬，也代表项目对 可靠性 和 优雅解决复杂问题 的追求！
(bushi~~

---

## **✨ 功能特色**

* ✅ 自动化签到：支持工学云平台签到，省时省力。
* 📫 邮箱支持: 完成任务点时自动发送邮件提示
* 🛡️ 滑块验证支持：已解决工学云最新版滑块验证问题，自动突破验证机制。
* 📖 学习通功能：自动刷课、考试和答题（开发中，敬请期待）。
* 🗂️ 数据库支持：安全管理用户信息和签到记录。
* ⚙️ 高扩展性：轻松集成更多平台和功能。
* 🔒 安全可靠：使用加密技术保障用户信息安全。

---

## 🚀 快速开始

### 1.环境要求

* 🛠️ Go 版本： 1.23.3 或更高版本
* 🗄️ 数据库支持：MySQL

### 2. 配置文件

在项目根目录下 [config](configuration/application.yaml) 文件，并根据需求填写以下字段：

```yaml

mysql:
  dataBase: # 数据库名
  userName: # 数据库用户名
  passWord: # 数据库密码
  port: # 数据库端口，默认为3306
  driverName: # 数据库驱动名，例如：mysql
  host: # 数据库主机地址，例如：127.0.0.1 或 localhost
  log-level: debug # MySQL 日志级别，可选值：debug, info, warn, error
  max-idle-conns: 10 # 最大空闲连接数
  max-open-conns: 100 # 最大打开连接数
  conn-max-life-time: 5 # 连接最大生存时间，单位：分钟

logger:
  level: info # 日志级别，可选值：debug, info, warn, error
  prefix: '[🛠️]' # 日志前缀
  director: log # 日志文件存放目录
  showLine: true # 是否显示日志行号
  logInConsole: true # 是否在控制台输出日志

mail:
  host: # 邮件服务器主机地址
  port: # 邮件服务器端口
  user: # 邮件发送用户
  password: # 邮件发送密码
  ssl: true # 是否使用 SSL 加密
  local-home:           # 可选，本地邮件存储目录 (可选配置)

ai:
  token:  # 使用讯飞星火大模型生成报告所需的 token
account:
  gongxueyun: #对工学云单用户无数据库服务新开账号密码登录
    off: true                                #true为开启本地用户加载，将关闭数据库连接
    phone: 1111111111                          #phone: 手机号码
    password: 11dwhdiwhdi                    #password: 密码
    country: 中国                            #country: 国家
    province: 浙江省                         #province: 省份
    city: 宁波市                               #city: 城市
    area: xx区                              #area: 地区
    latitude: 79.783787878                        #latitude: 纬度
    longitude: 120.898978978978                      #longitude: 经度
    email: xxxxxx@qq.com                  #email: 邮箱地址
    address: 浙江省·xxxx.xxxxx  #address: 详细地址
  cx:
    phone: xxxxxx
    password: xxxxxxx
  tasks:
    - cron: "0 8 * * *"
      description: "每天早上8点签到"
      taskType: "sign"
    - cron: "0 18 * * *"
      description: "每天晚上6点签到"
      taskType: "sign"
    - cron: "0 10 * * 5"
      description: "每周周五早上10点提交周报"
      taskType: "week"
    - cron: "0 10 L * *"
      description: "每月最后一天早上10点提交月报"
      taskType: "month"




```

---

## 🗂️ 数据库表设计

以下是用于存储工学云用户签到信息的表结构定义，包含基本的用户信息和签到状态字段,
在数据库中执行以下 SQL 命令创建 工学云用户签到表：

```sql
-- 工学云用户签到表
CREATE TABLE sign
(
    id        INT AUTO_INCREMENT PRIMARY KEY, -- 签到记录 ID (自增主键)
    username  VARCHAR(255) NOT NULL,          -- 用户名
    password  VARCHAR(255) NOT NULL,          -- 密码 (加密存储建议)
    country   VARCHAR(255) DEFAULT NULL,      -- 国家
    province  VARCHAR(255) DEFAULT NULL,      -- 省份
    city      VARCHAR(255) DEFAULT NULL,      -- 城市
    area      VARCHAR(255) DEFAULT NULL,      -- 区域
    latitude  VARCHAR(255) DEFAULT NULL,      -- 纬度
    longitude VARCHAR(255) DEFAULT NULL,      -- 经度
    email     VARCHAR(255) DEFAULT NULL,      -- 用户邮箱
    address   VARCHAR(255) DEFAULT NULL,      -- 完整地址
    type      INT          DEFAULT 0,         -- 签到类型 (0: 普通签到, 1: 特殊签到等)
    state     INT          DEFAULT 0          -- 签到状态 (0: 未签到, 1: 已签到)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

```

---

## 🌟 主要模块说明

| 模块      | 描述              | 状态     |
|---------|-----------------|--------|
| 工学云签到   | 模拟登录并完成每日签到任务   | 已完成 ✅  |
| 滑块验证破解	 | 自动处理工学云最新版滑块验证	 | 已完成 ✅  |
| 工学云定时上报 | 自动完成上报周报月报      | 开发中✅  |
| 超星学习通   | 自动刷课、考试、答题      | 开发中 🚧 |
| 用户管理    | 基于数据库管理用户信息     | 已完成 ✅  |
| 日志记录    | 自动记录签到状态和错误信息   | 已完成 ✅  |
| 发送邮件    | 完成任务点自动发送当前状态   | 已完成 ✅  |

---

## 📦 项目结构

```plaintext
├───.idea
├───config
├───configuration               # 配置文件
├───core                        # 核心文件
├───global                      # 全局配置
├───internal                   
│   ├───api                     # 平台接口
│   ├───entity                  # 数据库实体
│   └───service                 # 服务
│       ├───cx_service
│       │   └───data
│       └───gongxueyun_service
│           └───data
└───utils                       # 工具库
    └───blockPuzzle
```

---

## 🛠️ 开发计划

* ✅工学云平台自动签到
* ✅解决工学云最新版滑块验证
* ✅完成对邮箱发送的支持
* ✅工学云平台自动周报月报
* ⬜超星学习通自动刷课
* ⬜添加 Docker 支持
* ⬜增强错误处理和日志记录

---

## 📜 使用须知

* 本项目仅供学习与研究使用，请勿用于非法用途。
* 在使用时，请遵守相关平台的使用规则。

---

## 📈 GitHub Stars 统计

Star 趋势图
让 BronyaBot 成为你学习路上的得力助手！🎉 签到从此无忧！

![GitHub Stars](https://img.shields.io/github/stars/mirai-MIC/BronyaBot?style=flat&label=Stars)
[![Star History Chart](https://starchart.cc/mirai-MIC/BronyaBot.svg?variant=adaptive)](https://starchart.cc/mirai-MIC/BronyaBot)

---

## 📧 联系我们

如有问题或建议，请 [点击链接](https://qm.qq.com/q/z8mab8YVm8)。

**不支持定制不接单，需要的可以联系我代挂**

快来试试吧！🎉 简化签到，从现在开始！

---

## 💵赞助我们

<img src="https://github.com/user-attachments/assets/017c8856-ea97-4ab4-84e8-caed61b33268" width="400" height="450" />

