import type {
  Task,
  CreateTaskRequest,
  UpdateTaskRequest,
  Note,
  CreateNoteRequest,
  UpdateNoteRequest,
  Wishlist,
  CreateWishlistRequest,
  UpdateWishlistRequest,
} from './types'

const API_BASE = process.env.NEXT_PUBLIC_API_URL || ''

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

export async function listNotes(): Promise<Note[]> {
  const response = await fetch(`${API_BASE}/api/notes`)
  return handleResponse<Note[]>(response)
}

export async function getNote(id: number): Promise<Note> {
  const response = await fetch(`${API_BASE}/api/notes/${id}`)
  return handleResponse<Note>(response)
}

export async function createNote(data: CreateNoteRequest): Promise<Note> {
  const response = await fetch(`${API_BASE}/api/notes`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse<Note>(response)
}

export async function updateNote(id: number, data: UpdateNoteRequest): Promise<Note> {
  const response = await fetch(`${API_BASE}/api/notes/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse<Note>(response)
}

export async function deleteNote(id: number): Promise<void> {
  const response = await fetch(`${API_BASE}/api/notes/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    const error = await response.text()
    throw new Error(error || `HTTP ${response.status}`)
  }
}

export async function listWishlists(): Promise<Wishlist[]> {
  const response = await fetch(`${API_BASE}/api/wishlists`)
  return handleResponse<Wishlist[]>(response)
}

export async function getWishlist(id: number): Promise<Wishlist> {
  const response = await fetch(`${API_BASE}/api/wishlists/${id}`)
  return handleResponse<Wishlist>(response)
}

export async function createWishlist(data: CreateWishlistRequest): Promise<Wishlist> {
  const response = await fetch(`${API_BASE}/api/wishlists`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse<Wishlist>(response)
}

export async function updateWishlist(id: number, data: UpdateWishlistRequest): Promise<Wishlist> {
  const response = await fetch(`${API_BASE}/api/wishlists/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  return handleResponse<Wishlist>(response)
}

export async function deleteWishlist(id: number): Promise<void> {
  const response = await fetch(`${API_BASE}/api/wishlists/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    const error = await response.text()
    throw new Error(error || `HTTP ${response.status}`)
  }
}
