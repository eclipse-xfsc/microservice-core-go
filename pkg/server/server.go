package server

import (
	"fmt"
	swaggerfiles "github.com/swaggo/files"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Environment interface {
	// IsHealthy function to be used for the service readiness probe
	IsHealthy() bool
	// SetSwaggerBasePath sets the base path that will be used by swagger ui for requests url generation
	SetSwaggerBasePath(path string)
	// SwaggerOptions swagger config options. See https://github.com/swaggo/gin-swagger?tab=readme-ov-file#configuration
	SwaggerOptions() []func(config *ginSwagger.Config)
}

type HealthCheckResponse struct {
	Status bool `json:"status"`
}

type Server = GinServer

type GinServer struct {
	mode            string
	environment     Environment
	router          *gin.Engine
	routerGroups    sync.Map
	healthHandlerFn func(ctx *gin.Context)

	// initOnce uses sync.OnceFunc to call
	// GinServer.resetRoutes
	initOnce func()
}

type ServerMode string

const (
	ModeDebug      ServerMode = "debug"
	ModeProduction ServerMode = "production"
	ModeTesting    ServerMode = "testing"
)

const (
	StatusHealthy   = "healthy"
	StatusUnhealthy = "unhealthy"

	RouteParamTenantID                = ":tenantId"
	RouteParamTenantIDSwaggerNotation = "{tenantId}"

	routerGroupV1      = "v1"
	routerGroupMetrics = "metrics"
	routerGroupTenants = "tenants"
)

// New creates default gin server with basic route configuration.
func New(environment Environment, serverMode ...ServerMode) *GinServer {
	server := &GinServer{
		environment: environment,
	}
	var mode string
	if len(serverMode) == 0 {
		mode = string(ModeProduction)
	} else {
		mode = string(serverMode[0])
	}
	server.SetMode(mode)
	server.initOnce = sync.OnceFunc(server.resetRoutes)
	server.healthHandlerFn = getHealthHandler(server)

	return server
}

// getHealthHandler godoc
//
// @Summary		Basic health check
// @Description	whether service can be used. Can be used for readiness probe
// @Tags		docs
// @Produce		json
// @Success		200 {object} HealthCheckResponse
// @Router		/health [get]
// x-servers basePath=/v1/metrics
func getHealthHandler(server *GinServer) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthCheckResponse{
			Status: server.environment.IsHealthy(),
		})
	}
}

func (s *GinServer) SetMode(mode string) {
	switch ServerMode(mode) {
	case ModeProduction:
		gin.SetMode(gin.ReleaseMode)
	case ModeDebug:
		gin.SetMode(gin.DebugMode)
	case ModeTesting:
		gin.SetMode(gin.TestMode)
	default:
		panic("invalid mode specified")
	}

	s.mode = mode
}

func (s *GinServer) GetMode() string {
	return s.mode
}

// Run starts up the server on the given port (listening on all interfaces,
// if there are no explicitly defined ones). If >= 1 interfaces are specified,
// the intf will be combined with the given port and the server will listen
// on all these addresses (intf:port)
func (s *GinServer) Run(port int, interfaces ...string) error {
	s.initOnce()

	if len(interfaces) == 0 {
		return s.router.Run(fmt.Sprintf(":%d", port))
	}

	var addresses []string
	for _, currentIntf := range interfaces {
		addresses = append(addresses, fmt.Sprintf("%s:%d", currentIntf, port))
	}

	return s.router.Run(addresses...)
}

// Add runs given function to add endpoints to gin RouterGroup.
//
// Example of addFunc:
//
//	func addUserRoutes(rg *gin.RouterGroup) {
//		users := rg.Group("/users")
//
//		users.GET("/", func(c *gin.Context) {
//			c.JSON(http.StatusOK, "users")
//		})
//		users.GET("/comments", func(c *gin.Context) {
//			c.JSON(http.StatusOK, "users comments")
//		})
//		users.GET("/pictures", func(c *gin.Context) {
//			c.JSON(http.StatusOK, "users pictures")
//		})
//	}
func (s *GinServer) Add(addFunc func(tenantsGrp *gin.RouterGroup)) {
	s.initOnce()

	tenantsGrp := s.getRouterGroup(routerGroupTenants)
	if tenantsGrp == nil {
		panic("router is missing required route")
	}

	addFunc(tenantsGrp)
}

// AddHandler adds a handler for the given http method and route to the
// base group (/v1/tenants/:tenantID)
func (s *GinServer) AddHandler(method, route string, handler gin.HandlerFunc) {
	s.initOnce()

	tenantsGrp := s.getRouterGroup(routerGroupTenants)
	if tenantsGrp == nil {
		panic("router is missing required route")
	}

	tenantsGrp.Handle(method, route, handler)
}

// SetHealthHandler overwrites the current handler called for
// GET requests to /metrics/health
func (s *GinServer) SetHealthHandler(fn func(ctx *gin.Context)) {
	s.healthHandlerFn = fn
}

func (s *GinServer) resetRoutes() {
	s.router = gin.Default()

	v1 := s.router.Group("/v1")
	s.routerGroups.Store(routerGroupV1, v1)

	tenants := v1.Group(fmt.Sprintf("/tenants/%s", RouteParamTenantID))
	s.routerGroups.Store(routerGroupTenants, tenants)

	metrics := v1.Group("/metrics")
	s.routerGroups.Store(routerGroupMetrics, metrics)
	s.environment.SetSwaggerBasePath(strings.Replace(tenants.BasePath(), RouteParamTenantID, RouteParamTenantIDSwaggerNotation, 1))

	metrics.GET("/health", s.healthHandlerFn)

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, s.environment.SwaggerOptions()...))

}

func (s *GinServer) getRouterGroup(name string) *gin.RouterGroup {
	if group, isSet := s.routerGroups.Load(name); isSet {
		return group.(*gin.RouterGroup)
	}

	return nil
}
