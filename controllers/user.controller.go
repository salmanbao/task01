package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/salmanbao/server/db/sqlc"
)

type UserController struct {
	db *db.Queries
}

func NewUserController(db *db.Queries) *UserController {
	return &UserController{db}
}

func (ac *UserController) CreateUser(ctx *gin.Context) {
	var user db.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	args := &db.CreateUserParams{
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
	}

	user, err := ac.db.CreateUser(ctx, *args)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "data": "duplicate record"})
			return
		}
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": user}})
}

func (ac *UserController) GenerateOTP(ctx *gin.Context) {
	var p struct {
		PhoneNumber string `json:"phone_number"`
	}

	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := ac.db.GetUserByPhone(ctx, p.PhoneNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "data": err.Error()})
		return
	}
	expiration := time.Now().Add(time.Minute * 1).UTC()
	t := pgtype.Timestamp{}
	t.Scan(expiration)
	otp := fmt.Sprintf("%d", rand.Intn(9000)+1000)

	args := &db.UpdateUserParams{
		ID:                user.ID,
		Name:              user.Name,
		PhoneNumber:       user.PhoneNumber,
		Opt:               otp,
		OptExpirationTime: t,
	}

	user, err = ac.db.UpdateUser(ctx, *args)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "data": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": user}})
}

func (ac *UserController) VerifyOTP(ctx *gin.Context) {
	var otp struct {
		PhoneNumber string `json:"phone_number"`
		OTP         string `json:"otp"`
	}

	if err := ctx.ShouldBindJSON(&otp); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := ac.db.GetUserByPhone(ctx, otp.PhoneNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "data": err.Error()})
		return
	}
	if user.Opt != otp.OTP {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "data": "invalid otp"})
		return
	}
	now := time.Now().UTC().Unix()
	if user.Opt == otp.OTP {
		if now > user.OptExpirationTime.Time.UTC().Unix() {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "data": "otp expired"})
			return
		}
	}
	ctx.JSON(http.StatusAccepted, gin.H{"status": "success", "data": "Verified"})

}
