package collector

import (
	"crypto/x509"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	pb "github.com/icphalanx/rpc"
)

type Ingestor struct {
	DBConn *sql.DB
}

type DBCertificate struct {
	ID           int
	SerialNumber []byte
	HostID       int
	Revoked      bool

	i *Ingestor
}

type DBHost struct {
	ID        int
	HumanName string

	i *Ingestor
}

type DBMetric struct {
	ID        int
	Name      string
	Type      int
	HumanName string

	i *Ingestor
}

type DBReport struct {
	ID        int
	HostID    int
	Timestamp time.Time

	i *Ingestor
}

type DBUser struct {
	ID        int
	Username  string
	Password  string
	SuperUser bool
	HostIDs   []int

	i *Ingestor
}

func (i *Ingestor) GetOrAddSSLCertificate(cert *x509.Certificate) (*DBCertificate, error) {
	serialNumber := cert.SerialNumber.Bytes()

	dbc := new(DBCertificate)
	dbc.i = i

	err := i.DBConn.QueryRow(
		"SELECT id, serial_no, host_id, revoked FROM certificates WHERE serial_no = $1", serialNumber,
	).Scan(&dbc.ID, &dbc.SerialNumber, &dbc.HostID, &dbc.Revoked)
	if err == sql.ErrNoRows {
		err = nil

		var hostID *int
		if cert.Subject.CommonName != ProvisioningCertCommonName {
			var host *DBHost
			host, err = i.GetOrAddHost(cert.Subject.CommonName)
			hostID = &host.ID
		}
		if err == nil {
			err = i.DBConn.QueryRow(
				"INSERT INTO certificates (serial_no, host_id) VALUES ($1, $2) RETURNING id, serial_no, host_id, revoked", serialNumber, hostID,
			).Scan(&dbc.ID, &dbc.SerialNumber, &dbc.HostID, &dbc.Revoked)
		}
	}

	return dbc, err
}

func (i *Ingestor) GetOrAddHost(hostname string) (*DBHost, error) {
	dbh := new(DBHost)
	dbh.i = i

	err := i.DBConn.QueryRow(
		"SELECT id, human_name FROM hosts WHERE human_name = $1", hostname,
	).Scan(&dbh.ID, &dbh.HumanName)
	if err == sql.ErrNoRows {
		err = i.DBConn.QueryRow(
			"INSERT INTO hosts (human_name) VALUES ($1) RETURNING id, human_name",
			hostname,
		).Scan(&dbh.ID, &dbh.HumanName)
	}

	return dbh, err
}

type DBLogLine struct {
	ID        int
	HostID    int
	Timestamp time.Time
	Reporter  string
	LogLine   string
	Tags      []string
}

func (i *Ingestor) GetLogsForHost(hostID int, limit int) ([]*DBLogLine, error) {
	pbl := make([]*DBLogLine, 0)

	dir := "ASC"
	if limit < 0 {
		limit = -limit
		dir = "DESC"
	}
	rows, err := i.DBConn.Query("SELECT ll.id, ll.host_id, ll.timestamp, ll.reporter, ll.log_line, ARRAY(SELECT tag_str FROM logtags lt WHERE lt.log_id=ll.id) FROM loglines ll ORDER BY id "+dir+" LIMIT $1", limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tagsStr string
		ll := new(DBLogLine)
		err := rows.Scan(&ll.ID, &ll.HostID, &ll.Timestamp, &ll.Reporter, &ll.LogLine, &tagsStr)
		if err != nil {
			return nil, err
		}

		pbl = append(pbl, ll)
	}

	return pbl, nil
}

func (i *Ingestor) GetHosts(hostIDs []int) ([]*DBHost, error) {
	dbhs := make([]*DBHost, 0)

	var rows *sql.Rows
	var err error
	if len(hostIDs) > 0 {
		x := ""
		for n := 1; n <= len(hostIDs); n++ {
			x += fmt.Sprintf("$%d, ", n)
		}
		x = x[:len(x)-2]
		rows, err = i.DBConn.Query("SELECT id, human_name FROM hosts WHERE id IN ($1)", x)
	} else {
		rows, err = i.DBConn.Query("SELECT id, human_name FROM hosts")
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		dbh := new(DBHost)
		if err := rows.Scan(&dbh.ID, &dbh.HumanName); err != nil {
			return nil, err
		}

		dbhs = append(dbhs, dbh)
	}

	return dbhs, nil
}

