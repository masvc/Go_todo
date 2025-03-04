# Go言語で作るToDoアプリ

## 概要
Go言語を使って、シンプルなToDoアプリを作成するプロジェクトです。
基本的なCRUD操作（作成・読取・更新・削除）を実装しながら、Go言語の基礎を学びます。

## 機能
- ToDoの追加
- ToDoの一覧表示
- ToDoの完了/未完了の切り替え
- ToDoの削除

## プロジェクト構成
```
go_test/
├── main.go        # メインプログラム
├── todo/          # ToDoの処理関連
└── templates/     # 画面テンプレート
```

## 使用技術
- Go言語
- 標準ライブラリ
  - net/http（Webサーバー）
  - html/template（画面表示）

## 開発環境
- Go（最新版）
- VS Code

## はじめ方
1. Goのインストール
   - [Go公式サイト](https://golang.org/dl/)からダウンロード

2. プロジェクトの初期化
```bash
go mod init go_test
```

## 学習できること
- Go言語の基本文法
- 構造体の使い方
- Webアプリケーションの基礎
- HTMLテンプレートの使用方法

## 参考リソース
- [Go公式チュートリアル](https://golang.org/doc/tutorial/getting-started)
- [Go by Example](https://gobyexample.com/)

## ライセンス
MIT

## 貢献について
プルリクエストや提案は大歓迎です。
