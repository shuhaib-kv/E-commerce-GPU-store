package controllers

import (
	"errors"
	"fmt"
	"ga/middleware"
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserSignUp(c *gin.Context) {

	First_Name := c.PostForm("firstname")
	Last_Name := c.PostForm("lastname")
	User_Name := c.PostForm("username")
	Email := c.PostForm("email")
	Password := c.PostForm("password")
	HashPass := HashPassword(Password)
	Phone_Number := c.PostForm("phoneno")
	user := models.Users{FirstName: First_Name, LastName: Last_Name, UserName: User_Name, Email: Email, Password: HashPass, Phone: Phone_Number, Block_status: false}
	var check []models.Users
	database.Db.Find(&check)
	for _, i := range check {
		if i.Email == user.Email {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "email Already Exist",
			})
			return
		}
	}
	for _, i := range check {
		if i.UserName == user.UserName {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Username Already Exist",
			})
			return
		}
	}
	if user.FirstName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Name is required",
		})
		return
	}
	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email is required",
		})
		return
	}
	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password is required",
		})

		return
	}
	result := database.Db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Account Created",
	})
	var users models.Users
	database.Db.First(&users, "email = ?", Email)
	wallet := models.Wallet{UsersID: users.ID, Balance: 0}
	database.Db.Create(&wallet)
	var cart = models.Cart{
		User_id: user.ID,
	}
	database.Db.Create(&cart)

}
func UserLogin(c *gin.Context) {
	Email := c.PostForm("email")
	Password := c.PostForm("password")
	HashPassword(Password)
	var user models.Users
	record := database.Db.Raw("select * from users where email=?", Email).Scan(&user)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	if user.Block_status {
		c.JSON(404, gin.H{"msg": "user has been blocked By admin"})
		c.Abort()
		return
	}
	credentialcheck := user.CheckPassword(Password)
	if credentialcheck != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Password"})
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

	c.JSON(200, gin.H{"email": Email, "password": Password, "token": tokenString})
}

func UserHome(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "welcome User Home"})

}

func AddAddress(c *gin.Context) {
	// Name := c.PostForm("name")
	// Phone_number := c.PostForm("phone")
	// pho, _ := strconv.Atoi(Phone_number)
	// Pincodeu := c.PostForm("pincode")
	// pin, _ := strconv.Atoi(Pincodeu)
	// House := c.PostForm("house")
	// Area := c.PostForm("area")
	// Landmark := c.PostForm("landmark")
	// City := c.PostForm("city")
	useremail := c.GetString("user")
	fmt.Println(useremail)
	var UsersID int
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})
	}
	var body struct {
		Name         string
		Phone_number int
		Pincode      int
		House        string
		Area         string
		Landmark     string
		City         string
	}
	c.Bind(&body)
	address := models.Address{
		UserId:       UsersID,
		Name:         body.Name,
		Phone_number: body.Phone_number,
		Pincode:      body.Pincode,
		Area:         body.Area,
		House:        body.House,
		Landmark:     body.Landmark,
		City:         body.City,
	}
	record := database.Db.Create(&address)
	if record.Error != nil {
		c.JSON(404, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"Adress Added": address,
	})
}

func ShowAddress(c *gin.Context) {
	useremail := c.GetString("user")
	var UsersID int
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})

	}
	var Adress []struct {
		Address_id   int
		Name         string
		Phone_number int
		Pincode      int
		House        string
		Area         string
		Landmark     string
		City         string
	}
	record := database.Db.Select("address_id", "user_id", "name", "phone_number", "pincode", "house", "area", "landmark", "city").Table("addresses").Where("user_id=?", UsersID).Find(&Adress)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"error": record.Error.Error(),
		})
	}

	c.JSON(200, gin.H{
		"address": Adress,
	})

}
func SelectAddress(c *gin.Context) {
	useremail := c.GetString("user")
	var UsersID int
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})

	}
	var address []models.Address
	database.Db.Where("addresses.users_id = ?", UsersID).Find(&address)
	for _, i := range address {
		c.JSON(200, gin.H{
			"Name":          i.Name,
			"Phone Number":  i.Phone_number,
			"Pincode":       i.Pincode,
			"House Address": i.House,
			"Area":          i.Area,
			"Landmark":      i.Landmark,
			"City":          i.City,
			"id":            i.Address_id,
		})
	}
}

func EditAddress(c *gin.Context) {
	useremail := c.GetString("user")
	adress := c.Param("id")
	var UsersID int
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})

	}
	var body struct {
		Name         string
		Phone_number int
		Pincode      int
		House_Adress string
		Area         string
		Landmark     string
		City         string
	}

	c.Bind(&body)

	var Address []models.Address

	results := database.Db.First(&Address, adress)
	if results.Error != nil {
		c.JSON(400, gin.H{
			"ststus":  false,
			"message": " Address id doesn't exist ",
		})
		return
	}

	result := database.Db.Model(&Address).Updates(models.Address{
		Name:         body.Name,
		Phone_number: body.Phone_number,
		Pincode:      body.Pincode,
		House:        body.House_Adress,
		Area:         body.Area,
		Landmark:     body.Landmark,
		City:         body.City,
	})
	if result != nil {
		c.JSON(400, gin.H{
			"ststus":  false,
			"message": "  can,t edit your database  ",
		})
	} else {
		c.JSON(200, gin.H{
			"ststus":  true,
			"message": "Address updated",
		})
	}

}
