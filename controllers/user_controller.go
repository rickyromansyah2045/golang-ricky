package controllers

import (
	"content_portal/helpers"
	"content_portal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

type UserRegisterRequest struct {
	Age      uint   `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type PointGetResponseByUser struct {
	Id        uint       `json:"id"`
	Point     uint       `json:"point"`
	ContentId uint       `json:"content_id"`
	UserId    uint       `json:"user_id"`
	UpdatedAt *time.Time `json:"updated_at"`
	CreatedAt *time.Time `json:"created_at"`
	User      UserPointResponse
}

type UserLoginRequest struct {
	Email    string `json:"email" valid:"required~email is required, email~Invalid format email"`
	Password string `json:"password" valid:"required~password is required, minstringlength(6)~password has to have minimum length of 6 characters"`
}

type UserUpdateRequest struct {
	Email    string `json:"email" valid:"email~Invalid format email"`
	Username string `json:"username"`
}

type UserRegisterResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

type UserUpdateResponse struct {
	Id        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Age       int        `json:"age"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		db: db,
	}
}

func (u *UserController) Register(ctx *gin.Context) {
	var userReq UserRegisterRequest

	err := ctx.ShouldBindJSON(&userReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	newUser := models.User{
		Age:      int(userReq.Age),
		Email:    userReq.Email,
		Username: userReq.Username,
		Password: userReq.Password,
	}

	err = u.db.Create(&newUser).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err)
			return
		}
		if err.Error() == `ERROR: duplicate key value violates unique constraint "idx_users_username" (SQLSTATE 23505)` {
			helpers.BadRequestResponse(ctx, "username is duplicated")
			return
		}
		if err.Error() == `ERROR: duplicate key value violates unique constraint "idx_users_email" (SQLSTATE 23505)` {
			helpers.BadRequestResponse(ctx, "email is duplicated")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	response := UserRegisterResponse{
		Id:       newUser.Id,
		Username: newUser.Username,
		Email:    newUser.Email,
		Age:      newUser.Age,
	}

	helpers.WriteJsonResponse(ctx, http.StatusCreated, response)
}

func (u *UserController) Login(ctx *gin.Context) {
	var userReq UserLoginRequest

	err := ctx.ShouldBindJSON(&userReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	loginUser := models.User{
		Email:    userReq.Email,
		Password: userReq.Password,
	}

	err = u.db.First(&loginUser, "email=?", userReq.Email).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, "username / password is not match")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	isValid := helpers.ComparePassword(loginUser.Password, userReq.Password)

	if !isValid {
		helpers.UnauthorizeJsonResponse(ctx, "username / password is not match")
		return
	}

	token, err := helpers.GenerateToken(loginUser.Id, loginUser.Email)
	if err != nil {
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"token": token,
	})
}

func (u *UserController) Update(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	var userReq UserUpdateRequest
	var user models.User

	err := ctx.ShouldBindJSON(&userReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	updateUser := models.User{
		Email:    userReq.Email,
		Username: userReq.Username,
	}

	// Ga perlu awal
	err = u.db.First(&user, userId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, "User data not found")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}
	// Ga perlu akhir

	err = u.db.Model(&user).Updates(updateUser).Error
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	response := UserUpdateResponse{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Age:       user.Age,
		UpdatedAt: user.UpdatedAt,
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (u *UserController) Delete(ctx *gin.Context) {

	userId, _ := ctx.Get("id")
	var user models.User

	err := u.db.First(&user, userId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.BadRequestResponse(ctx, "User not found")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	err = u.db.Delete(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.NotFoundResponse(ctx, err)
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"message": "Your account has been successfully deleted",
	})
}

func (p *UserController) GetPoint(ctx *gin.Context) {

	var points []models.Point

	err := p.db.Preload("User").Find(&points).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.NotFoundResponse(ctx, err)
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	var response []PointGetResponseByUser
	for _, point := range points {
		var userData UserPointResponse
		if point.User != nil {
			userData = UserPointResponse{
				Id:       point.User.Id,
				Username: point.User.Username,
				Email:    point.User.Email,
			}
		}
		response = append(response, PointGetResponseByUser{
			Id:        point.Id,
			Point:     point.Point,
			ContentId: point.ContentId,
			UserId:    point.UserId,
			UpdatedAt: point.UpdatedAt,
			CreatedAt: point.CreatedAt,
			User:      userData,
		})
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}
