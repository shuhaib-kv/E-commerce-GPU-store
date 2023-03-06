package controllers

import (
	"errors"
	"ga/middleware"
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserSignUp(c *gin.Context) {
	type userInput struct {
		FirstName   string `json:"firstname" binding:"required"`
		LastName    string `json:"lastname" binding:"required"`
		UserName    string `json:"username" binding:"required"`
		Email       string `json:"email" binding:"required"`
		Password    string `json:"password" binding:"required"`
		PhoneNumber string `json:"phone" binding:"required"`
	}

	var input userInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  false,
			"message": "check json input",
			"error":   err.Error(),
		})
		return
	}

	HashPass := HashPassword(input.Password)
	user := models.Users{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		UserName:     input.UserName,
		Email:        input.Email,
		Password:     HashPass,
		Phone:        input.PhoneNumber,
		Block_status: false,
	}

	var check []models.Users
	database.Db.Find(&check)
	for _, i := range check {
		if i.Email == user.Email {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status":  false,
				"message": "change email or email already exist",
				"error":   "email Already Exist",
			})
			return
		}
	}
	for _, i := range check {
		if i.UserName == user.UserName {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status":  false,
				"error":   "Username Already Exist",
				"message": "Try another user name this user name is taken",
			})
			return
		}
	}

	result := database.Db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  false,
			"error":   "Blank input",
			"message": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status":  true,
		"message": "Account Created",
		"data":    "welcome to GPU-Ecom",
	})

	var users models.Users
	database.Db.First(&users, "email = ?", input.Email)
	wallet := models.Wallet{UsersID: users.ID, Balance: 0}
	database.Db.Create(&wallet)

	var cart = models.Cart{
		User_id: user.ID,
	}
	database.Db.Create(&cart)
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func UserLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error":   err.Error(),
			"message": "input error",
			"status":  false,
		})
		return
	}

	var user models.Users
	record := database.Db.Raw("select * from users where email=?", req.Email).Scan(&user)
	if record.Error != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error":   record.Error.Error(),
			"message": "user not found",
			"status":  false,
		})
		c.Abort()
		return
	}
	if user.Block_status {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error":   "user blocked",
			"status":  false,
			"message": "user has been blocked By admin",
		})
		c.Abort()
		return
	}
	credentialcheck := user.CheckPassword(req.Password)
	if credentialcheck != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid Password",
			"status":  false,
			"message": "try again",
		})
		c.Abort()
		return
	}
	tokenString, err := middleware.GenerateJWT(user.Email, user.ID)
	c.SetCookie("UserAuth", tokenString, 3600*24*30, "", "", false, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	res := LoginResponse{
		Email: req.Email,
		Token: tokenString,
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status":  true,
		"data":    res,
		"message": "welcome",
	})

}

func UserHome(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "welcome User Home"})

}

func AddAddress(c *gin.Context) {

	type AddressInput struct {
		Name         string `json:"name" binding:"required"`
		Phone_number uint   `json:"phone_number" binding:"required"`
		Pincode      uint   `json:"pincode" binding:"required"`
		House        string `json:"house" binding:"required"`
		Area         string `json:"area" binding:"required"`
		Landmark     string `json:"landmark" binding:"required"`
		City         string `json:"city" binding:"required"`
	}

	var input AddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	useremail := c.GetString("user")
	var UsersID uint
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "user coudnt find",
			"error":   err,
			"status":  false,
		})
		return
	}

	address := models.Address{
		UserId:       UsersID,
		Name:         input.Name,
		Phone_number: input.Phone_number,
		Pincode:      input.Pincode,
		Area:         input.Area,
		House:        input.House,
		Landmark:     input.Landmark,
		City:         input.City,
	}

	record := database.Db.Create(&address)
	if record.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": record.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Address added successfully",
		"data":    address})
}

func ShowAddress(c *gin.Context) {
	useremail := c.GetString("user")
	var UsersID uint
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})

	}
	var Adress []struct {
		Address_id   uint
		Name         string
		Phone_number uint
		Pincode      uint
		House        string
		Area         string
		Landmark     string
		City         string
	}
	record := database.Db.Select("address_id", "user_id", "name", "phone_number", "pincode", "house", "area", "landmark", "city").Table("addresses").Where("user_id=?", UsersID).Find(&Adress)
	if record.Error != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  false,
			"message": "change fields",
			"error":   record.Error.Error(),
		})
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status":  true,
		"message": "added address",
		"data":    Adress,
	})
}

func EditAddress(c *gin.Context) {
	var address models.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.Db.Model(&models.Address{}).Where("id = ?", address.ID).Updates(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address updated successfully"})

}
