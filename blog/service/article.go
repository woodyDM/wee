package service

import (
	"fmt"
	"wee-server/blog/entity"
	"wee-server/blog/repository"
	"wee-server/blog/vo"
	"wee-server/support/cache"
	"wee-server/support/database"
)

var CacheService *cache.RedisService
var CacheTimeout int

const (
	articleListCache = "articleList_%d_%d_%d"
	articleCache     = "article_%d"
)

func GetArticleListCached(userId, page, size int) *database.Page {
	key := CacheService.WrapCacheKey(fmt.Sprintf(articleListCache, userId, page, size))
	result, ok := CacheService.GetCache(key, CacheTimeout, database.Page{}, func() interface{} {
		return GetArticleList(userId, page, size)
	})
	if ok {
		return result.(*database.Page)
	} else {
		return nil
	}
}

func GetArticleList(userId, page, size int) *database.Page {
	p := repository.ArticleList(userId, page, size)
	list := make([]interface{}, 0)
	for _, it := range p.Data {
		article := it.(*entity.Article)
		articleVo := vo.NewArticleVO(article)
		list = append(list, articleVo)
	}
	return &database.Page{
		Page:         p.Page,
		PageSize:     p.PageSize,
		TotalElement: p.TotalElement,
		Data:         list,
	}
}

func GetArticleView(id int, currentUserId int) *vo.OneArticleVO {
	cacheKey := CacheService.WrapCacheKey(fmt.Sprintf(articleCache, id))
	article, ok := CacheService.GetCache(cacheKey, CacheTimeout, entity.Article{}, func() interface{} {
		return repository.Article(id)
	})
	if !ok {
		return nil
	}
	return vo.NewOneArticleVO(article.(*entity.Article), currentUserId)
}

func SaveArticle(r *vo.SaveOrUpdateArticleRequest) int {
	article := r.ToArticle()
	repository.SaveArticle(article)
	deleteCache(r.UserId, 0)
	return article.Id
}

func UpdateArticle(r *vo.SaveOrUpdateArticleRequest, a *entity.Article) {
	r.UpdateArticle(a)
	repository.UpdateArticle(a)
	deleteCache(a.AuthorId, a.Id)
}

func deleteCache(userId int, articleId int) {
	if articleId != 0 {
		articleKey := CacheService.WrapCacheKey(fmt.Sprintf(articleCache, articleId))
		CacheService.Del(articleKey)
	}
	pageSize := 10
	for page := 0; page < 10; page++ {
		listKey := CacheService.WrapCacheKey(fmt.Sprintf(articleListCache, userId, page, pageSize))
		CacheService.Del(listKey)
	}
}
