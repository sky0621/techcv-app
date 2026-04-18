"use client";

import { FormEvent, startTransition, useEffect, useState } from "react";

type Profile = {
  id: string;
  userId: string;
  fullName: string;
  nickname: string;
  location: string;
  email: string;
  phone: string;
  summary: string;
  githubUrl: string;
  zennUrl: string;
  qiitaUrl: string;
  websiteUrl: string;
  preferredWorkStyle: string;
  visibilitySettings: {
    email?: boolean;
    phone?: boolean;
  };
  createdAt: string;
  updatedAt: string;
};

type ProfileResponse = {
  profile: Profile;
};

type ProfileForm = {
  fullName: string;
  nickname: string;
  location: string;
  email: string;
  phone: string;
  summary: string;
  githubUrl: string;
  zennUrl: string;
  qiitaUrl: string;
  websiteUrl: string;
  preferredWorkStyle: string;
  showEmail: boolean;
  showPhone: boolean;
};

const emptyForm: ProfileForm = {
  fullName: "",
  nickname: "",
  location: "",
  email: "",
  phone: "",
  summary: "",
  githubUrl: "",
  zennUrl: "",
  qiitaUrl: "",
  websiteUrl: "",
  preferredWorkStyle: "",
  showEmail: false,
  showPhone: false
};

const fieldClassName =
  "mt-2 w-full rounded-2xl border border-slate-300 bg-white px-4 py-3 text-sm text-slate-900 outline-none transition focus:border-sky-500 focus:ring-4 focus:ring-sky-100";

function toForm(profile: Profile): ProfileForm {
  return {
    fullName: profile.fullName,
    nickname: profile.nickname,
    location: profile.location,
    email: profile.email,
    phone: profile.phone,
    summary: profile.summary,
    githubUrl: profile.githubUrl,
    zennUrl: profile.zennUrl,
    qiitaUrl: profile.qiitaUrl,
    websiteUrl: profile.websiteUrl,
    preferredWorkStyle: profile.preferredWorkStyle,
    showEmail: Boolean(profile.visibilitySettings?.email),
    showPhone: Boolean(profile.visibilitySettings?.phone)
  };
}

function formatDateTime(value: string) {
  if (!value) {
    return "-";
  }

  return new Intl.DateTimeFormat("ja-JP", {
    dateStyle: "medium",
    timeStyle: "short"
  }).format(new Date(value));
}

async function loadProfile() {
  const response = await fetch("/api/profile", {
    cache: "no-store"
  });

  if (!response.ok) {
    throw new Error("failed to load profile");
  }

  return (await response.json()) as ProfileResponse;
}

async function saveProfile(form: ProfileForm) {
  const response = await fetch("/api/profile", {
    method: "PUT",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      fullName: form.fullName,
      nickname: form.nickname,
      location: form.location,
      email: form.email,
      phone: form.phone,
      summary: form.summary,
      githubUrl: form.githubUrl,
      zennUrl: form.zennUrl,
      qiitaUrl: form.qiitaUrl,
      websiteUrl: form.websiteUrl,
      preferredWorkStyle: form.preferredWorkStyle,
      visibilitySettings: {
        email: form.showEmail,
        phone: form.showPhone
      }
    })
  });

  if (!response.ok) {
    throw new Error("failed to save profile");
  }

  return (await response.json()) as ProfileResponse;
}

