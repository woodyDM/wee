package repository

import (
	"database/sql"
	"wee-server/blog/entity"
	"wee-server/support/database"
)

var Mysql *database.Mysql

func ArticleList(userId, page, pageSize int) *database.Page {
	query := "select id,subtract,content,create_time,title from article where author_id = ? and is_show = ? order by id desc"
	return Mysql.DBQueryPage(func() interface{} {
		return new(entity.Article)
	}, func(rows *sql.Rows, e interface{}) error {
		b := e.(*entity.Article)
		return rows.Scan(&b.Id, &b.Subtract, &b.Content, &b.CreateTime, &b.Title)
	}, page, pageSize, query, userId, true)
}

func Article(id int) *entity.Article {
	query := "select id,title,content, subtract,click_number,praised_number,create_time,update_time,is_show,author_id from article where id = ? "
	v := Mysql.DBQueryRow(func() interface{} {
		return new(entity.Article)
	}, func(row *sql.Rows, e interface{}) error {
		r := e.(*entity.Article)
		return row.Scan(&r.Id, &r.Title, &r.Content, &r.Subtract, &r.ClickNumber, &r.PraisedNumber, &r.CreateTime, &r.UpdateTime, &r.IsShow, &r.AuthorId)
	}, query, id)
	if v == nil {
		return nil
	}
	return v.(*entity.Article)
}

func GetUserByName(name string) *entity.User {
	query := "select id,name,salt,password,email,del_flag,create_time,avatar from user where name = ? "
	v := Mysql.DBQueryRow(func() interface{} {
		return new(entity.User)
	}, func(rows *sql.Rows, en interface{}) error {
		e := en.(*entity.User)
		return rows.Scan(&e.Id, &e.Name, &e.Salt, &e.Password, &e.Email, &e.DelFlag, &e.CreateTime, &e.Avatar)
	}, query, name)
	if v != nil {
		return v.(*entity.User)
	}
	return nil
}

func SaveArticle(e *entity.Article) *entity.Article {
	query := "insert into article (title, content,subtract, click_number," +
		" praised_number, create_time,update_time,is_show,author_id)values(?,?,?,?,?,?,?,?,?)"
	id := Mysql.Save(query, e.Title, e.Content, e.Subtract, e.ClickNumber, e.PraisedNumber,
		e.CreateTime, e.UpdateTime, e.IsShow, e.AuthorId)
	e.Id = int(id)
	return e
}
func UpdateArticle(e *entity.Article) {
	query := "update article set title=?, content=?,subtract=?, click_number=?," +
		" praised_number=?, update_time=?,is_show=? where id = ? "
	Mysql.Update(query, e.Title, e.Content, e.Subtract, e.ClickNumber, e.PraisedNumber,
		e.UpdateTime, e.IsShow, e.Id)
}

func SaveHistory(e *entity.History) *entity.History {
	query := "insert into visit_history(create_time,user_agent,full_ip,trim_ip,visit_path,user_hash)values(?,?,?,?,?,?)"
	id := Mysql.Save(query, e.CreateTime, e.UserAgent, e.FullIp, e.TrimIp, e.VisitPath, e.UserHash)
	e.Id = int(id)
	return e
}
