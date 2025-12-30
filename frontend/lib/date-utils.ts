import {
  format,
  isSameDay as dateFnsIsSameDay,
  getDaysInMonth as dateFnsGetDaysInMonth,
  startOfMonth,
  getDay,
  eachDayOfInterval,
  endOfMonth,
} from 'date-fns'

export function getDaysInMonth(year: number, month: number): Date[] {
  const start = new Date(year, month, 1)
  const end = endOfMonth(start)
  return eachDayOfInterval({ start, end })
}

export function getCalendarDays(year: number, month: number): (Date | null)[] {
  const days = getDaysInMonth(year, month)
  const firstDay = getDay(days[0])
  // Adjust for Monday start (0 = Monday, 6 = Sunday)
  const startPadding = firstDay === 0 ? 6 : firstDay - 1

  const calendar: (Date | null)[] = []
  for (let i = 0; i < startPadding; i++) {
    calendar.push(null)
  }
  calendar.push(...days)

  return calendar
}

export function formatDateKey(date: Date): string {
  return format(date, 'yyyy-MM-dd')
}

export function isSameDay(date1: Date, date2: Date): boolean {
  return dateFnsIsSameDay(date1, date2)
}

export function getMonthName(month: number): string {
  return format(new Date(2000, month, 1), 'MMMM')
}

export function getDayName(date: Date): string {
  return format(date, 'EEEE')
}
