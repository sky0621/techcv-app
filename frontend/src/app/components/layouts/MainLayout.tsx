import { Outlet, Link, useLocation } from "react-router";
import {
  LayoutDashboard,
  User,
  Briefcase,
  FolderKanban,
  Code,
  FileText,
  LogOut,
  Settings
} from "lucide-react";
import { Button } from "../ui/button";

const navItems = [
  { path: "/app", label: "ダッシュボード", icon: LayoutDashboard },
  { path: "/app/profile", label: "プロフィール", icon: User },
  { path: "/app/job-history", label: "職歴", icon: Briefcase },
  { path: "/app/projects", label: "案件", icon: FolderKanban },
  { path: "/app/skills", label: "スキル", icon: Code },
  { path: "/app/resumes", label: "経歴書", icon: FileText },
];

export function MainLayout() {
  const location = useLocation();

  return (
    <div className="flex h-screen bg-gray-50">
      {/* Sidebar */}
      <aside className="w-64 bg-white border-r border-gray-200 flex flex-col">
        <div className="p-6 border-b border-gray-200">
          <h1 className="font-bold text-xl text-gray-900">経歴書管理</h1>
        </div>

        <nav className="flex-1 p-4 space-y-1">
          {navItems.map((item) => {
            const Icon = item.icon;
            const isActive = location.pathname === item.path ||
              (item.path !== "/app" && location.pathname.startsWith(item.path));

            return (
              <Link
                key={item.path}
                to={item.path}
                className={`flex items-center gap-3 px-4 py-2.5 rounded-lg transition-colors ${
                  isActive
                    ? "bg-blue-50 text-blue-700"
                    : "text-gray-700 hover:bg-gray-100"
                }`}
              >
                <Icon className="w-5 h-5" />
                <span className="font-medium">{item.label}</span>
              </Link>
            );
          })}
        </nav>

        <div className="p-4 border-t border-gray-200 space-y-1">
          <Button variant="ghost" className="w-full justify-start gap-3">
            <Settings className="w-5 h-5" />
            <span>設定</span>
          </Button>
          <Button variant="ghost" className="w-full justify-start gap-3 text-red-600 hover:text-red-700 hover:bg-red-50">
            <LogOut className="w-5 h-5" />
            <span>ログアウト</span>
          </Button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 overflow-auto">
        <Outlet />
      </main>
    </div>
  );
}
