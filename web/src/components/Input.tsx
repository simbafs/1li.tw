import React from 'react'

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
	label: string
	validate?: boolean
	required?: boolean
}

export function Input({ label, validate, required = false, ...props }: InputProps) {
	return (
		<div className="form-control w-full">
			<label className="label">
				<span className="label-text">{label}</span>
				{required && (
					<span className="label-text-alt tooltip text-red-500" data-tip="required">
						*
					</span>
				)}
			</label>
			<input
				required={required}
				className={`input w-full ${validate ? 'input-bordered' : ''} ${props.className || ''}`}
				{...props}
			/>
		</div>
	)
}
