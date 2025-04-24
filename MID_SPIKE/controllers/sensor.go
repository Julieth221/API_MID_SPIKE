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
	fmt.Println("Registrando sensor a través del API MID")

	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la solicitud", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Verificar que los campos necesarios estén presentes
	requiredFields := []string{"NombreTipoSensor", "Descripcion", "Nombre", "Ubicacion", "Cultivo", "FechaInstalacion", "Latitud", "Longitud"}
	for _, field := range requiredFields {
		if _, ok := body[field]; !ok {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = map[string]interface{}{"error": "Faltan campos requeridos", "details": fmt.Sprintf("El campo %s es obligatorio", field)}
			c.ServeJSON()
			return
		}
	}

	// Organizar los datos para el API CRUD de Tipo_sensor
	tipoSensorData := map[string]interface{}{
		"NombreTipoSensor": body["NombreTipoSensor"],
		"Descripcion":      body["Descripcion"],
	}

	// Convertir a JSON para enviar al API CRUD de Tipo_sensor
	tipoSensorJson, err := json.Marshal(tipoSensorData)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"error": "Error al convertir a JSON (TipoSensor)", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Llamar al servicio Metodo_post para registrar el tipo de sensor en el API CRUD
	tipoSensorResponse, err := services.Metodo_post("API_CRUD_SENSOR", "/v1/Tipo_sensor", tipoSensorJson)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"error": "Error al registrar tipo de sensor en el API CRUD", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Parsear la respuesta del API CRUD de Tipo_sensor
	var respuestaTipoSensor map[string]interface{}
	if err := json.Unmarshal(tipoSensorResponse, &respuestaTipoSensor); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la respuesta del API CRUD (TipoSensor)", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Verificar si hubo error al registrar el tipo de sensor
	if _, ok := respuestaTipoSensor["error"]; ok {
		c.Ctx.Output.SetStatus(500) // Ajusta el código de estado según el API CRUD
		c.Data["json"] = respuestaTipoSensor
		c.ServeJSON()
		return
	}

	// Extraer el ID del tipo de sensor registrado
	tipoSensorID, ok := respuestaTipoSensor["Data"].(map[string]interface{})["Id"].(float64)
	if !ok {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"error": "Error al obtener el ID del tipo de sensor", "details": "El API CRUD no devolvió un ID válido"}
		c.ServeJSON()
		return
	}

	// Organizar los datos para el API CRUD de GeolocalizacionSensor
	geolocalizacionData := map[string]interface{}{
		"Latitud":  body["Latitud"],
		"Longitud": body["Longitud"],
	}

	// Convertir a JSON para enviar al API CRUD de GeolocalizacionSensor
	geolocalizacionJson, err := json.Marshal(geolocalizacionData)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"error": "Error al convertir a JSON (GeolocalizacionSensor)", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Llamar al servicio Metodo_post para registrar la geolocalizacion del sensor en el API CRUD
	geolocalizacionResponse, err := services.Metodo_post("API_CRUD_SENSOR", "/v1/GeolocalizacionSensor", geolocalizacionJson)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"error": "Error al registrar geolocalizacion en el API CRUD", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Parsear la respuesta del API CRUD de GeolocalizacionSensor
	var respuestaGeolocalizacion map[string]interface{}
	if err := json.Unmarshal(geolocalizacionResponse, &respuestaGeolocalizacion); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la respuesta del API CRUD (GeolocalizacionSensor)", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Verificar si hubo error al registrar la geolocalizacion del sensor
	if _, ok := respuestaGeolocalizacion["error"]; ok {
		c.Ctx.Output.SetStatus(500) // Ajusta el código de estado según el API CRUD
		c.Data["json"] = respuestaGeolocalizacion
		c.ServeJSON()
		return
	}

	// Extraer el ID de la geolocalizacion registrada
	geolocalizacionID, ok := respuestaGeolocalizacion["Data"].(map[string]interface{})["Id"].(float64)
	if !ok {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"error": "Error al obtener el ID de la geolocalizacion", "details": "El API CRUD no devolvió un ID válido"}
		c.ServeJSON()
		return
	}

	// Organizar los datos para el registro del sensor, incluyendo la FK de Geolocalizacion y el nombre del cultivo
	sensorData := map[string]interface{}{
		"NombreSensor":       body["Nombre"],
		"FkTipoSensor":       map[string]interface{}{"Id": int(tipoSensorID)},
		"FkGeolocalizacionSensor": map[string]interface{}{"Id": int(geolocalizacionID)},
		"FechaInstalacion":   body["FechaInstalacion"],
		//  "FkRegistroCultivo":  map[string]interface{}{"Nombre": body["Cultivo"]}, //TODO: Esto es lo que voy a comentar
		"FkRegistroCultivo":  map[string]interface{}{"Id": 1}, //TODO: Esto es lo que voy a agregar para que no me de error mientras tanto
	}

	// Convertir a JSON para enviar al API CRUD de Sensor
	sensorJson, err := json.Marshal(sensorData)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"error": "Error al convertir a JSON (Sensor)", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Llamar al servicio Metodo_post para enviar los datos del sensor al API CRUD
	sensorResponse, err := services.Metodo_post("API_CRUD_SENSOR", "/v1/Sensor", sensorJson)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"error": "Error al registrar sensor en el API CRUD", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Parsear la respuesta del API CRUD de Sensor
	var respuestaSensor map[string]interface{}
	if err := json.Unmarshal(sensorResponse, &respuestaSensor); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la respuesta del API CRUD (Sensor)", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Verificar el estado de la respuesta del API CRUD y devolver el resultado adecuado
	if _, ok := respuestaSensor["error"]; ok {
		c.Ctx.Output.SetStatus(500) // Ajusta según el código de estado real del error del API CRUD
		c.Data["json"] = respuestaSensor
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = map[string]interface{}{"mensaje": "Sensor registrado con éxito", "data": respuestaSensor["Data"]} //Devuelvo la respuesta del create
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
