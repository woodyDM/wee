package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Page struct {
	Page         int
	PageSize     int
	TotalElement int64
	Data         []interface{}
}

type Mysql struct {
	db *sql.DB
}

const (
	alias = "CANT_USED_BY_USER"
)

func NewMysql(source string) *Mysql {

	driver := "mysql"
	DB, err := sql.Open(driver, source)
	if err != nil {
		panic(fmt.Errorf("Unable to init database of %s. Reason is %v. ", source, err))
	}
	log.Printf("Init database [%s] successful. ", driver)
	mysql := &Mysql{db: DB}
	mysql.tryConnect()
	return mysql
}

func (m *Mysql) Close() {
	m.db.Close()
}

func (m *Mysql) tryConnect() {
	var dummy int
	m.DBQueryScalar(&dummy, "select 1")
}

func (m *Mysql) DBQuery(provider func() interface{}, mapper func(rows *sql.Rows, entity interface{}) error, sql string, args ...interface{}) []interface{} {
	rows, e := m.db.Query(sql, args...)
	defer rows.Close()
	if e != nil {
		panic(e)
	}
	var result []interface{}
	for rows.Next() {
		v := provider()
		err := mapper(rows, v)
		if err != nil {
			panic(err)
		}
		result = append(result, v)
	}
	return result
}

func (m *Mysql) DBQueryRow(provider func() interface{}, mapper func(rows *sql.Rows, entity interface{}) error, sql string, args ...interface{}) interface{} {
	query := fmt.Sprintf("select %s.* from (%s) %s limit 1", alias, sql, alias)
	rows, e := m.db.Query(query, args...)
	defer rows.Close()
	if e != nil {
		panic(e)
	}

	if rows.Next() {
		v := provider()
		err := mapper(rows, v)
		if err != nil {
			panic(err)
		}
		return v
	}
	return nil
}

func (m *Mysql) dBQueryRow(mapper func(row *sql.Row) error, query string, args ...interface{}) {
	row := m.db.QueryRow(query, args...)
	e := mapper(row)
	if e != nil {
		panic(e)
	}
}

func (m *Mysql) DBQueryScalar(r interface{}, sqlstring string, args ...interface{}) {
	m.dBQueryRow(func(row *sql.Row) error {
		return row.Scan(r)
	}, sqlstring, args...)
}

func (m *Mysql) DBQueryPage(provider func() interface{}, mapper func(rows *sql.Rows, entity interface{}) error, page int, pageSize int, sql string, args ...interface{}) *Page {
	checkParams(page, pageSize)

	countSql := fmt.Sprintf("select count(*)  from (%s) %s;", sql, alias)
	var total int64
	m.DBQueryScalar(&total, countSql, args...)
	var data []interface{}
	start := page * pageSize
	if int64(start) < total {
		querySql := fmt.Sprintf("select %s.* from (%s) %s limit %d,%d", alias, sql, alias, start, pageSize)
		data = m.DBQuery(provider, mapper, querySql, args...)
	}
	return &Page{
		Page:         page,
		PageSize:     pageSize,
		TotalElement: total,
		Data:         data,
	}

}

func (m *Mysql) Save(query string, args ...interface{}) int64 {
	return m.exec(query, func(r sql.Result) (i int64, e error) {
		return r.LastInsertId()
	}, args...)
}

func (m *Mysql) Update(query string, args ...interface{}) int64 {
	return m.exec(query, func(r sql.Result) (i int64, e error) {
		return r.RowsAffected()
	}, args...)
}

func (m *Mysql) exec(query string, resultMapper func(r sql.Result) (int64, error), args ...interface{}) int64 {
	stmt, err := m.db.Prepare(query)
	if err != nil {
		panic(err)
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		panic(err)
	}
	rows, err := resultMapper(result)
	if err != nil {
		panic(err)
	}
	return rows
}
func checkParams(page int, size int) {
	if page < 0 {
		panic(fmt.Errorf("page should >=0. "))
	}
	if size <= 0 {
		panic(fmt.Errorf("size should >0. "))
	}
}
