import { useEffect } from 'react'
import {
	LineChart,
	Line,
	XAxis,
	YAxis,
	CartesianGrid,
	Tooltip,
	Legend,
	ResponsiveContainer,
	PieChart,
	Pie,
} from 'recharts'

export type Stats = {
	total: number
	by_time: {
		bucketStart: string
		count: number
	}[]
	by_country: {
		key: string
		count: number
	}[]
	by_os: {
		key: string
		count: number
	}[]
	by_browser: {
		key: string
		count: number
	}[]
}

function DrawPieChart({
	title,
	data,
	fill,
	nameKey = 'key',
	dataKey = 'count',
}: {
	title: string
	data: any[]
	fill: string
	nameKey?: string
	dataKey?: string
}) {
	return (
		<div className="card bg-base-100 shadow-xl">
			<div className="card-body">
				<h2 className="card-title">{title}</h2>
				<ResponsiveContainer width="100%" height={300}>
					<PieChart data={data}>
						<Pie dataKey={dataKey} nameKey={nameKey} fill={fill || '#82ca9d'} />
						<legend />
						<Tooltip />
					</PieChart>
				</ResponsiveContainer>
			</div>
		</div>
	)
}

export function StatsCharts({ stats }: { stats: Stats }) {
	useEffect(() => console.log(stats), [stats])

	return (
		<div className="grid grid-cols-1 gap-8 lg:grid-cols-2">
			<div className="card bg-base-100 shadow-xl">
				<div className="card-body">
					<h2 className="card-title">Clicks by Time</h2>
					<ResponsiveContainer width="100%" height={300}>
						<LineChart data={stats.by_time}>
							<CartesianGrid strokeDasharray="3 3" />
							<XAxis dataKey="bucketStart" tickFormatter={time => new Date(time).toLocaleDateString()} />
							<YAxis />
							<Tooltip />
							<Legend />
							<Line type="monotone" dataKey="count" stroke="#8884d8" />
						</LineChart>
					</ResponsiveContainer>
				</div>
			</div>
			<DrawPieChart title="Clicks by Country" data={stats.by_country} fill="#82ca9d" />
			<DrawPieChart title="Clicks by OS" data={stats.by_os} fill="#8884d8" />
			<DrawPieChart title="Clicks by Browser" data={stats.by_browser} fill="#ffc658" />
		</div>
	)
}
