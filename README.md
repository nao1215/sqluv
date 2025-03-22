![dbms_demo](doc/image/dbms_demo.gif)

The **sqluv (sql + love)** provides a simple text user interface for multiple DBMSs and CSV/TSV/LTSV (local/http/https/s3) files. The sqluv automatically decompresses compressed files in gz, bz2, xz, and zst formats. You execute SQL queries for the connected DBMS or local/http/https files. The sqluv has the color theme feature, so you can change the color theme to your liking.

The sqluv is a command derived from [nao1215/sqly](https://github.com/nao1215/sqly). Its starting point is to provide a more user-friendly interface for writing SQL compared to sqly.

>[!WARNING]
> sqluv is under development. You use sqluv for **viewer**. Do not execute UPDATE or DELETE in the production environment. sqluv can not update or delete data in the local file, but it can update or delete data in the connected DBMS.

## How to install
### Use "go install"

```shell
go install github.com/nao1215/sqluv@latest
```

### Use homebrew

```shell
brew install nao1215/tap/sqluv
```

## Supported OS, File Format, Compressed Format, DBMS, go version

- Windows/macOS/Linux
- CSV/TSV/LTSV (file://, http://, https://)
- gz/bz2/xz/zst
- MySQL/PostgreSQL/SQLite3/SQL Server
- go1.24 or later

## How to use

### Syntax

```shell
sqluv [FILE_PATHS/HTTP URL/HTTPS URL]  ※ Supported file formats: CSV, TSV, LTSV
```

### Connect to DBMS

When you start the sqluv command without specifying a file path, the following screen will appear. 

![first_screen](doc/image/dbms_first.png)

Please enter the connection information for the DBMS you want to connect to.

![dbms_connection](doc/image/dbms_info.png)

If the connection is successful, database connection information will be saved in the configuration file. The next time you start the sqluv command, you will be able to select the DBMS you want to connect to from the list.

![dbms_list](doc/image/dbms_list.png)

![home_screen](doc/image/dbms_home.png)

## SQL query history

If you execute a SQL query, the history will be saved in the `~/.config/sqluv/history.db`. So, you can look up the history by pressing the history button.

![history_button](./doc/image/history_button.png)

If you select a history, the SQL query will be copied to the query text area.

![history_list](./doc/image/sql_query_history.png)


### Read from a file

Please specify a file path (or url) when executing the sqluv command. The file will be loaded before launching the TUI. When the sqluv import csv/tsv/ltsv, the sqluv checks the file extension and determines the file format. If the file extension is not csv/tsv/ltsv, the sqluv will display an error message. The sqluv does not automatically detect the file format.

![sqluv_demo](./doc/image/demo.gif)

## Key bindings

| Key | Description |
| --- | --- |
| ESC | Quit |
| Ctrl + d | Quit |
| Ctrl + c | Copy the selected sql query |
| Ctrl + v | Paste the copied text |
| Ctrl + x | Cut the selected text |
| Ctrl + t | Change the theme |
| TAB | Move to the next field |
| Shift + TAB | Move to the previous field |

## Color theme

### Defaulut
![color_default](./doc/image/color_default.png)

### Sublime
![color_sublime](./doc/image/color_sublime.png)

### VS Code
![color_vscode](./doc/image/color_vscode.png)

### Atom
![color_atom](./doc/image/color_atom.png)

### Dark
![color_dark](./doc/image/color_dark.png)

### Light
![color_light](./doc/image/color_light.png)

### Solarized
![color_solarized](./doc/image/color_solarized.png)

### Monokai
![color_monokai](./doc/image/color_monokai.png)

### Nord
![color_nord](./doc/image/color_nord.png)

### Cappuccino
![color_cappuccino](./doc/image/color_cappuccino.png)

### Gruvbox
![color_gruvbox](./doc/image/color_gruvbox.png)

### Tokyo Night
![color_tokyo_night](./doc/image/color_tokyo_night.png)

### Dracula
![color_dracula](./doc/image/color_dracula.png)

## Altenative Tools

|Name | Description |
|:----|:------------|
| [jorgerojas26/lazysql](https://github.com/jorgerojas26/lazysql) |A cross-platform TUI database management tool written in Go.|
| [vladbalmos/mitzasql](https://github.com/vladbalmos/mitzasql) | MySQL command line / text based interface client |
| [TaKO8Ki/gobang](https://github.com/TaKO8Ki/gobang) | A cross-platform TUI database management tool written in Rust |


## Contributing

First off, thanks for taking the time to contribute! See [CONTRIBUTING.md](./CONTRIBUTING.md) for more information. Contributions are not only related to development. For example, GitHub Star motivates me to develop! 


[![Star History Chart](https://api.star-history.com/svg?repos=nao1215/sqluv&type=Date)](https://star-history.com/#nao1215/sqluv&Date)

## Contact
If you would like to send comments such as "find a bug" or "request for additional features" to the developer, please use one of the following contacts.

- [GitHub Issue](https://github.com/nao1215/sqluv/issues)

## LICENSE

[MIT License](./LICENSE)

