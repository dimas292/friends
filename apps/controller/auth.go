package controller

import (
	"database/sql"
	"firens/apps/pkg/tokens"
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
	queryFindByEmail = `
		SELECT id, email, password
		FROM auth
		WHERE email=$1
	`
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=7"`
	ImgUrl   string `json:"img_url" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=7"`
}

type Auth struct {
	Id       int
	Email    string
	Password string
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

func (a *AuthContorller) Login(ctx *gin.Context) {

	var req LoginRequest
	// binding
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"err":     err.Error(),
		})
		return
	}
	// prepare
	stmt, err := a.DB.Prepare(queryFindByEmail)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"err":     err.Error(),
		})
		return
	}
	// row untuk 1 data
	row := stmt.QueryRow(req.Email)

	var auth Auth
	// scan untuk memasukan data ke dalam auth struct
	err = row.Scan(
		&auth.Id,
		&auth.Email,
		&auth.Password,
	)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"err":     err.Error(),
		})
		return
	}
	// hash dan compare untuk menocokan pasword
	err = bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(req.Password))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"success": false,
			"err":     err.Error(),
		})
		return
	}
	// jwt
	tok := tokens.PayLoadToken{
		AuthId: auth.Id,
	}

	token, err := tokens.GenerateToken(&tok)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"err":     err.Error(),
		})
		return
	}

	// response
	resp := response.ResponseAPI{
		StatusCode: http.StatusOK,
		Message:    "LOGIN SUCCESS",
		Payload: gin.H{
			"token": token,
		},
	}

	ctx.JSON(resp.StatusCode, resp)
}

func (a *AuthContorller) Profile(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"id" : ctx.GetInt("authId"),
	})
}
