import axios from 'axios'
import type {
  ApiResponse,
  ItemDetailResult,
  QueueOverviewResult,
  ResolveUserResult,
  UserCommentResult,
  UserVO,
} from '../types/api'

const http = axios.create({
  timeout: 15000,
})

function unwrap<T>(res: ApiResponse<T>): T {
  if (!res.success) {
    throw new Error(res.msg || '请求失败')
  }
  return res.result
}

export async function resolveUsers(q: string) {
  const { data } = await http.get<ApiResponse<ResolveUserResult>>('/guest/resolve_user', { params: { q } })
  return unwrap(data)
}

export async function getUser(id: number) {
  const { data } = await http.get<ApiResponse<UserVO>>('/guest/check_user', { params: { id } })
  return unwrap(data)
}

export async function getUserComments(id: number, type: 'book' | 'movie' | 'game' | 'song', action: 'wish' | 'do' | 'collect') {
  const endpointMap = {
    book: '/guest/user_book',
    movie: '/guest/user_movie',
    game: '/guest/user_game',
    song: '/guest/user_song',
  }
  const { data } = await http.get<ApiResponse<UserCommentResult>>(endpointMap[type], {
    params: { id, action },
  })
  return unwrap(data)
}

export async function getItemDetail(type: 'book' | 'movie' | 'game' | 'song', id: number) {
  const { data } = await http.get<ApiResponse<ItemDetailResult>>('/guest/item_detail', {
    params: { type, id },
  })
  return unwrap(data)
}

export async function getQueueOverview() {
  const { data } = await http.get<ApiResponse<QueueOverviewResult>>('/explore/queue_overview', {
    params: { _: Date.now() },
  })
  return unwrap(data)
}

export async function refreshUser(id: number) {
  const { data } = await http.get<ApiResponse<unknown>>('/admin/refresh_user', { params: { id } })
  if (!data.success) {
    throw new Error(data.msg || '发起用户强制更新失败')
  }
}

export async function refreshItem(type: 'book' | 'movie' | 'game' | 'song', id: number) {
  const typeMap = {
    book: 1,
    movie: 2,
    game: 3,
    song: 4,
  } as const
  const { data } = await http.get<ApiResponse<unknown>>('/admin/refresh_item', {
    params: { type: typeMap[type], id },
  })
  if (!data.success) {
    throw new Error(data.msg || '发起条目强制更新失败')
  }
}
