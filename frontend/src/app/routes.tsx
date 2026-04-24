import { createBrowserRouter, Navigate } from "react-router";
import { MainLayout } from "./components/layouts/MainLayout";
import { DashboardPage } from "./components/pages/DashboardPage";
import { ProfilePage } from "./components/pages/ProfilePage";
import { JobHistoryPage } from "./components/pages/JobHistoryPage";
import { ProjectsPage } from "./components/pages/ProjectsPage";
import { SkillsPage } from "./components/pages/SkillsPage";
import { ResumesPage } from "./components/pages/ResumesPage";
import { ResumePreviewPage } from "./components/pages/ResumePreviewPage";
import { NotFoundPage } from "./components/pages/NotFoundPage";

export const router = createBrowserRouter([
  { path: "/", element: <Navigate to="/app" replace /> },
  {
    path: "/app",
    Component: MainLayout,
    children: [
      { index: true, Component: DashboardPage },
      { path: "profile", Component: ProfilePage },
      { path: "job-history", Component: JobHistoryPage },
      { path: "projects", Component: ProjectsPage },
      { path: "skills", Component: SkillsPage },
      { path: "resumes", Component: ResumesPage },
      { path: "resumes/:id/preview", Component: ResumePreviewPage },
    ],
  },
  {
    path: "*",
    Component: NotFoundPage,
  },
]);
