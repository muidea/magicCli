# magicCli 开发规则

本文件直接规定本仓的职责、边界和开发要求；单独检出本仓即可阅读执行。仓库内位置均相对于仓库根目录，协作项目通过项目名和 Git 地址识别，不依赖固定工作区或外部规范文件才能理解规则。

## 功能定位与边界

- 是基于 Cobra 的通用命令执行工具，选择 Local／Docker／SSH Executor，传递 context、命令输出及退出结果。
- 不拥有应用定义、安装审批、CAS 授权、数据迁移或业务恢复；命令退出0不能代替平台安装／业务提交证明。
- `internal/modules/kernel/executor` 是执行器工厂，不是 framework Module，不因目录名称引入插件生命周期。
- 保留目标选择、取消、错误及退出码合同；新导出接口须有真实调用者，不引入平台名或业务阶段特例，不记录秘密。

## 协作项目

- 本仓直接 Go 依赖及工具链以 go.mod 为准；不因同处工作区而假定依赖 magicCommon／magicBase。
- magicRunner — https://github.com/muidea/magicRunner.git：平台安装／协调 owner；magicPlatform — https://github.com/muidea/magicPlatform.git：部署交付 owner。本工具不替代两者的流程。

## 验证与交付

- 从仓库根检查实际入口、Makefile 和测试，再运行受影响测试及 `go build ./...`、`go vet ./...`；Go 代码用 gofmt／goimports，测试覆盖不同目标、取消和失败退出。
- Docker／SSH 和远端命令只使用任务已授权的目标，不为规则或文档修改执行实际远端动作；不硬编码目标主机、用户名或秘密。
- 若处于多仓工作区，按工作区 `upgrade_vendor.sh` 保护／同步所需依赖并审查差异，禁止手改 vendor；独立检出明确未验的消费者范围。
- 公开行为变化同步命令帮助及交付材料，工具验证和产品 Live 分开。规则直接维护在本文件，移除已删除入口或脚本；不提交凭据、构建产物或运行数据，提交按仓库分开。
