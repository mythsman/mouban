import { Alert, App, Avatar, Button, Card, Empty, Form, Input, List, Space, Spin, Table, Typography } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useMemo, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { resolveUsers } from '../api/client'
import type { UserVO } from '../types/api'

const { Title, Text } = Typography

type StatRow = {
  key: string
  typeLabel: string
  wish: number
  doing: number
  collect: number
}

const statColumns: ColumnsType<StatRow> = [
  { title: '条目', dataIndex: 'typeLabel', width: 56 },
  { title: '想', dataIndex: 'wish', width: 52, align: 'right' },
  { title: '在', dataIndex: 'doing', width: 52, align: 'right' },
  { title: '过', dataIndex: 'collect', width: 52, align: 'right' },
]

export default function UsersPage() {
  const { message } = App.useApp()
  const [searchParams, setSearchParams] = useSearchParams()
  const [loading, setLoading] = useState(false)
  const [users, setUsers] = useState<UserVO[]>([])
  const [error, setError] = useState('')

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

  return (
    <Space direction="vertical" size={12} style={{ width: '100%' }}>
      <Card size="small">
        <Title level={4} style={{ marginTop: 0, marginBottom: 12 }}>
          用户查询
        </Title>
        <Form layout="inline" onFinish={(v) => onSearch(v.q)} initialValues={{ q }}>
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

      <Spin spinning={loading}>
        {users.length === 0 ? (
          <Card size="small">{q ? <Empty description="没有匹配结果" /> : <Text type="secondary">请输入关键词开始查询</Text>}</Card>
        ) : (
          <List
            grid={{ gutter: 12, xs: 1, sm: 2, md: 2, lg: 3, xl: 3, xxl: 4 }}
            dataSource={users}
            pagination={{ pageSize: 12, showSizeChanger: true, pageSizeOptions: [12, 24, 48], size: 'small' }}
            renderItem={(u) => (
              <List.Item>
                <UserCard user={u} />
              </List.Item>
            )}
          />
        )}
      </Spin>
    </Space>
  )
}

function UserCard({ user }: { user: UserVO }) {
  const rows = useMemo<StatRow[]>(
    () => [
      { key: 'book', typeLabel: '图书', wish: user.book_wish, doing: user.book_do, collect: user.book_collect },
      { key: 'movie', typeLabel: '电影', wish: user.movie_wish, doing: user.movie_do, collect: user.movie_collect },
      { key: 'game', typeLabel: '游戏', wish: user.game_wish, doing: user.game_do, collect: user.game_collect },
      { key: 'song', typeLabel: '音乐', wish: user.song_wish, doing: user.song_do, collect: user.song_collect },
    ],
    [user],
  )

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
      <Table<StatRow>
        rowKey="key"
        columns={statColumns}
        dataSource={rows}
        size="small"
        pagination={false}
        showHeader
      />
    </Card>
  )
}
