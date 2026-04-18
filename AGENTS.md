# AGENTS.md

## 適用範囲
- より深い階層にある `AGENTS.md` で上書きされない限り、これらの指示は `techcv-app` リポジトリ全体に適用されます。

## レイヤリング規則
- transport に関する責務は transport レイヤに留めてください。
- JSON、HTTP、シリアライズ用のタグは domain model ではなく handler レイヤの DTO に持たせてください。
- domain 層や usecase 層の struct に `json`、`form` などの transport 用タグを付けないでください。
- HTTP の payload と usecase の入出力の変換には、handler レイヤの request/response struct を使ってください。

## バックエンド構成
- API の成長が見込まれる場合は、`handler`、`usecase`、`repository` を分離する構成を優先してください。
- domain model は framework や protocol の詳細ではなく、業務上の意味に集中させてください。
- shared helper の追加は、明確に横断的な責務に限ってください。
- データベースの schema 変更は `backend/migrations/schema.sql` で管理してください。
- MySQL の schema 変更は、`mysqldef` 互換の schema ファイルとコマンドで管理してください。
- repository の SQL は実装内にインラインで書かず、`sqlc` の query file で管理してください。
- `sqlc` で管理する SQL の正本は `backend/sqlc.yaml` と `backend/internal/**/repository/queries/*.sql` として扱ってください。

## 作業スタイル
- 初期実装は最小限に留め、あとから拡張しやすい形を保ってください。
- framework の暗黙的な magic よりも、レイヤ間の明示的な mapping を優先してください。
