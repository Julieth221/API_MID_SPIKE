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
	

	// "github.com/astaxie/beego"
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
		// beego.NSNamespace("/usuarios",
		// 	beego.NSInclude(
		// 		&controllers.UsuarioController{},
		// 	),
		// ),

		// Rutas para el controlador de autenticación
<<<<<<< HEAD
		// beego.NSNamespace("/auth",
		// 	beego.NSInclude(
		// 		&controllers.Autenticacion_usuarioController{},
		// 	),
		// ),
=======
		beego.NSNamespace("/auth",
			beego.NSInclude(
				&controllers.Autenticacion_usuarioController{},
			),
		),

		beego.NSNamespace("/finca",
			beego.NSInclude(
				&controllers.Gestion_fincaController{},
			),
		),
>>>>>>> 78cd5e00d6587a3492928b10f609d22e433b7c16
	)
	// beego.NSNamespace("/object",

	beego.AddNamespace(ns)
}
