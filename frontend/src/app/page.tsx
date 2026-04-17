export default function HomePage() {
  return (
    <main className="mx-auto flex min-h-screen max-w-5xl flex-col gap-8 px-6 py-16">
      <section className="rounded-xl border border-slate-200 bg-white p-8 shadow-sm">
        <p className="text-sm font-medium text-blue-600">TechCV</p>
        <h1 className="mt-3 text-4xl font-bold tracking-tight text-slate-900">
          Webエンジニア向け経歴書管理システム
        </h1>
        <p className="mt-4 max-w-2xl text-slate-600">
          案件中心で職務経歴を管理し、応募先ごとの経歴書バリエーションを
          プレビュー・出力できるプロダクトの実装リポジトリです。
        </p>
      </section>

      <section className="grid gap-4 md:grid-cols-3">
        <div className="rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="text-lg font-semibold text-slate-900">Frontend</h2>
          <p className="mt-2 text-sm text-slate-600">
            Next.js, Tailwind CSS, shadcn/ui を採用します。
          </p>
        </div>
        <div className="rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="text-lg font-semibold text-slate-900">Backend</h2>
          <p className="mt-2 text-sm text-slate-600">
            Go, Chi, sqlc, PostgreSQL を採用します。
          </p>
        </div>
        <div className="rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
          <h2 className="text-lg font-semibold text-slate-900">Auth</h2>
          <p className="mt-2 text-sm text-slate-600">
            MVP では Cookie ベースの SessionAuth を主軸にします。
          </p>
        </div>
      </section>
    </main>
  );
}
