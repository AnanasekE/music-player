import {Card} from "@/components/ui/card.tsx";
import {Skeleton} from "@/components/ui/skeleton.tsx";
import {Button} from "@/components/ui/button.tsx";
import {useQueue} from "@/components/providers/queueProvider.tsx";

export interface TrackInfo {
    id: number;
    title: string;
    author: string;
    coverImg: string | null;
    lengthSec: number;
    filePath: string;
}

interface SongProps {
    track: TrackInfo;
}

const Track = ({track}: SongProps) => {
    const {setCurrent} = useQueue()

    return (
        <Card className="max-w-lg m-2 p-2 flex flex-row items-center justify-between min-w-96">
            <div className="flex flex-row items-center truncate">
                <div className="ml-0.5 mr-2 min-w-10 min-h-10 max-w-10 max-h-10">
                    {track.coverImg ? (
                        <img src={track.coverImg} alt="image" className="w-10 h-10"/>
                    ) : (
                        <Skeleton className="w-10 h-10"/>
                    )}Current song:
                </div>
                <div className="flex flex-col truncate text-left">
                    <h3 className="truncate">{track.title}</h3>
                    <h3 className="italic">{track.author}</h3>
                </div>
            </div>
            <div className={"flex flex-row items-center ml-2"}>
                <h2>{Math.floor(track.lengthSec / 60)}:{(track.lengthSec % 60).toString().padStart(2, "0")}</h2>
                <Button className={"ml-2"} onClick={() => setCurrent(track)}>Play</Button>
            </div>
        </Card>
    );
};

export default Track;
