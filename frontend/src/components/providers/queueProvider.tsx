"use client"
import React, { createContext, useContext, useState } from "react";
import type { TrackInfo } from "@/components/app/track";

interface QueueContextType {
    queue: TrackInfo[];
    current: TrackInfo | null;
    addToQueue: (track: TrackInfo) => void;
    removeFromQueue: (track: TrackInfo) => void;
    clearQueue: () => void;
    setCurrent: (track: TrackInfo) => void;
}

const QueueContext = createContext<QueueContextType>({
    queue: [],
    current: null,
    addToQueue: () => {},
    removeFromQueue: () => {},
    clearQueue: () => {},
    setCurrent: () => {},
});

export function QueueProvider({ children }: { children: React.ReactNode }) {
    const [queue, setQueue] = useState<TrackInfo[]>([]);
    const [current, setCurrent] = useState<TrackInfo | null>(null);

    const addToQueue = (track: TrackInfo) => {
        setQueue((prev) => [...prev, track]);
    };

    const removeFromQueue = (track: TrackInfo) => {
        setQueue((prev) => prev.filter((t) => t.filePath !== track.filePath));
        if (current?.filePath === track.filePath) {
            setCurrent(null); // clear current if it was removed
        }
    };

    const clearQueue = () => {
        setQueue([]);
        setCurrent(null);
    };

    return (
        <QueueContext.Provider
            value={{ queue, current, addToQueue, removeFromQueue, clearQueue, setCurrent }}
        >
            {children}
        </QueueContext.Provider>
    );
}

export function useQueue() {
    return useContext(QueueContext);
}
