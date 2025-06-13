package api

import (

	// "context"
	"fmt"
	"larsa-tourism-microservices/pkg/caching"
	"larsa-tourism-microservices/pkg/db"
	"larsa-tourism-microservices/pkg/nats"

	home "larsa-tourism-microservices/pkg/services/home/di"
	messaging "larsa-tourism-microservices/pkg/services/messaging/di"
	ourService "larsa-tourism-microservices/pkg/services/our-service/di"
	statistics "larsa-tourism-microservices/pkg/services/statistics/di"

	gateway "larsa-tourism-microservices/pkg/gateway/di"
	dbsvcs "larsa-tourism-microservices/pkg/services/db/di"
	interactions "larsa-tourism-microservices/pkg/services/interactions/di"
	marketing "larsa-tourism-microservices/pkg/services/marketing/di"
	member "larsa-tourism-microservices/pkg/services/member/di"
	picklist "larsa-tourism-microservices/pkg/services/picklist/di"

	"larsa-tourism-microservices/pkg/util"
	"log"
	"net/http"
	"reflect"
	"strings"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
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

	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "x-client", "x-access-token", "x-service", "x-expire", "x-service-token", "x-user-id", "Content-Disposition"},
		ExposedHeaders:   []string{"Link", "Content-Disposition"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(common.ClientSubscriptionInstance.CheckSubscription)

	//ORDER MATTERS . DO NOT CHANGE

	//database
	do.ProvideValue(injector, conn)
	do.Provide(injector, db.NewWithTxn)

	// validator
	do.ProvideValue(injector, validateInstance)

	//messaging
	messaging.Init(injector, r)
	//getway
	gateway.Init(injector)
	//db service
	dbsvcs.Init(injector, r)

	picklist.Init(injector, r)
	home.Init(injector, r)
	member.Init(injector, r)
	ourService.Init(injector, r)
	interactions.Init(injector, r)
	marketing.Init(injector, r)
	statistics.Init(injector, r)
	//deprecated
	//travelreq.Init(injector, r)

	fmt.Println("start server")

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}

	return nil
}
