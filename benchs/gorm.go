package benchs

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

var gormdb *gorm.DB

var tableName string

func init() {
	st := NewSuite("gorm")
	tableName = "model"
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, GormInsert)
		//st.AddBenchmark("MultiInsert 100 row", 500*ORM_MULTI, GormMultiInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, GormUpdate)
		st.AddBenchmark("Read", 2000*ORM_MULTI, GormRead)
		st.AddBenchmark("MultiRead limit 100", 2000*ORM_MULTI, GormReadSlice)

		var err error
		gormdb, err = gorm.Open("mysql", ORM_SOURCE)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func GormInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})

	for i := 1; i <= b.N; i++ {
		m.Id = 0
		err := gormdb.Table(tableName).Create(m).Error
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func GormMultiInsertMulti(b *B) {
	panic(fmt.Errorf("Not support multi insert"))
}

func GormUpdate(b *B) {
	whereField := map[string]interface{}{}
	updateField := map[string]interface{}{}

	for i := 1; i <= b.N; i++ {
		fmt.Sprintf("%v\n", i)
		whereField["id"] =  i
		updateField["age"] = i
			err := gormdb.Table(tableName).Where(whereField).Update(updateField).Error
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func GormRead(b *B) {
	whereField := map[string]interface{}{}
	for i := 1; i <= b.N; i++ {
		var m Model
		whereField["id"] =  i
		err := gormdb.Table(tableName).Where(whereField).Find(&m).Error
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func GormReadSlice(b *B) {
	for i := 1; i <= b.N; i++ {
		var models []Model
		err := gormdb.Table(tableName).Where("id >= ?", 1).Find(&models).Limit(100).Error
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

