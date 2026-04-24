import { Outlet } from "react-router";

export function CenteredLayout() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="w-full max-w-md">
        <Outlet />
      </div>
    </div>
  );
}
