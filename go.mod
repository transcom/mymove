module github.com/transcom/mymove

go 1.20

require (
	github.com/DATA-DOG/go-txdb v0.1.7
	github.com/XSAM/otelsql v0.23.0
	github.com/alexedwards/scs/redisstore v0.0.0-20221223131519-238b052508b6
	github.com/alexedwards/scs/v2 v2.5.1
	github.com/aws/aws-sdk-go-v2 v1.21.0
	github.com/aws/aws-sdk-go-v2/config v1.18.39
	github.com/aws/aws-sdk-go-v2/credentials v1.13.37
	github.com/aws/aws-sdk-go-v2/feature/rds/auth v1.2.19
	github.com/aws/aws-sdk-go-v2/service/cloudwatchevents v1.15.13
	github.com/aws/aws-sdk-go-v2/service/ecr v1.19.5
	github.com/aws/aws-sdk-go-v2/service/ecs v1.30.1
	github.com/aws/aws-sdk-go-v2/service/rds v1.54.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.38.5
	github.com/aws/aws-sdk-go-v2/service/ses v1.16.7
	github.com/aws/aws-sdk-go-v2/service/ssm v1.37.5
	github.com/aws/aws-sdk-go-v2/service/sts v1.21.5
	github.com/aws/smithy-go v1.14.2
	github.com/benbjohnson/clock v1.3.5
	github.com/codegangsta/gin v0.0.0-20211113050330-71f90109db02
	github.com/disintegration/imaging v1.6.2
	github.com/dustin/go-humanize v1.0.1
	github.com/felixge/httpsnoop v1.0.4
	github.com/gabriel-vasile/mimetype v1.4.2
	github.com/go-chi/chi/v5 v5.0.10
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df
	github.com/go-logr/zapr v1.3.0
	github.com/go-openapi/errors v0.21.0
	github.com/go-openapi/loads v0.21.5
	github.com/go-openapi/runtime v0.27.0
	github.com/go-openapi/spec v0.20.14
	github.com/go-openapi/strfmt v0.22.0
	github.com/go-openapi/swag v0.22.9
	github.com/go-openapi/validate v0.23.0
	github.com/go-playground/validator/v10 v10.15.4
	github.com/go-swagger/go-swagger v0.30.5
	github.com/gobuffalo/envy v1.10.2
	github.com/gobuffalo/fizz v1.14.4
	github.com/gobuffalo/flect v1.0.2
	github.com/gobuffalo/pop/v6 v6.1.1
	github.com/gobuffalo/validate/v3 v3.3.3
	github.com/gocarina/gocsv v0.0.0-20221216233619-1fea7ae8d380
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/gomodule/redigo v1.8.9
	github.com/google/go-github/v31 v31.0.0
	github.com/gorilla/csrf v1.7.1
	github.com/imdario/mergo v0.3.16
	github.com/jackc/pgerrcode v0.0.0-20220416144525-469b46aa5efa
	github.com/jessevdk/go-flags v1.5.0
	github.com/jinzhu/copier v0.4.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/jung-kurt/gofpdf v1.16.2
	github.com/lib/pq v1.10.9
	github.com/markbates/goth v1.79.0
	github.com/namsral/flag v1.7.4-pre
	github.com/pdfcpu/pdfcpu v0.6.0
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.13.6
	github.com/pterm/pterm v0.12.66
	github.com/rickar/cal/v2 v2.1.13
	github.com/spf13/afero v1.11.0
	github.com/spf13/cobra v1.8.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.18.2
	github.com/stretchr/testify v1.9.0
	github.com/tcnksm/go-input v0.0.0-20180404061846-548a7d7a8ee8
	github.com/tealeg/xlsx/v3 v3.3.0
	github.com/tiaguinho/gosoap v1.4.4
	github.com/vektra/mockery/v2 v2.33.2
	go.flipt.io/flipt/rpc/flipt v1.25.0
	go.flipt.io/flipt/sdk/go v0.5.0
	go.mozilla.org/pkcs7 v0.0.0-20210826202110-33d05740a352
	go.opentelemetry.io/contrib/detectors/aws/ecs v1.18.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.43.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.43.0
	go.opentelemetry.io/contrib/propagators/aws v1.18.0
	go.opentelemetry.io/otel v1.18.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.40.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.18.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.18.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.40.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.17.0
	go.opentelemetry.io/otel/metric v1.18.0
	go.opentelemetry.io/otel/sdk v1.18.0
	go.opentelemetry.io/otel/sdk/metric v0.40.0
	go.opentelemetry.io/otel/trace v1.18.0
	go.uber.org/zap v1.26.0
	golang.org/x/crypto v0.21.0
	golang.org/x/net v0.21.0
	golang.org/x/oauth2 v0.17.0
	golang.org/x/text v0.14.0
	golang.org/x/tools v0.13.0
	google.golang.org/grpc v1.62.0
	gopkg.in/dnaeon/go-vcr.v3 v3.1.2
	gotest.tools/gotestsum v1.10.1
	pault.ag/go/pksigner v1.0.2
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx v1.2.29 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/patrickmn/go-cache v0.0.0-20180815053127-5633e0862627 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
)

