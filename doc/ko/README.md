![dbms_demo](../image/dbms_demo.gif)

[English](../../README.md) | [日本語](../ja/README.md) | [Русский](../ru/README.md) | [中文](../zh-cn/README.md) | [Español](../es/README.md) | [Français](../fr/README.md)


**sqluv (sql + love)**는 다양한 관계형 데이터베이스 관리 시스템(RDBMS) 및 CSV, TSV, LTSV 파일과 상호 작용하도록 설계된 간단한 텍스트 기반 사용자 인터페이스입니다. sqluv는 로컬 저장소, HTTPS, Amazon S3에서 CSV, TSV, LTSV 파일을 읽습니다. sqluv는 .gz, .bz2, .xz, .zst 형식의 압축 파일을 자동으로 압축 해제합니다. sqluv를 사용하면 SQL 쿼리 실행이 사용자 친화적인 경험이 되어 데이터베이스나 로컬 파일에 쉽게 연결할 수 있습니다.

sqluv는 [nao1215/sqly](https://github.com/nao1215/sqly)에서 파생된 명령입니다. 시작점은 sqly에 비해 더 사용자 친화적인 SQL 작성 인터페이스를 제공하는 것입니다.

>[!WARNING]
> sqluv는 개발 중입니다. sqluv를 **뷰어**로 사용하세요. 프로덕션 환경에서 UPDATE나 DELETE를 실행하지 마세요.

## 주요 기능
- 다중 DBMS 지원: MySQL, PostgreSQL, SQLite3, SQL Server와 같은 인기 있는 데이터베이스 시스템에 연결하고 상호 작용합니다.
- 파일 호환성: HTTPS, S3, 로컬 저장소에서 CSV, TSV, LTSV 형식의 파일에서 데이터를 읽습니다.
- 압축 파일 지원: .gz, .bz2, .xz, .zst 형식의 압축 파일을 자동으로 압축 해제합니다.
- 쿼리 히스토리: 쉬운 참조를 위해 SQL 쿼리 히스토리를 저장하고 액세스합니다. 퍼지 검색도 사용할 수 있습니다.
- 사용자 정의 가능한 테마: sqluv는 여러 색상 테마를 지원하여 사용자 선호도에 따라 인터페이스를 사용자 정의할 수 있습니다.

## 설치 방법
### "go install" 사용

```shell
go install github.com/nao1215/sqluv@latest
```

### homebrew 사용

```shell
brew install nao1215/tap/sqluv
```

## 지원되는 OS, 파일 형식, 압축 형식, DBMS, go 버전

- Windows/macOS/Linux
- CSV/TSV/LTSV (file://, http://, https://)
- gz/bz2/xz/zst
- MySQL/PostgreSQL/SQLite3/SQL Server
- go1.24 이상

## 사용 방법

### 구문

sqluv 인터페이스는 사용 편의성을 우선시합니다. 파일 경로를 지정하지 않고 시작하면 사용자에게 데이터베이스의 연결 세부 정보를 입력하라는 메시지가 표시됩니다. 구성이 저장되어 향후 쉽게 다시 연결할 수 있습니다. 기능의 간단한 개요는 다음과 같습니다:

```shell
sqluv [FILE_PATHS/HTTPS URL/S3 URL]  ※ 지원되는 파일 형식: CSV, TSV, LTSV
```

관련 파일 경로로 이 명령을 실행하여 사용자는 파일과의 상호 작용을 시작할 수 있습니다.

### DBMS에 연결

파일 경로를 지정하지 않고 sqluv 명령을 시작하면 다음 화면이 나타납니다.

![first_screen](../image/dbms_first.png)

연결하려는 DBMS의 연결 정보를 입력하세요.

![dbms_connection](../image/dbms_info.png)

연결이 성공하면 데이터베이스 연결 정보가 구성 파일에 저장됩니다. 다음에 sqluv 명령을 시작할 때 목록에서 연결하려는 DBMS를 선택할 수 있습니다.

![dbms_list](../image/dbms_list.png)

![home_screen](../image/dbms_home.png)

### SQL 쿼리 실행

SQL 쿼리를 실행하려면 쿼리 텍스트 영역에 SQL 쿼리를 입력하고 실행 버튼을 누르거나 `Ctrl + e`를 누르세요. 사이드바에서 테이블 이름을 선택하고 `Ctrl + e`를 누르면 sqluv는 `SELECT * FROM ${TABLE_NAME} LIMIT 100` 쿼리를 실행합니다.

테이블 이름을 검색하려면 사이드바에서 `/` 키를 누르세요. sqluv는 바닥글에 검색 필드를 표시합니다. `ESC` 키를 누르면 검색 필드가 지워집니다.

![sql_query](../image/search_tables.png)

열을 표시/숨기려면 사이드바에서 `Space` 키를 누르세요. 사이드바에서 `Space` 키를 누르면 sqluv는 테이블 DDL을 표시합니다.

![ddl](../image/ddl_info.png)

## SQL 쿼리 히스토리

SQL 쿼리를 실행하면 히스토리가 `~/.config/sqluv/history.db`에 저장됩니다. 따라서 히스토리 버튼을 눌러 히스토리를 조회할 수 있습니다.

![history_button](../image/history_button.png)

sqluv는 퍼지 검색을 지원합니다. 키워드를 입력하여 히스토리를 검색할 수 있습니다. 히스토리를 선택하면 SQL 쿼리가 입력 필드에 복사됩니다.

![history_list](../image/sql_query_history.gif)

### CSV/TSV/LTSV 가져오기

sqluv 명령을 실행할 때 파일 경로(또는 URL)를 지정하세요:

```shell
※ http/https의 파일
sqluv https://raw.githubusercontent.com/nao1215/sqluv/refs/heads/main/testdata/actor.csv
 
※ S3의 파일. 파일은 gz 형식으로 압축되어 있습니다.
sqluv s3://not-exist-s3-bucket/user.csv.gz
 
※ 여러 파일
sqluv https://raw.githubusercontent.com/nao1215/sqluv/refs/heads/main/testdata/actor.csv s3://not-exist-s3-bucket/user.tsv testdata/sample.ltsv
```

파일은 TUI를 시작하기 전에 로드됩니다. sqluv가 csv/tsv/ltsv를 가져올 때 sqluv는 파일 확장자를 확인하고 파일 형식을 결정합니다. 파일 확장자가 csv/tsv/ltsv가 아니면 sqluv는 오류 메시지를 표시합니다. sqluv는 파일 형식을 자동으로 감지하지 않습니다.

![sqluv_demo](../image/demo.gif)

### 결과를 파일에 저장

`Ctrl + s` 키를 눌러 결과를 파일에 저장할 수 있습니다. sqluv는 파일 경로를 입력하라고 요청합니다. 지원되는 파일 형식은 CSV, TSV, LTSV입니다.

![save_result](../image/file_save.png)

## 키 바인딩

| 키 | 설명 |
| --- | --- |
| Ctrl + d | 종료 |
| Ctrl + e | SQL 쿼리 실행 |
| Ctrl + h | SQL 쿼리 히스토리 표시 |
| Ctrl + c | 선택한 SQL 쿼리 복사 |
| Ctrl + v | 복사한 텍스트 붙여넣기 |
| Ctrl + x | 선택한 텍스트 잘라내기 |
| Ctrl + s | 결과를 파일에 저장 |
| Ctrl + t | 테마 변경 |
| /        | 테이블 이름 검색 (포커스가 사이드바에 있을 때)|
| ESC      | 검색 필드 지우기 (포커스가 사이드바에 있을 때)|
| Space    | 열 표시/숨기기 (포커스가 사이드바에 있을 때)|
| Enter    | 테이블 DDL 표시 (포커스가 사이드바에 있을 때)|
| F1       | 사이드바에 포커스 |
| F2       | 쿼리 텍스트 영역에 포커스 |
| F3       | 쿼리 결과 테이블에 포커스 |
| TAB | 다음 필드로 이동 |
| Shift + TAB | 이전 필드로 이동 |

## 색상 테마

### 기본값
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

## 대안 도구

|이름 | 설명 |
|:----|:------------|
| [jorgerojas26/lazysql](https://github.com/jorgerojas26/lazysql) |Go로 작성된 크로스 플랫폼 TUI 데이터베이스 관리 도구.|
| [vladbalmos/mitzasql](https://github.com/vladbalmos/mitzasql) | MySQL 명령줄 / 텍스트 기반 인터페이스 클라이언트 |
| [TaKO8Ki/gobang](https://github.com/TaKO8Ki/gobang) | Rust로 작성된 크로스 플랫폼 TUI 데이터베이스 관리 도구 |

## 기여하기

먼저 기여에 시간을 내주셔서 감사합니다! 자세한 정보는 [CONTRIBUTING.md](../../CONTRIBUTING.md)를 참조하세요. 기여는 개발에만 관련된 것이 아닙니다. 예를 들어, GitHub Star는 저에게 개발 동기를 부여합니다!

소셜 미디어나 블로그에서 sqluv를 소개해주시면 더 많은 사용자들이 발견할 수 있고, 개선을 위한 더 많은 아이디어를 모을 수 있습니다. sqluv 사용을 즐겨주시기 바랍니다!

[![Star History Chart](https://api.star-history.com/svg?repos=nao1215/sqluv&type=Date)](https://star-history.com/#nao1215/sqluv&Date)

## 연락처
"버그 발견" 또는 "추가 기능 요청"과 같은 의견을 개발자에게 보내고 싶으시다면 다음 연락처 중 하나를 사용하세요.

- [GitHub Issue](https://github.com/nao1215/sqluv/issues)

## 라이선스

[MIT License](../../LICENSE)