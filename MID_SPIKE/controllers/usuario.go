package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
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
	var reponseUsuario, responseCredencial []byte

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
		responseCredencial, _ = services.Metodo_post("API_CRUD", "/v1/Credenciales", json_credencial_byte)
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
		reponseUsuario, err = services.Metodo_post("API_CRUD", "/v1/Usuario", json_usuario_byte)
		if err != nil {
			fmt.Println(" Error al crear usuario:", err)
			return
		}

		fmt.Println("Respuesta de la API (Usuario):", string(reponseUsuario))
	}

	c.Data["json"] = map[string]interface{}{
		"Message": "¡Usuario creado exitosamente!",
	}
	c.ServeJSON()

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
// @router /:correo [put]
func (c *UsuarioController) Put() {
	fmt.Println("Función PUT")
	var body map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		fmt.Println("Error al leer el cuerpo de la solicitud:", err)
		c.Data["json"] = map[string]interface{}{"error": "Solicitud inválida"}
		c.ServeJSON()
		return
	}
	fmt.Println("Este es el body de ingreso:", body)
	correo, ok := body["CorreoElectronico"].(string)
	if !ok {
		fmt.Println("El campo CorreoElectronico no está en el body o tiene un formato incorrecto")
		c.Data["json"] = map[string]string{"error": "Campo CorreoElectronico inválido"}
		c.ServeJSON()
		return
	}
	// Buscar usuario por correo
	queryParam := "?query=CorreoElectronico:" + correo
	responseUsuario, err := services.Metodo_get("API_CRUD", "/v1/Usuario", queryParam)
	if err != nil {
		fmt.Println("Error al obtener usuario:", err)
		c.Data["json"] = map[string]string{"error": "Error al buscar usuario"}
		c.ServeJSON()
		return
	}

	// fmt.Println("Respuesta de la API (Usuario):", string(responseUsuario))

	// Parsear la respuesta
	var resultado map[string]interface{}
	if err := json.Unmarshal(responseUsuario, &resultado); err != nil {
		fmt.Println("Error al parsear la respuesta del usuario:", err)
		c.Data["json"] = map[string]string{"error": "Error en el procesamiento de la respuesta"}
		c.ServeJSON()
		return
	}
	// Validar si el usuario existe
	data, ok := resultado["Data"].([]interface{})
	if !ok || len(data) == 0 {
		fmt.Println("El usuario con correo", correo, "no existe")
		c.Data["json"] = map[string]string{"mensaje": "El usuario no existe"}
		c.ServeJSON()
		return
	}
	fmt.Println("Validacion del usuario:", data)

	// 2. Extraer el FkCredencial del usuario
	usuario := data[0].(map[string]interface{})
	fkCredencial, ok := usuario["FkCredencial"].(map[string]interface{})
	if !ok {
		fmt.Println("El usuario no tiene FkCredencial")
		c.Data["json"] = map[string]string{"error": "El usuario no tiene credenciales asociadas"}
		c.ServeJSON()
		return
	}
	fmt.Println("ID Credencial: ", fkCredencial)

	credencialID, ok := fkCredencial["Id"]
	if !ok {
		fmt.Println("El campo Id de FkCredencial no es válido")
		c.Data["json"] = map[string]string{"error": "Id de credencial inválido"}
		c.ServeJSON()
		return
	}
	fmt.Println("ID de la credencial encontrada:", (credencialID))

	// **Generar token de recuperación**
	token, hashedToken, err := services.GenerarToken()
	if err != nil {
		fmt.Println("Error al generar el token:", err)
		c.Data["json"] = map[string]string{"error": "Error al generar el token de recuperación"}
		c.ServeJSON()
		return
	}
	fmt.Println("Token generado: ", token)

	// **Guardar el token hasheado en la base de datos de credenciales**
	updateData := map[string]interface{}{
		"token": hashedToken,
	}
	updateJSON, _ := json.Marshal(updateData)

	// **Actualizar la credencial con el token**
	_, err = services.Metodo_put("API_CRUD", "/v1/Credenciales", fmt.Sprintf("%v", credencialID), updateJSON)
	if err != nil {
		fmt.Println("Error al actualizar el token de recuperación:", err)
		c.Data["json"] = map[string]string{"error": "Error al guardar el token de recuperación"}
		c.ServeJSON()
		return
	}
	// **Enviar el correo con el token al usuario**
	err = services.EnviarCorreo(correo, token)
	if err != nil {
		fmt.Println("Error al enviar el correo:", err)
		c.Data["json"] = map[string]string{"error": "Error al enviar el correo de recuperación"}
		c.ServeJSON()
		return
	}

	// **Respuesta de éxito**
	c.Data["json"] = map[string]string{"mensaje": "Correo de recuperación enviado correctamente"}
	c.ServeJSON()
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
