package controller

import (
	"errors"
	"fmt"
	"github.com/borankux/resource/request"
	"github.com/borankux/resource/response"
	"github.com/borankux/resource/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type ResourceController[T any] interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Get(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
}

type ResourceConfiguration struct {
	Authenticate  bool
	AuthMethods   []string
	ParentKey     string
	ParentMethods []string
	With          []string
	WithMethods   []string
}

// MethodContains add a method to check if those array contains the given string
func (rc *ResourceConfiguration) MethodContains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

type GeneralResourceController[T any] struct {
	config *ResourceConfiguration
}

func (gc GeneralResourceController[T]) Get(c *gin.Context) {
	r := response.R(c)
	service := services.NewGeneralService[T](c)
	var with []string
	if gc.config != nil {
		if gc.config.MethodContains(gc.config.WithMethods, "GET") {
			with = gc.config.With
		}
	}
	resource, err := service.Get(with)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.NotFound()
			return
		}
		r.Error("无法获取当前资源", 40001)
		return
	}
	r.Success(resource)
}

func (gc GeneralResourceController[T]) List(c *gin.Context) {
	r := response.R(c)
	service := services.NewGeneralService[T](c)
	queryParams := c.Request.URL.Query()
	var whereClauses []string
	whereClauseString := ""
	filterValue, filterExists := queryParams["filter"]
	if !filterExists {
		for key, values := range queryParams {
			for _, value := range values {
				whereClause := fmt.Sprintf("%s='%s'", key, value)
				whereClauses = append(whereClauses, whereClause)
			}
		}
		whereClauseString = strings.Join(whereClauses, " AND ")
	}

	if filterExists {
		whereClauseString = filterValue[0]
	}

	var with []string
	if gc.config != nil {
		if gc.config.MethodContains(gc.config.WithMethods, "LIST") {
			with = gc.config.With
		}
	}
	resource, err := service.List(whereClauseString, with)
	if err != nil {
		r.Error("无法获取此资源列表", 40000)
		return
	}
	r.Success(resource)
}

func (gc GeneralResourceController[T]) Create(c *gin.Context) {
	r := response.R(c)
	v := request.V(c)
	s := services.NewGeneralService[T](c)
	data := v.Validate(new(T))
	if data == nil {
		r.Error("参数验证失败", 40001)
		return
	}

	res, err := s.Create(data.(*T))
	if err != nil {
		r.Error("创建资源失败", 40002)
		return
	}

	r.Success(res)
}

func structToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name
		if fieldName == "Model" {
			continue
		}

		if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			result[fieldName] = field.Interface()
		}
	}

	return result
}

func (gc GeneralResourceController[T]) Update(c *gin.Context) {
	r := response.R(c)
	v := request.V(c)
	s := services.NewGeneralService[T](c)
	data := v.Validate(new(T))
	if data == nil {
		r.Error("参数验证失败", 40001)
		return
	}
	res, err := s.Update(structToMap(data.(*T)))
	if err != nil {
		r.Error("创建资源失败", 40002)
		return
	}

	r.Success(res)
}

func (gc GeneralResourceController[T]) Delete(c *gin.Context) {
	r := response.R(c)
	s := services.NewGeneralService[T](c)
	err := s.Delete()
	if err != nil {
		r.Error("删除资源失败", 40002)
		return
	}
	r.Success(nil)
}

func RegisterResource[T any](name string, g *gin.RouterGroup, config ...interface{}) {
	var cfg *ResourceConfiguration
	if len(config) > 0 {
		fmt.Println(config[0])
		cfg = config[0].(*ResourceConfiguration)
	}
	gc := GeneralResourceController[T]{
		config: cfg,
	}

	if gc.config != nil {
		fmt.Println(gc.config)
	}

	g.POST(fmt.Sprintf("/%s", name), gc.Create)
	g.GET(fmt.Sprintf("/%s/:id", name), gc.Get)
	g.GET(fmt.Sprintf("/%s", name), gc.List)
	g.PUT(fmt.Sprintf("/%s/:id", name), gc.Update)
	g.DELETE(fmt.Sprintf("/%s/:id", name), gc.Delete)
}
