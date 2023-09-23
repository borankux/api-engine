package response

import "github.com/gin-gonic/gin"

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Response interface {
	Success(data interface{})
	Error(message string, errorCode int)
	Exception(message string, errorCode int)
	NotFound()
	UnAuthorized()
	Forbidden()
}

type GeneralResponse struct {
	Gin *gin.Context
}

func (r *GeneralResponse) Success(data interface{}) {
	r.Gin.JSON(200, ApiResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func (r *GeneralResponse) Error(message string, errorCode int) {
	r.Gin.JSON(400, ApiResponse{
		Code:    errorCode,
		Message: message,
		Data:    nil,
	})
}

func (r *GeneralResponse) Exception(message string, errorCode int) {
	r.Gin.JSON(500, ApiResponse{
		Code:    errorCode,
		Message: message,
		Data:    nil,
	})
}

func (r *GeneralResponse) NotFound() {
	r.Gin.JSON(404, ApiResponse{
		Code:    404,
		Message: "Not Found",
		Data:    nil,
	})
}

func (r *GeneralResponse) UnAuthorized() {
	r.Gin.JSON(401, ApiResponse{
		Code:    401,
		Message: "UnAuthorized",
		Data:    nil,
	})
}

func (r *GeneralResponse) Forbidden() {
	r.Gin.JSON(403, ApiResponse{
		Code:    403,
		Message: "Forbidden",
		Data:    nil,
	})
}

func R(c *gin.Context) Response {
	return &GeneralResponse{
		Gin: c,
	}
}
