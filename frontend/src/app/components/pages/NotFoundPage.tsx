import { Link } from "react-router";
import { Home } from "lucide-react";
import { Button } from "../ui/button";

export function NotFoundPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center space-y-6 p-8">
        <div className="space-y-2">
          <h1 className="text-6xl font-bold text-gray-900">404</h1>
          <p className="text-xl text-gray-600">ページが見つかりませんでした</p>
        </div>
        <p className="text-gray-500">
          お探しのページは存在しないか、移動した可能性があります。
        </p>
        <div className="flex justify-center">
          <Button asChild>
            <Link to="/app">
              <Home className="w-4 h-4 mr-2" />
              ダッシュボードへ
            </Link>
          </Button>
        </div>
      </div>
    </div>
  );
}
