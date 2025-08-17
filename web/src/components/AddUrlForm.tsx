import React, { useState } from 'react'
import { createUrl } from '../lib/api'
import { canCreateAny, canCreatePrefix } from '../lib/permissions'
import { Input } from './Input'
import { useUser } from '../hooks/useUser'
import { formatShortPath } from '../lib/formatShortPath'
import { mutate } from 'swr'

export function AddUrlForm({ canCollapse = false }: { canCollapse?: boolean }) {
	const [collapsed, setCollapsed] = useState(canCollapse)

	const [originalUrl, setOriginalUrl] = useState('https://')
	const [withPrefix, setWithPrefix] = useState(true)
	const [customPath, setCustomPath] = useState('')

	const [error, setError] = useState('')
	const [success, setSuccess] = useState('')

	const user = useUser()

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault()

		setError('')
		setSuccess('')

		if (!originalUrl) {
			setError('Original URL is required.')
			return
		}

		let path = customPath.trim().replaceAll('/', '')
		if (withPrefix && user) {
			path = `@${user.username}/${customPath}`
		}

		try {
			const data = await createUrl({
				original_url: originalUrl,
				...(customPath && { custom_path: path }),
			})
			setSuccess(`Success! Short URL is: ${formatShortPath(data.ShortPath)}`)
			setOriginalUrl('https://')
			setCustomPath('')
			mutate(['list-urls', false])
		} catch (err: any) {
			setError(err.info?.message || 'Failed to create short URL.')
		}
	}

	return (
		<div className="collapse w-full">
			<input type="checkbox" checked={collapsed || !canCollapse} onChange={() => setCollapsed(!collapsed)} />
			<h2 className="collapse-title">
				{user ? 'Create a new Short URL' : 'Create a quick, anonymous Short URL'}
				{canCollapse && <span className="collapse-arrow">{collapsed ? '▼' : '▲'}</span>}
			</h2>
			<div className="collapse-content">
				<form onSubmit={handleSubmit} className="mt-6 flex flex-col gap-4">
					<Input
						label="Original URL"
						type="url"
						placeholder="https://example.com"
						value={originalUrl}
						onChange={e => setOriginalUrl(e.target.value)}
						pattern="^(https?://)?([a-zA-Z0-9]([a-zA-Z0-9\-].*[a-zA-Z0-9])?\.)+[a-zA-Z].*$"
						required
						validate
					/>
					{user && canCreatePrefix(user.permissions) && (
						<>
							<label className="label">Custom Path</label>
							{user && canCreateAny(user.permissions) && (
								<label className="label">
									<input
										type="checkbox"
										className="toggle toggle-primary toggle-xs"
										checked={withPrefix}
										onChange={e => setWithPrefix(e.target.checked)}
									/>
									With Prefix
								</label>
							)}
							<label className="input join-item flex w-full">
								{withPrefix && <span className="label">{`@${user.username}/`}</span>}
								<input
									className="w-full"
									type="text"
									placeholder="my-custom-path"
									value={customPath}
									onChange={e => setCustomPath(e.target.value)}
								/>
							</label>
						</>
					)}
					<button type="submit" className="btn btn-primary w-full">
						Shorten
					</button>
				</form>
				{error && <div className="alert alert-error mt-4">{error}</div>}
				{success && <div className="alert alert-success mt-4">{success}</div>}
			</div>
		</div>
	)
}
