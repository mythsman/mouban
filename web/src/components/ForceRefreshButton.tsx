import { ReloadOutlined } from '@ant-design/icons'
import { Button, Tooltip } from 'antd'

interface ForceRefreshButtonProps {
  onClick: () => void
  loading?: boolean
  tooltip?: string
}

export default function ForceRefreshButton({ onClick, loading, tooltip = '强制更新' }: ForceRefreshButtonProps) {
  return (
    <Tooltip title={tooltip}>
      <Button
        type="text"
        shape="circle"
        className="force-refresh-btn"
        icon={<ReloadOutlined />}
        loading={loading}
        onClick={onClick}
      />
    </Tooltip>
  )
}
