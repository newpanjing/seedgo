# Smart Butler Frontend

## 简介
这是一个基于 Vue 3 + TypeScript + Vite 构建的现代化前端项目。项目集成了 Tailwind CSS 进行样式开发，使用 Pinia 进行状态管理，并实现了完善的权限管理系统（RBAC）。

## 技术栈
- **核心框架**: [Vue 3](https://vuejs.org/)
- **构建工具**: [Vite](https://vitejs.dev/)
- **语言**: [TypeScript](https://www.typescriptlang.org/)
- **路由**: [Vue Router](https://router.vuejs.org/)
- **状态管理**: [Pinia](https://pinia.vuejs.org/)
- **UI 样式**: [Tailwind CSS](https://tailwindcss.com/)
- **UI 组件库**: 基于 [Radix Vue](https://www.radix-vue.com/) 定制 (Shadcn-like)
- **图标**: [Lucide Vue](https://lucide.dev/)
- **图表**: [ECharts](https://echarts.apache.org/)
- **HTTP 请求**: [Axios](https://axios-http.com/)
- **工具库**: [VueUse](https://vueuse.org/)

## 项目结构

```
frontend/
├── src/
│   ├── api/             # API 接口定义 (Auth, User, Role 等)
│   ├── assets/          # 静态资源文件
│   ├── components/      # 组件库
│   │   ├── business/    # 业务组件 (如 DictSelect, RoleSelect)
│   │   ├── common/      # 通用高阶组件 (DataTable, Form, Tree)
│   │   └── ui/          # 基础 UI 原子组件 (Button, Input, Dialog 等)
│   ├── composables/     # 组合式函数 (Hooks)
│   ├── directives/      # 自定义指令 (v-permission, v-super)
│   ├── layouts/         # 页面布局 (MainLayout, Sidebar, Header)
│   ├── lib/             # 核心库代码与工具
│   ├── router/          # 路由配置与守卫
│   ├── stores/          # Pinia 状态仓库 (Auth, Permission, Users 等)
│   ├── types/           # TypeScript 类型定义
│   ├── utils/           # 工具函数 (Request, Date, Tenant)
│   ├── views/           # 页面视图代码
│   │   ├── auth/        # 权限相关页面 (Users, Roles, Permissions)
│   │   ├── system/      # 系统管理页面 (Tenants, Dicts, Logs)
│   │   ├── login/       # 登录页
│   │   └── ...          # 其他业务页面 (Dashboard, Products 等)
│   ├── App.vue          # 应用根组件
│   ├── main.ts          # 应用入口文件
│   └── style.css        # 全局样式入口
├── package.json         # 项目依赖与脚本配置
├── vite.config.ts       # Vite 构建配置
├── tailwind.config.js   # Tailwind CSS 配置
└── tsconfig.json        # TypeScript 配置
```

## 核心功能

### 1. 权限管理
- **路由权限**: 基于 `src/router/index.ts` 中的全局前置守卫 (`router.beforeEach`) 实现。结合 `useAuthStore` 和 `usePermissionStore` 进行登录状态校验和动态菜单加载。
- **动态菜单**: 根据后端返回的权限数据动态生成侧边栏菜单。
- **指令控制**: 
  - `v-permission`: 用于细粒度的功能按钮权限控制。
  - `v-super`: 超级管理员权限控制。

### 2. 系统管理模块
- **租户管理**: 支持多租户架构，提供租户的增删改查。
- **用户与角色**: 完整的 RBAC 模型，支持用户分配角色，角色分配权限。
- **字典管理**: 系统通用字典配置，便于前端下拉框等数据源管理。
- **日志审计**: 包含登录日志和操作日志，记录系统关键行为。

### 3. 组件化开发
- **DataTable**: 封装了基于配置的表格组件，支持分页、自定义列、操作栏等功能。
- **Form System**: 封装了 `Form`、`FormItem` 以及集成的 `FormDialog` (弹窗表单) 和 `FormSheet` (抽屉表单)，简化表单开发。
- **UI Components**: 大量使用基于 Radix Vue 的无头组件，配合 Tailwind CSS 实现高度可定制的 UI 组件。

## 开发与构建

### 前置要求
- Node.js (推荐 v18+)
- npm 或 pnpm

### 安装依赖
```bash
npm install
```

### 启动开发服务器
```bash
npm run dev
```
启动后访问控制台输出的本地地址 (通常为 `http://localhost:5173`)。

### 构建生产版本
```bash
npm run build
```
构建产物将输出到 `dist` 目录。

### 代码预览
```bash
npm run preview
```

## 配置说明
- **API 代理**: 开发环境下，在 `vite.config.ts` 中配置了 `/api` 前缀的代理，默认转发到 `http://localhost:3000`。
- **自动导入**: 项目配置了 `unplugin-auto-import` 和 `unplugin-vue-components`：
  - 自动导入 Vue API (`ref`, `computed` 等)。
  - 自动导入 `src/components` 下的组件。
  - 自动生成 `auto-imports.d.ts` 和 `components.d.ts` 以支持 TypeScript 类型推导。
