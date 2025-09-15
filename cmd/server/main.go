package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"github.com/ariyaagustian/gin-crud-boilerplate/internal/config"
	"github.com/ariyaagustian/gin-crud-boilerplate/internal/db"
	"github.com/ariyaagustian/gin-crud-boilerplate/internal/domain"
	"github.com/ariyaagustian/gin-crud-boilerplate/internal/repository"
	"github.com/ariyaagustian/gin-crud-boilerplate/internal/service"
	transport "github.com/ariyaagustian/gin-crud-boilerplate/internal/transport/http"
	"github.com/ariyaagustian/gin-crud-boilerplate/internal/transport/http/handler"
)

func main() {
	// load config & db
	cfg := config.Load()
	var gdb *gorm.DB = db.Open(cfg.DSN)
	if err := gdb.AutoMigrate(&domain.User{}); err != nil {
		log.Fatal("auto migrate:", err)
	}

	// wiring dependency
	v := validator.New()
	userRepo := repository.NewUserRepository(gdb)

	userSvc := service.NewUserSvc(userRepo, v)
	authSvc := service.NewAuthSvc(userRepo, v, cfg.JWTSecret, cfg.JWTAccessTTL)

	userH := handler.NewUserHandler(userSvc)
	authH := handler.NewAuthHandler(authSvc)

	// router (public + protected)
	r := transport.NewRouter(userH, authH, cfg.JWTSecret)

	log.Printf("listening at :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
