import {Card, CardContent, CardTitle} from "@/components/ui/card.tsx";
import {useQueue} from "@/components/providers/queueProvider.tsx";

function AudioPlayer() {
    const {current} = useQueue()

    const title = current?.title
    const modifiedTitle = `${title} - ${current?.author}`

    return <Card className="mt-4">
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
                ></audio>
            )}
        </CardContent>
    </Card>;
}

export default AudioPlayer;