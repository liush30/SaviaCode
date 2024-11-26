package fabric

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/pkg/errors"
)

type MySQLUserStore struct {
	db *sql.DB
}

// NewMySQLUserStore 创建一个新的MySQLUserStore
func NewMySQLUserStore(dsn string) (*MySQLUserStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to open database connection")
	}
	// Optionally, you might want to test the connection here
	if err := db.Ping(); err != nil {
		return nil, errors.WithMessage(err, "failed to ping database")
	}
	return &MySQLUserStore{db: db}, nil
}

// Store 将用户存储在MySQL中
func (s *MySQLUserStore) Store(user *msp.UserData) error {
	query := `INSERT INTO t_users (USER_ID, msp_id, enrollment_cert) VALUES (?, ?, ?)
	          ON DUPLICATE KEY UPDATE enrollment_cert = VALUES(enrollment_cert)`
	_, err := s.db.Exec(query, user.ID, user.MSPID, user.EnrollmentCertificate)
	if err != nil {
		return errors.WithMessage(err, "failed to store user")
	}
	return nil
}

// Load 从MySQL中加载用户
func (s *MySQLUserStore) Load(id msp.IdentityIdentifier) (*msp.UserData, error) {
	log.Println("Load user from MySQL")
	query := `SELECT enrollment_cert FROM t_users WHERE USER_ID = ? AND msp_id = ?`
	row := s.db.QueryRow(query, id.ID, id.MSPID)

	var cert []byte
	if err := row.Scan(&cert); err != nil {
		if err == sql.ErrNoRows {
			return nil, msp.ErrUserNotFound
		}
		return nil, errors.WithMessage(err, "failed to load user")
	}

	return &msp.UserData{
		ID:                    id.ID,
		MSPID:                 id.MSPID,
		EnrollmentCertificate: cert,
	}, nil
}

// Delete 从MySQL中删除用户
func (s *MySQLUserStore) Delete(id msp.IdentityIdentifier) error {
	query := `DELETE FROM t_users WHERE USER_ID = ? AND msp_id = ?`
	_, err := s.db.Exec(query, id.ID, id.MSPID)
	if err != nil {
		return errors.WithMessage(err, "failed to delete user")
	}
	return nil
}
