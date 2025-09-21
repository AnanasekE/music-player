import "./index.css";
import {useEffect, useState} from "react";
import {type TrackInfo} from "@/components/app/track.tsx";
import TrackList from "@/components/app/TrackList.tsx";
import AudioPlayer from "@/components/app/AudioPlayer.tsx";
import Playlists from "@/components/app/Playlists.tsx";
import AddSong from "@/components/app/AddSong.tsx";
import {ThemeSwitcher} from "@/components/app/ThemeSwitcher.tsx";
import {NavigationMenu, NavigationMenuLink} from "@/components/ui/navigation-menu.tsx";
import {Dialog, DialogContent, DialogTitle, DialogTrigger} from "@/components/ui/dialog.tsx";
import SongQueue from "@/components/app/SongQueue.tsx";


export function App() {
    const [tracklist, setTracklist] = useState<TrackInfo[]>();

    const fetchTracks = () => {
        fetch("http://localhost:8080/tracks")
            .then((res) => res.json())
            .then((json: TrackInfo[]) => {
                const tracks: TrackInfo[] = json.map((track) => ({
                    id: track.id,
                    title: track.title,
                    author: track.author,
                    coverImg: track.coverImg ?? "https://placehold.net/default.png",
                    lengthSec: track.lengthSec,
                    filePath: track.filePath,
                }));
                setTracklist(tracks);
            })
            .catch((err) => console.error("Failed to fetch tracks:", err));
    };

    useEffect(() => {
        fetchTracks();
        document.documentElement.classList.add("dark");
    }, []);

    return (
        <div className="h-screen flex flex-col px-8 pb-8 text-center relative z-10 bg-background text-foreground w-screen">
            <ThemeSwitcher/>
            <NavigationMenu className={"max-h-10 p-6"}>
                <NavigationMenuLink>
                    <Dialog>
                        <DialogTitle className={"sr-only"}></DialogTitle>
                        <DialogTrigger>Add Song</DialogTrigger>
                        <DialogContent className="bg-transparent shadow-none p-0 border-0">
                            <AddSong onAdd={fetchTracks}/>
                        </DialogContent>
                    </Dialog>
                </NavigationMenuLink>
                <NavigationMenuLink>Themes</NavigationMenuLink>
            </NavigationMenu>
            <div className="flex-1 grid grid-cols-3 gap-4 overflow-hidden w-full">
                <div className="overflow-y-auto">
                    <Playlists/>
                </div>

                <div className="overflow-y-auto">
                    <TrackList tracklist={tracklist}/>
                </div>

                <div className="overflow-y-auto">
                    <SongQueue/>
                </div>
            </div>

            <div className="mt-4">
                <AudioPlayer/>
            </div>
        </div>
    );
}

export default App;
