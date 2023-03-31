package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiError struct {
	Code    int
	Message string
}

func (err ApiError) Error() string {
	return fmt.Sprintf("error_code = %v, error_message = %v", err.Code, err.Message)
}

func CreateError(code int, message string) ApiError {
	return ApiError{Code: code, Message: message}
}

func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

var Success = ApiError{Code: 200, Message: "success"}

var ErrAccountNotFound = ApiError{Code: 422, Message: "Invalid account!"}
var ErrInvalidPayerAccount = ApiError{Code: 422, Message: "Invalid payer account!"}
var ErrInvalidBeneficiaryAccount = ApiError{Code: 422, Message: "Invalid beneficiary account!"}
var ErrNoEnoughBalance = ApiError{Code: 422, Message: "Sorry, you do not have enough balance! Please add funds."}
var ErrInvalidCurrency = ApiError{Code: 422, Message: "Mismatch currency!"}
var ErrTokenExpired = ApiError{Code: 401, Message: "The login session has been expired, please login again!"}
var ErrDecodingRequest = ApiError{Code: 422, Message: "One or few parameters are not valid!"}
var ErrUnauthorisedRequest = ApiError{Code: 401, Message: "Unauthorised request!"}
var ErrCurrencyRequired = ApiError{Code: 422, Message: "The currency is required!"}
var ErrAmountRequired = ApiError{Code: 422, Message: "Please enter valid amount!"}
var ErrToAccountRequired = ApiError{Code: 422, Message: "Please choose beneficiary!"}
var ErrCurrencyMatch = ApiError{Code: 422, Message: "Please enter valid currency, it doesnt match with existing!"}
var SystemError = ApiError{Code: 500, Message: "Something went wrong. Please try again later"}
