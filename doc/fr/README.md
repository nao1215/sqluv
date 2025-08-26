![dbms_demo](../image/dbms_demo.gif)

[English](../../README.md) | [日本語](../ja/README.md) | [Русский](../ru/README.md) | [中文](../zh-cn/README.md) | [한국어](../ko/README.md) | [Español](../es/README.md)

**sqluv (sql + love)** est une interface utilisateur textuelle simple conçue pour interagir avec divers systèmes de gestion de bases de données relationnelles (SGBDR) ainsi qu'avec des fichiers CSV, TSV et LTSV. sqluv lit les fichiers CSV, TSV et LTSV depuis le stockage local, HTTPS et Amazon S3. sqluv décompresse automatiquement les fichiers compressés aux formats .gz, .bz2, .xz et .zst. Avec sqluv, l'exécution de requêtes SQL devient une expérience conviviale, permettant des connexions transparentes aux bases de données ou aux fichiers locaux avec facilité.

sqluv est une commande dérivée de [nao1215/sqly](https://github.com/nao1215/sqly). Son point de départ est de fournir une interface plus conviviale pour écrire du SQL par rapport à sqly.

>[!WARNING]
> sqluv est en développement. Utilisez sqluv en tant que **visualiseur**. N'exécutez pas UPDATE ou DELETE dans l'environnement de production.

## Fonctionnalités Principales
- Support Multi-SGBD : Connectez-vous et interagissez avec les systèmes de bases de données populaires comme MySQL, PostgreSQL, SQLite3 et SQL Server.
- Compatibilité des Fichiers : Lisez les données de fichiers aux formats CSV, TSV et LTSV depuis HTTPS, S3 et le stockage local.
- Support des Fichiers Compressés : Décompresse automatiquement les fichiers compressés aux formats .gz, .bz2, .xz et .zst.
- Historique des Requêtes : Sauvegardez et accédez à l'historique des requêtes SQL pour référence facile. La recherche floue est également disponible.
- Thèmes Personnalisables : sqluv supporte plusieurs thèmes de couleurs, permettant la personnalisation de l'interface selon les préférences de l'utilisateur.

## Comment installer
### Utiliser "go install"

```shell
go install github.com/nao1215/sqluv@latest
```

### Utiliser homebrew

```shell
brew install nao1215/tap/sqluv
```

## OS Supportés, Format de Fichier, Format Compressé, SGBD, version go

- Windows/macOS/Linux
- CSV/TSV/LTSV (file://, http://, https://)
- gz/bz2/xz/zst
- MySQL/PostgreSQL/SQLite3/SQL Server
- go1.24 ou ultérieur

## Comment utiliser

### Syntaxe

L'interface sqluv privilégie la facilité d'utilisation. Lors du lancement sans spécifier un chemin de fichier, les utilisateurs sont invités à saisir les détails de connexion de leur base de données. La configuration est sauvegardée, permettant des reconnexions faciles à l'avenir. Voici un bref aperçu des capacités :

```shell
sqluv [FILE_PATHS/HTTPS URL/S3 URL]  ※ Formats de fichiers supportés : CSV, TSV, LTSV
```

En exécutant cette commande avec les chemins de fichiers pertinents, les utilisateurs peuvent initier des interactions avec les fichiers.

### Se connecter au SGBD

Lorsque vous démarrez la commande sqluv sans spécifier un chemin de fichier, l'écran suivant apparaîtra.

![first_screen](../image/dbms_first.png)

Veuillez saisir les informations de connexion pour le SGBD auquel vous souhaitez vous connecter.

![dbms_connection](../image/dbms_info.png)

Si la connexion réussit, les informations de connexion à la base de données seront sauvegardées dans le fichier de configuration. La prochaine fois que vous démarrez la commande sqluv, vous pourrez sélectionner le SGBD auquel vous souhaitez vous connecter depuis la liste.

![dbms_list](../image/dbms_list.png)

![home_screen](../image/dbms_home.png)

### Exécuter une requête SQL

Pour exécuter une requête SQL, saisissez la requête SQL dans la zone de texte de requête et appuyez sur le bouton d'exécution ou `Ctrl + e`. Lorsque vous sélectionnez le nom de table dans la barre latérale et appuyez sur `Ctrl + e`, sqluv exécute la requête `SELECT * FROM ${TABLE_NAME} LIMIT 100`.

Pour rechercher un nom de table, appuyez sur la touche `/` dans la barre latérale. sqluv affichera le champ de recherche dans le pied de page. Si vous appuyez sur la touche `ESC`, le champ de recherche sera effacé.

![sql_query](../image/search_tables.png)

Pour afficher/masquer les colonnes, appuyez sur la touche `Space` dans la barre latérale. Si vous appuyez sur la touche `Space` dans la barre latérale, sqluv affiche le DDL de la table.

![ddl](../image/ddl_info.png)

## Historique des requêtes SQL

Si vous exécutez une requête SQL, l'historique sera sauvegardé dans `~/.config/sqluv/history.db`. Ainsi, vous pouvez consulter l'historique en appuyant sur le bouton historique.

![history_button](../image/history_button.png)

sqluv supporte la recherche floue. Vous pouvez rechercher dans l'historique en tapant le mot-clé. Si vous sélectionnez l'historique, la requête SQL sera copiée dans le champ de saisie.

![history_list](../image/sql_query_history.gif)

### Importer CSV/TSV/LTSV

Veuillez spécifier un chemin de fichier (ou une url) lors de l'exécution de la commande sqluv :

```shell
※ fichier sur http/https
sqluv https://raw.githubusercontent.com/nao1215/sqluv/refs/heads/main/testdata/actor.csv
 
※ fichier sur s3. le fichier est compressé au format gz.
sqluv s3://not-exist-s3-bucket/user.csv.gz
 
※ Fichiers multiples
sqluv https://raw.githubusercontent.com/nao1215/sqluv/refs/heads/main/testdata/actor.csv s3://not-exist-s3-bucket/user.tsv testdata/sample.ltsv
```

Le fichier sera chargé avant de lancer l'IUT. Quand sqluv importe csv/tsv/ltsv, sqluv vérifie l'extension du fichier et détermine le format du fichier. Si l'extension du fichier n'est pas csv/tsv/ltsv, sqluv affichera un message d'erreur. sqluv ne détecte pas automatiquement le format du fichier.

![sqluv_demo](../image/demo.gif)

### Sauvegarder le résultat dans un fichier

Vous pouvez sauvegarder le résultat dans un fichier en appuyant sur la touche `Ctrl + s`. sqluv vous demandera de saisir le chemin du fichier. Les formats de fichiers supportés sont CSV, TSV et LTSV.

![save_result](../image/file_save.png)

## Raccourcis clavier

| Touche | Description |
| --- | --- |
| Ctrl + d | Quitter |
| Ctrl + e | Exécuter la requête SQL |
| Ctrl + h | Afficher l'historique des requêtes SQL |
| Ctrl + c | Copier la requête SQL sélectionnée |
| Ctrl + v | Coller le texte copié |
| Ctrl + x | Couper le texte sélectionné |
| Ctrl + s | Sauvegarder le résultat dans un fichier |
| Ctrl + t | Changer le thème |
| /        | Rechercher le nom de la table (quand le focus est sur la barre latérale)|
| ESC      | Effacer le champ de recherche (quand le focus est sur la barre latérale)|
| Space    | Afficher/Masquer les Colonnes (quand le focus est sur la barre latérale)|
| Enter    | Afficher le DDL de la table (quand le focus est sur la barre latérale)|
| F1       | Focus sur la barre latérale |
| F2       | Focus sur la zone de texte de requête |
| F3       | Focus sur le tableau de résultats de requête |
| TAB | Passer au champ suivant |
| Shift + TAB | Passer au champ précédent |

## Thème de couleur

### Par défaut
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

## Outils Alternatifs

|Nom | Description |
|:----|:------------|
| [jorgerojas26/lazysql](https://github.com/jorgerojas26/lazysql) |Un outil de gestion de bases de données TUI multi-plateforme écrit en Go.|
| [vladbalmos/mitzasql](https://github.com/vladbalmos/mitzasql) | Client d'interface de ligne de commande / basé sur le texte MySQL |
| [TaKO8Ki/gobang](https://github.com/TaKO8Ki/gobang) | Un outil de gestion de bases de données TUI multi-plateforme écrit en Rust |

## Contribuer

Tout d'abord, merci de prendre le temps de contribuer ! Voir [CONTRIBUTING.md](../../CONTRIBUTING.md) pour plus d'informations. Les contributions ne sont pas seulement liées au développement. Par exemple, GitHub Star me motive à développer !

Si vous présentez sqluv sur les réseaux sociaux ou les blogs, plus d'utilisateurs le découvriront, et nous pourrons rassembler plus d'idées d'amélioration. Nous espérons que vous apprécierez utiliser sqluv !

[![Star History Chart](https://api.star-history.com/svg?repos=nao1215/sqluv&type=Date)](https://star-history.com/#nao1215/sqluv&Date)

## Contact
Si vous souhaitez envoyer des commentaires tels que "trouver un bug" ou "demande de fonctionnalités supplémentaires" au développeur, veuillez utiliser l'un des contacts suivants.

- [GitHub Issue](https://github.com/nao1215/sqluv/issues)

## LICENCE

[MIT License](../../LICENSE)