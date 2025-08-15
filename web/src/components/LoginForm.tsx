import React, { useState } from 'react'
import { login } from '../lib/api'
import { Input } from './Input'

export function LoginForm() {
	const [username, setUsername] = useState('')
	const [password, setPassword] = useState('')
	const [error, setError] = useState('')

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault()
		setError('')

		try {
			const data = await login({ username, password })
			localStorage.setItem('user', JSON.stringify(data))
			window.location.href = '/dashboard'
		} catch (err: any) {
			setError(err.info?.message || 'Failed to login.')
		}
	}

	return (
		<fieldset className="fieldset bg-base-200 border-base-300 rounded-box w-xs border px-8 py-4">
			<legend className="card-title">Login</legend>

			<form onSubmit={handleSubmit} className="flex flex-col gap-4">
				<Input
					label="Username"
					type="text"
					placeholder="your username"
					value={username}
					onChange={e => setUsername(e.target.value)}
					required
				/>
				<Input
					label="Password"
					type="password"
					placeholder="your password"
					value={password}
					onChange={e => setPassword(e.target.value)}
					required
				/>
				<button type="submit" className="btn btn-primary w-full">
					Login
				</button>
			</form>
			{error && <div className="alert alert-error mt-4">{error}</div>}
			<div className="mt-4 text-center">
				Do not have an account?{' '}
				<a href="/register" className="link">
					Register
				</a>
			</div>
		</fieldset>
	)
}
