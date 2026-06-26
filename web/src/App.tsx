import { Layout, Menu, Spin } from 'antd'
import { lazy, Suspense, useMemo } from 'react'
import { Link, Route, Routes, useLocation } from 'react-router-dom'

const HomePage = lazy(() => import('./pages/HomePage'))
const UsersPage = lazy(() => import('./pages/UsersPage'))
const UserDetailPage = lazy(() => import('./pages/UserDetailPage'))
const ItemDetailPage = lazy(() => import('./pages/ItemDetailPage'))
const QueuePage = lazy(() => import('./pages/QueuePage'))

const { Header, Content } = Layout

export default function RootApp() {
  const location = useLocation()
  const selectedKey = useMemo(() => {
    if (location.pathname.startsWith('/users') || location.pathname.startsWith('/user/')) return '/users'
    if (location.pathname.startsWith('/queue') || location.pathname.startsWith('/explore/queue')) return '/queue'
    if (location.pathname.startsWith('/explore/users')) return '/users'
    return '/'
  }, [location.pathname])

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ position: 'sticky', top: 0, zIndex: 10, width: '100%' }}>
        <Menu
          theme="dark"
          mode="horizontal"
          selectedKeys={[selectedKey]}
          items={[
            { key: '/', label: <Link to="/">首页</Link> },
            { key: '/users', label: <Link to="/users">用户查询</Link> },
            { key: '/queue', label: <Link to="/queue">调度队列</Link> },
          ]}
        />
      </Header>
      <Content style={{ padding: 24, maxWidth: 1280, margin: '0 auto', width: '100%' }}>
        <Suspense
          fallback={
            <div style={{ paddingTop: 40, textAlign: 'center' }}>
              <Spin />
            </div>
          }
        >
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/users" element={<UsersPage />} />
            <Route path="/users/:id" element={<UserDetailPage />} />
            <Route path="/items/:type/:id" element={<ItemDetailPage />} />
            <Route path="/queue" element={<QueuePage />} />

            <Route path="/explore" element={<HomePage />} />
            <Route path="/explore/users" element={<UsersPage />} />
            <Route path="/explore/queue" element={<QueuePage />} />
            <Route path="/user/:id" element={<UserDetailPage />} />
            <Route path="/item/:type/:id" element={<ItemDetailPage />} />
          </Routes>
        </Suspense>
      </Content>
    </Layout>
  )
}
