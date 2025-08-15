import {
	BarChart,
	Bar,
	LineChart,
	Line,
	XAxis,
	YAxis,
	CartesianGrid,
	Tooltip,
	Legend,
	ResponsiveContainer,
} from 'recharts'

export type Stats = {
	total: number
	by_time: { time: string; count: number }[]
	by_country: { country: string; count: number }[]
	by_os: { os: string; count: number }[]
	by_browser: { browser: string; count: number }[]
}

export function StatsCharts({ stats }: { stats: Stats }) {
	return (
		<div className="grid grid-cols-1 gap-8 lg:grid-cols-2">
			<div className="card bg-base-100 shadow-xl">
				<div className="card-body">
					<h2 className="card-title">Clicks by Time</h2>
					<ResponsiveContainer width="100%" height={300}>
						<LineChart data={stats.by_time}>
							<CartesianGrid strokeDasharray="3 3" />
							<XAxis dataKey="time" tickFormatter={time => new Date(time).toLocaleDateString()} />
							<YAxis />
							<Tooltip />
							<Legend />
							<Line type="monotone" dataKey="count" stroke="#8884d8" />
						</LineChart>
					</ResponsiveContainer>
				</div>
			</div>
			<div className="card bg-base-100 shadow-xl">
				<div className="card-body">
					<h2 className="card-title">Clicks by Country</h2>
					<ResponsiveContainer width="100%" height={300}>
						<BarChart data={stats.by_country}>
							<CartesianGrid strokeDasharray="3 3" />
							<XAxis dataKey="country" />
							<YAxis />
							<Tooltip />
							<Legend />
							<Bar dataKey="count" fill="#82ca9d" />
						</BarChart>
					</ResponsiveContainer>
				</div>
			</div>
			<div className="card bg-base-100 shadow-xl">
				<div className="card-body">
					<h2 className="card-title">Clicks by OS</h2>
					<ResponsiveContainer width="100%" height={300}>
						<BarChart data={stats.by_os}>
							<CartesianGrid strokeDasharray="3 3" />
							<XAxis dataKey="os" />
							<YAxis />
							<Tooltip />
							<Legend />
							<Bar dataKey="count" fill="#8884d8" />
						</BarChart>
					</ResponsiveContainer>
				</div>
			</div>
			<div className="card bg-base-100 shadow-xl">
				<div className="card-body">
					<h2 className="card-title">Clicks by Browser</h2>
					<ResponsiveContainer width="100%" height={300}>
						<BarChart data={stats.by_browser}>
							<CartesianGrid strokeDasharray="3 3" />
							<XAxis dataKey="browser" />
							<YAxis />
							<Tooltip />
							<Legend />
							<Bar dataKey="count" fill="#ffc658" />
						</BarChart>
					</ResponsiveContainer>
				</div>
			</div>
		</div>
	)
}
