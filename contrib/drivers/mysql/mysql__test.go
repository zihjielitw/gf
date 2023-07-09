// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

const (
	TableSize        = 10
	TableName        = "user"
	TestSchema1      = "test1"
	TestSchema2      = "test2"
	TableNamePrefix1 = "gf_"
	TestDbUser       = "root"
	TestDbPass       = "12345678"
	CreateTime       = "2018-10-24 10:00:00"
)

var (
	db        gdb.DB
	db2       gdb.DB
	dbPrefix  gdb.DB
	dbInvalid gdb.DB
	ctx       = context.TODO()
)

func init() {
	nodeDefault := gdb.ConfigNode{
		Link: "mysql:root:12345678@tcp(127.0.0.1:3306)/?loc=Local&parseTime=true",
	}

	nodePrefix := gdb.ConfigNode{
		Link: "mysql:root:12345678@tcp(127.0.0.1:3306)/?loc=Local&parseTime=true",
	}
	nodePrefix.Prefix = TableNamePrefix1

	nodeInvalid := gdb.ConfigNode{
		Link: "mysql:root:12345678@tcp(127.0.0.1:3307)/?loc=Local&parseTime=true",
	}

	gdb.AddConfigNode("test", nodeDefault)
	gdb.AddConfigNode("prefix", nodePrefix)
	gdb.AddConfigNode("nodeinvalid", nodeInvalid)
	gdb.AddConfigNode(gdb.DefaultGroupName, nodeDefault)

	// Default db.
	if r, err := gdb.NewByGroup(); err != nil {
		gtest.Error(err)
	} else {
		db = r
	}
	schemaTemplate := "CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET UTF8"
	if _, err := db.Exec(ctx, fmt.Sprintf(schemaTemplate, TestSchema1)); err != nil {
		gtest.Error(err)
	}
	if _, err := db.Exec(ctx, fmt.Sprintf(schemaTemplate, TestSchema2)); err != nil {
		gtest.Error(err)
	}
	db = db.Schema(TestSchema1)
	db2 = db.Schema(TestSchema2)

	// Prefix db.
	if r, err := gdb.NewByGroup("prefix"); err != nil {
		gtest.Error(err)
	} else {
		dbPrefix = r
	}
	if _, err := dbPrefix.Exec(ctx, fmt.Sprintf(schemaTemplate, TestSchema1)); err != nil {
		gtest.Error(err)
	}
	if _, err := dbPrefix.Exec(ctx, fmt.Sprintf(schemaTemplate, TestSchema2)); err != nil {
		gtest.Error(err)
	}
	dbPrefix = dbPrefix.Schema(TestSchema1)

	// Invalid db.
	if r, err := gdb.NewByGroup("nodeinvalid"); err != nil {
		gtest.Error(err)
	} else {
		dbInvalid = r
	}
	dbInvalid = dbInvalid.Schema(TestSchema1)
}

func createTable(table ...string) string {
	return createTableWithDb(db, table...)
}

func createInitTable(table ...string) string {
	return createInitTableWithDb(db, table...)
}

func dropTable(table string) {
	dropTableWithDb(db, table)
}

func createTableWithDb(db gdb.DB, table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`%s_%d`, TableName, gtime.TimestampNano())
	}
	dropTableWithDb(db, name)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	    CREATE TABLE %s (
	        id          int(10) unsigned NOT NULL AUTO_INCREMENT,
	        passport    varchar(45) NULL,
	        password    char(32) NULL,
	        nickname    varchar(45) NULL,
	        create_time timestamp(6) NULL,
	        PRIMARY KEY (id)
	    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	    `, name,
	)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

func createInitTableWithDb(db gdb.DB, table ...string) (name string) {
	name = createTableWithDb(db, table...)
	array := garray.New(true)
	for i := 1; i <= TableSize; i++ {
		array.Append(g.Map{
			"id":          i,
			"passport":    fmt.Sprintf(`user_%d`, i),
			"password":    fmt.Sprintf(`pass_%d`, i),
			"nickname":    fmt.Sprintf(`name_%d`, i),
			"create_time": gtime.NewFromStr(CreateTime).String(),
		})
	}

	result, err := db.Insert(ctx, name, array.Slice())
	gtest.AssertNil(err)

	n, e := result.RowsAffected()
	gtest.Assert(e, nil)
	gtest.Assert(n, TableSize)
	return
}

func dropTableWithDb(db gdb.DB, table string) {
	if _, err := db.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)); err != nil {
		gtest.Error(err)
	}
}
