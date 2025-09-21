import {Card, CardContent, CardTitle} from "@/components/ui/card.tsx";
import Track, {type TrackInfo} from "@/components/app/track.tsx";

function TrackList(props: { tracklist: TrackInfo[] | undefined}) {
    return <Card className="flex-1 overflow-y-auto h-full bg-secondary">
        <CardTitle className="text-4xl m-4">Tracks</CardTitle>
        <CardContent>
            <ol className="flex flex-col items-center">
                {props.tracklist?.map(track =>
                    <li key={track.title} className={"min-w-48"}>
                        <Track track={track}/>
                    </li>
                )}
            </ol>
        </CardContent>
    </Card>;
}

export default TrackList;
