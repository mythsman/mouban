import { Alert, App, Avatar, Button, Card, Empty, Form, Input, List, Select, Space, Tag, Typography } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import { useMemo, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { resolveUsers } from '../api/client'
import type { UserVO } from '../types/api'

const { Title, Text } = Typography

type SortType = 'relevance' | 'total_desc' | 'sync_desc'

export default function UsersPage() {
  const { message } = App.useApp()
  const [searchParams, setSearchParams] = useSearchParams()
  const [loading, setLoading] = useState(false)
  const [users, setUsers] = useState<UserVO[]>([])
  const [error, setError] = useState('')
  const [sortBy, setSortBy] = useState<SortType>('relevance')

  const q = searchParams.get('q') || ''

  async function onSearch(keyword: string) {
    const query = keyword.trim()
    if (!query) {
      message.warning('请输入关键词')
      return
    }
    setLoading(true)
    setError('')
    try {
      const result = await resolveUsers(query)
      setUsers(result.users || [])
      setSearchParams({ q: query })
    } catch (e) {
      setError(e instanceof Error ? e.message : '查询失败')
      setUsers([])
    } finally {
      setLoading(false)
    }
  }

  const sortedUsers = useMemo(() => {
    const list = [...users]
    if (sortBy === 'total_desc') {
      list.sort((a, b) => userTotalCount(b) - userTotalCount(a))
    }
    if (sortBy === 'sync_desc') {
      list.sort((a, b) => b.sync_at - a.sync_at)
    }
    return list
  }, [sortBy, users])

  return (
    <Space direction="vertical" size={12} style={{ width: '100%' }}>
      <Card size="small">
        <Title level={4} style={{ marginTop: 0, marginBottom: 8 }}>
          用户查询
        </Title>
        <Text type="secondary">支持 UID / domain / 用户名 精确匹配</Text>
        <Form layout="inline" onFinish={(v) => onSearch(v.q)} initialValues={{ q }} style={{ marginTop: 12 }}>
          <Form.Item name="q" rules={[{ required: true, message: '请输入关键词' }]} style={{ flex: 1, minWidth: 320 }}>
            <Input placeholder="例如: 1000001 / ahbei / 阿北" allowClear />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              查询
            </Button>
          </Form.Item>
        </Form>
      </Card>

      {error ? <Alert type="error" message={error} /> : null}

      <Card size="small">
        <Space style={{ width: '100%', justifyContent: 'space-between' }} wrap>
          <Text type="secondary">结果数：{sortedUsers.length}</Text>
          <Space>
            <Text type="secondary">排序</Text>
            <Select<SortType>
              value={sortBy}
              onChange={setSortBy}
              style={{ width: 170 }}
              options={[
                { label: '相关性（默认）', value: 'relevance' },
                { label: '总条目数降序', value: 'total_desc' },
                { label: '最近同步时间', value: 'sync_desc' },
              ]}
            />
          </Space>
        </Space>
      </Card>

      {sortedUsers.length === 0 ? (
        <Card size="small">{q ? <Empty description="没有匹配结果" /> : <Text type="secondary">请输入关键词开始查询</Text>}</Card>
      ) : (
        <List
          grid={{ gutter: 12, xs: 1, sm: 2, md: 2, lg: 3, xl: 3, xxl: 4 }}
          dataSource={sortedUsers}
          loading={loading}
          pagination={{ pageSize: 12, showSizeChanger: true, pageSizeOptions: [12, 24, 48], size: 'small' }}
          renderItem={(u) => (
            <List.Item>
              <UserCard user={u} />
            </List.Item>
          )}
        />
      )}
    </Space>
  )
}

function UserCard({ user }: { user: UserVO }) {
  const mediaRows = [
    { label: '图书', wish: user.book_wish, doing: user.book_do, collect: user.book_collect },
    { label: '电影', wish: user.movie_wish, doing: user.movie_do, collect: user.movie_collect },
    { label: '游戏', wish: user.game_wish, doing: user.game_do, collect: user.game_collect },
    { label: '音乐', wish: user.song_wish, doing: user.song_do, collect: user.song_collect },
  ]

  return (
    <Card
      size="small"
      styles={{ body: { padding: 12 } }}
      title={
        <Space size={10}>
          <Avatar src={user.thumbnail} icon={<UserOutlined />} size={52} />
          <div>
            <div style={{ lineHeight: 1.2, fontWeight: 600 }}>{user.name}</div>
            <Text type="secondary" style={{ fontSize: 12 }}>
              ID: {user.id}
            </Text>
          </div>
        </Space>
      }
      extra={
        <Link to={`/users/${user.id}`}>
          <Button type="link" size="small" style={{ padding: 0 }}>
            详情
          </Button>
        </Link>
      }
    >
      <Space wrap size={6} style={{ marginBottom: 8 }}>
        <Tag color="blue">总条目 {userTotalCount(user)}</Tag>
        {user.sync_at ? <Tag>同步 {new Date(user.sync_at * 1000).toLocaleDateString('zh-CN')}</Tag> : null}
      </Space>

      <div className="user-stat-list">
        {mediaRows.map((row) => {
          const total = row.wish + row.doing + row.collect
          return (
            <div key={row.label} className="user-stat-row">
              <Text strong>{row.label}</Text>
              <Text type="secondary">总 {total}</Text>
              <Text type="secondary">过 {row.collect}</Text>
              <Text type="secondary">在 {row.doing}</Text>
              <Text type="secondary">想 {row.wish}</Text>
            </div>
          )
        })}
      </div>
    </Card>
  )
}

function userTotalCount(user: UserVO) {
  return (
    user.book_wish +
    user.book_do +
    user.book_collect +
    user.movie_wish +
    user.movie_do +
    user.movie_collect +
    user.game_wish +
    user.game_do +
    user.game_collect +
    user.song_wish +
    user.song_do +
    user.song_collect
  )
}
