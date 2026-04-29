import { useEffect, useState } from "react";
import { Plus, Code, Database, Cloud, Wrench, Pencil, Trash2 } from "lucide-react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select";
import { Badge } from "../ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";

type Skill = {
  id: number;
  name: string;
  category: string;
  experience: string;
  proficiency: string;
};

type SkillOption = {
  id: string;
  name: string;
  icon?: string;
  sortOrder: number;
};

type SkillOptionsResponse = {
  categories: SkillOption[];
  proficiencyLevels: SkillOption[];
};

const iconComponents: Record<string, any> = {
  cloud: Cloud,
  code: Code,
  database: Database,
  wrench: Wrench,
};

const initialSkills: Skill[] = [
  { id: 1, name: "TypeScript", category: "言語", experience: "3年", proficiency: "上級" },
  { id: 2, name: "React", category: "フレームワーク", experience: "3年", proficiency: "上級" },
  { id: 3, name: "Node.js", category: "フレームワーク", experience: "2年", proficiency: "中級" },
  { id: 4, name: "PostgreSQL", category: "データベース", experience: "2年", proficiency: "中級" },
  { id: 5, name: "Docker", category: "インフラ", experience: "2年", proficiency: "中級" },
  { id: 6, name: "AWS", category: "インフラ", experience: "1年", proficiency: "初級" },
  { id: 7, name: "Git", category: "ツール", experience: "4年", proficiency: "上級" },
];

const defaultCategory = "言語";
const defaultProficiency = "中級";

