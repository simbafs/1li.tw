import { type Stats } from '../components/StatsCharts'

export const BASE = () => location.origin // TODO: is there any other solution?
export const API_URL = () => `${BASE()}/api`

export type URL = {
	ID: number
	ShortPath: string
	OriginalURL: string
	TotalClicks: number
	CreatedAt: string
	Username?: string // this will show in some url endpoints  // TODO: make this presistent
}

type Method = 'POST' | 'GET' | 'PUT' | 'DELETE'

// A generic fetch function
export async function api<T>(path: string, method: Method, body?: any) {
	const res = await fetch(`/api${path}`, {
		method,
		headers: {
			'Content-Type': 'application/json',
		},
		credentials: 'include',
		body: body ? JSON.stringify(body) : undefined,
	})

	const responseBody = await res.json().catch(() => ({}))

	if (!res.ok || (responseBody && responseBody.error)) {
		const errorMessage = (responseBody && responseBody.error) || 'An error occurred while fetching the data.'
		console.error(errorMessage)
		throw new Error(errorMessage)
	}

	return responseBody as T
}

// routes about authentication
export const register = (data: { username: string; password: string }) => api(`/auth/register`, 'POST', data)
export const login = (data: { username: string; password: string }) => api(`/auth/login`, 'POST', data)
export const logout = () => api(`/auth/logout`, 'POST')
export const linkTelegram = (token: string) => api(`/auth/telegram/link`, 'POST', { token })

// route about user itself
export const getMe = () => api('/me', 'GET')

// routes about a short URL
export const createUrl = (original_url: string, custom_path?: string) =>
	api<URL>(`/url`, 'POST', { original_url, custom_path })
export const getUrls = () => api<URL[]>(`/url`, 'GET')
export const deleteUrl = (id: number) => api(`/url/${id}`, 'DELETE')
export const getUrlStats = (id: number) => api<Stats>(`/url/${id}/stats`, 'GET')

// routes about managge users
export const listUsers = () => api('/user', 'GET')
export const updateUserPermission = (id: number, permission: number) =>
	api(`/user/${id}/permission`, 'PUT', { permission })
export const deleteUser = (id: number) => api(`/user/${id}`, 'DELETE')

// routes about admin
export const adminGetUrls = () => api<URL[]>(`/admin/url`, 'GET')
