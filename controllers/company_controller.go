package controllers

import (
	"content_portal/helpers"
	"content_portal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompanyController struct {
	db *gorm.DB
}

type AllUserRegisterResponse struct {
	Id        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Age       int        `json:"age"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Point     PointResponse
}

type PointResponse struct {
	Point     uint `json:"point"`
	ContentId uint `json:"content_id"`
}

type CompanyRegisterRequest struct {
	Role     bool   `json:"role"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type CompanyLoginRequest struct {
	Email    string `json:"email" valid:"required~email is required, email~Invalid format email"`
	Password string `json:"password" valid:"required~password is required, minstringlength(6)~password has to have minimum length of 6 characters"`
}

type CompanyUpdateRequest struct {
	Email    string `json:"email" valid:"email~Invalid format email"`
	Username string `json:"username"`
}

type CompanyRegisterResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     bool   `json:"role"`
}

type CompanyUpdateResponse struct {
	Id        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      bool       `json:"role"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func NewCompanyController(db *gorm.DB) *CompanyController {
	return &CompanyController{
		db: db,
	}
}

func (u *CompanyController) Register(ctx *gin.Context) {
	var companyReq CompanyRegisterRequest

	err := ctx.ShouldBindJSON(&companyReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	newUser := models.Company{
		Role:     bool(companyReq.Role),
		Email:    companyReq.Email,
		Username: companyReq.Username,
		Password: companyReq.Password,
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

	response := CompanyRegisterResponse{
		Id:       newUser.Id,
		Username: newUser.Username,
		Email:    newUser.Email,
		Role:     newUser.Role,
	}

	helpers.WriteJsonResponse(ctx, http.StatusCreated, response)
}

func (u *CompanyController) Login(ctx *gin.Context) {
	var companyReq CompanyLoginRequest

	err := ctx.ShouldBindJSON(&companyReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	loginCompany := models.Company{
		Email:    companyReq.Email,
		Password: companyReq.Password,
	}

	err = u.db.First(&loginCompany, "email=?", companyReq.Email).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, "username / password is not match")
			return
		}
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	isValid := helpers.ComparePassword(loginCompany.Password, companyReq.Password)

	if !isValid {
		helpers.UnauthorizeJsonResponse(ctx, "username / password is not match")
		return
	}

	token, err := helpers.GenerateToken(loginCompany.Id, loginCompany.Email)
	if err != nil {
		helpers.InternalServerJsonResponse(ctx, err)
		return
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"token": token,
	})
}

func (u *CompanyController) Update(ctx *gin.Context) {

	var companyReq CompanyUpdateRequest
	var company models.Company

	err := ctx.ShouldBindJSON(&companyReq)
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	updateUser := models.Company{
		Email:    companyReq.Email,
		Username: companyReq.Username,
	}

	err = u.db.Model(&company).Updates(updateUser).Error
	if err != nil {
		helpers.BadRequestResponse(ctx, err)
		return
	}

	response := CompanyUpdateResponse{
		Id:        company.Id,
		Username:  company.Username,
		Email:     company.Email,
		Role:      company.Role,
		UpdatedAt: company.UpdatedAt,
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (u *CompanyController) Delete(ctx *gin.Context) {
	var user models.Company

	err := u.db.Delete(&user).Error
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

func (p *CompanyController) GetUserRegister(ctx *gin.Context) {

	var users []models.User

	// Get all records user
	err := p.db.Find(&users).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helpers.NotFoundResponse(ctx, err)
			return
		}
		helpers.BadRequestResponse(ctx, err)
		return
	}

	var response []AllUserRegisterResponse
	for _, user := range users {

		var pointData PointResponse
		if user.Point != nil {
			pointData = PointResponse{
				Point:     user.Point.Point,
				ContentId: user.Point.ContentId,
			}
		}

		response = append(response, AllUserRegisterResponse{
			Id:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			Age:       user.Age,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Point:     pointData,
		})
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}
