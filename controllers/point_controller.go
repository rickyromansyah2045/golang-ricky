package controllers

import (
	"content_portal/helpers"
	"content_portal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PointController struct {
	db *gorm.DB
}

type PointCreateRequest struct {
	Point     uint `json:"point"`
	ContentId uint `json:"content_id"`
	UserId    uint `json:"user_id"`
}

type PointCreateResponse struct {
	Id        uint       `json:"id"`
	Point     uint       `json:"point"`
	ContentId uint       `json:"content_id"`
	UserId    uint       `json:"user_id"`
	CreatedAt *time.Time `json:"created_at"`
}

type PointUpdateResponse struct {
	Id        uint       `json:"id"`
	Point     uint       `json:"point"`
	ContentId uint       `json:"content_id"`
	UserId    uint       `json:"user_id"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type PointGetResponse struct {
	Id        uint       `json:"id"`
	Point     uint       `json:"point"`
	ContentId uint       `json:"content_id"`
	UserId    uint       `json:"user_id"`
	UpdatedAt *time.Time `json:"updated_at"`
	CreatedAt *time.Time `json:"created_at"`
	User      UserPointResponse
	Content   ContentPointResponse
}

type UserPointResponse struct {
	Id       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type ContentPointResponse struct {
	Id     uint   `json:"id"`
	Title  string `json:"title"`
	UserId uint   `json:"user_id"`
}

func NewPointController(db *gorm.DB) *PointController {
	return &PointController{
		db: db,
	}
}

func (c *PointController) Create(ctx *gin.Context) {
	// userId, _ := ctx.Get("id")
	var PointReq PointCreateRequest

	err := ctx.ShouldBindJSON(&PointReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	newPoint := models.Point{
		Point:     PointReq.Point,
		ContentId: PointReq.ContentId,
		UserId:    PointReq.UserId,
		// UserId:    uint(userId.(float64)),
	}

	err = c.db.Create(&newPoint).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err)
			return
		}
		if err.Error() == `ERROR: insert or update on table "points" violates foreign key constraint "fk_contents_point" (SQLSTATE 23503)` {
			helpers.NotFoundResponse(ctx, "Content not found")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	response := PointCreateResponse{
		Id:        newPoint.Id,
		Point:     newPoint.Point,
		ContentId: newPoint.ContentId,
		UserId:    newPoint.UserId,
		CreatedAt: newPoint.CreatedAt,
	}

	helpers.WriteJsonResponse(ctx, http.StatusCreated, response)
}

func (c *PointController) Get(ctx *gin.Context) {
	var points []models.Point

	err := c.db.Preload("User").Preload("Content").Find(&points).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.NotFoundResponse(ctx, err)
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	var response []PointGetResponse
	for _, point := range points {
		var userData UserPointResponse
		if point.User != nil {
			userData = UserPointResponse{
				Id:       point.User.Id,
				Username: point.User.Username,
				Email:    point.User.Email,
			}
		}
		var contentData ContentPointResponse
		if point.Content != nil {
			contentData = ContentPointResponse{
				Id:     point.Content.Id,
				Title:  point.Content.Title,
				UserId: point.Content.UserId,
			}
		}
		response = append(response, PointGetResponse{
			Id:        point.Id,
			Point:     point.Point,
			ContentId: point.ContentId,
			UserId:    point.UserId,
			UpdatedAt: point.UpdatedAt,
			CreatedAt: point.CreatedAt,
			User:      userData,
			Content:   contentData,
		})
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (c *PointController) Update(ctx *gin.Context) {
	// userId, _ := ctx.Get("id")
	pointId := ctx.Param("pointId")
	var pointReq PointCreateRequest
	var point models.Point

	err := ctx.ShouldBindJSON(&pointReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	updatePoint := models.Point{
		Point: pointReq.Point,
	}

	err = c.db.First(&point, pointId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.NotFoundResponse(ctx, "data not found")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	err = c.db.Model(&point).Updates(updatePoint).Error
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	response := PointUpdateResponse{
		Id:        point.Id,
		Point:     point.Point,
		ContentId: point.ContentId,
		UserId:    point.UserId,
		UpdatedAt: point.UpdatedAt,
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (c *PointController) Delete(ctx *gin.Context) {

	// userId, _ := ctx.Get("id")
	pointId := ctx.Param("pointId")
	var point models.Point

	err := c.db.First(&point, pointId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.NotFoundResponse(ctx, "data not found")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	// if point.UserId != uint(userId.(float64)) {
	// 	helpers.BadRequestResponse(ctx, "you're not allowed to delete this point")
	// 	return
	// }

	err = c.db.Delete(&point).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.NotFoundResponse(ctx, err)
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"message": "Your point has been successfully deleted",
	})
}
