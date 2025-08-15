import React from 'react'

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
	label: string
	validate?: boolean
	required?: boolean
	prefix?: string
}

export function Input({ label, validate, required = false, prefix, ...props }: InputProps) {
	return (
		<>
			<label className="label">
				<span className="label-text">{label}</span>
				{required && (
					<span className="label-text-alt tooltip text-red-500" data-tip="required">
						*
					</span>
				)}
			</label>
			{prefix ? (
				<label className="input flex w-full">
					<span className="label">{prefix}</span>
					<input required={required} className={`w-full ${props.className || ''}`} {...props} />
				</label>
			) : (
				<input required={required} className={`input w-full ${props.className || ''}`} {...props} />
			)}
		</>
	)
}
