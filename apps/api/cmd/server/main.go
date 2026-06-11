package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/agamtech/owncommerce/apps/api/internal/config"
	"github.com/agamtech/owncommerce/apps/api/internal/core/audit"
	"github.com/agamtech/owncommerce/apps/api/internal/core/auth"
	"github.com/agamtech/owncommerce/apps/api/internal/core/iam"
	"github.com/agamtech/owncommerce/apps/api/internal/core/tenant"
	"github.com/agamtech/owncommerce/apps/api/internal/handler"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/database"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/middleware"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/server"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL, cfg.AppEnv)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	if err := database.Seed(db); err != nil {
		log.Fatalf("seed: %v", err)
	}

	if _, err := storage.NewLocalStorage(cfg.StoragePath); err != nil {
		log.Fatalf("storage: %v", err)
	}

	tenantRepo := tenant.NewRepository(db)
	tenantSvc := tenant.NewService(tenantRepo, cfg.PlatformDomain)

	authRepo := auth.NewRepository(db)
	iamRepo := iam.NewRepository(db)
	auditSvc := audit.NewService(db)
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTAccessExpiry, cfg.JWTRefreshExpiry)
	authSvc := auth.NewService(authRepo, jwtManager, tenantSvc, iamRepo, auditSvc)

	authMiddleware := middleware.NewAuthMiddleware(jwtManager, iamRepo)
	tenantMiddleware := middleware.NewTenantMiddleware(tenantSvc)

	app := server.New()
	server.RegisterRoutes(app, server.Handlers{
		Health: handler.NewHealthHandler(db),
		Auth:   handler.NewAuthHandler(authSvc),
		Tenant: handler.NewTenantHandler(tenantSvc, tenantRepo),
		Audit:  handler.NewAuditHandler(auditSvc),
		IAM:    handler.NewIAMHandler(iamRepo),
	}, server.Middlewares{
		Auth:   authMiddleware,
		Tenant: tenantMiddleware,
	})

	go func() {
		addr := ":" + cfg.AppPort
		log.Printf("OwnCommerce API listening on %s (env=%s)", addr, cfg.AppEnv)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")
	_ = app.Shutdown()
}
