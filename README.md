![sqluv_demo](./doc/image/demo.gif)

The **sqluv (sql + love)** provides a TUI for executing SQL queries on local CSV, TSV, and LTSV files. Currently, it supports only simple SELECT queries. In the future, it is planned to offer functionality to operate as a client for multiple DBMSs.

The sqluv is a command derived from [nao1215/sqly](https://github.com/nao1215/sqly). Its starting point is to provide a more user-friendly interface for writing SQL compared to sqly.

## How to install
### Use "go install"

```shell
go install github.com/nao1215/sqluv@latest
```

### Use homebrew

```shell
brew install nao1215/tap/sqluv
```

## Supported OS & go version

- Windows
- macOS
- Linux
- go1.24 or later

## How to use

### Syntax

```shell
sqluv [FILE_PATH]
```

※ Supported file formats: CSV, TSV, LTSV

## Key bindings

| Key | Description |
| --- | --- |
| ESC | Quit |
| Ctrl + D | Quit |
| TAB | Move to the next field |
| Shift + TAB | Move to the previous field |

### Contact
If you would like to send comments such as "find a bug" or "request for additional features" to the developer, please use one of the following contacts.

- [GitHub Issue](https://github.com/nao1215/sqluv/issues)

## LICENSE

[MIT License](./LICENSE)

