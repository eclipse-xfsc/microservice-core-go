package server

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eclipse-xfsc/microservice-core-go/pkg/server/environment"
)

var srv *GinServer

const (
	iface = "127.0.0.1"
	port  = 23438
)

func TestMain(m *testing.M) {
	srv = New(environment.NewDefaultEnv())
	go func() {
		srv.Run(port, iface)
	}()

	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	srv := New(environment.NewDefaultEnv())
	assert.IsType(t, &GinServer{}, srv)
}

func TestServer_GetMode(t *testing.T) {
	srv := New(environment.NewDefaultEnv())

	assert.Equal(t, string(ModeProduction), srv.GetMode())
}

func TestServer_SetMode(t *testing.T) {
	srv := New(environment.NewDefaultEnv())

	srv.SetMode(string(ModeDebug))

	assert.Equal(t, string(ModeDebug), srv.GetMode())
}

func TestServer_Add(t *testing.T) {

	srv.Add(func(tenantsGrp *gin.RouterGroup) {
		tenantsGrp.GET("/test-add", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusAccepted)
		})
	})

	resp, err := http.Get(fmt.Sprintf("http://%s:%d%s", iface, port, "/v1/tenants/123/test-add"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
}

func TestGinServer_AddHandler(t *testing.T) {
	const iface = "127.0.0.1"
	const port = 23438

	srv.AddHandler(http.MethodGet, "/test-add-handler", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusTeapot)
	})

	resp, err := http.Get(fmt.Sprintf("http://%s:%d%s", iface, port, "/v1/tenants/123/test-add-handler"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
}

func TestGinServer_SetHealthHandler(t *testing.T) {
	const iface = "127.0.0.1"
	const port = 23438

	srv.SetHealthHandler(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})

	resp, err := http.Get(fmt.Sprintf("http://%s:%d%s", iface, port, "/v1/metrics/health"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
