package routers

import (
	"github.com/beego/beego"
	"github.com/beego/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:SensorController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:SensorController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:SensorController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:SensorController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:SensorController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:SensorController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:SensorController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:SensorController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:SensorController"] = append(beego.GlobalControllerRouter["github.com/sena_2824182/API_MID_SPIKE/controllers:SensorController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
