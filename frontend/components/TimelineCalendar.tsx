"use client";

import { useRef, useState, useEffect } from "react";

const MONTHS = [
  "JAN", "FEB", "MAR", "APR", "MAY", "JUN",
  "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"
];

type TimelineDay = {
  date: Date;
  label: string;
  hours: number[];
};

function generateDays(centerDate: Date, daysToShow: number): TimelineDay[] {
  const days: TimelineDay[] = [];
  const halfDays = Math.floor(daysToShow / 2);

  for (let i = -halfDays; i <= halfDays; i++) {
    const date = new Date(centerDate);
    date.setDate(date.getDate() + i);

    const month = MONTHS[date.getMonth()];
    const day = String(date.getDate()).padStart(2, "0");
    const year = date.getFullYear();

    days.push({
      date,
      label: `${month} ${day} ${year}`,
      hours: [0, 6, 12, 18],
    });
  }

  return days;
}

function formatHour(hour: number): string {
  return `${String(hour).padStart(2, "0")}:00`;
}

export function TimelineCalendar() {
  const containerRef = useRef<HTMLDivElement>(null);
  const scrollRef = useRef<HTMLDivElement>(null);
  const [centerDate, setCenterDate] = useState(new Date());
  const [isDragging, setIsDragging] = useState(false);
  const [startX, setStartX] = useState(0);
  const [scrollLeft, setScrollLeft] = useState(0);

  const days = generateDays(centerDate, 15);

  // Center scroll on mount
  useEffect(() => {
    if (scrollRef.current) {
      const scrollWidth = scrollRef.current.scrollWidth;
      const clientWidth = scrollRef.current.clientWidth;
      scrollRef.current.scrollLeft = (scrollWidth - clientWidth) / 2;
    }
  }, []);

  // Horizontal scroll with mouse wheel
  useEffect(() => {
    const el = scrollRef.current;
    if (!el) return;

    const handleWheel = (e: WheelEvent) => {
      e.preventDefault();
      el.scrollLeft += e.deltaY + e.deltaX;
    };

    el.addEventListener("wheel", handleWheel, { passive: false });
    return () => el.removeEventListener("wheel", handleWheel);
  }, []);

  // Drag to scroll
  const handleMouseDown = (e: React.MouseEvent) => {
    if (!scrollRef.current) return;
    setIsDragging(true);
    setStartX(e.pageX - scrollRef.current.offsetLeft);
    setScrollLeft(scrollRef.current.scrollLeft);
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!isDragging || !scrollRef.current) return;
    e.preventDefault();
    const x = e.pageX - scrollRef.current.offsetLeft;
    const walk = (x - startX) * 1.5;
    scrollRef.current.scrollLeft = scrollLeft - walk;
  };

  const handleMouseUp = () => {
    setIsDragging(false);
  };

  const handleMouseLeave = () => {
    setIsDragging(false);
  };

  const jumpDays = (delta: number) => {
    setCenterDate((prev) => {
      const next = new Date(prev);
      next.setDate(next.getDate() + delta);
      return next;
    });
  };

  const goToToday = () => {
    setCenterDate(new Date());
    // Re-center scroll
    setTimeout(() => {
      if (scrollRef.current) {
        const scrollWidth = scrollRef.current.scrollWidth;
        const clientWidth = scrollRef.current.clientWidth;
        scrollRef.current.scrollLeft = (scrollWidth - clientWidth) / 2;
      }
    }, 0);
  };

  return (
    <div
      ref={containerRef}
      className="relative w-full bg-slate-900 rounded-xl overflow-hidden select-none"
      style={{ height: "100px" }}
    >
      {/* Navigation buttons */}
      <div className="absolute left-3 top-1/2 -translate-y-1/2 z-10 flex flex-col gap-1">
        <button
          onClick={() => jumpDays(-7)}
          className="w-8 h-8 bg-slate-800/80 hover:bg-slate-700 text-slate-300 rounded-full flex items-center justify-center transition-colors"
          title="Previous week"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <polyline points="15 18 9 12 15 6" />
          </svg>
        </button>
        <button
          onClick={goToToday}
          className="w-8 h-8 bg-slate-800/80 hover:bg-slate-700 text-slate-300 rounded-full flex items-center justify-center transition-colors text-xs font-medium"
          title="Today"
        >
          ●
        </button>
      </div>

      <div className="absolute right-3 top-1/2 -translate-y-1/2 z-10">
        <button
          onClick={() => jumpDays(7)}
          className="w-8 h-8 bg-slate-800/80 hover:bg-slate-700 text-slate-300 rounded-full flex items-center justify-center transition-colors"
          title="Next week"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <polyline points="9 18 15 12 9 6" />
          </svg>
        </button>
      </div>

      {/* Center indicator */}
      <div className="absolute left-1/2 top-0 bottom-0 -translate-x-1/2 z-10 pointer-events-none">
        <div className="relative h-full flex flex-col items-center">
          <div className="w-3 h-3 rounded-full bg-white border-2 border-slate-900 mt-1" />
          <div className="w-0.5 flex-1 bg-white/80" />
        </div>
      </div>

      {/* Scrollable timeline */}
      <div
        ref={scrollRef}
        className={`h-full overflow-x-auto overflow-y-hidden scrollbar-hide px-12 ${
          isDragging ? "cursor-grabbing" : "cursor-grab"
        }`}
        onMouseDown={handleMouseDown}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        onMouseLeave={handleMouseLeave}
        style={{ scrollbarWidth: "none", msOverflowStyle: "none" }}
      >
        <div className="h-full flex items-end pb-3" style={{ width: "max-content" }}>
          {days.map((day, dayIndex) => {
            const isToday =
              day.date.toDateString() === new Date().toDateString();

            return (
              <div key={dayIndex} className="flex flex-col">
                {/* Date label */}
                <div className="h-6 flex items-center px-2">
                  <span
                    className={`text-xs font-medium tracking-wide ${
                      isToday ? "text-white" : "text-slate-400"
                    }`}
                  >
                    {day.label}
                  </span>
                </div>

                {/* Hour ticks */}
                <div className="flex items-end h-12">
                  {day.hours.map((hour, hourIndex) => (
                    <div
                      key={hourIndex}
                      className="flex flex-col items-center"
                      style={{ width: "60px" }}
                    >
                      <span className="text-[10px] text-slate-500 mb-1">
                        {formatHour(hour)}
                      </span>
                      <div
                        className={`w-px ${
                          hour === 0
                            ? "h-6 bg-slate-500"
                            : "h-4 bg-slate-600"
                        }`}
                      />
                    </div>
                  ))}
                </div>

                {/* Minor tick marks */}
                <div className="flex h-3">
                  {Array.from({ length: 24 }).map((_, i) => (
                    <div
                      key={i}
                      className="flex justify-center"
                      style={{ width: "10px" }}
                    >
                      <div className="w-px h-2 bg-slate-700" />
                    </div>
                  ))}
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Gradient edges */}
      <div className="absolute inset-y-0 left-0 w-16 bg-gradient-to-r from-slate-900 to-transparent pointer-events-none" />
      <div className="absolute inset-y-0 right-0 w-16 bg-gradient-to-l from-slate-900 to-transparent pointer-events-none" />
    </div>
  );
}
