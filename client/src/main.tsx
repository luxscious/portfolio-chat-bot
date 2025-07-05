import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./styles/index.css";
import App from "./App";

// Dynamically set favicon from env variable
const favicon = document.querySelector(
  "link[rel='icon']"
) as HTMLLinkElement | null;
if (favicon) {
  favicon.href = `${import.meta.env.VITE_BLOB_BASE_URL}/g_icon.png`;
}

createRoot(document.getElementById("root") as HTMLElement).render(
  <StrictMode>
    <App />
  </StrictMode>
);
