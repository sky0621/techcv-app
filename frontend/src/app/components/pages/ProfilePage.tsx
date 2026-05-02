import { useEffect, useState } from "react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Label } from "../ui/label";
import { Textarea } from "../ui/textarea";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import { Switch } from "../ui/switch";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select";
import * as LucideIcons from "lucide-react";
import type { LucideIcon } from "lucide-react";
import { Award, Plus, Trash2 } from "lucide-react";

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
  links: ProfileLinkForm[];
};

type ProfileLinkForm = {
  id: string;
  linkMasterId: string;
  key: string;
  name: string;
  icon: string;
  placeholder: string;
  url: string;
  sortOrder: number;
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
  links?: ProfileLinkForm[];
  visibilitySettings?: Partial<VisibilityForm>;
};

type ProfileResponse = {
  profile: ProfilePayload;
};

type ProfileLinkMastersResponse = {
  linkMasters: ProfileLinkMaster[];
};

type ProfileLinkMaster = {
  id: string;
  key: string;
  name: string;
  icon: string;
  placeholder: string;
  sortOrder: number;
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
  links: [],
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
    links: (profile.links ?? []).map(toProfileLinkForm),
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
    links: profile.links,
    visibilitySettings: visibility,
  };
}

function toProfileLinkForm(link: ProfileLinkForm): ProfileLinkForm {
  return {
    id: link.id ?? "",
    linkMasterId: link.linkMasterId ?? "",
    key: link.key ?? "",
    name: link.name ?? "",
    icon: link.icon ?? "link",
    placeholder: link.placeholder ?? "",
    url: link.url ?? "",
    sortOrder: link.sortOrder ?? 0,
  };
}

function normalizeIconName(icon?: string) {
  return (icon ?? "").trim().replace(/^lucide-/i, "");
}

function toLucideExportName(icon?: string) {
  const normalizedIcon = normalizeIconName(icon);
  if (!normalizedIcon) return "Link";

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

  return isLucideIcon(candidate) ? candidate : LucideIcons.Link;
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
  const [linkMasters, setLinkMasters] = useState<ProfileLinkMaster[]>([]);
  const [newLinkMasterId, setNewLinkMasterId] = useState("");
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
        const [profileResponse, linkMastersResponse] = await Promise.all([
          fetch("/api/profile", { signal: controller.signal }),
          fetch("/api/profile/link-masters", { signal: controller.signal }),
        ]);
        if (!profileResponse.ok || !linkMastersResponse.ok) {
          throw new Error("プロフィールの取得に失敗しました");
        }

        const data = (await profileResponse.json()) as ProfileResponse;
        const linkMasterData = (await linkMastersResponse.json()) as ProfileLinkMastersResponse;
        const masters = linkMasterData.linkMasters ?? [];
        setLinkMasters(masters);
        setProfile(toProfileForm(data.profile));
        const registeredMasterIds = new Set((data.profile.links ?? []).map((link) => link.linkMasterId));
        setNewLinkMasterId(masters.find((master) => !registeredMasterIds.has(master.id))?.id ?? "");
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
      const registeredMasterIds = new Set((data.profile.links ?? []).map((link) => link.linkMasterId));
      setNewLinkMasterId(linkMasters.find((master) => !registeredMasterIds.has(master.id))?.id ?? "");
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

  const updateLink = (index: number, url: string) => {
    setProfile((current) => ({
      ...current,
      links: current.links.map((link, linkIndex) =>
        linkIndex === index ? { ...link, url } : link
      ),
    }));
  };

  const registeredLinkMasterIds = new Set(profile.links.map((link) => link.linkMasterId));
  const availableLinkMasters = linkMasters.filter((master) => !registeredLinkMasterIds.has(master.id));
  const selectedLinkMasterId = availableLinkMasters.some((master) => master.id === newLinkMasterId)
    ? newLinkMasterId
    : availableLinkMasters[0]?.id ?? "";

  const addLink = () => {
    const master = linkMasters.find((candidate) => candidate.id === selectedLinkMasterId);
    if (!master) return;

    setProfile((current) => ({
      ...current,
      links: [
        ...current.links,
        {
          id: "",
          linkMasterId: master.id,
          key: master.key,
          name: master.name,
          icon: master.icon,
          placeholder: master.placeholder,
          url: "",
          sortOrder: master.sortOrder,
        },
      ].sort((left, right) => left.sortOrder - right.sortOrder),
    }));
    setNewLinkMasterId(
      availableLinkMasters.find((candidate) => candidate.id !== master.id)?.id ?? "",
    );
  };

  const removeLink = (index: number) => {
    setProfile((current) => ({
      ...current,
      links: current.links.filter((_, linkIndex) => linkIndex !== index),
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
              <CardHeader className="flex flex-row items-center justify-between gap-4">
                <div>
                  <CardTitle>SNS・外部リンク</CardTitle>
                  <CardDescription>登録済みの外部リンクを管理</CardDescription>
                </div>
                <div className="flex min-w-80 items-center gap-2">
                  <Select
                    value={selectedLinkMasterId}
                    onValueChange={setNewLinkMasterId}
                    disabled={availableLinkMasters.length === 0}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="追加するリンク種別" />
                    </SelectTrigger>
                    <SelectContent>
                      {availableLinkMasters.map((master) => (
                        <SelectItem key={master.id} value={master.id}>
                          {master.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <Button
                    type="button"
                    variant="outline"
                    onClick={addLink}
                    disabled={availableLinkMasters.length === 0}
                  >
                    <Plus className="w-4 h-4 mr-2" />
                    追加
                  </Button>
                </div>
              </CardHeader>
              <CardContent className="space-y-4">
                {profile.links.length === 0 ? (
                  <div className="rounded-md border border-dashed p-6 text-center text-sm text-gray-500">
                    SNS・外部リンクはまだ登録されていません。
                  </div>
                ) : (
                  profile.links.map((link, index) => {
                    const Icon = getIcon(link.icon);
                    const inputID = `profile-link-${link.linkMasterId}`;

                    return (
                      <div key={link.linkMasterId} className="space-y-3 rounded-md border p-4">
                        <div className="flex items-center justify-between gap-3">
                          <Label htmlFor={inputID} className="flex items-center gap-2 text-base font-medium">
                            <Icon className="w-4 h-4" />
                            {link.name}
                          </Label>
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            onClick={() => removeLink(index)}
                          >
                            <Trash2 className="w-4 h-4 mr-2" />
                            削除
                          </Button>
                        </div>
                        <Input
                          id={inputID}
                          placeholder={link.placeholder}
                          value={link.url}
                          onChange={(e) => updateLink(index, e.target.value)}
                        />
                      </div>
                    );
                  })
                )}
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
