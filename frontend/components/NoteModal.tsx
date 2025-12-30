'use client'

import { useState, useEffect, useRef } from 'react'
import type { Note, NoteColor, CreateNoteRequest } from '@/lib/types'

interface NoteModalProps {
  note?: Note | null
  onSave: (data: CreateNoteRequest) => void
  onDelete?: () => void
  onClose: () => void
}

const colorOptions: { value: NoteColor; label: string; className: string }[] = [
  { value: 'yellow', label: 'Yellow', className: 'bg-yellow-200 border-yellow-400' },
  { value: 'pink', label: 'Pink', className: 'bg-pink-200 border-pink-400' },
  { value: 'green', label: 'Green', className: 'bg-green-200 border-green-400' },
]

export function NoteModal({ note, onSave, onDelete, onClose }: NoteModalProps) {
  const [title, setTitle] = useState(note?.title ?? '')
  const [color, setColor] = useState<NoteColor>(note?.color ?? 'yellow')
  const editorRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (editorRef.current && note?.content) {
      editorRef.current.innerHTML = note.content
    }
  }, [note?.content])

  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', handleEsc)
    return () => window.removeEventListener('keydown', handleEsc)
  }, [onClose])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSave({
      title,
      content: editorRef.current?.innerHTML ?? '',
      color,
    })
  }

  const execCommand = (command: string, value?: string) => {
    document.execCommand(command, false, value)
    editorRef.current?.focus()
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-lg mx-4">
        <form onSubmit={handleSubmit}>
          <div className="flex items-center justify-between p-4 border-b">
            <h2 className="text-lg font-medium">
              {note ? 'Edit Note' : 'New Note'}
            </h2>
            <button
              type="button"
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600"
            >
              ✕
            </button>
          </div>

          <div className="p-4 space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Title
              </label>
              <input
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                required
                autoFocus
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-200"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Color
              </label>
              <div className="flex gap-2">
                {colorOptions.map((opt) => (
                  <button
                    key={opt.value}
                    type="button"
                    onClick={() => setColor(opt.value)}
                    className={`w-8 h-8 rounded-full border-2 transition-transform ${opt.className} ${
                      color === opt.value ? 'scale-110 ring-2 ring-offset-2 ring-gray-400' : ''
                    }`}
                    aria-label={opt.label}
                  />
                ))}
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Content
              </label>
              <div className="border border-gray-300 rounded-md overflow-hidden">
                <div className="flex gap-1 p-1 border-b border-gray-200 bg-gray-50">
                  <button
                    type="button"
                    onClick={() => execCommand('bold')}
                    className="px-2 py-1 text-sm font-bold hover:bg-gray-200 rounded"
                    title="Bold (Ctrl+B)"
                  >
                    B
                  </button>
                  <button
                    type="button"
                    onClick={() => execCommand('italic')}
                    className="px-2 py-1 text-sm italic hover:bg-gray-200 rounded"
                    title="Italic (Ctrl+I)"
                  >
                    I
                  </button>
                  <div className="w-px bg-gray-300 mx-1" />
                  <button
                    type="button"
                    onClick={() => execCommand('insertUnorderedList')}
                    className="px-2 py-1 text-sm hover:bg-gray-200 rounded"
                    title="Bullet list"
                  >
                    • List
                  </button>
                  <button
                    type="button"
                    onClick={() => execCommand('insertOrderedList')}
                    className="px-2 py-1 text-sm hover:bg-gray-200 rounded"
                    title="Numbered list"
                  >
                    1. List
                  </button>
                </div>
                <div
                  ref={editorRef}
                  contentEditable
                  className="min-h-[120px] p-3 focus:outline-none prose prose-sm max-w-none"
                  onKeyDown={(e) => {
                    if (e.key === 'b' && (e.ctrlKey || e.metaKey)) {
                      e.preventDefault()
                      execCommand('bold')
                    }
                    if (e.key === 'i' && (e.ctrlKey || e.metaKey)) {
                      e.preventDefault()
                      execCommand('italic')
                    }
                  }}
                />
              </div>
            </div>
          </div>

          <div className="flex items-center justify-between p-4 border-t bg-gray-50 rounded-b-lg">
            {note && onDelete ? (
              <button
                type="button"
                onClick={onDelete}
                className="px-4 py-2 text-red-600 hover:bg-red-50 rounded-md"
              >
                Delete
              </button>
            ) : (
              <div />
            )}
            <div className="flex gap-2">
              <button
                type="button"
                onClick={onClose}
                className="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-md"
              >
                Cancel
              </button>
              <button
                type="submit"
                className="px-4 py-2 bg-amber-500 text-white rounded-md hover:bg-amber-600"
              >
                {note ? 'Save' : 'Create'}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  )
}
