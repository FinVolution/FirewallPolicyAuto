import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { RouterProvider } from 'react-router-dom'
import router from '@/routes/index'
import zhCN from 'antd/locale/zh_CN';
import { ConfigProvider } from 'antd'
import { TasksProvider } from './context/menu'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ConfigProvider locale={zhCN}>
      <TasksProvider>
        <RouterProvider router={router} />
      </TasksProvider>
    </ConfigProvider>
  </StrictMode>,
)
