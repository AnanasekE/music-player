import { useTheme } from "../providers/themeProvider.tsx";
import { Button } from "@/components/ui/button";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card.tsx";

const themes = ["light", "dark", "violet"] as const;

export function ThemeSwitcher() {
    const { theme, setTheme } = useTheme();

    return (
        <Card>
            <CardHeader>Change Theme</CardHeader>
            <CardContent className={"flex flex-col"}>
                {themes.map((t) => (
                    <Button
                        className={"m-1"}
                        key={t}
                        variant={theme === t ? "default" : "outline"}
                        size="sm"
                        onClick={() => setTheme(t)}
                    >
                        {t}
                    </Button>
                ))}
            </CardContent>
        </Card>
    );
}
