import { useEffect, useState } from "react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Textarea } from "../ui/textarea";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import { Switch } from "../ui/switch";
import { Award, BookOpen, Github, Globe, Plus, Trash2 } from "lucide-react";

type QualificationForm = {
  id: string;
  name: string;
  acquiredDate: string;
  organization: string;
  url: string;
  memo: string;
};

type ProfileForm = {
  displayName: string;
  email: string;
  location: string;
  bio: string;
  githubUrl: string;
  zennUrl: string;
  qiitaUrl: string;
  websiteUrl: string;
  occupation: string;
  employmentType: string;
  workStyle: string;
  qualifications: QualificationForm[];
};

type VisibilityForm = {
  location: boolean;
  email: boolean;
};

type ProfilePayload = {
  displayName?: string;
  location?: string;
  email?: string;
  bio?: string;
  githubUrl?: string;
  zennUrl?: string;
  qiitaUrl?: string;
  websiteUrl?: string;
  occupation?: string;
  employmentType?: string;
  workStyle?: string;
  qualifications?: QualificationForm[];
  visibilitySettings?: Partial<VisibilityForm>;
};

type ProfileResponse = {
  profile: ProfilePayload;
};

const initialProfile: ProfileForm = {
  displayName: "",
  email: "",
  location: "",
  bio: "",
  githubUrl: "",
  zennUrl: "",
  qiitaUrl: "",
  websiteUrl: "",
  occupation: "",
  employmentType: "",
  workStyle: "",
  qualifications: [],
};

const initialVisibility: VisibilityForm = {
  location: true,
  email: true,
};

function toProfileForm(profile: ProfilePayload): ProfileForm {
  return {
    displayName: profile.displayName ?? "",
    email: profile.email ?? "",
    location: profile.location ?? "",
    bio: profile.bio ?? "",
    githubUrl: profile.githubUrl ?? "",
    zennUrl: profile.zennUrl ?? "",
    qiitaUrl: profile.qiitaUrl ?? "",
    websiteUrl: profile.websiteUrl ?? "",
    occupation: profile.occupation ?? "",
    employmentType: profile.employmentType ?? "",
    workStyle: profile.workStyle ?? "",
    qualifications: (profile.qualifications ?? []).map(toQualificationForm),
  };
}

function toProfilePayload(profile: ProfileForm, visibility: VisibilityForm): ProfilePayload {
  return {
    displayName: profile.displayName,
    email: profile.email,
    location: profile.location,
    bio: profile.bio,
    githubUrl: profile.githubUrl,
    zennUrl: profile.zennUrl,
    qiitaUrl: profile.qiitaUrl,
    websiteUrl: profile.websiteUrl,
    occupation: profile.occupation,
    employmentType: profile.employmentType,
    workStyle: profile.workStyle,
    qualifications: profile.qualifications,
    visibilitySettings: visibility,
  };
}

function toQualificationForm(qualification: QualificationForm): QualificationForm {
  return {
    id: qualification.id ?? "",
    name: qualification.name ?? "",
    acquiredDate: qualification.acquiredDate ?? "",
    organization: qualification.organization ?? "",
    url: qualification.url ?? "",
    memo: qualification.memo ?? "",
  };
}

function emptyQualification(): QualificationForm {
  return {
    id: "",
    name: "",
    acquiredDate: "",
    organization: "",
    url: "",
    memo: "",
  };
}

