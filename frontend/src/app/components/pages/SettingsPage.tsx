import { useEffect, useState } from "react";
import * as LucideIcons from "lucide-react";
import type { LucideIcon } from "lucide-react";
import { Button } from "../ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "../ui/dialog";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select";

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

type JobEmploymentType = {
  id: string;
  name: string;
  sortOrder: number;
};

type JobHistoryOptionsResponse = {
  employmentTypes: JobEmploymentType[];
};

type CategoryForm = {
  id: string;
  name: string;
  icon: string;
};

type SkillMasterForm = {
  id: string;
  name: string;
  categoryId: string;
};

type EmploymentTypeForm = {
  id: string;
  name: string;
};

const iconExamples = "code, database, cloud, wrench, frame";
const emptyCategoryForm: CategoryForm = {
  id: "",
  name: "",
  icon: "wrench",
};
const emptySkillMasterForm: SkillMasterForm = {
  id: "",
  name: "",
  categoryId: "",
};
const emptyEmploymentTypeForm: EmploymentTypeForm = {
  id: "",
  name: "",
};

function normalizeIconName(icon?: string) {
  return (icon ?? "").trim().replace(/^lucide-/i, "");
}

function getLucideIconLabel(icon?: string) {
  const normalizedIcon = normalizeIconName(icon);

  return `lucide-${normalizedIcon || "wrench"}`;
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

function getIcon(icon?: string): LucideIcon {
  const candidate = (LucideIcons as Record<string, unknown>)[toLucideExportName(icon)];

  return isLucideIcon(candidate) ? candidate : LucideIcons.Wrench;
}

export function SettingsPage() {
  const [categories, setCategories] = useState<SkillOption[]>([]);
  const [proficiencyLevels, setProficiencyLevels] = useState<SkillOption[]>([]);
  const [skillMasters, setSkillMasters] = useState<SkillMaster[]>([]);
  const [employmentTypes, setEmploymentTypes] = useState<JobEmploymentType[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSavingCategory, setIsSavingCategory] = useState(false);
  const [isSavingSkillMaster, setIsSavingSkillMaster] = useState(false);
  const [isSavingEmploymentType, setIsSavingEmploymentType] = useState(false);
  const [categoryDialogOpen, setCategoryDialogOpen] = useState(false);
  const [skillMasterDialogOpen, setSkillMasterDialogOpen] = useState(false);
  const [employmentTypeDialogOpen, setEmploymentTypeDialogOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState<SkillOption | null>(null);
  const [editingSkillMaster, setEditingSkillMaster] = useState<SkillMaster | null>(null);
  const [editingEmploymentType, setEditingEmploymentType] = useState<JobEmploymentType | null>(null);
  const [categoryForm, setCategoryForm] = useState<CategoryForm>(emptyCategoryForm);
  const [skillMasterForm, setSkillMasterForm] = useState<SkillMasterForm>(emptySkillMasterForm);
  const [employmentTypeForm, setEmploymentTypeForm] = useState<EmploymentTypeForm>(emptyEmploymentTypeForm);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  async function loadOptions(signal?: AbortSignal) {
    setIsLoading(true);
    setError("");

    try {
      const [skillOptionsResponse, jobHistoryOptionsResponse] = await Promise.all([
        fetch("/api/skills/options", { signal }),
        fetch("/api/job-histories/options", { signal }),
      ]);
      if (!skillOptionsResponse.ok || !jobHistoryOptionsResponse.ok) {
        throw new Error("設定情報の取得に失敗しました");
      }

      const skillOptions = (await skillOptionsResponse.json()) as SkillOptionsResponse;
      const jobHistoryOptions = (await jobHistoryOptionsResponse.json()) as JobHistoryOptionsResponse;
      setCategories(skillOptions.categories ?? []);
      setProficiencyLevels(skillOptions.proficiencyLevels ?? []);
      setSkillMasters(skillOptions.skillMasters ?? []);
      setEmploymentTypes(jobHistoryOptions.employmentTypes ?? []);
    } catch (caught) {
      if (caught instanceof DOMException && caught.name === "AbortError") {
        return;
      }
      setError(caught instanceof Error ? caught.message : "設定情報の取得に失敗しました");
    } finally {
      if (!signal?.aborted) {
        setIsLoading(false);
      }
    }
  }

  useEffect(() => {
    const controller = new AbortController();

    void loadOptions(controller.signal);

    return () => {
      controller.abort();
    };
  }, []);

  const openAddCategoryDialog = () => {
    setEditingCategory(null);
    setCategoryForm(emptyCategoryForm);
    setCategoryDialogOpen(true);
    setMessage("");
    setError("");
  };

  const openEditCategoryDialog = (category: SkillOption) => {
    setEditingCategory(category);
    setCategoryForm({
      id: category.id,
      name: category.name,
      icon: category.icon || "wrench",
    });
    setCategoryDialogOpen(true);
    setMessage("");
    setError("");
  };

  const openAddEmploymentTypeDialog = () => {
    setEditingEmploymentType(null);
    setEmploymentTypeForm(emptyEmploymentTypeForm);
    setEmploymentTypeDialogOpen(true);
    setMessage("");
    setError("");
  };

  const openAddSkillMasterDialog = () => {
    setEditingSkillMaster(null);
    setSkillMasterForm({
      ...emptySkillMasterForm,
      categoryId: categories[0]?.id ?? "",
    });
    setSkillMasterDialogOpen(true);
    setMessage("");
    setError("");
  };

  const openEditSkillMasterDialog = (skillMaster: SkillMaster) => {
    setEditingSkillMaster(skillMaster);
    setSkillMasterForm({
      id: skillMaster.id,
      name: skillMaster.name,
      categoryId: skillMaster.categoryId,
    });
    setSkillMasterDialogOpen(true);
    setMessage("");
    setError("");
  };

  const openEditEmploymentTypeDialog = (employmentType: JobEmploymentType) => {
    setEditingEmploymentType(employmentType);
    setEmploymentTypeForm({
      id: employmentType.id,
      name: employmentType.name,
    });
    setEmploymentTypeDialogOpen(true);
    setMessage("");
    setError("");
  };

  const handleSaveCategory = async () => {
    setIsSavingCategory(true);
    setMessage("");
    setError("");

    try {
      const response = await fetch(
        editingCategory
          ? `/api/skills/categories/${encodeURIComponent(editingCategory.id)}`
          : "/api/skills/categories",
        {
          method: editingCategory ? "PUT" : "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(categoryForm),
        },
      );
      if (!response.ok) {
        throw new Error(editingCategory ? "カテゴリの更新に失敗しました" : "カテゴリの追加に失敗しました");
      }

      setCategoryDialogOpen(false);
      await loadOptions();
      setMessage(editingCategory ? "カテゴリを更新しました" : "カテゴリを追加しました");
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "カテゴリの保存に失敗しました");
    } finally {
      setIsSavingCategory(false);
    }
  };

  const canSaveCategory =
    categoryForm.name.trim() !== "" &&
    categoryForm.icon.trim() !== "" &&
    (editingCategory !== null || categoryForm.id.trim() !== "");

  const handleSaveSkillMaster = async () => {
    setIsSavingSkillMaster(true);
    setMessage("");
    setError("");

    try {
      const response = await fetch(
        editingSkillMaster
          ? `/api/skills/masters/${encodeURIComponent(editingSkillMaster.id)}`
          : "/api/skills/masters",
        {
          method: editingSkillMaster ? "PUT" : "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(skillMasterForm),
        },
      );
      if (!response.ok) {
        throw new Error(editingSkillMaster ? "スキルマスタの更新に失敗しました" : "スキルマスタの追加に失敗しました");
      }

      setSkillMasterDialogOpen(false);
      await loadOptions();
      setMessage(editingSkillMaster ? "スキルマスタを更新しました" : "スキルマスタを追加しました");
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "スキルマスタの保存に失敗しました");
    } finally {
      setIsSavingSkillMaster(false);
    }
  };

  const canSaveSkillMaster =
    skillMasterForm.name.trim() !== "" &&
    skillMasterForm.categoryId.trim() !== "" &&
    (editingSkillMaster !== null || skillMasterForm.id.trim() !== "");

  const handleSaveEmploymentType = async () => {
    setIsSavingEmploymentType(true);
    setMessage("");
    setError("");

    try {
      const response = await fetch(
        editingEmploymentType
          ? `/api/job-employment-types/${encodeURIComponent(editingEmploymentType.id)}`
          : "/api/job-employment-types",
        {
          method: editingEmploymentType ? "PUT" : "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(employmentTypeForm),
        },
      );
      if (!response.ok) {
        throw new Error(editingEmploymentType ? "雇用形態の更新に失敗しました" : "雇用形態の追加に失敗しました");
      }

      setEmploymentTypeDialogOpen(false);
      await loadOptions();
      setMessage(editingEmploymentType ? "雇用形態を更新しました" : "雇用形態を追加しました");
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "雇用形態の保存に失敗しました");
    } finally {
      setIsSavingEmploymentType(false);
    }
  };

  const canSaveEmploymentType =
    employmentTypeForm.name.trim() !== "" &&
    (editingEmploymentType !== null || employmentTypeForm.id.trim() !== "");

  return (
    <div className="p-8">
      <div className="max-w-6xl mx-auto space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">設定</h1>
          <p className="text-gray-600 mt-1">スキル管理で利用するマスタ情報</p>
        </div>

        {isLoading && (
          <div className="rounded-md border border-gray-200 bg-gray-50 px-4 py-3 text-sm text-gray-600">
            設定情報を読み込み中です。
          </div>
        )}
        {error && (
          <div className="rounded-md border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {error}
          </div>
        )}
        {message && (
          <div className="rounded-md border border-green-200 bg-green-50 px-4 py-3 text-sm text-green-700">
            {message}
          </div>
        )}

        <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between gap-4">
              <div>
                <CardTitle>スキルカテゴリ</CardTitle>
                <CardDescription>skill_categories テーブルの内容</CardDescription>
              </div>
              <Button type="button" variant="outline" onClick={openAddCategoryDialog}>
                <LucideIcons.Plus className="w-4 h-4 mr-2" />
                追加
              </Button>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>表示順</TableHead>
                    <TableHead>名称</TableHead>
                    <TableHead>アイコン</TableHead>
                    <TableHead>ID</TableHead>
                    <TableHead className="text-right">操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {categories.map((category) => {
                    const Icon = getIcon(category.icon);

                    return (
                      <TableRow key={category.id}>
                        <TableCell>{category.sortOrder}</TableCell>
                        <TableCell className="font-medium">{category.name}</TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <Icon className="w-4 h-4 text-gray-600" />
                            <span>{getLucideIconLabel(category.icon)}</span>
                          </div>
                        </TableCell>
                        <TableCell className="font-mono text-xs text-gray-600">
                          {category.id}
                        </TableCell>
                        <TableCell className="text-right">
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            onClick={() => openEditCategoryDialog(category)}
                          >
                            <LucideIcons.Pencil className="w-4 h-4 mr-2" />
                            編集
                          </Button>
                        </TableCell>
                      </TableRow>
                    );
                  })}
                  {!isLoading && categories.length === 0 && (
                    <TableRow>
                      <TableCell colSpan={5} className="text-center text-gray-500">
                        スキルカテゴリは登録されていません。
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between gap-4">
              <div>
                <CardTitle>スキルマスタ</CardTitle>
                <CardDescription>skill_masters テーブルの内容</CardDescription>
              </div>
              <Button
                type="button"
                variant="outline"
                onClick={openAddSkillMasterDialog}
                disabled={categories.length === 0}
              >
                <LucideIcons.Plus className="w-4 h-4 mr-2" />
                追加
              </Button>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>表示順</TableHead>
                    <TableHead>名称</TableHead>
                    <TableHead>カテゴリ</TableHead>
                    <TableHead>ID</TableHead>
                    <TableHead className="text-right">操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {skillMasters.map((skillMaster) => (
                    <TableRow key={skillMaster.id}>
                      <TableCell>{skillMaster.sortOrder}</TableCell>
                      <TableCell className="font-medium">{skillMaster.name}</TableCell>
                      <TableCell>{skillMaster.category}</TableCell>
                      <TableCell className="font-mono text-xs text-gray-600">
                        {skillMaster.id}
                      </TableCell>
                      <TableCell className="text-right">
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => openEditSkillMasterDialog(skillMaster)}
                        >
                          <LucideIcons.Pencil className="w-4 h-4 mr-2" />
                          編集
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                  {!isLoading && skillMasters.length === 0 && (
                    <TableRow>
                      <TableCell colSpan={5} className="text-center text-gray-500">
                        スキルマスタは登録されていません。
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>習熟度</CardTitle>
              <CardDescription>skill_proficiency_levels テーブルの内容</CardDescription>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>表示順</TableHead>
                    <TableHead>名称</TableHead>
                    <TableHead>ID</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {proficiencyLevels.map((level) => (
                    <TableRow key={level.id}>
                      <TableCell>{level.sortOrder}</TableCell>
                      <TableCell className="font-medium">{level.name}</TableCell>
                      <TableCell className="font-mono text-xs text-gray-600">
                        {level.id}
                      </TableCell>
                    </TableRow>
                  ))}
                  {!isLoading && proficiencyLevels.length === 0 && (
                    <TableRow>
                      <TableCell colSpan={3} className="text-center text-gray-500">
                        習熟度は登録されていません。
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between gap-4">
              <div>
                <CardTitle>雇用形態</CardTitle>
                <CardDescription>job_employment_types テーブルの内容</CardDescription>
              </div>
              <Button type="button" variant="outline" onClick={openAddEmploymentTypeDialog}>
                <LucideIcons.Plus className="w-4 h-4 mr-2" />
                追加
              </Button>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>表示順</TableHead>
                    <TableHead>名称</TableHead>
                    <TableHead>ID</TableHead>
                    <TableHead className="text-right">操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {employmentTypes.map((employmentType) => (
                    <TableRow key={employmentType.id}>
                      <TableCell>{employmentType.sortOrder}</TableCell>
                      <TableCell className="font-medium">{employmentType.name}</TableCell>
                      <TableCell className="font-mono text-xs text-gray-600">
                        {employmentType.id}
                      </TableCell>
                      <TableCell className="text-right">
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => openEditEmploymentTypeDialog(employmentType)}
                        >
                          <LucideIcons.Pencil className="w-4 h-4 mr-2" />
                          編集
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                  {!isLoading && employmentTypes.length === 0 && (
                    <TableRow>
                      <TableCell colSpan={4} className="text-center text-gray-500">
                        雇用形態は登録されていません。
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </div>

        <Dialog open={categoryDialogOpen} onOpenChange={setCategoryDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editingCategory ? "カテゴリを編集" : "カテゴリを追加"}</DialogTitle>
              <DialogDescription>
                {editingCategory
                  ? "カテゴリの名称とアイコンを変更します。"
                  : "カテゴリのID、名称、アイコンを登録します。"}
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              {!editingCategory && (
                <div className="space-y-2">
                  <Label htmlFor="category-id">ID</Label>
                  <Input
                    id="category-id"
                    value={categoryForm.id}
                    onChange={(e) =>
                      setCategoryForm({ ...categoryForm, id: e.target.value })
                    }
                    placeholder="skill_category_backend"
                  />
                </div>
              )}

              <div className="space-y-2">
                <Label htmlFor="category-name">名称</Label>
                <Input
                  id="category-name"
                  value={categoryForm.name}
                  onChange={(e) =>
                    setCategoryForm({ ...categoryForm, name: e.target.value })
                  }
                  placeholder="バックエンド"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="category-icon">アイコン</Label>
                <div className="flex items-center gap-3">
                  <Input
                    id="category-icon"
                    value={categoryForm.icon}
                    onChange={(e) =>
                      setCategoryForm({ ...categoryForm, icon: e.target.value })
                    }
                    placeholder="code"
                  />
                  <div className="flex h-9 min-w-28 items-center gap-2 rounded-md border bg-gray-50 px-3 text-sm text-gray-700">
                    {(() => {
                      const Icon = getIcon(categoryForm.icon);

                      return <Icon className="w-4 h-4 text-gray-600" />;
                    })()}
                    <span>{getLucideIconLabel(categoryForm.icon)}</span>
                  </div>
                </div>
                <p className="text-xs text-gray-500">
                  例: {iconExamples}。プレビューでは lucide- と組み合わせて表示します。
                </p>
              </div>
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setCategoryDialogOpen(false)}
              >
                キャンセル
              </Button>
              <Button
                type="button"
                onClick={handleSaveCategory}
                disabled={!canSaveCategory || isSavingCategory}
              >
                {isSavingCategory ? "保存中" : "保存"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        <Dialog open={skillMasterDialogOpen} onOpenChange={setSkillMasterDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editingSkillMaster ? "スキルマスタを編集" : "スキルマスタを追加"}</DialogTitle>
              <DialogDescription>
                {editingSkillMaster
                  ? "スキル名とカテゴリを変更します。"
                  : "スキルマスタのID、名称、カテゴリを登録します。"}
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              {!editingSkillMaster && (
                <div className="space-y-2">
                  <Label htmlFor="skill-master-id">ID</Label>
                  <Input
                    id="skill-master-id"
                    value={skillMasterForm.id}
                    onChange={(e) =>
                      setSkillMasterForm({ ...skillMasterForm, id: e.target.value })
                    }
                    placeholder="skill_master_kotlin"
                  />
                </div>
              )}

              <div className="space-y-2">
                <Label htmlFor="skill-master-name">名称</Label>
                <Input
                  id="skill-master-name"
                  value={skillMasterForm.name}
                  onChange={(e) =>
                    setSkillMasterForm({ ...skillMasterForm, name: e.target.value })
                  }
                  placeholder="Kotlin"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="skill-master-category">カテゴリ</Label>
                <Select
                  value={skillMasterForm.categoryId}
                  onValueChange={(value) =>
                    setSkillMasterForm({ ...skillMasterForm, categoryId: value })
                  }
                  disabled={categories.length === 0}
                >
                  <SelectTrigger id="skill-master-category">
                    <SelectValue placeholder="カテゴリを選択" />
                  </SelectTrigger>
                  <SelectContent>
                    {categories.map((category) => (
                      <SelectItem key={category.id} value={category.id}>
                        {category.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setSkillMasterDialogOpen(false)}
              >
                キャンセル
              </Button>
              <Button
                type="button"
                onClick={handleSaveSkillMaster}
                disabled={!canSaveSkillMaster || isSavingSkillMaster}
              >
                {isSavingSkillMaster ? "保存中" : "保存"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        <Dialog open={employmentTypeDialogOpen} onOpenChange={setEmploymentTypeDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editingEmploymentType ? "雇用形態を編集" : "雇用形態を追加"}</DialogTitle>
              <DialogDescription>
                {editingEmploymentType
                  ? "雇用形態の名称を変更します。"
                  : "雇用形態のIDと名称を登録します。"}
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              {!editingEmploymentType && (
                <div className="space-y-2">
                  <Label htmlFor="employment-type-id">ID</Label>
                  <Input
                    id="employment-type-id"
                    value={employmentTypeForm.id}
                    onChange={(e) =>
                      setEmploymentTypeForm({ ...employmentTypeForm, id: e.target.value })
                    }
                    placeholder="job_employment_type_intern"
                  />
                </div>
              )}

              <div className="space-y-2">
                <Label htmlFor="employment-type-name">名称</Label>
                <Input
                  id="employment-type-name"
                  value={employmentTypeForm.name}
                  onChange={(e) =>
                    setEmploymentTypeForm({ ...employmentTypeForm, name: e.target.value })
                  }
                  placeholder="インターン"
                />
              </div>
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setEmploymentTypeDialogOpen(false)}
              >
                キャンセル
              </Button>
              <Button
                type="button"
                onClick={handleSaveEmploymentType}
                disabled={!canSaveEmploymentType || isSavingEmploymentType}
              >
                {isSavingEmploymentType ? "保存中" : "保存"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
}
