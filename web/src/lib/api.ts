export const BASE = import.meta.env.PROD ? '' : 'http://localhost:8080'
export const API_URL = `${BASE}/api`

// A generic fetch function
async function fetcher(url: string, options: RequestInit = {}) {
	const res = await fetch(`${API_URL}${url}`, {
		...options,
		headers: {
			'Content-Type': 'application/json',
			...options.headers,
		},
		credentials: 'include',
	})

	if (!res.ok) {
		throw new Error('An error occurred while fetching the data.')
	}

	// if res.status is 204, return null
	if (res.status === 204) {
		return null
	}

	return res.json()
}

// Auth APIs
export const login = (data: any) => fetcher(`/auth/login`, { method: 'POST', body: JSON.stringify(data) })
export const register = (data: any) => fetcher(`/auth/register`, { method: 'POST', body: JSON.stringify(data) })
export const linkTelegram = (token: string) =>
	fetcher(`/auth/telegram/link`, { method: 'POST', body: JSON.stringify({ token }) })

// URL APIs
export const getUrls = () => fetcher(`/url`)
export const createUrl = (data: { original_url: string; custom_path?: string }) =>
	fetcher(`/url`, { method: 'POST', body: JSON.stringify(data) })
export const deleteUrl = (id: number) => fetcher(`/url/${id}`, { method: 'DELETE' })
export const getUrlStats = (id: number) => fetcher(`/url/${id}/stats`)

// Admin APIs
export const adminGetUrls = () => fetcher(`/admin/urls`)
export const adminDeleteUrl = (id: number) => fetcher(`/admin/urls/${id}`, { method: 'DELETE' })
