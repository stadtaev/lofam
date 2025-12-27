import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import type { Task, CreateTaskRequest, TaskStatus, TaskPriority } from '../types/task'
import { listTasks, createTask, updateTask, deleteTask } from '../api/tasks'

export const Route = createFileRoute('/')({
  component: TasksPage,
})

function TasksPage() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [editingTask, setEditingTask] = useState<Task | null>(null)
  const [isCreating, setIsCreating] = useState(false)

  async function fetchTasks() {
    try {
      setLoading(true)
      const data = await listTasks()
      setTasks(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch tasks')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchTasks()
  }, [])

  async function handleCreate(data: CreateTaskRequest) {
    try {
      await createTask(data)
      setIsCreating(false)
      fetchTasks()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create task')
    }
  }

  async function handleUpdate(id: number, data: Partial<Task>) {
    try {
      await updateTask(id, data)
      setEditingTask(null)
      fetchTasks()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update task')
    }
  }

  async function handleDelete(id: number) {
    if (!confirm('Delete this task?')) return
    try {
      await deleteTask(id)
      fetchTasks()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete task')
    }
  }

  async function handleStatusChange(task: Task, status: TaskStatus) {
    await handleUpdate(task.id, { status })
  }

  if (loading) return <p aria-busy="true">Loading tasks...</p>

  return (
    <>
      <header>
        <h1>Tasks</h1>
        <button onClick={() => setIsCreating(true)}>New Task</button>
      </header>

      {error && <p style={{ color: 'var(--pico-color-red-500)' }}>{error}</p>}

      {isCreating && (
        <TaskForm
          onSubmit={handleCreate}
          onCancel={() => setIsCreating(false)}
        />
      )}

      {editingTask && (
        <TaskForm
          task={editingTask}
          onSubmit={(data) => handleUpdate(editingTask.id, data)}
          onCancel={() => setEditingTask(null)}
        />
      )}

      {tasks.length === 0 ? (
        <p>No tasks yet. Create one to get started.</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>Title</th>
              <th>Status</th>
              <th>Priority</th>
              <th>Due Date</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {tasks.map((task) => (
              <tr key={task.id}>
                <td>
                  <strong>{task.title}</strong>
                  {task.description && (
                    <small style={{ display: 'block', opacity: 0.7 }}>
                      {task.description}
                    </small>
                  )}
                </td>
                <td>
                  <StatusBadge
                    status={task.status}
                    onChange={(status) => handleStatusChange(task, status)}
                  />
                </td>
                <td>
                  <PriorityBadge priority={task.priority} />
                </td>
                <td>{task.dueDate ? formatDate(task.dueDate) : '-'}</td>
                <td>
                  <div style={{ display: 'flex', gap: '0.5rem' }}>
                    <button
                      className="outline"
                      onClick={() => setEditingTask(task)}
                    >
                      Edit
                    </button>
                    <button
                      className="outline secondary"
                      onClick={() => handleDelete(task.id)}
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </>
  )
}

interface TaskFormProps {
  task?: Task
  onSubmit: (data: CreateTaskRequest) => void
  onCancel: () => void
}

function TaskForm({ task, onSubmit, onCancel }: TaskFormProps) {
  const [title, setTitle] = useState(task?.title ?? '')
  const [description, setDescription] = useState(task?.description ?? '')
  const [status, setStatus] = useState<TaskStatus>(task?.status ?? 'todo')
  const [priority, setPriority] = useState<TaskPriority>(task?.priority ?? 'medium')
  const [dueDate, setDueDate] = useState(task?.dueDate ? task.dueDate.split('T')[0] : '')

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    onSubmit({
      title,
      description: description || undefined,
      status,
      priority,
      dueDate: dueDate ? `${dueDate}T00:00:00Z` : undefined,
    })
  }

  return (
    <dialog open>
      <article>
        <header>
          <button aria-label="Close" rel="prev" onClick={onCancel} />
          <h2>{task ? 'Edit Task' : 'New Task'}</h2>
        </header>
        <form onSubmit={handleSubmit}>
          <label>
            Title
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
              autoFocus
            />
          </label>
          <label>
            Description
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
            />
          </label>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '1rem' }}>
            <label>
              Status
              <select value={status} onChange={(e) => setStatus(e.target.value as TaskStatus)}>
                <option value="todo">To Do</option>
                <option value="in_progress">In Progress</option>
                <option value="done">Done</option>
              </select>
            </label>
            <label>
              Priority
              <select value={priority} onChange={(e) => setPriority(e.target.value as TaskPriority)}>
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
              </select>
            </label>
            <label>
              Due Date
              <input
                type="date"
                value={dueDate}
                onChange={(e) => setDueDate(e.target.value)}
              />
            </label>
          </div>
          <footer>
            <button type="button" className="secondary" onClick={onCancel}>
              Cancel
            </button>
            <button type="submit">{task ? 'Save' : 'Create'}</button>
          </footer>
        </form>
      </article>
    </dialog>
  )
}

function StatusBadge({
  status,
  onChange,
}: {
  status: TaskStatus
  onChange: (status: TaskStatus) => void
}) {
  const colors: Record<TaskStatus, string> = {
    todo: 'var(--pico-color-grey-500)',
    in_progress: 'var(--pico-color-blue-500)',
    done: 'var(--pico-color-green-500)',
  }
  const labels: Record<TaskStatus, string> = {
    todo: 'To Do',
    in_progress: 'In Progress',
    done: 'Done',
  }

  return (
    <select
      value={status}
      onChange={(e) => onChange(e.target.value as TaskStatus)}
      style={{
        backgroundColor: colors[status],
        color: 'white',
        padding: '0.25rem 0.5rem',
        border: 'none',
        borderRadius: '4px',
        fontSize: '0.875rem',
      }}
    >
      {Object.entries(labels).map(([value, label]) => (
        <option key={value} value={value}>
          {label}
        </option>
      ))}
    </select>
  )
}

function PriorityBadge({ priority }: { priority: TaskPriority }) {
  const colors: Record<TaskPriority, string> = {
    low: 'var(--pico-color-grey-500)',
    medium: 'var(--pico-color-amber-500)',
    high: 'var(--pico-color-red-500)',
  }

  return (
    <span
      style={{
        backgroundColor: colors[priority],
        color: 'white',
        padding: '0.25rem 0.5rem',
        borderRadius: '4px',
        fontSize: '0.875rem',
        textTransform: 'capitalize',
      }}
    >
      {priority}
    </span>
  )
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleDateString()
}
