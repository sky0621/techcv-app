import { useEffect, useState } from "react";
import * as LucideIcons from "lucide-react";
import type { LucideIcon } from "lucide-react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select";
import { Badge } from "../ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";

type Skill = {
  id: string;
  name: string;
  categoryId: string;
  category: string;
  experience: number;
  proficiencyLevelId: string;
  proficiency: string;
  sortOrder: number;
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
  skillMasters: SkillMaster[];
};

type SkillMaster = {
  id: string;
  name: string;
  categoryId: string;
  category: string;
  sortOrder: number;
};

type SkillsResponse = {
  skills: Skill[];
};

type SkillResponse = {
  skill: Skill;
};

type SkillFormData = {
  skillMasterId: string;
  name: string;
  categoryId: string;
  experience: string;
  proficiencyLevelId: string;
};

function normalizeIconName(icon?: string) {
  return (icon ?? "").trim().replace(/^lucide-/i, "");
}

function toLucideExportName(icon?: string) {
  const normalizedIcon = normalizeIconName(icon);
  if (!normalizedIcon) return "Wrench";

  return normalizedIcon
    .split(/[^a-zA-Z0-9]+/)
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join("");
}

function isLucideIcon(candidate: unknown): candidate is LucideIcon {
  return (
    candidate !== null &&
    typeof candidate === "object" &&
    "displayName" in candidate &&
    typeof candidate.displayName === "string"
  );
}

function getCategoryIcon(icon?: string): LucideIcon {
  const candidate = (LucideIcons as Record<string, unknown>)[toLucideExportName(icon)];

  return isLucideIcon(candidate) ? candidate : LucideIcons.Wrench;
}

