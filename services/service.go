package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RequestParams interface{}
type Service[T any] interface {
	Create(data *T) (T, error)
	Update(params map[string]interface{}) (*T, error)
	Get(with []string) (T, error)
	Delete() error
	List(conditions interface{}, with []string) ([]T, error)
}

type GeneralService[T any] struct {
	Gin *gin.Context
}

func (s *GeneralService[T]) Create(data *T) (*T, error) {
	db := s.Gin.MustGet("database").(*gorm.DB)
	tx := db.Create(data)
	return data, tx.Error
}

func (s *GeneralService[T]) Update(params map[string]interface{}) (*T, error) {
	db := s.Gin.MustGet("database").(*gorm.DB)
	model := new(T)
	id := s.Gin.Param("id")
	tx := db.Model(model).First(model, id)
	if tx.Error != nil {
		return model, tx.Error
	}
	tx = db.Model(model).Updates(params)
	return model, tx.Error
}

func (s *GeneralService[T]) Get(with []string) (*T, error) {
	db := s.Gin.MustGet("database").(*gorm.DB)
	var model = new(T)
	id := s.Gin.Param("id")
	tx := db.Debug().Model(model)
	fmt.Println(with)
	if len(with) > 0 {
		for _, v := range with {
			tx = tx.Preload(v)
		}
	}
	tx = tx.First(model, id)
	if tx.Error != nil {
		return model, tx.Error
	}
	return model, nil
}

func (s *GeneralService[T]) Delete() error {
	db := s.Gin.MustGet("database").(*gorm.DB)
	var model = new(T)
	id := s.Gin.Param("id")
	tx := db.Model(model).Unscoped().Delete(model, id)
	return tx.Error
}

func (s *GeneralService[T]) List(conditions interface{}, with []string) ([]T, error) {
	tx := s.Gin.MustGet("database").(*gorm.DB)
	var model []T
	tx = tx.Debug()
	tx = tx.Model(model)
	tx = tx.Where(conditions)
	fmt.Println(with)
	if len(with) > 0 {
		for _, v := range with {
			tx = tx.Preload(v)
		}
	}
	tx = tx.Find(&model)
	return model, tx.Error
}

func NewGeneralService[T any](c *gin.Context) *GeneralService[T] {
	return &GeneralService[T]{
		Gin: c,
	}
}
