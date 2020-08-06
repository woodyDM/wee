package repository

import (
	"testing"
	"time"
	"wee-server/blog/entity"
	"wee-server/support/database"
)

func TestArticleList(t *testing.T) {

	list := ArticleList(2, 1, 5)
	t.Log(list)
}

func TestGetUserByName(t *testing.T) {

	u1 := GetUserByName("w")
	u2 := GetUserByName("wk")
	t.Log(u1)
	t.Log(u2)

}

func TestSaveHistory(t *testing.T) {
	mysql := database.NewMysql("root:123456@tcp(127.0.0.1:3306)/blog?charset=utf8mb4")
	defer mysql.Close()
	Mysql = mysql
	h := &entity.History{
		CreateTime: time.Now(),
		UserAgent:  "UA",
		FullIp:     "12",
		TrimIp:     "23",
		VisitPath:  "/path",
		UserHash:   12345,
	}
	SaveHistory(h)
	if h.Id == 0 {
		t.Fail()
	}
}
