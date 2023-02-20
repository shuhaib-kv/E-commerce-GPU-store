package controllers

import (
	"fmt"
	"ga/initializers"
	"ga/middleware"
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

func init() {
	initializers.LoadEnvVariables()

}

var (
	accountSid string
	authToken  string
	fromPhone  string

	client *twilio.RestClient
)

func OtpLog(c *gin.Context) {
	accountSid = os.Getenv("ACCOUNT_SID")
	authToken = os.Getenv("AUTH_TOCKEN")
	fromPhone = os.Getenv("SID")
	fmt.Println(accountSid)
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	Mob := c.PostForm("number")

	result := ChekNumber(Mob)
	fmt.Println(result)

	if !result {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "Mobile number doesnt exist! Please SignUp",
		})
		return
	}
	mobile := "+91" + Mob
	params := &verify.CreateVerificationParams{}
	params.SetTo(mobile)
	params.SetChannel("sms")
	resp, err := client.VerifyV2.CreateVerification(fromPhone, params)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(400, gin.H{
			"status":  false,
			"message": "error sending OTP",
		})
	} else {
		fmt.Printf("Sent verification '%s'\n", *resp.Sid)
		c.JSON(200, gin.H{
			"status":  true,
			"message": "OTP Sent Succesfully",
			"data":    "check your phone",
		})
	}

}

func ChekNumber(str string) bool {

	mobilenumber := str
	var checkOtp models.Users
	database.Db.Raw("SELECT phone FROM users WHERE phone=?", mobilenumber).Scan(&checkOtp)
	return checkOtp.Phone == mobilenumber

}
func CheckOtp(c *gin.Context) {
	accountSid = os.Getenv("ACCOUNT_SID")
	authToken = os.Getenv("AUTH_TOCKEN")
	fromPhone = os.Getenv("SID")
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	Mob := c.PostForm("number")
	code := c.PostForm("otps")
	fmt.Println(Mob)
	ChekNumber(Mob)
	var user models.Users
	database.Db.First(&user, "phone = ?", Mob)

	mobile := "+91" + Mob
	fromPhone = os.Getenv("SID")
	fmt.Println(mobile)
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(mobile)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(fromPhone, params)

	if err != nil {
		fmt.Println(err.Error())
	} else if *resp.Status == "approved" {
		tokenstring, err := middleware.GenerateJWT(user.Email, user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create token",
			})

			return
		}
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("UserAuth", tokenstring, 3600*24*30, "", "", false, true)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "ok",
			"data":    tokenstring,
		})
	} else {
		c.JSON(404, gin.H{
			"status":  false,
			"error":   "otp is invalid",
			"message": "check otp",
		})
	}
}
