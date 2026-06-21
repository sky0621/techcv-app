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
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select";

type SkillOption = {
  id: string;
  name: string;
  icon?: string;
  sortOrder: number;
};

type SkillOptionsResponse = {
  categories: SkillOption[];
  skillMasters: SkillMaster[];
};

type SkillMaster = {
  id: string;
  name: string;
  categoryId: string;
  category: string;
  url: string;
};

type JobEmploymentType = {
  id: string;
  name: string;
  sortOrder: number;
};

type JobCompany = {
  id: string;
  name: string;
  url: string;
};

type ProjectPhase = {
  id: string;
  name: string;
  sortOrder: number;
};

type JobHistoryOptionsResponse = {
  employmentTypes: JobEmploymentType[];
  companies: JobCompany[];
};

type ProjectOptionsResponse = {
  phases: ProjectPhase[];
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
  name: string;
  categoryId: string;
  url: string;
};

type EmploymentTypeForm = {
  id: string;
  name: string;
  sortOrder: number;
};

type CompanyForm = {
  name: string;
  url: string;
};

type ProjectPhaseForm = {
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
  name: "",
  categoryId: "",
  url: "",
};
const emptyEmploymentTypeForm: EmploymentTypeForm = {
  id: "",
  name: "",
  sortOrder: 0,
};
const emptyCompanyForm: CompanyForm = {
  name: "",
  url: "",
};
const emptyProjectPhaseForm: ProjectPhaseForm = {
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
  const [skillMasters, setSkillMasters] = useState<SkillMaster[]>([]);
  const [employmentTypes, setEmploymentTypes] = useState<JobEmploymentType[]>([]);
  const [companies, setCompanies] = useState<JobCompany[]>([]);
  const [projectPhases, setProjectPhases] = useState<ProjectPhase[]>([]);
  const [profileLinkMasters, setProfileLinkMasters] = useState<ProfileLinkMaster[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSavingCategory, setIsSavingCategory] = useState(false);
  const [isSavingSkillMaster, setIsSavingSkillMaster] = useState(false);
  const [isSavingEmploymentType, setIsSavingEmploymentType] = useState(false);
  const [isSavingCompany, setIsSavingCompany] = useState(false);
  const [isSavingProjectPhase, setIsSavingProjectPhase] = useState(false);
  const [isSavingProfileLinkMaster, setIsSavingProfileLinkMaster] = useState(false);
  const [categoryDialogOpen, setCategoryDialogOpen] = useState(false);
  const [skillMasterDialogOpen, setSkillMasterDialogOpen] = useState(false);
  const [employmentTypeDialogOpen, setEmploymentTypeDialogOpen] = useState(false);
  const [companyDialogOpen, setCompanyDialogOpen] = useState(false);
  const [projectPhaseDialogOpen, setProjectPhaseDialogOpen] = useState(false);
  const [profileLinkMasterDialogOpen, setProfileLinkMasterDialogOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState<SkillOption | null>(null);
  const [editingSkillMaster, setEditingSkillMaster] = useState<SkillMaster | null>(null);
  const [editingEmploymentType, setEditingEmploymentType] = useState<JobEmploymentType | null>(null);
  const [editingCompany, setEditingCompany] = useState<JobCompany | null>(null);
  const [editingProjectPhase, setEditingProjectPhase] = useState<ProjectPhase | null>(null);
  const [editingProfileLinkMaster, setEditingProfileLinkMaster] = useState<ProfileLinkMaster | null>(null);
  const [categoryForm, setCategoryForm] = useState<CategoryForm>(emptyCategoryForm);
  const [skillMasterForm, setSkillMasterForm] = useState<SkillMasterForm>(emptySkillMasterForm);
  const [employmentTypeForm, setEmploymentTypeForm] = useState<EmploymentTypeForm>(emptyEmploymentTypeForm);
  const [companyForm, setCompanyForm] = useState<CompanyForm>(emptyCompanyForm);
  const [projectPhaseForm, setProjectPhaseForm] = useState<ProjectPhaseForm>(emptyProjectPhaseForm);
  const [profileLinkMasterForm, setProfileLinkMasterForm] = useState<ProfileLinkMasterForm>(emptyProfileLinkMasterForm);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  async function loadOptions(signal?: AbortSignal) {
    setIsLoading(true);
    setError("");

    try {
      const [
        skillOptionsResponse,
        jobHistoryOptionsResponse,
        profileLinkMastersResponse,
        projectOptionsResponse,
      ] = await Promise.all([
        fetch("/api/skills/options", { signal }),
        fetch("/api/job-histories/options", { signal }),
        fetch("/api/profile/link-masters", { signal }),
        fetch("/api/projects/options", { signal }),
      ]);
      if (
        !skillOptionsResponse.ok ||
        !jobHistoryOptionsResponse.ok ||
        !profileLinkMastersResponse.ok ||
        !projectOptionsResponse.ok
      ) {
        throw new Error("設定情報の取得に失敗しました");
      }

      const skillOptions = (await skillOptionsResponse.json()) as SkillOptionsResponse;
      const jobHistoryOptions = (await jobHistoryOptionsResponse.json()) as JobHistoryOptionsResponse;
      const profileLinkMastersData = (await profileLinkMastersResponse.json()) as ProfileLinkMastersResponse;
      const projectOptions = (await projectOptionsResponse.json()) as ProjectOptionsResponse;
      setCategories(skillOptions.categories ?? []);
      setSkillMasters(skillOptions.skillMasters ?? []);
      setEmploymentTypes(jobHistoryOptions.employmentTypes ?? []);
      setCompanies(jobHistoryOptions.companies ?? []);
      setProjectPhases(projectOptions.phases ?? []);
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

  const openAddCompanyDialog = () => {
    setEditingCompany(null);
    setCompanyForm(emptyCompanyForm);
    setCompanyDialogOpen(true);
    setMessage("");
    setError("");
  };

  const openAddProjectPhaseDialog = () => {
    setEditingProjectPhase(null);
    setProjectPhaseForm(emptyProjectPhaseForm);
    setProjectPhaseDialogOpen(true);
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
      name: skillMaster.name,
      categoryId: skillMaster.categoryId,
      url: skillMaster.url,
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
      sortOrder: employmentType.sortOrder,
    });
    setEmploymentTypeDialogOpen(true);
    setMessage("");
    setError("");
  };

  const openEditCompanyDialog = (company: JobCompany) => {
    setEditingCompany(company);
    setCompanyForm({
      name: company.name,
      url: company.url,
    });
    setCompanyDialogOpen(true);
    setMessage("");
    setError("");
  };

  const openEditProjectPhaseDialog = (phase: ProjectPhase) => {
    setEditingProjectPhase(phase);
    setProjectPhaseForm({
      id: phase.id,
      name: phase.name,
      sortOrder: phase.sortOrder,
    });
    setProjectPhaseDialogOpen(true);
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
    skillMasterForm.categoryId.trim() !== "";

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

  const handleSaveCompany = async () => {
    setIsSavingCompany(true);
    setMessage("");
    setError("");

    try {
      const response = await fetch(
        editingCompany
          ? `/api/job-companies/${encodeURIComponent(editingCompany.id)}`
          : "/api/job-companies",
        {
          method: editingCompany ? "PUT" : "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(companyForm),
        },
      );
      if (!response.ok) {
        throw new Error(editingCompany ? "会社の更新に失敗しました" : "会社の追加に失敗しました");
      }

      setCompanyDialogOpen(false);
      await loadOptions();
      setMessage(editingCompany ? "会社を更新しました" : "会社を追加しました");
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "会社の保存に失敗しました");
    } finally {
      setIsSavingCompany(false);
    }
  };

  const canSaveCompany = companyForm.name.trim() !== "";

  const handleSaveProjectPhase = async () => {
    setIsSavingProjectPhase(true);
    setMessage("");
    setError("");

    try {
      const response = await fetch(
        editingProjectPhase
          ? `/api/project-phases/${encodeURIComponent(editingProjectPhase.id)}`
          : "/api/project-phases",
        {
          method: editingProjectPhase ? "PUT" : "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(projectPhaseForm),
        },
      );
      if (!response.ok) {
        throw new Error(editingProjectPhase ? "担当工程の更新に失敗しました" : "担当工程の追加に失敗しました");
      }

      setProjectPhaseDialogOpen(false);
      await loadOptions();
      setMessage(editingProjectPhase ? "担当工程を更新しました" : "担当工程を追加しました");
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "担当工程の保存に失敗しました");
    } finally {
      setIsSavingProjectPhase(false);
    }
  };

  const canSaveProjectPhase =
    projectPhaseForm.name.trim() !== "";

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

  const scrollToTop = () => {
    const scrollContainer = document.querySelector("main");
    if (scrollContainer) {
      scrollContainer.scrollTo({ top: 0, behavior: "smooth" });
      return;
    }

    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  return (
    <div className="p-8">
      <div className="max-w-6xl mx-auto space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">設定</h1>
          <p className="text-gray-600 mt-1">各画面で利用するマスタ情報</p>
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

        <Tabs defaultValue="profile" className="w-full">
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="profile">プロフィール</TabsTrigger>
            <TabsTrigger value="skills">スキル</TabsTrigger>
            <TabsTrigger value="job-history">職歴</TabsTrigger>
            <TabsTrigger value="projects">案件</TabsTrigger>
          </TabsList>

          <TabsContent value="profile" className="mt-6">
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
          </TabsContent>

          <TabsContent value="skills" className="mt-6 grid grid-cols-1 xl:grid-cols-2 gap-6">
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
                    <TableHead>名称</TableHead>
                    <TableHead>カテゴリ</TableHead>
                    <TableHead>ID</TableHead>
                    <TableHead className="text-right">操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {skillMasters.map((skillMaster) => (
                    <TableRow key={skillMaster.id}>
                      <TableCell className="font-medium">
                        {skillMaster.url ? (
                          <a
                            href={skillMaster.url}
                            target="_blank"
                            rel="noreferrer"
                            className="text-blue-600 hover:underline"
                          >
                            {skillMaster.name}
                          </a>
                        ) : (
                          skillMaster.name
                        )}
                      </TableCell>
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
                      <TableCell colSpan={4} className="text-center text-gray-500">
                        スキルマスタは登録されていません。
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>

          </TabsContent>

          <TabsContent value="job-history" className="mt-6 grid grid-cols-1 xl:grid-cols-2 gap-6">
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

          <Card>
            <CardHeader className="flex flex-row items-center justify-between gap-4">
              <div>
                <CardTitle>会社</CardTitle>
                <CardDescription>job_companies テーブルの内容</CardDescription>
              </div>
              <Button type="button" variant="outline" onClick={openAddCompanyDialog}>
                <LucideIcons.Plus className="w-4 h-4 mr-2" />
                追加
              </Button>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>会社名</TableHead>
                    <TableHead>URL</TableHead>
                    <TableHead>ID</TableHead>
                    <TableHead className="text-right">操作</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {companies.map((company) => (
                    <TableRow key={company.id}>
                      <TableCell className="font-medium">{company.name}</TableCell>
                      <TableCell>
                        {company.url ? (
                          <a
                            href={company.url}
                            target="_blank"
                            rel="noreferrer"
                            className="text-sm text-blue-600 hover:underline"
                          >
                            {company.url}
                          </a>
                        ) : (
                          <span className="text-sm text-gray-400">未設定</span>
                        )}
                      </TableCell>
                      <TableCell className="font-mono text-xs text-gray-600">{company.id}</TableCell>
                      <TableCell className="text-right">
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => openEditCompanyDialog(company)}
                        >
                          <LucideIcons.Pencil className="w-4 h-4 mr-2" />
                          編集
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                  {!isLoading && companies.length === 0 && (
                    <TableRow>
                      <TableCell colSpan={4} className="text-center text-gray-500">
                        会社は登録されていません。
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
          </TabsContent>

          <TabsContent value="projects" className="mt-6">
            <Card>
            <CardHeader className="flex flex-row items-center justify-between gap-4">
              <div>
                <CardTitle>案件の担当工程</CardTitle>
                <CardDescription>project_phases テーブルの内容</CardDescription>
              </div>
              <Button type="button" variant="outline" onClick={openAddProjectPhaseDialog}>
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
                  {projectPhases.map((phase) => (
                    <TableRow key={phase.id}>
                      <TableCell>{phase.sortOrder}</TableCell>
                      <TableCell className="font-medium">{phase.name}</TableCell>
                      <TableCell className="font-mono text-xs text-gray-600">{phase.id}</TableCell>
                      <TableCell className="text-right">
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => openEditProjectPhaseDialog(phase)}
                        >
                          <LucideIcons.Pencil className="w-4 h-4 mr-2" />
                          編集
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                  {!isLoading && projectPhases.length === 0 && (
                    <TableRow>
                      <TableCell colSpan={4} className="text-center text-gray-500">
                        案件の担当工程は登録されていません。
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
          </TabsContent>
        </Tabs>

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
                  ? "スキル名、URL、カテゴリを変更します。"
                  : "スキルマスタの名称、URL、カテゴリを登録します。"}
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
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
                <Label htmlFor="skill-master-url">URL</Label>
                <Input
                  id="skill-master-url"
                  value={skillMasterForm.url}
                  onChange={(e) =>
                    setSkillMasterForm({ ...skillMasterForm, url: e.target.value })
                  }
                  placeholder="https://example.com"
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

        <Dialog open={companyDialogOpen} onOpenChange={setCompanyDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editingCompany ? "会社を編集" : "会社を追加"}</DialogTitle>
              <DialogDescription>
                職歴管理で選択する会社名と会社URLを登録します。
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="company-name">会社名</Label>
                <Input
                  id="company-name"
                  value={companyForm.name}
                  onChange={(e) =>
                    setCompanyForm({ ...companyForm, name: e.target.value })
                  }
                  placeholder="株式会社〇〇"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="company-url">会社URL</Label>
                <Input
                  id="company-url"
                  value={companyForm.url}
                  onChange={(e) =>
                    setCompanyForm({ ...companyForm, url: e.target.value })
                  }
                  placeholder="https://example.com"
                />
              </div>
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setCompanyDialogOpen(false)}
              >
                キャンセル
              </Button>
              <Button
                type="button"
                onClick={handleSaveCompany}
                disabled={!canSaveCompany || isSavingCompany}
              >
                {isSavingCompany ? "保存中" : "保存"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        <Dialog open={projectPhaseDialogOpen} onOpenChange={setProjectPhaseDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editingProjectPhase ? "担当工程を編集" : "担当工程を追加"}</DialogTitle>
              <DialogDescription>
                {editingProjectPhase
                  ? "案件登録・編集フォームで選択する担当工程の名称と表示順を変更します。"
                  : "案件登録・編集フォームで選択する担当工程の名称と表示順を登録します。IDは自動で採番されます。"}
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="project-phase-name">名称</Label>
                <Input
                  id="project-phase-name"
                  value={projectPhaseForm.name}
                  onChange={(e) =>
                    setProjectPhaseForm({ ...projectPhaseForm, name: e.target.value })
                  }
                  placeholder="要件定義"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="project-phase-sort-order">表示順</Label>
                <Input
                  id="project-phase-sort-order"
                  type="number"
                  value={projectPhaseForm.sortOrder}
                  onChange={(e) =>
                    setProjectPhaseForm({ ...projectPhaseForm, sortOrder: Number(e.target.value) })
                  }
                  placeholder="1"
                />
              </div>
            </div>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setProjectPhaseDialogOpen(false)}
              >
                キャンセル
              </Button>
              <Button
                type="button"
                onClick={handleSaveProjectPhase}
                disabled={!canSaveProjectPhase || isSavingProjectPhase}
              >
                {isSavingProjectPhase ? "保存中" : "保存"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
      <Button
        type="button"
        variant="outline"
        size="icon"
        className="fixed bottom-6 right-6 h-11 w-11 rounded-full bg-white shadow-md"
        onClick={scrollToTop}
        aria-label="画面上部へ移動"
      >
        <LucideIcons.ArrowUp className="h-5 w-5" />
      </Button>
    </div>
  );
}
