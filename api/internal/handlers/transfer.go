package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strings"
	"thunes-api/errors"
	"thunes-api/internal/transfer"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/prometheus/client_golang/prometheus"
)

type transferHandler struct {
	service transfer.Service
}

func NewTransferHandler(s transfer.Service) (*transferHandler, error) {
	return &transferHandler{
		service: s,
	}, nil
}

type TransferRequest struct {
	ToAccountID int64  `json:"to_account_id" binding:"required,min=1"`
	Amount      int64  `json:"amount" binding:"required,gt=0"`
	Currency    string `json:"currency" binding:"required,currency"`
}

type ApiSuccess struct {
	Code    int
	Message string
	Details *transfer.UserAccount
}

//custom metric of type histogram
var getTransferLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_request_get_transfer_api_duration_seconds",
	Help:    "Duration of the transfer request.",
	Buckets: prometheus.DefBuckets,
}, []string{"path", "status", "method", "error"})

func init() {
	//register the counter so prometheus can collect this metric
	prometheus.MustRegister(getTransferLatency)
}

//make new transfer
func (h transferHandler) Transfer(writer http.ResponseWriter, request *http.Request) {
	now := time.Now()

	reqDump, _ := httputil.DumpRequest(request, true)
	log.Printf("REQUEST:\n%s", string(reqDump))

	//validate params
	transferReq, err := processParams(request)
	if err != nil {
		errors.JSONError(writer, err, http.StatusUnprocessableEntity)
		return
	}

	//check if user logged in
	claims, err := verifyAuth(request)
	if err != nil {
		errors.JSONError(writer, err, http.StatusUnauthorized)
		return
	}

	//validate account
	userDetails, err := h.service.ValidateAccount(request.Context(), claims.Username, transferReq.ToAccountID, transferReq.Amount, transferReq.Currency)
	if err != nil {
		fmt.Println("validate accounts err")
		//add data to prometheus
		getTransferLatency.With(prometheus.Labels{"path": request.URL.Path, "method": request.Method, "status": "422", "error": "validate accounts err"}).Observe(time.Since(now).Seconds())

		errors.JSONError(writer, err, http.StatusUnprocessableEntity)
		return
	}

	//transfer money
	arg := transfer.TransferTxParams{
		FromAccountID: userDetails.ID,
		ToAccountID:   transferReq.ToAccountID,
		Amount:        transferReq.Amount,
		Currency:      transferReq.Currency,
		Username:      claims.Username,
	}
	fmt.Println(arg)
	details, err := h.service.TransferTx(request.Context(), &arg)
	if err != nil {
		fmt.Println("transfer err")
		//add data to prometheus
		getTransferLatency.With(prometheus.Labels{"path": request.URL.Path, "method": request.Method, "status": "422", "error": "transfer err"}).Observe(time.Since(now).Seconds())

		errors.JSONError(writer, err, http.StatusUnprocessableEntity)
		return
	}

	//add data to prometheus
	//sleep(200)
	getTransferLatency.With(prometheus.Labels{"path": request.URL.Path, "method": request.Method, "status": "200", "error": "nil"}).Observe(time.Since(now).Seconds())

	//prepare output
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(http.StatusOK)
	data := ApiSuccess{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Details: details,
	}
	json.NewEncoder(writer).Encode(data)
}

//get list of all beneficiaries
func (h transferHandler) Beneficiaries(writer http.ResponseWriter, request *http.Request) {
	//check if user logged in
	claims, err := verifyAuth(request)
	if err != nil {
		errors.JSONError(writer, err, http.StatusUnauthorized)
		return
	}

	//get all beneficiaries
	benefDetails, err := h.service.GetBeneficiaries(request.Context(), claims.Username)
	fmt.Println(benefDetails)
	if err != nil {
		fmt.Println("err getting beneficiaries")
		errors.JSONError(writer, err, http.StatusUnprocessableEntity)
		return
	}

	//prepare output
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(benefDetails)
}

func verifyAuth(request *http.Request) (*Claims, error) {
	//check if user is authorized and logged in
	auth := request.Header.Get("Authorization")
	if auth == "" {
		return nil, errors.ErrUnauthorisedRequest
	}

	splitToken := strings.Split(auth, "Bearer ")
	auth = splitToken[1]

	token, err := jwt.ParseWithClaims(auth, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	//wrong token or expired
	if err != nil || !token.Valid {
		fmt.Println("invalid token")
		return nil, errors.ErrTokenExpired
	}

	claims, _ := token.Claims.(*Claims)
	return claims, nil
}

func processParams(request *http.Request) (TransferRequest, error) {
	var params TransferRequest

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		return params, errors.ErrDecodingRequest
	}

	if params.ToAccountID == 0 {
		return params, errors.ErrToAccountRequired
	}

	if params.Amount == 0 {
		return params, errors.ErrAmountRequired
	}

	if params.Amount < 0 {
		return params, errors.ErrAmountRequired
	}

	var amount interface{} = params.Amount
	if _, ok := amount.(float64); ok {
		return params, errors.ErrAmountRequired
	}

	if len(params.Currency) == 0 {
		return params, errors.ErrCurrencyRequired
	}
	return params, nil
}

func sleep(ms int) {
	rand.Seed(time.Now().UnixNano())
	now := time.Now()
	n := rand.Intn(ms + now.Second())
	time.Sleep(time.Duration(n) * time.Millisecond)
}
