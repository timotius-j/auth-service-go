package main

import (
	"log"

	"github.com/TimX-21/auth-service-go/internal/auth/handler"
	"github.com/TimX-21/auth-service-go/internal/auth/repository"
	"github.com/TimX-21/auth-service-go/internal/auth/route"
	"github.com/TimX-21/auth-service-go/internal/auth/service"
	"github.com/TimX-21/auth-service-go/internal/config"
	"github.com/TimX-21/auth-service-go/internal/util"
	"github.com/TimX-21/auth-service-go/pkg"
)

func main() {

	err := config.InitZapSugaredLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
		return
	}

	defer config.Log.Sync()

	db, err := pkg.ConnectDB()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
		return
	}

	// config.RunMigrations(db)

	txManager := repository.NewTransactionManager(db)
	cfg := config.LoadResetConfig()
	emailSender := util.NewDummyEmailSender()

	authR := repository.NewAuthRepository(db)
	authS := service.NewAuthService(authR, txManager, cfg, emailSender)
	authH := handler.NewAuthHandler(authS)

	routeConfig := route.NewRouteConfig(
		authH,
	)

	router := route.Setup(routeConfig)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
