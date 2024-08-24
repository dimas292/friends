package controller

import (
	"database/sql"
	"firens/apps/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type AuthContorller struct {
	DB *sql.DB
}

var (
	queryCreate = `
		INSERT INTO auth (email, password, img_url)
		VALUES ($1, $2, $3)
		`
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=7"`
	ImgUrl   string `json:"img_url" validate:"required"`
}

func (a *AuthContorller) Register(ctx *gin.Context) {

	var req RegisterRequest
	// binding
	err := ctx.ShouldBind(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"err":     err.Error(),
		})
	}
	// validasi
	val := validator.New()
	err = val.Struct(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"err":     err.Error(),
		})
		return
	}
	// crypto / enkripsi password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	req.Password = string(hash)
	// prepare agar tidal gampang kena sql inject
	stmt, err := a.DB.Prepare(queryCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"err":     err.Error(),
		})
		return
	}

	// execute
	_, err = stmt.Exec(req.Email, req.Password, req.ImgUrl)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"err":     err.Error(),
		})
		return
	}
	// response custom contract
	resp := response.ResponseAPI{
		StatusCode: http.StatusCreated,
		Message:    "CREATED SUCCESS",
		Payload:    req,
	}
	ctx.JSON(resp.StatusCode, resp)
}
