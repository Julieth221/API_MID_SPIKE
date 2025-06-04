// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		// Rutas para el controlador de sensores
		beego.NSNamespace("/sensores",
			beego.NSInclude(
				&controllers.SensorController{},
			),
		),
		// Rutas para el controlador de usuarios
		beego.NSNamespace("/usuarios",
			beego.NSInclude(
				&controllers.UsuarioController{},
			),
		),
		// Rutas para el controlador de gestión de finca
		beego.NSNamespace("/finca",
			beego.NSInclude(
				&controllers.Gestion_fincaController{},
			),
		),

		beego.NSNamespace("/arrendamiento",
			beego.NSInclude(
				&controllers.Gestion_arrendamientoController{},
			),
		),

		beego.NSNamespace("/historial_parcela",
			beego.NSInclude(
				&controllers.Owner_historial_parcelaController{},
			),
		),
	)

	beego.AddNamespace(ns)
}
