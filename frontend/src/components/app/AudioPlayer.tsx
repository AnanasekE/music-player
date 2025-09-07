import {Card, CardContent, CardTitle} from "@/components/ui/card.tsx";
import type {TrackInfo} from "@/components/app/track.tsx";

function AudioPlayer(props: {
    currentTrack: TrackInfo | undefined
}) {

    const title = props.currentTrack?.title
    const modifiedTitle = `${title} - ${props.currentTrack?.author}`

    return <Card className="mt-4">
        <CardTitle className="text-xl m-2">Audio Player</CardTitle>
        <CardContent>
            <h2 className="mb-2">
                Currently
                playing: {props.currentTrack?.title && props.currentTrack.author ? modifiedTitle : "Nothing selected"}
            </h2>
            {props.currentTrack && (
                <audio
                    controls
                    autoPlay
                    className="w-full"
                    src={`http://localhost:8080/music/${props.currentTrack.filePath}`}
                ></audio>
            )}
        </CardContent>
    </Card>;
}

export default AudioPlayer;