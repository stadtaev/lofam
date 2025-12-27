export function getDaysInMonth(year: number, month: number): Date[] {
  const days: Date[] = []
  const date = new Date(year, month, 1)
  while (date.getMonth() === month) {
    days.push(new Date(date))
    date.setDate(date.getDate() + 1)
  }
  return days
}

export function getCalendarDays(year: number, month: number): (Date | null)[] {
  const days = getDaysInMonth(year, month)
  const firstDay = days[0].getDay()
  const startPadding = firstDay === 0 ? 6 : firstDay - 1

  const calendar: (Date | null)[] = []
  for (let i = 0; i < startPadding; i++) {
    calendar.push(null)
  }
  calendar.push(...days)

  return calendar
}

export function formatDateKey(date: Date): string {
  return date.toISOString().split('T')[0]
}

export function isSameDay(date1: Date, date2: Date): boolean {
  return formatDateKey(date1) === formatDateKey(date2)
}

export function getMonthName(month: number): string {
  return new Date(2000, month).toLocaleString('en', { month: 'long' })
}

export function getDayName(date: Date): string {
  return date.toLocaleString('en', { weekday: 'long' })
}
