import {Card, CardContent, CardTitle} from "@/components/ui/card.tsx";
import {useQueue} from "@/components/providers/queueProvider.tsx";
import {useEffect, useRef} from "react";

function AudioPlayer() {
    const {current, playNext, queue, clearQueue} = useQueue()
    const audioRef = useRef<HTMLAudioElement | null>(null);

    useEffect(() => {
        if (audioRef.current && current) {
            audioRef.current.src = current.filePath;
            audioRef.current.play().catch((err) => {
                console.error("Failed to play:", err);
            });
        }
    }, [current]);

    const title = current?.title
    const modifiedTitle = `${title} - ${current?.author}`

    return <Card className="mt-1 bg-secondary">
        <CardTitle className="text-xl m-2">Audio Player</CardTitle>
        <CardContent>
            <h2 className="mb-2">
                Currently
                playing: {current?.title && current.author ? modifiedTitle : "Nothing selected"}
            </h2>
            {current && (
                <audio
                    controls
                    autoPlay
                    className="w-full"
                    src={`http://localhost:8080/music/${current.filePath}`}
                    onEnded={() => {
                        if (queue.length > 0) {
                            playNext()
                        } else {
                            clearQueue()
                        }
                    }}
                ></audio>
            )}
        </CardContent>
    </Card>;
}

export default AudioPlayer;