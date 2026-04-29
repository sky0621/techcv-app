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
  startDate: string;
  endDate: string;
  employmentType: string;
  projectCount: number;
  sortOrder: number;
};

type JobHistoriesResponse = {
  jobHistories: JobHistory[];
};

type JobHistoryResponse = {
  jobHistory: JobHistory;
};

export function JobHistoryPage() {
  const [jobs, setJobs] = useState<JobHistory[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingJob, setEditingJob] = useState<JobHistory | null>(null);
  const [formData, setFormData] = useState({
    company: "",
    startDate: "",
    endDate: "",
    employmentType: "正社員",
  });

  useEffect(() => {
    const controller = new AbortController();

    async function loadJobHistories() {
      setIsLoading(true);
      setError("");

      try {
        const response = await fetch("/api/job-histories", {
          signal: controller.signal,
        });
        if (!response.ok) {
          throw new Error("職歴の取得に失敗しました");
        }

        const data = (await response.json()) as JobHistoriesResponse;
        setJobs(data.jobHistories ?? []);
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
    setFormData({ company: "", startDate: "", endDate: "", employmentType: "正社員" });
    setIsDialogOpen(true);
  };

  const handleEdit = (job: JobHistory) => {
    setEditingJob(job);
    setFormData({
      company: job.company,
      startDate: job.startDate,
      endDate: job.endDate,
      employmentType: job.employmentType,
    });
    setIsDialogOpen(true);
  };

  const handleSave = async () => {
    setIsSaving(true);
    setError("");

    try {
      const response = await fetch(
        editingJob
          ? `/api/job-histories/${encodeURIComponent(editingJob.id)}`
          : "/api/job-histories",
        {
          method: editingJob ? "PUT" : "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(formData),
        },
      );
      if (!response.ok) {
        throw new Error(editingJob ? "職歴の更新に失敗しました" : "職歴の追加に失敗しました");
      }

      const data = (await response.json()) as JobHistoryResponse;
      if (editingJob) {
        setJobs(jobs.map(j => j.id === editingJob.id ? data.jobHistory : j));
      } else {
        setJobs([data.jobHistory, ...jobs]);
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
    formData.company.trim() !== "" &&
    formData.startDate.trim() !== "" &&
    formData.employmentType.trim() !== "";

  return (
    <div className="p-8">
      <div className="max-w-6xl mx-auto space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">職歴管理</h1>
            <p className="text-gray-600 mt-1">勤務先の履歴を時系列で管理</p>
          </div>
          <Button onClick={handleAdd} disabled={isLoading}>
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
                      <CardTitle className="text-xl">{job.company}</CardTitle>
                      <div className="flex items-center gap-3 mt-2 text-sm text-gray-600">
                        <span>{job.startDate} 〜 {job.endDate}</span>
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
                会社名、所属期間、雇用形態を入力してください
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
                <Label htmlFor="employmentType">雇用形態</Label>
                <Select
                  value={formData.employmentType}
                  onValueChange={(value) => setFormData({ ...formData, employmentType: value })}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="正社員">正社員</SelectItem>
                    <SelectItem value="契約社員">契約社員</SelectItem>
                    <SelectItem value="業務委託">業務委託</SelectItem>
                    <SelectItem value="派遣">派遣</SelectItem>
                    <SelectItem value="アルバイト">アルバイト</SelectItem>
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
