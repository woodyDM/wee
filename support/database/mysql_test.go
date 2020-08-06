package database

import (
	"database/sql"
	"fmt"
	"testing"
	"wee-server/blog/entity"
)

func TestQuery(t *testing.T) {

	result := DBQuery(func() interface{} {
		return new(entity.Article)
	}, func(rows *sql.Rows, e interface{}) error {
		bloger := e.(*entity.Article)
		return rows.Scan(&bloger.Id, &bloger.Title, &bloger.Content, &bloger.Subtract)
	}, "select id,title,content,subtract from article where id > 36")
	fmt.Println(result)

}
func TestDBQueryPage(t *testing.T) {

	result := DBQueryPage(func() interface{} {
		return new(entity.Article)
	}, func(rows *sql.Rows, e interface{}) error {
		bloger := e.(*entity.Article)
		return rows.Scan(&bloger.Id, &bloger.Title, &bloger.Content, &bloger.Subtract)
	}, 1, 2, "select id,title,content,subtract from article where id > ? and id < ? ", 36, 40)
	fmt.Println(result)

}

func TestDBQueryRow(t *testing.T) {
	var r int
	DBQueryScalar(&r, "select count(*) from article where id > ?", 36)
	fmt.Println(r)
}
