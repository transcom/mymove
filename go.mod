module github.com/transcom/mymove

go 1.17

require (
	github.com/0xAX/notificator v0.0.0-20191016112426-3962a5ea8da1 // indirect
	github.com/99designs/aws-vault v4.5.1+incompatible
	github.com/99designs/keyring v1.1.6
	github.com/DATA-DOG/go-txdb v0.1.5
	github.com/XSAM/otelsql v0.10.0
	github.com/alexedwards/scs/redisstore v0.0.0-20200225172727-3308e1066830
	github.com/alexedwards/scs/v2 v2.5.0
	github.com/aws/aws-sdk-go v1.42.53
	github.com/benbjohnson/clock v1.3.0
	github.com/codegangsta/envy v0.0.0-20141216192214-4b78388c8ce4 // indirect
	github.com/codegangsta/gin v0.0.0-20171026143024-cafe2ce98974
	github.com/disintegration/imaging v1.6.2
	github.com/dustin/go-humanize v1.0.0
	github.com/felixge/httpsnoop v1.0.2
	github.com/getlantern/deepcopy v0.0.0-20160317154340-7f45deb8130a
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df
	github.com/go-ini/ini v1.49.0 // indirect
	github.com/go-openapi/errors v0.20.2
	github.com/go-openapi/loads v0.21.1
	github.com/go-openapi/runtime v0.23.0
	github.com/go-openapi/spec v0.20.4
	github.com/go-openapi/strfmt v0.21.2
	github.com/go-openapi/swag v0.21.1
	github.com/go-openapi/validate v0.20.3
	github.com/go-playground/locales v0.14.0
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator/v10 v10.10.0
	github.com/gobuffalo/envy v1.10.1
	github.com/gobuffalo/fizz v1.14.0
	github.com/gobuffalo/flect v0.2.4
	github.com/gobuffalo/nulls v0.4.1 // indirect
	github.com/gobuffalo/pop/v5 v5.3.4
	github.com/gobuffalo/validate/v3 v3.3.1
	github.com/gocarina/gocsv v0.0.0-20190927101021-3ecffd272576
	github.com/gofrs/uuid v4.2.0+incompatible
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-github/v31 v31.0.0
	github.com/gorilla/csrf v1.7.1
	github.com/gorilla/mux v1.8.0
	github.com/imdario/mergo v0.3.12
	github.com/jackc/pgerrcode v0.0.0-20201024163028-a0d42d470451
	github.com/jessevdk/go-flags v1.5.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/jstemmer/go-junit-report v0.9.1
	github.com/jung-kurt/gofpdf v1.16.2
	github.com/leodido/go-urn v1.2.1
	github.com/lib/pq v1.10.4
	github.com/luna-duclos/instrumentedsql v1.1.3
	github.com/markbates/goth v1.69.0
	github.com/mattn/go-shellwords v1.0.6 // indirect
	github.com/mitchellh/mapstructure v1.4.3
	github.com/namsral/flag v1.7.4-pre
	github.com/pdfcpu/pdfcpu v0.2.5
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.13.4
	github.com/pterm/pterm v0.12.36
	github.com/spf13/afero v1.8.1
	github.com/spf13/cobra v1.3.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.1
	github.com/stretchr/objx v0.3.0
	github.com/stretchr/testify v1.7.0
	github.com/tcnksm/go-input v0.0.0-20180404061846-548a7d7a8ee8
	github.com/tealeg/xlsx/v3 v3.2.4
	github.com/tiaguinho/gosoap v1.4.4
	github.com/vektra/mockery/v2 v2.9.4
	go.mozilla.org/pkcs7 v0.0.0-20181213175627-3cffc6fbfe83
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20211108221036-ceb1ce70b4fa
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	golang.org/x/sys v0.0.0-20211210111614-af8b64212486 // indirect
	golang.org/x/text v0.3.7
	golang.org/x/tools v0.1.8
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
	gopkg.in/urfave/cli.v1 v1.20.0 // indirect
	pault.ag/go/pksigner v1.0.2
)

