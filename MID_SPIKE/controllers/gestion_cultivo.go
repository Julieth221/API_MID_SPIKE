package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	auth_JWT "github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/auth_jwt"
	"github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/services"
)

// Gestion_cultivoController operations for Gestion_cultivo
type Gestion_cultivoController struct {
	beego.Controller
}

// URLMapping ...
func (c *Gestion_cultivoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Registrar nuevo cultivo
// @Title Create
// @Description create Gestion_cultivo
// @Param	body		body 	models.Gestion_cultivo	true		"body for Gestion_cultivo content"
// @Success 201 {object} models.Gestion_cultivo
// @Failure 403 body is empty
// @router / [post]
func (c *Gestion_cultivoController) Post() {
	fmt.Println("Registrar cultivo")

	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la solicitud"}
		c.ServeJSON()
		return
	}

	// Obtener el ID del usuario desde el token JWT
	token := c.Ctx.Input.Header("Authorization")
	if token == "" {
		c.Data["json"] = map[string]interface{}{"error": "No se proporcionó token"}
		c.ServeJSON()
		return
	}

	claims, err := auth_JWT.ValidarJWT(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Token inválido o expirado"}
		c.ServeJSON()
		return
	}

	userID := claims.UserID
	fmt.Println(userID)

	//obtener el id de tipo de arroz
	tipoArrozID, err := getTipoArroz(body["tipo_arroz"].(string))
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al obtener el tipo de arroz"}
		c.ServeJSON()
		return
	}
	fmt.Println("Id de Tipo de arroz: ", tipoArrozID)

	jsonData, _ := json.Marshal(body)
	fmt.Println("Enviando datos al cultivo", string(jsonData))
	response, err := services.Metodo_post("API_CRUD_FINCA", "/v1/Registro_Cultivo", jsonData)
	if err != nil {
		c.respondWithError("Error al registrar un cultivo")
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		c.respondWithError("Error al leer respuesta de cultivo")
		return
	}

}

func getTipoArroz(nombreTipo string) (int, error) {
	response, err := services.Metodo_get("API_CRUD_CULTIVO", "/v1/tipo_arroz?query=nombre:", nombreTipo)
	if err != nil {
		return 0, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response, &data); err != nil {
		return 0, err
	}

	if resultados, ok := data["Data"].([]interface{}); ok && len(resultados) > 0 {
		tipoArroz := resultados[0].(map[string]interface{})
		return int(tipoArroz["Id"].(float64)), nil
	}
	return 0, fmt.Errorf("tipo de arroz no encontrado")

}

func getMetodoSiembra(nombreTipo string) (int, error) {
	response, err := services.Metodo_get("API_CRUD_CULTIVO", "/v1/metodo_siembra?query=nombre:", nombreTipo)
	if err != nil {
		return 0, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(response, &data); err != nil {
		return 0, err
	}

	if resultados, ok := data["Data"].([]interface{}); ok && len(resultados) > 0 {
		metodo_siembra := resultados[0].(map[string]interface{})
		return int(metodo_siembra["Id"].(float64)), nil
	}
	return 0, fmt.Errorf("Metodo de siembra no encontrado")

}

func getEstadoFenologico(nombreTipo string) (int, error) {
	response, err := services.Metodo_get("API_CRUD_FINCA", "/v1/estado_fenologico_cultivo?query=nombre:", nombreTipo)
	if err != nil {
		return 0, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response, &data); err != nil {
		return 0, err
	}

	if resultados, ok := data["Data"].([]interface{}); ok && len(resultados) > 0 {
		estado_fenologico := resultados[0].(map[string]interface{})
		return int(estado_fenologico["Id"].(float64)), nil
	}
	return 0, fmt.Errorf("Estado fenologico no encontrado")

}

// Registrar tipo Arroz
// @Title Create
// @Description create Gestion_cultivo
// @Param	body		body 	models.Gestion_cultivo	true		"body for Gestion_cultivo content"
// @Success 201 {object} models.Gestion_cultivo
// @Failure 403 body is empty
// @router tipo/arroz [post]
func (c *Gestion_cultivoController) PostTipoArroz() {
	fmt.Println("Registrar tipo de arroz")

	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la solicitud"}
		c.ServeJSON()
		return
	}

	jsonData, _ := json.Marshal(body)

	response, err := services.Metodo_post("API_CRUD_CULTIVO", "/v1/tipo_arroz", jsonData)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al crear tipo de arroz", "detalle": err.Error()}
		c.ServeJSON()
		return
	}
	fmt.Println("Respuesta de la API : ", string(response))

	c.Data["json"] = map[string]interface{}{
		"menssage": "Tipo de arroz creado con exito",
	}
	c.ServeJSON()
}

// Registrar MetodO Siembra
// @Title Create
// @Description create Gestion_cultivo
// @Param	body		body 	models.Gestion_cultivo	true		"body for Gestion_cultivo content"
// @Success 201 {object} models.Gestion_cultivo
// @Failure 403 body is empty
// @router crear/metodosiembra/ [post]
func (c *Gestion_cultivoController) PostMetodoSiembra() {
	fmt.Println("Registrar tipo de arroz")

	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la solicitud"}
		c.ServeJSON()
		return
	}

	jsonData, _ := json.Marshal(body)

	response, err := services.Metodo_post("API_CRUD_CULTIVO", "/v1/metodo_siembra", jsonData)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al crear tipo de arroz", "detalle": err.Error()}
		c.ServeJSON()
		return
	}
	fmt.Println("Respuesta de la API : ", string(response))

	c.Data["json"] = map[string]interface{}{
		"menssage": "Tipo de arroz creado con exito",
	}
	c.ServeJSON()
}

// Registrar MetodO Siembra
// @Title Create
// @Description create Gestion_cultivo
// @Param	body		body 	models.Gestion_cultivo	true		"body for Gestion_cultivo content"
// @Success 201 {object} models.Gestion_cultivo
// @Failure 403 body is empty
// @router crear/estadofenolofico/cultivo/ [post]
func (c *Gestion_cultivoController) PostEstadoFenologico() {
	fmt.Println("Registrar tipo de arroz")

	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la solicitud"}
		c.ServeJSON()
		return
	}

	jsonData, _ := json.Marshal(body)

	response, err := services.Metodo_post("API_CRUD_CULTIVO", "/v1/Estado_fenologico_cultivo", jsonData)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al crear tipo de arroz", "detalle": err.Error()}
		c.ServeJSON()
		return
	}
	fmt.Println("Respuesta de la API : ", string(response))

	c.Data["json"] = map[string]interface{}{
		"menssage": "Tipo de arroz creado con exito",
	}
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Gestion_cultivo by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Gestion_cultivo
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Gestion_cultivoController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Gestion_cultivo
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Gestion_cultivo
// @Failure 403
// @router / [get]
func (c *Gestion_cultivoController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Gestion_cultivo
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Gestion_cultivo	true		"body for Gestion_cultivo content"
// @Success 200 {object} models.Gestion_cultivo
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Gestion_cultivoController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Gestion_cultivo
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Gestion_cultivoController) Delete() {

}

func (c *Gestion_cultivoController) respondWithError(msg string, details ...string) {
	errorResponse := map[string]string{"error": msg}
	if len(details) > 0 {
		errorResponse["detalle"] = details[0]
	}
	c.Data["json"] = errorResponse
	c.ServeJSON()
}
