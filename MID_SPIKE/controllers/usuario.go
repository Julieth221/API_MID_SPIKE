package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/models"
	"github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/services"
)

// Usuario_controller.GoController operations for Usuario_controller.Go
type UsuarioController struct {
	beego.Controller
}

// URLMapping ...
func (c *UsuarioController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Usuario_controller.Go
// @Param	body		body 	models.Usuario_controller.Go	true		"body for Usuario_controller.Go content"
// @Success 201 {object} models.Usuario_controller.Go
// @Failure 403 body is empty
// @router / [post]
func (c *UsuarioController) Post() {
	var body_ingreso map[string]interface{}
	var alerta models.Alert
	var reponseUsuario, responseCredencial []byte
	fmt.Println("primer print: ", alerta, reponseUsuario)

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body_ingreso); err == nil {
		fmt.Println("Body que ingresa", body_ingreso)

		jsonData, err := json.MarshalIndent(body_ingreso, "", " ")
		if err != nil {
			fmt.Println("Error al convertir a JSON", err)
		}
		fmt.Println("Body de ingreso en JSON:", string(jsonData))

		jsonCredencial := map[string]interface{}{
			"contraseña": body_ingreso["contraseña"],
		}
		fmt.Println("este es el json para credenciales: ", jsonCredencial)

		jsonUsuario := map[string]interface{}{
			"nombre":             body_ingreso["Nombre"],
			"apellido":           body_ingreso["Apellido"],
			"contacto":           body_ingreso["Contacto"],
			"correo_electronico": body_ingreso["CorreoElectronico"],
		}
		fmt.Println("este es el json usuario: ", jsonUsuario)

		json_credencial_byte, _ := json.Marshal(jsonCredencial)
		// json_usuario_byte, _ := json.Marshal(jsonUsuario)

		fmt.Println("json credencial: ", string(json_credencial_byte))
		responseCredencial, _ = services.Metodo_post("API_CRUD_USUARIO", "/v1/Credenciales", json_credencial_byte)
		if err != nil {
			fmt.Println("Error al crear credenciales:", err)
			return
		}
		fmt.Println("Respuesta de la API (Credenciales): ", string(responseCredencial))

		// Obtener el ID de credencial creada
		var credencialresponse map[string]interface{}
		if err := json.Unmarshal(responseCredencial, &credencialresponse); err != nil {
			fmt.Println("Error al parsear respuesta de credenciales:", err)
			return
		}

		// Extraer el ID de la credencial
		var credencialID float64
		if data, ok := credencialresponse["Data"].(map[string]interface{}); ok {
			if id, exists := data["Id"].(float64); exists {
				credencialID = id
			} else {
				fmt.Println("Error: No se encontró el ID en la respuesta de Credenciales")
				return
			}
		} else {
			fmt.Println("Error: Estructura de respuesta de credenciales no válida")
			return
		}
		fmt.Println("Id de credenciales: ", credencialID)

		// Crear JSON para Usuario con fk_credencial
		jsonUsuario = map[string]interface{}{
			"nombre":             body_ingreso["Nombre"],
			"apellido":           body_ingreso["Apellido"],
			"contacto":           body_ingreso["Contacto"],
			"correo_electronico": body_ingreso["CorreoElectronico"],
			"FkCredencial": map[string]interface{}{
				"Id": credencialID,
			},
		}

		// Convertir a JSON y enviar POST a /v1/Usuario
		json_usuario_byte, _ := json.Marshal(jsonUsuario)
		fmt.Println("Enviando JSON a /v1/Usuario:", string(json_usuario_byte))
		reponseUsuario, err = services.Metodo_post("API_CRUD_USUARIO", "/v1/Usuario", json_usuario_byte)
		if err != nil {
			fmt.Println(" Error al crear usuario:", err)
			return
		}

		fmt.Println("Respuesta de la API (Usuario):", string(reponseUsuario))
	}

}

// GetOne ...
// @Title GetOne
// @Description get Usuario_controller.Go by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Usuario_controller.Go
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UsuarioController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Usuario_controller.Go
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Usuario_controller.Go
// @Failure 403
// @router / [get]
func (c *UsuarioController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Usuario_controller.Go
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Usuario_controller.Go	true		"body for Usuario_controller.Go content"
// @Success 200 {object} models.Usuario_controller.Go
// @Failure 403 :id is not int
// @router /:id [put]
func (c *UsuarioController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Usuario_controller.Go
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *UsuarioController) Delete() {

}
