package rds_test

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

// A table with all data types
const MySQLCreateTableQuery = "CREATE TABLE IF NOT EXISTS `all_types` (" +
	"`id` MEDIUMINT NOT NULL AUTO_INCREMENT," +

	// common mysql types
	"`sql_tiny_int` TINYINT," +
	"`sql_small_int` SMALLINT," +
	"`sql_medium_int` MEDIUMINT," +
	"`sql_int` bigint," +
	"`sql_big_int` BIGINT," +
	"`sql_decimal` DECIMAL(5,2)," +
	"`sql_float` double," +
	"`sql_double` DOUBLE," +
	//"`sql_bit` BIT," +
	"`sql_boolean` TINYINT(1)," +
	"`sql_char` CHAR," +
	"`sql_varchar` VARCHAR(100)," +
	"`sql_binary` BINARY," +
	"`sql_varbinary` VARBINARY(100)," +
	"`sql_tinyblob` TINYBLOB," +
	"`sql_blob` BLOB," +
	"`sql_mediumblob` MEDIUMBLOB," +
	"`sql_longblob` LONGBLOB," +
	"`sql_tinytext` TINYTEXT," +
	"`sql_text` TEXT," +
	"`sql_mediumtext` MEDIUMTEXT," +
	"`sql_enum` ENUM('ONE', 'TWO')," +
	"`sql_set` SET('ONE', 'TWO')," +
	"`sql_date` DATE," +
	"`sql_time` TIME," +
	"`sql_datetime` DATETIME," +
	"`sql_timestamp` TIMESTAMP," +
	"`sql_year` YEAR," +

	// types as automapped by gorm
	"`string` varchar(256)," +
	"`bytes` longblob," +
	"`byte` tinyint unsigned," +
	"`int8` tinyint," +
	"`int16` smallint," +
	"`int32` int," +
	"`int64` bigint," +
	"`uint` bigint unsigned," +
	"`uint8` tinyint unsigned," +
	"`uint16` smallint unsigned," +
	"`uint32` int unsigned," +
	"`uint64` bigint unsigned," +
	"`float32` float," +
	"`float64` double," +
	"PRIMARY KEY (`id`))"

const MySQLDropTableQuery = "DROP TABLE IF EXISTS `all_types`;"

// TestMySQLRow of data persisted to mysql
type TestMySQLRow struct {
	ID         int32
	TinyInt    int8
	SmallInt   int8
	MediumInt  int8
	Int        int8
	BigInt     int8
	Decimal    float64
	Float      float64
	Double     float64
	//Bit        []uint8
	Boolean    bool
	Char       string
	Varchar    string
	Binary     bool
	Varbinary  []byte
	Tinyblob   []byte
	Blob       []byte
	Mediumblob []byte
	Longblob   []byte
	Tinytext   string
	Text       string
	Mediumtext string
	Enum       string
	Set        string
	Date       string
	Time       string
	Datetime   string
	Timestamp  string
	Year       string
	String     string
	Bytes      []byte
	Byte       byte
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint       uint
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Float32    float32
	Float64    float64
}

// NewTestMySQLRow to insert into the database
func NewTestMySQLRow() *TestMySQLRow {
	bytes := make([]byte, 4)
	_, _ = rand.Read(bytes)

	t := &TestMySQLRow{
		TinyInt:    1,
		SmallInt:   2,
		MediumInt:  3,
		Int:        4,
		BigInt:     5,
		Decimal:    5.11,
		Float:      5.1111,
		Double:     1234.5678,
		//Bit:        []uint8{1},
		Boolean:    true,
		Char:       "1",
		Varchar:    "varchar",
		Binary:     true,
		Varbinary:  bytes,
		Tinyblob:   bytes,
		Blob:       bytes,
		Mediumblob: bytes,
		Longblob:   bytes,
		Tinytext:   "tiny",
		Text:       "Text",
		Mediumtext: "Mediumtext",
		Enum:       "ONE",
		Set:        "TWO",
		Date:       time.Now().Format("2006-01-02"),
		Time:       time.Now().Format("2006-01-02 15:04:05"),
		Datetime:   time.Now().Format("2006-01-02 15:04:05"),
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		Year:       time.Now().Format("2006"),
		String:     "This is a string",
		Bytes:      bytes,
		Byte:       bytes[0],
		Int8:       1,
		Int16:      2,
		Int32:      3,
		Int64:      4,
		Uint:       4,
		Uint8:      5,
		Uint16:     6,
		Uint32:     7,
		Uint64:     8,
		Float32:    1234.5678,
		Float64:    1234.5678,
	}

	return t
}

