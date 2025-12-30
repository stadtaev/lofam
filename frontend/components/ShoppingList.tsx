'use client'

import { useState } from 'react'
import type { ShoppingItem } from '@/lib/types'

interface ShoppingListProps {
  items: ShoppingItem[]
  onAdd: (title: string) => void
  onDelete: (id: number) => void
}

export function ShoppingList({ items, onAdd, onDelete }: ShoppingListProps) {
  const [newItem, setNewItem] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (newItem.trim()) {
      onAdd(newItem.trim())
      setNewItem('')
    }
  }

  return (
    <div className="bg-yellow-50 rounded-lg p-4 h-full max-h-[400px] flex flex-col">
      <h2 className="text-lg font-semibold text-gray-900 mb-3">Shopping List</h2>

      {/* Add item form */}
      <form onSubmit={handleSubmit} className="flex gap-2 mb-3">
        <input
          type="text"
          value={newItem}
          onChange={(e) => setNewItem(e.target.value)}
          placeholder="Add item..."
          className="flex-1 px-3 py-2 text-sm border border-yellow-200 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-yellow-400"
        />
        <button
          type="submit"
          className="px-3 py-2 bg-yellow-400 hover:bg-yellow-500 text-yellow-900 rounded-lg transition-colors"
          aria-label="Add item"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <line x1="12" y1="5" x2="12" y2="19" />
            <line x1="5" y1="12" x2="19" y2="12" />
          </svg>
        </button>
      </form>

      {/* Items table */}
      <div className="flex-1 overflow-y-auto">
        {items.length === 0 ? (
          <p className="text-sm text-gray-500 text-center py-4">
            No items yet
          </p>
        ) : (
          <table className="w-full">
            <tbody>
              {items.map((item) => (
                <tr
                  key={item.id}
                  className="border-b border-yellow-200 last:border-b-0"
                >
                  <td className="py-2 text-sm text-gray-800">{item.title}</td>
                  <td className="py-2 w-10 text-right">
                    <button
                      onClick={() => onDelete(item.id)}
                      className="p-1 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded transition-colors"
                      aria-label="Remove item"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="2"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                      >
                        <line x1="5" y1="12" x2="19" y2="12" />
                      </svg>
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
