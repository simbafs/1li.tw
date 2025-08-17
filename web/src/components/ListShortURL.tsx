import { adminGetUrls, deleteUrl, getUrls, type URL } from '../lib/api'
import { canDeleteAny, canDeleteOwn, canViewAnyStats, canViewOwnStats } from '../lib/permissions'
import { formatShortPath } from '../lib/formatShortPath'
import { toast } from 'react-toastify'
import type { User } from '../hooks/useUser'
import useSWR from 'swr'
import { useEffect, useState } from 'react'

// format date to YYYY/MM/DD format
function formatDate(dateString: string): string {
	const date = new Date(dateString)
	const year = date.getFullYear()
	const month = String(date.getMonth() + 1).padStart(2, '0') // Months are zero-based
	const day = String(date.getDate()).padStart(2, '0')
	return `${year}/${month}/${day}`
}

export function ListShortURL({ user }: { user: User }) {
	const [showOthers, setShowOthers] = useState(false)
	const {
		data: urls,
		error,
		mutate,
	} = useSWR<URL[]>(['list-urls', showOthers], ([, showOthers]) => (showOthers ? adminGetUrls : getUrls)())

	useEffect(() => console.log(urls), [urls])

	const handleDelete = async (id: number) => {
		if (window.confirm('Are you sure you want to delete this URL?')) {
			try {
				await deleteUrl(id)
				mutate()
			} catch (err: any) {
				toast.error('Failed to delete URL.')
			}
		}
	}

	if (error) {
		return <div className="alert alert-error">{error.message}</div>
	}

	if (!urls) {
		return <span className="loading loading-spinner loading-lg" />
	}

	return (
		<div>
			{(canDeleteAny(user.permissions) || canViewAnyStats(user.permissions)) && (
				<label className="label w-full px-4">
					<input
						type="checkbox"
						className="toggle toggle-primary toggle-xs"
						checked={showOthers}
						onChange={e => {
							setShowOthers(e.target.checked)
						}}
					/>
					<span>Show others urls</span>
				</label>
			)}
			<div className="overflow-x-auto">
				<table className="table w-full">
					<thead>
						<tr>
							{showOthers && <th>Owner</th>}
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
								{showOthers && <td>{url.Username || 'Unknown'}</td>}
								<td>
									<a
										href={formatShortPath(url.ShortPath)}
										target="_blank"
										rel="noopener noreferrer"
										className="link link-primary"
									>
										{url.ShortPath}
									</a>
								</td>
								<td className="max-w-xs truncate">{url.OriginalURL}</td>
								<td>{url.TotalClicks}</td>
								<td>{formatDate(url.CreatedAt)}</td>
								<td className="join">
									{canViewOwnStats(user.permissions) && (
										<a href={`/dashboard/stats?id=${url.ID}`} className="btn btn-sm join-item">
											Stats
										</a>
									)}
									{canDeleteOwn(user.permissions) && (
										<button
											onClick={() => handleDelete(url.ID)}
											className="btn btn-sm text-error join-item"
										>
											Delete
										</button>
									)}
								</td>
							</tr>
						))}
					</tbody>
				</table>
			</div>
		</div>
	)
}
