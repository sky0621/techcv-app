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

type ProfileLinkMaster = {
  id: string;
  key: string;
  name: string;
  icon: string;
  placeholder: string;
  sortOrder: number;
};

type ProfileLinkMastersResponse = {
  linkMasters: ProfileLinkMaster[];
};

type CategoryForm = {
  id: string;
  name: string;
  icon: string;
  sortOrder: number;
};

type SkillMasterForm = {
  id: string;
  name: string;
  categoryId: string;
  sortOrder: number;
};

type ProficiencyLevelForm = {
  name: string;
  sortOrder: number;
};

type EmploymentTypeForm = {
  id: string;
  name: string;
  sortOrder: number;
};

type ProfileLinkMasterForm = {
  id: string;
  key: string;
  name: string;
  icon: string;
  placeholder: string;
  sortOrder: number;
};

const iconExamples = "code, database, cloud, wrench, frame";
const emptyCategoryForm: CategoryForm = {
  id: "",
  name: "",
  icon: "wrench",
  sortOrder: 0,
};
const emptySkillMasterForm: SkillMasterForm = {
  id: "",
  name: "",
  categoryId: "",
  sortOrder: 0,
};
const emptyProficiencyLevelForm: ProficiencyLevelForm = {
  name: "",
  sortOrder: 0,
};
const emptyEmploymentTypeForm: EmploymentTypeForm = {
  id: "",
  name: "",
  sortOrder: 0,
};
const emptyProfileLinkMasterForm: ProfileLinkMasterForm = {
  id: "",
  key: "",
  name: "",
  icon: "link",
  placeholder: "",
  sortOrder: 0,
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
  const [profileLinkMasters, setProfileLinkMasters] = useState<ProfileLinkMaster[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSavingCategory, setIsSavingCategory] = useState(false);
  const [isSavingSkillMaster, setIsSavingSkillMaster] = useState(false);
  const [isSavingProficiencyLevel, setIsSavingProficiencyLevel] = useState(false);
  const [isSavingEmploymentType, setIsSavingEmploymentType] = useState(false);
  const [isSavingProfileLinkMaster, setIsSavingProfileLinkMaster] = useState(false);
  const [categoryDialogOpen, setCategoryDialogOpen] = useState(false);
  const [skillMasterDialogOpen, setSkillMasterDialogOpen] = useState(false);
  const [proficiencyLevelDialogOpen, setProficiencyLevelDialogOpen] = useState(false);
  const [employmentTypeDialogOpen, setEmploymentTypeDialogOpen] = useState(false);
  const [profileLinkMasterDialogOpen, setProfileLinkMasterDialogOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState<SkillOption | null>(null);
  const [editingSkillMaster, setEditingSkillMaster] = useState<SkillMaster | null>(null);
  const [editingProficiencyLevel, setEditingProficiencyLevel] = useState<SkillOption | null>(null);
  const [editingEmploymentType, setEditingEmploymentType] = useState<JobEmploymentType | null>(null);
  const [editingProfileLinkMaster, setEditingProfileLinkMaster] = useState<ProfileLinkMaster | null>(null);
  const [categoryForm, setCategoryForm] = useState<CategoryForm>(emptyCategoryForm);
  const [skillMasterForm, setSkillMasterForm] = useState<SkillMasterForm>(emptySkillMasterForm);
  const [proficiencyLevelForm, setProficiencyLevelForm] = useState<ProficiencyLevelForm>(emptyProficiencyLevelForm);
  const [employmentTypeForm, setEmploymentTypeForm] = useState<EmploymentTypeForm>(emptyEmploymentTypeForm);
  const [profileLinkMasterForm, setProfileLinkMasterForm] = useState<ProfileLinkMasterForm>(emptyProfileLinkMasterForm);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  async function loadOptions(signal?: AbortSignal) {
    setIsLoading(true);
    setError("");

    try {
      const [skillOptionsResponse, jobHistoryOptionsResponse, profileLinkMastersResponse] = await Promise.all([
        fetch("/api/skills/options", { signal }),
        fetch("/api/job-histories/options", { signal }),
        fetch("/api/profile/link-masters", { signal }),
      ]);
      if (!skillOptionsResponse.ok || !jobHistoryOptionsResponse.ok || !profileLinkMastersResponse.ok) {
        throw new Error("設定情報の取得に失敗しました");
      }

      const skillOptions = (await skillOptionsResponse.json()) as SkillOptionsResponse;
      const jobHistoryOptions = (await jobHistoryOptionsResponse.json()) as JobHistoryOptionsResponse;
      const profileLinkMastersData = (await profileLinkMastersResponse.json()) as ProfileLinkMastersResponse;
      setCategories(skillOptions.categories ?? []);
      setProficiencyLevels(skillOptions.proficiencyLevels ?? []);
      setSkillMasters(skillOptions.skillMasters ?? []);
      setEmploymentTypes(jobHistoryOptions.employmentTypes ?? []);
      setProfileLinkMasters(profileLinkMastersData.linkMasters ?? []);
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

  const openAddProfileLinkMasterDialog = () => {
    setEditingProfileLinkMaster(null);
    setProfileLinkMasterForm(emptyProfileLinkMasterForm);
    setProfileLinkMasterDialogOpen(true);
    setMessage("");
    setError("");
  };

  const openEditProfileLinkMasterDialog = (linkMaster: ProfileLinkMaster) => {
    setEditingProfileLinkMaster(linkMaster);
    setProfileLinkMasterForm({
      id: linkMaster.id,
      key: linkMaster.key,
      name: linkMaster.name,
      icon: linkMaster.icon,
      placeholder: linkMaster.placeholder,
      sortOrder: linkMaster.sortOrder,
    });
    setProfileLinkMasterDialogOpen(true);
    setMessage("");
    setError("");
  };

  const openEditCategoryDialog = (category: SkillOption) => {
    setEditingCategory(category);
    setCategoryForm({
      id: category.id,
      name: category.name,
      icon: category.icon || "wrench",
      sortOrder: category.sortOrder,
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
      sortOrder: skillMaster.sortOrder,
    });
    setSkillMasterDialogOpen(true);
    setMessage("");
    setError("");
  };

  const openEditProficiencyLevelDialog = (level: SkillOption) => {
    setEditingProficiencyLevel(level);
    setProficiencyLevelForm({
      name: level.name,
      sortOrder: level.sortOrder,
    });
    setProficiencyLevelDialogOpen(true);
    setMessage("");
    setError("");
  };

  const openEditEmploymentTypeDialog = (employmentType: JobEmploymentType) => {
    setEditingEmploymentType(employmentType);
    setEmploymentTypeForm({
      id: employmentType.id,
      name: employmentType.name,
      sortOrder: employmentType.sortOrder,
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

  const handleSaveProficiencyLevel = async () => {
    if (!editingProficiencyLevel) {
      return;
    }

    setIsSavingProficiencyLevel(true);
    setMessage("");
    setError("");

    try {
      const response = await fetch(
        `/api/skills/proficiency-levels/${encodeURIComponent(editingProficiencyLevel.id)}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(proficiencyLevelForm),
        },
      );
      if (!response.ok) {
        throw new Error("習熟度の更新に失敗しました");
      }

      setProficiencyLevelDialogOpen(false);
      await loadOptions();
      setMessage("習熟度を更新しました");
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "習熟度の保存に失敗しました");
    } finally {
      setIsSavingProficiencyLevel(false);
    }
  };

  const canSaveProficiencyLevel =
    proficiencyLevelForm.name.trim() !== "";

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

  const handleSaveProfileLinkMaster = async () => {
    setIsSavingProfileLinkMaster(true);
    setMessage("");
    setError("");

    try {
      const response = await fetch(
        editingProfileLinkMaster
          ? `/api/profile/link-masters/${encodeURIComponent(editingProfileLinkMaster.id)}`
          : "/api/profile/link-masters",
        {
          method: editingProfileLinkMaster ? "PUT" : "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(profileLinkMasterForm),
        },
      );
      if (!response.ok) {
        throw new Error(editingProfileLinkMaster ? "SNSリンク種別の更新に失敗しました" : "SNSリンク種別の追加に失敗しました");
      }

      setProfileLinkMasterDialogOpen(false);
      await loadOptions();
      setMessage(editingProfileLinkMaster ? "SNSリンク種別を更新しました" : "SNSリンク種別を追加しました");
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "SNSリンク種別の保存に失敗しました");
    } finally {
      setIsSavingProfileLinkMaster(false);
    }
  };

  const canSaveProfileLinkMaster =
    profileLinkMasterForm.key.trim() !== "" &&
    profileLinkMasterForm.name.trim() !== "" &&
    profileLinkMasterForm.icon.trim() !== "" &&
    (editingProfileLinkMaster !== null || profileLinkMasterForm.id.trim() !== "");

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
                <CardTitle>SNS・外部リンク</CardTitle>
                <CardDescription>profile_link_masters テーブルの内容</CardDescription>
              </div>
              <Button type="button" variant="outline" onClick={openAddProfileLinkMasterDialog}>
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
                    <TableHead>キー</TableHead>
                    <TableHead>ID</TableHead>
                    <TableHead className="text-right">操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {profileLinkMasters.map((linkMaster) => {
                    const Icon = getIcon(linkMaster.icon);

                    return (
                      <TableRow key={linkMaster.id}>
                        <TableCell>{linkMaster.sortOrder}</TableCell>
                        <TableCell className="font-medium">{linkMaster.name}</TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <Icon className="w-4 h-4 text-gray-600" />
                            <span>{getLucideIconLabel(linkMaster.icon)}</span>
                          </div>
                        </TableCell>
                        <TableCell className="font-mono text-xs text-gray-600">{linkMaster.key}</TableCell>
                        <TableCell className="font-mono text-xs text-gray-600">{linkMaster.id}</TableCell>
                        <TableCell className="text-right">
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            onClick={() => openEditProfileLinkMasterDialog(linkMaster)}
                          >
                            <LucideIcons.Pencil className="w-4 h-4 mr-2" />
                            編集
                          </Button>
                        </TableCell>
                      </TableRow>
                    );
                  })}
                  {!isLoading && profileLinkMasters.length === 0 && (
                    <TableRow>
                      <TableCell colSpan={6} className="text-center text-gray-500">
                        SNS・外部リンク種別は登録されていません。
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
                    <TableHead className="text-right">操作</TableHead>
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
                      <TableCell className="text-right">
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => openEditProficiencyLevelDialog(level)}
                        >
                          <LucideIcons.Pencil className="w-4 h-4 mr-2" />
                          編集
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                  {!isLoading && proficiencyLevels.length === 0 && (
                    <TableRow>
                      <TableCell colSpan={4} className="text-center text-gray-500">
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

        <Dialog open={profileLinkMasterDialogOpen} onOpenChange={setProfileLinkMasterDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editingProfileLinkMaster ? "SNSリンク種別を編集" : "SNSリンク種別を追加"}</DialogTitle>
              <DialogDescription>
                {editingProfileLinkMaster
                  ? "プロフィールのSNS・リンクタブに表示する名称、アイコン、プレースホルダー、表示順を変更します。"
                  : "プロフィールのSNS・リンクタブに表示するID、キー、名称、アイコン、プレースホルダー、表示順を登録します。"}
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              {!editingProfileLinkMaster && (
                <div className="space-y-2">
                  <Label htmlFor="profile-link-master-id">ID</Label>
                  <Input
                    id="profile-link-master-id"
                    type="number"
                    min={1}
                    value={profileLinkMasterForm.id}
                    onChange={(e) =>
                      setProfileLinkMasterForm({ ...profileLinkMasterForm, id: e.target.value })
                    }
                    placeholder="5"
                  />
                </div>
              )}

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="profile-link-master-key">キー</Label>
                  <Input
                    id="profile-link-master-key"
                    value={profileLinkMasterForm.key}
                    onChange={(e) =>
                      setProfileLinkMasterForm({ ...profileLinkMasterForm, key: e.target.value })
                    }
                    placeholder="x"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="profile-link-master-name">名称</Label>
                  <Input
                    id="profile-link-master-name"
                    value={profileLinkMasterForm.name}
                    onChange={(e) =>
                      setProfileLinkMasterForm({ ...profileLinkMasterForm, name: e.target.value })
                    }
                    placeholder="X"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="profile-link-master-icon">アイコン</Label>
                <div className="flex items-center gap-3">
                  <Input
                    id="profile-link-master-icon"
                    value={profileLinkMasterForm.icon}
                    onChange={(e) =>
                      setProfileLinkMasterForm({ ...profileLinkMasterForm, icon: e.target.value })
                    }
                    placeholder="link"
                  />
                  <div className="flex h-9 min-w-28 items-center gap-2 rounded-md border bg-gray-50 px-3 text-sm text-gray-700">
                    {(() => {
                      const Icon = getIcon(profileLinkMasterForm.icon);

                      return <Icon className="w-4 h-4 text-gray-600" />;
                    })()}
                    <span>{getLucideIconLabel(profileLinkMasterForm.icon)}</span>
                  </div>
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="profile-link-master-placeholder">未入力時の文言</Label>
                <Input
                  id="profile-link-master-placeholder"
                  value={profileLinkMasterForm.placeholder}
                  onChange={(e) =>
                    setProfileLinkMasterForm({ ...profileLinkMasterForm, placeholder: e.target.value })
                  }
                  placeholder="https://example.com/username"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="profile-link-master-sort-order">表示順</Label>
                <Input
                  id="profile-link-master-sort-order"
                  type="number"
                  value={profileLinkMasterForm.sortOrder}
                  onChange={(e) =>
                    setProfileLinkMasterForm({ ...profileLinkMasterForm, sortOrder: Number(e.target.value) })
                  }
                  placeholder="5"
                />
              </div>
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setProfileLinkMasterDialogOpen(false)}
              >
                キャンセル
              </Button>
              <Button
                type="button"
                onClick={handleSaveProfileLinkMaster}
                disabled={!canSaveProfileLinkMaster || isSavingProfileLinkMaster}
              >
                {isSavingProfileLinkMaster ? "保存中" : "保存"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        <Dialog open={categoryDialogOpen} onOpenChange={setCategoryDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editingCategory ? "カテゴリを編集" : "カテゴリを追加"}</DialogTitle>
              <DialogDescription>
                {editingCategory
                  ? "カテゴリの名称、アイコン、表示順を変更します。"
                  : "カテゴリのID、名称、アイコン、表示順を登録します。"}
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              {!editingCategory && (
                <div className="space-y-2">
                  <Label htmlFor="category-id">ID</Label>
                  <Input
                    id="category-id"
                    type="number"
                    min={1}
                    value={categoryForm.id}
                    onChange={(e) =>
                      setCategoryForm({ ...categoryForm, id: e.target.value })
                    }
                    placeholder="7"
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

              <div className="space-y-2">
                <Label htmlFor="category-sort-order">表示順</Label>
                <Input
                  id="category-sort-order"
                  type="number"
                  value={categoryForm.sortOrder}
                  onChange={(e) =>
                    setCategoryForm({ ...categoryForm, sortOrder: Number(e.target.value) })
                  }
                  placeholder="7"
                />
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
                  ? "スキル名、カテゴリ、表示順を変更します。"
                  : "スキルマスタのID、名称、カテゴリ、表示順を登録します。"}
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              {!editingSkillMaster && (
                <div className="space-y-2">
                  <Label htmlFor="skill-master-id">ID</Label>
                  <Input
                    id="skill-master-id"
                    type="number"
                    min={1}
                    value={skillMasterForm.id}
                    onChange={(e) =>
                      setSkillMasterForm({ ...skillMasterForm, id: e.target.value })
                    }
                    placeholder="12"
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

              <div className="space-y-2">
                <Label htmlFor="skill-master-sort-order">表示順</Label>
                <Input
                  id="skill-master-sort-order"
                  type="number"
                  value={skillMasterForm.sortOrder}
                  onChange={(e) =>
                    setSkillMasterForm({ ...skillMasterForm, sortOrder: Number(e.target.value) })
                  }
                  placeholder="12"
                />
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

        <Dialog open={proficiencyLevelDialogOpen} onOpenChange={setProficiencyLevelDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>習熟度を編集</DialogTitle>
              <DialogDescription>
                習熟度の名称と表示順を変更します。
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="proficiency-level-name">名称</Label>
                <Input
                  id="proficiency-level-name"
                  value={proficiencyLevelForm.name}
                  onChange={(e) =>
                    setProficiencyLevelForm({ ...proficiencyLevelForm, name: e.target.value })
                  }
                  placeholder="中級"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="proficiency-level-sort-order">表示順</Label>
                <Input
                  id="proficiency-level-sort-order"
                  type="number"
                  value={proficiencyLevelForm.sortOrder}
                  onChange={(e) =>
                    setProficiencyLevelForm({ ...proficiencyLevelForm, sortOrder: Number(e.target.value) })
                  }
                  placeholder="2"
                />
              </div>
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setProficiencyLevelDialogOpen(false)}
              >
                キャンセル
              </Button>
              <Button
                type="button"
                onClick={handleSaveProficiencyLevel}
                disabled={!canSaveProficiencyLevel || isSavingProficiencyLevel}
              >
                {isSavingProficiencyLevel ? "保存中" : "保存"}
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
                  ? "雇用形態の名称と表示順を変更します。"
                  : "雇用形態のID、名称、表示順を登録します。"}
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              {!editingEmploymentType && (
                <div className="space-y-2">
                  <Label htmlFor="employment-type-id">ID</Label>
                  <Input
                    id="employment-type-id"
                    type="number"
                    min={1}
                    value={employmentTypeForm.id}
                    onChange={(e) =>
                      setEmploymentTypeForm({ ...employmentTypeForm, id: e.target.value })
                    }
                    placeholder="6"
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

              <div className="space-y-2">
                <Label htmlFor="employment-type-sort-order">表示順</Label>
                <Input
                  id="employment-type-sort-order"
                  type="number"
                  value={employmentTypeForm.sortOrder}
                  onChange={(e) =>
                    setEmploymentTypeForm({ ...employmentTypeForm, sortOrder: Number(e.target.value) })
                  }
                  placeholder="6"
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
