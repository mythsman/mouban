import { Alert, App, Avatar, Card, Descriptions, Empty, Rate, Space, Spin, Table, Tabs, Typography } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { useEffect, useMemo, useState } from 'react'
import { Link, useParams, useSearchParams } from 'react-router-dom'
import { getUser, getUserComments } from '../api/client'
import type { UserComment, UserVO } from '../types/api'

const { Text } = Typography

type MediaType = 'book' | 'movie' | 'game' | 'song'
type ActionType = 'wish' | 'do' | 'collect'

const mediaOptions: Array<{ key: MediaType; label: string }> = [
  { key: 'book', label: '图书' },
  { key: 'movie', label: '电影' },
  { key: 'game', label: '游戏' },
  { key: 'song', label: '音乐' },
]

const actionOptions: Array<{ key: ActionType; label: string }> = [
  { key: 'collect', label: '已完成' },
  { key: 'do', label: '进行中' },
  { key: 'wish', label: '想要' },
]

function commentKey(type: MediaType, action: ActionType) {
  return `${type}_${action}`
}

export default function UserDetailPage() {
  const { message } = App.useApp()
  const { id } = useParams()
  const [searchParams] = useSearchParams()

  const userId = Number(id)
  const keyword = searchParams.get('q') || ''

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
        width: 360,
        render: (_, row) => {
          const itemId = row.item?.douban_id || row.item?.DoubanId
          const title = row.item?.title || row.item?.Title || '-'
          const thumbnail = row.item?.thumbnail || row.item?.Thumbnail
          return (
            <Space>
              <Avatar shape="square" src={thumbnail as string | undefined} size={40} />
              {itemId ? <Link to={`/items/${activeType}/${itemId}`}>{title}</Link> : <Text>{title}</Text>}
            </Space>
          )
        },
      },
      {
        title: '评分',
        dataIndex: 'rate',
        width: 140,
        render: (rate) => <Rate disabled value={Number(rate) || 0} count={5} />,
      },
      { title: '标签', dataIndex: 'label', width: 180 },
      { title: '评论', dataIndex: 'comment', ellipsis: true },
      { title: '标注日期', dataIndex: 'mark_date', width: 140 },
    ],
    [activeType],
  )

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      {keyword ? (
        <Card>
          <Link to={`/users?q=${encodeURIComponent(keyword)}`}>返回搜索结果</Link>
        </Card>
      ) : null}

      <Spin spinning={loading}>
        {error ? <Alert type="error" message={error} /> : null}
        {!loading && !error && !user ? <Empty description="未找到用户" /> : null}

        {user ? (
          <>
            <Card>
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
            </Card>

            <Card>
              <Tabs
                activeKey={activeType}
                onChange={(k) => setActiveType(k as MediaType)}
                items={mediaOptions.map((x) => ({ key: x.key, label: x.label }))}
              />
              <Tabs
                activeKey={activeAction}
                onChange={(k) => setActiveAction(k as ActionType)}
                items={actionOptions.map((x) => ({ key: x.key, label: x.label }))}
              />
              <Spin spinning={commentLoading}>
                <Table
                  className="nowrap-table"
                  rowKey={(_, idx) => `${activeType}_${activeAction}_${idx}`}
                  columns={columns}
                  dataSource={comments}
                  locale={{ emptyText: <Empty description="暂无数据" /> }}
                  size="small"
                  pagination={{ defaultPageSize: 20, showSizeChanger: true, pageSizeOptions: [20, 50, 100] }}
                  scroll={{ x: 'max-content' }}
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