export function SkillsPage() {
  const [skills, setSkills] = useState<Skill[]>([]);
  const [categories, setCategories] = useState<SkillOption[]>([]);
  const [proficiencyLevels, setProficiencyLevels] = useState<SkillOption[]>([]);
  const [skillMasters, setSkillMasters] = useState<SkillMaster[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingSkill, setEditingSkill] = useState<Skill | null>(null);
  const [formData, setFormData] = useState<SkillFormData>({
    skillMasterId: "",
    name: "",
    categoryId: "",
    experience: "",
    proficiencyLevelId: "",
  });

  useEffect(() => {
    const controller = new AbortController();

    async function loadSkillsPageData() {
      setIsLoading(true);
      setError("");

      try {
        const [optionsResponse, skillsResponse] = await Promise.all([
          fetch("/api/skills/options", { signal: controller.signal }),
          fetch("/api/skills", { signal: controller.signal }),
        ]);
        if (!optionsResponse.ok) {
          throw new Error("スキルの選択肢の取得に失敗しました");
        }
        if (!skillsResponse.ok) {
          throw new Error("スキルの取得に失敗しました");
        }

        const optionsData = (await optionsResponse.json()) as SkillOptionsResponse;
        const skillsData = (await skillsResponse.json()) as SkillsResponse;
        setCategories(optionsData.categories ?? []);
        setProficiencyLevels(optionsData.proficiencyLevels ?? []);
        setSkillMasters(optionsData.skillMasters ?? []);
        setSkills(skillsData.skills ?? []);
      } catch (caught) {
        if (caught instanceof DOMException && caught.name === "AbortError") {
          return;
        }
        setError(caught instanceof Error ? caught.message : "スキル情報の取得に失敗しました");
      } finally {
        if (!controller.signal.aborted) {
          setIsLoading(false);
        }
      }
    }

    void loadSkillsPageData();

    return () => {
      controller.abort();
    };
  }, []);

  const selectedDefaultSkillMaster = skillMasters[0];
  const selectedDefaultProficiencyLevelId =
    proficiencyLevels.find((level) => level.name === "中級")?.id ?? proficiencyLevels[0]?.id ?? "";

  const skillsByCategory = categories.reduce((acc, category) => {
    acc[category.id] = skills.filter(s => s.categoryId === category.id);
    return acc;
  }, {} as Record<string, Skill[]>);

  const handleAdd = () => {
    setEditingSkill(null);
    setFormData({
      skillMasterId: selectedDefaultSkillMaster?.id ?? "",
      name: selectedDefaultSkillMaster?.name ?? "",
      categoryId: selectedDefaultSkillMaster?.categoryId ?? "",
      experience: "",
      proficiencyLevelId: selectedDefaultProficiencyLevelId,
    });
    setIsDialogOpen(true);
  };

  const handleEdit = (skill: Skill) => {
    const skillMaster = skillMasters.find((master) => master.name === skill.name);

    setEditingSkill(skill);
    setFormData({
      skillMasterId: skillMaster?.id ?? "",
      name: skill.name,
      categoryId: skill.categoryId,
      experience: String(skill.experience),
      proficiencyLevelId: skill.proficiencyLevelId,
    });
    setIsDialogOpen(true);
  };

  const handleSkillMasterChange = (skillMasterId: string) => {
    const skillMaster = skillMasters.find((master) => master.id === skillMasterId);
    if (!skillMaster) {
      return;
    }

    setFormData({
      ...formData,
      skillMasterId: skillMaster.id,
      name: skillMaster.name,
      categoryId: skillMaster.categoryId,
    });
  };

  const handleSave = async () => {
    setIsSaving(true);
    setError("");

    try {
      const response = await fetch(
        editingSkill ? `/api/skills/${encodeURIComponent(editingSkill.id)}` : "/api/skills",
        {
          method: editingSkill ? "PUT" : "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            name: formData.name,
            categoryId: formData.categoryId,
            experience: Number(formData.experience),
            proficiencyLevelId: formData.proficiencyLevelId,
          }),
        },
      );
      if (!response.ok) {
        throw new Error(editingSkill ? "スキルの更新に失敗しました" : "スキルの追加に失敗しました");
      }

      const data = (await response.json()) as SkillResponse;
      if (editingSkill) {
        setSkills(skills.map(s => s.id === editingSkill.id ? data.skill : s));
      } else {
        setSkills([...skills, data.skill]);
      }
      setIsDialogOpen(false);
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "スキルの保存に失敗しました");
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (confirm("本当に削除しますか？")) {
      setError("");

      try {
        const response = await fetch(`/api/skills/${encodeURIComponent(id)}`, {
          method: "DELETE",
        });
        if (!response.ok) {
          throw new Error("スキルの削除に失敗しました");
        }

        setSkills(skills.filter(s => s.id !== id));
      } catch (caught) {
        setError(caught instanceof Error ? caught.message : "スキルの削除に失敗しました");
      }
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

  const canSave =
    formData.skillMasterId !== "" &&
    formData.name.trim() !== "" &&
    formData.categoryId !== "" &&
    formData.proficiencyLevelId !== "";

  return (
    <div className="p-8">
      <div className="max-w-6xl mx-auto space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">スキル管理</h1>
            <p className="text-gray-600 mt-1">技術スキルと習熟度を管理</p>
          </div>
          <Button onClick={handleAdd} disabled={isLoading || skillMasters.length === 0 || proficiencyLevels.length === 0}>
            <LucideIcons.Plus className="w-4 h-4 mr-2" />
            スキルを追加
          </Button>
        </div>

        {isLoading && (
          <div className="rounded-md border border-gray-200 bg-gray-50 px-4 py-3 text-sm text-gray-600">
            スキル情報を読み込み中です。
          </div>
        )}
        {error && (
          <div className="rounded-md border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {error}
          </div>
        )}
        {!isLoading && skillMasters.length === 0 && (
          <div className="rounded-md border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800">
            スキルを登録するには、先に設定画面でスキルマスタを登録してください。
          </div>
        )}

        <Tabs defaultValue="all" className="w-full">
          <TabsList>
            <TabsTrigger value="all">すべて ({skills.length})</TabsTrigger>
            {categories.map(category => (
              <TabsTrigger key={category.id} value={category.id}>
                {category.name} ({skillsByCategory[category.id]?.length ?? 0})
              </TabsTrigger>
            ))}
          </TabsList>

          <TabsContent value="all" className="space-y-4 mt-6">
            {categories.map(category => {
              const Icon = getCategoryIcon(category.icon);
              const categoryName = category.name;
              const categorySkills = skillsByCategory[category.id] ?? [];

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
                              <span className="text-xs text-gray-600">{skill.experience}年</span>
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
                              <LucideIcons.Pencil className="w-3.5 h-3.5" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleDelete(skill.id)}
                            >
                              <LucideIcons.Trash2 className="w-3.5 h-3.5 text-red-600" />
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
                  <LucideIcons.Code className="w-12 h-12 text-gray-400 mb-4" />
                  <p className="text-gray-600 mb-4">スキルがまだ登録されていません</p>
                  <Button onClick={handleAdd} disabled={skillMasters.length === 0 || proficiencyLevels.length === 0}>
                    <LucideIcons.Plus className="w-4 h-4 mr-2" />
                    最初のスキルを追加
                  </Button>
                </CardContent>
              </Card>
            )}
          </TabsContent>

          {categories.map(category => {
            const Icon = getCategoryIcon(category.icon);
            const categoryName = category.name;
            const categorySkills = skillsByCategory[category.id] ?? [];

            return (
              <TabsContent key={category.id} value={category.id} className="mt-6">
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
                                <span className="text-xs text-gray-600">{skill.experience}年</span>
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
                                <LucideIcons.Pencil className="w-3.5 h-3.5" />
                              </Button>
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleDelete(skill.id)}
                              >
                                <LucideIcons.Trash2 className="w-3.5 h-3.5 text-red-600" />
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
                <Label htmlFor="skillMaster">スキル名</Label>
                <Select
                  value={formData.skillMasterId}
                  onValueChange={handleSkillMasterChange}
                  disabled={skillMasters.length === 0}
                >
                  <SelectTrigger id="skillMaster">
                    <SelectValue placeholder="スキルマスタから選択" />
                  </SelectTrigger>
                  <SelectContent>
                    {skillMasters.map(skillMaster => (
                      <SelectItem key={skillMaster.id} value={skillMaster.id}>
                        {skillMaster.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="category">カテゴリ</Label>
                <Select
                  value={formData.categoryId}
                  onValueChange={(value) => setFormData({ ...formData, categoryId: value })}
                  disabled
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {categories.map(category => (
                      <SelectItem key={category.id} value={category.id}>{category.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="experience">経験年数</Label>
                <Input
                  id="experience"
                  type="number"
                  min={0}
                  value={formData.experience}
                  onChange={(e) => setFormData({ ...formData, experience: e.target.value })}
                  placeholder="3"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="proficiency">習熟度</Label>
                <Select
                  value={formData.proficiencyLevelId}
                  onValueChange={(value) => setFormData({ ...formData, proficiencyLevelId: value })}
                  disabled={proficiencyLevels.length === 0}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {proficiencyLevels.map(level => (
                      <SelectItem key={level.id} value={level.id}>{level.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            <DialogFooter>
              <Button variant="outline" onClick={() => setIsDialogOpen(false)}>
                キャンセル
              </Button>
              <Button onClick={handleSave} disabled={!canSave || isSaving}>
                {isSaving ? "保存中" : "保存"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
}
