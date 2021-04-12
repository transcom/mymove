module github.com/transcom/mymove

go 1.15

require (
	github.com/0xAX/notificator v0.0.0-20191016112426-3962a5ea8da1 // indirect
	github.com/99designs/aws-vault v4.5.1+incompatible
	github.com/99designs/keyring v1.1.6
	github.com/alexedwards/scs/redisstore v0.0.0-20200225172727-3308e1066830
	github.com/alexedwards/scs/v2 v2.4.0
	github.com/aws/aws-sdk-go v1.38.17
	github.com/benbjohnson/clock v1.1.0
	github.com/codegangsta/envy v0.0.0-20141216192214-4b78388c8ce4 // indirect
	github.com/codegangsta/gin v0.0.0-20171026143024-cafe2ce98974
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/disintegration/imaging v1.6.2
	github.com/dustin/go-humanize v1.0.0
	github.com/felixge/httpsnoop v1.0.1
	github.com/getlantern/deepcopy v0.0.0-20160317154340-7f45deb8130a
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df
	github.com/go-ini/ini v1.49.0 // indirect
	github.com/go-openapi/errors v0.20.0
	github.com/go-openapi/loads v0.20.2
	github.com/go-openapi/runtime v0.19.27
	github.com/go-openapi/spec v0.20.3
	github.com/go-openapi/strfmt v0.20.1
	github.com/go-openapi/swag v0.19.15
	github.com/go-openapi/validate v0.20.2
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.5.0
	github.com/gobuffalo/envy v1.9.0
	github.com/gobuffalo/fizz v1.13.0
	github.com/gobuffalo/flect v0.2.2
	github.com/gobuffalo/nulls v0.4.0 // indirect
	github.com/gobuffalo/pop/v5 v5.3.3
	github.com/gobuffalo/validate/v3 v3.3.0
	github.com/gocarina/gocsv v0.0.0-20190927101021-3ecffd272576
	github.com/gofrs/flock v0.7.3
	github.com/gofrs/uuid v3.4.0+incompatible
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-github/v31 v31.0.0
	github.com/gorilla/csrf v1.7.0
	github.com/imdario/mergo v0.3.12
	github.com/jackc/pgerrcode v0.0.0-20190803225404-afa3381909a6
	github.com/jessevdk/go-flags v1.5.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/jstemmer/go-junit-report v0.9.1
	github.com/jung-kurt/gofpdf v1.16.2
	github.com/leodido/go-urn v1.2.1
	github.com/lib/pq v1.10.0
	github.com/markbates/goth v1.67.1
	github.com/mattn/go-shellwords v1.0.6 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/namsral/flag v1.7.4-pre
	github.com/pdfcpu/pdfcpu v0.2.5
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.13.0
	github.com/rickar/cal v1.0.5
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v0.0.7
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/objx v0.3.0
	github.com/stretchr/testify v1.7.0
	github.com/tcnksm/go-input v0.0.0-20180404061846-548a7d7a8ee8
	github.com/tealeg/xlsx v1.0.5
	github.com/tiaguinho/gosoap v1.4.4
	github.com/vektra/mockery v0.0.0-20181123154057-e78b021dcbb5
	go.mozilla.org/pkcs7 v0.0.0-20181213175627-3cffc6fbfe83
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	goji.io v2.0.2+incompatible
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	golang.org/x/text v0.3.6
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
	gopkg.in/urfave/cli.v1 v1.20.0 // indirect
	pault.ag/go/pksigner v1.0.2
)

// transcom/sqlx v1.2.1 is just jmoiron's 1.2.0 with custom driver fixes
// This is a temporary solution till https://github.com/jmoiron/sqlx/pull/560
// is merged or a better solution is completed as mentioned in
// https://github.com/jmoiron/sqlx/pull/520
replace github.com/jmoiron/sqlx v1.2.0 => github.com/transcom/sqlx v1.2.1

// https://github.com/codegangsta/gin/issues/154#issuecomment-544391671
// This fixes an issue that was being caused due to urfave/cli v1.21.0
// being renamed.
replace gopkg.in/urfave/cli.v1 => github.com/urfave/cli v1.21.0

// Update to ignore compiler warnings on macOS catalina
// https://github.com/keybase/go-keychain/pull/55
// https://github.com/99designs/aws-vault/pull/427
replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
