import { useEffect, useState } from "react";
import { Plus, Pencil, Trash2, Briefcase } from "lucide-react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select";
import { Badge } from "../ui/badge";

type JobHistory = {
  id: string;
  company: string;
  displayName: string;
  startYear: number;
  startMonth: number;
  endYear: number | null;
  endMonth: number | null;
  employmentTypeId: string;
  employmentType: string;
  projectCount: number;
};

type JobEmploymentType = {
  id: string;
  name: string;
  sortOrder: number;
};

type JobHistoriesResponse = {
  jobHistories: JobHistory[];
};

type JobHistoryOptionsResponse = {
  employmentTypes: JobEmploymentType[];
};

type JobHistoryResponse = {
  jobHistory: JobHistory;
};

function sortJobHistoriesByStartDateDesc(values: JobHistory[]) {
  return [...values].sort((left, right) => {
    if (right.startYear !== left.startYear) {
      return right.startYear - left.startYear;
    }
    if (right.startMonth !== left.startMonth) {
      return right.startMonth - left.startMonth;
    }

    return Number(right.id) - Number(left.id);
  });
}

function formatYearMonth(year: number, month: number) {
  return `${year}-${String(month).padStart(2, "0")}`;
}

function formatEndYearMonth(year: number | null, month: number | null) {
  if (year === null || month === null) {
    return "現在";
  }

  return formatYearMonth(year, month);
}

function formatDurationMonths(totalMonths: number) {
  const years = Math.floor(totalMonths / 12);
  const months = totalMonths % 12;

  if (years === 0) {
    return `${months}ヶ月`;
  }
  if (months === 0) {
    return `${years}年`;
  }

  return `${years}年${months}ヶ月`;
}

function formatJobDuration(job: JobHistory) {
  const now = new Date();
  const endYear = job.endYear ?? now.getFullYear();
  const endMonth = job.endMonth ?? now.getMonth() + 1;
  const totalMonths = Math.max(
    1,
    (endYear - job.startYear) * 12 + endMonth - job.startMonth + 1,
  );

  return formatDurationMonths(totalMonths);
}

function parseYearMonth(value: string) {
  const [year, month] = value.split("-").map(Number);

  return {
    year: Number.isFinite(year) ? year : 0,
    month: Number.isFinite(month) ? month : 0,
  };
}

