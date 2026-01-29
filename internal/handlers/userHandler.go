package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/helper"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/repository"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/utils"
)

// RegisterRequest - request body for user registration
// @Description User registration payload
type RegisterRequest struct {
	Fullname string `json:"fullname" binding:"required,min=2,max=50" example:"John Doe"`
	Username string `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

// LoginRequest - request body for user login
// @Description User login payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

// Register godoc
// @Summary     Register a new user
// @Description Create a new user account
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       user  body      RegisterRequest  true  "User registration data"
// @Success     200   {object}  map[string]interface{}
// @Failure     400   {object}  map[string]interface{}
// @Failure     409   {object}  map[string]interface{}
// @Failure     500   {object}  map[string]interface{}
// @Router      /users/register [post]
func Register(c *gin.Context) {
	// ✅ Simple Go pattern: bind + validate in one step
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	// Business logic
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to hash password"})
		return
	}

	user := models.User{
		Fullname: req.Fullname,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}
	userData, err := repository.Register(&user)

	if err != nil {
		if msg, ok := utils.ParsePostgresError(err); ok {
			utils.ResponseError(c, http.StatusConflict, msg, nil)
			return
		}
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		return
	}
	userResponse := utils.ToUserResponse(userData)
	utils.ResponseSuccess(c, http.StatusOK, "user registered successfully", userResponse)

}

// Login godoc
// @Summary     User login
// @Description Authenticate user and return JWT token
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       credentials  body      LoginRequest  true  "Login credentials"
// @Success     200          {object}  map[string]interface{}
// @Failure     400          {object}  map[string]interface{}
// @Failure     500          {object}  map[string]interface{}
// @Router      /users/login [post]
func Login(c *gin.Context) {
	// ✅ Simple Go pattern: bind + validate in one step
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid Credential", nil)
		return
	}

	if isPasswordValid := utils.CheckPassword(req.Password, user.Password); isPasswordValid == false {
		utils.ResponseError(c, http.StatusUnauthorized, "Invalid Credential", nil)
		return
	}
	claims := utils.JWTClaims{
		UserID: user.ID.String(),
		Roles:  []string{"user", "admin"}, // dynamic from DB
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "e-commerce-app",
		},
	}
	token, err := utils.CreateToken(claims)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "could not create token", nil)
		return
	}
	userResponse := utils.ToUserResponse(user)
	utils.ResponseSuccess(c, http.StatusOK, "loggedin successfully", gin.H{
		"data":  userResponse,
		"token": token,
	})
}

// GetAllUsers godoc
// @Summary     Get all users
// @Description Retrieve list of all users (requires auth)
// @Tags        Users
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Success     200  {object}  map[string]interface{}
// @Failure     500  {object}  map[string]interface{}
// @Router      /users/all [get]
func GetAllUsers(c *gin.Context) {
	users, err := repository.GetAllUsers()
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}
	utils.ResponseSuccess(c, http.StatusOK, "data fetched successfully", users)
}

// GetUser godoc
// @Summary     Get user by ID
// @Description Retrieve a single user by UUID
// @Tags        Users
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       id   path      string  true  "User UUID"
// @Success     200  {object}  map[string]interface{}
// @Failure     400  {object}  map[string]interface{}
// @Failure     404  {object}  map[string]interface{}
// @Router      /users/user/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("id") // returns string
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user, err := repository.GetUserByUUID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err, "status": "Failure", "message": "Something went wrong"})
		return
	}
	userResponse := utils.ToUserResponse(user)
	utils.ResponseSuccess(c, http.StatusOK, "data fetched successfully", userResponse)
}

// GetUserByEmail godoc
// @Summary     Get user by email
// @Description Retrieve a user by email query param
// @Tags        Users
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       email  query     string  true  "User email"
// @Success     200    {object}  map[string]interface{}
// @Failure     404    {object}  map[string]interface{}
// @Router      /users/user [get]
func GetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err, "status": "Failure", "message": "Something went wrong"})
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "data fetched successfully", user)
}

func GetFilterAndSearchUsers(c *gin.Context) {
	// 1. Parse and set defaults for pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}

	// 2. Create the params object
	params := helper.UserFilterParams{
		ProductName: c.Query("productName"),
		FullName:    c.Query("fullname"),
		Page:        page,
		Limit:       limit,
	}
	users, total, err := repository.FilterAndSearchUsers(params)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}
	// 3. Return response with metadata
	utils.ResponseSuccess(c, http.StatusOK, "data fetched successfully", gin.H{
		"users": users,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}
