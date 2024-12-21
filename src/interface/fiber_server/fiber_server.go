package fiber_server

import (
	"boilerplate-demo/src/interface/fiber_server/config"
	"boilerplate-demo/src/interface/fiber_server/helper"
	"boilerplate-demo/src/interface/fiber_server/middleware"
	"boilerplate-demo/src/interface/fiber_server/route"
	player_spec "boilerplate-demo/src/interface/fiber_server/spec/player"
	system_spec "boilerplate-demo/src/interface/fiber_server/spec/system"
	"boilerplate-demo/src/use_case"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

type FiberServer struct {
	useCase *use_case.UseCase
	server  *fiber.App
	config  *config.ServerConfig
}

func New(uc *use_case.UseCase, sc *config.ServerConfig) *FiberServer {
	server := fiber.New(fiber.Config{
		CaseSensitive:         false,
		StrictRouting:         false,
		DisableStartupMessage: true,
		ReadTimeout:           30 * time.Second,
	})

	if sc.CorsAllowAll {
		server.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "*",
			AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		}))
	}

	f := &FiberServer{uc, server, sc}
	server.Use(middleware.Recovery)

	system_spec.RegisterHandlersWithOptions(server, route.NewRouteSystem(sc, uc), system_spec.FiberServerOptions{BaseURL: "/system"})

	helper.AddSwaggerUI(server, system_spec.GetSwagger, "/system")
	helper.AddSwaggerUI(server, player_spec.GetSwagger, "/v1")

	if sc.RequestLog {
		server.Use(middleware.LoggerMiddleware)
	}

	//Custom route here
	player_spec.RegisterHandlersWithOptions(server, route.NewRouteV1(uc), player_spec.FiberServerOptions{BaseURL: "/v1"})
	return f
}

func (f FiberServer) Start(wg *sync.WaitGroup) {
	wg.Add(2)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer wg.Done()
		<-exit
		zap.L().Info("Shutting down server...")

		err := f.server.Shutdown()
		if err != nil {
			zap.L().Info("Server shutdown with error", zap.Error(err))
		} else {
			zap.L().Info("Server gracefully shutdown")
		}
	}()

	go func() {
		defer wg.Done()
		zap.L().Info("Server is starting...")
		err := f.server.Listen(f.config.ListenAddress)
		if err != nil {
			zap.L().Info("Server error", zap.Error(err))
		}
		zap.L().Info("Server has been shutdown")
	}()

}
