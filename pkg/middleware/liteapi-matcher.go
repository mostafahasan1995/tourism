package middleware

import (
	"encoding/json"
	"fmt"
	liteApiSdk "larsa-tourism-microservices/liteapi-sdk"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/util"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func LiteApiMatcher(router *chi.Mux, i *do.Injector) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			method := r.Method
			path := r.URL.Path

			fmt.Println("the method", method)
			fmt.Println("the path", path)

			rctx := chi.NewRouteContext()
			result := router.Match(rctx, method, path)
			if result {
				fmt.Println("router matched")
				next.ServeHTTP(w, r)
			} else {
				fmt.Println("try matching liteapi route")
				b := strings.HasPrefix(path, "/liteapi")
				if b {

					ctx, _ := util.AddCtxAppCfg(r)

					liteApiIntFunc := do.MustInvoke[liteApiSdk.LiteApiInitFunc](i)
					liteApiSdk, err := liteApiIntFunc(ctx)
					if err != nil {
						err := map[string]any{
							"statusCode": http.StatusBadRequest,
							"message":    err.Error(),
						}
						helpers.WriteJsonCtx(ctx, w, http.StatusBadRequest, err)
					} else {
						url := r.URL.RequestURI()

						resp, err := liteApiSdk.Request(method, url, nil)
						if err != nil {
							err := map[string]any{
								"statusCode": http.StatusBadRequest,
								"message":    err.Error(),
							}
							helpers.WriteJsonCtx(ctx, w, http.StatusBadRequest, err)

						} else if resp.Status == "failed" {
							err := map[string]any{
								"statusCode": resp.Code,
								"message":    resp.Err,
							}
							helpers.WriteJsonCtx(ctx, w, resp.Code, err)
						} else {
							var result map[string]any
							err := json.Unmarshal(resp.Data, &result)
							if err != nil {
								err := map[string]any{
									"statusCode": resp.Code,
									"message":    err.Error(),
								}
								helpers.WriteJsonCtx(ctx, w, resp.Code, err)
							} else {
								helpers.WriteJsonCtx(ctx, w, resp.Code, result)
							}
						}

					}

				} else {
					w.Write([]byte("404 not found"))
				}
			}

		})
	}
}
