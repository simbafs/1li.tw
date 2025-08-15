import { useState, useEffect } from 'react'
import { getUrls, deleteUrl } from '../lib/api'
import { AddUrlForm } from './AddUrlForm'

interface Url {
	ID: number
	ShortPath: string
	OriginalURL: string
	TotalClicks: number
	CreatedAt: string
}

export function Dashboard() {
	const [urls, setUrls] = useState<Url[]>([])
	const [error, setError] = useState('')
	const [loading, setLoading] = useState(true)

	const fetchUrls = async () => {
		try {
			const data = await getUrls()
			setUrls(data)
		} catch (err: any) {
			if (err.status === 401) {
				// Unauthorized, redirect to login
				localStorage.removeItem('user')
				window.location.href = '/login'
			} else {
				setError('Failed to fetch URLs.')
			}
		} finally {
			setLoading(false)
		}
	}

	useEffect(() => {
		fetchUrls()
	}, [])

	const handleDelete = async (id: number) => {
		if (window.confirm('Are you sure you want to delete this URL?')) {
			try {
				await deleteUrl(id)
				setUrls(urls.filter(url => url.ID !== id))
			} catch (err: any) {
				setError('Failed to delete URL.')
			}
		}
	}

	if (loading) {
		return (
			<div className="text-center">
				<span className="loading loading-spinner loading-lg"></span>
			</div>
		)
	}

	return (
		<div>
			<div className="mb-8">
				<AddUrlForm />
			</div>

			<h2 className="mb-4 text-2xl font-bold">My URLs</h2>
			{error && <div className="alert alert-error">{error}</div>}
			<div className="overflow-x-auto">
				<table className="table w-full">
					<thead>
						<tr>
							<th>Short URL</th>
							<th>Original URL</th>
							<th>Clicks</th>
							<th>Created At</th>
							<th>Actions</th>
						</tr>
					</thead>
					<tbody>
						{urls.map(url => (
							<tr key={url.ID}>
								<td>
									<a
										href={`/${url.ShortPath}`}
										target="_blank"
										rel="noopener noreferrer"
										className="link link-primary"
									>
										{url.ShortPath}
									</a>
								</td>
								<td className="max-w-xs truncate">{url.OriginalURL}</td>
								<td>{url.TotalClicks}</td>
								<td>{new Date(url.CreatedAt).toLocaleDateString()}</td>
								<td>
									<a href={`/dashboard/stats?id=${url.ID}`} className="btn btn-ghost btn-sm">
										Stats
									</a>
									<button
										onClick={() => handleDelete(url.ID)}
										className="btn btn-ghost btn-sm text-error"
									>
										Delete
									</button>
								</td>
							</tr>
						))}
					</tbody>
				</table>
			</div>
		</div>
	)
}
