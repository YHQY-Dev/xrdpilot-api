# XrdPilot 定制说明

本仓库基于 [QuantumNous/new-api](https://github.com/QuantumNous/new-api) 的 **v1.0.0-rc.26** 做站点定制，当前生产镜像是 `calciumion/new-api:v1.0.0-rc.26`。

**不要**用官方镜像直接覆盖正在运行的 XrdPilot 实例。上线时应构建并推送本仓库的自定义镜像，再替换服务器上的镜像标签。

## Git

- `upstream`：官方仓库 `https://github.com/QuantumNous/new-api.git`
- `origin`：本定制仓库
- 工作分支：`xrdpilot`（从 tag `v1.0.0-rc.26` 拉出）

## 产品行为

- 公开注册立即开通，不做人工审核。
- 对外用邮箱 + 密码登录；老账号仍可用原来的用户名登录。
- 注册不再收集用户名。内部 `users.username` 由邮箱 local-part 自动生成（最长 20，冲突加数字后缀）。
- 真实姓名写入 `display_name`（最长 20，按 rune 截断），完整姓名另存 `real_name`。
- 科研资料单独列存储，不写入 `remark`。

## 升级时必碰文件

跟官方后续 RC 时，优先复查这些定制点：

- `model/user.go`：User 新增列、搜索字段、登录按用户名或邮箱查找
- `model/user_profile.go`：用户名生成、资料校验（尽量把新逻辑放这里）
- `controller/user.go`：`Register`
- `web/src/features/auth/lib/registration-profile.ts`
- `web/src/features/auth/sign-up/components/sign-up-form.tsx`
- `web/src/features/auth/sign-in/components/user-auth-form.tsx`
- `web/src/features/users/components/users-columns.tsx`
- `web/src/features/users/components/users-mutate-drawer.tsx`
- `web/src/i18n/locales/zh.json`（以及 backend `i18n/locales/*.yaml` 中的 `user.profile_invalid`）

不要改官方品牌 / 项目名等受保护信息。
