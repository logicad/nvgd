package help

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assetsf586eea2876d83b41022dafcc2e615003dfdcce3 = "# NVGD - Night Vision Goggles Daemon\n\nHTTP file server to help DevOps.\n\n[![Build Status](https://travis-ci.org/koron/nvgd.svg?branch=master)](https://travis-ci.org/koron/nvgd)\n[![Go Report Card](https://goreportcard.com/badge/github.com/koron/nvgd)](https://goreportcard.com/report/github.com/koron/nvgd)\n\nIndex:\n\n  * [How to use](#how-to-use)\n  * [Acceptable path](#acceptable-path)\n  * [Configuration file](#configuration-file)\n    * [Command Protocol Handlers](#command-protocol-handlers)\n    * [S3 Protocol Handlers](#config-s3-protocol-handlers)\n    * [Config DB Protocol Handler](#config-db-protocol-handler)\n    * [Default Filters](#default-filters)\n  * [Filters](#filters)\n  * [Prefix Aliases](#prefix-aliases)\n\n## How to use\n\nInstall:\n\n    $ go get github.com/koron/nvgd\n\nRun:\n\n    $ nvgd\n\nAccess:\n\n    $ curl http://127.0.0.1:9280/file:///var/log/message/httpd.log?tail=limit:25\n\nUpdate:\n\n    $ go get -u github.com/koron/nvgd\n\n\n## Acceptable path\n\nNvgd accepts path in these like format:\n\n    /{protocol}://{args/for/protocol}[?{filters}]\n\nNvgd supports these `protocol`s:\n\n  * `file` - `/file:///path/to/source`\n  * `command` - result of pre-defined commands\n  * `s3obj`\n    * get object: `/s3obj://bucket-name/key/to/object`\n  * `s3list`\n    * list common prefixes and objects: `/s3list://bucket-name/prefix/of/key`\n  * `db` - query pre-defined databases\n    * query `id` and `email` form users in `db_pq`:\n\n        ```\n        /db://db_pq/select id,email from users\n        ```\n\n    * support multiple databases:\n\n        ```\n        /db://db_pq2/foo/select id,email from users\n        /db://db_pq2/bar/select id,email from users\n        ```\n\n        This searchs from `foo` and `bar` databases.\n\n  * `config` - current nvgd's configuration\n\n      `/config://` or `/config/` (alias)\n\n  * `help` - show help (README.md) of nvgd.\n\n      `/help://` or `/help/` (alias)\n\n      It would be better that combining with `markdown` filter.\n\n      ```\n      /help/?markdown\n      ```\n\nSee also:\n\n  * [Filters](#filters)\n\n\n## Configuration file\n\nNvgd takes a configuration file in YAML format.  A file `nvgd.conf.yml` in\ncurrent directory or given file with `-c` option is loaded at start.\n\n`nvgd.conf.yml` consist from these parts:\n\n```yml\n# Listen IP address and port (OPTIONAL, default is \"127.0.0.1:9280\")\naddr: \"0.0.0.0:8080\"\n\n# Configuratio for protocols (OPTIONAL)\nprotocols:\n\n  # Pre-defined command handlers.\n  command:\n    ...\n\n  # AWS S3 protocol handler configuration (see other section, OPTIONAL).\n  s3:\n    ...\n\n  # DB protocol handler configuration (OPTIONAL, see below)\n  db:\n    ...\n\n# Default filters: pair of path prefix and filter description.\ndefualt_filters:\n  ...\n```\n\n### Commnad Protocol Handlers\n\nConfiguration of pre-defined command protocol handler maps a key to\ncorresponding command source.\n\nExample:\n\n```yml\ncommand:\n  \"df\": \"df -h\"\n  \"lstmp\": \"ls -l /tmp\"\n```\n\nThis enables two resources under `command` protocol.\n\n  * `/command://df`\n  * `/command://lstmp`\n\nYou could add filters of course, like: `/command://df?grep=re:foo`\n\n### Config S3 Protocol Handlers\n\nConfiguration of S3 protocor handlers consist from 2 parts: `default` and\n`buckets`.  `default` part cotains default configuration to connect S3.  And\n`buckets` part could contain configuration for each buckets specific.\n\n```yml\ns3:\n\n  # IANA timezone to show times (optional).  \"Asia/Tokyo\" for JST.\n  timezone: Asia/Tokyo\n\n  # default configuration to connect to S3 (REQUIRED for S3)\n  default:\n\n    # Access key ID for S3 (REQUIRED)\n    access_key_id: xxxxxxxxxxxxxxxxxxxx\n\n    # Secret access key (REQUIRED)\n    secret_access_key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\n\n    # Access point to connect (OPTIONAL, default is \"ap-northeast-1\")\n    region: ap-northeast-1\n\n    # Session token to connect (OPTIONAL, default is empty: not used)\n    session_token: xxxxxxx\n\n    # MaxKeys for S3 object listing. valid between 1 to 1000.\n    # (OPTIONAL, default is 1000)\n    max_keys: 10\n\n  # bucket specific configurations (OPTIONAL)\n  buckets:\n\n    # bucket name can be specified as key.\n    \"your_bucket_name1\":\n      # same properties with \"default\" can be placed at here.\n\n    # other buckets can be added here.\n    \"your_bucket_name2\":\n      ...\n```\n\n### Config DB Protocol Handler\n\nSample of configuration for DB protocol handler.\n\n```yml\ndb:\n  # key could be set favorite name for your database\n  db_pq:\n    # driver supports 'postgres' or 'mysql' for now\n    driver: 'postgres'\n    # name is driver-specific source name (DSN)\n    name: 'postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full'\n    # limit number of rows for a query (default: 100)\n    max_rows: 50\n\n  # sample of connecting to MySQL\n  db_mysql:\n    driver: 'mysql'\n    name:   'user:password@/dbname'\n```\n\nWith above configuration, you will be able to access those databases with below URLs or commands.\n\n  * `curl 'http://127.0.0.1:9280/db://db_pq/select%20email%20from%20users'`\n  * `curl 'http://127.0.0.1:9280/db://db_mysql/select%20email%20from%20users'`\n\n#### Mutiple Databases in an instance\n\nTo make DB protocol handler connect with multiple databases in an instance,\nthere are 3 steps to make it enable.\n\n1.  Add `multiple_database: true` property to DB configuration.\n2.  Add `{{.dbname}}` placeholder in value of `name`.\n3.  Access to URL `/db://db_pq/DBNAME/you query`.\n\n    DBNAME is used to expand `{{.dbname}}` in above.\n\nAs a result, your configuration would be like this:\n\n```yml\ndb:\n  db_pq:\n    driver: 'postgres'\n    name: 'postgres://pqgotest:password@localhost/{{.dbname}}?sslmode=verify-full'\n    multiple_database: true\n\n  # sample of connecting to MySQL\n  db_mysql:\n    driver: 'mysql'\n    name:   'user:password@/{{.dbname}}'\n    multiple_database: true\n```\n\n### Default Filters\n\nDefault filters provide a capability to apply implicit filters depending on\npath prefixes. See [Filters](#filters) for detail of filters.\n\nTo apply `tail` filter for under `/file:///var/log/` path:\n\n```yaml\ndefault_filters:\n  \"file:///var/log/\":\n    - \"tail\"\n```\n\nIf you want to show last 100 lines, change like this:\n\n```yaml\ndefault_filters:\n  \"file:///var/log/\":\n    - \"tail=limit:100\"\n```\n\nYou can specify different filters for paths.\n\n```yaml\ndefault_filters:\n  \"file:///var/log/\":\n    - \"tail\"\n  \"file:///tmp/\":\n    - \"head\"\n```\n\nDefault filters can be ignored separately by [all (pseudo) filter](#all-pseudo-filter).\n\nDefualt filters are ignored for directories source of file protocols.\n\n\n## Filters\n\nNvgd supports these filters:\n\n  * [Grep filter](#grep-filter)\n  * [Head filter](#head-filter)\n  * [Tail filter](#tail-filter)\n  * [Cut filter](#cut-filter)\n  * [Hash filter](#hash-filter)\n  * [LTSV filter](#ltsv-filter)\n  * [Index HTML filter](#index-html-filter)\n  * [HTML Table filter](#html-table-filter)\n  * [Text Table filter](#text-table-filter)\n  * [Markdown filter](#markdown-filter)\n  * [Refresh filter](#refresh-filter)\n  * [Download filter](#download-filter)\n  * [All (pseudo) filter](#all-pseudo-filter)\n\n### Filter Spec\n\nWhere `{filters}` is:\n\n    {filter}[&{filter}...]\n\nWhere `{filter}` is:\n\n    {filter_name}[={options}]\n\nWhere `{options}` is:\n\n    {option_name}:{value}[;{option_name}:{value}...]\n\nSee other section for detail of each filters.\n\nExample: get last 50 lines except empty lines.\n\n    /file:///var/log/messages?grep=re:^$;match:false&tail=limit:50\n\n### Grep filter\n\nOutput lines which matches against regular expression.\n\nAs default, matching is made for whole line.  But when valid option `field` is\ngiven, then matching is made for specified a field, which is splitted by\n`delim` character.\n\n`grep` command equivalent.\n\n  * filter\\_name: `grep`\n  * options\n    * `re` - regular expression used for match.\n    * `match` - output when match or not match.  default is true.\n    * `field` - a match target N'th field counted from 1.\n      default is none (whole line).\n    * `delim` - field delimiter string (default: TAB character).\n\n### Head filter\n\nOutput the first N lines.\n\n`head` command equivalent.\n\n  * filter\\_name: `head`\n  * options\n    * `start` - start line number for output.  begging 0.  default is 0.\n    * `limit` - line number for output.  default is 10.\n\n### Tail filter\n\nOutput the last N lines.\n\n`tail` command equivalent.\n\n  * filter\\_name: `tail`\n  * options\n    * `limit` - line number for output.  default is 10.\n\n### Cut filter\n\nOutput selected fields of lines.\n\n`cut` command equivalent.\n\n  * filter\\_name: `cut`\n  * options:\n    * `delim` - field delimiter string (default: TAB character).\n    * `list` - selected fields, combinable by comma `,`.\n      * `N` - N'th field counted from 1.\n      * `N-M` - from N'th, to M'th field (included).\n      * `N-` - from N'th field, to end of line.\n      * `N-` - from first, to N'th field.\n\n### Hash filter\n\nOutput hash value.\n\n  * filter\\_name: `hash`\n  * options:\n    * `algorithm` - one of `md5` (default), `sha1`, `sha256` or `sha512`\n    * `encoding` - one of `hex` (default), `base64` or `binary`\n\n### Count filter\n\nCount lines.\n\n  * filter\\_name: `count`\n  * options: (none)\n\n### LTSV filter\n\nOutput, match to value of specified label, and output selected labels.\n\n  * filter\\_name: `ltsv`\n  * options:\n    * `grep` - match parameter: `{label},{pattern}`\n    * `match` - output when match or not match.  default is true.\n    * `cut` - selected labels, combinable by comma `,`.\n\n### Index HTML filter\n\nConvert LTSV to Index HTML.\n(limited for s3list and files (dir) source for now)\n\n  * filter\\_name: `indexhtml`\n  * options: (none)\n\nExample: list objects in S3 bucket \"foo\" with Index HTML.\n\n    http://127.0.0.1:9280/s3list://foo/?indexhtml\n\nThis filter should be the last of filters.\n\n### HTML Table filter\n\nConvert LTSV to HTML table.\n\n  * filter\\_name: `htmltable`\n  * options: (none)\n\nExample: query id and email column from users table on mine database.\n\n    http://127.0.0.1:9280/db://mine/select%20id,email%20from%20users?htmltable\n\nThis filter should be the last of filters.\n\n### Text Table filter\n\nConvert LTSV to plain text table.\n\n  * filter\\_name: `texttable`\n  * options: (none)\n\nExample: query id and email column from users table on mine database.\n\n    http://127.0.0.1:9280/db://mine/select%20id,email%20from%20users?texttable\n\nAbove query generate this table.\n\n    +-----+-----------------------+\n    |  id |        email          |\n    +-----+-----------------------+\n    |    0|foo@example.com        |\n    |    1|bar@example.com        |\n    +-----+-----------------------+\n\nThis filter should be the last of filters.\n\n### Markdown filter\n\nConvert markdown text to HTML.\n\n  * filter\\_name: `markdown`\n  * options: (none)\n\nExample: show help in HTML.\n\n    http://127.0.0.1:9280/help://?markdown\n    http://127.0.0.1:9280/help/?markdown\n\n### Refresh filter\n\nAdd \"Refresh\" header with specified time (sec).\n\n  * filter\\_name: `refresh`\n  * options: interval seconds to refresh.  0 for disable.\n\nExample: Open below URL using WEB browser, it refresh in each 5 seconds\nautomatically.\n\n    http://127.0.0.1:9280/file:///var/log/messages?tail&refresh=5\n\n### Download filter\n\nAdd \"Content-Disposition: attachment\" header to the response.  It make the\nbrowser to download the resource instead of showing in it.\n\n  * filter\\_name: `download`\n  * options: (none)\n\nExample: download the file \"messages\" and would be saved as file.\n\n    http://127.0.0.1:9280/file:///var/log/messages?download\n\n\n### All (pseudo) filter\n\nIgnore [default filters](#default-filters)\n\n  * filter\\_name: `all`\n  * options: (none)\n\nExample: if specified some default filters for `file:///var/`, this ignore\nthose.\n\n    http://127.0.0.1:9280/file:///var/log/messages?all\n\n\n## Prefix Aliases\n\nnvgd supports prefix aliases to keep compatiblities with [koron/night][night].\nCurrently these aliases are registered.\n\n* `files/` -> `file:///`\n* `commands/` -> `command://`\n* `config/` -> `config://`\n* `help/` -> `help://`\n\nFor example this URL:\n\n    http://127.0.0.1:9280/files/var/log/messages\n\nworks same as below URL:\n\n    http://127.0.0.1:9280/file:///var/log/messages\n\n\n## References\n\n  * [koron/night][night] previous impl in NodeJS.\n\n[night]:https://github.com/koron/night\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"README.md"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ff,
		Mtime:    time.Unix(1490155950, 1490155950832968900),
		Data:     nil,
	}, "/README.md": &assets.File{
		Path:     "/README.md",
		FileMode: 0x1b6,
		Mtime:    time.Unix(1490155950, 1490155950836973900),
		Data:     []byte(_Assetsf586eea2876d83b41022dafcc2e615003dfdcce3),
	}}, "")
