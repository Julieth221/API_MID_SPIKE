package test

import (
	"net/http"
	"net/http/httptest"

	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/routers"

	"github.com/beego/beego"
	"github.com/beego/beego/v2/core/logs"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(0)

	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))

	beego.TestBeegoInit(apppath)
}

// TestGet prueba un endpoint específico
func TestGet(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/object", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	logs.Info("Probando TestGet - Código[%d]\n%s", w.Code, w.Body.String())

	Convey("Asunto: Prueba del Endpoint de Estación", t, func() {
		Convey("El Código de Estado Debe Ser 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("El Resultado No Debe Estar Vacío", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}


