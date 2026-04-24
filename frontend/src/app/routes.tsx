import { createBrowserRouter } from "react-router";
import { AuthLayout } from "./components/layouts/AuthLayout";
import { MainLayout } from "./components/layouts/MainLayout";
import { LoginPage } from "./components/pages/LoginPage";
import { SignupPage } from "./components/pages/SignupPage";
import { DashboardPage } from "./components/pages/DashboardPage";
import { ProfilePage } from "./components/pages/ProfilePage";
import { JobHistoryPage } from "./components/pages/JobHistoryPage";
import { ProjectsPage } from "./components/pages/ProjectsPage";
import { SkillsPage } from "./components/pages/SkillsPage";
import { ResumesPage } from "./components/pages/ResumesPage";
import { ResumePreviewPage } from "./components/pages/ResumePreviewPage";
import { NotFoundPage } from "./components/pages/NotFoundPage";

export const router = createBrowserRouter([
  {
    path: "/",
    Component: AuthLayout,
    children: [
      { index: true, Component: LoginPage },
      { path: "signup", Component: SignupPage },
    ],
  },
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
