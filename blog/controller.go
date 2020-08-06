package blog

import (
	"fmt"
	"log"
	"strconv"
	"wee-server/blog/repository"
	"wee-server/blog/service"
	"wee-server/blog/vo"
	"wee-server/wee"
)

func ApiHealthIndex(ctx *wee.Context) {
	ctx.String("Service is on\nV1.0.0\n")
}

func ArticleListController(ctx *wee.Context) {
	userId := ctx.GetPathVariableInt("userId")
	page := ctx.GetRequestFormInt("page", 0)
	//size := ctx.GetRequestFormInt("size", 10)
	size := 10 //now fix size to 10
	result := vo.ToPageResp(service.GetArticleListCached(userId, page, size))
	resp := vo.NewResponse(result)
	ctx.Json(resp)
}

func ArticleViewController(ctx *wee.Context) {
	id := ctx.GetPathVariableInt("id")
	currentUserId := 0
	user := service.GetCurrentUser(ctx)
	if user != nil {
		currentUserId = user.Id
	}
	view := service.GetArticleView(id, currentUserId)
	if view != nil && view.Show {
		resp := vo.NewResponse(view)
		ctx.Json(resp)
	} else {
		vo.WriteToResponse(ctx, "no article found")
	}
}

func UserArticleViewController(ctx *wee.Context) {
	id := ctx.GetPathVariableInt("id")
	user := service.GetCurrentUser(ctx)
	view := service.GetArticleView(id, user.Id)
	if view == nil {
		vo.WriteToResponse(ctx, fmt.Sprintf("No article found for id %d", id))
		return
	}
	if view.AuthorId != user.Id {
		vo.WriteToResponse(ctx, "User mismatch.")
		return
	}
	resp := vo.NewResponse(view)
	ctx.Json(resp)
}

func UserSaveArticleController(ctx *wee.Context) {
	request := new(vo.SaveOrUpdateArticleRequest)
	ctx.JsonRequest(request)
	if e, ok := request.GetError(); !ok {
		vo.WriteToResponse(ctx, e)
		return
	}
	user := service.GetCurrentUser(ctx)
	request.UserId = user.Id
	articleId := service.SaveArticle(request)
	ctx.Json(vo.NewResponse(strconv.Itoa(articleId)))
}
func UserUpdateArticleController(ctx *wee.Context) {
	articleId := ctx.GetPathVariableInt("id")
	article := repository.Article(articleId)
	if article == nil {
		vo.WriteToResponse(ctx, "Not found")
		return
	}
	request := new(vo.SaveOrUpdateArticleRequest)
	ctx.JsonRequest(request)
	if e, ok := request.GetError(); !ok {
		vo.WriteToResponse(ctx, e)
		return
	}
	user := service.GetCurrentUser(ctx)
	request.UserId = user.Id
	if user.Id != article.AuthorId {
		vo.WriteToResponse(ctx, "User not Match！")
		return
	}
	service.UpdateArticle(request, article)
	ctx.Json(vo.NewResponse(strconv.Itoa(articleId)))
}

func Login(ctx *wee.Context) {
	request := new(vo.LoginRequest)
	ctx.JsonRequest(request)
	if request.Name == "" || request.Password == "" {
		ctx.Json(vo.NewFailResponse("Fail"))
		return
	}
	log.Printf("[%s] is trying to login.", request.Name)
	user := service.TryLogin(ctx, request.Name, request.Password)
	if user != nil {
		log.Printf("User [%s] login success! \n", request.Name)
		userVO := vo.NewUser(user)
		ctx.Json(vo.NewResponse(userVO))
	} else {
		ctx.Json(vo.NewFailResponse("Fail"))
	}

}

func Logout(ctx *wee.Context) {
	service.Logout(ctx)
	ctx.Json(vo.NewResponse("OK"))
}
