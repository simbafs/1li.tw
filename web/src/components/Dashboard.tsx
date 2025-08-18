import { useState, useEffect } from 'react'
import { AddUrlForm } from './AddUrlForm'
import type { User } from '../hooks/useUser'
import { ListShortURL } from './ListShortURL'
import { toast } from 'react-toastify'

export function Dashboard() {
	const [user, setUser] = useState<User | null>(null)

	useEffect(() => {
		const userStr = localStorage.getItem('user')
		if (userStr) {
			try {
				setUser(JSON.parse(userStr))
			} catch (e) {
				toast.error(`Failed to parse user from localStorage: ${e}`)
				setUser(null)
				// Redirect to login if user data is corrupted
				localStorage.removeItem('user')
				window.location.href = '/login'
			}
		}
	}, [])

	if (!user) {
		return <span className="loading loading-spinner loading-lg" />
	}

	return (
		<div className="flex w-full flex-col gap-8">
			<AddUrlForm canCollapse />

			<h2 className="text-2xl font-bold">My URLs</h2>
			<ListShortURL user={user} />
		</div>
	)
}
