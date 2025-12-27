'use client'

import type { Task } from '@/lib/types'
import { formatDateKey } from '@/lib/date-utils'

interface TodaySectionProps {
  tasks: Task[]
  onAddTask: () => void
}

export function TodaySection({ tasks, onAddTask }: TodaySectionProps) {
  const today = formatDateKey(new Date())
  const todayTasks = tasks.filter(
    (task) => task.dueDate && task.dueDate.split('T')[0] === today
  )

  return (
    <div className="mt-8 pt-6 border-t">
      <div className="flex items-center justify-between mb-4">
        <h3 className="font-medium">Today</h3>
        <button
          onClick={onAddTask}
          className="flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700"
        >
          <span>+</span>
          <span>ADD NEW EVENT</span>
        </button>
      </div>

      {todayTasks.length === 0 ? (
        <p className="text-sm text-gray-400">No events are planned</p>
      ) : (
        <div className="space-y-2">
          {todayTasks.map((task) => (
            <div key={task.id} className="text-sm">
              <span
                className={
                  task.status === 'done' ? 'line-through text-gray-400' : ''
                }
              >
                {task.title}
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
