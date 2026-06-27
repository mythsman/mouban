import { Breadcrumb, Layout, Menu } from 'antd'
import { useMemo } from 'react'
import { Link, Navigate, Route, Routes, useLocation } from 'react-router-dom'
import UsersPage from './pages/UsersPage'
import UserDetailPage from './pages/UserDetailPage'
import ItemDetailPage from './pages/ItemDetailPage'
import QueuePage from './pages/QueuePage'

const { Header, Content } = Layout

export default function RootApp() {
  const location = useLocation()

  const selectedKey = useMemo(() => {
    if (location.pathname.startsWith('/queue') || location.pathname.startsWith('/explore/queue')) return '/queue'
    return '/users'
  }, [location.pathname])

  const breadcrumbItems = useMemo(() => {
    const path = location.pathname
    if (path.startsWith('/queue') || path.startsWith('/explore/queue')) {
      return [{ title: '调度队列' }]
    }
    if (path.startsWith('/users/') || path.startsWith('/user/')) {
      const id = path.split('/').filter(Boolean).pop()
      return [{ title: <Link to="/users">用户查询</Link> }, { title: `用户详情 ${id || ''}` }]
    }
    if (path.startsWith('/items/') || path.startsWith('/item/')) {
      const parts = path.split('/').filter(Boolean)
      const type = parts[1] || ''
      const id = parts[2] || ''
      const typeLabelMap: Record<string, string> = {
        book: '图书',
        movie: '电影',
        game: '游戏',
        song: '音乐',
      }
      return [
        { title: <Link to="/users">用户查询</Link> },
        { title: `${typeLabelMap[type] || '条目'}详情 ${id}` },
      ]
    }
    return [{ title: '用户查询' }]
  }, [location.pathname])

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ position: 'sticky', top: 0, zIndex: 10, width: '100%' }}>
        <Menu
          theme="dark"
          mode="horizontal"
          selectedKeys={[selectedKey]}
          items={[
            { key: '/users', label: <Link to="/users">用户查询</Link> },
            { key: '/queue', label: <Link to="/queue">调度队列</Link> },
          ]}
        />
      </Header>
      <Content style={{ padding: 20, maxWidth: 1360, margin: '0 auto', width: '100%' }}>
        <Breadcrumb items={breadcrumbItems} style={{ marginBottom: 12 }} />
        <Routes>
          <Route path="/" element={<Navigate to="/users" replace />} />
          <Route path="/users" element={<UsersPage />} />
          <Route path="/users/:id" element={<UserDetailPage />} />
          <Route path="/items/:type/:id" element={<ItemDetailPage />} />
          <Route path="/queue" element={<QueuePage />} />

          <Route path="/explore" element={<Navigate to="/users" replace />} />
          <Route path="/explore/users" element={<UsersPage />} />
          <Route path="/explore/queue" element={<QueuePage />} />
          <Route path="/user/:id" element={<UserDetailPage />} />
          <Route path="/item/:type/:id" element={<ItemDetailPage />} />
          <Route path="*" element={<Navigate to="/users" replace />} />
        </Routes>
      </Content>
    </Layout>
  )
}
