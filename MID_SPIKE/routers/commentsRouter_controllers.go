package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Autenticacion_usuarioController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Autenticacion_usuarioController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Autenticacion_usuarioController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Autenticacion_usuarioController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Autenticacion_usuarioController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Autenticacion_usuarioController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Autenticacion_usuarioController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Autenticacion_usuarioController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Autenticacion_usuarioController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Autenticacion_usuarioController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/nombrefinca",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"],
        beego.ControllerComments{
            Method: "GeneraryEnviarToken",
            Router: "/:correo",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"],
        beego.ControllerComments{
            Method: "ActualizarContraseña",
            Router: "/recuperar/:token",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"],
        beego.ControllerComments{
            Method: "Login",
            Router: "/sistem/login",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:UsuarioController"],
        beego.ControllerComments{
            Method: "ValidarToken",
            Router: "/validartoken",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
