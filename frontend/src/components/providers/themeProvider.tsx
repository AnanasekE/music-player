"use client"
import { createContext, useContext, useEffect, useState, type ReactNode } from "react";

type Theme = "light" | "dark" | "violet";

const availableThemes:Theme[] = ["light", "dark", "violet"]

const ThemeContext = createContext<{
    theme: Theme;
    setTheme: (t: Theme) => void;
}>({
    theme: "light",
    setTheme: () => {},
});

export function ThemeProvider({ children }: { children: ReactNode }) {
    const [theme, setTheme] = useState<Theme>(() => {
        const saved = localStorage.getItem("theme");
        if (saved && availableThemes.includes(saved as Theme)) {
            return saved as Theme;
        }
        return "violet";
    });


    useEffect(() => {
        const root = document.documentElement;

        if (theme === "light") {
            root.classList.remove("dark");
            root.removeAttribute("data-theme");
        } else if (theme === "dark") {
            root.classList.add("dark");
            root.removeAttribute("data-theme");
        } else {
            root.classList.remove("dark");
            root.dataset.theme = theme;
        }
        
        localStorage.setItem("theme", theme)
    }, [theme]);

    return (
        <ThemeContext.Provider value={{ theme, setTheme }}>
            {children}
        </ThemeContext.Provider>
    );
}

export function useTheme() {
    return useContext(ThemeContext);
}

