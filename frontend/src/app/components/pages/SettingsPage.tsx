import { useEffect, useState } from "react";
import { Code, Database, Cloud, Pencil, Plus, Wrench } from "lucide-react";
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
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table";

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

type CategoryForm = {
  id: string;
  name: string;
  icon: string;
};

const iconComponents: Record<string, any> = {
  cloud: Cloud,
  code: Code,
  database: Database,
  wrench: Wrench,
};

const iconOptions = ["code", "database", "cloud", "wrench"];
const emptyCategoryForm: CategoryForm = {
  id: "",
  name: "",
  icon: "wrench",
};

function getIcon(icon?: string) {
  if (!icon) {
    return Wrench;
  }

  return iconComponents[icon] ?? Wrench;
}

export function SettingsPage() {
  const [categories, setCategories] = useState<SkillOption[]>([]);
  const [proficiencyLevels, setProficiencyLevels] = useState<SkillOption[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSavingCategory, setIsSavingCategory] = useState(false);
  const [categoryDialogOpen, setCategoryDialogOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState<SkillOption | null>(null);
  const [categoryForm, setCategoryForm] = useState<CategoryForm>(emptyCategoryForm);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  async function loadOptions(signal?: AbortSignal) {
    setIsLoading(true);
    setError("");

    try {
      const response = await fetch("/api/skills/options", {
        signal,
      });
      if (!response.ok) {
        throw new Error("設定情報の取得に失敗しました");
      }

      const data = (await response.json()) as SkillOptionsResponse;
      setCategories(data.categories ?? []);
      setProficiencyLevels(data.proficiencyLevels ?? []);
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
                <Plus className="w-4 h-4 mr-2" />
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
                            <span>{category.icon || "wrench"}</span>
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
                            <Pencil className="w-4 h-4 mr-2" />
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
                <Select
                  value={categoryForm.icon}
                  onValueChange={(value) =>
                    setCategoryForm({ ...categoryForm, icon: value })
                  }
                >
                  <SelectTrigger id="category-icon">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {iconOptions.map((icon) => {
                      const Icon = getIcon(icon);

                      return (
                        <SelectItem key={icon} value={icon}>
                          <span className="flex items-center gap-2">
                            <Icon className="w-4 h-4" />
                            {icon}
                          </span>
                        </SelectItem>
                      );
                    })}
                  </SelectContent>
                </Select>
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
      </div>
    </div>
  );
}