func (i *Ingestor) RevokeCertificate(cert *x509.Certificate) error {
	res, err := i.DBConn.Exec("UPDATE certificates SET revoked=true WHERE serial_no=$1", cert.SerialNumber.Bytes())
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (i *Ingestor) GetOrAddMetric(m *pb.Metric) (*DBMetric, error) {
	dbm := new(DBMetric)
	dbm.i = i

	err := i.DBConn.QueryRow("SELECT id, name, type, human_name FROM metrics WHERE name = $1 AND type = $2", m.Id, m.Type).Scan(&dbm.ID, &dbm.Name, &dbm.Type, &dbm.HumanName)
	if err == sql.ErrNoRows {
		err = i.DBConn.QueryRow(
			"INSERT INTO metrics (name, type, human_name) VALUES ($1, $2, $3) RETURNING id, name, type, human_name", m.Id, m.Type, m.HumanName).Scan(&dbm.ID, &dbm.Name, &dbm.Type, &dbm.HumanName)
	}

	return dbm, err
}

func (i *Ingestor) AddReport(dbh *DBHost) (*DBReport, error) {
	dbr := new(DBReport)
	dbr.i = i

	err := i.DBConn.QueryRow("INSERT INTO reports (host_id) VALUES ($1) RETURNING id, host_id, timestamp", dbh.ID).Scan(&dbr.ID, &dbr.HostID, &dbr.Timestamp)

	return dbr, err
}

func (i *Ingestor) AddMetricReading(dbr *DBReport, dbm *DBMetric, m *pb.Metric) error {
	var readingID int

	err := i.DBConn.QueryRow("INSERT INTO readings (report_id, metric_id) VALUES ($1, $2) RETURNING id", dbr.ID, dbm.ID).Scan(&readingID)
	if err != nil {
		return err
	}

	switch m.Type {
	case pb.Metric_UNCOUNTABLE:
		_, err = i.DBConn.Exec("INSERT INTO readings_uncountable (reading_id, value) VALUES ($1, $2)", readingID, m.GetIntValue())
		break
	case pb.Metric_STRINGARRAY:
		sa := m.GetStringArrayValue()
		if sa != nil {
			arr := make([]interface{}, len(sa.Value)+1)
			arr[0] = readingID

			sqls := "INSERT INTO readings_stringarray (reading_id, value) VALUES ($1, ARRAY["
			for n := 2; n < len(sa.Value)+2; n++ {
				if n != 2 {
					sqls += ", "
				}
				sqls += fmt.Sprintf("$%d", n)
				arr[n-1] = sa.Value[n-2]
			}
			sqls += "])"

			_, err = i.DBConn.Exec(sqls, arr...)
		}
		break
	}

	return err
}

func (i *Ingestor) IngestReport(repreq *pb.ReportRequest) error {
	dbh, err := i.GetOrAddHost(repreq.Host.HumanName)
	if err != nil {
		return err
	}

	dbr, err := i.AddReport(dbh)
	if err != nil {
		return err
	}

	for _, rep := range repreq.Reporters {
		for _, m := range rep.Metrics {
			dbm, err := i.GetOrAddMetric(m)
			if err != nil {
				return err
			}

			err = i.AddMetricReading(dbr, dbm, m)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *Ingestor) IngestLine(ln *pb.LogLine) error {
	ts := GoogleTimestampToTime(ln.Timestamp)

	dbh, err := i.GetOrAddHost(ln.Host)
	if err != nil {
		return err
	}

	var logID int

	err = i.DBConn.QueryRow("INSERT INTO loglines (host_id, timestamp, reporter, log_line) VALUES ($1, $2, $3, $4) RETURNING id", dbh.ID, ts, ln.Reporter, ln.Line).Scan(&logID)
	if err != nil {
		return err
	}

	for _, tag := range ln.Tags {
		i.DBConn.Exec("INSERT INTO logtags (log_id, tag_str) VALUES ($1, $2)", logID, tag)
	}

	return err
}

func (i *Ingestor) getUserBySQL(sql string, args ...interface{}) (*DBUser, error) {
	dbu := new(DBUser)
	dbu.i = i

	var hosts string

	err := i.DBConn.QueryRow(sql, args...).Scan(&dbu.ID, &dbu.Username, &dbu.Password, &dbu.SuperUser, &hosts)
	if err != nil {
		return nil, err
	}

	if hosts[0] != '{' || hosts[len(hosts)-1] != '}' {
		return nil, fmt.Errorf("`hosts` on user didn't match expected pattern")
	}

	hosts = hosts[1 : len(hosts)-1]
	if len(hosts) > 0 {
		hostsSlice := strings.Split(hosts, ",")
		dbu.HostIDs = make([]int, len(hostsSlice))
		for n, hstr := range hostsSlice {
			dbu.HostIDs[n], err = strconv.Atoi(hstr)
			if err != nil {
				return nil, err
			}
		}
	}

	return dbu, nil
}

func (i *Ingestor) GetUser(username string) (*DBUser, error) {
	return i.getUserBySQL("SELECT id, username, password, superuser, hosts FROM users WHERE username = $1", username)
}

func (i *Ingestor) GetUserByID(userID int) (*DBUser, error) {
	return i.getUserBySQL("SELECT id, username, password, superuser, hosts FROM users WHERE id = $1", userID)

}
