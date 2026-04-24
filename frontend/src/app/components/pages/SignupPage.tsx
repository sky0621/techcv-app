import { Navigate } from "react-router";

export function AppRedirectPage() {
  return <Navigate to="/app" replace />;
}
