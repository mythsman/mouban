import { Button, Card, Flex, Typography } from 'antd'
import { Link } from 'react-router-dom'

const { Title, Paragraph } = Typography

export default function HomePage() {
  return (
    <Flex vertical gap={16}>
      <Card>
        <Title level={2} style={{ marginTop: 0 }}>
          mouban
        </Title>
        <Paragraph>豆瓣书 / 影 / 游 / 音数据查询服务</Paragraph>
        <Flex gap={8}>
          <Link to="/users">
            <Button type="primary">进入用户查询</Button>
          </Link>
          <Link to="/queue">
            <Button>查看调度队列</Button>
          </Link>
        </Flex>
      </Card>
    </Flex>
  )
}
