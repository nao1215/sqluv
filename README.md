![dbms_demo](doc/image/dbms_demo.gif)

The **sqluv (sql + love)** provides a simple text user interface for multiple DBMSs and local CSV/TSV/LTSV files. You execute SQL queries for the connected DBMS or local files.

The sqluv is a command derived from [nao1215/sqly](https://github.com/nao1215/sqly). Its starting point is to provide a more user-friendly interface for writing SQL compared to sqly.

>[!WARNING]
> sqluv is under development. Do not execute UPDATE or DELETE in the production environment. Currently sqluv starts transactions and executes commits on a per SQL query basis. Therefore, commits may be executed at unexpected times.

## How to install
### Use "go install"

```shell
go install github.com/nao1215/sqluv@latest
```

### Use homebrew

```shell
brew install nao1215/tap/sqluv
```

## Supported OS, DBMS, go version

- Windows
- macOS
- Linux
- MySQL
- PostgreSQL
- SQLite3
- SQL Server
- go1.24 or later

## How to use

### Syntax

```shell
sqluv [FILE_PATHS]
```

※ Supported file formats: CSV, TSV, LTSV

### Connect to DBMS

When you start the sqluv command without specifying a file path, the following screen will appear. 

![first_screen](doc/image/dbms_first.png)

Please enter the connection information for the DBMS you want to connect to.

![dbms_connection](doc/image/dbms_info.png)

If the connection is successful, database connection information will be saved in the configuration file. The next time you start the sqluv command, you will be able to select the DBMS you want to connect to from the list.

![dbms_list](doc/image/dbms_list.png)

## SQL query history

If you execute a SQL query, the history will be saved in the `~/.config/sqluv/history.db`. So, you can look up the history by pressing the history button.

![history_button](./doc/image/history_button.png)

If you select a history, the SQL query will be copied to the query text area.

![history_list](./doc/image/sql_query_history.png)


### Read from a file

Please specify a file path when executing the sqluv command. The file will be loaded before launching the TUI.
If you start the sqluv command without specifying a file path, it will enter DBMS connection mode, and local files cannot be loaded.

![sqluv_demo](./doc/image/demo.gif)

## Key bindings

| Key | Description |
| --- | --- |
| ESC | Quit |
| Ctrl + D | Quit |
| TAB | Move to the next field |
| Shift + TAB | Move to the previous field |

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

