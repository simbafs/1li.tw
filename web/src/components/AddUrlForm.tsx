import React, { useEffect, useState } from 'react'
import { BASE, createUrl } from '../lib/api'
import { Input } from './Input'

export function AddUrlForm() {
	const [originalUrl, setOriginalUrl] = useState('https://')
	const [customPath, setCustomPath] = useState('')
	const [error, setError] = useState('')
	const [success, setSuccess] = useState('')

	const [loggedIn, setLoggedIn] = useState(false)

	useEffect(() => {
		const user = localStorage.getItem('user')
		if (user) {
			setLoggedIn(true)
		} else {
			setLoggedIn(false)
		}
	}, [])

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault()
		setError('')
		setSuccess('')

		if (!originalUrl) {
			setError('Original URL is required.')
			return
		}

		try {
			const data = await createUrl({
				original_url: originalUrl,
				...(customPath && { custom_path: customPath }),
			})
			setSuccess(`Success! Short URL is: ${BASE}/${data.ShortPath}`)
			setOriginalUrl('')
			setCustomPath('')
		} catch (err: any) {
			setError(err.info?.message || 'Failed to create short URL.')
		}
	}

	return (
		<div className="w-full px-20">
			<h2 className="card-title">
				{loggedIn ? 'Create a new Short URL' : 'Create a quick, anonymous Short URL'}
			</h2>
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
				{loggedIn && (
					<Input
						label="Custom Path"
						type="text"
						placeholder="my-custom-path"
						value={customPath}
						onChange={e => setCustomPath(e.target.value)}
						optional
					/>
				)}
				<button type="submit" className="btn btn-primary w-full">
					Shorten
				</button>
			</form>
			{error && <div className="alert alert-error mt-4">{error}</div>}
			{success && <div className="alert alert-success mt-4">{success}</div>}
			{/* </div> */}
		</div>
	)
}
