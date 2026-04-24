import { useParams, Link } from "react-router";
import { ArrowLeft, Download, FileText, Github, Globe, Mail, MapPin, Phone } from "lucide-react";
import { Button } from "../ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";
import { Badge } from "../ui/badge";
import { Separator } from "../ui/separator";

const mockProfile = {
  name: "山田 太郎",
  email: "yamada@example.com",
  phone: "090-1234-5678",
  location: "東京都",
  github: "https://github.com/username",
  website: "https://example.com",
};

const mockProjects = [
  {
    id: 1,
    name: "ECサイトリニューアル",
    company: "株式会社A",
    period: "2024/01 - 現在",
    description: "大手ECサイトのフロントエンド刷新プロジェクト。モダンな技術スタックを採用し、ユーザー体験を大幅に改善。",
    role: "フロントエンドエンジニア",
    teamSize: "8名",
    technologies: ["React", "TypeScript", "Next.js", "Tailwind CSS", "GraphQL"],
    phases: ["要件定義", "設計", "実装", "テスト"],
    achievements: "ページ表示速度を50%改善、コンバージョン率15%向上。レスポンシブデザインの実装により、モバイルでのユーザー満足度が大幅に向上。",
  },
  {
    id: 2,
    name: "業務管理システム開発",
    company: "株式会社B",
    period: "2023/06 - 2023/12",
    description: "社内業務効率化のためのWebアプリケーション開発。複雑なワークフローを可視化し、承認プロセスを自動化。",
    role: "フルスタックエンジニア",
    teamSize: "5名",
    technologies: ["Vue.js", "Node.js", "PostgreSQL", "Docker", "AWS"],
    phases: ["設計", "実装", "テスト", "運用保守"],
    achievements: "業務時間を30%削減、ユーザー満足度90%以上。API設計とDB最適化により高速なレスポンスを実現。",
  },
];

const mockSkills = {
  "言語": ["TypeScript", "JavaScript", "Python", "Go"],
  "フレームワーク": ["React", "Vue.js", "Next.js", "Node.js"],
  "データベース": ["PostgreSQL", "MySQL", "MongoDB"],
  "インフラ": ["Docker", "AWS", "GCP"],
};

export function ResumePreviewPage() {
  const { id } = useParams();

  const handleDownloadMarkdown = () => {
    alert("Markdown形式でダウンロード（実装予定）");
  };

  const handleDownloadPDF = () => {
    alert("PDF形式でダウンロード（実装予定）");
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b sticky top-0 z-10">
        <div className="max-w-5xl mx-auto px-8 py-4 flex items-center justify-between">
          <Button asChild variant="ghost">
            <Link to="/app/resumes">
              <ArrowLeft className="w-4 h-4 mr-2" />
              経歴書一覧に戻る
            </Link>
          </Button>
          <div className="flex gap-2">
            <Button variant="outline" onClick={handleDownloadMarkdown}>
              <FileText className="w-4 h-4 mr-2" />
              Markdown
            </Button>
            <Button onClick={handleDownloadPDF}>
              <Download className="w-4 h-4 mr-2" />
              PDFダウンロード
            </Button>
          </div>
        </div>
      </div>

      {/* Preview Content */}
      <div className="max-w-5xl mx-auto p-8">
        <Card className="shadow-lg">
          <CardContent className="p-12">
            {/* Header Section */}
            <div className="space-y-6">
              <div>
                <h1 className="text-4xl font-bold text-gray-900 mb-2">{mockProfile.name}</h1>
                <p className="text-lg text-gray-600">Webエンジニア</p>
              </div>

              <div className="grid grid-cols-2 gap-3 text-sm">
                <div className="flex items-center gap-2 text-gray-700">
                  <Mail className="w-4 h-4" />
                  {mockProfile.email}
                </div>
                <div className="flex items-center gap-2 text-gray-700">
                  <Phone className="w-4 h-4" />
                  {mockProfile.phone}
                </div>
                <div className="flex items-center gap-2 text-gray-700">
                  <MapPin className="w-4 h-4" />
                  {mockProfile.location}
                </div>
                <div className="flex items-center gap-2 text-gray-700">
                  <Github className="w-4 h-4" />
                  <a href={mockProfile.github} className="text-blue-600 hover:underline">
                    GitHub
                  </a>
                </div>
              </div>
            </div>

            <Separator className="my-8" />

            {/* Summary */}
            <div className="space-y-3">
              <h2 className="text-2xl font-bold text-gray-900">概要</h2>
              <p className="text-gray-700 leading-relaxed">
                Webエンジニアとして5年以上の経験があり、フロントエンドからバックエンドまで幅広く対応できます。
                React、TypeScriptを用いたモダンなフロントエンド開発が得意で、ユーザー体験を重視した実装を心がけています。
                チーム開発での経験も豊富で、アジャイル開発手法を活用したプロジェクト推進が可能です。
              </p>
            </div>

            <Separator className="my-8" />

            {/* Skills */}
            <div className="space-y-4">
              <h2 className="text-2xl font-bold text-gray-900">スキル</h2>
              {Object.entries(mockSkills).map(([category, skills]) => (
                <div key={category} className="space-y-2">
                  <h3 className="font-semibold text-gray-800">{category}</h3>
                  <div className="flex flex-wrap gap-2">
                    {skills.map(skill => (
                      <Badge key={skill} variant="secondary">{skill}</Badge>
                    ))}
                  </div>
                </div>
              ))}
            </div>

            <Separator className="my-8" />

            {/* Projects */}
            <div className="space-y-6">
              <h2 className="text-2xl font-bold text-gray-900">プロジェクト実績</h2>
              {mockProjects.map((project, index) => (
                <div key={project.id} className="space-y-3">
                  <div>
                    <h3 className="text-xl font-bold text-gray-900">{project.name}</h3>
                    <div className="text-sm text-gray-600 mt-1">
                      {project.company} • {project.period}
                    </div>
                  </div>

                  <p className="text-gray-700">{project.description}</p>

                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <span className="font-semibold text-gray-800">役割: </span>
                      <span className="text-gray-700">{project.role}</span>
                    </div>
                    <div>
                      <span className="font-semibold text-gray-800">チーム規模: </span>
                      <span className="text-gray-700">{project.teamSize}</span>
                    </div>
                  </div>

                  <div>
                    <div className="font-semibold text-gray-800 mb-2">使用技術</div>
                    <div className="flex flex-wrap gap-1.5">
                      {project.technologies.map(tech => (
                        <Badge key={tech} variant="outline">{tech}</Badge>
                      ))}
                    </div>
                  </div>

                  <div>
                    <div className="font-semibold text-gray-800 mb-2">担当工程</div>
                    <div className="flex flex-wrap gap-1.5">
                      {project.phases.map(phase => (
                        <Badge key={phase} variant="outline">{phase}</Badge>
                      ))}
                    </div>
                  </div>

                  <div>
                    <div className="font-semibold text-gray-800 mb-2">成果・工夫</div>
                    <p className="text-gray-700">{project.achievements}</p>
                  </div>

                  {index < mockProjects.length - 1 && <Separator className="my-6" />}
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
