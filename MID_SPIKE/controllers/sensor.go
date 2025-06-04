package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	requiredFields := []string{"NombreTipoSensor", "FechaInstalacion", "Latitud", "Longitud", "FkCultivo", "FkUsuario", "IdentificadorSensor"}
	for _, field := range requiredFields {
		if _, ok := body[field]; !ok {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = map[string]interface{}{"error": "Faltan campos requeridos", "details": fmt.Sprintf("El campo %s es obligatorio", field)}
			c.ServeJSON()
			return
		}
	}

	// Validar si el tipo de sensor ya existe antes de crear uno nuevo
	nombreTipoSensor, okNombre := body["NombreTipoSensor"].(string)
	if !okNombre || nombreTipoSensor == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"error": "NombreTipoSensor inválido o ausente"}
		c.ServeJSON()
		return
	}

	// Consultar si ya existe el TipoSensor por nombre
	queryUrl := fmt.Sprintf("?query=NombreTipoSensor:%s", nombreTipoSensor)
	tipoSensorExistenteResp, err := services.Metodo_get("API_CRUD_SENSOR", "/v1/Tipo_sensor", queryUrl)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"error": "Error consultando TipoSensor existente", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Parsear la respuesta
	var tipoSensorExistente map[string]interface{}
	if err := json.Unmarshal(tipoSensorExistenteResp, &tipoSensorExistente); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"error": "Error procesando respuesta de consulta TipoSensor", "details": err.Error()}
		c.ServeJSON()
		return
	}

	var tipoSensorID float64

	// Verificar si ya existe
	crearTipoSensor := false

	if data, ok := tipoSensorExistente["Data"].([]interface{}); ok && len(data) > 0 {
		if tipoExistente, ok := data[0].(map[string]interface{}); ok {
			if id, ok := tipoExistente["Id"].(float64); ok {
				tipoSensorID = id
			} else {
				crearTipoSensor = true
			}
		} else {
			crearTipoSensor = true
		}
	} else {
		crearTipoSensor = true
	}

	if crearTipoSensor {
		tipoSensorData := map[string]interface{}{
			"NombreTipoSensor": nombreTipoSensor,
		}
		tipoSensorJson, err := json.Marshal(tipoSensorData)
		if err != nil {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = map[string]interface{}{"error": "Error al convertir a JSON (TipoSensor)", "details": err.Error()}
			c.ServeJSON()
			return
		}
		fmt.Println("json que se envia a Tipo sensor:", string(tipoSensorJson))
		tipoSensorResponse, err := services.Metodo_post("API_CRUD_SENSOR", "/v1/Tipo_sensor", tipoSensorJson)
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = map[string]interface{}{"error": "Error al registrar nuevo TipoSensor", "details": err.Error()}
			c.ServeJSON()
			return
		}

		var respuestaTipoSensor map[string]interface{}
		if err := json.Unmarshal(tipoSensorResponse, &respuestaTipoSensor); err != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = map[string]interface{}{"error": "Error procesando respuesta al registrar TipoSensor", "details": err.Error()}
			c.ServeJSON()
			return
		}

		dataTipo, ok := respuestaTipoSensor["Data"].(map[string]interface{})
		if !ok {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = map[string]interface{}{"error": "Error: la respuesta del API CRUD (TipoSensor) no contiene datos válidos"}
			c.ServeJSON()
			return
		}

		tipoSensorIDFloat, ok := dataTipo["Id"].(float64)
		if !ok {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = map[string]interface{}{"error": "Error: el ID del TipoSensor no es válido"}
			c.ServeJSON()
			return
		}
		tipoSensorID = tipoSensorIDFloat
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
	geolocalizacionResponse, err := services.Metodo_post("API_CRUD_SENSOR", "/v1/Geolocalizacion_sensor", geolocalizacionJson)
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
	geolocalizacionID, okGeolocalizacionID := respuestaGeolocalizacion["Data"].(map[string]interface{})["Id"].(float64)
	if !okGeolocalizacionID {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"error": "Error al obtener el ID de la geolocalizacion", "details": "El API CRUD no devolvió un ID válido"}
		c.ServeJSON()
		return
	}

	// Organizar los datos para el registro del sensor, incluyendo las FKs de Geolocalizacion y TipoSensor
	sensorData := map[string]interface{}{
		"FkTipoSensor":            map[string]interface{}{"Id": int(tipoSensorID)},
		"FkGeolocalizacionSensor": map[string]interface{}{"Id": int(geolocalizacionID)},
		"FechaInstalacion":        body["FechaInstalacion"],
		"FkCultivo": map[string]interface{}{
			"Id": body["FkCultivo"],
		},
		"FkUsuario":           map[string]interface{}{"Id": body["FkUsuario"]},
		"IdentificadorSensor": body["IdentificadorSensor"],
	}

	// Convertir a JSON para enviar al API CRUD de Sensor
	sensorJson, err := json.Marshal(sensorData)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"error": "Error al convertir a JSON (Sensor)", "details": err.Error()}
		c.ServeJSON()
		return
	}
	fmt.Println("json que se envia a Sensor:", string(sensorJson))
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

	// Extraer el ID del sensor registrado
	sensorID, okSensorID := respuestaSensor["Data"].(map[string]interface{})["Id"].(float64)
	if !okSensorID {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"error": "Error al obtener el ID del sensor", "details": "El API CRUD no devolvió un ID válido"}
		c.ServeJSON()
		return
	}

	// Organizar los datos para el registro del sensor_geolocalizacion
	sensorGeolocalizacionData := map[string]interface{}{
		"FkSensor":                map[string]interface{}{"Id": int(sensorID)}, // Usar el ID del sensor registrado
		"FkGeolocalizacionSensor": map[string]interface{}{"Id": int(geolocalizacionID)},
		"Activo":                  true,
	}

	// Convertir a JSON para enviar al API CRUD de SensorGeolocalizacion
	sensorGeolocalizacionJson, err := json.Marshal(sensorGeolocalizacionData)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{"error": "Error al convertir a JSON (SensorGeolocalizacion)", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Llamar al servicio Metodo_post para registrar la relación sensor-geolocalizacion en el API CRUD
	sensorGeolocalizacionResponse, err := services.Metodo_post("API_CRUD_SENSOR", "/v1/sensor_geolocalizacion", sensorGeolocalizacionJson)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"error": "Error al registrar la relación sensor-geolocalización en el API CRUD", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Parsear la respuesta del API CRUD de SensorGeolocalizacion
	var respuestaSensorGeolocalizacion map[string]interface{}
	if err := json.Unmarshal(sensorGeolocalizacionResponse, &respuestaSensorGeolocalizacion); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la respuesta del API CRUD (SensorGeolocalizacion)", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Verificar si hubo error al registrar la relación sensor-geolocalizacion
	if _, ok := respuestaSensorGeolocalizacion["error"]; ok {
		c.Ctx.Output.SetStatus(500) // Ajusta el código de estado según el API CRUD
		c.Data["json"] = respuestaSensorGeolocalizacion
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = map[string]interface{}{"mensaje": "Sensor y geolocalización registrados y relacionados con éxito", "data": respuestaSensor["Data"]} //Devuelvo la respuesta del create
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
// @Success 200 {object} []map[string]interface{}
// @Failure 500
// @router / [get]
func (c *SensorController) GetAll() {
	// falta llamar la tabla tipo_sensor
	fmt.Println("Obteniendo la relación entre sensores y geolocalizaciones")

	// Llamada al servicio para obtener la tabla de relación sensor-geolocalizacion
	sensorGeoResponse, err := services.Metodo_get("API_CRUD_SENSOR", "/v1/sensor_geolocalizacion", "")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{"error": "Error al obtener la relación sensor-geolocalización", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	// Imprimir la respuesta del servicio para inspección
	fmt.Printf("Respuesta del servicio: %s\n", sensorGeoResponse)

	// Parsear la respuesta de la tabla de relación
	var sensorGeoData map[string]interface{}
	if err := json.Unmarshal(sensorGeoResponse, &sensorGeoData); err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la respuesta de la relación sensor-geolocalización", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	// Validar que hay datos en la tabla de relación
	sensorGeolocalizaciones, okSensorGeo := sensorGeoData["Data"].([]interface{})
	if !okSensorGeo {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{"mensaje": "No hay relaciones sensor-geolocalización registradas"}
		c.ServeJSON()
		return
	}

	// Obtener la estructura de la tabla sensor_geolocalizacion desde el CRUD
	sensorGeoStructResponse, err := services.Metodo_get("API_CRUD_SENSOR", "/v1/sensor_geolocalizacion", "?limit=0") // Usar limit=0 para obtener solo la estructura
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{"error": "Error al obtener la estructura de la tabla sensor_geolocalizacion", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	// Parsear la respuesta de la estructura de la tabla
	var sensorGeoStructData map[string]interface{}
	if err := json.Unmarshal(sensorGeoStructResponse, &sensorGeoStructData); err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la respuesta de la estructura de la tabla sensor_geolocalizacion", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	// Extraer la estructura del primer elemento de los datos (asumiendo que el primer elemento tiene la estructura)
	estructura, okEstructura := sensorGeoStructData["Data"].([]interface{})[0].(map[string]interface{})
	if !okEstructura {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{"error": "Error al extraer la estructura de la tabla sensor_geolocalizacion"}
		c.ServeJSON()
		return
	}

	// Crear un mapa para almacenar los datos organizados
	respuesta := make([]map[string]interface{}, 0)

	// Iterar sobre los datos de sensor_geolocalizacion
	for _, sg := range sensorGeolocalizaciones {
		sgMap, ok := sg.(map[string]interface{})
		if !ok {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Data["json"] = map[string]interface{}{"error": "Error al procesar datos de sensor_geolocalizacion"}
			c.ServeJSON()
			return
		}

		// Crear un nuevo mapa para cada elemento, usando la estructura del CRUD
		nuevoElemento := make(map[string]interface{})
		for k := range estructura {
			// Convertir los valores al tipo adecuado basado en la estructura del CRUD
			switch v := sgMap[k].(type) {
			case float64:
				nuevoElemento[k] = int(v) // Convertir float64 a int para los campos ID
			default:
				nuevoElemento[k] = v
			}
		}
		respuesta = append(respuesta, nuevoElemento)
	}

	// Responder con la lista de relaciones sensor-geolocalización
	c.Data["json"] = map[string]interface{}{
		"mensaje":           "Lista de relaciones sensor-geolocalización",
		"data":              respuesta,
		"cantidad sensores": len(respuesta),
	}
	c.ServeJSON()
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