require (
	github.com/golang-jwt/jwt/v4 v4.3.0
	github.com/rickar/cal/v2 v2.1.2
	github.com/trussworks/otelhttp v0.0.0-20211201183707-a6d8db57ed9b
	go.opentelemetry.io/contrib/detectors/aws/ecs v1.3.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.29.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.27.0
	go.opentelemetry.io/contrib/propagators/aws v1.3.0
	go.opentelemetry.io/otel v1.4.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.26.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.26.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.4.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.4.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.26.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.4.0
	go.opentelemetry.io/otel/metric v0.26.0
	go.opentelemetry.io/otel/sdk v1.4.0
	go.opentelemetry.io/otel/sdk/export/metric v0.26.0
	go.opentelemetry.io/otel/sdk/metric v0.26.0
	go.opentelemetry.io/otel/trace v1.4.0
	gotest.tools/gotestsum v1.7.0
)

require (
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/atomicgo/cursor v0.0.1 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/danieljoos/wincred v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dnephin/pflag v1.0.7 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dvsekhvalnov/jose2go v0.0.0-20200901110807-248326c1351b // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/frankban/quicktest v1.11.2 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-logr/logr v1.2.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/analysis v0.21.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gobuffalo/attrs v1.0.1 // indirect
	github.com/gobuffalo/genny/v2 v2.0.8 // indirect
	github.com/gobuffalo/github_flavored_markdown v1.1.1 // indirect
	github.com/gobuffalo/helpers v0.6.4 // indirect
	github.com/gobuffalo/logger v1.0.6 // indirect
	github.com/gobuffalo/packd v1.0.1 // indirect
	github.com/gobuffalo/packr/v2 v2.8.0 // indirect
	github.com/gobuffalo/plush/v4 v4.1.9 // indirect
	github.com/gobuffalo/tags/v3 v3.1.2 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/gookit/color v1.4.2 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hhrutter/lzw v0.0.0-20190827003112-58b82c5a41cc // indirect
	github.com/hhrutter/tiff v0.0.0-20190827003322-d08e2ad45835 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.10.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.1.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.8.1 // indirect
	github.com/jackc/pgx/v4 v4.13.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/keybase/go-keychain v0.0.0-20190712205309-48d3d31d256d // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/markbates/errx v1.1.0 // indirect
	github.com/markbates/oncer v1.0.0 // indirect
	github.com/markbates/safe v1.0.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/microcosm-cc/bluemonday v1.0.16 // indirect
	github.com/miekg/pkcs11 v1.0.3 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rogpeppe/fastuuid v1.2.0 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/rs/zerolog v1.18.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/shabbyrobe/xmlwriter v0.0.0-20200208144257-9fca06d00ffa // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/sourcegraph/annotate v0.0.0-20160123013949-f4cad6c6324d // indirect
	github.com/sourcegraph/syntaxhighlight v0.0.0-20170531221838-bd320f5d308e // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	go.mongodb.org/mongo-driver v1.7.5 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.4.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.26.0 // indirect
	go.opentelemetry.io/proto/otlp v0.12.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8 // indirect
	golang.org/x/mod v0.5.1 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa // indirect
	google.golang.org/grpc v1.44.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	pault.ag/go/fasc v0.0.0-20190505145209-c337c3c0bbf0 // indirect
	pault.ag/go/othername v0.0.0-20190316144542-859caba4369b // indirect
	pault.ag/go/piv v0.0.0-20190320181422-d9d61c70919c // indirect
)

// https://github.com/codegangsta/gin/issues/154#issuecomment-544391671
// This fixes an issue that was being caused due to urfave/cli v1.21.0
// being renamed.
replace gopkg.in/urfave/cli.v1 => github.com/urfave/cli v1.21.0

// Update to ignore compiler warnings on macOS catalina
// https://github.com/keybase/go-keychain/pull/55
// https://github.com/99designs/aws-vault/pull/427
replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