export default function HomePage() {
  const [profile, setProfile] = useState<Profile | null>(null);
  const [form, setForm] = useState<ProfileForm>(emptyForm);
  const [loadError, setLoadError] = useState("");
  const [saveMessage, setSaveMessage] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    let cancelled = false;

    loadProfile()
      .then((data) => {
        if (cancelled) {
          return;
        }

        startTransition(() => {
          setProfile(data.profile);
          setForm(toForm(data.profile));
          setLoadError("");
        });
      })
      .catch(() => {
        if (!cancelled) {
          setLoadError(
            "プロフィールの取得に失敗しました。バックエンドと MySQL の起動状態を確認してください。"
          );
        }
      })
      .finally(() => {
        if (!cancelled) {
          setIsLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, []);

  function updateField<K extends keyof ProfileForm>(key: K, value: ProfileForm[K]) {
    setForm((current) => ({
      ...current,
      [key]: value
    }));
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSaveMessage("");
    setIsSaving(true);

    saveProfile(form)
      .then((data) => {
        startTransition(() => {
          setProfile(data.profile);
          setForm(toForm(data.profile));
          setSaveMessage("プロフィールを保存しました。");
        });
      })
      .catch(() => {
        setSaveMessage("保存に失敗しました。バックエンド接続を確認してください。");
      })
      .finally(() => {
        setIsSaving(false);
      });
  }

  return (
    <main className="min-h-screen bg-[radial-gradient(circle_at_top,_#e0f2fe,_#f8fafc_45%,_#e2e8f0)] px-6 py-10 text-slate-900">
      <div className="mx-auto flex max-w-7xl flex-col gap-8">
        <section className="overflow-hidden rounded-[32px] border border-white/70 bg-slate-950 px-8 py-10 text-white shadow-[0_30px_80px_rgba(15,23,42,0.25)]">
          <p className="text-sm uppercase tracking-[0.3em] text-sky-300">TechCV</p>
          <div className="mt-4 grid gap-8 lg:grid-cols-[1.5fr_1fr]">
            <div>
              <h1 className="max-w-3xl text-4xl font-semibold tracking-tight md:text-5xl">
                profile エンドポイント経由でプロフィールを参照・登録・更新する画面
              </h1>
              <p className="mt-4 max-w-2xl text-sm leading-7 text-slate-300 md:text-base">
                初回表示時に GET /api/profile でプロフィールを取得し、フォーム保存時に
                PUT /api/profile へ反映します。Next.js 側は同一オリジンの /api を
                バックエンドへプロキシします。
              </p>
            </div>

            <div className="rounded-[28px] border border-white/10 bg-white/10 p-6 backdrop-blur">
              <p className="text-sm text-slate-300">状態</p>
              <dl className="mt-4 space-y-3 text-sm">
                <div className="flex items-center justify-between gap-4">
                  <dt className="text-slate-300">Backend</dt>
                  <dd className="rounded-full bg-emerald-400/20 px-3 py-1 text-emerald-200">
                    MySQL + Go API
                  </dd>
                </div>
                <div className="flex items-center justify-between gap-4">
                  <dt className="text-slate-300">Endpoint</dt>
                  <dd className="font-mono text-sky-200">/api/profile</dd>
                </div>
                <div className="flex items-center justify-between gap-4">
                  <dt className="text-slate-300">User ID</dt>
                  <dd className="font-mono text-slate-100">{profile?.userId ?? "user_01"}</dd>
                </div>
              </dl>
            </div>
          </div>
        </section>

        <section className="grid gap-8 xl:grid-cols-[1.3fr_0.9fr]">
          <form
            className="rounded-[28px] border border-slate-200/80 bg-white/90 p-8 shadow-[0_20px_60px_rgba(148,163,184,0.2)] backdrop-blur"
            onSubmit={handleSubmit}
          >
            <div className="flex flex-col gap-3 border-b border-slate-200 pb-6 md:flex-row md:items-end md:justify-between">
              <div>
                <p className="text-sm font-medium text-sky-700">Profile Editor</p>
                <h2 className="mt-1 text-2xl font-semibold">プロフィール編集</h2>
                <p className="mt-2 text-sm text-slate-600">
                  画面から入力した値をそのまま profile API へ保存します。
                </p>
              </div>
              <button
                className="inline-flex h-12 items-center justify-center rounded-full bg-slate-950 px-6 text-sm font-medium text-white transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:bg-slate-400"
                disabled={isLoading || isSaving}
                type="submit"
              >
                {isSaving ? "保存中..." : "保存する"}
              </button>
            </div>

            {loadError ? (
              <p className="mt-6 rounded-2xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-700">
                {loadError}
              </p>
            ) : null}

            {saveMessage ? (
              <p className="mt-6 rounded-2xl border border-sky-200 bg-sky-50 px-4 py-3 text-sm text-sky-700">
                {saveMessage}
              </p>
            ) : null}

            <div className="mt-8 grid gap-5 md:grid-cols-2">
              <label className="block text-sm font-medium text-slate-700">
                氏名
                <input
                  className={fieldClassName}
                  disabled={isLoading}
                  onChange={(event) => updateField("fullName", event.target.value)}
                  value={form.fullName}
                />
              </label>

              <label className="block text-sm font-medium text-slate-700">
                表示名
                <input
                  className={fieldClassName}
                  disabled={isLoading}
                  onChange={(event) => updateField("nickname", event.target.value)}
                  value={form.nickname}
                />
              </label>

              <label className="block text-sm font-medium text-slate-700">
                居住地
                <input
                  className={fieldClassName}
                  disabled={isLoading}
                  onChange={(event) => updateField("location", event.target.value)}
                  value={form.location}
                />
              </label>

              <label className="block text-sm font-medium text-slate-700">
                働き方の希望
                <input
                  className={fieldClassName}
                  disabled={isLoading}
                  onChange={(event) => updateField("preferredWorkStyle", event.target.value)}
                  value={form.preferredWorkStyle}
                />
              </label>

              <label className="block text-sm font-medium text-slate-700">
                メールアドレス
                <input
                  className={fieldClassName}
                  disabled={isLoading}
                  onChange={(event) => updateField("email", event.target.value)}
                  type="email"
                  value={form.email}
                />
              </label>

              <label className="block text-sm font-medium text-slate-700">
                電話番号
                <input
                  className={fieldClassName}
                  disabled={isLoading}
                  onChange={(event) => updateField("phone", event.target.value)}
                  value={form.phone}
                />
              </label>

              <label className="block text-sm font-medium text-slate-700">
                GitHub URL
                <input
                  className={fieldClassName}
                  disabled={isLoading}
                  onChange={(event) => updateField("githubUrl", event.target.value)}
                  value={form.githubUrl}
                />
              </label>

              <label className="block text-sm font-medium text-slate-700">
                Zenn URL
                <input
                  className={fieldClassName}
                  disabled={isLoading}
                  onChange={(event) => updateField("zennUrl", event.target.value)}
                  value={form.zennUrl}
                />
              </label>

              <label className="block text-sm font-medium text-slate-700">
                Qiita URL
                <input
                  className={fieldClassName}
                  disabled={isLoading}
                  onChange={(event) => updateField("qiitaUrl", event.target.value)}
                  value={form.qiitaUrl}
                />
              </label>

              <label className="block text-sm font-medium text-slate-700">
                Web サイト URL
                <input
                  className={fieldClassName}
                  disabled={isLoading}
                  onChange={(event) => updateField("websiteUrl", event.target.value)}
                  value={form.websiteUrl}
                />
              </label>
            </div>

            <label className="mt-5 block text-sm font-medium text-slate-700">
              概要
              <textarea
                className={`${fieldClassName} min-h-40 resize-y`}
                disabled={isLoading}
                onChange={(event) => updateField("summary", event.target.value)}
                value={form.summary}
              />
            </label>

            <div className="mt-6 grid gap-4 rounded-[24px] border border-slate-200 bg-slate-50 p-5 md:grid-cols-2">
              <label className="flex items-start gap-3 rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-700">
                <input
                  checked={form.showEmail}
                  className="mt-1 size-4 accent-sky-600"
                  disabled={isLoading}
                  onChange={(event) => updateField("showEmail", event.target.checked)}
                  type="checkbox"
                />
                <span>
                  メールアドレスを公開する
                </span>
              </label>

              <label className="flex items-start gap-3 rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-700">
                <input
                  checked={form.showPhone}
                  className="mt-1 size-4 accent-sky-600"
                  disabled={isLoading}
                  onChange={(event) => updateField("showPhone", event.target.checked)}
                  type="checkbox"
                />
                <span>
                  電話番号を公開する
                </span>
              </label>
            </div>
          </form>

          <aside className="rounded-[28px] border border-slate-200/80 bg-white/90 p-8 shadow-[0_20px_60px_rgba(148,163,184,0.2)] backdrop-blur">
            <p className="text-sm font-medium text-sky-700">Profile Snapshot</p>
            <h2 className="mt-1 text-2xl font-semibold">保存済みデータ</h2>

            {profile ? (
              <div className="mt-6 space-y-6">
                <div className="rounded-[24px] bg-slate-950 p-6 text-white">
                  <p className="text-xs uppercase tracking-[0.3em] text-sky-300">Profile</p>
                  <h3 className="mt-4 text-3xl font-semibold">
                    {profile.fullName || "未登録"}
                  </h3>
                  <p className="mt-2 text-sm text-slate-300">
                    {profile.nickname || "表示名未登録"}
                  </p>
                  <p className="mt-6 text-sm leading-7 text-slate-300">
                    {profile.summary || "プロフィール概要はまだ入力されていません。"}
                  </p>
                </div>

                <dl className="grid gap-4 text-sm text-slate-700">
                  <div className="rounded-2xl border border-slate-200 p-4">
                    <dt className="text-slate-500">所在地</dt>
                    <dd className="mt-1 font-medium">{profile.location || "-"}</dd>
                  </div>
                  <div className="rounded-2xl border border-slate-200 p-4">
                    <dt className="text-slate-500">メール</dt>
                    <dd className="mt-1 font-medium">
                      {profile.email || "-"}
                      <span className="ml-2 text-xs text-slate-500">
                        {profile.visibilitySettings?.email ? "公開" : "非公開"}
                      </span>
                    </dd>
                  </div>
                  <div className="rounded-2xl border border-slate-200 p-4">
                    <dt className="text-slate-500">電話番号</dt>
                    <dd className="mt-1 font-medium">
                      {profile.phone || "-"}
                      <span className="ml-2 text-xs text-slate-500">
                        {profile.visibilitySettings?.phone ? "公開" : "非公開"}
                      </span>
                    </dd>
                  </div>
                  <div className="rounded-2xl border border-slate-200 p-4">
                    <dt className="text-slate-500">リンク</dt>
                    <dd className="mt-1 space-y-1">
                      <p>{profile.githubUrl || "GitHub URL 未登録"}</p>
                      <p>{profile.zennUrl || "Zenn URL 未登録"}</p>
                      <p>{profile.qiitaUrl || "Qiita URL 未登録"}</p>
                      <p>{profile.websiteUrl || "Website URL 未登録"}</p>
                    </dd>
                  </div>
                  <div className="rounded-2xl border border-slate-200 p-4">
                    <dt className="text-slate-500">更新日時</dt>
                    <dd className="mt-1 font-medium">{formatDateTime(profile.updatedAt)}</dd>
                  </div>
                </dl>
              </div>
            ) : (
              <div className="mt-6 rounded-2xl border border-dashed border-slate-300 px-5 py-8 text-sm text-slate-500">
                {isLoading
                  ? "プロフィールを読み込んでいます..."
                  : "プロフィールを取得できませんでした。"}
              </div>
            )}
          </aside>
        </section>
      </div>
    </main>
  );
}
