import { useState } from "react";
import { Plus, Search, Pencil, Trash2, FolderKanban, Save } from "lucide-react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Textarea } from "../ui/textarea";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";
import { Badge } from "../ui/badge";
import { Checkbox } from "../ui/checkbox";

type Project = {
  id: number;
  name: string;
  company: string;
  startDate: string;
  endDate: string;
  description: string;
  role: string;
  teamSize: string;
  technologies: string[];
  phases: string[];
  achievements: string;
  isDraft: boolean;
};

const initialProjects: Project[] = [
  {
    id: 1,
    name: "ECサイトリニューアル",
    company: "株式会社A",
    startDate: "2024-01",
    endDate: "現在",
    description: "大手ECサイトのフロントエンド刷新プロジェクト",
    role: "フロントエンドエンジニア",
    teamSize: "8名",
    technologies: ["React", "TypeScript", "Next.js", "Tailwind CSS"],
    phases: ["要件定義", "設計", "実装", "テスト"],
    achievements: "ページ表示速度を50%改善、コンバージョン率15%向上",
    isDraft: false,
  },
  {
    id: 2,
    name: "業務管理システム開発",
    company: "株式会社B",
    startDate: "2023-06",
    endDate: "2023-12",
    description: "社内業務効率化のためのWebアプリケーション開発",
    role: "フルスタックエンジニア",
    teamSize: "5名",
    technologies: ["Vue.js", "Node.js", "PostgreSQL", "Docker"],
    phases: ["設計", "実装", "テスト", "運用保守"],
    achievements: "業務時間を30%削減、ユーザー満足度90%以上",
    isDraft: false,
  },
];

const allTechnologies = [
  "React", "Vue.js", "Angular", "TypeScript", "JavaScript",
  "Node.js", "Python", "Go", "Ruby", "PHP",
  "PostgreSQL", "MySQL", "MongoDB", "Redis",
  "Docker", "Kubernetes", "AWS", "GCP", "Azure",
  "Next.js", "Nuxt.js", "Tailwind CSS", "GraphQL", "REST API"
];

const phases = ["要件定義", "設計", "実装", "テスト", "運用保守", "リリース"];

