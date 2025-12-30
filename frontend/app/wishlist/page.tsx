"use client";

import { useState, useEffect, useCallback } from "react";
import Link from "next/link";
import { WishlistList } from "@/components/WishlistList";
import { WishlistModal } from "@/components/WishlistModal";
import {
  listWishlists,
  createWishlist,
  updateWishlist,
  deleteWishlist,
} from "@/lib/api";
import type { Wishlist, CreateWishlistRequest } from "@/lib/types";

export default function WishlistPage() {
  const [items, setItems] = useState<Wishlist[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [showModal, setShowModal] = useState(false);
  const [editingItem, setEditingItem] = useState<Wishlist | null>(null);

  const fetchItems = useCallback(async () => {
    try {
      setLoading(true);
      const data = await listWishlists();
      setItems(data || []);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch wishlist");
      setItems([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchItems();
  }, [fetchItems]);

  const handleAddItem = () => {
    setEditingItem(null);
    setShowModal(true);
  };

  const handleEditItem = (item: Wishlist) => {
    setEditingItem(item);
    setShowModal(true);
  };

  const handleSaveItem = async (data: CreateWishlistRequest) => {
    try {
      if (editingItem) {
        await updateWishlist(editingItem.id, {
          title: data.title,
          content: data.content ?? "",
          color: data.color,
        });
      } else {
        await createWishlist(data);
      }
      setShowModal(false);
      setEditingItem(null);
      fetchItems();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save item");
    }
  };

  const handleDeleteItem = async () => {
    if (!editingItem) return;
    try {
      await deleteWishlist(editingItem.id);
      setShowModal(false);
      setEditingItem(null);
      fetchItems();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete item");
    }
  };

  const handleDeleteItemDirectly = async (item: Wishlist) => {
    try {
      await deleteWishlist(item.id);
      fetchItems();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete item");
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-gray-400">Loading...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white p-4 lg:p-8">
      {error && (
        <div className="fixed top-4 right-4 bg-red-100 text-red-600 px-4 py-2 rounded-lg z-40">
          {error}
          <button onClick={() => setError(null)} className="ml-2">
            ✕
          </button>
        </div>
      )}

      <div className="max-w-2xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-light">Wishlist</h1>
          <Link
            href="/"
            className="text-gray-500 hover:text-gray-700 text-sm"
          >
            ← Back to Tasks
          </Link>
        </div>

        <div className="bg-gray-50 rounded-lg p-4">
          <WishlistList
            items={items}
            onAdd={handleAddItem}
            onEdit={handleEditItem}
            onDelete={handleDeleteItemDirectly}
          />
        </div>
      </div>

      {showModal && (
        <WishlistModal
          item={editingItem}
          onSave={handleSaveItem}
          onDelete={editingItem ? handleDeleteItem : undefined}
          onClose={() => {
            setShowModal(false);
            setEditingItem(null);
          }}
        />
      )}
    </div>
  );
}
