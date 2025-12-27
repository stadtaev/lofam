'use client'

import { useState, useEffect, useCallback } from 'react'
import { Calendar } from '@/components/Calendar'
import { TaskList } from '@/components/TaskList'
import { TaskModal } from '@/components/TaskModal'
import { TodaySection } from '@/components/TodaySection'
import { listTasks, createTask, updateTask, deleteTask } from '@/lib/api'
import type { Task, CreateTaskRequest } from '@/lib/types'

export default function Home() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const today = new Date()
  const [year, setYear] = useState(today.getFullYear())
  const [month, setMonth] = useState(today.getMonth())
  const [selectedDate, setSelectedDate] = useState<Date | null>(null)
  const [searchQuery, setSearchQuery] = useState('')

  const [showModal, setShowModal] = useState(false)
  const [editingTask, setEditingTask] = useState<Task | null>(null)

  const fetchTasks = useCallback(async () => {
    try {
      setLoading(true)
      const data = await listTasks()
      setTasks(data || [])
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch tasks')
      setTasks([])
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchTasks()
  }, [fetchTasks])

  const handlePrevMonth = () => {
    if (month === 0) {
      setMonth(11)
      setYear(year - 1)
    } else {
      setMonth(month - 1)
    }
  }

  const handleNextMonth = () => {
    if (month === 11) {
      setMonth(0)
      setYear(year + 1)
    } else {
      setMonth(month + 1)
    }
  }

  const handleDateSelect = (date: Date) => {
    setSelectedDate(date)
  }

  const handleTaskClick = (task: Task) => {
    setEditingTask(task)
    setShowModal(true)
  }

  const handleAddTask = () => {
    setEditingTask(null)
    setShowModal(true)
  }

  const handleSaveTask = async (data: CreateTaskRequest) => {
    try {
      if (editingTask) {
        await updateTask(editingTask.id, data)
      } else {
        await createTask(data)
      }
      setShowModal(false)
      setEditingTask(null)
      fetchTasks()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save task')
    }
  }

  const handleDeleteTask = async () => {
    if (!editingTask) return
    try {
      await deleteTask(editingTask.id)
      setShowModal(false)
      setEditingTask(null)
      fetchTasks()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete task')
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-gray-400">Loading...</p>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-white p-8">
      {error && (
        <div className="fixed top-4 right-4 bg-red-100 text-red-600 px-4 py-2 rounded-lg">
          {error}
          <button onClick={() => setError(null)} className="ml-2">✕</button>
        </div>
      )}

      <div className="max-w-6xl mx-auto flex flex-col lg:flex-row gap-8">
        <div className="flex-shrink-0">
          <Calendar
            year={year}
            month={month}
            tasks={tasks}
            selectedDate={selectedDate}
            onDateSelect={handleDateSelect}
            onPrevMonth={handlePrevMonth}
            onNextMonth={handleNextMonth}
          />
          <TodaySection tasks={tasks} onAddTask={handleAddTask} />
        </div>

        <TaskList
          tasks={tasks}
          month={month}
          searchQuery={searchQuery}
          onSearchChange={setSearchQuery}
          onTaskClick={handleTaskClick}
        />
      </div>

      {showModal && (
        <TaskModal
          task={editingTask}
          initialDate={selectedDate}
          onSave={handleSaveTask}
          onDelete={editingTask ? handleDeleteTask : undefined}
          onClose={() => {
            setShowModal(false)
            setEditingTask(null)
          }}
        />
      )}
    </div>
  )
}
