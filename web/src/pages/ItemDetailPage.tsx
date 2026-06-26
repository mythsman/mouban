import { Alert, Card, Descriptions, Empty, Space, Spin, Statistic, Typography } from 'antd'
import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { getItemDetail } from '../api/client'
import type { ItemDetailResult } from '../types/api'

const { Title, Text, Paragraph, Link } = Typography

type MediaType = 'book' | 'movie' | 'game' | 'song'

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

  const itemData = data?.book || data?.movie || data?.game || data?.song || null

  return (
    <Spin spinning={loading}>
      <Space direction="vertical" size={16} style={{ width: '100%' }}>
        {error ? <Alert type="error" message={error} /> : null}
        {!loading && !error && !data ? <Empty description="暂无数据" /> : null}

        {data ? (
          <>
            <Card>
              <Title level={3} style={{ marginTop: 0 }}>
                {String((itemData as Record<string, unknown>)?.Title || '-')}
              </Title>
              <Space split={<Text type="secondary">|</Text>}>
                <Text>{data.type_name}</Text>
                <Text>DoubanID: {data.item_id}</Text>
                <Text>最近爬取: {data.crawled_at_text}</Text>
                <Text>数据更新: {data.data_updated_text}</Text>
              </Space>
              <div style={{ marginTop: 8 }}>
                <Link href={data.douban_url} target="_blank">
                  在豆瓣打开
                </Link>
              </div>
            </Card>

            {data.rating ? (
              <Card title="评分信息">
                <Space wrap>
                  <Statistic title="均分" value={data.rating.rating} precision={1} />
                  <Statistic title="人数" value={data.rating.total} />
                  <Statistic title="五星" value={data.rating.star5} suffix="%" precision={1} />
                  <Statistic title="四星" value={data.rating.star4} suffix="%" precision={1} />
                  <Statistic title="三星" value={data.rating.star3} suffix="%" precision={1} />
                  <Statistic title="二星" value={data.rating.star2} suffix="%" precision={1} />
                  <Statistic title="一星" value={data.rating.star1} suffix="%" precision={1} />
                </Space>
              </Card>
            ) : null}

            {itemData ? (
              <Card title={`${data.type_name}详情`}>
                <Descriptions column={1} size="small" bordered>
                  {Object.entries(itemData).map(([key, value]) => {
                    if (!value || ['ID', 'CreatedAt', 'UpdatedAt'].includes(key)) return null
                    if (typeof value === 'string' && value.length > 200) {
                      return (
                        <Descriptions.Item key={key} label={key}>
                          <Paragraph style={{ marginBottom: 0, whiteSpace: 'pre-wrap' }}>{value}</Paragraph>
                        </Descriptions.Item>
                      )
                    }
                    return (
                      <Descriptions.Item key={key} label={key}>
                        {String(value)}
                      </Descriptions.Item>
                    )
                  })}
                </Descriptions>
              </Card>
            ) : null}
          </>
        ) : null}
      </Space>
    </Spin>
  )
}
