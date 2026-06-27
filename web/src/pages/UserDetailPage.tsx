import { App, Avatar, Button, Card, Descriptions, Empty, Space, Spin, Table, Tabs, Tooltip, Typography } from 'antd'
import { ExportOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useEffect, useMemo, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { getUser, getUserComments } from '../api/client'
import type { UserComment, UserVO } from '../types/api'

const { Text, Paragraph } = Typography

type MediaType = 'book' | 'movie' | 'game' | 'song'
type ActionType = 'wish' | 'do' | 'collect'

const mediaOptions: Array<{ key: MediaType; label: string }> = [
  { key: 'book', label: '图书' },
  { key: 'movie', label: '电影' },
  { key: 'game', label: '游戏' },
  { key: 'song', label: '音乐' },
]

const actionLabelMap: Record<MediaType, Record<ActionType, string>> = {
  book: { wish: '想读', do: '在读', collect: '读过' },
  movie: { wish: '想看', do: '在看', collect: '看过' },
  game: { wish: '想玩', do: '在玩', collect: '玩过' },
  song: { wish: '想听', do: '在听', collect: '听过' },
}

function commentKey(type: MediaType, action: ActionType) {
  return `${type}_${action}`
}

export default function UserDetailPage() {
  const { message } = App.useApp()
  const { id } = useParams()

  const userId = Number(id)

  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [user, setUser] = useState<UserVO | null>(null)
  const [activeType, setActiveType] = useState<MediaType>('book')
  const [activeAction, setActiveAction] = useState<ActionType>('collect')
  const [commentMap, setCommentMap] = useState<Record<string, UserComment[]>>({})
  const [commentLoading, setCommentLoading] = useState(false)

  useEffect(() => {
    if (!userId) return
    setLoading(true)
    setError('')
    getUser(userId)
      .then(setUser)
      .catch((e) => {
        setError(e instanceof Error ? e.message : '加载失败')
        message.error('用户加载失败')
      })
      .finally(() => setLoading(false))
  }, [message, userId])

  useEffect(() => {
    if (!userId || !user) return
    const key = commentKey(activeType, activeAction)
    if (commentMap[key]) return

    setCommentLoading(true)
    getUserComments(userId, activeType, activeAction)
      .then((res) => {
        setCommentMap((prev) => ({ ...prev, [key]: res.comment || [] }))
      })
      .catch((e) => {
        message.error(e instanceof Error ? e.message : '评论加载失败')
      })
      .finally(() => setCommentLoading(false))
  }, [activeAction, activeType, commentMap, message, user, userId])

  const comments = commentMap[commentKey(activeType, activeAction)] || []

  const columns: ColumnsType<UserComment> = useMemo(
    () => [
      {
        title: '条目',
        dataIndex: 'item',
        render: (_, row) => {
          const itemId = row.item?.douban_id || row.item?.DoubanId
          const title = row.item?.title || row.item?.Title || '-'
          const thumbnail = row.item?.thumbnail || row.item?.Thumbnail
          return (
            <Space size={8} align="start">
              <Avatar shape="square" src={thumbnail as string | undefined} size={36} />
              {itemId ? <Link to={`/items/${activeType}/${itemId}`}>{title}</Link> : <Text>{title}</Text>}
            </Space>
          )
        },
      },
      {
        title: '评分',
        dataIndex: 'rate',
        width: 88,
        render: (rate) => {
          const n = Number(rate) || 0
          return n > 0 ? `${'★'.repeat(Math.min(n, 5))}` : '-'
        },
      },
      { title: '标签', dataIndex: 'label', width: 140 },
      {
        title: '评论',
        dataIndex: 'comment',
        render: (value) => (
          <Paragraph style={{ marginBottom: 0, whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{value || '-'}</Paragraph>
        ),
      },
      { title: '标注日期', dataIndex: 'mark_date', width: 120 },
    ],
    [activeType],
  )

  const userProfileUrl = user ? `https://www.douban.com/people/${user.domain || user.id}/` : '#'

  return (
    <Space direction="vertical" size={12} style={{ width: '100%' }}>
      <Spin spinning={loading}>
        {error ? <Card size="small"><Text type="danger">{error}</Text></Card> : null}
        {!loading && !error && !user ? <Empty description="未找到用户" /> : null}

        {user ? (
          <>
            <Card size="small">
              <Space align="start" style={{ width: '100%', justifyContent: 'space-between' }}>
                <Space align="start">
                  <Avatar src={user.thumbnail} size={72} />
                  <Descriptions title={user.name} column={1} size="small">
                    <Descriptions.Item label="ID">{user.id}</Descriptions.Item>
                    <Descriptions.Item label="Domain">{user.domain || '-'}</Descriptions.Item>
                    <Descriptions.Item label="最近发布">{formatUnix(user.publish_at)}</Descriptions.Item>
                    <Descriptions.Item label="最近同步">{formatUnix(user.sync_at)}</Descriptions.Item>
                    <Descriptions.Item label="最近检查">{formatUnix(user.check_at)}</Descriptions.Item>
                  </Descriptions>
                </Space>
                <Tooltip title="跳转豆瓣主页">
                  <Button
                    type="text"
                    shape="circle"
                    icon={<ExportOutlined />}
                    onClick={() => window.open(userProfileUrl, '_blank', 'noopener,noreferrer')}
                  />
                </Tooltip>
              </Space>
            </Card>

            <Card size="small">
              <Tabs
                activeKey={activeType}
                onChange={(k) => {
                  setActiveType(k as MediaType)
                  setActiveAction('collect')
                }}
                items={mediaOptions.map((x) => ({ key: x.key, label: x.label }))}
              />
              <Tabs
                activeKey={activeAction}
                onChange={(k) => setActiveAction(k as ActionType)}
                items={(['collect', 'do', 'wish'] as ActionType[]).map((key) => ({
                  key,
                  label: actionLabelMap[activeType][key],
                }))}
              />
              <Spin spinning={commentLoading}>
                <Table
                  rowKey={(_, idx) => `${activeType}_${activeAction}_${idx}`}
                  columns={columns}
                  dataSource={comments}
                  locale={{ emptyText: <Empty description="暂无数据" /> }}
                  size="small"
                  pagination={{ defaultPageSize: 20, showSizeChanger: true, pageSizeOptions: [20, 50, 100] }}
                />
              </Spin>
            </Card>
          </>
        ) : null}
      </Spin>
    </Space>
  )
}

function formatUnix(ts: number) {
  if (!ts) return '暂无'
  return new Date(ts * 1000).toLocaleString('zh-CN', { hour12: false })
}
