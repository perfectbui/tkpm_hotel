package delivery

import (
	"TKPM/common/errors"
	"TKPM/internals/models"
	"TKPM/utils"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-gorilla/gorilla"
	"github.com/opentracing/opentracing-go"
)

// NewHTTPHandler ...
func NewHTTPHandler(roomDelivery RoomDelivery, accountDelivery AccountDelivery, contractDelivery ContractDelivery, tracer opentracing.Tracer) http.Handler {
	router := mux.NewRouter()

	// account
	router.HandleFunc("/api/v1/sign-up", accountDelivery.SignUp).Methods("POST")

	router.HandleFunc("/api/v1/sign-in", accountDelivery.SignIn).Methods("POST")

	router.HandleFunc("/api/v1/account", accountDelivery.GetAccountById).Methods("GET")

	// room
	router.HandleFunc("/api/v1/room", http.HandlerFunc(roomDelivery.Create)).Methods("POST")

	router.HandleFunc("/api/v1/room", Adapt(http.HandlerFunc(roomDelivery.Update),
		CheckAuth(accountDelivery)).ServeHTTP).Methods("PUT")

	router.HandleFunc("/api/v1/room/list", Adapt(http.HandlerFunc(roomDelivery.GetRoomList),
		CheckAuth(accountDelivery)).ServeHTTP).Methods("GET")

	router.HandleFunc("/api/v1/room", Adapt(http.HandlerFunc(roomDelivery.GetRoomById),
		CheckAuth(accountDelivery)).ServeHTTP).Methods("GET")

	router.HandleFunc("/api/v1/upload", Adapt(http.HandlerFunc(roomDelivery.UploadImage),
		CheckAuth(accountDelivery)).ServeHTTP).Methods("POST")

	// contract
	router.HandleFunc("/api/v1/contract", Adapt(http.HandlerFunc(contractDelivery.Create),
		CheckAuth(accountDelivery)).ServeHTTP).Methods("POST")

	router.HandleFunc("/api/v1/contract", Adapt(http.HandlerFunc(contractDelivery.Update),
		CheckAuth(accountDelivery)).ServeHTTP).Methods("PUT")

	router.HandleFunc("/api/v1/contract/list", Adapt(http.HandlerFunc(contractDelivery.GetContractList),
		CheckAuth(accountDelivery)).ServeHTTP).Methods("GET")

	router.HandleFunc("/api/v1/contract", Adapt(http.HandlerFunc(contractDelivery.GetContractById),
		CheckAuth(accountDelivery)).ServeHTTP).Methods("GET")

	_ = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		route.Handler(
			gorilla.Middleware(tracer, route.GetHandler()))
		return nil
	})

	// subRouter.Use(mux.MiddlewareFunc(CheckAuth(accountDelivery)))

	return router
}

type Adapter func(http.Handler) http.Handler

func Adapt(handler http.Handler, adapters ...Adapter) http.Handler {
	for i := len(adapters); i > 0; i-- {
		handler = adapters[i-1](handler)
	}
	return handler
}

func AllowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
}

func CheckAuth(accountDelivery AccountDelivery) Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token := req.Header.Get("Authorization")
			sDec, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				utils.ResponseWithJson(w, http.StatusForbidden, map[string]string{"message": errors.ErrInvalidToken.Error()})
				return
			}

			var account models.Account
			err = json.Unmarshal([]byte(string(sDec)), &account)
			if err != nil {
				utils.ResponseWithJson(w, http.StatusForbidden, map[string]string{"message": errors.ErrInvalidToken.Error()})
				return
			}

			_, err = accountDelivery.CheckAuth(string(account.AccountID))
			if err != nil {
				utils.ResponseWithJson(w, http.StatusForbidden, map[string]string{"message": errors.ErrInvalidToken.Error()})
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

func CheckRole(accountDelivery AccountDelivery, role string) Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token := req.Header.Get("Authorization")
			sDec, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				utils.ResponseWithJson(w, http.StatusForbidden, map[string]string{"message": errors.ErrInvalidToken.Error()})
				return
			}

			var account models.Account
			err = json.Unmarshal([]byte(string(sDec)), &account)
			if err != nil {
				utils.ResponseWithJson(w, http.StatusForbidden, map[string]string{"message": errors.ErrInvalidToken.Error()})
				return
			}

			_, err = accountDelivery.CheckAuth(string(account.AccountID))
			if err != nil {
				utils.ResponseWithJson(w, http.StatusForbidden, map[string]string{"message": errors.ErrInvalidToken.Error()})
				return
			}

			if role != account.Role {
				utils.ResponseWithJson(w, http.StatusForbidden, map[string]string{"message": errors.ErrNoPermission.Error()})
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

func Tracing() Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// tracer := opentracing.GlobalTracer()
			// spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
			// serverSpan := tracer.StartSpan("hotel root", ext.RPCServerOption(spanCtx))
			// defer serverSpan.Finish()
			next.ServeHTTP(w, req)
		})
	}

}
