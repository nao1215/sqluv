![dbms_demo](../image/dbms_demo.gif)

[English](../../README.md) | [日本語](../ja/README.md) | [Русский](../ru/README.md) | [中文](../zh-cn/README.md) | [한국어](../ko/README.md) | [Français](../fr/README.md)

**sqluv (sql + love)** es una interfaz de usuario basada en texto sencilla diseñada para interactuar con varios Sistemas de Gestión de Bases de Datos Relacionales (RDBMS) así como archivos CSV, TSV y LTSV. sqluv lee archivos CSV, TSV y LTSV desde almacenamiento local, HTTPS y Amazon S3. sqluv descomprime automáticamente archivos comprimidos en formatos .gz, .bz2, .xz y .zst. Con sqluv, ejecutar consultas SQL se convierte en una experiencia amigable para el usuario, permitiendo conexiones sencillas a bases de datos o archivos locales con facilidad.

sqluv es un comando derivado de [nao1215/sqly](https://github.com/nao1215/sqly). Su punto de partida es proporcionar una interfaz más amigable para escribir SQL comparado con sqly.

>[!WARNING]
> sqluv está en desarrollo. Use sqluv como **visor**. No ejecute UPDATE o DELETE en el entorno de producción.

## Características Principales
- Soporte Multi-DBMS: Conecte e interactúe con sistemas de bases de datos populares como MySQL, PostgreSQL, SQLite3 y SQL Server.
- Compatibilidad de Archivos: Lea datos de archivos en formatos CSV, TSV y LTSV desde HTTPS, S3 y almacenamiento local.
- Soporte de Archivos Comprimidos: Descomprime automáticamente archivos comprimidos en formatos .gz, .bz2, .xz y .zst.
- Historial de Consultas: Guarde y acceda al historial de consultas SQL para referencia fácil. También está disponible la búsqueda difusa.
- Temas Personalizables: sqluv soporta múltiples temas de color, permitiendo la personalización de la interfaz basada en las preferencias del usuario.

## Cómo instalar
### Usar "go install"

```shell
go install github.com/nao1215/sqluv@latest
```

### Usar homebrew

```shell
brew install nao1215/tap/sqluv
```

## SO Soportados, Formato de Archivo, Formato Comprimido, DBMS, versión de go

- Windows/macOS/Linux
- CSV/TSV/LTSV (file://, http://, https://)
- gz/bz2/xz/zst
- MySQL/PostgreSQL/SQLite3/SQL Server
- go1.24 o posterior

## Cómo usar

### Sintaxis

La interfaz de sqluv prioriza la facilidad de uso. Al iniciar sin especificar una ruta de archivo, se solicita a los usuarios que ingresen detalles de conexión para su base de datos. La configuración se guarda, permitiendo fáciles reconexiones en el futuro. A continuación se presenta un breve resumen de las capacidades:

```shell
sqluv [FILE_PATHS/HTTPS URL/S3 URL]  ※ Formatos de archivo soportados: CSV, TSV, LTSV
```

Al ejecutar este comando con las rutas de archivo relevantes, los usuarios pueden iniciar interacciones con archivos.

### Conectar a DBMS

Cuando inicie el comando sqluv sin especificar una ruta de archivo, aparecerá la siguiente pantalla.

![first_screen](../image/dbms_first.png)

Por favor ingrese la información de conexión para el DBMS al que desea conectarse.

![dbms_connection](../image/dbms_info.png)

Si la conexión es exitosa, la información de conexión de la base de datos se guardará en el archivo de configuración. La próxima vez que inicie el comando sqluv, podrá seleccionar el DBMS al que desea conectarse desde la lista.

![dbms_list](../image/dbms_list.png)

![home_screen](../image/dbms_home.png)

### Ejecutar consulta SQL

Para ejecutar una consulta SQL, ingrese la consulta SQL en el área de texto de consulta y presione el botón ejecutar o `Ctrl + e`. Cuando seleccione el nombre de la tabla en la barra lateral y presione `Ctrl + e`, sqluv ejecuta la consulta `SELECT * FROM ${TABLE_NAME} LIMIT 100`.

Para buscar un nombre de tabla, presione la tecla `/` en la barra lateral. sqluv mostrará el campo de búsqueda en el pie de página. Si presiona la tecla `ESC`, el campo de búsqueda se limpiará.

![sql_query](../image/search_tables.png)

Para mostrar/ocultar columnas, presione la tecla `Space` en la barra lateral. Si presiona la tecla `Space` en la barra lateral, sqluv muestra el DDL de la tabla.

![ddl](../image/ddl_info.png)

## Historial de consultas SQL

Si ejecuta una consulta SQL, el historial se guardará en `~/.config/sqluv/history.db`. Por lo tanto, puede consultar el historial presionando el botón de historial.

![history_button](../image/history_button.png)

sqluv soporta búsqueda difusa. Puede buscar en el historial escribiendo la palabra clave. Si selecciona el historial, la consulta SQL se copiará al campo de entrada.

![history_list](../image/sql_query_history.gif)

### Importar CSV/TSV/LTSV

Por favor especifique una ruta de archivo (o url) al ejecutar el comando sqluv:

```shell
※ archivo en http/https
sqluv https://raw.githubusercontent.com/nao1215/sqluv/refs/heads/main/testdata/actor.csv
 
※ archivo en s3. el archivo está comprimido en formato gz.
sqluv s3://not-exist-s3-bucket/user.csv.gz
 
※ Múltiples archivos
sqluv https://raw.githubusercontent.com/nao1215/sqluv/refs/heads/main/testdata/actor.csv s3://not-exist-s3-bucket/user.tsv testdata/sample.ltsv
```

El archivo se cargará antes de lanzar la TUI. Cuando sqluv importa csv/tsv/ltsv, sqluv verifica la extensión del archivo y determina el formato del archivo. Si la extensión del archivo no es csv/tsv/ltsv, sqluv mostrará un mensaje de error. sqluv no detecta automáticamente el formato del archivo.

![sqluv_demo](../image/demo.gif)

### Guardar el resultado en un archivo

Puede guardar el resultado en un archivo presionando la tecla `Ctrl + s`. sqluv le pedirá que ingrese la ruta del archivo. Los formatos de archivo soportados son CSV, TSV y LTSV.

![save_result](../image/file_save.png)

## Combinaciones de teclas

| Tecla | Descripción |
| --- | --- |
| Ctrl + d | Salir |
| Ctrl + e | Ejecutar la consulta SQL |
| Ctrl + h | Mostrar el historial de consultas SQL |
| Ctrl + c | Copiar la consulta SQL seleccionada |
| Ctrl + v | Pegar el texto copiado |
| Ctrl + x | Cortar el texto seleccionado |
| Ctrl + s | Guardar el resultado en un archivo |
| Ctrl + t | Cambiar el tema |
| /        | Buscar el nombre de la tabla (cuando el foco está en la barra lateral)|
| ESC      | Limpiar el campo de búsqueda (cuando el foco está en la barra lateral)|
| Space    | Mostrar/Ocultar Columnas (cuando el foco está en la barra lateral)|
| Enter    | Mostrar el DDL de la tabla (cuando el foco está en la barra lateral)|
| F1       | Enfocar en la barra lateral |
| F2       | Enfocar en el área de texto de consulta |
| F3       | Enfocar en la tabla de resultados de consulta |
| TAB | Moverse al siguiente campo |
| Shift + TAB | Moverse al campo anterior |

## Tema de color

### Por defecto
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

## Herramientas Alternativas

|Nombre | Descripción |
|:----|:------------|
| [jorgerojas26/lazysql](https://github.com/jorgerojas26/lazysql) |Una herramienta de gestión de bases de datos TUI multiplataforma escrita en Go.|
| [vladbalmos/mitzasql](https://github.com/vladbalmos/mitzasql) | Cliente de interfaz de línea de comandos / basado en texto MySQL |
| [TaKO8Ki/gobang](https://github.com/TaKO8Ki/gobang) | Una herramienta de gestión de bases de datos TUI multiplataforma escrita en Rust |

## Contribuir

En primer lugar, ¡gracias por tomarse el tiempo para contribuir! Vea [CONTRIBUTING.md](../../CONTRIBUTING.md) para más información. Las contribuciones no están relacionadas solo con el desarrollo. ¡Por ejemplo, GitHub Star me motiva a desarrollar!

Si presenta sqluv en redes sociales o blogs, más usuarios lo descubrirán, y podemos reunir más ideas para mejorar. ¡Esperamos que disfrute usando sqluv!

[![Star History Chart](https://api.star-history.com/svg?repos=nao1215/sqluv&type=Date)](https://star-history.com/#nao1215/sqluv&Date)

## Contacto
Si desea enviar comentarios como "encontrar un error" o "solicitud de características adicionales" al desarrollador, por favor use uno de los siguientes contactos.

- [GitHub Issue](https://github.com/nao1215/sqluv/issues)

## LICENCIA

[MIT License](../../LICENSE)