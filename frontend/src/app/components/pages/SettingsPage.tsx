import { useEffect, useState } from "react";
import { Code, Database, Cloud, Wrench } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
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

const iconComponents: Record<string, any> = {
  cloud: Cloud,
  code: Code,
  database: Database,
  wrench: Wrench,
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
  const [error, setError] = useState("");

  useEffect(() => {
    const controller = new AbortController();

    async function loadOptions() {
      setIsLoading(true);
      setError("");

      try {
        const response = await fetch("/api/skills/options", {
          signal: controller.signal,
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
        if (!controller.signal.aborted) {
          setIsLoading(false);
        }
      }
    }

    void loadOptions();

    return () => {
      controller.abort();
    };
  }, []);

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

        <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
          <Card>
            <CardHeader>
              <CardTitle>スキルカテゴリ</CardTitle>
              <CardDescription>skill_categories テーブルの内容</CardDescription>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>表示順</TableHead>
                    <TableHead>名称</TableHead>
                    <TableHead>アイコン</TableHead>
                    <TableHead>ID</TableHead>
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
                      </TableRow>
                    );
                  })}
                  {!isLoading && categories.length === 0 && (
                    <TableRow>
                      <TableCell colSpan={4} className="text-center text-gray-500">
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
      </div>
    </div>
  );
}
