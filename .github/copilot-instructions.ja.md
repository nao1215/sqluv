# プロジェクト概要

本プロジェクトで開発するアプリは、`sqluv`と呼ばれ、Golang製のクロスプラットホームアプリです。　`sqluv`は、Text User Interface（以下、TUIと略す）を持つターミナルアプリであり、基本機能は「複数のRDBMSに接続」および「SQLクエリを実行」です。RDBMSだけでなく、ローカル／HTTPS／Amazon S3にあるCSV、TSV、LTSVのファイルをSQLite3（インメモリ）に読み込み、SQLクエリを実行することができます。RDBMS接続、ファイル読み込みを同時に実行することはできず、起動時引数がなければRDBMS接続を試み、起動時引数があればファイル読み込みを試みます。

# プロジェクトのディレクトリ構成

- config：環境変数、起動時引数、設定ファイル（DB接続情報ファイル、カラーテーマ情報ファイル、SQLクエリ履歴ファイル）を管理する
- di：Dependency Injectionを管理する
- doc：ドキュメントを管理する
- docker：Dockerfileが利用するファイルを管理する
- domein/model：ドメインオブジェクトを管理する
- domain/repository：一時的な記録、永続化に関するインターフェースを管理する
- infrastructure/memory：domain/repositoryで定義されたインターフェースの中で、一時的な記録に関するものを実装する
- infrastructure/mock：domain/repositoryで定義されたインターフェースをモックする
- infrastructure/persistence：domain/repositoryで定義されたインターフェースの中で、永続化に関するものを実装する
- infrastructure：infrastructure以下のディレクトリで共通利用する関数やエラーを定義する
- interactor/mock：usecaseで定義されたインターフェースをモックする
- interactor：usecaseで定義されたインターフェースを実装する
- testdata：テストデータを管理する
- tui：TUIを管理する
- usecase：tuiから呼び出されるユースケースインターフェースを管理する

# Text User Interfaceの概要

`github.com/rivo/tview`のライブラリを用いて、TUIを提供している。TUIは、以下のコンポーネントで構成されている。

- TUIを管理する構造体（`tui.TUI`構造体）
- 各コンポーネントを管理する構造体（`tui.home`構造体）
- テーブルおよびカラム一覧を表示するサイドバー（`tui.sidebar`構造体）
- 実行したいSQLクエリを書くクエリテキストエリア（`tui.queryTextArea`構造体）
- SQLクエリ実行結果を表示するテーブルビュー（`tui.queryResultTable`構造体）
- SQLクエリ実行結果を表示するテーブルビュー（`tui.queryResultTable`構造体）
- SQLクエリ実行ボタン（`tui.executeButton`構造体）
- SQLクエリヒストリーボタン（`tui.historyButton`構造体）
- SQLクエリ実行時の統計結果を表示するビュー（`tui.rowStatistics`構造体）
- エラーやinformationを表示するダイアログ（`tui.dialog`構造体）
- キーバインド情報を示すフッター（`tui.footer`構造体）
- カラーテーマを管理する構造体（`tui.Theme`構造体）

上記は、RDBMS接続後、もしくはファイル読み込み後に表示されるTUIの構成である。RDBMS接続前は、モーダル（`tui.connectionModal`構造体）を表示する。

# SQLクエリ履歴の仕様

`history.db` ファイルは、`sqluv` の設定ファイルを管理するディレクトリ（例: `~/.config/sqluv`）に存在する。`sqluv` がSQLクエリを正常に実行するたびに、そのクエリ履歴をSQLite3を使用して `history.db` に保存する。`sqluv` の設定ファイルに関連するコードは `config/config_file.go` にある。

`sqluv` は、ホーム画面で [History] ボタンを押下するか、`Ctrl-h` を入力すると、SQLクエリ履歴リスト画面を表示する。また、`sqluv` はクエリ履歴のファジー検索をサポートしており、`github.com/lithammer/fuzzysearch/fuzzy` を使用して実装されている。ファジー検索の実装は `tui/tui.go` にある。

# インターフェースの仕様

インターフェース名はGoの慣習に従い、動詞 + "er"を結び付けた名称とする。一つのインターフェースは、一つのメソッドのみを持つ。インターフェース名は名詞+動詞+"er"の形式とし、インターフェースが持つメソッド名は動詞+"名詞"の形式となる。インターフェース名が決まると、メソッド名は自動的に決定できる。具体的には、`FileGetter`インターフェースは、`GetFile()`メソッドを持つ。

# ホーム画面におけるキーバインド

|キー	| 説明 |
| --- | --- |
| Ctrl + d	|アプリケーションを終了|
| Ctrl + e	|SQLクエリを実行 |
| Ctrl + h	|SQLクエリ履歴を表示 |
| Ctrl + c	|選択したSQLクエリをコピー |
| Ctrl + v	|コピーされたテキストを貼り付け |
| Ctrl + x	|選択したテキストをカット |
| Ctrl + s	|結果をファイルに保存 |
| Ctrl + t	|テーマを変更 |
| F1        |サイドバーにフォーカス |
| F2        |クエリ入力エリアにフォーカス |
| F3        |クエリ結果テーブルにフォーカス |
| TAB	|次のフィールドに移動 |
| Shift + TAB	|前のフィールドに移動 |
