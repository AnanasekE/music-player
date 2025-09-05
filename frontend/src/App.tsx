import {Card, CardContent, CardTitle} from "@/components/ui/card";
import "./index.css";
import {useEffect, useState} from "react";
import Track, {type TrackInfo} from "@/components/app/track.tsx";

interface FileInfo {
    FileName: string;
    Path: string;
}

export function App() {
    const [tracklist, setTracklist] = useState<TrackInfo[]>();
    const [currentTrack, setCurrentTrack] = useState<TrackInfo>();

    useEffect(() => {
        fetch("http://localhost:8080/tracks")
            .then((res) => res.json())
            .then((json: FileInfo[]) => {
                const tracks: TrackInfo[] = json.map((track) => ({
                    title: track.FileName.replace(/\.[^/.]+$/, ""),
                    author: "Author",
                    imgSrc: "https://placehold.net/default.png",
                    lengthSec: 231,
                    path: track.Path,
                    fileName: track.FileName,
                }));
                setTracklist(tracks);
            })
            .catch((err) => console.error("Failed to fetch tracks:", err));
    }, []);

    useEffect(() => {
        document.documentElement.classList.add("dark");
    }, []);

    return (
        <div className="h-screen flex flex-col p-8 text-center relative z-10 bg-background text-foreground">
            <Card className="flex-1 overflow-y-auto max-w-lg">
                <CardTitle className="text-4xl m-4">Tracks</CardTitle>
                <CardContent>
                    <ol className="flex flex-col">
                        {tracklist?.map((track) => (
                            <li key={track.fileName}>
                                <Track track={track} onPlay={() => setCurrentTrack(track)}/>
                            </li>
                        ))}
                    </ol>
                </CardContent>
            </Card>

            <Card className="mt-4">
                <CardTitle className="text-xl m-2">Audio Player</CardTitle>
                <CardContent>
                    <h2 className="mb-2">
                        Currently playing: {currentTrack?.title ?? "Nothing selected"}
                    </h2>
                    {currentTrack && (
                        <audio
                            controls
                            autoPlay
                            className="w-full"
                            src={`http://localhost:8080/${currentTrack.path}`}
                        ></audio>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}

export default App;
