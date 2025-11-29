package rds_test

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	. "github.com/smartystreets/goconvey/convey"
)

// A table with all data types
const PostGreSQLCreateTableQuery = "CREATE TABLE IF NOT EXISTS all_types (" +
	"id SERIAL PRIMARY KEY," +

	// common postgres types
	"sql_boolean BOOLEAN," +
	"sql_char CHAR(1)," +
	"sql_varchar VARCHAR(100)," +
	"sql_text TEXT," +
	"sql_small_int SMALLINT," +
	"sql_medium_int INTEGER," +
	"sql_int BIGINT," +
	"sql_decimal DECIMAL," +
	"sql_numeric NUMERIC(5,2)," +
	"sql_real REAL," +
	"sql_byte BYTEA," +
	"sql_date DATE NULL," +
	"sql_time TIME NULL," +
	// "sql_timestampz TIMESTAMPZ NULL," +
	// "sql_interval INTERVAL NULL," +
	"sql_timestamp TIMESTAMP NULL)"

const PostgreSQLDropTableQuery = "DROP TABLE all_types;"

// TestPostgreSQLRow of data persisted to postgres
type TestPostgreSQLRow struct {
	ID        int32
	Boolean   bool
	Char      string
	Varchar   string
	Text      string
	SmallInt  int8
	MediumInt int32
	Bigint    int64
	Decimal   float64
	Numeric   float64
	Real      float64
	Byte      []byte
	Date      string
	Time      string
	Timestamp string
}

// NewTestPostgreSQLRow to insert into the database
func NewTestPostgreSQLRow() *TestPostgreSQLRow {
	bytes := make([]byte, 4)
	_, _ = rand.Read(bytes)

	t := &TestPostgreSQLRow{
		Boolean:   true,
		Char:      "C",
		Varchar:   "Varchar",
		Text:      "This is some text",
		SmallInt:  1,
		MediumInt: 2,
		Bigint:    3,
		Decimal:   2.22,
		Numeric:   23,
		Real:      33,
		Byte:      bytes,
		Date:      time.Now().Format("2006-01-02"),
		Time:      time.Now().Format("15:04:05"),
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	return t
}

// Scan applies the results of an SQL queries to this row
func (r *TestPostgreSQLRow) Scan(row *sql.Rows) error {
	return row.Scan(
		&r.ID,        // ID         int32
		&r.Boolean,   // Boolean   bool
		&r.Char,      // Char      string
		&r.Varchar,   // Varchar   string
		&r.Text,      // Text      string
		&r.SmallInt,  // SmallInt  int8
		&r.MediumInt, // MediumInt int32
		&r.Bigint,    // Bigint    int64
		&r.Decimal,   // Decimal   float64
		&r.Numeric,   // Numeric   float64
		&r.Real,      // Real      float64
		&r.Byte,      // Byte      []byte
		&r.Date,      // Date      string
		&r.Time,      // Time      string
		&r.Timestamp, // Timestamp string
	)
}

// Insert applies the results of an SQL queries to this row
func (r *TestPostgreSQLRow) Insert(db *sql.DB) (sql.Result, error) {
	params := []interface{}{
		r.Boolean,   // Boolean   bool
		r.Char,      // Char      string
		r.Varchar,   // Varchar   string
		r.Text,      // Text      string
		r.SmallInt,  // SmallInt  int8
		r.MediumInt, // MediumInt int32
		r.Bigint,    // Bigint    int64
		r.Decimal,   // Decimal   float64
		r.Numeric,   // Numeric   float64
		r.Real,      // Real      float64
		r.Byte,      // Byte      []byte
		r.Date,      // Date      string
		r.Time,      // Time      string
		r.Timestamp, // Timestamp string
	}
	query := "INSERT INTO all_types (" +
		"sql_boolean," +
		"sql_char,sql_varchar,sql_text," +
		"sql_small_int,sql_medium_int,sql_int," +
		"sql_decimal,sql_numeric,sql_real,sql_byte," +
		"sql_date,sql_time,sql_timestamp) " +
		"VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12::date,$13::time,$14::date)"
	return db.Exec(query, params...)
}

// Full suite of postgres tests, starting with research queries about which data types are supportyed
func Test_Postgresql(t *testing.T) {
	dsn := TestPostgresConfig.ToDSN()
	rdsDB, err := sql.Open("rds", dsn)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := rdsDB.Close()
		if err != nil {
			panic(err)
		}
	}()

	// Create the Local DB Instance
	dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", "root", "supersecret", "127.0.0.1", "5432", TestPostgresConfig.Database)
	localDB, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := localDB.Close()
		if err != nil {
			panic(err)
		}
	}()

	Convey("Postgresql", t, func() {
		result, err := rdsDB.Exec(PostGreSQLCreateTableQuery)
		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
		defer func() {
			result, err := rdsDB.Exec(PostgreSQLDropTableQuery)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
		}()

		result, err = localDB.Exec(PostGreSQLCreateTableQuery)
		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
		defer func() {
			result, err := localDB.Exec(PostgreSQLDropTableQuery)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
		}()

		Convey("Ping()", func() {
			row, err := rdsDB.Query("SELECT 1")
			So(err, ShouldBeNil)
			So(row, ShouldNotBeNil)

			row, err = localDB.Query("SELECT 1")
			So(err, ShouldBeNil)
			So(row, ShouldNotBeNil)
		})

		Convey("Multi Statement", func() {
			row, err := rdsDB.Query("SELECT 1; SELECT 2;")
			So(err, ShouldBeNil)
			So(row, ShouldNotBeNil)
		})

		Convey("Table", func() {

			for i := 0; i < 10; i++ {
				row := NewTestPostgreSQLRow()

				result, err = row.Insert(localDB)
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)
				count, err := result.RowsAffected()
				So(err, ShouldBeNil)
				So(count, ShouldEqual, 1)

				result, err := row.Insert(rdsDB)
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)
				count, err = result.RowsAffected()
				So(err, ShouldBeNil)
				So(count, ShouldEqual, 1)
			}

			Convey("Compare results", func() {
				localRows, err := localDB.Query("SELECT * FROM all_types")
				So(err, ShouldBeNil)
				localCols, err := localRows.ColumnTypes()
				So(err, ShouldBeNil)

				rdsRows, err := rdsDB.Query("SELECT * FROM all_types")
				So(err, ShouldBeNil)
				rdsCols, err := rdsRows.ColumnTypes()
				So(err, ShouldBeNil)

				So(len(rdsCols), ShouldEqual, len(localCols))

				for localRows.Next() {
					rdsRow := &TestPostgreSQLRow{}
					localRow := &TestPostgreSQLRow{}

					rdsRows.Next()

					err = localRow.Scan(localRows)
					So(err, ShouldBeNil)

					err := rdsRow.Scan(rdsRows)
					So(err, ShouldBeNil)

					So(rdsRow, ShouldResemble, localRow)
				}
			})
		})
	})
}
