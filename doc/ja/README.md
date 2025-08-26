![dbms_demo](../image/dbms_demo.gif)

[English](../../README.md) | [Русский](../ru/README.md) | [中文](../zh-cn/README.md) | [한국어](../ko/README.md) | [Español](../es/README.md) | [Français](../fr/README.md)

**sqluv (sql + love)** は、様々なリレーショナルデータベース管理システム（RDBMS）およびCSV、TSV、LTSVファイルとやり取りするための直感的なテキストベースユーザーインターフェースです。sqluvは、ローカルストレージ、HTTPS、Amazon S3からCSV、TSV、LTSVファイルを読み込みます。sqluvは、.gz、.bz2、.xz、.zst形式の圧縮ファイルを自動的に展開します。sqluvを使用すると、SQLクエリの実行がユーザーフレンドリーな体験となり、データベースやローカルファイルへのシームレスな接続が簡単に行えます。

sqluvは[nao1215/sqly](https://github.com/nao1215/sqly)から派生したコマンドです。その出発点は、sqlyと比較してよりユーザーフレンドリーなSQL記述インターフェースを提供することです。

>[!WARNING]
> sqluvは開発中です。sqluvは**ビューワー**として使用してください。本番環境でUPDATEやDELETEを実行しないでください。

## 主要機能
- マルチDBMSサポート：MySQL、PostgreSQL、SQLite3、SQL Serverなどの主要なデータベースシステムに接続・操作が可能。
- ファイル互換性：HTTPS、S3、ローカルストレージからCSV、TSV、LTSV形式のファイルからデータを読み取り。
- 圧縮ファイルサポート：.gz、.bz2、.xz、.zst形式の圧縮ファイルを自動的に展開。
- クエリ履歴：SQLクエリ履歴の保存とアクセスが可能。ファジー検索も利用可能。
- カスタマイズ可能なテーマ：sqluvは複数のカラーテーマをサポートし、ユーザーの好みに基づいてインターフェースをカスタマイズ可能。

## インストール方法
### "go install"を使用

```shell
go install github.com/nao1215/sqluv@latest
```

### homebrewを使用

```shell
brew install nao1215/tap/sqluv
```

## サポートされているOS、ファイル形式、圧縮形式、DBMS、goバージョン

- Windows/macOS/Linux
- CSV/TSV/LTSV (file://, http://, https://)
- gz/bz2/xz/zst
- MySQL/PostgreSQL/SQLite3/SQL Server
- go1.24以降

## 使用方法

### 構文

sqluvインターフェースは使いやすさを優先しています。ファイルパスを指定せずに起動すると、データベースの接続詳細を入力するよう求められます。設定は保存され、将来の再接続が簡単になります。以下に機能の簡単な概要を示します：

```shell
sqluv [FILE_PATHS/HTTPS URL/S3 URL]  ※ サポートされているファイル形式：CSV、TSV、LTSV
```

関連するファイルパスでこのコマンドを実行することで、ユーザーはファイルとのやり取りを開始できます。

### DBMSへの接続

ファイルパスを指定せずにsqluvコマンドを開始すると、以下の画面が表示されます。

![first_screen](../image/dbms_first.png)

接続したいDBMSの接続情報を入力してください。

![dbms_connection](../image/dbms_info.png)

接続が成功すると、データベース接続情報が設定ファイルに保存されます。次回sqluvコマンドを開始する際には、リストから接続したいDBMSを選択できるようになります。

![dbms_list](../image/dbms_list.png)

![home_screen](../image/dbms_home.png)

### SQLクエリの実行

SQLクエリを実行するには、クエリテキストエリアにSQLクエリを入力し、実行ボタンを押すか`Ctrl + e`を押します。サイドバーでテーブル名を選択して`Ctrl + e`を押すと、sqluvは`SELECT * FROM ${TABLE_NAME} LIMIT 100`クエリを実行します。

テーブル名を検索するには、サイドバーで`/`キーを押します。sqluvはフッターに検索フィールドを表示します。`ESC`キーを押すと、検索フィールドがクリアされます。

![sql_query](../image/search_tables.png)

列を表示/非表示にするには、サイドバーで`Space`キーを押します。サイドバーで`Space`キーを押すと、sqluvはテーブルDDLを表示します。

![ddl](../image/ddl_info.png)

## SQLクエリ履歴

SQLクエリを実行すると、履歴は`~/.config/sqluv/history.db`に保存されます。そのため、履歴ボタンを押すことで履歴を確認できます。

![history_button](../image/history_button.png)

sqluvはファジー検索をサポートしています。キーワードを入力して履歴を検索できます。履歴を選択すると、SQLクエリが入力フィールドにコピーされます。

![history_list](../image/sql_query_history.gif)

### CSV/TSV/LTSVのインポート

sqluvコマンドを実行する際にファイルパス（またはURL）を指定してください：

```shell
※ http/httpsのファイル
sqluv https://raw.githubusercontent.com/nao1215/sqluv/refs/heads/main/testdata/actor.csv
 
※ S3のファイル。ファイルはgz形式で圧縮されています。
sqluv s3://not-exist-s3-bucket/user.csv.gz
 
※ 複数ファイル
sqluv https://raw.githubusercontent.com/nao1215/sqluv/refs/heads/main/testdata/actor.csv s3://not-exist-s3-bucket/user.tsv testdata/sample.ltsv
```

ファイルはTUIを起動する前にロードされます。sqluvがcsv/tsv/ltsvをインポートする際、sqluvはファイル拡張子をチェックしてファイル形式を決定します。ファイル拡張子がcsv/tsv/ltsvでない場合、sqluvはエラーメッセージを表示します。sqluvはファイル形式を自動検出しません。

![sqluv_demo](../image/demo.gif)

### 結果をファイルに保存

`Ctrl + s`キーを押すことで結果をファイルに保存できます。sqluvはファイルパスの入力を求めます。サポートされているファイル形式は、CSV、TSV、LTSVです。

![save_result](../image/file_save.png)

## キーバインディング

| キー | 説明 |
| --- | --- |
| Ctrl + d | 終了 |
| Ctrl + e | SQLクエリの実行 |
| Ctrl + h | SQLクエリ履歴の表示 |
| Ctrl + c | 選択したSQLクエリをコピー |
| Ctrl + v | コピーしたテキストを貼り付け |
| Ctrl + x | 選択したテキストをカット |
| Ctrl + s | 結果をファイルに保存 |
| Ctrl + t | テーマの変更 |
| /        | テーブル名の検索（サイドバーにフォーカスがある場合）|
| ESC      | 検索フィールドのクリア（サイドバーにフォーカスがある場合）|
| Space    | 列の表示/非表示（サイドバーにフォーカスがある場合）|
| Enter    | テーブルDDLの表示（サイドバーにフォーカスがある場合）|
| F1       | サイドバーにフォーカス |
| F2       | クエリテキストエリアにフォーカス |
| F3       | クエリ結果テーブルにフォーカス |
| TAB | 次のフィールドに移動 |
| Shift + TAB | 前のフィールドに移動 |

## カラーテーマ

### デフォルト
![color_default](../image/color_default.png)

### Sublime
![color_sublime](../image/color_sublime.png)

### VS Code
![color_vscode](../image/color_vscode.png)

### Atom
![color_atom](../image/color_atom.png)

### Dark
![color_dark](../image/color_dark.png)

### Light
![color_light](../image/color_light.png)

### Solarized
![color_solarized](../image/color_solarized.png)

### Monokai
![color_monokai](../image/color_monokai.png)

### Nord
![color_nord](../image/color_nord.png)

### Cappuccino
![color_cappuccino](../image/color_cappuccino.png)

### Gruvbox
![color_gruvbox](../image/color_gruvbox.png)

### Tokyo Night
![color_tokyo_night](../image/color_tokyo_night.png)

### Dracula
![color_dracula](../image/color_dracula.png)

### Cyber Neon
![color_cyber_neon](../image/color_cyber_neon.png)

### Earthy Tones
![color_earthy_tones](../image/color_earthy_tones.png)

### Royal Inferno
![color_royal_inferno](../image/color_royal_inferno.png)

## 代替ツール

|名前 | 説明 |
|:----|:------------|
| [jorgerojas26/lazysql](https://github.com/jorgerojas26/lazysql) |Goで書かれたクロスプラットフォームTUIデータベース管理ツール。|
| [vladbalmos/mitzasql](https://github.com/vladbalmos/mitzasql) | MySQLコマンドライン/テキストベースインターフェースクライアント |
| [TaKO8Ki/gobang](https://github.com/TaKO8Ki/gobang) | Rustで書かれたクロスプラットフォームTUIデータベース管理ツール |

## コントリビュート

まず、コントリビュートに時間を割いていただき、ありがとうございます！詳細については[CONTRIBUTING.md](../../CONTRIBUTING.md)を参照してください。コントリビュートは開発に関連するものだけではありません。例えば、GitHub Starは私の開発モチベーションになります！

ソーシャルメディアやブログでsqluvを紹介していただけると、より多くのユーザーに発見され、改善のためのより多くのアイデアを集めることができます。sqluvを楽しんでお使いいただけることを願っています！

[![Star History Chart](https://api.star-history.com/svg?repos=nao1215/sqluv&type=Date)](https://star-history.com/#nao1215/sqluv&Date)

## 連絡先
「バグを見つけた」や「追加機能のリクエスト」などのコメントを開発者に送信したい場合は、以下の連絡先のいずれかをご利用ください。

- [GitHub Issue](https://github.com/nao1215/sqluv/issues)

## ライセンス

[MIT License](../../LICENSE)