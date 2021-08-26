package rds_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/golang/mock/gomock"
	"github.com/krotscheck/go-rds-driver"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Connection(t *testing.T) {
	d := rds.NewDriver()
	dsn, err := TestConfig.ToDSN()
	ctx := context.Background()
	if err != nil {
		t.Fatal(err)
	}
	connector, err := d.OpenConnector(dsn)
	if err != nil {
		t.Fatal(err)
	}

	Convey("Connection", t, func() {
		c, err := connector.Connect(ctx)
		So(err, ShouldBeNil)
		connection, ok := c.(*rds.Connection)
		So(ok, ShouldBeTrue)

		Convey("Ping", func() {
			err := connection.Ping(ctx)
			So(err, ShouldBeNil)
		})

		Convey("Prepare", func() {
			stmt, err := connection.Prepare("select 1")
			So(err, ShouldBeNil)
			So(stmt, ShouldNotBeNil)
		})

		Convey("PrepareContext", func() {
			stmt, err := connection.PrepareContext(ctx, "select 1")
			So(err, ShouldBeNil)
			So(stmt, ShouldNotBeNil)
		})

		Convey("Close", func() {
			err := connection.Close()
			So(err, ShouldBeNil)
			err = connection.Close()
			So(err, ShouldNotBeNil)
		})

		Convey("Transactions", func() {
			Convey("Begin()", func() {
				tx, err := connection.Begin()
				So(err, ShouldBeNil)
				So(tx, ShouldNotBeNil)

				defer func() {
					So(tx.Rollback(), ShouldBeNil)
				}()

				dTx, ok := tx.(*rds.Tx)
				So(ok, ShouldBeTrue)
				So(dTx.TransactionID, ShouldNotBeEmpty)
			})

			Convey("Isolation Levels", func() {
				Convey("Unsupported", func() {
					Convey(sql.LevelLinearizable.String(), func() {
						tx, err := connection.BeginTx(ctx, driver.TxOptions{
							Isolation: driver.IsolationLevel(sql.LevelLinearizable),
							ReadOnly:  false,
						})
						So(err, ShouldNotBeNil)
						So(tx, ShouldBeNil)
					})
					Convey(sql.LevelSnapshot.String(), func() {
						tx, err := connection.BeginTx(ctx, driver.TxOptions{
							Isolation: driver.IsolationLevel(sql.LevelSnapshot),
							ReadOnly:  false,
						})
						So(err, ShouldNotBeNil)
						So(tx, ShouldBeNil)
					})
					Convey(sql.LevelWriteCommitted.String(), func() {
						tx, err := connection.BeginTx(ctx, driver.TxOptions{
							Isolation: driver.IsolationLevel(sql.LevelWriteCommitted),
							ReadOnly:  false,
						})
						So(err, ShouldNotBeNil)
						So(tx, ShouldBeNil)
					})
				})
				Convey("Supported", func() {
					Convey(sql.LevelSerializable.String(), func() {
						tx, err := connection.BeginTx(ctx, driver.TxOptions{
							Isolation: driver.IsolationLevel(sql.LevelSerializable),
							ReadOnly:  false,
						})
						So(err, ShouldBeNil)
						So(tx, ShouldNotBeNil)
						defer func() {
							So(tx.Rollback(), ShouldBeNil)
						}()
					})
					Convey(sql.LevelReadUncommitted.String(), func() {
						tx, err := connection.BeginTx(ctx, driver.TxOptions{
							Isolation: driver.IsolationLevel(sql.LevelReadUncommitted),
							ReadOnly:  false,
						})
						So(err, ShouldBeNil)
						So(tx, ShouldNotBeNil)
						defer func() {
							So(tx.Rollback(), ShouldBeNil)
						}()
					})
					Convey(sql.LevelReadCommitted.String(), func() {
						tx, err := connection.BeginTx(ctx, driver.TxOptions{
							Isolation: driver.IsolationLevel(sql.LevelReadCommitted),
							ReadOnly:  false,
						})
						So(err, ShouldBeNil)
						So(tx, ShouldNotBeNil)
						defer func() {
							So(tx.Rollback(), ShouldBeNil)
						}()
					})
					Convey(sql.LevelRepeatableRead.String(), func() {
						tx, err := connection.BeginTx(ctx, driver.TxOptions{
							Isolation: driver.IsolationLevel(sql.LevelRepeatableRead),
							ReadOnly:  false,
						})
						So(err, ShouldBeNil)
						So(tx, ShouldNotBeNil)
						defer func() {
							So(tx.Rollback(), ShouldBeNil)
						}()
					})
				})
			})

			Convey("Read Only On/Off", func() {
				tx, err := connection.BeginTx(ctx, driver.TxOptions{
					Isolation: driver.IsolationLevel(sql.LevelDefault),
					ReadOnly:  true,
				})
				So(err, ShouldBeNil)
				So(tx, ShouldNotBeNil)
				defer func() {
					So(tx.Rollback(), ShouldBeNil)
				}()
			})

			Convey("Commit()", func() {
				tx, err := connection.Begin()
				So(err, ShouldBeNil)
				So(tx, ShouldNotBeNil)

				defer func() {
					So(tx.Commit(), ShouldBeNil)
				}()

				dTx, ok := tx.(*rds.Tx)
				So(ok, ShouldBeTrue)
				So(dTx.TransactionID, ShouldNotBeEmpty)
			})
		})

		Convey("ResetSession", func() {

			Convey("Clean", func() {
				err := connection.ResetSession(ctx)
				So(err, ShouldBeNil)
			})

			Convey("Not clean", func() {
				_, err := connection.Begin()
				So(err, ShouldBeNil)
				err = connection.ResetSession(ctx)
				So(err, ShouldEqual, driver.ErrBadConn)
			})
		})

		Convey("IsValid", func() {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRDS := NewMockRDSDataServiceAPI(ctrl)

			Convey("Closed", func() {
				c := rds.NewConnection(ctx, mockRDS, "resourceARN", "secretARN", "database", rds.DialectMySQL)
				conn, ok := c.(*rds.Connection)
				So(ok, ShouldBeTrue)
				err := conn.Close()
				So(err, ShouldBeNil)
				So(conn.IsValid(), ShouldBeFalse)
			})

			Convey("Misconfigured", func() {
				c := rds.NewConnection(ctx, mockRDS, "", "secretARN", "database", rds.DialectMySQL)
				conn, ok := c.(*rds.Connection)
				So(ok, ShouldBeTrue)
				So(conn.IsValid(), ShouldBeFalse)

				c = rds.NewConnection(ctx, mockRDS, "resourceARN", "", "database", rds.DialectMySQL)
				conn, ok = c.(*rds.Connection)
				So(ok, ShouldBeTrue)
				So(conn.IsValid(), ShouldBeFalse)

				c = rds.NewConnection(ctx, mockRDS, "resourceARN", "secretARN", "", rds.DialectMySQL)
				conn, ok = c.(*rds.Connection)
				So(ok, ShouldBeTrue)
				So(conn.IsValid(), ShouldBeFalse)

				c = rds.NewConnection(ctx, mockRDS, "resourceARN", "secretARN", "database", rds.DialectMySQL)
				conn, ok = c.(*rds.Connection)
				So(ok, ShouldBeTrue)
				So(conn.IsValid(), ShouldBeTrue)
			})
		})

		Convey("QueryContext", func() {
			rows, err := connection.QueryContext(ctx, "select 1", []driver.NamedValue{})
			So(err, ShouldBeNil)
			So(len(rows.Columns()), ShouldEqual, 1)
			values := make([]driver.Value, len(rows.Columns()))
			err = rows.Next(values)
			So(err, ShouldBeNil)
			So(len(values), ShouldEqual, 1)
			So(values[0], ShouldEqual, 1)
		})

		Convey("ExecContext", func() {
			result, err := connection.ExecContext(ctx, "select 1", []driver.NamedValue{})
			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
			rows, err := result.RowsAffected()
			So(err, ShouldBeNil)
			So(rows, ShouldEqual, 0)
		})
	})
}
