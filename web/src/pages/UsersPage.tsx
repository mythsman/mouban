import { Alert, App, Avatar, Button, Card, Col, Empty, Form, Input, List, Row, Space, Spin, Typography } from 'antd'
import { UserOutlined } from '@ant-design/icons'
import { useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { resolveUsers } from '../api/client'
import type { UserVO } from '../types/api'

const { Title, Text } = Typography

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
                <Card
                  size="small"
                  styles={{ body: { padding: 12 } }}
                  title={
                    <Space>
                      <Avatar src={u.thumbnail} icon={<UserOutlined />} size={36} />
                      <div>
                        <div style={{ lineHeight: 1.2 }}>{u.name}</div>
                        <Text type="secondary" style={{ fontSize: 12 }}>
                          ID: {u.id}
                        </Text>
                      </div>
                    </Space>
                  }
                  extra={
                    <Link to={`/users/${u.id}`}>
                      <Button type="link" size="small" style={{ padding: 0 }}>
                        详情
                      </Button>
                    </Link>
                  }
                >
                  <Row gutter={[8, 8]}>
                    <Col span={12}><Text style={{ fontSize: 12 }}>图书 想读 {u.book_wish}</Text></Col>
                    <Col span={12}><Text style={{ fontSize: 12 }}>图书 在读 {u.book_do}</Text></Col>
                    <Col span={12}><Text style={{ fontSize: 12 }}>图书 读过 {u.book_collect}</Text></Col>
                    <Col span={12}><Text style={{ fontSize: 12 }}>电影 想看 {u.movie_wish}</Text></Col>
                    <Col span={12}><Text style={{ fontSize: 12 }}>电影 在看 {u.movie_do}</Text></Col>
                    <Col span={12}><Text style={{ fontSize: 12 }}>电影 看过 {u.movie_collect}</Text></Col>
                    <Col span={12}><Text style={{ fontSize: 12 }}>游戏 想玩 {u.game_wish}</Text></Col>
                    <Col span={12}><Text style={{ fontSize: 12 }}>游戏 在玩 {u.game_do}</Text></Col>
                    <Col span={12}><Text style={{ fontSize: 12 }}>游戏 玩过 {u.game_collect}</Text></Col>
                    <Col span={12}><Text style={{ fontSize: 12 }}>音乐 想听 {u.song_wish}</Text></Col>
                    <Col span={12}><Text style={{ fontSize: 12 }}>音乐 在听 {u.song_do}</Text></Col>
                    <Col span={12}><Text style={{ fontSize: 12 }}>音乐 听过 {u.song_collect}</Text></Col>
                  </Row>
                </Card>
              </List.Item>
            )}
          />
        )}
      </Spin>
    </Space>
  )
}
