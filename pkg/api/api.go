package api

import (

	// "context"
	"context"
	"fmt"
	liteApiSdk "larsa-tourism-microservices/liteapi-sdk"
	"larsa-tourism-microservices/pkg/caching"
	"larsa-tourism-microservices/pkg/db"
	"larsa-tourism-microservices/pkg/nats"
	"os"
	"os/signal"
	"syscall"
	"time"

	home "larsa-tourism-microservices/pkg/services/home/di"
	messaging "larsa-tourism-microservices/pkg/services/messaging/di"
	ourService "larsa-tourism-microservices/pkg/services/our-service/di"
	statistics "larsa-tourism-microservices/pkg/services/statistics/di"

	gateway "larsa-tourism-microservices/pkg/gateway/di"
	dbsvcs "larsa-tourism-microservices/pkg/services/db/di"
	exhibitionmanagement "larsa-tourism-microservices/pkg/services/exhibition-management/di"
	interactions "larsa-tourism-microservices/pkg/services/interactions/di"
	liteapiPkg "larsa-tourism-microservices/pkg/services/liteapi"
	liteapi "larsa-tourism-microservices/pkg/services/liteapi/di"
	marketing "larsa-tourism-microservices/pkg/services/marketing/di"
	member "larsa-tourism-microservices/pkg/services/member/di"
	picklist "larsa-tourism-microservices/pkg/services/picklist/di"

	//
	transtest "larsa-tourism-microservices/pkg/services/trans-test/di"

	// Import for custom validation
	picklistModels "larsa-tourism-microservices/pkg/services/picklist/models"

	"larsa-tourism-microservices/pkg/util"
	"log"
	"net/http"
	"reflect"
	"strings"

	middlewares "larsa-tourism-microservices/pkg/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

func Start() error {
	port := util.GetEnv("PORT", "3277")

	fmt.Println("start server on port:", port)

	injector := do.New()

	conn, err := db.InitDB()
	if err != nil {
		return err
	}

	caching.InitRedis()
	nats.InitNats()

	r := chi.NewRouter()
	validateInstance := validator.New()
	validateInstance.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validations
	picklistModels.RegisterCustomValidations(validateInstance)

	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5))

	//r.Match(chi.RouteContext(),"","")

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "x-client", "x-access-token", "x-service", "x-expire", "x-service-token", "x-user-id", "Content-Disposition"},
		ExposedHeaders:   []string{"Link", "Content-Disposition"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middlewares.LiteApiMatcher(r, injector))

	//r.Use(common.ClientSubscriptionInstance.CheckSubscription)

	//ORDER MATTERS . DO NOT CHANGE

	//database
	do.ProvideValue(injector, conn)
	do.Provide(injector, db.NewWithTxn)

	// validator
	do.ProvideValue(injector, validateInstance)

	do.ProvideValue(injector, liteApiSdk.LiteApiSdkInit())

	//messaging
	messaging.Init(injector, r)
	//getway
	gateway.Init(injector)
	//db service
	dbsvcs.Init(injector, r)

	// interactions must be initialized before modules that depend on FaveSvcs
	interactions.Init(injector, r)

	picklist.Init(injector, r)
	home.Init(injector, r)
	member.Init(injector, r)
	ourService.Init(injector, r)
	marketing.Init(injector, r)
	exhibitionmanagement.Init(injector, r)
	statistics.Init(injector, r)
	transtest.Init(injector, r)
	liteapi.Init(injector, r)

	updater := liteapiPkg.NewHotelsUpdater(
		injector,
		[]string{
			"d0lfvck5drjs739an280",
		},
	)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	fmt.Println("start server")

	// Channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := updater.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for signal
	<-stop
	log.Println("Shutting down server...")

	updater.Stop()

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server exited gracefully")

	return nil
}
