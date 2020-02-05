module github.com/transcom/mymove

go 1.13

require (
	github.com/0xAX/notificator v0.0.0-20191016112426-3962a5ea8da1 // indirect
	github.com/99designs/aws-vault v4.5.1+incompatible
	github.com/99designs/keyring v1.1.3
	github.com/aws/aws-sdk-go v1.28.10
	github.com/cockroachdb/cockroach-go v0.0.0-20190925194419-606b3d062051 // indirect
	github.com/codegangsta/envy v0.0.0-20141216192214-4b78388c8ce4 // indirect
	github.com/codegangsta/gin v0.0.0-20171026143024-cafe2ce98974
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/disintegration/imaging v1.6.2
	github.com/dustin/go-humanize v1.0.0
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a
	github.com/fatih/color v1.7.0 // indirect
	github.com/felixge/httpsnoop v1.0.1
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df
	github.com/go-ini/ini v1.49.0 // indirect
	github.com/go-openapi/errors v0.19.3
	github.com/go-openapi/loads v0.19.4
	github.com/go-openapi/runtime v0.19.10
	github.com/go-openapi/spec v0.19.5
	github.com/go-openapi/strfmt v0.19.4
	github.com/go-openapi/swag v0.19.7
	github.com/go-openapi/validate v0.19.6
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-swagger/go-swagger v0.21.0
	github.com/gobuffalo/envy v1.9.0
	github.com/gobuffalo/fizz v1.9.5
	github.com/gobuffalo/flect v0.2.0
	github.com/gobuffalo/nulls v0.1.0 // indirect
	github.com/gobuffalo/packr v1.30.1
	github.com/gobuffalo/pop v4.13.1+incompatible
	github.com/gobuffalo/uuid v2.0.5+incompatible
	github.com/gobuffalo/validate v2.0.4+incompatible
	github.com/gocarina/gocsv v0.0.0-20190927101021-3ecffd272576
	github.com/gofrs/flock v0.7.1
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/gorilla/csrf v1.6.2
	github.com/imdario/mergo v0.3.8
	github.com/jackc/pgconn v1.1.0 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/jstemmer/go-junit-report v0.9.1
	github.com/jung-kurt/gofpdf v1.16.2
	github.com/leodido/go-urn v1.2.0
	github.com/lib/pq v1.3.0
	github.com/markbates/goth v1.61.1
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/mattn/go-shellwords v1.0.6 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/namsral/flag v1.7.4-pre
	github.com/pdfcpu/pdfcpu v0.2.5
	github.com/pkg/errors v0.9.1
	github.com/rickar/cal v1.0.1
	github.com/rogpeppe/go-internal v1.5.1 // indirect
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.2
	github.com/stretchr/objx v0.2.0
	github.com/stretchr/testify v1.4.0
	github.com/tcnksm/go-input v0.0.0-20180404061846-548a7d7a8ee8
	github.com/tealeg/xlsx v1.0.5
	github.com/transcom/nom v0.0.0-20190822223231-fc93620f238d
	github.com/vektra/mockery v0.0.0-20181123154057-e78b021dcbb5
	go.mozilla.org/pkcs7 v0.0.0-20181213175627-3cffc6fbfe83
	go.uber.org/zap v1.13.0
	goji.io v2.0.2+incompatible
	golang.org/x/crypto v0.0.0-20191028145041-f83a4685e152
	golang.org/x/net v0.0.0-20191028085509-fe3aa8a45271
	golang.org/x/text v0.3.2
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
	gopkg.in/urfave/cli.v1 v1.20.0 // indirect
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
