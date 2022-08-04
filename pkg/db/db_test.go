package db

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var pem = flag.String("pem", "", "PEM key location")
var postgresPass = flag.String("postgres-pass", "", "Database password")

// 10.0.8.145 is the private IP of DNS: postgres.staging.darwinrevolution.com. This DNS is only available from within the VPC.
var host = flag.String("host", "10.0.8.145", "Private IP of the database. If using a hostname, the hostname must resolve from the host's PC")

func TestPostgresSSH(t *testing.T) {
	// to run this test: go test ./pkg/db -run TestPostgresSSH -pem "/path/to/bastion.pem"
	if *pem == "" {
		t.Skip()
	}

	key, err := os.ReadFile(*pem)

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	user := "darwin"
	database := "odin"
	port := "5432"
	bastion := "54.198.236.124"

	client, err := SSH("ec2-user", bastion, "22", key)

	_, err = Postgres(*host, port, user, *postgresPass, database, client)

	assert.NoError(t, err)
}
