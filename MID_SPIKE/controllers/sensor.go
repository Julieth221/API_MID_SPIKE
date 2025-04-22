package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/services"
)

// SensorController operations for Sensor
type SensorController struct {
	beego.Controller
}

// URLMapping ...
func (c *SensorController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description Registra un nuevo tipo de sensor a través del API CRUD
// @Param	body		body 	models.TipoSensor	true		"Objeto TipoSensor a ser creado"
// @Success 201 {object} models.TipoSensor
// @Failure 400 Bad request
// @router / [post]
func (c *SensorController) Post() {
	fmt.Println("Registrando tipo de sensor a través del API CRUD")

	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la solicitud", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Convertir el body a JSON para enviarlo al API CRUD
	jsonData, err := json.Marshal(body)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"error": "Error al convertir a JSON", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Llamar al servicio Metodo_post para enviar los datos al API CRUD
	response, err := services.Metodo_post("API_CRUD_SENSOR", "/v1/Tipo_sensor", jsonData) // Ajusta el nombre del servicio y el endpoint
	if err != nil {
		c.Ctx.Output.SetStatus(500) // O el código de estado adecuado devuelto por el API CRUD
		c.Data["json"] = map[string]interface{}{"error": "Error al registrar tipo de sensor en el API CRUD", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Parsear la respuesta del API CRUD
	var respuestaCRUD map[string]interface{}
	if err := json.Unmarshal(response, &respuestaCRUD); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la respuesta del API CRUD", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Verificar el estado de la respuesta del API CRUD y devolver el resultado adecuado
	if _, ok := respuestaCRUD["error"]; ok {
		c.Ctx.Output.SetStatus(500) // Ajusta según el código de estado real del error del API CRUD
		c.Data["json"] = respuestaCRUD
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(201)                                                                                              // Created
	c.Data["json"] = map[string]interface{}{"mensaje": "Tipo de sensor registrado con éxito", "data": respuestaCRUD["Data"]} //Devuelvo la respuesta del create
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Sensor by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Sensor
// @Failure 403 :id is empty
// @router /:id [get]
func (c *SensorController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Sensor
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Sensor
// @Failure 403
// @router / [get]
func (c *SensorController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Sensor
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Sensor	true		"body for Sensor content"
// @Success 200 {object} models.Sensor
// @Failure 403 :id is not int
// @router /:id [put]
func (c *SensorController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Sensor
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *SensorController) Delete() {

}
