package controllers

import (
	"content_portal/helpers"
	"content_portal/models"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContentController struct {
	db *gorm.DB
}

type ContentCreateRequest struct {
	Title      string `json:"title"`
	IsiContent string `json:"isi_content"`
	PhotoUrl   string `json:"photo_url"`
}

type ContentCreateResponse struct {
	Id         uint       `json:"id"`
	Title      string     `json:"title"`
	IsiContent string     `json:"isi_content"`
	PhotoUrl   string     `json:"photo_url"`
	UserId     uint       `json:"user_id"`
	CreatedAt  *time.Time `json:"created_at"`
}

type ContentUpdateResponse struct {
	Id         uint       `json:"id"`
	Title      string     `json:"title"`
	IsiContent string     `json:"isi_content"`
	PhotoUrl   string     `json:"photo_url"`
	UserId     uint       `json:"user_id"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type ContentGetResponse struct {
	Id         uint       `json:"id"`
	Title      string     `json:"title"`
	IsiContent string     `json:"isi_content"`
	PhotoUrl   string     `json:"photo_url"`
	UserId     uint       `json:"user_id"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	User       UserDataResponse
}

type UserDataResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func NewContentController(db *gorm.DB) *ContentController {
	return &ContentController{
		db: db,
	}
}

func (p *ContentController) Create(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	var contentReq ContentCreateRequest

	err := ctx.ShouldBindJSON(&contentReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	newContent := models.Content{
		Title:      contentReq.Title,
		IsiContent: contentReq.IsiContent,
		PhotoUrl:   contentReq.PhotoUrl,
		UserId:     uint(userId.(float64)),
	}

	err = p.db.Create(&newContent).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err)
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	response := ContentCreateResponse{
		Id:         newContent.Id,
		Title:      newContent.Title,
		IsiContent: newContent.IsiContent,
		PhotoUrl:   newContent.PhotoUrl,
		UserId:     newContent.UserId,
		CreatedAt:  newContent.CreatedAt,
	}

	helpers.WriteJsonResponse(ctx, http.StatusCreated, response)
}

func (p *ContentController) Get(ctx *gin.Context) {

	var contents []models.Content

	err := p.db.Preload("User").Find(&contents).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err)
			return
		}
		helpers.BadRequestResponse(ctx, err)
		return
	}

	var response []ContentGetResponse
	for _, content := range contents {
		var userData UserDataResponse
		if content.User != nil {
			userData = UserDataResponse{
				Username: content.User.Username,
				Email:    content.User.Email,
			}
		}
		response = append(response, ContentGetResponse{
			Id:         content.Id,
			Title:      content.Title,
			IsiContent: content.IsiContent,
			PhotoUrl:   content.PhotoUrl,
			UserId:     content.UserId,
			CreatedAt:  content.CreatedAt,
			UpdatedAt:  content.UpdatedAt,
			User:       userData,
		})
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (p *ContentController) GetById(ctx *gin.Context) {

	userId, _ := ctx.Get("id")
	contentId := ctx.Param("contentId")
	var contents models.Content

	err := p.db.First(&contents, contentId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, "data not found")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	if contents.UserId != uint(userId.(float64)) {
		helpers.UnauthorizeJsonResponse(ctx, "you're not allowed to access this content")
		return
	}

	response := ContentUpdateResponse{
		Id:         contents.Id,
		Title:      contents.Title,
		IsiContent: contents.IsiContent,
		PhotoUrl:   contents.PhotoUrl,
		UserId:     contents.UserId,
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (p *ContentController) Update(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	contentId := ctx.Param("contentId")
	var contentReq ContentUpdateResponse
	var content models.Content

	err := ctx.ShouldBindJSON(&contentReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	updatedContent := models.Content{
		Title:      contentReq.Title,
		IsiContent: contentReq.IsiContent,
		PhotoUrl:   contentReq.PhotoUrl,
	}

	_, errCreate := govalidator.ValidateStruct(&updatedContent)
	if errCreate != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	err = p.db.First(&content, contentId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, "data not found")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	if content.UserId != uint(userId.(float64)) {
		helpers.UnauthorizeJsonResponse(ctx, "you're not allowed to update or edit this content")
		return
	}

	err = p.db.Model(&content).Updates(updatedContent).Error
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	response := ContentUpdateResponse{
		Id:         content.Id,
		Title:      content.Title,
		IsiContent: content.IsiContent,
		PhotoUrl:   content.PhotoUrl,
		UserId:     content.UserId,
		UpdatedAt:  content.UpdatedAt,
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (p *ContentController) Delete(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	contentId := ctx.Param("contentId")
	var photo models.Content

	err := p.db.First(&photo, contentId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, "data not found")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	if photo.UserId != uint(userId.(float64)) {
		helpers.UnauthorizeJsonResponse(ctx, "you're not allowed to delete this content")
		return
	}

	err = p.db.Delete(&photo).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err.Error())
			return
		}
		helpers.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"message": "Your content has been successfully deleted",
	})

}
