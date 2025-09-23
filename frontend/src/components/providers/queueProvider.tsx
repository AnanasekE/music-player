"use client"
import React, {createContext, useContext, useState} from "react";
import type {TrackInfo} from "@/components/app/track";

interface QueueContextType {
    queue: TrackInfo[];
    current: TrackInfo | null;
    addToQueue: (track: TrackInfo) => void;
    removeAllFromQueue: (track: TrackInfo) => void;
    removeByIndex: (index: number) => void;
    clearQueue: () => void;
    setCurrent: (track: TrackInfo) => void;
    playNext: () => void;
}

const QueueContext = createContext<QueueContextType>({
    queue: [],
    current: null,
    addToQueue: () => {
    },
    removeAllFromQueue: () => {
    },
    clearQueue: () => {
    },
    setCurrent: () => {
    },
    playNext: () => {
    },
    removeByIndex: () => {

    }
});

export function QueueProvider({children}: { children: React.ReactNode }) {
    const [queue, setQueue] = useState<TrackInfo[]>([]);
    const [current, setCurrent] = useState<TrackInfo | null>(null);

    const addToQueue = (track: TrackInfo) => {
        if (queue.length === 0 && !current) {
            setCurrent(track)
        } else {
            setQueue((prev) => [...prev, track]);
        }
    };

    const removeAllFromQueue = (track: TrackInfo) => {
        setQueue((prev) => prev.filter((t) => t.id !== track.id));
        if (current?.id === track.id) {
            setCurrent(null); // clear current if it was removed
        }
    };

    const clearQueue = () => {
        setQueue([]);
        setCurrent(null);
    };
    const playNext = () => {
        if (queue.length === 0) return
        setCurrent(queue[0] as TrackInfo) // force for typescript as cannot be undefined
        setQueue(queue.slice(1, queue.length))
    }

    const removeByIndex = (index:number) => {
        setQueue((prev) => prev.filter((_, i) => i !== index));
    }

    return (
        <QueueContext.Provider
            value={{queue, current, addToQueue, removeAllFromQueue, clearQueue, setCurrent, playNext, removeByIndex}}
        >
            {children}
        </QueueContext.Provider>
    );
}

export function useQueue() {
    return useContext(QueueContext);
}
