package routes

import (
	"github.com/daisuzuki829/run-together-towards-goals/api"
	"github.com/daisuzuki829/run-together-towards-goals/controllers"
	"github.com/daisuzuki829/run-together-towards-goals/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"time"
)

func Handler(dbConn *gorm.DB) {

	handler := controllers.Handler{
		Db: dbConn,
	}
	apiHandler := api.Handler{
		Db: dbConn,
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "DELETE", "POST", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.LoadHTMLGlob("templates/*.html")
	r.Static("/assets", "./assets")

	// セッションの設定
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("my_session", store))

	// login
	r.POST("/login", controllers.Login)
	r.GET("/logout", controllers.Logout)

	r.GET("/registration", func(c *gin.Context) {
		rg := models.NewGenreRepository()
		genre := rg.GetAll()
		c.HTML(http.StatusOK, "registration.html", gin.H{
			"genre": genre,
		})
	})
	r.POST("/registration", controllers.NewRegistration)

	r.Group("/")
	r.Use(controllers.SessionCheck)
	{
		// user info
		r.GET("/_users", handler.GetAllUsers)
		rUser := r.Group("/user")
		{
			rUser.POST("add", handler.AddUser)
			rUser.GET("edit/:id", handler.GetUser)
			rUser.POST("edit_ok/:id", handler.EditUser)
			rUser.POST("delete/:id", handler.DeleteUser)
		}

		// genre info
		r.GET("/_genres", handler.GetAllGenres)
		rGenre := r.Group("/genre")
		{
			rGenre.POST("add", handler.AddGenre)
			rGenre.GET("edit/:id", handler.GetGenre)
			rGenre.POST("edit_ok/:id", handler.EditGenre)
			rGenre.GET("delete/:id", handler.DeleteGenre)
			rGenre.DELETE("delete/:id", handler.DeleteGenre)
		}

		// daily kpt info
		r.GET("/_daily_kpts", handler.GetAllDailyKpts)
		rDailyKpt := r.Group("/daily_kpt")
		{
			rDailyKpt.POST("add", handler.AddDailyKpt)
			rDailyKpt.POST("good/:id", handler.IncreaseGood)
			rDailyKpt.POST("fight/:id", handler.IncreaseFight)
			rDailyKpt.POST("delete/:id", handler.DeleteDailyKpt)
		}
		r.GET("/index", func(c *gin.Context) {
			c.HTML(http.StatusOK, "welcome.html", gin.H{
				"title": "title",
			})
		})

		r.NoRoute(func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/index")
		})
	}

	rApi := r.Group("/api")
	{
		rApiUser := rApi.Group("/user")
		{
			rApiUser.GET("", apiHandler.GetUser)
			rApiUser.POST("add", apiHandler.AddUser)
			rApiUser.PUT("edit", apiHandler.EditUser)
		}
		rApiDailyKpt := rApi.Group("/daily_kpt")
		{
			rApiDailyKpt.GET("", apiHandler.GetDailyKpts)
			rApiDailyKpt.POST("add", apiHandler.PostDailyKpt)
			rApiDailyKpt.PUT("good", apiHandler.IncreaseGood)
			rApiDailyKpt.PUT("fight", apiHandler.IncreaseFight)
		}
		rApiMyGoal := rApi.Group("/my_goals")
		{
			rApiMyGoal.POST("add", apiHandler.SetMyGoal)
			rApiMyGoal.PUT("edit", apiHandler.EditMyGoal)
		}
	}

	//spew.Dump(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}
