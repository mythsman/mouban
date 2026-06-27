import { Card, Progress, Space, Typography } from 'antd'

const { Text, Title } = Typography

interface StatCardProps {
  title: string
  value: number | string
  subtitle?: string
  percent?: number
}

export default function StatCard({ title, value, subtitle, percent }: StatCardProps) {
  return (
    <Card size="small" className="stat-card">
      <Space direction="vertical" size={4} style={{ width: '100%' }}>
        <Text type="secondary">{title}</Text>
        <Title level={4} style={{ margin: 0 }}>
          {value}
        </Title>
        {subtitle ? <Text type="secondary">{subtitle}</Text> : null}
        {typeof percent === 'number' ? <Progress percent={Math.min(100, Math.max(0, percent))} size="small" /> : null}
      </Space>
    </Card>
  )
}
