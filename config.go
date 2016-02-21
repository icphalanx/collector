package collector

import (
	"crypto/tls"
	"database/sql"
)

type Config struct {
	CARemote  string
	TLSConfig *tls.Config
	DBConn    *sql.DB
}
