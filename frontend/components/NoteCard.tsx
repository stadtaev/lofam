'use client'

import type { Note, NoteColor } from '@/lib/types'

interface NoteCardProps {
  note: Note
  onClick: () => void
}

const colorClasses: Record<NoteColor, string> = {
  yellow: 'bg-yellow-100 hover:bg-yellow-200 border-yellow-300',
  pink: 'bg-pink-100 hover:bg-pink-200 border-pink-300',
  green: 'bg-green-100 hover:bg-green-200 border-green-300',
}

export function NoteCard({ note, onClick }: NoteCardProps) {
  return (
    <button
      onClick={onClick}
      className={`w-full text-left p-3 rounded-lg border transition-colors ${colorClasses[note.color]}`}
    >
      <h3 className="font-medium text-gray-900 truncate">{note.title}</h3>
      {note.content && (
        <p
          className="mt-1 text-sm text-gray-600 line-clamp-3"
          dangerouslySetInnerHTML={{ __html: note.content }}
        />
      )}
    </button>
  )
}
