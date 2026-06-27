import { Button, Card, Descriptions, Empty, Image, Space, Spin, Statistic, Tooltip, Typography } from 'antd'
import { ExportOutlined, PictureOutlined } from '@ant-design/icons'
import { useEffect, useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'
import { getItemDetail } from '../api/client'
import type { ItemDetailResult } from '../types/api'

const { Title, Paragraph, Text } = Typography

type MediaType = 'book' | 'movie' | 'game' | 'song'

type FieldDef = {
  key: string
  label: string
}

const fieldMap: Record<MediaType, FieldDef[]> = {
  book: [
    { key: 'Subtitle', label: '副标题' },
    { key: 'Author', label: '作者' },
    { key: 'Translator', label: '译者' },
    { key: 'Press', label: '出版社' },
    { key: 'PublishDate', label: '出版日期' },
    { key: 'ISBN', label: 'ISBN' },
    { key: 'Producer', label: '出品方' },
  ],
  movie: [
    { key: 'Director', label: '导演' },
    { key: 'Writer', label: '编剧' },
    { key: 'Actor', label: '主演' },
    { key: 'Style', label: '类型' },
    { key: 'Country', label: '制片国家/地区' },
    { key: 'Language', label: '语言' },
    { key: 'PublishDate', label: '上映日期' },
    { key: 'Alias', label: '又名' },
  ],
  game: [
    { key: 'Platform', label: '平台' },
    { key: 'Genre', label: '类型' },
    { key: 'Developer', label: '开发商' },
    { key: 'Publisher', label: '发行商' },
    { key: 'PublishDate', label: '发行日期' },
    { key: 'Alias', label: '别名' },
  ],
  song: [
    { key: 'Musician', label: '表演者' },
    { key: 'AlbumType', label: '专辑类型' },
    { key: 'Genre', label: '流派' },
    { key: 'Media', label: '介质' },
    { key: 'Publisher', label: '出版者' },
    { key: 'PublishDate', label: '发行日期' },
    { key: 'Alias', label: '又名' },
  ],
}

const introMap: Record<MediaType, FieldDef[]> = {
  book: [
    { key: 'BookIntro', label: '内容简介' },
    { key: 'AuthorIntro', label: '作者简介' },
  ],
  movie: [{ key: 'Intro', label: '简介' }],
  game: [{ key: 'Intro', label: '简介' }],
  song: [
    { key: 'Intro', label: '简介' },
    { key: 'TrackList', label: '曲目列表' },
  ],
}

export default function ItemDetailPage() {
  const { type, id } = useParams()
  const itemId = Number(id)
  const itemType = (type || '') as MediaType

  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [data, setData] = useState<ItemDetailResult | null>(null)

  useEffect(() => {
    if (!itemId || !['book', 'movie', 'game', 'song'].includes(itemType)) {
      setError('参数错误')
      setLoading(false)
      return
    }

    setLoading(true)
    setError('')
    getItemDetail(itemType, itemId)
      .then(setData)
      .catch((e) => setError(e instanceof Error ? e.message : '加载失败'))
      .finally(() => setLoading(false))
  }, [itemId, itemType])

  const itemData = useMemo(() => {
    return (data?.book || data?.movie || data?.game || data?.song || null) as Record<string, unknown> | null
  }, [data])

  const title = String(itemData?.Title || '-')
  const thumbnail = (itemData?.Thumbnail as string) || ''

  return (
    <Spin spinning={loading}>
      <Space direction="vertical" size={12} style={{ width: '100%' }}>
        {error ? (
          <Card size="small">
            <Text type="danger">{error}</Text>
          </Card>
        ) : null}

        {!loading && !error && !data ? <Empty description="暂无数据" /> : null}

        {data && itemData ? (
          <>
            <Card size="small">
              <Space align="start" style={{ width: '100%', justifyContent: 'space-between' }}>
                <Space align="start" size={16}>
                  {thumbnail ? (
                    <Image src={thumbnail} width={160} preview={false} style={{ borderRadius: 6 }} />
                  ) : (
                    <div
                      style={{
                        width: 160,
                        height: 220,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        border: '1px solid #f0f0f0',
                        borderRadius: 6,
                      }}
                    >
                      <PictureOutlined style={{ fontSize: 24, color: '#bbb' }} />
                    </div>
                  )}
                  <div>
                    <Title level={3} style={{ marginTop: 0, marginBottom: 4, fontWeight: 600 }}>
                      {title}
                    </Title>
                    <Descriptions column={1} size="small" style={{ marginBottom: 8 }} items={[
                      { key: 'type', label: '类型', children: data.type_name },
                      { key: 'id', label: '豆瓣ID', children: data.item_id },
                      { key: 'crawl', label: '最近爬取', children: data.crawled_at_text },
                      { key: 'update', label: '数据更新', children: data.data_updated_text },
                    ]} />
                    <Descriptions
                      column={1}
                      size="small"
                      styles={{ label: { width: 96, color: 'rgba(0,0,0,0.65)' } }}
                      items={fieldMap[itemType]
                        .map((f) => ({ key: f.key, label: f.label, children: itemData[f.key] ? String(itemData[f.key]) : '-' }))}
                    />
                  </div>
                </Space>
                <Tooltip title="跳转豆瓣页面">
                  <Button
                    type="text"
                    shape="circle"
                    icon={<ExportOutlined />}
                    onClick={() => window.open(data.douban_url, '_blank', 'noopener,noreferrer')}
                  />
                </Tooltip>
              </Space>
            </Card>

            {data.rating ? (
              <Card size="small" title="评分信息">
                <Space wrap>
                  <Statistic title="平均分" value={data.rating.rating} precision={1} />
                  <Statistic title="评分人数" value={data.rating.total} />
                  <Statistic title="五星" value={data.rating.star5} suffix="%" precision={1} />
                  <Statistic title="四星" value={data.rating.star4} suffix="%" precision={1} />
                  <Statistic title="三星" value={data.rating.star3} suffix="%" precision={1} />
                  <Statistic title="二星" value={data.rating.star2} suffix="%" precision={1} />
                  <Statistic title="一星" value={data.rating.star1} suffix="%" precision={1} />
                </Space>
              </Card>
            ) : null}

            {introMap[itemType].map((f) => {
              const content = itemData[f.key]
              if (!content) return null
              return (
                <Card key={f.key} size="small" title={f.label}>
                  <Paragraph style={{ marginBottom: 0, whiteSpace: 'pre-wrap' }}>{String(content)}</Paragraph>
                </Card>
              )
            })}
          </>
        ) : null}
      </Space>
    </Spin>
  )
}
