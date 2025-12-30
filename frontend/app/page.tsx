"use client";

import { useState, useEffect, useCallback } from "react";
import { Calendar } from "@/components/Calendar";
import { TaskList } from "@/components/TaskList";
import { TaskModal } from "@/components/TaskModal";
import { TodaySection } from "@/components/TodaySection";
import { NoteList } from "@/components/NoteList";
import { NoteModal } from "@/components/NoteModal";
import { Timeline } from "@/components/Timeline";
import {
  listTasks,
  createTask,
  updateTask,
  deleteTask,
  listNotes,
  createNote,
  updateNote,
  deleteNote,
} from "@/lib/api";
import type { Task, CreateTaskRequest, Note, CreateNoteRequest } from "@/lib/types";

export default function Home() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const today = new Date();
  const [year, setYear] = useState(today.getFullYear());
  const [month, setMonth] = useState(today.getMonth());
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);
  const [searchQuery, setSearchQuery] = useState("");

  const [showTaskModal, setShowTaskModal] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);

  const [showNoteModal, setShowNoteModal] = useState(false);
  const [editingNote, setEditingNote] = useState<Note | null>(null);

  const fetchTasks = useCallback(async () => {
    try {
      const data = await listTasks();
      setTasks(data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch tasks");
      setTasks([]);
    }
  }, []);

  const fetchNotes = useCallback(async () => {
    try {
      const data = await listNotes();
      setNotes(data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch notes");
      setNotes([]);
    }
  }, []);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      await Promise.all([fetchTasks(), fetchNotes()]);
      setLoading(false);
    };
    fetchData();
  }, [fetchTasks, fetchNotes]);

  const handlePrevMonth = () => {
    if (month === 0) {
      setMonth(11);
      setYear(year - 1);
    } else {
      setMonth(month - 1);
    }
  };

  const handleNextMonth = () => {
    if (month === 11) {
      setMonth(0);
      setYear(year + 1);
    } else {
      setMonth(month + 1);
    }
  };

  const handleToday = () => {
    const now = new Date();
    setYear(now.getFullYear());
    setMonth(now.getMonth());
  };

  const handleDateSelect = (date: Date) => {
    setSelectedDate(date);
    // Sync calendar month/year to show the selected date
    setYear(date.getFullYear());
    setMonth(date.getMonth());
  };

  const handleTaskClick = (task: Task) => {
    setEditingTask(task);
    setShowTaskModal(true);
  };

  const handleAddTask = () => {
    setEditingTask(null);
    setShowTaskModal(true);
  };

  const handleSaveTask = async (data: CreateTaskRequest) => {
    try {
      if (editingTask) {
        await updateTask(editingTask.id, data);
      } else {
        await createTask(data);
      }
      setShowTaskModal(false);
      setEditingTask(null);
      fetchTasks();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save task");
    }
  };

  const handleDeleteTask = async () => {
    if (!editingTask) return;
    try {
      await deleteTask(editingTask.id);
      setShowTaskModal(false);
      setEditingTask(null);
      fetchTasks();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete task");
    }
  };

  const handleAddNote = () => {
    setEditingNote(null);
    setShowNoteModal(true);
  };

  const handleEditNote = (note: Note) => {
    setEditingNote(note);
    setShowNoteModal(true);
  };

  const handleSaveNote = async (data: CreateNoteRequest) => {
    try {
      if (editingNote) {
        await updateNote(editingNote.id, {
          title: data.title,
          content: data.content ?? "",
          color: data.color,
        });
      } else {
        await createNote(data);
      }
      setShowNoteModal(false);
      setEditingNote(null);
      fetchNotes();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save note");
    }
  };

  const handleDeleteNote = async () => {
    if (!editingNote) return;
    try {
      await deleteNote(editingNote.id);
      setShowNoteModal(false);
      setEditingNote(null);
      fetchNotes();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete note");
    }
  };

  const handleDeleteNoteDirectly = async (note: Note) => {
    try {
      await deleteNote(note.id);
      fetchNotes();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete note");
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

      <div className="max-w-7xl mx-auto">
        {/* Horizontal Timeline */}
        <div className="mb-6">
          <Timeline
            selectedDate={selectedDate}
            onDateSelect={handleDateSelect}
            tasks={tasks}
          />
        </div>

        <div className="flex flex-col xl:flex-row gap-6">
          {/* Left: Calendar + Today */}
          <div className="shrink-0">
            <Calendar
              year={year}
              month={month}
              tasks={tasks}
              selectedDate={selectedDate}
              onDateSelect={handleDateSelect}
              onPrevMonth={handlePrevMonth}
              onNextMonth={handleNextMonth}
              onToday={handleToday}
            />
            <TodaySection tasks={tasks} onAddTask={handleAddTask} />
          </div>

          {/* Middle: Task List */}
          <div className="flex-1 min-w-0">
            <TaskList
              tasks={tasks}
              month={month}
              searchQuery={searchQuery}
              onSearchChange={setSearchQuery}
              onTaskClick={handleTaskClick}
            />
          </div>

          {/* Right: Notes Sidebar (desktop) / Below (mobile) */}
          <div className="w-full xl:w-72 shrink-0">
            <div className="bg-gray-50 rounded-lg p-4 h-full max-h-[600px]">
              <NoteList
                notes={notes}
                onAdd={handleAddNote}
                onEdit={handleEditNote}
                onDelete={handleDeleteNoteDirectly}
              />
            </div>
          </div>
        </div>
      </div>

      {showTaskModal && (
        <TaskModal
          task={editingTask}
          initialDate={selectedDate}
          onSave={handleSaveTask}
          onDelete={editingTask ? handleDeleteTask : undefined}
          onClose={() => {
            setShowTaskModal(false);
            setEditingTask(null);
          }}
        />
      )}

      {showNoteModal && (
        <NoteModal
          note={editingNote}
          onSave={handleSaveNote}
          onDelete={editingNote ? handleDeleteNote : undefined}
          onClose={() => {
            setShowNoteModal(false);
            setEditingNote(null);
          }}
        />
      )}
    </div>
  );
}
