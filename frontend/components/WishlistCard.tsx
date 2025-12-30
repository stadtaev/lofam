'use client'

import type { Wishlist, WishlistColor } from '@/lib/types'

interface WishlistCardProps {
  item: Wishlist
  onClick: () => void
  onDelete: () => void
}

const colorClasses: Record<WishlistColor, string> = {
  yellow: 'bg-yellow-100 hover:bg-yellow-200 border-yellow-300',
  pink: 'bg-pink-100 hover:bg-pink-200 border-pink-300',
  green: 'bg-green-100 hover:bg-green-200 border-green-300',
}

export function WishlistCard({ item, onClick, onDelete }: WishlistCardProps) {
  return (
    <div className={`relative w-full rounded-lg border transition-colors ${colorClasses[item.color]}`}>
      <button
        onClick={(e) => {
          e.stopPropagation()
          onDelete()
        }}
        className="absolute top-1 right-1 p-1 text-gray-400 hover:text-gray-600 hover:bg-black/10 rounded"
        aria-label="Delete item"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="14"
          height="14"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <line x1="18" y1="6" x2="6" y2="18" />
          <line x1="6" y1="6" x2="18" y2="18" />
        </svg>
      </button>
      <button
        onClick={onClick}
        className="w-full text-left p-3 pr-7"
      >
        <h3 className="font-medium text-gray-900 truncate">{item.title}</h3>
        {item.content && (
          <p
            className="mt-1 text-sm text-gray-600 line-clamp-3"
            dangerouslySetInnerHTML={{ __html: item.content }}
          />
        )}
      </button>
    </div>
  )
}