require (
	atomicgo.dev/cursor v0.2.0 // indirect
	atomicgo.dev/keyboard v0.2.9 // indirect
	atomicgo.dev/schedule v0.1.0 // indirect
	github.com/0xAX/notificator v0.0.0-20220220101646-ee9b8921e557 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/Masterminds/sprig/v3 v3.2.3 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.13 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.41 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.35 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.42 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.1.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.36 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.35 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.15.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.13.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.15.6 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/brunoscheufler/aws-ecs-metadata-go v0.0.0-20221221133751-67e37ae746cd // indirect
	github.com/chigopher/pathlib v1.0.0 // indirect
	github.com/codegangsta/envy v0.0.0-20141216192214-4b78388c8ce4 // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dnephin/pflag v1.0.7 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/frankban/quicktest v1.14.6 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/analysis v0.22.2 // indirect
	github.com/go-openapi/inflect v0.19.0 // indirect
	github.com/go-openapi/jsonpointer v0.20.2 // indirect
	github.com/go-openapi/jsonreference v0.20.4 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/gobuffalo/attrs v1.0.3 // indirect
	github.com/gobuffalo/genny/v2 v2.1.0 // indirect
	github.com/gobuffalo/github_flavored_markdown v1.1.4 // indirect
	github.com/gobuffalo/helpers v0.6.7 // indirect
	github.com/gobuffalo/logger v1.0.7 // indirect
	github.com/gobuffalo/nulls v0.4.2 // indirect
	github.com/gobuffalo/packd v1.0.2 // indirect
	github.com/gobuffalo/plush/v4 v4.1.18 // indirect
	github.com/gobuffalo/tags/v3 v3.1.4 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.17.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hhrutter/lzw v1.0.0 // indirect
	github.com/hhrutter/tiff v1.0.1 // indirect
	github.com/huandu/xstrings v1.3.3 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.2 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/pgx/v4 v4.18.1 // indirect
	github.com/jarcoal/httpmock v1.3.0
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/lithammer/fuzzysearch v1.1.8 // indirect
	github.com/luna-duclos/instrumentedsql v1.1.3 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mattn/go-shellwords v1.0.12 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/microcosm-cc/bluemonday v1.0.23 // indirect
	github.com/miekg/pkcs11 v1.0.3 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/okta/okta-jwt-verifier-golang v1.3.1
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/peterbourgon/diskv/v3 v3.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/rogpeppe/fastuuid v1.2.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/rs/zerolog v1.29.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sergi/go-diff v1.3.1 // indirect
	github.com/shabbyrobe/xmlwriter v0.0.0-20200208144257-9fca06d00ffa // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/sourcegraph/annotate v0.0.0-20160123013949-f4cad6c6324d // indirect
	github.com/sourcegraph/syntaxhighlight v0.0.0-20170531221838-bd320f5d308e // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/toqueteos/webbrowser v1.2.0 // indirect
	github.com/urfave/cli v1.22.10 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.flipt.io/flipt/errors v1.19.3 // indirect
	go.mongodb.org/mongo-driver v1.13.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.40.0 // indirect
	go.opentelemetry.io/proto/otlp v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9
	golang.org/x/image v0.12.0 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/term v0.18.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	pault.ag/go/fasc v0.0.0-20190505145209-c337c3c0bbf0 // indirect
	pault.ag/go/othername v0.0.0-20190316144542-859caba4369b // indirect
	pault.ag/go/piv v0.0.0-20190320181422-d9d61c70919c // indirect
)
