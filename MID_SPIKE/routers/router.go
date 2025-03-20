// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",

		// Rutas para el controlador de usuarios
		beego.NSNamespace("/usuarios",
			beego.NSInclude(
				&controllers.UsuarioController{},
			),
		),

		// Rutas para el controlador de autenticación
		beego.NSNamespace("/auth",
			beego.NSInclude(
				&controllers.Autenticacion_usuarioController{},
			),
		),
	)
	// beego.NSNamespace("/object",

	beego.AddNamespace(ns)
}
