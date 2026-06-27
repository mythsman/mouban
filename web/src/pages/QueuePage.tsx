import { Card, Progress, Space, Table } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { getQueueOverview } from '../api/client'
import type { QueueCompletedTask, QueueOverviewResult, QueueRunningTask, QueueType } from '../types/api'

export default function QueuePage() {
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState<QueueOverviewResult | null>(null)

  useEffect(() => {
    let mounted = true

    async function load() {
      try {
        const res = await getQueueOverview()
        if (mounted) setData(res)
      } finally {
        if (mounted) setLoading(false)
      }
    }

    load()
    const timer = setInterval(load, 20000)
    return () => {
      mounted = false
      clearInterval(timer)
    }
  }, [])

  const typeColumns: ColumnsType<QueueType> = [
    { title: '类型', dataIndex: 'type_label', width: 120 },
    { title: 'to_crawl', dataIndex: 'to_crawl', width: 110 },
    { title: 'crawling', dataIndex: 'crawling', width: 110 },
    { title: 'can_crawl', dataIndex: 'can_crawl', width: 120 },
    { title: 'unready', dataIndex: 'unready', width: 110 },
    { title: 'invalid', dataIndex: 'invalid', width: 110 },
    { title: '最老等待(秒)', dataIndex: 'oldest_wait_seconds', width: 140 },
  ]

  const runningColumns: ColumnsType<QueueRunningTask> = [
    { title: '类型', dataIndex: 'type_label', width: 120 },
    { title: 'DoubanID', dataIndex: 'douban_id', width: 140 },
    { title: '状态', dataIndex: 'status', width: 120 },
    { title: '开始时间', dataIndex: 'updated_at_text', width: 180 },
    { title: '运行时长(秒)', dataIndex: 'running_for_seconds', width: 140 },
  ]

  const completedColumns: ColumnsType<QueueCompletedTask> = [
    { title: '完成时间', dataIndex: 'updated_at_text', width: 180 },
    { title: '类型', dataIndex: 'type_label', width: 120 },
    {
      title: 'DoubanID',
      dataIndex: 'douban_id',
      width: 140,
      render: (_, row) => (row.detail_url && row.detail_url !== '#' ? <Link to={row.detail_url}>{row.douban_id}</Link> : row.douban_id),
    },
    { title: '状态', dataIndex: 'status', width: 120 },
    { title: '结果', dataIndex: 'result', width: 120 },
  ]

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <Space wrap>
        {(data?.pools || []).map((pool) => (
          <Card key={pool.pool} title={pool.pool_label} style={{ minWidth: 280 }}>
            <div>并发上限：{pool.concurrency}</div>
            <div>运行中：{pool.running}</div>
            <Progress percent={Math.min(100, Number((pool.utilization * 100).toFixed(1)))} status="active" />
          </Card>
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

      <Card title="最近完成（最新 50 条）">
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