export function ProjectsPage() {
  const [projects, setProjects] = useState<Project[]>(initialProjects);
  const [searchQuery, setSearchQuery] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingProject, setEditingProject] = useState<Project | null>(null);
  const [formData, setFormData] = useState<Partial<Project>>({
    name: "",
    company: "",
    startDate: "",
    endDate: "",
    description: "",
    role: "",
    teamSize: "",
    technologies: [],
    phases: [],
    achievements: "",
    isDraft: false,
  });

  const filteredProjects = projects.filter(p =>
    p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    p.company.toLowerCase().includes(searchQuery.toLowerCase()) ||
    p.technologies.some(t => t.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  const handleAdd = () => {
    setEditingProject(null);
    setFormData({
      name: "",
      company: "",
      startDate: "",
      endDate: "",
      description: "",
      role: "",
      teamSize: "",
      technologies: [],
      phases: [],
      achievements: "",
      isDraft: false,
    });
    setIsDialogOpen(true);
  };

  const handleEdit = (project: Project) => {
    setEditingProject(project);
    setFormData(project);
    setIsDialogOpen(true);
  };

  const handleSave = (asDraft: boolean = false) => {
    const projectData = { ...formData, isDraft: asDraft } as Project;

    if (editingProject) {
      setProjects(projects.map(p => p.id === editingProject.id ? { ...p, ...projectData } : p));
    } else {
      setProjects([{ id: Date.now(), ...projectData }, ...projects]);
    }
    setIsDialogOpen(false);
  };

  const handleDelete = (id: number) => {
    if (confirm("本当に削除しますか？")) {
      setProjects(projects.filter(p => p.id !== id));
    }
  };

  const toggleTechnology = (tech: string) => {
    const current = formData.technologies || [];
    setFormData({
      ...formData,
      technologies: current.includes(tech)
        ? current.filter(t => t !== tech)
        : [...current, tech]
    });
  };

  const togglePhase = (phase: string) => {
    const current = formData.phases || [];
    setFormData({
      ...formData,
      phases: current.includes(phase)
        ? current.filter(p => p !== phase)
        : [...current, phase]
    });
  };

  return (
    <div className="p-8">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">案件管理</h1>
            <p className="text-gray-600 mt-1">プロジェクト実績を詳細に記録</p>
          </div>
          <Button onClick={handleAdd}>
            <Plus className="w-4 h-4 mr-2" />
            案件を追加
          </Button>
        </div>

        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
          <Input
            placeholder="案件名、会社名、技術で検索..."
            className="pl-10"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>

        <div className="grid grid-cols-1 gap-4">
          {filteredProjects.map((project) => (
            <Card key={project.id} className={project.isDraft ? "border-yellow-300 bg-yellow-50" : ""}>
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div className="flex items-start gap-4 flex-1">
                    <div className="p-3 bg-purple-100 rounded-lg">
                      <FolderKanban className="w-6 h-6 text-purple-600" />
                    </div>
                    <div className="flex-1">
                      <div className="flex items-center gap-2">
                        <CardTitle className="text-xl">{project.name}</CardTitle>
                        {project.isDraft && <Badge variant="outline" className="bg-yellow-100">下書き</Badge>}
                      </div>
                      <CardDescription className="mt-1">
                        {project.company} • {project.startDate} 〜 {project.endDate}
                      </CardDescription>
                      <p className="text-sm text-gray-700 mt-2">{project.description}</p>
                      <div className="mt-3 space-y-2">
                        <div className="flex items-center gap-2 text-sm">
                          <span className="font-medium text-gray-700">役割:</span>
                          <span className="text-gray-600">{project.role}</span>
                          <span className="text-gray-400">•</span>
                          <span className="text-gray-600">チーム{project.teamSize}</span>
                        </div>
                        <div className="flex flex-wrap gap-1.5">
                          {project.technologies.map(tech => (
                            <Badge key={tech} variant="secondary">{tech}</Badge>
                          ))}
                        </div>
                      </div>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleEdit(project)}
                    >
                      <Pencil className="w-4 h-4" />
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleDelete(project.id)}
                    >
                      <Trash2 className="w-4 h-4 text-red-600" />
                    </Button>
                  </div>
                </div>
              </CardHeader>
            </Card>
          ))}

          {filteredProjects.length === 0 && (
            <Card>
              <CardContent className="flex flex-col items-center justify-center py-12">
                <FolderKanban className="w-12 h-12 text-gray-400 mb-4" />
                <p className="text-gray-600 mb-4">
                  {searchQuery ? "検索結果が見つかりませんでした" : "案件がまだ登録されていません"}
                </p>
                {!searchQuery && (
                  <Button onClick={handleAdd}>
                    <Plus className="w-4 h-4 mr-2" />
                    最初の案件を追加
                  </Button>
                )}
              </CardContent>
            </Card>
          )}
        </div>

        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
            <DialogHeader>
              <DialogTitle>
                {editingProject ? "案件を編集" : "案件を追加"}
              </DialogTitle>
              <DialogDescription>
                プロジェクトの詳細情報を入力してください
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="name">案件名</Label>
                  <Input
                    id="name"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    placeholder="ECサイトリニューアル"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="company">会社名</Label>
                  <Input
                    id="company"
                    value={formData.company}
                    onChange={(e) => setFormData({ ...formData, company: e.target.value })}
                    placeholder="株式会社〇〇"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="startDate">開始日</Label>
                  <Input
                    id="startDate"
                    type="month"
                    value={formData.startDate}
                    onChange={(e) => setFormData({ ...formData, startDate: e.target.value })}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="endDate">終了日</Label>
                  <Input
                    id="endDate"
                    value={formData.endDate}
                    onChange={(e) => setFormData({ ...formData, endDate: e.target.value })}
                    placeholder="現在 または 2024-12"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="description">プロダクト概要</Label>
                <Textarea
                  id="description"
                  rows={3}
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder="プロジェクトの概要を記述"
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="role">担当役割</Label>
                  <Input
                    id="role"
                    value={formData.role}
                    onChange={(e) => setFormData({ ...formData, role: e.target.value })}
                    placeholder="フロントエンドエンジニア"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="teamSize">チーム規模</Label>
                  <Input
                    id="teamSize"
                    value={formData.teamSize}
                    onChange={(e) => setFormData({ ...formData, teamSize: e.target.value })}
                    placeholder="8名"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label>使用技術</Label>
                <div className="border rounded-lg p-4 max-h-48 overflow-y-auto">
                  <div className="flex flex-wrap gap-2">
                    {allTechnologies.map(tech => (
                      <Badge
                        key={tech}
                        variant={formData.technologies?.includes(tech) ? "default" : "outline"}
                        className="cursor-pointer"
                        onClick={() => toggleTechnology(tech)}
                      >
                        {tech}
                      </Badge>
                    ))}
                  </div>
                </div>
              </div>

              <div className="space-y-2">
                <Label>担当工程</Label>
                <div className="grid grid-cols-3 gap-3">
                  {phases.map(phase => (
                    <div key={phase} className="flex items-center space-x-2">
                      <Checkbox
                        id={phase}
                        checked={formData.phases?.includes(phase)}
                        onCheckedChange={() => togglePhase(phase)}
                      />
                      <Label htmlFor={phase} className="cursor-pointer">{phase}</Label>
                    </div>
                  ))}
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="achievements">成果・工夫・改善内容</Label>
                <Textarea
                  id="achievements"
                  rows={4}
                  value={formData.achievements}
                  onChange={(e) => setFormData({ ...formData, achievements: e.target.value })}
                  placeholder="具体的な成果や工夫した点を記述"
                />
              </div>
            </div>

            <DialogFooter className="gap-2">
              <Button variant="outline" onClick={() => setIsDialogOpen(false)}>
                キャンセル
              </Button>
              <Button variant="outline" onClick={() => handleSave(true)}>
                <Save className="w-4 h-4 mr-2" />
                下書き保存
              </Button>
              <Button onClick={() => handleSave(false)}>保存</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
}
