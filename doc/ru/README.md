![dbms_demo](../image/dbms_demo.gif)

[English](../../README.md) | [日本語](../ja/README.md) | [中文](../zh-cn/README.md) | [한국어](../ko/README.md) | [Español](../es/README.md) | [Français](../fr/README.md)

**sqluv (sql + love)** — это простой текстовый пользовательский интерфейс, предназначенный для взаимодействия с различными системами управления реляционными базами данных (СУБД), а также с файлами CSV, TSV и LTSV. sqluv читает файлы CSV, TSV и LTSV из локального хранилища, HTTPS и Amazon S3. sqluv автоматически распаковывает сжатые файлы в форматах .gz, .bz2, .xz и .zst. С помощью sqluv выполнение SQL-запросов становится удобным, обеспечивая простое подключение к базам данных или локальным файлам.

sqluv — это команда, производная от [nao1215/sqly](https://github.com/nao1215/sqly). Её отправной точкой является предоставление более удобного интерфейса для написания SQL по сравнению с sqly.

>[!WARNING]
> sqluv находится в разработке. Используйте sqluv в качестве **просмотрщика**. Не выполняйте UPDATE или DELETE в производственной среде.

## Ключевые функции
- Поддержка множественных СУБД: Подключение и взаимодействие с популярными системами баз данных, такими как MySQL, PostgreSQL, SQLite3 и SQL Server.
- Совместимость файлов: Чтение данных из файлов в форматах CSV, TSV и LTSV с HTTPS, S3 и локального хранилища.
- Поддержка сжатых файлов: Автоматическая распаковка сжатых файлов в форматах .gz, .bz2, .xz и .zst.
- История запросов: Сохранение и доступ к истории SQL-запросов для удобства. Также доступен нечёткий поиск.
- Настраиваемые темы: sqluv поддерживает множественные цветовые темы, позволяя настраивать интерфейс в соответствии с предпочтениями пользователя.

## Как установить
### Использование "go install"

```shell
go install github.com/nao1215/sqluv@latest
```

### Использование homebrew

```shell
brew install nao1215/tap/sqluv
```

## Поддерживаемые ОС, форматы файлов, форматы сжатия, СУБД, версии go

- Windows/macOS/Linux
- CSV/TSV/LTSV (file://, http://, https://)
- gz/bz2/xz/zst
- MySQL/PostgreSQL/SQLite3/SQL Server
- go1.24 или новее

## Как использовать

### Синтаксис

Интерфейс sqluv ставит в приоритет простоту использования. При запуске без указания пути к файлу пользователям предлагается ввести детали подключения к их базе данных. Конфигурация сохраняется, что позволяет легко переподключаться в будущем. Ниже приведён краткий обзор возможностей:

```shell
sqluv [FILE_PATHS/HTTPS URL/S3 URL]  ※ Поддерживаемые форматы файлов: CSV, TSV, LTSV
```

Запустив эту команду с соответствующими путями к файлам, пользователи могут начать взаимодействие с файлами.

### Подключение к СУБД

Когда вы запускаете команду sqluv без указания пути к файлу, появится следующий экран.

![first_screen](../image/dbms_first.png)

Пожалуйста, введите информацию о подключении для СУБД, к которой вы хотите подключиться.

![dbms_connection](../image/dbms_info.png)

При успешном подключении информация о подключении к базе данных будет сохранена в файле конфигурации. В следующий раз при запуске команды sqluv вы сможете выбрать СУБД для подключения из списка.

![dbms_list](../image/dbms_list.png)

![home_screen](../image/dbms_home.png)

### Выполнение SQL-запроса

Для выполнения SQL-запроса введите SQL-запрос в области текста запроса и нажмите кнопку выполнения или `Ctrl + e`. Когда вы выбираете имя таблицы на боковой панели и нажимаете `Ctrl + e`, sqluv выполняет запрос `SELECT * FROM ${TABLE_NAME} LIMIT 100`.

Для поиска имени таблицы нажмите клавишу `/` на боковой панели. sqluv отобразит поле поиска в нижней части. Если нажать клавишу `ESC`, поле поиска очистится.

![sql_query](../image/search_tables.png)

Для показа/скрытия столбцов нажмите клавишу `Space` на боковой панели. Если нажать клавишу `Space` на боковой панели, sqluv отобразит DDL таблицы.

![ddl](../image/ddl_info.png)

## История SQL-запросов

При выполнении SQL-запроса история сохраняется в `~/.config/sqluv/history.db`. Таким образом, вы можете просмотреть историю, нажав кнопку истории.

![history_button](../image/history_button.png)

sqluv поддерживает нечёткий поиск. Вы можете искать по истории, вводя ключевое слово. Если выбрать историю, SQL-запрос скопируется в поле ввода.

![history_list](../image/sql_query_history.gif)

### Импорт CSV/TSV/LTSV

Пожалуйста, укажите путь к файлу (или URL) при выполнении команды sqluv:

```shell
※ файл по http/https
sqluv https://raw.githubusercontent.com/nao1215/sqluv/refs/heads/main/testdata/actor.csv
 
※ файл в S3. Файл сжат в формате gz.
sqluv s3://not-exist-s3-bucket/user.csv.gz
 
※ Множественные файлы
sqluv https://raw.githubusercontent.com/nao1215/sqluv/refs/heads/main/testdata/actor.csv s3://not-exist-s3-bucket/user.tsv testdata/sample.ltsv
```

Файл будет загружен перед запуском TUI. Когда sqluv импортирует csv/tsv/ltsv, sqluv проверяет расширение файла и определяет формат файла. Если расширение файла не csv/tsv/ltsv, sqluv отобразит сообщение об ошибке. sqluv не определяет формат файла автоматически.

![sqluv_demo](../image/demo.gif)

### Сохранение результата в файл

Вы можете сохранить результат в файл, нажав клавишу `Ctrl + s`. sqluv попросит вас ввести путь к файлу. Поддерживаемые форматы файлов: CSV, TSV и LTSV.

![save_result](../image/file_save.png)

## Сочетания клавиш

| Клавиша | Описание |
| --- | --- |
| Ctrl + d | Выход |
| Ctrl + e | Выполнение SQL-запроса |
| Ctrl + h | Отображение истории SQL-запросов |
| Ctrl + c | Копирование выбранного SQL-запроса |
| Ctrl + v | Вставка скопированного текста |
| Ctrl + x | Вырезание выбранного текста |
| Ctrl + s | Сохранение результата в файл |
| Ctrl + t | Смена темы |
| /        | Поиск имени таблицы (когда фокус на боковой панели)|
| ESC      | Очистка поля поиска (когда фокус на боковой панели)|
| Space    | Показать/скрыть столбцы (когда фокус на боковой панели)|
| Enter    | Показать DDL таблицы (когда фокус на боковой панели)|
| F1       | Фокус на боковой панели |
| F2       | Фокус на области текста запроса |
| F3       | Фокус на таблице результата запроса |
| TAB | Переход к следующему полю |
| Shift + TAB | Переход к предыдущему полю |

## Цветовая тема

### По умолчанию
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

## Альтернативные инструменты

|Название | Описание |
|:----|:------------|
| [jorgerojas26/lazysql](https://github.com/jorgerojas26/lazysql) |Кроссплатформенный TUI-инструмент управления базами данных, написанный на Go.|
| [vladbalmos/mitzasql](https://github.com/vladbalmos/mitzasql) | Клиент интерфейса командной строки / текстового интерфейса MySQL |
| [TaKO8Ki/gobang](https://github.com/TaKO8Ki/gobang) | Кроссплатформенный TUI-инструмент управления базами данных, написанный на Rust |

## Вклад в разработку

Прежде всего, спасибо, что нашли время для вклада в проект! Смотрите [CONTRIBUTING.md](../../CONTRIBUTING.md) для получения дополнительной информации. Вклады касаются не только разработки. Например, GitHub Star мотивирует меня к разработке!

Если вы представите sqluv в социальных сетях или блогах, больше пользователей обнаружит его, и мы сможем собрать больше идей для улучшения. Надеемся, вам понравится использовать sqluv!

[![Star History Chart](https://api.star-history.com/svg?repos=nao1215/sqluv&type=Date)](https://star-history.com/#nao1215/sqluv&Date)

## Контакты
Если вы хотите отправить разработчику комментарии, такие как «найти ошибку» или «запрос на дополнительные функции», пожалуйста, используйте один из следующих контактов.

- [GitHub Issue](https://github.com/nao1215/sqluv/issues)

## ЛИЦЕНЗИЯ

[MIT License](../../LICENSE)