// Scan applies the results of an SQL query to this row
func (r *TestMySQLRow) Scan(row *sql.Rows) error {
	return row.Scan(
		&r.ID,         //ID         int32
		&r.TinyInt,    //TinyInt    int8
		&r.SmallInt,   //SmallInt   int8
		&r.Mediumtext, //MediumInt  int8
		&r.Int,        //Int        int8
		&r.BigInt,     //BigInt     int8
		&r.Decimal,    //Decimal    float64
		&r.Float,      //Float      float64
		&r.Double,     //Double     float64
		//&r.Bit,        //Bit        []uint8
		&r.Boolean,    //Boolean    tinyint(1)
		&r.Char,       //Char       string
		&r.Varchar,    //Varchar    string
		&r.Binary,     //Binary     bool
		&r.Varbinary,  //Varbinary  []byte
		&r.Tinyblob,   //Tinyblob   []byte
		&r.Blob,       //Blob       []byte
		&r.Mediumblob, //Mediumblob []byte
		&r.Longblob,   //Longblob   []byte
		&r.Tinytext,   //Tinytext   string
		&r.Text,       //Text       string
		&r.Mediumtext, //Mediumtext string
		&r.Enum,       //Enum       string
		&r.Set,        //Set        string
		&r.Date,       //Date       time.Time
		&r.Time,       //Time       time.Time
		&r.Datetime,   //Datetime   time.Time
		&r.Timestamp,  //Timestamp  time.Time
		&r.Year,       //Year       time.Time
		&r.String,     //String     string
		&r.Bytes,      //Bytes      []byte
		&r.Byte,       //Byte       byte
		&r.Int8,       //Int8       int8
		&r.Int16,      //Int16      int16
		&r.Int32,      //Int32      int32
		&r.Int64,      //Int64      int64
		&r.Uint,       //Uint       uint
		&r.Uint8,      //Uint8      uint8
		&r.Uint16,     //Uint16     uint16
		&r.Uint32,     //Uint32     uint32
		&r.Uint64,     //Uint64     uint64
		&r.Float32,    //Float32    float32
		&r.Float64,    //Float64    float64
	)
}

// Insert applies the results of an SQL query to this row
func (r *TestMySQLRow) Insert(db *sql.DB) (sql.Result, error) {
	params := []interface{}{
		r.TinyInt, r.SmallInt, r.MediumInt, r.Int, r.BigInt,
		r.Decimal, r.Float, r.Double,
		//r.Bit,
		r.Boolean,
		r.Char, r.Varchar,
		r.Binary, r.Varbinary,
		r.Tinyblob, r.Blob, r.Mediumblob, r.Longblob,
		r.Tinytext, r.Text, r.Mediumtext,
		r.Enum,
		r.Set,
		r.Date, r.Time, r.Datetime, r.Timestamp, r.Year,
		r.String,
		r.Bytes,
		r.Byte,
		r.Int8,
		r.Int16,
		r.Int32,
		r.Int64,
		r.Uint,
		r.Uint8,
		r.Uint16,
		r.Uint32,
		r.Uint64,
		r.Float32,
		r.Float64,
	}
	query := "INSERT INTO `all_types` SET" +
		"`sql_tiny_int` = ?,`sql_small_int` = ?,`sql_medium_int` = ?,`sql_int` = ?,`sql_big_int` = ?," +
		"`sql_decimal` = ?,`sql_float` = ?,`sql_double` = ?," +
		//"`sql_bit` = ?," +
		"`sql_boolean` = ?," +
		"`sql_char` = ?,`sql_varchar` = ?," +
		"`sql_binary` = ?,`sql_varbinary` = ?," +
		"`sql_tinyblob` = ?,`sql_blob` = ?,`sql_mediumblob` = ?,`sql_longblob` = ?," +
		"`sql_tinytext` = ?,`sql_text` = ?,`sql_mediumtext` = ?," +
		"`sql_enum` = ?," +
		"`sql_set` = ?," +
		"`sql_date` = ?,`sql_time` = ?,`sql_datetime` = ?,`sql_timestamp` = ?,`sql_year` = ?," +

		// types as automapped by gorm
		"`string` = ?," +
		"`bytes` = ?," +
		"`byte` = ?," +
		"`int8` = ?," +
		"`int16` = ?," +
		"`int32` = ?," +
		"`int64` = ?," +
		"`uint` = ?," +
		"`uint8` = ?," +
		"`uint16` = ?," +
		"`uint32` = ?," +
		"`uint64` = ?," +
		"`float32` = ?," +
		"`float64` = ?"
	return db.Exec(query, params...)
}

// Full suite of mysql tests, starting with research queries about which data types are supportyed
func Test_Mysql(t *testing.T) {
	// Create the RDS DB Instance
	dsn, err := TestMysqlConfig.ToDSN()
	if err != nil {
		panic(err)
	}
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

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "root", "supersecret", "127.0.0.1", "3306", TestMysqlConfig.Database)
	// Create the Local DB Instance
	localDB, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := localDB.Close()
		if err != nil {
			panic(err)
		}
	}()

	// Clean
	if _, err = localDB.Exec(MySQLDropTableQuery); err != nil {
		t.Fatal(err)
	}
	if _, err = rdsDB.Exec(MySQLDropTableQuery); err != nil {
		t.Fatal(err)
	}

	Convey("Mysql", t, func() {
		result, err := rdsDB.Exec(MySQLCreateTableQuery)
		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
		defer func() {
			result, err := rdsDB.Exec(MySQLDropTableQuery)
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
		}()

		result, err = localDB.Exec(MySQLCreateTableQuery)
		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
		defer func() {
			result, err := localDB.Exec(MySQLDropTableQuery)
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

		Convey("Table", func() {

			for i := 0; i < 10; i++ {
				row := NewTestMySQLRow()

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

				for rdsRows.Next() {
					rdsRow := &TestMySQLRow{}
					localRow := &TestMySQLRow{}

					localRows.Next()

					err := rdsRow.Scan(rdsRows)
					So(err, ShouldBeNil)

					err = localRow.Scan(localRows)
					So(err, ShouldBeNil)

					So(rdsRow, ShouldResemble, localRow)
				}
			})
		})
	})
}