export function JobHistoryPage() {
  const [jobs, setJobs] = useState<JobHistory[]>([]);
  const [employmentTypes, setEmploymentTypes] = useState<JobEmploymentType[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingJob, setEditingJob] = useState<JobHistory | null>(null);
  const [formData, setFormData] = useState({
    company: "",
    displayName: "",
    startDate: "",
    endDate: "",
    employmentTypeId: "",
  });

  useEffect(() => {
    const controller = new AbortController();

    async function loadJobHistories() {
      setIsLoading(true);
      setError("");

      try {
        const [jobHistoriesResponse, optionsResponse] = await Promise.all([
          fetch("/api/job-histories", { signal: controller.signal }),
          fetch("/api/job-histories/options", { signal: controller.signal }),
        ]);
        if (!jobHistoriesResponse.ok) {
          throw new Error("職歴の取得に失敗しました");
        }
        if (!optionsResponse.ok) {
          throw new Error("職歴の選択肢の取得に失敗しました");
        }

        const jobHistoriesData = (await jobHistoriesResponse.json()) as JobHistoriesResponse;
        const optionsData = (await optionsResponse.json()) as JobHistoryOptionsResponse;
        setJobs(sortJobHistoriesByStartDateDesc(jobHistoriesData.jobHistories ?? []));
        setEmploymentTypes(optionsData.employmentTypes ?? []);
      } catch (caught) {
        if (caught instanceof DOMException && caught.name === "AbortError") {
          return;
        }
        setError(caught instanceof Error ? caught.message : "職歴の取得に失敗しました");
      } finally {
        if (!controller.signal.aborted) {
          setIsLoading(false);
        }
      }
    }

    void loadJobHistories();

    return () => {
      controller.abort();
    };
  }, []);

  const handleAdd = () => {
    setEditingJob(null);
    setFormData({
      company: "",
      displayName: "",
      startDate: "",
      endDate: "",
      employmentTypeId: employmentTypes[0]?.id ?? "",
    });
    setIsDialogOpen(true);
  };

  const handleEdit = (job: JobHistory) => {
    setEditingJob(job);
    setFormData({
      company: job.company,
      displayName: job.displayName,
      startDate: formatYearMonth(job.startYear, job.startMonth),
      endDate:
        job.endYear === null || job.endMonth === null
          ? ""
          : formatYearMonth(job.endYear, job.endMonth),
      employmentTypeId: job.employmentTypeId,
    });
    setIsDialogOpen(true);
  };

  const handleSave = async () => {
    setIsSaving(true);
    setError("");

    try {
      const start = parseYearMonth(formData.startDate);
      const end = parseYearMonth(formData.endDate);
      const requestBody = {
        company: formData.company,
        displayName: formData.displayName,
        startYear: start.year,
        startMonth: start.month,
        endYear: formData.endDate.trim() === "" ? null : end.year,
        endMonth: formData.endDate.trim() === "" ? null : end.month,
        employmentTypeId: formData.employmentTypeId,
      };
      const response = await fetch(
        editingJob
          ? `/api/job-histories/${encodeURIComponent(editingJob.id)}`
          : "/api/job-histories",
        {
          method: editingJob ? "PUT" : "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(requestBody),
        },
      );
      if (!response.ok) {
        throw new Error(editingJob ? "職歴の更新に失敗しました" : "職歴の追加に失敗しました");
      }

      const data = (await response.json()) as JobHistoryResponse;
      if (editingJob) {
        setJobs(sortJobHistoriesByStartDateDesc(
          jobs.map(j => j.id === editingJob.id ? data.jobHistory : j),
        ));
      } else {
        setJobs(sortJobHistoriesByStartDateDesc([data.jobHistory, ...jobs]));
      }
      setIsDialogOpen(false);
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "職歴の保存に失敗しました");
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (confirm("本当に削除しますか？")) {
      setError("");

      try {
        const response = await fetch(`/api/job-histories/${encodeURIComponent(id)}`, {
          method: "DELETE",
        });
        if (!response.ok) {
          throw new Error("職歴の削除に失敗しました");
        }

        setJobs(jobs.filter(j => j.id !== id));
      } catch (caught) {
        setError(caught instanceof Error ? caught.message : "職歴の削除に失敗しました");
      }
    }
  };

  const canSave =
    (formData.company.trim() !== "" || formData.displayName.trim() !== "") &&
    formData.startDate.trim() !== "" &&
    formData.employmentTypeId.trim() !== "";

  return (
    <div className="p-8">
      <div className="max-w-6xl mx-auto space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">職歴管理</h1>
            <p className="text-gray-600 mt-1">勤務先の履歴を時系列で管理</p>
          </div>
          <Button onClick={handleAdd} disabled={isLoading || employmentTypes.length === 0}>
            <Plus className="w-4 h-4 mr-2" />
            職歴を追加
          </Button>
        </div>

        {isLoading && (
          <div className="rounded-md border border-gray-200 bg-gray-50 px-4 py-3 text-sm text-gray-600">
            職歴を読み込み中です。
          </div>
        )}
        {error && (
          <div className="rounded-md border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {error}
          </div>
        )}

        <div className="space-y-4">
          {jobs.map((job) => (
            <Card key={job.id}>
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div className="flex items-start gap-4">
                    <div className="p-3 bg-blue-100 rounded-lg">
                      <Briefcase className="w-6 h-6 text-blue-600" />
                    </div>
                    <div>
                      <CardTitle className="text-xl">{job.displayName || job.company}</CardTitle>
                      {job.company && job.displayName && (
                        <div className="mt-1 text-sm text-gray-500">{job.company}</div>
                      )}
                      <div className="flex items-center gap-3 mt-2 text-sm text-gray-600">
                        <span>
                          {formatYearMonth(job.startYear, job.startMonth)} 〜{" "}
                          {formatEndYearMonth(job.endYear, job.endMonth)}
                        </span>
                        <span>期間: {formatJobDuration(job)}</span>
                        <Badge variant="outline">{job.employmentType}</Badge>
                        <span>{job.projectCount}件の案件</span>
                      </div>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleEdit(job)}
                    >
                      <Pencil className="w-4 h-4" />
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleDelete(job.id)}
                    >
                      <Trash2 className="w-4 h-4 text-red-600" />
                    </Button>
                  </div>
                </div>
              </CardHeader>
            </Card>
          ))}

          {!isLoading && jobs.length === 0 && (
            <Card>
              <CardContent className="flex flex-col items-center justify-center py-12">
                <Briefcase className="w-12 h-12 text-gray-400 mb-4" />
                <p className="text-gray-600 mb-4">職歴がまだ登録されていません</p>
                <Button onClick={handleAdd}>
                  <Plus className="w-4 h-4 mr-2" />
                  最初の職歴を追加
                </Button>
              </CardContent>
            </Card>
          )}
        </div>

        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>
                {editingJob ? "職歴を編集" : "職歴を追加"}
              </DialogTitle>
              <DialogDescription>
                会社名または表示名、所属期間、雇用形態を入力してください
              </DialogDescription>
            </DialogHeader>

            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="company">会社名</Label>
                <Input
                  id="company"
                  value={formData.company}
                  onChange={(e) => setFormData({ ...formData, company: e.target.value })}
                  placeholder="株式会社〇〇"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="displayName">表示名</Label>
                <Input
                  id="displayName"
                  value={formData.displayName}
                  onChange={(e) => setFormData({ ...formData, displayName: e.target.value })}
                  placeholder="現職、屋号、表示用の勤務先名など"
                />
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
                    type="month"
                    value={formData.endDate}
                    onChange={(e) => setFormData({ ...formData, endDate: e.target.value })}
                    placeholder="空欄なら現在"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="employmentType">雇用形態</Label>
                <Select
                  value={formData.employmentTypeId}
                  onValueChange={(value) => setFormData({ ...formData, employmentTypeId: value })}
                  disabled={employmentTypes.length === 0}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {employmentTypes.map((employmentType) => (
                      <SelectItem key={employmentType.id} value={employmentType.id}>
                        {employmentType.name}
                      </SelectItem>
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
