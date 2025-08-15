import { useEffect, useState } from 'react'

const stroageKey = 'user'

export type User = {
	username: string
	permissions: number
}

export function useUser() {
	const [user, setUser] = useState<User | null>(null)

	useEffect(() => {
		const user = localStorage.getItem(stroageKey)
		if (user) {
			try {
				setUser(JSON.parse(user))
			} catch (e) {
				console.error('Failed to parse user from localStorage', e)
				setUser(null)
			}
		}
	}, [])

	return user
}
