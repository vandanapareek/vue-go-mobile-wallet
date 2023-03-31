package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"thunes-api/errors"
	"thunes-api/internal/users"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/prometheus/client_golang/prometheus"
)

type userHandler struct {
	service users.Service
}

func NewUserHandler(s users.Service) (*userHandler, error) {
	return &userHandler{
		service: s,
	}, nil
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var jwtKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Resp struct {
	UserInfo *users.UserInfo `json:"user_info"`
	Token    string          `json:"token"`
}

//custom metric of type counter
var userStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_get_user_status_count",
		Help: "Count of status returned by user.",
	},
	[]string{"user", "status"},
)

//custom metric of type histogram
var getLoginLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_request_get_login_api_duration_seconds",
	Help:    "Duration of the login request.",
	Buckets: []float64{0.1, 0.15, 0.2, 0.25, 0.3},
}, []string{"path", "status", "method", "error"})

func init() {
	// register the counter so prometheus can collect this metric
	prometheus.MustRegister(userStatus)
	prometheus.MustRegister(getLoginLatency)

}

func (h userHandler) Login(writer http.ResponseWriter, request *http.Request) {
	now := time.Now()

	reqDump, _ := httputil.DumpRequest(request, true)
	log.Printf("REQUEST:\n%s", string(reqDump))

	decoder := json.NewDecoder(request.Body)

	var req loginRequest
	err := decoder.Decode(&req)

	if err != nil {
		log.Printf("Error decoding request: %s", err)
		//add data to prometheus
		getLoginLatency.With(prometheus.Labels{"path": request.URL.Path, "method": request.Method, "status": "500", "error": "Error decoding request"}).Observe(time.Since(now).Seconds())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := h.service.Login(request.Context(), req.Username, req.Password)
	if err != nil {
		log.Printf("Login issue: %s", err)
		//add data to prometheus
		getLoginLatency.With(prometheus.Labels{"path": request.URL.Path, "method": request.Method, "status": "401", "error": "Login issue"}).Observe(time.Since(now).Seconds())
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("Unauthorized"))
		return
	}

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error generating JWT token %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	userInfo, err := h.service.GetUserInfo(request.Context(), claims.Username)
	if err != nil {
		fmt.Println("err getting userInfo")
		errors.JSONError(writer, err, http.StatusUnprocessableEntity)
		return
	}
	resp := Resp{
		UserInfo: userInfo,
		Token:    tokenString,
	}

	//add data to prometheus
	getLoginLatency.With(prometheus.Labels{"path": request.URL.Path, "method": request.Method, "status": "200", "error": "nil"}).Observe(time.Since(now).Seconds())

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(resp)

	log.Printf("User %s succesfully logged in", req.Username)
	log.Printf("Generated JWT Token: %s", tokenString)
}
