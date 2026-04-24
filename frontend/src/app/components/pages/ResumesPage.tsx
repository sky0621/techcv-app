import { useState } from "react";
import { Plus, FileText, Eye, Pencil, Copy, Trash2, Download } from "lucide-react";
import { Link } from "react-router";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Textarea } from "../ui/textarea";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";
import { Checkbox } from "../ui/checkbox";
import { Badge } from "../ui/badge";

type Resume = {
  id: number;
  name: string;
  description: string;
  selectedProjects: number[];
  summary: string;
  createdAt: string;
  updatedAt: string;
};

const mockProjects = [
  { id: 1, name: "ECサイトリニューアル", company: "株式会社A", period: "2024/01 - 現在" },
  { id: 2, name: "業務管理システム開発", company: "株式会社B", period: "2023/06 - 2023/12" },
  { id: 3, name: "モバイルアプリ開発", company: "株式会社C", period: "2023/01 - 2023/05" },
  { id: 4, name: "データ分析基盤構築", company: "株式会社D", period: "2022/08 - 2022/12" },
  { id: 5, name: "社内システム刷新", company: "株式会社E", period: "2022/01 - 2022/07" },
];

const initialResumes: Resume[] = [
  {
    id: 1,
    name: "ベース経歴書",
    description: "全案件を含む標準の経歴書",
    selectedProjects: [1, 2, 3, 4, 5],
    summary: "Webエンジニアとして5年以上の経験があり、フロントエンドからバックエンドまで幅広く対応できます。",
    createdAt: "2024-04-01",
    updatedAt: "2024-04-20",
  },
  {
    id: 2,
    name: "フロントエンド特化版",
    description: "フロントエンド案件を中心に構成",
    selectedProjects: [1, 3],
    summary: "React、TypeScriptを用いたモダンなフロントエンド開発が得意です。ユーザー体験を重視した実装を心がけています。",
    createdAt: "2024-04-10",
    updatedAt: "2024-04-18",
  },
];

