import { Outlet } from "react-router";

export function Layout() {
  return (
    <div className="h-screen max-h-screen bg-zinc-800 text-white">
      <Outlet />
    </div>
  );
}
