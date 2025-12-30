'use client'

import type { Wishlist } from '@/lib/types'
import { WishlistCard } from './WishlistCard'

interface WishlistListProps {
  items: Wishlist[]
  onAdd: () => void
  onEdit: (item: Wishlist) => void
  onDelete: (item: Wishlist) => void
}

export function WishlistList({ items, onAdd, onEdit, onDelete }: WishlistListProps) {
  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold text-gray-900">Wishlist</h2>
        <button
          onClick={onAdd}
          className="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-full transition-colors"
          aria-label="Add item"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="20"
            height="20"
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
      </div>

      <div className="flex-1 overflow-y-auto space-y-2">
        {items.length === 0 ? (
          <p className="text-sm text-gray-500 text-center py-4">
            No wishlist items yet
          </p>
        ) : (
          items.map((item) => (
            <WishlistCard
              key={item.id}
              item={item}
              onClick={() => onEdit(item)}
              onDelete={() => onDelete(item)}
            />
          ))
        )}
      </div>
    </div>
  )
}
