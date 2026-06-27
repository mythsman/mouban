import { Tag } from 'antd'

type StatusTone = 'success' | 'processing' | 'warning' | 'error' | 'default'

function normalizeTone(status: string): StatusTone {
  const s = status.toLowerCase()
  if (s.includes('ready') || s.includes('crawled')) return 'success'
  if (s.includes('crawling')) return 'processing'
  if (s.includes('to crawl') || s.includes('unready')) return 'warning'
  if (s.includes('invalid') || s.includes('fail')) return 'error'
  return 'default'
}

export default function StatusTag({ text }: { text: string }) {
  const tone = normalizeTone(text)
  return <Tag color={tone}>{text}</Tag>
}