export function ResumesPage() {
  const [resumes, setResumes] = useState<Resume[]>(initialResumes);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingResume, setEditingResume] = useState<Resume | null>(null);
  const [formData, setFormData] = useState({
    name: "",
    description: "",
    selectedProjects: [] as number[],
    summary: "",
  });

  const handleAdd = () => {
    setEditingResume(null);
    setFormData({
      name: "",
      description: "",
      selectedProjects: [],
      summary: "",
    });
    setIsDialogOpen(true);
  };

  const handleEdit = (resume: Resume) => {
    setEditingResume(resume);
    setFormData({
      name: resume.name,
      description: resume.description,
      selectedProjects: resume.selectedProjects,
      summary: resume.summary,
    });
    setIsDialogOpen(true);
  };

  const handleCopy = (resume: Resume) => {
    const newResume: Resume = {
      ...resume,
      id: Date.now(),
      name: `${resume.name}のコピー`,
      createdAt: new Date().toISOString().split('T')[0],
      updatedAt: new Date().toISOString().split('T')[0],
    };
    setResumes([newResume, ...resumes]);
  };

  const handleSave = () => {
    const now = new Date().toISOString().split('T')[0];

    if (editingResume) {
      setResumes(resumes.map(r =>
        r.id === editingResume.id
          ? { ...r, ...formData, updatedAt: now }
          : r
      ));
    } else {
      const newResume: Resume = {
        id: Date.now(),
        ...formData,
        createdAt: now,
        updatedAt: now,
      };
      setResumes([newResume, ...resumes]);
    }
    setIsDialogOpen(false);
  };

  const handleDelete = (id: number) => {
    if (confirm("本当に削除しますか？")) {
      setResumes(resumes.filter(r => r.id !== id));
    }
  };

  const toggleProject = (projectId: number) => {
    setFormData({
      ...formData,
      selectedProjects: formData.selectedProjects.includes(projectId)
        ? formData.selectedProjects.filter(id => id !== projectId)
        : [...formData.selectedProjects, projectId]
    });
  };

  return (
    <div className="p-8">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">経歴書バリエーション</h1>
            <p className="text-gray-600 mt-1">応募先に合わせた経歴書を作成</p>
          </div>
          <Button onClick={handleAdd}>
            <Plus className="w-4 h-4 mr-2" />
            新規作成
          </Button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {resumes.map((resume) => (
            <Card key={resume.id} className="flex flex-col">
              <CardHeader>
                <div className="flex items-start justify-between">
                  <FileText className="w-8 h-8 text-blue-600" />
                  <Badge variant="outline">{resume.selectedProjects.length}件</Badge>
                </div>
                <CardTitle className="mt-4">{resume.name}</CardTitle>
                <CardDescription>{resume.description}</CardDescription>
                <div className="text-xs text-gray-500 mt-2">
                  更新: {resume.updatedAt}
                </div>
              </CardHeader>
              <CardContent className="flex-1 flex flex-col gap-2">
                <div className="flex-1">
                  <p className="text-sm text-gray-600 line-clamp-3">{resume.summary}</p>
                </div>
                <div className="grid grid-cols-2 gap-2 mt-auto">
                  <Button asChild variant="outline" size="sm">
                    <Link to={`/app/resumes/${resume.id}/preview`}>
                      <Eye className="w-4 h-4 mr-1" />
                      プレビュー
                    </Link>
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleEdit(resume)}
                  >
                    <Pencil className="w-4 h-4 mr-1" />
                    編集
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleCopy(resume)}
                  >
                    <Copy className="w-4 h-4 mr-1" />
                    複製
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleDelete(resume.id)}
                  >
                    <Trash2 className="w-4 h-4 mr-1 text-red-600" />
                    削除
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}

          {resumes.length === 0 && (
            <Card className="col-span-full">
              <CardContent className="flex flex-col items-center justify-center py-12">
                <FileText className="w-12 h-12 text-gray-400 mb-4" />
                <p className="text-gray-600 mb-4">経歴書がまだ作成されていません</p>
                <Button onClick={handleAdd}>
                  <Plus className="w-4 h-4 mr-2" />
                  最初の経歴書を作成
                </Button>
              </CardContent>
            </Card>
          )}
        </div>

        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
            <DialogHeader>
              <DialogTitle>
                {editingResume ? "経歴書を編集" : "経歴書を作成"}
              </DialogTitle>
              <DialogDescription>
                表示する案件を選択し、自己PRを記入してください
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-6 py-4">
              <div className="space-y-2">
                <Label htmlFor="name">経歴書名</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="フロントエンド特化版"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="description">説明</Label>
                <Input
                  id="description"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder="この経歴書の用途や特徴"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="summary">自己PR・要約</Label>
                <Textarea
                  id="summary"
                  rows={4}
                  value={formData.summary}
                  onChange={(e) => setFormData({ ...formData, summary: e.target.value })}
                  placeholder="あなたの強みや経験を簡潔に記述してください"
                />
              </div>

              <div className="space-y-3">
                <Label>表示する案件を選択</Label>
                <div className="border rounded-lg p-4 space-y-2 max-h-64 overflow-y-auto">
                  {mockProjects.map(project => (
                    <div
                      key={project.id}
                      className="flex items-start space-x-3 p-3 rounded-lg hover:bg-gray-50 transition-colors"
                    >
                      <Checkbox
                        id={`project-${project.id}`}
                        checked={formData.selectedProjects.includes(project.id)}
                        onCheckedChange={() => toggleProject(project.id)}
                      />
                      <div className="flex-1">
                        <Label
                          htmlFor={`project-${project.id}`}
                          className="font-medium cursor-pointer"
                        >
                          {project.name}
                        </Label>
                        <div className="text-sm text-gray-600">
                          {project.company} • {project.period}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
                <p className="text-sm text-gray-600">
                  {formData.selectedProjects.length}件の案件を選択中
                </p>
              </div>
            </div>

            <DialogFooter>
              <Button variant="outline" onClick={() => setIsDialogOpen(false)}>
                キャンセル
              </Button>
              <Button onClick={handleSave}>保存</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
}
