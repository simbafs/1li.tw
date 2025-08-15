import React, { useState } from 'react'
import { register } from '../lib/api'
import { Input } from './Input'

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
		<fieldset className="fieldset bg-base-200 border-base-300 rounded-box w-xs border px-8 py-4">
			<legend className="card-title">Register</legend>
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
					Register
				</button>
			</form>
			{error && <div className="alert alert-error mt-4">{error}</div>}
			{success && <div className="alert alert-success mt-4">{success}</div>}
			<div className="mt-4 text-center">
				Already have an account?{' '}
				<a href="/login" className="link">
					Login
				</a>
			</div>
		</fieldset>
	)
}
