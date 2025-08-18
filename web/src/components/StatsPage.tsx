import { getUrlStats } from '../lib/api'
import { StatsCharts } from './StatsCharts'
import useSWR from 'swr'

async function fetchStats() {
	const params = new URLSearchParams(window.location.search)
	const id = params.get('id')
	if (!id) {
		return
	}

	return getUrlStats(Number(id))
}

export function StatsPage() {
	const { data: stats, error } = useSWR('get-stats', fetchStats)

	if (error) {
		return <div className="alert alert-error">{error}</div>
	}

	if (!stats) {
		return (
			<div className="text-center">
				<span className="loading loading-spinner loading-lg"></span>
			</div>
		)
	}

	return (
		<div>
			<h1 className="mb-4 text-3xl font-bold">Statistics</h1>
			<div className="md:stats stats-vertical mb-8 shadow">
				<div className="stat">
					<div className="stat-title">Owner</div>
					<div className="stat-value">{stats.owner_name}</div>
				</div>
				<div className="stat">
					<div className="stat-title">Original URL</div>
					<div className="stat-value overflow-scroll">
						<a href={stats.url.OriginalURL} target="_blank">
							{stats.url.OriginalURL}
						</a>
					</div>
				</div>
				<div className="stat">
					<div className="stat-title">Total Clicks</div>
					<div className="stat-value">{stats.total}</div>
				</div>
			</div>
			<StatsCharts stats={stats} />
		</div>
	)
}
