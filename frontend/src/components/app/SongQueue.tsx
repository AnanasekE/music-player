import {Card, CardContent, CardTitle} from "@/components/ui/card.tsx";
import {useQueue} from "@/components/providers/queueProvider.tsx";
import {Button} from "@/components/ui/button.tsx";
import Track from "@/components/app/track.tsx";

function SongQueue() {
    const {queue, clearQueue, current} = useQueue()
    return <Card className={"flex-1 overflow-y-auto h-full bg-secondary"}>
        <CardTitle className="text-4xl m-4">Song Queue</CardTitle>
        <CardContent>
            <div className={"flex flex-row justify-center"}>
                <Button onClick={clearQueue}>Clear Queue</Button>
            </div>
            {current && (
                <>
                    <h2>Currently Playing:</h2>
                    <Track track={current}/>
                </>
            )}
            <br/>
            <div>
                {queue.map(track =>
                    <Track track={track}/>
                )}
            </div>
        </CardContent>

    </Card>
}

export default SongQueue
