import { useState, useEffect } from 'react'
import { getUrlStats } from '../lib/api'
import { StatsCharts, type Stats } from './StatsCharts'

export function StatsPage() {
	const [urlId, setUrlId] = useState<number | null>(null)
	const [stats, setStats] = useState<Stats | null>(null)
	const [error, setError] = useState('')
	const [loading, setLoading] = useState(true)

	useEffect(() => {
		const urlParams = new URLSearchParams(window.location.search)
		const id = urlParams.get('id')
		if (id) {
			setUrlId(Number(id))
		} else {
			setError('URL ID is missing.')
			setLoading(false)
		}
	}, [])

	useEffect(() => {
		if (!urlId) return

		const fetchStats = async () => {
			try {
				const data = await getUrlStats(urlId)
				setStats(data)
			} catch (err: any) {
				if (err.status === 401) {
					localStorage.removeItem('user')
					window.location.href = '/login'
				} else {
					setError('Failed to fetch stats.')
				}
			} finally {
				setLoading(false)
			}
		}
		fetchStats()
	}, [urlId])

	if (loading) {
		return (
			<div className="text-center">
				<span className="loading loading-spinner loading-lg"></span>
			</div>
		)
	}

	if (error) {
		return <div className="alert alert-error">{error}</div>
	}

	if (!stats) {
		return <div>No stats available.</div>
	}

	return (
		<div>
			<h1 className="mb-4 text-3xl font-bold">Statistics</h1>
			<div className="stats mb-8 shadow">
				<div className="stat">
					<div className="stat-title">Total Clicks</div>
					<div className="stat-value">{stats.total}</div>
				</div>
			</div>
			<StatsCharts stats={stats} />
		</div>
	)
}
