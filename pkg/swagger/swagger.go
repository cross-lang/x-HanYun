package swagger

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

//go:embed swagger-ui
var swaggerUI embed.FS

//go:embed swagger.json
var swaggerJSON []byte

// RegisterRoutes 注册 Swagger UI 路由
func RegisterRoutes(server *rest.Server) {
	swaggerFS, err := fs.Sub(swaggerUI, "swagger-ui")
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(swaggerFS))

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/swagger/",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.StripPrefix("/swagger/", fileServer).ServeHTTP(w, r)
		}),
	})

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/swagger/swagger.json",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write(swaggerJSON)
		}),
	})
}
