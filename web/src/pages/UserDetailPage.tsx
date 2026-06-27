import { Alert, App, Avatar, Button, Card, Descriptions, Empty, Segmented, Space, Spin, Table, Tabs, Tooltip, Typography } from 'antd'
import { ExportOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import { useEffect, useMemo, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { getUser, getUserComments } from '../api/client'
import StatCard from '../components/StatCard'
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
    if (!userId) {
      setError('参数错误')
      setLoading(false)
      return
    }

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

  const columns: ColumnsType<UserComment> = useMemo(
    () => [
      {
        title: '条目',
        dataIndex: 'item',
        width: 480,
        render: (_, row) => {
          const itemId = row.item?.douban_id || row.item?.DoubanId
          const title = row.item?.title || row.item?.Title || '-'
          const thumbnail = row.item?.thumbnail || row.item?.Thumbnail
          return (
            <Space size={10} align="start">
              <Avatar shape="square" src={thumbnail as string | undefined} size={52} />
              {itemId ? (
                <Link
                  to={`/items/${activeType}/${itemId}`}
                  style={{ maxWidth: 380, display: 'inline-block', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}
                  title={title}
                >
                  {title}
                </Link>
              ) : (
                <Text>{title}</Text>
              )}
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
        render: (value) => <Paragraph style={{ marginBottom: 0, whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{value || '-'}</Paragraph>,
      },
      { title: '标注日期', dataIndex: 'mark_date', width: 120 },
    ],
    [activeType],
  )

  const userProfileUrl = user ? `https://www.douban.com/people/${user.domain || user.id}/` : '#'
  const currentList = commentMap[commentKey(activeType, activeAction)] || []

  return (
    <Space direction="vertical" size={12} style={{ width: '100%' }}>
      {error ? <Alert type="error" message={error} showIcon /> : null}

      <Spin spinning={loading}>
        {!loading && !error && !user ? <Empty description="未找到用户" /> : null}

        {user ? (
          <>
            <Card size="small">
              <Space align="start" style={{ width: '100%', justifyContent: 'space-between' }}>
                <Space align="start">
                  <Avatar src={user.thumbnail} size={92} />
                  <Descriptions title={user.name} column={1} size="small">
                    <Descriptions.Item label="ID">{user.id}</Descriptions.Item>
                    <Descriptions.Item label="Domain">{user.domain || '-'}</Descriptions.Item>
                    <Descriptions.Item label="最近发布">{formatUnix(user.publish_at)}</Descriptions.Item>
                    <Descriptions.Item label="最近同步">{formatUnix(user.sync_at)}</Descriptions.Item>
                    <Descriptions.Item label="最近检查">{formatUnix(user.check_at)}</Descriptions.Item>
                  </Descriptions>
                </Space>
                <Tooltip title="跳转豆瓣主页">
                  <Button type="text" shape="circle" icon={<ExportOutlined />} onClick={() => window.open(userProfileUrl, '_blank', 'noopener,noreferrer')} />
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
                items={mediaOptions.map((x) => ({ key: x.key, label: `${x.label} (${typeTotalCount(user, x.key)})` }))}
              />

              <Space wrap style={{ marginBottom: 12 }}>
                <StatCard title={actionLabelMap[activeType].collect} value={typeActionCount(user, activeType, 'collect')} />
                <StatCard title={actionLabelMap[activeType].do} value={typeActionCount(user, activeType, 'do')} />
                <StatCard title={actionLabelMap[activeType].wish} value={typeActionCount(user, activeType, 'wish')} />
              </Space>

              <Segmented
                block
                value={activeAction}
                onChange={(v) => setActiveAction(v as ActionType)}
                options={(['collect', 'do', 'wish'] as ActionType[]).map((action) => ({
                  label: `${actionLabelMap[activeType][action]} (${typeActionCount(user, activeType, action)})`,
                  value: action,
                }))}
              />

              <div style={{ marginTop: 12 }}>
                <Spin spinning={commentLoading}>
                  <Table
                    rowKey={(_, idx) => `${activeType}_${activeAction}_${idx}`}
                    columns={columns}
                    dataSource={currentList}
                    locale={{ emptyText: <Empty description="暂无数据" /> }}
                    size="small"
                    pagination={{ defaultPageSize: 20, showSizeChanger: true, pageSizeOptions: [20, 50, 100] }}
                    scroll={{ x: 'max-content' }}
                  />
                </Spin>
              </div>
            </Card>
          </>
        ) : null}
      </Spin>
    </Space>
  )
}

function typeActionCount(user: UserVO, type: MediaType, action: ActionType): number {
  if (type === 'book' && action === 'wish') return user.book_wish
  if (type === 'book' && action === 'do') return user.book_do
  if (type === 'book' && action === 'collect') return user.book_collect

  if (type === 'movie' && action === 'wish') return user.movie_wish
  if (type === 'movie' && action === 'do') return user.movie_do
  if (type === 'movie' && action === 'collect') return user.movie_collect

  if (type === 'game' && action === 'wish') return user.game_wish
  if (type === 'game' && action === 'do') return user.game_do
  if (type === 'game' && action === 'collect') return user.game_collect

  if (type === 'song' && action === 'wish') return user.song_wish
  if (type === 'song' && action === 'do') return user.song_do
  return user.song_collect
}

function typeTotalCount(user: UserVO, type: MediaType): number {
  return typeActionCount(user, type, 'wish') + typeActionCount(user, type, 'do') + typeActionCount(user, type, 'collect')
}

function formatUnix(ts: number) {
  if (!ts) return '暂无'
  return new Date(ts * 1000).toLocaleString('zh-CN', { hour12: false })
}
