package environment

import (
	ginSwagger "github.com/swaggo/gin-swagger"
)

type DefaultEnv struct {
	healthFunc      func() bool
	swaggerBasePath string
	swaggerOptions  []func(config *ginSwagger.Config)
}

func NewDefaultEnv() *DefaultEnv {
	return &DefaultEnv{
		healthFunc: func() bool {
			return true
		},
	}
}

func (e *DefaultEnv) IsHealthy() bool {
	return e.healthFunc()
}

func (e *DefaultEnv) SetHealthFunc(fn func() bool) {
	e.healthFunc = fn
}

func (e *DefaultEnv) SetSwaggerBasePath(path string) {
	e.swaggerBasePath = path
}

func (e *DefaultEnv) SwaggerOptions() []func(config *ginSwagger.Config) {
	return e.swaggerOptions
}
