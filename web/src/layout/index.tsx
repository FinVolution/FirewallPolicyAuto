import { FC } from 'react';
import Icon from '@ant-design/icons';
import { Breadcrumb, Layout, Menu } from 'antd';
import type { MenuProps } from 'antd';
import { Outlet, useNavigate } from 'react-router-dom'
import './index.css';
import { useMenus } from '@/context/menu';
import PolicyIcon from '@/assets/policy.svg?react'

const { Header, Content } = Layout

type MenuItem = Required<MenuProps>['items'][number];

type MenuClickEventHandler = Required<MenuProps>['onClick'];

const items: MenuItem[] = [
  {
    key: 'policy',
    icon: <Icon component={PolicyIcon} />,
    label: '策略查询',
  },
];

interface Props {

}
const LayoutPage: FC<Props> = () => {
  const menus = useMenus();
  const nav = useNavigate()
  const onClick: MenuClickEventHandler = (info) => {
    nav(info.key)
  }

  return <Layout className='layout'>
    <Header className="header">
      <div className="logo" />
      <Menu
        theme='dark'
        selectedKeys={[menus.selectedKeys]}
        mode="horizontal"
        style={{ marginLeft: '40px' }}
        items={items}
        onClick={onClick}
      />
    </Header>
    <Layout>
      {/* <Sider width={200} style={{ background: '#fff' }} >
        <Menu
          selectedKeys={[menus.selectedKeys]}
          mode="inline"
          style={{ height: '100%', borderRight: 0 }}
          items={items}
          onClick={onClick}
        />
      </Sider> */}
      <Layout style={{ padding: '0 48px 24px' }} >
        <Breadcrumb style={{ margin: '16px 0' }} items={menus.items} />
        <Content
          className="site-layout-background layout-content"
          style={{
            padding: 24,
            margin: 0,
            minHeight: 280,
            borderRadius: '8px',
            boxShadow: '0 0px 10px rgba(0, 0, 0, 0.1)'
          }}
        >
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  </Layout>;
}
export default LayoutPage;