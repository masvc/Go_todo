# シンプルなToDoアプリ（Go + Vercel）

## 概要
シンプルで使いやすいToDoアプリケーションです。
フロントエンドとバックエンドの分離された設計で、Vercelでホスティングされています。

## 機能
- タスクの追加
- タスクの削除
- タスク一覧の表示

## 技術スタック
### フロントエンド
- HTML/CSS：シンプルなレイアウト
- JavaScript (Vanilla)：
  - Fetch APIによる非同期通信
  - DOMの動的な操作

### バックエンド
- Go言語：
  - HTTPサーバーの実装
  - JSONのエンコード/デコード
  - 並行処理（mutex）
- Vercel Serverless Functions

## プロジェクト構成
```
go_todo/
├── api/
│   └── todos/
│       └── index.go    # サーバーレス関数
├── public/
│   └── index.html      # フロントエンド
├── go.mod             # Go依存関係
└── vercel.json        # Vercel設定
```

## 開発環境のセットアップ
1. リポジトリのクローン
```bash
git clone https://github.com/masvc/Go_todo.git
cd Go_todo
```

2. Vercel CLIのインストール
```bash
npm install -g vercel
```

3. 開発サーバーの起動
```bash
vercel dev
```

## デプロイ
1. Vercel CLIでログイン
```bash
vercel login
```

2. デプロイの実行
```bash
vercel --prod
```

## 技術的なポイント
1. **シンプルな設計**
   - 必要最小限の機能に絞り込み
   - メンテナンスしやすいコード構造

2. **フロントエンド**
   - 直感的なUI/UX
   - バニラJavaScriptでの実装
   - 非同期通信の適切な処理

3. **バックエンド**
   - RESTful APIの設計
   - 並行処理の考慮（mutex）
   - エラーハンドリング

4. **デプロイ**
   - サーバーレスアーキテクチャ
   - 静的ファイルの配信
   - APIエンドポイントの管理

## ライセンス
MIT
