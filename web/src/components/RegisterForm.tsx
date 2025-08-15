import React, { useState } from 'react'
import { register } from '../lib/api'

export function RegisterForm() {
	const [username, setUsername] = useState('')
	const [password, setPassword] = useState('')
	const [error, setError] = useState('')
	const [success, setSuccess] = useState('')

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault()
		setError('')
		setSuccess('')

		if (!username || !password) {
			setError('Username and password are required.')
			return
		}

		try {
			await register({ username, password })
			setSuccess('Registration successful! You can now log in.')
			setUsername('')
			setPassword('')
			setTimeout(() => {
				window.location.href = '/login'
			}, 2000)
		} catch (err: any) {
			setError(err.info?.message || 'Failed to register.')
		}
	}

	return (
		<div className="card bg-base-200 w-full max-w-sm shadow-xl">
			<div className="card-body">
				<h2 className="card-title">Register</h2>
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
							Register
						</button>
					</div>
				</form>
				{error && <div className="alert alert-error mt-4">{error}</div>}
				{success && <div className="alert alert-success mt-4">{success}</div>}
				<div className="mt-4 text-center">
					<a href="/login" className="link">
						Already have an account? Login
					</a>
				</div>
			</div>
		</div>
	)
}
