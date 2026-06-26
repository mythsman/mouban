import { Alert, App, Button, Card, Col, Descriptions, Empty, Form, Input, Row, Space, Spin, Typography } from 'antd'
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
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <Card>
        <Title level={3} style={{ marginTop: 0 }}>
          用户查询
        </Title>
        <Form layout="inline" onFinish={(v) => onSearch(v.q)} initialValues={{ q }}>
          <Form.Item name="q" rules={[{ required: true, message: '请输入关键词' }]} style={{ flex: 1 }}>
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
          <Card>{q ? <Empty description="没有匹配结果" /> : <Text type="secondary">请输入关键词开始查询</Text>}</Card>
        ) : (
          <Row gutter={[16, 16]}>
            {users.map((u) => (
              <Col xs={24} md={12} key={u.id}>
                <Card
                  title={u.name}
                  extra={
                    <Link to={`/users/${u.id}${q ? `?q=${encodeURIComponent(q)}` : ''}`}>
                      <Button type="link" style={{ padding: 0 }}>
                        查看详情
                      </Button>
                    </Link>
                  }
                >
                  <Descriptions column={1} size="small">
                    <Descriptions.Item label="ID">{u.id}</Descriptions.Item>
                    <Descriptions.Item label="Domain">{u.domain || '-'}</Descriptions.Item>
                    <Descriptions.Item label="图书">想读 {u.book_wish} / 在读 {u.book_do} / 读过 {u.book_collect}</Descriptions.Item>
                    <Descriptions.Item label="电影">想看 {u.movie_wish} / 在看 {u.movie_do} / 看过 {u.movie_collect}</Descriptions.Item>
                    <Descriptions.Item label="游戏">想玩 {u.game_wish} / 在玩 {u.game_do} / 玩过 {u.game_collect}</Descriptions.Item>
                    <Descriptions.Item label="音乐">想听 {u.song_wish} / 在听 {u.song_do} / 听过 {u.song_collect}</Descriptions.Item>
                  </Descriptions>
                </Card>
              </Col>
            ))}
          </Row>
        )}
      </Spin>
    </Space>
  )
}
