# 项目依赖

本项目由`React + TypeScript + Vite`构建，`Vite`版本 5.4.1，本地运行需要 `Node.js` 版本 14.18+。

因部分模块限制，打包需要 `Node.js` 版本 16+ 。

推荐使用环境 `Node.js` 16+

## 使用方式

### Development

```js
// 配置service地址
// vite.config.ts
target: "http://xxx.xxx.xxx.xxx:8080"
```

```js
// 安装依赖包
npm install
// 运行
npm run dev
```

### Build

```js
// 命令
npm run build
// 产物路径
dist
```

## 其他依赖项
### 线上依赖
```json
  "@ant-design/icons": "^5.4.0",        # ICON
  "antd": "^5.20.1",                    # UI库
  "axios": "^1.7.4",                    # 请求库
  "lodash": "^4.17.21",                 # 工具库
  "qs": "^6.13.0",                      # 路由参数处理
  "react": "^18.3.1",                   # React
  "react-dom": "^18.3.1",               
  "react-router-dom": "^6.26.1"         # 路由
```

### 开发依赖
```json
  "@eslint/js": "^9.9.0",                         # eslint
  "@types/lodash": "^4.17.7",                     # lodash类型声明
  "@types/node": "^22.4.1",                       # node类型声明
  "@types/qs": "^6.9.15",                         # qs类型声明
  "@types/react": "^18.3.3",                      # react类型声明
  "@types/react-dom": "^18.3.0",                  # react类型声明
  "@vitejs/plugin-react": "^4.3.1",               # vite插件 - react
  "eslint": "^9.9.0",                             # eslint
  "eslint-plugin-react-hooks": "^5.1.0-rc.0",     
  "eslint-plugin-react-refresh": "^0.4.9",        
  "globals": "^15.9.0",                           
  "typescript": "^5.5.3",                         # ts
  "typescript-eslint": "^8.0.1",                  
  "vite": "^5.4.1",                               # vite
  "vite-plugin-svgr": "^4.2.0"                    # vite插件 - svg使用
  "prettier": "^3.3.3",
  "svgo": "^3.3.2",
```

# 目录结构

```
├─node_modules        # 依赖包
├─public              # 静态资源
├─src             
│ ├─assets            # 静态资源
│ ├─context           # 状态管理
│ ├─layout            # 布局
│ ├─pages             # 页面
│ │  └─Policy         # 策略查询
│ │     └─index.tsx
│ ├─routes            # 路由
│ ├─server            # api封装
│ │  ├─xxx.ts         # api
│ │  └─axios.ts       # axios封装
│ └─main.tsx          # 根组件
├─package.json        # 依赖包
├─tsconfig.json       # ts配置
├─eslint.config.js    # eslint配置
├─index.html          # 入口文件
└─vite.config.ts      # vite配置

```