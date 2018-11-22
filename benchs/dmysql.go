package benchs

import (
	"fmt"
	"context"
	"time"
	"git.xiaojukeji.com/gobiz/dmysql"
	"git.xiaojukeji.com/gobiz/dmysql/builder"
)



var mgr *dmysql.Manager
var conn *dmysql.MySQL
var ctx context.Context

func init() {
	st := NewSuite("dmysql")
	tableName = "model"
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, DmysqlInsert)
		//st.AddBenchmark("MultiInsert 100 row", 500*ORM_MULTI, DmysqlMultiInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, DmysqlUpdate)
		st.AddBenchmark("Read", 2000*ORM_MULTI, DmysqlRead)
		st.AddBenchmark("MultiRead limit 100", 2000*ORM_MULTI, DmysqlReadSlice)

		var err error
		mgr, err = dmysql.New(ORM_HOST, ORM_USER, ORM_PASSWD, ORM_DB, ORM_CHARSET,
			dmysql.WithMaxConnSize(ORM_MAX_CONN), // 最大连接数
			//WithDebug(true), // 启用debug
			dmysql.WithDialTimeout(time.Second*1), // 连接超时时间
			dmysql.WithReadTimeout(time.Second*2), // 读超时
			dmysql.WithWriteTimeout(time.Second*2), // 写超时
			dmysql.WithAutoCommit(true), // 启用自动提交
			dmysql.WithPoolSize(4))// 连接池大小

		conn, err = mgr.Get() // 获取一个连接
		ctx = context.Background()
		if err != nil {
			panic(err)
			defer mgr.Put(conn) // 将连接归还到连接池
		}
	}
}

func DmysqlInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})

	for i := 1; i <= b.N; i++ {
		m.Id = 0
		modelMap := GetModelMap(m)
		_, err := conn.Insert(ctx, tableName, modelMap)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func DmysqlMultiInsertMulti(b *B) {
	var m *Model
	var ms []map[string]interface{}
	wrapExecute(b, func() {
		initDB()
		m = NewModel()

		ms = make([]map[string]interface{}, 0, 100)
		for i := 0; i < 100; i++ {
			ms = append(ms, GetModelMap(m))
		}
	})
	for i := 0; i < b.N; i++ {
		cond, vals, err := builder.BuildInsert(tableName, ms)
		err = conn.Execute(ctx, cond, vals...)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func DmysqlUpdate(b *B) {
	for i := 1; i <= b.N; i++ {
		where := map[string]interface{} {
			"id" : i,
		}
		update := map[string]interface{}{
			"age" : i,
		}
		cond,vals,err := builder.BuildUpdate(tableName, where, update)
		err = conn.Execute(ctx, cond, vals...)
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func DmysqlRead(b *B) {
	selectFields := []string{"id", "name", "title", "fax", "web","age", "`right`", "`counter`"}
	for i := 1; i <= b.N; i++ {
		where := map[string]interface{}{
			"id" : i,
		}
		cond,vals,err := builder.BuildSelect(tableName, where, selectFields)
		err = conn.Query(ctx, cond, vals...)
		rows, err := conn.FetchAllMap(ctx)
		if err != nil || rows == nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func DmysqlReadSlice(b *B) {

	selectFields := []string{"id", "name", "title", "fax", "web","age", "`right`", "`counter`"}
	for i := 1; i <= b.N; i++ {
		where := map[string]interface{}{
			"id >": 1,
			"_limit" : []uint{0,100},
		}
		cond,vals,err := builder.BuildSelect(tableName, where, selectFields)
		err = conn.Query(ctx, cond, vals...)
		rows, err := conn.FetchAllMap(ctx)
		if err != nil || rows == nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}



