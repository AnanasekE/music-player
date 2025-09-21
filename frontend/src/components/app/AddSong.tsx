import {Card, CardContent, CardTitle} from "@/components/ui/card.tsx";
import {Button} from "@/components/ui/button.tsx";
import {Input} from "@/components/ui/input.tsx";
import {Label} from "@/components/ui/label.tsx";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select.tsx";
import {useEffect, useState} from "react";
import {set} from "react-hook-form";

interface AddSongProps {
    onAdd: () => void
}

function AddSong({onAdd}: AddSongProps) {
    const [possibleSongPaths, setPossibleSongPaths] = useState<string[]>([]);
    const [title, setTitle] = useState("");
    const [author, setAuthor] = useState("");
    const [cover, setCover] = useState<File | null>(null);
    const [filePath, setFilePath] = useState<string>(""); // from server
    const [newFile, setNewFile] = useState<File | null>(null); // user-uploaded file

    useEffect(() => {
        fetch("http://localhost:8080/audio-paths")
            .then((res) => res.json())
            .then((pathList: string[]) => setPossibleSongPaths(pathList))
            .catch((err) => console.error("Failed to fetch audio paths:", err));
    }, []);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();

        if (!filePath && !newFile) {
            console.error("No file selected!");
            return;
        }

        const formData = new FormData();
        formData.append("title", title);
        formData.append("author", author);

        if (cover) {
            formData.append("cover", cover);
        }

        let endpoint = "";

        if (newFile) {
            formData.append("file", newFile);
            endpoint = "/upload-track";
        } else if (filePath) {
            formData.append("filePath", filePath);
            endpoint = "/add-track";
        }

        fetch(`http://localhost:8080${endpoint}`, {
            method: "POST",
            body: formData,
        })
            .then((res) => {
                if (!res.ok) throw new Error("Failed to add/upload song");
                return res.text();
            })
            .then((msg) => console.log("Success:", msg))
            .catch((err) => console.error(err))
            .finally(() => {
                onAdd()
            });
    };

    return (
        <Card className="flex-1 overflow-y-auto h-full bg-secondary">
            <CardTitle className="text-4xl m-4">Add Song</CardTitle>
            <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4 text-left">
                    <div>
                        <Label htmlFor="title">Title</Label>
                        <Input id="title" value={title} onChange={(e) => setTitle(e.target.value)} required/>
                    </div>

                    <div>
                        <Label htmlFor="author">Author</Label>
                        <Input id="author" value={author} onChange={(e) => setAuthor(e.target.value)} required/>
                    </div>

                    <div>
                        <Label htmlFor="cover">Cover (optional)</Label>
                        <Input id="cover" type="file" accept="image/*"
                               onChange={(e) => setCover(e.target.files?.[0] || null)}/>
                    </div>

                    <div>
                        <Label>Pick from server paths</Label>
                        <Select onValueChange={setFilePath} disabled={!!newFile}>
                            <SelectTrigger>
                                <SelectValue placeholder="Choose existing file"/>
                            </SelectTrigger>
                            <SelectContent>
                                {possibleSongPaths && possibleSongPaths.map((path) => (
                                    <SelectItem key={path} value={path}>
                                        {path}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>

                    <div>
                        <Label htmlFor="newFile">Or upload new file</Label>
                        <Input
                            id="newFile"
                            type="file"
                            accept="audio/*"
                            onChange={(e) => setNewFile(e.target.files?.[0] || null)}
                            disabled={!!filePath}
                        />
                    </div>

                    <Button type="submit" className="w-full">
                        Add Song
                    </Button>
                </form>
            </CardContent>
        </Card>
    );
}

export default AddSong;
