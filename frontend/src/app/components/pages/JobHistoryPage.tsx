import { useState } from "react";
import { Plus, Pencil, Trash2, Briefcase } from "lucide-react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "../ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select";
import { Badge } from "../ui/badge";

type JobHistory = {
  id: number;
  company: string;
  startDate: string;
  endDate: string;
  employmentType: string;
  projectCount: number;
};

const initialJobs: JobHistory[] = [
  { id: 1, company: "株式会社A", startDate: "2023-01", endDate: "現在", employmentType: "正社員", projectCount: 5 },
  { id: 2, company: "株式会社B", startDate: "2021-04", endDate: "2022-12", employmentType: "正社員", projectCount: 4 },
  { id: 3, company: "フリーランス", startDate: "2020-01", endDate: "2021-03", employmentType: "業務委託", projectCount: 3 },
];

export function JobHistoryPage() {
  const [jobs, setJobs] = useState<JobHistory[]>(initialJobs);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingJob, setEditingJob] = useState<JobHistory | null>(null);
  const [formData, setFormData] = useState({
    company: "",
    startDate: "",
    endDate: "",
    employmentType: "正社員",
  });

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

  const handleSave = () => {
    if (editingJob) {
      setJobs(jobs.map(j => j.id === editingJob.id ? { ...j, ...formData } : j));
    } else {
      const newJob: JobHistory = {
        id: Date.now(),
        ...formData,
        projectCount: 0,
      };
      setJobs([newJob, ...jobs]);
    }
    setIsDialogOpen(false);
  };

  const handleDelete = (id: number) => {
    if (confirm("本当に削除しますか？")) {
      setJobs(jobs.filter(j => j.id !== id));
    }
  };

  return (
    <div className="p-8">
      <div className="max-w-6xl mx-auto space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">職歴管理</h1>
            <p className="text-gray-600 mt-1">勤務先の履歴を時系列で管理</p>
          </div>
          <Button onClick={handleAdd}>
            <Plus className="w-4 h-4 mr-2" />
            職歴を追加
          </Button>
        </div>

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

          {jobs.length === 0 && (
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
              <Button onClick={handleSave}>保存</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
}
