import { createRoot } from "react-dom/client";
import { HashRouter, Route, Routes } from "react-router";
import { Layout } from "./pages";
import { HomePage } from "./pages/home";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import "./index.css";

const queryClient = new QueryClient();

createRoot(document.getElementById("root")!).render(
  <>
    <QueryClientProvider client={queryClient}>
      <HashRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<HomePage />} />
          </Route>
        </Routes>
      </HashRouter>
    </QueryClientProvider>
  </>,
);
