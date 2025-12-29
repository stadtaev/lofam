'use client'

import { getCalendarDays, getMonthName, isSameDay, formatDateKey } from '@/lib/date-utils'
import type { Task } from '@/lib/types'

interface CalendarProps {
  year: number
  month: number
  tasks: Task[]
  selectedDate: Date | null
  onDateSelect: (date: Date) => void
  onPrevMonth: () => void
  onNextMonth: () => void
  onToday: () => void
}

const WEEKDAYS = ['M', 'T', 'W', 'T', 'F', 'S', 'S']

export function Calendar({
  year,
  month,
  tasks,
  selectedDate,
  onDateSelect,
  onPrevMonth,
  onNextMonth,
  onToday,
}: CalendarProps) {
  const days = getCalendarDays(year, month)
  const today = new Date()

  const tasksByDate = tasks.reduce((acc, task) => {
    if (task.dueDate) {
      const key = task.dueDate.split('T')[0]
      if (!acc[key]) acc[key] = []
      acc[key].push(task)
    }
    return acc
  }, {} as Record<string, Task[]>)

  return (
    <div className="w-full lg:w-80">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-3xl font-light">{getMonthName(month)}</h2>
          <p className="text-gray-400">{year}</p>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={onToday}
            className="px-3 py-1 text-sm hover:bg-gray-100 rounded text-gray-500"
          >
            Today
          </button>
          <button
            onClick={onPrevMonth}
            className="p-2 hover:bg-gray-100 rounded-full text-gray-400"
          >
            ‹
          </button>
          <button
            onClick={onNextMonth}
            className="p-2 hover:bg-gray-100 rounded-full text-gray-400"
          >
            ›
          </button>
        </div>
      </div>

      <div className="grid grid-cols-7 gap-1 mb-2">
        {WEEKDAYS.map((day, i) => (
          <div key={i} className="text-center text-sm text-gray-400 py-2">
            {day}
          </div>
        ))}
      </div>

      <div className="grid grid-cols-7 gap-1">
        {days.map((date, i) => {
          if (!date) {
            return <div key={`empty-${i}`} className="p-2" />
          }

          const dateKey = formatDateKey(date)
          const hasTasks = !!tasksByDate[dateKey]
          const isToday = isSameDay(date, today)
          const isSelected = selectedDate && isSameDay(date, selectedDate)

          return (
            <button
              key={dateKey}
              onClick={() => onDateSelect(date)}
              className={`
                p-2 text-center rounded-full text-sm relative
                ${isToday ? 'bg-red-100 text-red-600' : ''}
                ${isSelected ? 'ring-2 ring-red-400' : ''}
                ${hasTasks && !isToday ? 'bg-yellow-50' : ''}
                hover:bg-gray-100
              `}
            >
              {date.getDate()}
              {hasTasks && (
                <span className="absolute bottom-0 left-1/2 -translate-x-1/2 w-1 h-1 bg-red-400 rounded-full" />
              )}
            </button>
          )
        })}
      </div>
    </div>
  )
}