export function ProfilePage() {
  const [profile, setProfile] = useState<ProfileForm>(initialProfile);
  const [visibility, setVisibility] = useState<VisibilityForm>(initialVisibility);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    const controller = new AbortController();

    async function loadProfile() {
      setIsLoading(true);
      setError("");

      try {
        const response = await fetch("/api/profile", {
          signal: controller.signal,
        });
        if (!response.ok) {
          throw new Error("プロフィールの取得に失敗しました");
        }

        const data = (await response.json()) as ProfileResponse;
        setProfile(toProfileForm(data.profile));
        setVisibility({
          ...initialVisibility,
          ...data.profile.visibilitySettings,
        });
      } catch (caught) {
        if (caught instanceof DOMException && caught.name === "AbortError") {
          return;
        }
        setError(caught instanceof Error ? caught.message : "プロフィールの取得に失敗しました");
      } finally {
        if (!controller.signal.aborted) {
          setIsLoading(false);
        }
      }
    }

    void loadProfile();

    return () => {
      controller.abort();
    };
  }, []);

  const handleSave = async () => {
    setIsSaving(true);
    setMessage("");
    setError("");

    try {
      const response = await fetch("/api/profile", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(toProfilePayload(profile, visibility)),
      });
      if (!response.ok) {
        throw new Error("プロフィールの保存に失敗しました");
      }

      const data = (await response.json()) as ProfileResponse;
      setProfile(toProfileForm(data.profile));
      setVisibility({
        ...initialVisibility,
        ...data.profile.visibilitySettings,
      });
      setMessage("保存しました");
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "プロフィールの保存に失敗しました");
    } finally {
      setIsSaving(false);
    }
  };

  const addQualification = () => {
    setProfile((current) => ({
      ...current,
      qualifications: [...current.qualifications, emptyQualification()],
    }));
  };

  const updateQualification = (
    index: number,
    field: keyof QualificationForm,
    value: string,
  ) => {
    setProfile((current) => ({
      ...current,
      qualifications: current.qualifications.map((qualification, qualificationIndex) =>
        qualificationIndex === index
          ? { ...qualification, [field]: value }
          : qualification
      ),
    }));
  };

  const removeQualification = (index: number) => {
    setProfile((current) => ({
      ...current,
      qualifications: current.qualifications.filter((_, qualificationIndex) =>
        qualificationIndex !== index
      ),
    }));
  };

  return (
    <div className="p-8">
      <div className="max-w-4xl mx-auto space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">プロフィール管理</h1>
          <p className="text-gray-600 mt-1">基本情報とSNSリンクを管理</p>
        </div>

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

        <Tabs defaultValue="basic" className="w-full">
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="basic">基本情報</TabsTrigger>
            <TabsTrigger value="links">SNS・リンク</TabsTrigger>
            <TabsTrigger value="preferences">働き方</TabsTrigger>
            <TabsTrigger value="qualifications">資格情報</TabsTrigger>
          </TabsList>

          <TabsContent value="basic" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>基本情報</CardTitle>
                <CardDescription>氏名、連絡先などの基本情報</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="displayName">表示名</Label>
                    <Input
                      id="displayName"
                      value={profile.displayName}
                      onChange={(e) =>
                        setProfile({ ...profile, displayName: e.target.value })
                      }
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="email">メールアドレス</Label>
                    <Input
                      id="email"
                      type="email"
                      value={profile.email}
                      onChange={(e) =>
                        setProfile({ ...profile, email: e.target.value })
                      }
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="location">居住地</Label>
                  <Input
                    id="location"
                    value={profile.location}
                    onChange={(e) =>
                      setProfile({ ...profile, location: e.target.value })
                    }
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="bio">自己紹介</Label>
                  <Textarea
                    id="bio"
                    rows={4}
                    value={profile.bio}
                    onChange={(e) =>
                      setProfile({ ...profile, bio: e.target.value })
                    }
                  />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>公開設定</CardTitle>
                <CardDescription>経歴書に表示する項目を選択</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <Label htmlFor="location-visibility">居住地を表示</Label>
                  <Switch
                    id="location-visibility"
                    checked={visibility.location}
                    onCheckedChange={(checked) =>
                      setVisibility({ ...visibility, location: checked })
                    }
                  />
                </div>
                <div className="flex items-center justify-between">
                  <Label htmlFor="email-visibility">メールアドレスを表示</Label>
                  <Switch
                    id="email-visibility"
                    checked={visibility.email}
                    onCheckedChange={(checked) =>
                      setVisibility({ ...visibility, email: checked })
                    }
                  />
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="links" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>SNS・外部リンク</CardTitle>
                <CardDescription>GitHub、Zenn、Qiitaなどのリンク</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="github" className="flex items-center gap-2">
                    <Github className="w-4 h-4" />
                    GitHub
                  </Label>
                  <Input
                    id="github"
                    placeholder="https://github.com/username"
                    value={profile.githubUrl}
                    onChange={(e) =>
                      setProfile({ ...profile, githubUrl: e.target.value })
                    }
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="zenn" className="flex items-center gap-2">
                    <BookOpen className="w-4 h-4" />
                    Zenn
                  </Label>
                  <Input
                    id="zenn"
                    placeholder="https://zenn.dev/username"
                    value={profile.zennUrl}
                    onChange={(e) =>
                      setProfile({ ...profile, zennUrl: e.target.value })
                    }
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="qiita" className="flex items-center gap-2">
                    <BookOpen className="w-4 h-4" />
                    Qiita
                  </Label>
                  <Input
                    id="qiita"
                    placeholder="https://qiita.com/username"
                    value={profile.qiitaUrl}
                    onChange={(e) =>
                      setProfile({ ...profile, qiitaUrl: e.target.value })
                    }
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="website" className="flex items-center gap-2">
                    <Globe className="w-4 h-4" />
                    個人サイト
                  </Label>
                  <Input
                    id="website"
                    placeholder="https://example.com"
                    value={profile.websiteUrl}
                    onChange={(e) =>
                      setProfile({ ...profile, websiteUrl: e.target.value })
                    }
                  />
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="preferences" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>職業・働き方</CardTitle>
                <CardDescription>職業、労働形態、働き方の希望</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="occupation">職業</Label>
                    <Input
                      id="occupation"
                      placeholder="例: ソフトウェアエンジニア"
                      value={profile.occupation}
                      onChange={(e) =>
                        setProfile({ ...profile, occupation: e.target.value })
                      }
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="employmentType">労働形態</Label>
                    <Input
                      id="employmentType"
                      placeholder="例: フリーランス、正社員、業務委託"
                      value={profile.employmentType}
                      onChange={(e) =>
                        setProfile({ ...profile, employmentType: e.target.value })
                      }
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="workStyle">希望する働き方</Label>
                  <Input
                    id="workStyle"
                    placeholder="例: リモートワーク希望、週3日出社可能"
                    value={profile.workStyle}
                    onChange={(e) =>
                      setProfile({ ...profile, workStyle: e.target.value })
                    }
                  />
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="qualifications" className="space-y-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between gap-4">
                <div>
                  <CardTitle>資格情報</CardTitle>
                  <CardDescription>資格名、取得日、管理団体、説明URL、メモを管理</CardDescription>
                </div>
                <Button type="button" variant="outline" onClick={addQualification}>
                  <Plus className="w-4 h-4 mr-2" />
                  追加
                </Button>
              </CardHeader>
              <CardContent className="space-y-4">
                {profile.qualifications.length === 0 ? (
                  <div className="rounded-md border border-dashed p-6 text-center text-sm text-gray-500">
                    資格情報はまだ登録されていません。
                  </div>
                ) : (
                  profile.qualifications.map((qualification, index) => (
                    <div
                      key={qualification.id || `qualification-${index}`}
                      className="space-y-4 rounded-md border p-4"
                    >
                      <div className="flex items-center justify-between gap-3">
                        <div className="flex items-center gap-2 font-medium text-gray-900">
                          <Award className="w-4 h-4" />
                          資格 {index + 1}
                        </div>
                        <Button
                          type="button"
                          variant="outline"
                          onClick={() => removeQualification(index)}
                        >
                          <Trash2 className="w-4 h-4 mr-2" />
                          削除
                        </Button>
                      </div>

                      <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                          <Label htmlFor={`qualification-name-${index}`}>資格名</Label>
                          <Input
                            id={`qualification-name-${index}`}
                            value={qualification.name}
                            onChange={(e) =>
                              updateQualification(index, "name", e.target.value)
                            }
                          />
                        </div>
                        <div className="space-y-2">
                          <Label htmlFor={`qualification-date-${index}`}>取得日</Label>
                          <Input
                            id={`qualification-date-${index}`}
                            type="date"
                            value={qualification.acquiredDate}
                            onChange={(e) =>
                              updateQualification(index, "acquiredDate", e.target.value)
                            }
                          />
                        </div>
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor={`qualification-organization-${index}`}>
                          資格を管理する団体名
                        </Label>
                        <Input
                          id={`qualification-organization-${index}`}
                          value={qualification.organization}
                          onChange={(e) =>
                            updateQualification(index, "organization", e.target.value)
                          }
                        />
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor={`qualification-url-${index}`}>説明URL</Label>
                        <Input
                          id={`qualification-url-${index}`}
                          type="url"
                          placeholder="https://example.com/certification"
                          value={qualification.url}
                          onChange={(e) =>
                            updateQualification(index, "url", e.target.value)
                          }
                        />
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor={`qualification-memo-${index}`}>メモ</Label>
                        <Textarea
                          id={`qualification-memo-${index}`}
                          rows={3}
                          value={qualification.memo}
                          onChange={(e) =>
                            updateQualification(index, "memo", e.target.value)
                          }
                        />
                      </div>
                    </div>
                  ))
                )}
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        <div className="flex justify-end gap-3">
          <Button variant="outline">キャンセル</Button>
          <Button onClick={handleSave} disabled={isLoading || isSaving}>
            {isSaving ? "保存中" : "保存"}
          </Button>
        </div>
      </div>
    </div>
  );
}
