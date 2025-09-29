package main

import (
	"Chumber-Workflow-System/docs"
	"Chumber-Workflow-System/internal/config"
	"Chumber-Workflow-System/internal/handler"
	"Chumber-Workflow-System/internal/middleware"
	"Chumber-Workflow-System/internal/repository"
	"Chumber-Workflow-System/internal/service"
	"Chumber-Workflow-System/pkg/auth"
	"Chumber-Workflow-System/pkg/database"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize auth managers
	passwordManager := auth.NewPasswordManager(nil)
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpirationTime, cfg.JWT.Issuer)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	requestRepo := repository.NewRequestRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	branchRepo := repository.NewBranchRepository(db)
	requestTypeRepo := repository.NewRequestTypeRepository(db)
	userActivityRepo := repository.NewUserActivityRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, companyRepo, branchRepo, roleRepo, passwordManager, jwtManager)
	requestService := service.NewRequestService(requestRepo, userRepo, requestTypeRepo)
	companyService := service.NewCompanyService(companyRepo)
	branchService := service.NewBranchService(branchRepo)
	requestTypeService := service.NewRequestTypeService(requestTypeRepo)
	userActivityService := service.NewUserActivityService(userActivityRepo, userRepo)
	roleService := service.NewRoleService(permissionRepo, roleRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	requestHandler := handler.NewRequestHandler(requestService, userActivityRepo, *cfg)
	companyHandler := handler.NewCompanyHandler(companyService)
	branchHandler := handler.NewBranchHandler(branchService)
	requestTypeHandler := handler.NewRequestTypeHandler(requestTypeService)
	userActivityHandler := handler.NewUserActivityHandler(userActivityService)
	roleHandler := handler.NewRoleBasedAccessHandler(roleService)

	// Setup router

	router := setupRouter(cfg, jwtManager, userService, permissionRepo, userHandler, requestHandler, companyHandler, branchHandler, requestTypeHandler, userActivityHandler, roleHandler)

	// Create server
	server := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	log.Printf("Starting server on %s", cfg.GetServerAddress())
	log.Printf("Environment: %s", cfg.App.Environment)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRouter(
	cfg *config.Config,
	jwtManager *auth.JWTManager,
	userService *service.UserService,
	permissionRepo *repository.PermissionRepository,

	userHandler *handler.UserHandler,
	requestHandler *handler.RequestHandler,
	companyHandler *handler.CompanyHandler,
	branchHandler *handler.BranchHandler,
	requestTypeHandler *handler.RequestTypeHandler,
	userActivityHandler *handler.UserActivityHandler,
	roleHandler *handler.RoleBasedAccessHandler,
) *gin.Engine {
	// Set gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// Swagger configuration
	docs.SwaggerInfo.BasePath = "/api/v1"
	// Optional: set Title/Description/Version/Host from config if desired
	// docs.SwaggerInfo.Title = cfg.App.Name
	// docs.SwaggerInfo.Version = cfg.App.Version
	// docs.SwaggerInfo.Host = cfg.Server.Host + ":" + cfg.Server.Port

	// Swagger UI endpoint
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
			"version":   cfg.App.Version,
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")

	// Public routes (no authentication required)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	// Public request routes (no authentication required)
	publicReq := v1.Group("")
	publicReq.Use(middleware.OptionalAuth(jwtManager))
	publicReq.GET("/request/:identifier", requestHandler.GetRequestByIdentifier)
	publicReq.GET("/request/serial/:serial", requestHandler.GetRequestBySerial)

	// Protected routes (authentication required)
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	protected.Use(middleware.RequirePasswordChangeMiddleware(userService))

	// User routes
	users := protected.Group("/users")
	{
		users.GET("/profile", userHandler.GetProfile)
		users.PUT("/profile", userHandler.UpdateProfile)
		users.POST("/change-password", userHandler.ChangePassword)
		users.GET("/role-permissions", userHandler.GetCurrentUserRoleAndPermissions)
	}

	// Request routes
	requests := protected.Group("/requests")
	{
		requests.POST("", requestHandler.CreateRequest)
		requests.GET("", requestHandler.ListRequests)
		requests.GET("/count", requestHandler.ListRequestsCount)
		requests.GET("/:id", requestHandler.GetRequest)
		requests.PUT("/:id", requestHandler.UpdateRequest)
		requests.POST("/:id/approve", requestHandler.ApproveRequest)
		requests.POST("/:id/reject", requestHandler.RejectRequest)
		requests.GET("/:id/history", requestHandler.GetRequestHistory)
	}

	// Company routes (read access for authenticated users)
	companies := protected.Group("/companies")
	{
		companies.GET("", companyHandler.ListCompanies)
		companies.GET("/all", companyHandler.GetAllCompanies)
		companies.GET("/:id", companyHandler.GetCompany)
	}

	// Branch routes (read access for authenticated users)
	branches := protected.Group("/branches")
	{
		branches.GET("", branchHandler.ListBranches)
		branches.GET("/all", branchHandler.GetAllBranches)
		branches.GET("/:id", branchHandler.GetBranch)
	}

	// Request type routes (read access for authenticated users)
	requestTypes := protected.Group("/request-types")
	{
		requestTypes.GET("", requestTypeHandler.ListRequestTypes)
		requestTypes.GET("/all", requestTypeHandler.GetAllRequestTypes)
	}

	// User activity routes (read access for authenticated users)
	userActivities := protected.Group("/activities")
	{
		userActivities.GET("", userActivityHandler.GetAllActivities)
		userActivities.GET("/user/:user_id", userActivityHandler.GetActivitiesByUserID)
		userActivities.GET("/branch/:branch_id", userActivityHandler.GetActivitiesByBranchID)
		userActivities.GET("/module/:module/entity/:entity_id", userActivityHandler.GetActivitiesByModuleAndEntityID)
	}

	// Admin routes (admin only)
	admin := protected.Group("/admin")
	admin.GET("/users/member-count", middleware.RequireAccountantOrAdminOrBranchAdmin(), userHandler.GetMemberCount)
	admin.GET("/requests/ratio-sum-per-branch", middleware.RequireAdmin(), requestHandler.GetRatioSumPerBranch)
	admin.GET("requests/ratio-sum", middleware.RequireAdmin(), requestHandler.GetRatioSum)
	admin.Use(middleware.RequireAdmin())
	{
		// User management

		admin.POST("/users", userHandler.CreateUser)
		admin.GET("/users", userHandler.ListUsers)
		admin.GET("/users/:id", userHandler.GetUser)
		admin.PUT("/users/:id", userHandler.UpdateUserAdmin)
		admin.PUT("/users/:id/role", userHandler.UpdateUserRole)
		admin.DELETE("/users/:id", userHandler.DeleteUser)
		admin.POST("/users/:id/reset-password", userHandler.ResetPassword)
		//admin.GET("/users/:id/permissions", userHandler.ListUserPermissions)

		// Request management
		adminRequests := admin.Group("/requests")
		{
			adminRequests.DELETE("/:id", requestHandler.DeleteRequest)

		}
		actionRequests := protected.Group("/requests")
		{
			actionRequests.PUT("/:id/approve", middleware.RequireStaffOrAdminOrBranchAdmin(), requestHandler.ApproveRequest)
			actionRequests.PUT("/:id/reject", middleware.RequireStaffOrAdminOrBranchAdmin(), requestHandler.RejectRequest)
			actionRequests.PUT("/:id/mark-paid", middleware.RequireAccountantOrAdminOrBranchAdmin(), requestHandler.MarkAsPaid)
		}
		// Company management
		adminCompanies := admin.Group("/companies")
		{
			adminCompanies.POST("", companyHandler.CreateCompany)
			adminCompanies.PUT("/:id", companyHandler.UpdateCompany)
			adminCompanies.DELETE("/:id", companyHandler.DeleteCompany)
		}

		// Branch management
		adminBranches := admin.Group("/branches")
		{
			adminBranches.POST("", branchHandler.CreateBranch)
			adminBranches.PUT("/:id", branchHandler.UpdateBranch)
			adminBranches.DELETE("/:id", branchHandler.DeleteBranch)
		}

		// Role management
		adminRoles := admin.Group("/roles")
		{
			adminRoles.GET("", roleHandler.ListRoles)
			adminRoles.POST("", roleHandler.CreateRole)
			adminRoles.DELETE("/:id", roleHandler.DeleteRole)
			adminRoles.GET("/:id/permissions", roleHandler.GetRolePermissions)
			adminRoles.POST("/:id/permissions", roleHandler.AssignPermissionsToRole)
		}

		// Permission management
		admin.GET("/permissions", roleHandler.ListPermissions)
		admin.GET("/users/:id/permissions", roleHandler.ListUserPermissions)
	}

	return router
}
