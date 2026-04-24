import { useState } from "react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Textarea } from "../ui/textarea";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import { Switch } from "../ui/switch";
import { Github, Globe, BookOpen } from "lucide-react";

export function ProfilePage() {
  const [profile, setProfile] = useState({
    displayName: "山田 太郎",
    email: "yamada@example.com",
    location: "東京都",
    phone: "090-1234-5678",
    bio: "Webエンジニアとして5年以上の経験があります。",
    githubUrl: "https://github.com/username",
    zennUrl: "https://zenn.dev/username",
    qiitaUrl: "https://qiita.com/username",
    websiteUrl: "https://example.com",
    workStyle: "リモートワーク希望",
    availability: "即日可能",
  });

  const [visibility, setVisibility] = useState({
    phone: true,
    location: true,
    email: true,
  });

  const handleSave = () => {
    // TODO: Supabaseへ保存
    alert("保存しました");
  };

  return (
    <div className="p-8">
      <div className="max-w-4xl mx-auto space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">プロフィール管理</h1>
          <p className="text-gray-600 mt-1">基本情報とSNSリンクを管理</p>
        </div>

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

                <div className="grid grid-cols-2 gap-4">
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
                    <Label htmlFor="phone">電話番号</Label>
                    <Input
                      id="phone"
                      value={profile.phone}
                      onChange={(e) =>
                        setProfile({ ...profile, phone: e.target.value })
                      }
                    />
                  </div>
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
                  <Label htmlFor="phone-visibility">電話番号を表示</Label>
                  <Switch
                    id="phone-visibility"
                    checked={visibility.phone}
                    onCheckedChange={(checked) =>
                      setVisibility({ ...visibility, phone: checked })
                    }
                  />
                </div>
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

                <div className="space-y-2">
                  <Label htmlFor="availability">稼働開始時期</Label>
                  <Input
                    id="availability"
                    placeholder="例: 即日可能、2024年6月から"
                    value={profile.availability}
                    onChange={(e) =>
                      setProfile({ ...profile, availability: e.target.value })
                    }
                  />
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        <div className="flex justify-end gap-3">
          <Button variant="outline">キャンセル</Button>
          <Button onClick={handleSave}>保存</Button>
        </div>
      </div>
    </div>
  );
}
