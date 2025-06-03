package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"],
        beego.ControllerComments{
            Method: "GetActivosPorFinca",
            Router: "/activos/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"],
        beego.ControllerComments{
            Method: "Post_Arrendamiento",
            Router: "/arrendamiento/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/arrendamiento/parcelas/porfinca/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"],
        beego.ControllerComments{
            Method: "PostNuevoArrendamiento",
            Router: "/arrendamiento/versionar",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"],
        beego.ControllerComments{
            Method: "GetDisponiblesPorFinca",
            Router: "/disponibles/:id/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_arrendamientoController"],
        beego.ControllerComments{
            Method: "GetParcelaPorArrendamiento",
            Router: "/parcelasarrendamiento/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_cultivoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_cultivoController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_cultivoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_cultivoController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_cultivoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_cultivoController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_cultivoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_cultivoController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_cultivoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_cultivoController"],
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
            Method: "GetFinca",
            Router: "/:buscarfinca/:porid/:id",
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
            Method: "Patch",
            Router: "/:id",
            AllowHTTPMethods: []string{"patch"},
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
            Method: "GetParcelas",
            Router: "/:parcelas/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"],
        beego.ControllerComments{
            Method: "DesactivarArrendamiento",
            Router: "/arrendamiento/desactivar/:id",
            AllowHTTPMethods: []string{"patch"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"],
        beego.ControllerComments{
            Method: "Post_Arrendatario",
            Router: "/arrendatario/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"],
        beego.ControllerComments{
            Method: "CrearParcelasParaFincaExistente",
            Router: "/crear_parcelas/finca_id",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"],
        beego.ControllerComments{
            Method: "GetDetalles",
            Router: "/detalles/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"],
        beego.ControllerComments{
            Method: "PostNuevaVersionParcela",
            Router: "/parcela/versionar/:id",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_fincaController"],
        beego.ControllerComments{
            Method: "GetTiposSueloUsuario",
            Router: "/tipos_suelo_usuario/usuario",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_insumo_cultivoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_insumo_cultivoController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_insumo_cultivoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_insumo_cultivoController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_insumo_cultivoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_insumo_cultivoController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_insumo_cultivoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_insumo_cultivoController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_insumo_cultivoController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Gestion_insumo_cultivoController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Owner_historial_parcelaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Owner_historial_parcelaController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Owner_historial_parcelaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Owner_historial_parcelaController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Owner_historial_parcelaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Owner_historial_parcelaController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Owner_historial_parcelaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Owner_historial_parcelaController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Owner_historial_parcelaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Owner_historial_parcelaController"],
        beego.ControllerComments{
            Method: "GetParcelasOriginalesPorFinca",
            Router: "/parcelas/originales/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Owner_historial_parcelaController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:Owner_historial_parcelaController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "parcelaversiones/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:SensorController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:SensorController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:SensorController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:SensorController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:SensorController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:SensorController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:SensorController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:SensorController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:SensorController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/controllers:SensorController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
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
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
