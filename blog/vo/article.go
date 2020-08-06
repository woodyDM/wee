package vo

import (
	"time"
	"unicode/utf8"
	"wee-server/blog/entity"
)

type ArticleVO struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Subtract   string `json:"subtract"`
	CreateTime string `json:"createTime"`
}

type OneArticleVO struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Subtract      string `json:"subtract"`
	Date          string `json:"date"`
	IsCurrentUser bool   `json:"currentUser"`
	Show          bool   `json:"show"`
	AuthorId      int    `json:"authorId"`
	ClickNum      int    `json:"clickNum"`
}

type SaveOrUpdateArticleRequest struct {
	Subtract string `json:"subtract"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Show     bool   `json:"show"`
	UserId   int
}

func NewArticleVO(a *entity.Article) *ArticleVO {

	return &ArticleVO{
		Id:         a.Id,
		Title:      a.Title,
		Content:    a.Content,
		Subtract:   a.Subtract,
		CreateTime: stampToDate(a.CreateTime),
	}
}

func NewOneArticleVO(a *entity.Article, currentId int) *OneArticleVO {
	return &OneArticleVO{
		Id:            a.Id,
		Title:         a.Title,
		Content:       a.Content,
		Subtract:      a.Subtract,
		Date:          stampToDate(a.CreateTime),
		IsCurrentUser: currentId == a.AuthorId,
		Show:          a.IsShow,
		AuthorId:      a.AuthorId,
		ClickNum:      a.ClickNumber,
	}
}

func stampToDate(stamp int) string {
	t := time.Unix(int64(stamp), 0)
	return t.Format("2006-01-02")
}

func (r *SaveOrUpdateArticleRequest) GetError() (e string, ok bool) {
	if r.Subtract == "" {
		return "摘要不能为空", false
	}
	if r.Title == "" {
		return "标题不能为空", false
	}
	if r.Content == "" {
		return "内容不能为空", false
	}
	if utf8.RuneCountInString(r.Subtract) > 100 {
		return "摘要最大长度不超过100", false
	}
	if utf8.RuneCountInString(r.Title) > 100 {
		return "标题最大长度不超过100", false
	}
	return "", true
}
func (r *SaveOrUpdateArticleRequest) ToArticle() *entity.Article {
	createTime := GetZone8UnixTime()
	return &entity.Article{
		Title:         r.Title,
		Content:       r.Content,
		Subtract:      r.Subtract,
		ClickNumber:   0,
		PraisedNumber: 0,
		CreateTime:    createTime,
		UpdateTime:    createTime,
		IsShow:        r.Show,
		AuthorId:      r.UserId,
	}
}

func (r *SaveOrUpdateArticleRequest) UpdateArticle(a *entity.Article) {
	a.Subtract = r.Subtract
	a.Title = r.Title
	a.UpdateTime = GetZone8UnixTime()
	a.Content = r.Content
	a.IsShow = r.Show
}
