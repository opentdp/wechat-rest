# Wrest Chat

智能聊天助手，是一个通用的聊天辅助程序，通过 [Nanomsg 协议](wcferry/proto/wcferry.proto) 与聊天软件互通，内置 WEB 管理界面，可接入GPT、Gemini、星火、文心、混元、通义千问等大语言模型。目前已适配 *PC微信*，更多聊天软件适配中，敬请期待！

> 本项目**未对微信程序进行任何破解或修改**，与微信互操作的能力均基于开源项目 [WeChatFerry RPC](https://github.com/lich0821/WeChatFerry/tree/master/WeChatFerry) 实现，感谢各位开源贡献者。

## 主要特性

- 使用 Go 语言编写，无运行时依赖
- 提供 HTTP 接口，便于对接各类编程语言
- 提供 Websocket 接口，接收推送的新消息
- 支持 HTTP/WS 接口授权，参见 [配置说明](#配置说明)
- 支持作为 SDK 使用，参见 [wcferry/README.md](./wcferry/README.md)
- 内置 AI 机器人，参见 [wclient/README.md](./wclient/README.md)
- 内置 Web 管理界面，可以管理机器人各项配置
- 内置 Api 调试工具，所有接口都可以在线调试
- 尽可能将消息中的 Xml 转为 Object，便于前端解析
- 支持计划任务、外部指令、指令插件等扩展功能，详见 [wrest-plugin](https://github.com/opentdp/wrest-plugin)

## 快速开始

请仔细阅读[免责声明](#免责声明)和[常见问题](#常见问题)后再开始使用，首次使用可参照下面的步骤开始：

- 下载并安装 [WeChatSetup-3.9.2.23.exe](https://github.com/opentdp/wrest-chat/releases/download/v0.0.1/WeChatSetup-3.9.2.23.exe) 和 [wrest-chat.zip](https://github.com/opentdp/wrest-chat/releases)

  - 非开发者请直接下载编译好的二进制文件，不要下载源码

- 双击 `wrest.exe` 将自动启动微信和接口服务，扫码登录微信

  - 启动成功后，浏览器访问 `http://localhost:7600` 配置机器人

- 若无人值守，可选择使用 `starter.bat` 启动服务，实现如下能力：
  
  - 写入禁止微信自动更新的注册表配置
  - 在 `wrest.exe` 崩溃后自动重启

## 配置说明

机器人相关参数均已支持从 WEB 界面管理，[config.yml](./config.yml) 用来配置一些核心能力，一般情况下保持默认即可。

- 修改 `config.yml` 中的参数，需重启 **wrest.exe** 才能生效

  - 请使用 `Ctrl + C` 终止 **wrest.exe**，切勿直接关闭任务窗口
  - 重启时，提示端口被占用，请退出微信后重试

- 设置 `Web.Token` 后，请求接口时必须携带 **header** 信息: `Authorization: Bearer $token`

## 开发指南

模块依赖示意：`WEB ---> API ---> BOT ---> SDK ---> Wcferry ---> WeChat`。其中 `BOT` 模块并非必须的，可根据自己的需求选择是否开启，`Wcferry` 模块为第三方开源依赖，必须和 `WeChat` 版本匹配使用。

查看和调试 *HTTP/WS* 接口，可使用浏览器访问 `http://localhost:7600/swagger/`。更多插件开发资源请查阅 [wrest-plugin](https://github.com/opentdp/wrest-plugin) 项目。

### API 模块

实现了 HTTP/WS 接口，详情查看 [httpd/README.md](./httpd/README.md)

### BOT 模块

实现了群聊机器人，详情查看 [wclient/README.md](./wclient/README.md)

### SDK 模块

实现了 WCF 客户端，详情查看 [wcferry/README.md](./wcferry/README.md)

### WEB 模块

实现了 WEB 控制台，详情查看 [webview/README.md](./webview/README.md)

## 代码提交

提交代码时请使用 `feat: something` 作为说明，支持的标识如下

- `feat` 新功能（feature）
- `fix` 错误修复
- `docs` 文档更改（documentation）
- `style` 格式（不影响代码含义的更改，空格、格式、缺少分号等）
- `refactor` 重构（即不是新功能，也不是修补bug的代码变动）
- `perf` 优化（提高性能的代码更改）
- `test` 测试（添加缺失的测试或更正现有测试）
- `chore` 构建过程或辅助工具的变动
- `revert` 还原以前的提交

## 常见问题

### Q1、提示注入失败

当前分支兼容的**PC微信**版本是 `3.9.2.23`，请在  [快速开始](#快速开始) 中点击下载。

### Q2、如何在群内 `@` 其他人

首先要在消息中添加 `@昵称`，然后在 `aters` 参数添加此人的 `wxid`。相关接口 `/wcf/send_txt`。

### Q3、如何更新机器人程序

从 [快速开始](#快速开始) 中下载新版本。关闭机器人后，将解压出来的文件覆盖过去。若修改过 `config.yml` 请排除，防止配置被重置。

### Q4、常用 AI 密钥获取地址

- 阿里 通义千问 <https://dashscope.console.aliyun.com/apiKey>
- 讯飞 星火 <https://console.xfyun.cn/services/bm3>
- 谷歌 Gemini <https://aistudio.google.com/app/apikey?hl=zh-cn>

### Q5、更改系统内置指令参数

登录控制台 `http://localhost:7600/#/keyword/list`，添加一个**指令别名**，修改需要的参数，保存后生效。

- 场景选择`全局`，修改级别为需要的值，可**覆盖内置指令级别**
- 所有内置指令修改级别为管理员以上的值，相当于在全局**禁用内置指令**
- 添加一个名为 `@xxx` 的指令别名，目标选择 `/ai`，可实现使用 `@xxx` 唤醒

### Q6、消息不推送或管理页卡住

点击CMD窗口的左上角，选择`属性`，取消`快速编辑模式`，保存后关闭。在窗口中按下任意键（理论上应该无反应），检查下机器人是否恢复正常，否则请重启机器人。

### Q7、最低配置是什么
#### 1. Windows 平台
##### 1.1 Windows PC
- Windows 10，CPU 1 核，内存 4GB
##### 1.2 Windows Server
- Windows Server 2012，CPU 1 核，内存 4GB
#### 2. Linux 平台
- Todo
#### 3. macOS 平台
- Todo

## 免责声明

[WrestChat](https://github.com/opentdp/wrest-chat) 和 [WeChatFerry](https://github.com/lich0821/WeChatFerry) 是供学习交流的开源项目，代码及其制品仅供参考，不保证质量，不构成任何商业承诺或担保，不得用于商业或非法用途，使用者自行承担后果。
