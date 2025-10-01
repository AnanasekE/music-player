/**
 * This file is the entry point for the React app, it sets up the root
 * element and renders the App component to the DOM.
 *
 * It is included in `src/index.html`.
 */

import {StrictMode} from "react";
import {createRoot} from "react-dom/client";
import {App} from "./App";
import {ThemeProvider} from "@/components/providers/themeProvider.tsx";
import {QueueProvider} from "@/components/providers/queueProvider.tsx";
import {TracksProvider} from "@/components/providers/tracksProvider.tsx";

const elem = document.getElementById("root")!;
const app = (
    <StrictMode>
        <ThemeProvider>
            <TracksProvider>
                <QueueProvider>
                    <App/>
                </QueueProvider>
            </TracksProvider>
        </ThemeProvider>
    </StrictMode>
);

if (import.meta.hot) {
    // With hot module reloading, `import.meta.hot.data` is persisted.
    const root = (import.meta.hot.data.root ??= createRoot(elem));
    root.render(app);
} else {
    // The hot module reloading API is not available in production.
    createRoot(elem).render(app);
}
