![dbms_demo](../image/dbms_demo.gif)

[English](../../README.md) | [日本語](../ja/README.md) | [Русский](../ru/README.md) | [한국어](../ko/README.md) | [Español](../es/README.md) | [Français](../fr/README.md)

**sqluv (sql + love)** 是一个简洁的基于文本的用户界面，专为与各种关系型数据库管理系统（RDBMS）以及 CSV、TSV 和 LTSV 文件进行交互而设计。sqluv 可以从本地存储、HTTPS 和 Amazon S3 读取 CSV、TSV 和 LTSV 文件。sqluv 会自动解压 .gz、.bz2、.xz 和 .zst 格式的压缩文件。使用 sqluv，执行 SQL 查询变得用户友好，能够轻松地与数据库或本地文件进行无缝连接。

sqluv 是从 [nao1215/sqly](https://github.com/nao1215/sqly) 派生的命令。其起点是提供比 sqly 更用户友好的 SQL 编写界面。

>[!WARNING]
> sqluv 正在开发中。您将 sqluv 用作**查看器**。请勿在生产环境中执行 UPDATE 或 DELETE。

## 主要功能
- 多 DBMS 支持：连接并与流行的数据库系统（如 MySQL、PostgreSQL、SQLite3 和 SQL Server）进行交互。
- 文件兼容性：从 HTTPS、S3 和本地存储读取 CSV、TSV 和 LTSV 格式的文件中的数据。
- 支持压缩文件：自动解压 .gz、.bz2、.xz 和 .zst 格式的压缩文件。
- 查询历史：保存并访问 SQL 查询历史以便参考。还支持模糊搜索。
- 可自定义主题：sqluv 支持多种颜色主题，使用户可以根据偏好自定义界面。

## 如何安装
### 使用 "go install"

```shell
go install github.com/nao1215/sqluv@latest
```

### 使用 homebrew

```shell
brew install nao1215/tap/sqluv
```

## 支持的操作系统、文件格式、压缩格式、DBMS、go 版本

- Windows/macOS/Linux
- CSV/TSV/LTSV (file://, http://, https://)
- gz/bz2/xz/zst
- MySQL/PostgreSQL/SQLite3/SQL Server
- go1.24 或更高版本

## 如何使用

### 语法

sqluv 界面注重易用性。在不指定文件路径的情况下启动时，用户会被提示输入数据库的连接详细信息。配置会被保存，方便将来重新连接。以下是功能的简要概述：

```shell
sqluv [FILE_PATHS/HTTPS URL/S3 URL]  ※ 支持的文件格式：CSV、TSV、LTSV
```

通过运行此命令并提供相关的文件路径，用户可以开始与文件进行交互。

### 连接到 DBMS

当您在不指定文件路径的情况下启动 sqluv 命令时，将出现以下屏幕。

![first_screen](../image/dbms_first.png)

请输入要连接的 DBMS 的连接信息。

![dbms_connection](../image/dbms_info.png)

如果连接成功，数据库连接信息将保存在配置文件中。下次启动 sqluv 命令时，您将能够从列表中选择要连接的 DBMS。

![dbms_list](../image/dbms_list.png)

![home_screen](../image/dbms_home.png)

### 执行 SQL 查询

要执行 SQL 查询，请在查询文本区域中输入 SQL 查询并按执行按钮或 `Ctrl + e`。当您在侧边栏中选择表名并按 `Ctrl + e` 时，sqluv 会执行 `SELECT * FROM ${TABLE_NAME} LIMIT 100` 查询。

要搜索表名，请在侧边栏中按 `/` 键。sqluv 将在页脚显示搜索字段。如果按 `ESC` 键，搜索字段将被清除。

![sql_query](../image/search_tables.png)

要显示/隐藏列，请在侧边栏中按 `Space` 键。如果在侧边栏中按 `Space` 键，sqluv 会显示表 DDL。

![ddl](../image/ddl_info.png)

## SQL 查询历史

如果执行了 SQL 查询，历史记录将保存在 `~/.config/sqluv/history.db` 中。因此，您可以通过按历史按钮查看历史记录。

![history_button](../image/history_button.png)

sqluv 支持模糊搜索。您可以通过输入关键词搜索历史记录。如果选择历史记录，SQL 查询将被复制到输入字段。

![history_list](../image/sql_query_history.gif)

### 导入 CSV/TSV/LTSV

在执行 sqluv 命令时请指定文件路径（或 URL）：

```shell
※ http/https 上的文件
sqluv https://raw.githubusercontent.com/nao1215/sqluv/refs/heads/main/testdata/actor.csv
 
※ S3 上的文件。文件以 gz 格式压缩。
sqluv s3://not-exist-s3-bucket/user.csv.gz
 
※ 多个文件
sqluv https://raw.githubusercontent.com/nao1215/sqluv/refs/heads/main/testdata/actor.csv s3://not-exist-s3-bucket/user.tsv testdata/sample.ltsv
```

文件将在启动 TUI 之前加载。当 sqluv 导入 csv/tsv/ltsv 时，sqluv 会检查文件扩展名并确定文件格式。如果文件扩展名不是 csv/tsv/ltsv，sqluv 会显示错误消息。sqluv 不会自动检测文件格式。

![sqluv_demo](../image/demo.gif)

### 将结果保存到文件

您可以通过按 `Ctrl + s` 键将结果保存到文件。sqluv 会要求您输入文件路径。支持的文件格式是 CSV、TSV 和 LTSV。

![save_result](../image/file_save.png)

## 键绑定

| 键 | 描述 |
| --- | --- |
| Ctrl + d | 退出 |
| Ctrl + e | 执行 SQL 查询 |
| Ctrl + h | 显示 SQL 查询历史 |
| Ctrl + c | 复制选定的 SQL 查询 |
| Ctrl + v | 粘贴复制的文本 |
| Ctrl + x | 剪切选定的文本 |
| Ctrl + s | 将结果保存到文件 |
| Ctrl + t | 更改主题 |
| /        | 搜索表名（当焦点在侧边栏上时）|
| ESC      | 清除搜索字段（当焦点在侧边栏上时）|
| Space    | 显示/隐藏列（当焦点在侧边栏上时）|
| Enter    | 显示表 DDL（当焦点在侧边栏上时）|
| F1       | 焦点在侧边栏 |
| F2       | 焦点在查询文本区域 |
| F3       | 焦点在查询结果表 |
| TAB | 移动到下一个字段 |
| Shift + TAB | 移动到上一个字段 |

## 颜色主题

### 默认
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

## 替代工具

|名称 | 描述 |
|:----|:------------|
| [jorgerojas26/lazysql](https://github.com/jorgerojas26/lazysql) |用 Go 编写的跨平台 TUI 数据库管理工具。|
| [vladbalmos/mitzasql](https://github.com/vladbalmos/mitzasql) | MySQL 命令行/基于文本的界面客户端 |
| [TaKO8Ki/gobang](https://github.com/TaKO8Ki/gobang) | 用 Rust 编写的跨平台 TUI 数据库管理工具 |

## 贡献

首先，感谢您抽出时间来贡献！有关更多信息，请参阅 [CONTRIBUTING.md](../../CONTRIBUTING.md)。贡献不仅仅与开发相关。例如，GitHub Star 激励我进行开发！

如果您在社交媒体或博客上介绍 sqluv，更多用户会发现它，我们可以收集更多改进的想法。希望您喜欢使用 sqluv！

[![Star History Chart](https://api.star-history.com/svg?repos=nao1215/sqluv&type=Date)](https://star-history.com/#nao1215/sqluv&Date)

## 联系
如果您想向开发者发送诸如"发现错误"或"请求附加功能"等评论，请使用以下联系方式之一。

- [GitHub Issue](https://github.com/nao1215/sqluv/issues)

## 许可证

[MIT License](../../LICENSE)