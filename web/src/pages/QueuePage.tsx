import { Alert, Card, Space, Table, Typography } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { getQueueOverview } from '../api/client'
import StatusTag from '../components/StatusTag'
import StatCard from '../components/StatCard'
import TimeText from '../components/TimeText'
import type { QueueCompletedTask, QueueOverviewResult, QueueRunningTask, QueueType } from '../types/api'

const { Text } = Typography

export default function QueuePage() {
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState<QueueOverviewResult | null>(null)
  const [error, setError] = useState('')
  const [updatedAt, setUpdatedAt] = useState<string>('')

  async function load() {
    try {
      setError('')
      const res = await getQueueOverview()
      setData(res)
      setUpdatedAt(new Date().toLocaleString('zh-CN', { hour12: false }))
    } catch (e) {
      setError(e instanceof Error ? e.message : '队列数据加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  const summary = useMemo(() => {
    const types = data?.types || []
    return types.reduce(
      (acc, row) => {
        acc.toCrawl += row.to_crawl
        acc.crawling += row.crawling
        acc.oldestWait = Math.max(acc.oldestWait, row.oldest_wait_seconds)
        return acc
      },
      { toCrawl: 0, crawling: 0, oldestWait: 0 },
    )
  }, [data?.types])

  const typeColumns: ColumnsType<QueueType> = [
    { title: '类型', dataIndex: 'type_label', width: 120 },
    { title: '待抓取', dataIndex: 'to_crawl', width: 110 },
    { title: '执行中', dataIndex: 'crawling', width: 110 },
    { title: '最老等待(秒)', dataIndex: 'oldest_wait_seconds', width: 140 },
  ]

  const runningColumns: ColumnsType<QueueRunningTask> = [
    { title: '类型', dataIndex: 'type_label', width: 120 },
    {
      title: 'DoubanID',
      dataIndex: 'douban_id',
      width: 140,
      render: (_, row) => (row.detail_url && row.detail_url !== '#' ? <Link to={row.detail_url}>{row.douban_id}</Link> : row.douban_id),
    },
    {
      title: '名称',
      dataIndex: 'title',
      width: 260,
      render: (_, row) => {
        const text = row.title || '-'
        return row.detail_url && row.detail_url !== '#' ? <Link to={row.detail_url}>{text}</Link> : text
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 120,
      render: (v) => <StatusTag text={v} />,
    },
    {
      title: '开始时间',
      dataIndex: 'updated_at_text',
      width: 180,
      render: (v) => <TimeText value={v} />,
    },
    { title: '运行时长(秒)', dataIndex: 'running_for_seconds', width: 140 },
  ]

  const completedColumns: ColumnsType<QueueCompletedTask> = [
    {
      title: '完成时间',
      dataIndex: 'updated_at_text',
      width: 180,
      render: (v) => <TimeText value={v} />,
    },
    { title: '类型', dataIndex: 'type_label', width: 120 },
    {
      title: 'DoubanID',
      dataIndex: 'douban_id',
      width: 140,
      render: (_, row) => (row.detail_url && row.detail_url !== '#' ? <Link to={row.detail_url}>{row.douban_id}</Link> : row.douban_id),
    },
    {
      title: '名称',
      dataIndex: 'title',
      width: 260,
      render: (_, row) => {
        const text = row.title || '-'
        return row.detail_url && row.detail_url !== '#' ? <Link to={row.detail_url}>{text}</Link> : text
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 120,
      render: (v) => <StatusTag text={v} />,
    },
    {
      title: '结果',
      dataIndex: 'result',
      width: 120,
      render: (v) => <StatusTag text={v} />,
    },
  ]

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <Text type="secondary">最近刷新：{updatedAt || '-'}</Text>

      {error ? <Alert type="error" message={error} showIcon /> : null}

      <Space wrap>
        <StatCard title="待抓取总量" value={summary.toCrawl} />
        <StatCard title="执行中总量" value={summary.crawling} />
        <StatCard title="最老等待" value={`${summary.oldestWait}s`} />
        {(data?.pools || []).map((pool) => (
          <StatCard
            key={pool.pool}
            title={pool.pool_label}
            value={`${pool.running}/${pool.concurrency}`}
            subtitle={`占用率 ${(pool.utilization * 100).toFixed(1)}%`}
            percent={Number((pool.utilization * 100).toFixed(1))}
          />
        ))}
      </Space>

      <Card title="按类型统计（队列视角）">
        <Table
          className="nowrap-table"
          rowKey="type_code"
          loading={loading}
          columns={typeColumns}
          dataSource={data?.types || []}
          size="small"
          pagination={false}
          scroll={{ x: 'max-content' }}
        />
      </Card>

      <Card title="当前执行中（crawling）">
        <Table
          className="nowrap-table"
          rowKey={(row) => `${row.type_label}_${row.douban_id}`}
          loading={loading}
          columns={runningColumns}
          dataSource={data?.running || []}
          size="small"
          pagination={{ defaultPageSize: 20, showSizeChanger: true }}
          scroll={{ x: 'max-content' }}
        />
      </Card>

      <Card title="最近完成">
        <Table
          className="nowrap-table"
          rowKey={(row) => `${row.type_label}_${row.douban_id}_${row.updated_at_text}`}
          loading={loading}
          columns={completedColumns}
          dataSource={data?.completed || []}
          size="small"
          pagination={{ defaultPageSize: 20, showSizeChanger: true }}
          scroll={{ x: 'max-content' }}
        />
      </Card>
    </Space>
  )
}
