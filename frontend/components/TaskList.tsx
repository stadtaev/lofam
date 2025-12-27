'use client'

import { getDayName, getMonthName } from '@/lib/date-utils'
import type { Task, TaskPriority } from '@/lib/types'

interface TaskListProps {
  tasks: Task[]
  month: number
  searchQuery: string
  onSearchChange: (query: string) => void
  onTaskClick: (task: Task) => void
}

const PRIORITY_COLORS: Record<TaskPriority, string> = {
  high: 'text-red-500',
  medium: 'text-amber-500',
  low: 'text-gray-400',
}

export function TaskList({
  tasks,
  month,
  searchQuery,
  onSearchChange,
  onTaskClick,
}: TaskListProps) {
  const filteredTasks = tasks.filter(
    (task) =>
      task.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      task.description.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const tasksByDate = filteredTasks
    .filter((task) => task.dueDate)
    .sort((a, b) => {
      if (!a.dueDate || !b.dueDate) return 0
      return new Date(a.dueDate).getTime() - new Date(b.dueDate).getTime()
    })
    .reduce((acc, task) => {
      const date = task.dueDate!.split('T')[0]
      if (!acc[date]) acc[date] = []
      acc[date].push(task)
      return acc
    }, {} as Record<string, Task[]>)

  const sortedDates = Object.keys(tasksByDate).sort()

  return (
    <div className="flex-1 border-t lg:border-t-0 lg:border-l pt-8 lg:pt-0 lg:pl-8">
      <div className="mb-6">
        <div className="relative">
          <input
            type="text"
            placeholder="Search"
            value={searchQuery}
            onChange={(e) => onSearchChange(e.target.value)}
            className="w-full px-4 py-2 pl-10 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-red-200"
          />
          <svg
            className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
        </div>
      </div>

      <div className="mb-4">
        <h3 className="text-gray-400">
          {getMonthName(month)} ({filteredTasks.length})
        </h3>
      </div>

      <div className="space-y-6 max-h-[calc(100vh-200px)] overflow-y-auto">
        {sortedDates.map((dateStr) => {
          const date = new Date(dateStr)
          const dayTasks = tasksByDate[dateStr]

          return (
            <div key={dateStr} className="flex gap-4">
              <div className="w-12 flex-shrink-0">
                <span className="text-2xl font-light">{date.getDate()}</span>
              </div>
              <div className="flex-1">
                <p className="text-sm text-gray-400 mb-2">{getDayName(date)}</p>
                <div className="space-y-2">
                  {dayTasks.map((task) => (
                    <button
                      key={task.id}
                      onClick={() => onTaskClick(task)}
                      className="flex items-start gap-2 w-full text-left hover:bg-gray-50 p-1 rounded"
                    >
                      <span className={`mt-1 ${PRIORITY_COLORS[task.priority]}`}>
                        ▸
                      </span>
                      <span
                        className={
                          task.status === 'done'
                            ? 'line-through text-gray-400'
                            : ''
                        }
                      >
                        {task.title}
                      </span>
                    </button>
                  ))}
                </div>
              </div>
            </div>
          )
        })}

        {sortedDates.length === 0 && (
          <p className="text-gray-400 text-center py-8">No tasks found</p>
        )}
      </div>
    </div>
  )
}
