import React from 'react';
import { LaptopOutlined, NotificationOutlined, UserOutlined } from '@ant-design/icons';
import { MenuProps, Typography } from 'antd';
import { Breadcrumb, Layout, Menu, theme } from 'antd';

const { Header, Content, Sider } = Layout;


const sideNav: MenuProps['items'] = [
  {
    key: String(1),
    icon: React.createElement(LaptopOutlined),
    label: 'Tasks'
  },
  {
    key: String(2),
    icon: React.createElement(LaptopOutlined),
    label: 'Artifacts'
  },
  {
    key: String(3),
    icon: React.createElement(LaptopOutlined),
    label: 'Settings'
  }
];

const Dashboard: React.FC = () => {
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  return (
    <Layout style={{ height: 'inherit' }}>
      <Header style={{ backgroundColor: 'white', borderBottom: '1px solid #cecece', display: 'flex', flexDirection: 'column' }} className="header">
        <div className="logo" />
        <Typography.Title level={3}>Chrononut</Typography.Title>
        {/* <Menu mode="horizontal" defaultSelectedKeys={['2']}  /> */}
      </Header>
      <Layout>
        <Sider width={200} style={{ background: colorBgContainer }}>
          <Menu
            mode="inline"
            defaultSelectedKeys={['1']}
            defaultOpenKeys={['sub1']}
            style={{ height: '100%', borderRight: 0 }}
            items={sideNav}
          />
        </Sider>
        <Layout style={{ padding: '0 24px 24px' }}>
          <Content
            style={{
              padding: 24,
              margin: 0,
              minHeight: 280,
              background: colorBgContainer,
            }}
          >
            Content
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
};

export default Dashboard;