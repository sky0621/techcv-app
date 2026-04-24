import { Link } from "react-router";
import { Plus, FileText, FolderKanban, Clock } from "lucide-react";
import { Button } from "../ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";

const recentProjects = [
  { id: 1, name: "ECサイトリニューアル", company: "株式会社A", period: "2024/01 - 現在" },
  { id: 2, name: "業務管理システム開発", company: "株式会社B", period: "2023/06 - 2023/12" },
  { id: 3, name: "モバイルアプリ開発", company: "株式会社C", period: "2023/01 - 2023/05" },
];

const resumes = [
  { id: 1, name: "ベース経歴書", updatedAt: "2024/04/20", projectCount: 12 },
  { id: 2, name: "フロントエンド特化版", updatedAt: "2024/04/18", projectCount: 8 },
  { id: 3, name: "バックエンド特化版", updatedAt: "2024/04/15", projectCount: 7 },
];

export function DashboardPage() {
  return (
    <div className="p-8">
      <div className="max-w-7xl mx-auto space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">ダッシュボード</h1>
          <p className="text-gray-600 mt-1">職務経歴の管理と経歴書作成</p>
        </div>

        {/* クイックアクション */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Link to="/app/projects">
            <Card className="hover:shadow-lg transition-shadow cursor-pointer">
              <CardHeader>
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-blue-100 rounded-lg">
                    <Plus className="w-6 h-6 text-blue-600" />
                  </div>
                  <CardTitle className="text-lg">新規案件を追加</CardTitle>
                </div>
              </CardHeader>
            </Card>
          </Link>

          <Link to="/app/resumes">
            <Card className="hover:shadow-lg transition-shadow cursor-pointer">
              <CardHeader>
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-green-100 rounded-lg">
                    <FileText className="w-6 h-6 text-green-600" />
                  </div>
                  <CardTitle className="text-lg">経歴書を作成</CardTitle>
                </div>
              </CardHeader>
            </Card>
          </Link>

          <Link to="/app/profile">
            <Card className="hover:shadow-lg transition-shadow cursor-pointer">
              <CardHeader>
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-purple-100 rounded-lg">
                    <FolderKanban className="w-6 h-6 text-purple-600" />
                  </div>
                  <CardTitle className="text-lg">プロフィール編集</CardTitle>
                </div>
              </CardHeader>
            </Card>
          </Link>
        </div>

        {/* 最近編集した案件 */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>最近編集した案件</CardTitle>
                <CardDescription>直近で更新された案件一覧</CardDescription>
              </div>
              <Button asChild variant="outline">
                <Link to="/app/projects">すべて見る</Link>
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {recentProjects.map((project) => (
                <div
                  key={project.id}
                  className="flex items-center justify-between p-4 border rounded-lg hover:bg-gray-50 transition-colors"
                >
                  <div className="flex items-center gap-3">
                    <FolderKanban className="w-5 h-5 text-gray-400" />
                    <div>
                      <div className="font-medium text-gray-900">{project.name}</div>
                      <div className="text-sm text-gray-600">{project.company}</div>
                    </div>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-gray-500">
                    <Clock className="w-4 h-4" />
                    {project.period}
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* 経歴書バリエーション */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>経歴書バリエーション</CardTitle>
                <CardDescription>作成済みの経歴書一覧</CardDescription>
              </div>
              <Button asChild>
                <Link to="/app/resumes">
                  <Plus className="w-4 h-4 mr-2" />
                  新規作成
                </Link>
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {resumes.map((resume) => (
                <Card key={resume.id} className="hover:shadow-md transition-shadow">
                  <CardHeader>
                    <CardTitle className="text-base">{resume.name}</CardTitle>
                    <CardDescription>
                      {resume.projectCount}件の案件 • {resume.updatedAt}更新
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="flex gap-2">
                    <Button asChild variant="outline" size="sm" className="flex-1">
                      <Link to={`/app/resumes/${resume.id}/preview`}>プレビュー</Link>
                    </Button>
                    <Button asChild variant="outline" size="sm" className="flex-1">
                      <Link to="/app/resumes">編集</Link>
                    </Button>
                  </CardContent>
                </Card>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