export function SkillsPage() {
  const [skills, setSkills] = useState<Skill[]>(initialSkills);
  const [categories, setCategories] = useState<SkillOption[]>([]);
  const [proficiencyLevels, setProficiencyLevels] = useState<SkillOption[]>([]);
  const [isOptionsLoading, setIsOptionsLoading] = useState(true);
  const [optionsError, setOptionsError] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingSkill, setEditingSkill] = useState<Skill | null>(null);
  const [formData, setFormData] = useState({
    name: "",
    category: defaultCategory,
    experience: "",
    proficiency: defaultProficiency,
  });

  useEffect(() => {
    const controller = new AbortController();

    async function loadOptions() {
      setIsOptionsLoading(true);
      setOptionsError("");

      try {
        const response = await fetch("/api/skills/options", {
          signal: controller.signal,
        });
        if (!response.ok) {
          throw new Error("スキルの選択肢の取得に失敗しました");
        }

        const data = (await response.json()) as SkillOptionsResponse;
        setCategories(data.categories ?? []);
        setProficiencyLevels(data.proficiencyLevels ?? []);
      } catch (caught) {
        if (caught instanceof DOMException && caught.name === "AbortError") {
          return;
        }
        setOptionsError(
          caught instanceof Error ? caught.message : "スキルの選択肢の取得に失敗しました",
        );
      } finally {
        if (!controller.signal.aborted) {
          setIsOptionsLoading(false);
        }
      }
    }

    void loadOptions();

    return () => {
      controller.abort();
    };
  }, []);

  const categoryNames = categories.map((category) => category.name);
  const proficiencyLevelNames = proficiencyLevels.map((level) => level.name);
  const selectedDefaultCategory = categoryNames[0] ?? defaultCategory;
  const selectedDefaultProficiency =
    proficiencyLevelNames.find((level) => level === defaultProficiency) ??
    proficiencyLevelNames[0] ??
    defaultProficiency;

  const skillsByCategory = categoryNames.reduce((acc, category) => {
    acc[category] = skills.filter(s => s.category === category);
    return acc;
  }, {} as Record<string, Skill[]>);

  const handleAdd = () => {
    setEditingSkill(null);
    setFormData({
      name: "",
      category: selectedDefaultCategory,
      experience: "",
      proficiency: selectedDefaultProficiency,
    });
    setIsDialogOpen(true);
  };

  const handleEdit = (skill: Skill) => {
    setEditingSkill(skill);
    setFormData(skill);
    setIsDialogOpen(true);
  };

  const handleSave = () => {
    if (editingSkill) {
      setSkills(skills.map(s => s.id === editingSkill.id ? { ...s, ...formData } : s));
    } else {
      setSkills([...skills, { id: Date.now(), ...formData }]);
    }
    setIsDialogOpen(false);
  };

  const handleDelete = (id: number) => {
    if (confirm("本当に削除しますか？")) {
      setSkills(skills.filter(s => s.id !== id));
    }
  };

  const getProficiencyColor = (proficiency: string) => {
    switch (proficiency) {
      case "初級": return "bg-blue-100 text-blue-700";
      case "中級": return "bg-green-100 text-green-700";
      case "上級": return "bg-purple-100 text-purple-700";
      case "エキスパート": return "bg-orange-100 text-orange-700";
      default: return "bg-gray-100 text-gray-700";
    }
  };

  const getCategoryIcon = (icon?: string) => {
    if (!icon) {
      return Wrench;
    }

    return iconComponents[icon] ?? Wrench;
  };

  return (
    <div className="p-8">
      <div className="max-w-6xl mx-auto space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">スキル管理</h1>
            <p className="text-gray-600 mt-1">技術スキルと習熟度を管理</p>
          </div>
          <Button onClick={handleAdd} disabled={isOptionsLoading || !!optionsError}>
            <Plus className="w-4 h-4 mr-2" />
            スキルを追加
          </Button>
        </div>

        {isOptionsLoading && (
          <div className="rounded-md border border-gray-200 bg-gray-50 px-4 py-3 text-sm text-gray-600">
            スキルの選択肢を読み込み中です。
          </div>
        )}
        {optionsError && (
          <div className="rounded-md border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {optionsError}
          </div>
        )}

        <Tabs defaultValue="all" className="w-full">
          <TabsList>
            <TabsTrigger value="all">すべて ({skills.length})</TabsTrigger>
            {categoryNames.map(cat => (
              <TabsTrigger key={cat} value={cat}>
                {cat} ({skillsByCategory[cat]?.length ?? 0})
              </TabsTrigger>
            ))}
          </TabsList>

          <TabsContent value="all" className="space-y-4 mt-6">
            {categories.map(category => {
              const Icon = getCategoryIcon(category.icon);
              const categoryName = category.name;
              const categorySkills = skillsByCategory[categoryName];

              if (categorySkills.length === 0) return null;

              return (
                <Card key={category.id}>
                  <CardHeader>
                    <div className="flex items-center gap-3">
                      <Icon className="w-5 h-5 text-gray-600" />
                      <CardTitle>{categoryName}</CardTitle>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                      {categorySkills.map(skill => (
                        <div
                          key={skill.id}
                          className="flex items-center justify-between p-3 border rounded-lg hover:bg-gray-50 transition-colors"
                        >
                          <div className="flex-1">
                            <div className="font-medium text-gray-900">{skill.name}</div>
                            <div className="flex items-center gap-2 mt-1">
                              <span className="text-xs text-gray-600">{skill.experience}</span>
                              <Badge variant="secondary" className={`text-xs ${getProficiencyColor(skill.proficiency)}`}>
                                {skill.proficiency}
                              </Badge>
                            </div>
                          </div>
                          <div className="flex gap-1">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleEdit(skill)}
                            >
                              <Pencil className="w-3.5 h-3.5" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleDelete(skill.id)}
                            >
                              <Trash2 className="w-3.5 h-3.5 text-red-600" />
                            </Button>
                          </div>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              );
            })}

            {skills.length === 0 && (
              <Card>
                <CardContent className="flex flex-col items-center justify-center py-12">
                  <Code className="w-12 h-12 text-gray-400 mb-4" />
                  <p className="text-gray-600 mb-4">スキルがまだ登録されていません</p>
                  <Button onClick={handleAdd}>
                    <Plus className="w-4 h-4 mr-2" />
                    最初のスキルを追加
                  </Button>
                </CardContent>
              </Card>
            )}
          </TabsContent>

          {categories.map(category => {
            const Icon = getCategoryIcon(category.icon);
            const categoryName = category.name;
            const categorySkills = skillsByCategory[categoryName] ?? [];

            return (
              <TabsContent key={category.id} value={categoryName} className="mt-6">
                <Card>
                  <CardHeader>
                    <div className="flex items-center gap-3">
                      <Icon className="w-5 h-5 text-gray-600" />
                      <div>
                        <CardTitle>{categoryName}</CardTitle>
                        <CardDescription>{categorySkills.length}件のスキル</CardDescription>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    {categorySkills.length > 0 ? (
                      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                        {categorySkills.map(skill => (
                          <div
                            key={skill.id}
                            className="flex items-center justify-between p-3 border rounded-lg hover:bg-gray-50 transition-colors"
                          >
                            <div className="flex-1">
                              <div className="font-medium text-gray-900">{skill.name}</div>
                              <div className="flex items-center gap-2 mt-1">
                                <span className="text-xs text-gray-600">{skill.experience}</span>
                                <Badge variant="secondary" className={`text-xs ${getProficiencyColor(skill.proficiency)}`}>
                                  {skill.proficiency}
                                </Badge>
                              </div>
                            </div>
                            <div className="flex gap-1">
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleEdit(skill)}
                              >
                                <Pencil className="w-3.5 h-3.5" />
                              </Button>
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleDelete(skill.id)}
                              >
                                <Trash2 className="w-3.5 h-3.5 text-red-600" />
                              </Button>
                            </div>
                          </div>
                        ))}
                      </div>
                    ) : (
                      <div className="text-center py-8 text-gray-600">
                        このカテゴリにスキルがありません
                      </div>
                    )}
                  </CardContent>
                </Card>
              </TabsContent>
            );
          })}
        </Tabs>

        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>
                {editingSkill ? "スキルを編集" : "スキルを追加"}
              </DialogTitle>
              <DialogDescription>
                スキル名、カテゴリ、経験年数、習熟度を入力してください
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="name">スキル名</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="TypeScript"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="category">カテゴリ</Label>
                <Select
                  value={formData.category}
                  onValueChange={(value) => setFormData({ ...formData, category: value })}
                  disabled={categoryNames.length === 0}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {categoryNames.map(cat => (
                      <SelectItem key={cat} value={cat}>{cat}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="experience">経験年数</Label>
                <Input
                  id="experience"
                  value={formData.experience}
                  onChange={(e) => setFormData({ ...formData, experience: e.target.value })}
                  placeholder="3年"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="proficiency">習熟度</Label>
                <Select
                  value={formData.proficiency}
                  onValueChange={(value) => setFormData({ ...formData, proficiency: value })}
                  disabled={proficiencyLevelNames.length === 0}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {proficiencyLevelNames.map(level => (
                      <SelectItem key={level} value={level}>{level}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
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
