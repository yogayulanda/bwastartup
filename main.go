package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//Connect DB
	dsn := "root:root@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		//Hentikan Program dan Munculkan Error
		log.Fatal(err.Error())
	}

	//Deklarasi Repository
	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)

	//Deklarasi Service
	campaignService := campaign.NewService(campaignRepository)
	userService := user.NewService(userRepository)
	authService := auth.NewService()

	//testing
	// campaign, err := campaignService.FindCampaigns(80)
	// fmt.Println(len(campaign))

	//Deklarasi Handler
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router := gin.Default()
	//handler images
	router.Static("/images", "./images")
	api := router.Group("/api/v1")

	// Handler user
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	// api.GET("/users/fetch", userHandler.Login)
	//Handler Campaigns
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)

	// Run Handler
	router.Run()
}

//Func Middlerware
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	//get middleware
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		tokenString := ""
		//bearer token = pisahkan sama spasi .
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		//token true . get data dalam token

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		// ambil nilai user_id
		userID := int(claim["user_id"].(float64))
		//get user by id token
		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}

//input dari user
//handler mapping input user menjadi requesr struct input
//service : melakukan mapping dari struct input ke struct user/model
//repository bikin function buat save (v)
//db save struct user ke db (v)

//middleware
// 1. ambil nilai header Authorizationc : bearer:token
// 2. ambil header Authorization ambil nilai tokennya : token
// 3. validasi token
// 4. dapat user_id dari valid token ,
// 5. ambil user dari db berdasarkan user_id lewat service
// 6. set context isinya user
