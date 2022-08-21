package swagger

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

// RegisterRoute 注册swagger路由，jsonPath参数是*.swagger.json所在路径
// 访问swagger-ui界面，例如 http://127.0.0.1:8080/swagger-ui/
// 访问swagger.json，只需在swagger-ui页面填写URL，例如 http://127.0.0.1:8080/swagger/helloworld.swagger.json
func RegisterRoute(mux *http.ServeMux, jsonPath string) {
	// swagger-ui 路由
	mux.Handle("/swagger-ui/", swaggerUI("/swagger-ui/"))
	// swagger file 文档路由
	mux.HandleFunc("/swagger/", swaggerFile(jsonPath))
}

func swaggerUI(path string) http.Handler {
	fileServer := http.FileServer(&assetfs.AssetFS{
		// 把前端文件压缩成swagger包，在IDE识别不了，会被IDE标记红色或波浪线，忽略即可，可以正常编译
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "pkg/swagger-ui",
	})

	return http.StripPrefix(path, fileServer)
}

func swaggerFile(jsonPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "swagger.json") {
			http.NotFound(w, r)
			return
		}

		p := strings.TrimPrefix(r.URL.Path, "/swagger/")
		p = path.Join(jsonPath, p) // 指定*.swagger.json所在路径
		fmt.Println("swagger json file =", p)

		http.ServeFile(w, r, p)
	}
}
