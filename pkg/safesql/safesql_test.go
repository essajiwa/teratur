package safesql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

const (
	driverName = "postgres"
)

func TestOpenMaster(t *testing.T) {
	master, err := OpenMasterDB("sqlite3", "file::memory:?mode=memory&cache=shared")
	require.NoError(t, err)

	err = master.Close()
	require.NoError(t, err)
}

func TestOpenMasterFailed(t *testing.T) {
	_, err := OpenMasterDB("yoursql", "file::memory:?mode=memory&cache=shared")
	require.Error(t, err)
}

func TestNewMaster(t *testing.T) {
	const (
		prepQuery = "insert into affiliate"
	)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectPrepare(prepQuery).ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))

	master := NewMasterDB(db, driverName)

	_, err = master.PreparexContext(context.Background(), prepQuery)
	require.NoError(t, err)
}

func TestOpenSlave(t *testing.T) {
	slave, err := OpenSlaveDB("sqlite3", "file::memory:?mode=memory&cache=shared")
	require.NoError(t, err)

	err = slave.Close()
	require.NoError(t, err)
}

func TestOpenSlaveFailed(t *testing.T) {
	_, err := OpenSlaveDB("yoursql", "file::memory:?mode=memory&cache=shared")
	require.Error(t, err)
}

func TestNewSlave(t *testing.T) {
	const (
		prepQuery = "select id from affiliate"
	)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectPrepare(prepQuery).ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{
		"id",
	}).AddRow(1))

	slave := NewSlaveDB(db, driverName)

	_, err = slave.PreparexContext(context.Background(), prepQuery)
	require.NoError(t, err)
}
