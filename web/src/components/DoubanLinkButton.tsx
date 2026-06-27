import { Button, Tooltip } from 'antd'

interface DoubanLinkButtonProps {
  url?: string
  tooltip?: string
}

export default function DoubanLinkButton({ url, tooltip = '跳转豆瓣' }: DoubanLinkButtonProps) {
  const disabled = !url || url === '#'

  return (
    <Tooltip title={tooltip}>
      <Button
        type="text"
        shape="circle"
        className="douban-link-btn"
        disabled={disabled}
        onClick={() => {
          if (!disabled) {
            window.open(url, '_blank', 'noopener,noreferrer')
          }
        }}
      >
        <span className="douban-link-btn__text">豆</span>
      </Button>
    </Tooltip>
  )
}
