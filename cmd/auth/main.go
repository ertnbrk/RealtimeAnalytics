package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ertnbrk/RealtimeAnalytics/internal/auth/controller"
	"github.com/ertnbrk/RealtimeAnalytics/internal/auth/middleware"
	"github.com/ertnbrk/RealtimeAnalytics/internal/auth/repository"
	"github.com/ertnbrk/RealtimeAnalytics/internal/auth/service"
)

func main() {
	// 1) Çevresel değişkenler yerine şimdilik sabit kullanalım:
	dbURL := "postgres://go_user:go_pass@localhost:5433/analytics_db?sslmode=disable"
	jwtSecret := "changeme123"
	tokenTTL := time.Hour

	// 2) DB bağlantısı
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("DB bağlantı hatası: %v", err)
	}
	defer db.Close()

	// 3) Repo ve servis
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, jwtSecret, tokenTTL)

	// 4) Controller
	authCtr := controller.NewAuthController(authSvc)

	// 5) Gin sunucusu
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 6) Public endpoint'ler
	r.POST("/api/v1/auth/signup", authCtr.SignUpHandler)
	r.POST("/api/v1/auth/login", authCtr.LoginHandler)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// 7) JWT korumalı grup
	authorized := r.Group("/api/v1")
	authorized.Use(middleware.JWTAuthMiddleware(jwtSecret))
	{
		authorized.GET("/auth/profile", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			userRole, _ := c.Get("userRole")
			c.JSON(200, gin.H{
				"id":   userID,
				"role": userRole,
			})
		})
	}

	log.Println("AuthService 8080 portunda çalışıyor...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Sunucu başlatılamadı: %v", err)
	}
}
