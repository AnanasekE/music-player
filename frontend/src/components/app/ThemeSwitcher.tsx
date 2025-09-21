import { useTheme } from "../providers/themeProvider.tsx";
import { Button } from "@/components/ui/button";

const themes = ["light", "dark", "violet"] as const;

export function ThemeSwitcher() {
    const { theme, setTheme } = useTheme();

    return (
        <div className="fixed top-4 right-4 flex gap-2 z-50">
            {themes.map((t) => (
                <Button
                    key={t}
                    variant={theme === t ? "default" : "outline"}
                    size="sm"
                    onClick={() => setTheme(t)}
                >
                    {t}
                </Button>
            ))}
        </div>
    );
}
