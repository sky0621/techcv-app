import { useEffect, useState } from "react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Textarea } from "../ui/textarea";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import { Switch } from "../ui/switch";
import { Github, Globe, BookOpen } from "lucide-react";

type ProfileForm = {
  displayName: string;
  email: string;
  location: string;
  bio: string;
  githubUrl: string;
  zennUrl: string;
  qiitaUrl: string;
  websiteUrl: string;
  workStyle: string;
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
  workStyle?: string;
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
  workStyle: "",
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
    workStyle: profile.workStyle ?? "",
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
    workStyle: profile.workStyle,
    visibilitySettings: visibility,
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
          <TabsList className="grid w-full grid-cols-3">
            <TabsTrigger value="basic">基本情報</TabsTrigger>
            <TabsTrigger value="links">SNS・リンク</TabsTrigger>
            <TabsTrigger value="preferences">働き方</TabsTrigger>
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
                <CardTitle>働き方の希望</CardTitle>
                <CardDescription>勤務形態、稼働開始時期など</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
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
