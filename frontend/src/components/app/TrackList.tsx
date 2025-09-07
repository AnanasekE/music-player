import {Card, CardContent, CardTitle} from "@/components/ui/card.tsx";
import Track, {type TrackInfo} from "@/components/app/track.tsx";

function TrackList(props: { tracklist: TrackInfo[] | undefined, setCurrentTrack: (track: TrackInfo) => void }) {
    return <Card className="flex-1 overflow-y-auto h-full">
        <CardTitle className="text-4xl m-4">Tracks</CardTitle>
        <CardContent>
            <ol className="flex flex-col">
                {props.tracklist?.map(track =>
                    <li key={track.title}>
                        <Track track={track} onPlay={() => props.setCurrentTrack(track)}/>
                    </li>
                )}
            </ol>
        </CardContent>
    </Card>;
}

export default TrackList;
