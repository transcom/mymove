package adminapi

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/iampostgres"
)

var adminDBURL = ""
var adminDB *pgx.Conn
var adminDBPassword = ""
var dbSSLRootCertFile = ""
var dbIAM = false

func NewAdminDB(v *viper.Viper) {
	// DREW DEBUG START
	db := v.GetString(cli.DbNameFlag)
	host := v.GetString(cli.DbHostFlag)
	port := v.GetUint16(cli.DbPortFlag)
	user := v.GetString(cli.DbUserFlag)
	adminDBPassword = v.GetString(cli.DbPasswordFlag)
	dbOptions := map[string]string{
		"sslmode": v.GetString(cli.DbSSLModeFlag),
	}

	s := "postgres://%s:%s@%s:%d/%s?sslmode=%s"
	adminDBURL = fmt.Sprintf(s, user, "", host, port, db, dbOptions["sslmode"])

	dbSSLRootCertFile = v.GetString(cli.DbSSLRootCertFlag)
	dbIAM = v.GetBool(cli.DbIamFlag)
}

func connectAdminDB(logger *zap.Logger) *pgx.Conn {
	connConfig, err := pgx.ParseConfig(adminDBURL)
	if err != nil {
		logger.Error("admin debug: Cannot parse dbURL",
			zap.String("dbURL", adminDBURL),
			zap.Error(err))
		return nil
	}
	if len(dbSSLRootCertFile) > 0 {
		caCert, cerr := os.ReadFile(dbSSLRootCertFile)
		if cerr != nil {
			logger.Error("admin debug: Cannot read ssl root cert",
				zap.String("DbSSLRootCertFlag", dbSSLRootCertFile),
				zap.Error(cerr))
			return nil
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		connConfig.TLSConfig.RootCAs = caCertPool
	}

	if dbIAM {
		connConfig.User = "ecs_user"
		token, perr := iampostgres.GetPgAdminPassword()
		if perr != nil {
			logger.Error("admin debug: cannot get pg admin password",
				zap.Error(perr))
			return nil
		}
		connConfig.Password = token
	} else {
		connConfig.User = "postgres"
		connConfig.Password = adminDBPassword
	}

	db, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		logger.Error("admin debug: pgx connect failed",
			zap.Error(err))
		return nil
	}

	return db
}

func NewAdminDebugMiddleware(logger *zap.Logger) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		key := os.Getenv("DB_HOST")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("X-Admin-Debug")
			if authHeader != key {
				logger.Error("admin debug: unauthorized user")
				http.Error(w, http.StatusText(403), http.StatusForbidden)
				return
			}
			inner.ServeHTTP(w, r)
		})
	}
}

func NewEnvHandler(logger *zap.Logger) http.HandlerFunc {
	type envInOut struct {
		Name string `json:"name"`
	}

	sendEnv := func(w http.ResponseWriter, req envInOut) {
		err := json.NewEncoder(w).Encode(req)
		if err != nil {
			logger.Error("admin debug: Error encoding env req", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req envInOut
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			req.Name = err.Error()
			sendEnv(w, req)
			return
		}
		v := os.Getenv(req.Name)
		req.Name = v
		sendEnv(w, req)
	}
}

func NewPGHandler(logger *zap.Logger) http.HandlerFunc {
	type pgCmd struct {
		Query string `json:"query"`
		Exec  string `json:"exec"`
	}

	type pgResp struct {
		Error        error      `json:"error"`
		RowsAffected int64      `json:"rowsAffected"`
		Columns      []string   `json:"columns"`
		Rows         [][]string `json:"rows"`
	}

	sendError := func(w http.ResponseWriter, pgerr error) {
		pr := pgResp{
			Error: pgerr,
		}
		err := json.NewEncoder(w).Encode(pr)
		if err != nil {
			logger.Error("admin debug: Error encoding error resp", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	sendResp := func(w http.ResponseWriter, resp pgResp) {
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			logger.Error("admin debug: Error encoding resp", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if adminDB == nil {
			adminDB = connectAdminDB(logger)
			if adminDB == nil {
				sendError(w, errors.New("No adminDB"))
				return
			}
		}
		err := adminDB.Ping(r.Context())
		if err != nil {
			adminDB = nil
			logger.Error("admin debug: cannot ping db", zap.Error(err))
			sendError(w, err)
		}
		var cmd pgCmd
		err = json.NewDecoder(r.Body).Decode(&cmd)
		if err != nil {
			sendError(w, err)
			return
		}
		if cmd.Exec != "" {
			result, err := adminDB.Exec(r.Context(), cmd.Exec)
			if err != nil {
				sendError(w, err)
				return
			}
			resp := pgResp{
				RowsAffected: result.RowsAffected(),
			}
			sendResp(w, resp)
			return
		}
		if cmd.Query != "" {
			rows, err := adminDB.Query(r.Context(), cmd.Query)
			if err != nil {
				sendError(w, err)
				return
			}
			defer rows.Close()
			desc := rows.FieldDescriptions()
			cols := []string{}
			for i := range desc {
				cols = append(cols, string(desc[i].Name))
			}
			resp := pgResp{
				Columns: cols,
				Rows:    [][]string{},
			}

			for rows.Next() {
				if rows.Err() != nil {
					sendError(w, err)
					return
				}
				vals, err := rows.Values()
				if err != nil {
					sendError(w, err)
					return
				}
				row := []string{}
				for i := range vals {
					row = append(row, fmt.Sprint(vals[i]))
				}
				resp.Rows = append(resp.Rows, row)
			}
			sendResp(w, resp)
			return
		}
	}
}
