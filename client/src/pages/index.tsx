import { Outlet } from "react-router";

export function Layout() {
  return (
    <div className="bg-zinc-800 text-white h-screen max-h-screen">
      <Outlet />
    </div>
  )
}
