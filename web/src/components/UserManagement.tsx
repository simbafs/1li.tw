import { useEffect, useState } from 'react'
import { listUsers, updateUserPermission, deleteUser } from '../lib/api'
import { Permission, permissionNames } from '../lib/permissions'
import { toast, ToastContainer } from 'react-toastify'

interface User {
	ID: number
	Username: string
	Permissions: number
	TelegramChatID: number
	CreatedAt: string
}

const UserManagement = () => {
	const [users, setUsers] = useState<User[]>([])
	const [error, setError] = useState<string | null>(null)

	const fetchUsers = async () => {
		try {
			const data = (await listUsers()) as User[]
			setUsers(data.filter(u => u.ID !== 1))
		} catch (err) {
			setError('Failed to fetch users.')
			console.error(err)
		}
	}

	useEffect(() => {
		fetchUsers()
	}, [])

	const handlePermissionChange = async (userId: number, permission: Permission, checked: boolean) => {
		const user = users.find(u => u.ID === userId)
		if (!user) return

		const newPermissions = checked ? user.Permissions | permission : user.Permissions & ~permission

		try {
			await updateUserPermission(userId, newPermissions)
			setUsers(users.map(u => (u.ID === userId ? { ...u, Permissions: newPermissions } : u)))
		} catch (err) {
			toast.error(`Failed to update permissions: ${err}`)
			console.error(err)
		}
	}

	const handleDelete = async (userId: number) => {
		if (window.confirm('Are you sure you want to delete this user?')) {
			try {
				await deleteUser(userId)
				setUsers(users.filter(u => u.ID !== userId))
			} catch (err) {
				toast.error(`Failed to delete user: ${err}`)
				console.error(err)
			}
		}
	}

	if (error) {
		return <div className="text-red-500">{error}</div>
	}

	return (
		<>
			<div className="overflow-x-auto">
				<table className="table w-full">
					<thead>
						<tr>
							<th>ID</th>
							<th>Username</th>
							<th>Permissions</th>
							<th>Created At</th>
							<th>Actions</th>
						</tr>
					</thead>
					<tbody>
						{users.map(user => (
							<tr key={user.ID}>
								<td>{user.ID}</td>
								<td>{user.Username}</td>
								<td>
									<div className="flex flex-col">
										{Object.entries(permissionNames).map(([value, name]) => {
											const p = parseInt(value) as Permission
											return (
												<label key={p} className="label cursor-pointer">
													<input
														type="checkbox"
														className="checkbox checkbox-primary checkbox-sm"
														checked={(user.Permissions & p) === p}
														onChange={e =>
															handlePermissionChange(user.ID, p, e.target.checked)
														}
													/>
													<span className="label-text">{name}</span>
												</label>
											)
										})}
									</div>
								</td>
								<td>{new Date(user.CreatedAt).toLocaleString()}</td>
								<td>
									<button className="btn btn-sm btn-error" onClick={() => handleDelete(user.ID)}>
										Delete
									</button>
								</td>
							</tr>
						))}
					</tbody>
				</table>
			</div>
			<ToastContainer position="bottom-right" theme="light" />
		</>
	)
}

export default UserManagement
