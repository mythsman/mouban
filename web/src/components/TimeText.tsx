export default function TimeText({ value }: { value: string }) {
  return <span className="time-text">{value || '-'}</span>
}
