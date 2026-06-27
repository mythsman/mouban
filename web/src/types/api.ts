export interface ApiResponse<T> {
  success: boolean
  msg?: string
  result: T
}

export interface UserVO {
  id: number
  domain: string
  name: string
  thumbnail: string
  book_wish: number
  book_do: number
  book_collect: number
  game_wish: number
  game_do: number
  game_collect: number
  movie_wish: number
  movie_do: number
  movie_collect: number
  song_wish: number
  song_do: number
  song_collect: number
  publish_at: number
  sync_at: number
  check_at: number
}

export interface ResolveUserResult {
  keyword: string
  users: UserVO[]
}

export interface UserCommentItem {
  DoubanId?: number
  douban_id?: number
  Title?: string
  title?: string
  Thumbnail?: string
  thumbnail?: string
}

export interface UserComment {
  item: UserCommentItem
  rate: number
  label: string
  comment: string
  mark_date: string
}

export interface UserCommentResult {
  user: UserVO
  comment: UserComment[]
}

export interface ItemRating {
  total: number
  rating: number
  star5: number
  star4: number
  star3: number
  star2: number
  star1: number
}

export interface ItemDetailResult {
  type: 'book' | 'movie' | 'game' | 'song'
  type_name: string
  item_id: number
  douban_url: string
  crawled_at_text: string
  data_updated_text: string
  rating?: ItemRating
  book?: Record<string, unknown> | null
  movie?: Record<string, unknown> | null
  game?: Record<string, unknown> | null
  song?: Record<string, unknown> | null
}

export interface QueuePool {
  pool: string
  pool_label: string
  concurrency: number
  running: number
  utilization: number
}

export interface QueueType {
  type_code: number
  type_label: string
  to_crawl: number
  crawling: number
  can_crawl: number
  unready: number
  invalid: number
  oldest_wait_seconds: number
}

export interface QueueRunningTask {
  type_label: string
  douban_id: number
  title: string
  detail_url: string
  status: string
  updated_at_text: string
  running_for_seconds: number
}

export interface QueueCompletedTask {
  type_label: string
  douban_id: number
  title: string
  status: string
  result: string
  detail_url: string
  updated_at_text: string
}

export interface QueueOverviewResult {
  pools: QueuePool[]
  types: QueueType[]
  running: QueueRunningTask[]
  completed: QueueCompletedTask[]
}
