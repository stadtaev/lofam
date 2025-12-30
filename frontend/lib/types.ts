export type TaskStatus = 'todo' | 'in_progress' | 'done'
export type TaskPriority = 'low' | 'medium' | 'high'

export interface Task {
  id: number
  title: string
  description: string
  status: TaskStatus
  priority: TaskPriority
  dueDate: string | null
  createdAt: string
}

export interface CreateTaskRequest {
  title: string
  description?: string
  status?: TaskStatus
  priority?: TaskPriority
  dueDate?: string
}

export interface UpdateTaskRequest {
  title?: string
  description?: string
  status?: TaskStatus
  priority?: TaskPriority
  dueDate?: string | null
}

export type NoteColor = 'yellow' | 'pink' | 'green'

export interface Note {
  id: number
  title: string
  content: string
  color: NoteColor
  createdAt: string
  updatedAt: string
}

export interface CreateNoteRequest {
  title: string
  content?: string
  color: NoteColor
}

export interface UpdateNoteRequest {
  title: string
  content: string
  color: NoteColor
}
