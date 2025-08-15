import React, { useState } from 'react'
import { login } from '../lib/api'

export function LoginForm() {
	const [username, setUsername] = useState('')
	const [password, setPassword] = useState('')
	const [error, setError] = useState('')

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault()
		setError('')

		if (!username || !password) {
			setError('Username and password are required.')
			return
		}

		try {
			const user = await login({ username, password })
			localStorage.setItem('user', JSON.stringify(user))
			window.location.href = '/dashboard'
		} catch (err: any) {
			setError(err.info?.message || 'Failed to log in.')
		}
	}

	return (
		<div className="card bg-base-200 w-full max-w-sm shadow-xl">
			<div className="card-body">
				<h2 className="card-title">Login</h2>
				<form onSubmit={handleSubmit}>
					<div className="form-control">
						<label className="label">
							<span className="label-text">Username</span>
						</label>
						<input
							type="text"
							placeholder="username"
							className="input input-bordered"
							value={username}
							onChange={e => setUsername(e.target.value)}
							required
						/>
					</div>
					<div className="form-control">
						<label className="label">
							<span className="label-text">Password</span>
						</label>
						<input
							type="password"
							placeholder="password"
							className="input input-bordered"
							value={password}
							onChange={e => setPassword(e.target.value)}
							required
						/>
					</div>
					<div className="form-control mt-6">
						<button type="submit" className="btn btn-primary">
							Login
						</button>
					</div>
				</form>
				{error && <div className="alert alert-error mt-4">{error}</div>}
				<div className="mt-4 text-center">
					<a href="/register" className="link">
						Don't have an account? Register
					</a>
				</div>
			</div>
		</div>
	)
}
