import type { Task, CreateTaskRequest, UpdateTaskRequest } from './types'

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error = await response.text()
    throw new Error(error || `HTTP ${response.status}`)
  }
  return response.json()
}

export async function listTasks(): Promise<Task[]> {
  const response = await fetch(`${API_BASE}/api/tasks`)
  return handleResponse<Task[]>(response)
}

export async function getTask(id: number): Promise<Task> {
  const response = await fetch(`${API_BASE}/api/tasks/${id}`)
  return handleResponse<Task>(response)
}

export async function createTask(data: CreateTaskRequest): Promise<Task> {
  const response = await fetch(`${API_BASE}/api/tasks`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse<Task>(response)
}

export async function updateTask(id: number, data: UpdateTaskRequest): Promise<Task> {
  const response = await fetch(`${API_BASE}/api/tasks/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse<Task>(response)
}

export async function deleteTask(id: number): Promise<void> {
  const response = await fetch(`${API_BASE}/api/tasks/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    const error = await response.text()
    throw new Error(error || `HTTP ${response.status}`)
  }
}
