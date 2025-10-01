"use client"

import type {TrackInfo} from "@/components/app/track.tsx";
import {createContext, useContext, useEffect, useState} from "react";

interface TracksContextType {
    tracks: TrackInfo[];
    refetch: () => void;
}

const TracksContext = createContext<TracksContextType>({
    tracks: [],
    refetch: () => {
    }
})

export function TracksProvider({children}: { children: React.ReactNode }) {
    const [tracks, setTracks] = useState<TrackInfo[]>([])

    useEffect(() => {
        getData()
    }, []);

    const getData = () => {
        fetch("http://localhost:8080/tracks")
            .then(data => data.json())
            .then((json: TrackInfo[]) => {
                const tracks: TrackInfo[] = json.map((track) => ({
                    id: track.id,
                    title: track.title,
                    author: track.author,
                    coverImg: track.coverImg ?? "https://placehold.net/default.png",
                    lengthSec: track.lengthSec,
                    filePath: track.filePath,
                }));
                setTracks(tracks);
            }).catch((err) => console.error("Failed to fetch tracks:", err));
    }

    const refetch = () => {
        getData()
    }

    return (
        <TracksContext.Provider value={{tracks, refetch}}>
            {children}
        </TracksContext.Provider>
    )
}

export function useTracks() {
    return useContext(TracksContext);
}