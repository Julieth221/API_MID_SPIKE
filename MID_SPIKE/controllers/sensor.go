package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
		"FkCultivo":               body["FkCultivo"],
		"FkUsuario":               body["FkUsuario"],
		"IdentificadorSensor":     body["IdentificadorSensor"],
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
// @router /sensoresPorCultivo/:id [get]
func (c *SensorController) GetAll() {
	cultivoID := c.Ctx.Input.Param(":id")

	// 1. Obtener registro de cultivo y extraer Id_Parcela dinámicamente
	cultivoResp, err := services.Metodo_get("API_CRUD_CULTIVO", "/v1/Registro_Cultivo/", cultivoID)
	if err != nil {
		c.setError(http.StatusInternalServerError, "Error al obtener RegistroCultivo", err)
		return
	}
	var rc map[string]interface{}
	if err := json.Unmarshal(cultivoResp, &rc); err != nil {
		c.setError(http.StatusInternalServerError, "Error parsing RegistroCultivo", err)
		return
	}
	// soportar Data como objeto único o slice
	var cultMap map[string]interface{}
	switch d := rc["Data"].(type) {
	case []interface{}:
		if len(d) == 0 {
			c.setError(http.StatusNotFound, "RegistroCultivo no encontrado", nil)
			return
		}
		cultMap = d[0].(map[string]interface{})
	case map[string]interface{}:
		cultMap = d
	default:
		c.setError(http.StatusInternalServerError, "Formato de Data inesperado", nil)
		return
	}
	parcelaIDf, ok := cultMap["Id_Parcela"].(float64)
	if !ok {
		c.setError(http.StatusInternalServerError, "Campo Id_Parcela no encontrado", nil)
		return
	}
	parcelaID := int(parcelaIDf)

	// 2. Obtener FincaParcela y coordenadas del parcel_bounds
	fpResp, err := services.Metodo_get("API_CRUD_FINCA",
		fmt.Sprintf("/v1/FincaParcela?query=FkParcelaFinca:%d", parcelaID), "")
	if err != nil {
		c.setError(http.StatusInternalServerError, "Error al obtener FincaParcela", err)
		return
	}
	var fp map[string]interface{}
	if err := json.Unmarshal(fpResp, &fp); err != nil {
		c.setError(http.StatusInternalServerError, "Error parsing FincaParcela", err)
		return
	}
	fpList, ok := fp["Data"].([]interface{})
	if !ok || len(fpList) == 0 {
		c.setError(http.StatusInternalServerError, "FincaParcela sin datos", nil)
		return
	}
	fpMap := fpList[0].(map[string]interface{})
	geoObj, _ := fpMap["FkGeolocalizacion"].(map[string]interface{})
	li, _ := strconv.ParseFloat(geoObj["LatitudInicial"].(string), 64)
	lf, _ := strconv.ParseFloat(geoObj["LatitudFinal"].(string), 64)
	lj, _ := strconv.ParseFloat(geoObj["LongitudInicial"].(string), 64)
	ljf, _ := strconv.ParseFloat(geoObj["LongitudFinal"].(string), 64)
	// midLat := (li + lf) / 2
	// midLng := (lj + ljf) / 2

	// 3. Obtener sensores asociados al cultivo
	sensResp, err := services.Metodo_get("API_CRUD_SENSOR",
		fmt.Sprintf("/v1/Sensor?query=FkCultivo:%s", cultivoID), "")
	if err != nil {
		c.setError(http.StatusInternalServerError, "Error al obtener sensores", err)
		return
	}
	var sd map[string]interface{}
	if err := json.Unmarshal(sensResp, &sd); err != nil {
		c.setError(http.StatusInternalServerError, "Error parsing lista sensores", err)
		return
	}
	sensData, ok := sd["Data"].([]interface{})
	if !ok || len(sensData) == 0 {
		c.Data["json"] = map[string]interface{}{"mensaje": "No hay sensores registrados para este cultivo", "data": []interface{}{}}
		c.ServeJSON()
		return
	}

	// 4. Construir respuesta iterando dinámicamente
	var result []map[string]interface{}
	for _, item := range sensData {
		sMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		sensorID := int(sMap["Id"].(float64))
		nombre, _ := sMap["IdentificadorSensor"].(string)
		fecha, _ := sMap["FechaInstalacion"].(string)
		activo, _ := sMap["Activo"].(bool)

		tipoObj, _ := sMap["FkTipoSensor"].(map[string]interface{})
		tipoName, _ := tipoObj["NombreTipoSensor"].(string)

		// 4.1 obtener posición cardinal del sensor
		sgResp, _ := services.Metodo_get("API_CRUD_SENSOR",
			fmt.Sprintf("/v1/sensor_geolocalizacion?query=FkSensor:%d", sensorID), "")
		var sgd map[string]interface{}
		json.Unmarshal(sgResp, &sgd)

		var lat, lng float64
		if sArr, ok := sgd["Data"].([]interface{}); ok && len(sArr) > 0 {
			sg0 := sArr[0].(map[string]interface{})
			if geoObj, ok := sg0["FkGeolocalizacionSensor"].(map[string]interface{}); ok {
				lat, _ = geoObj["Latitud"].(float64)
				lng, _ = geoObj["Longitud"].(float64)
			}
		}

		estado := "inactivo"
		if activo {
			estado = "activo"
		}

		elemento := map[string]interface{}{
			"id":            sensorID,
			"nombre_sensor": nombre,
			"tipo_sensor":   tipoName,
			"ubicacionSensor": map[string]float64{
				"lat": lat,
				"lng": lng,
			},
			"fecha_instalacion": fecha,
			"estado":            estado,
			"geolocalizacionParcela": map[string]float64{
				"lat_inicial": li,
				"lng_inicial": lj,
				"lat_final":   lf,
				"lng_final":   ljf,
			},
		}
		result = append(result, elemento)
	}

	c.Data["json"] = map[string]interface{}{"mensaje": "Sensores por cultivo", "data": result}
	c.ServeJSON()
}

// GetOne ...
// @Title GetGeolocalizacionParcela
// @Description get Sensor by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Sensor
// @Failure 403 :id is empty
// @router /geolocalizacionParcela/:id [get]
func (c *SensorController) GetGeolocalizacionParcela() {
	cultivoID := c.Ctx.Input.Param(":id")

	// 1. Obtener el Registro_Cultivo para extraer dinámicamente el Id_Parcela
	cultivoResp, err := services.Metodo_get("API_CRUD_CULTIVO", "/v1/Registro_Cultivo/", cultivoID)
	if err != nil {
		c.setError(http.StatusInternalServerError, "Error al obtener RegistroCultivo", err)
		return
	}
	var rc map[string]interface{}
	if err := json.Unmarshal(cultivoResp, &rc); err != nil {
		c.setError(http.StatusInternalServerError, "Error parseando RegistroCultivo", err)
		return
	}

	// Normalizar Data (objeto único o slice)
	var cultMap map[string]interface{}
	switch d := rc["Data"].(type) {
	case []interface{}:
		if len(d) == 0 {
			c.setError(http.StatusNotFound, "RegistroCultivo no encontrado", nil)
			return
		}
		cultMap = d[0].(map[string]interface{})
	case map[string]interface{}:
		cultMap = d
	default:
		c.setError(http.StatusInternalServerError, "Formato de Data inesperado", nil)
		return
	}

	parcelaIDf, ok := cultMap["Id_Parcela"].(float64)
	if !ok {
		c.setError(http.StatusInternalServerError, "Campo Id_Parcela no encontrado", nil)
		return
	}
	parcelaID := int(parcelaIDf)

	// 2. Obtener FincaParcela y coordenadas del parcel_bounds
	fpResp, err := services.Metodo_get("API_CRUD_FINCA",
		fmt.Sprintf("/v1/FincaParcela?query=FkParcelaFinca:%d", parcelaID), "")
	if err != nil {
		c.setError(http.StatusInternalServerError, "Error al obtener FincaParcela", err)
		return
	}
	var fp map[string]interface{}
	if err := json.Unmarshal(fpResp, &fp); err != nil {
		c.setError(http.StatusInternalServerError, "Error parseando FincaParcela", err)
		return
	}
	fpList, ok := fp["Data"].([]interface{})
	if !ok || len(fpList) == 0 {
		c.setError(http.StatusNotFound, "FincaParcela sin datos", nil)
		return
	}
	fpMap := fpList[0].(map[string]interface{})
	geoObj := fpMap["FkGeolocalizacion"].(map[string]interface{})
	parObj := fpMap["FkParcelaFinca"].(map[string]interface{})

	li, _ := strconv.ParseFloat(geoObj["LatitudInicial"].(string), 64)
	lf, _ := strconv.ParseFloat(geoObj["LatitudFinal"].(string), 64)
	lj, _ := strconv.ParseFloat(geoObj["LongitudInicial"].(string), 64)
	ljf, _ := strconv.ParseFloat(geoObj["LongitudFinal"].(string), 64)

	// 3. Obtener sensores asociados al cultivo
	sensResp, err := services.Metodo_get("API_CRUD_SENSOR",
		fmt.Sprintf("/v1/Sensor?query=FkCultivo:%s", cultivoID), "")
	if err != nil {
		c.setError(http.StatusInternalServerError, "Error al obtener sensores", err)
		return
	}
	var sd map[string]interface{}
	if err := json.Unmarshal(sensResp, &sd); err != nil {
		c.setError(http.StatusInternalServerError, "Error parseando lista sensores", err)
		return
	}
	sensData, ok := sd["Data"].([]interface{})
	if !ok {
		c.setError(http.StatusInternalServerError, "Formato de datos de sensores inesperado", nil)
		return
	}

	// 4. Construir respuesta incluyendo campos extra
	var sensors []map[string]interface{}
	for _, item := range sensData {
		sMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		sensorID := int(sMap["Id"].(float64))
		identificador, _ := sMap["IdentificadorSensor"].(string)
		activo, _ := sMap["Activo"].(bool)

		tipoObj, _ := sMap["FkTipoSensor"].(map[string]interface{})
		tipoName, _ := tipoObj["NombreTipoSensor"].(string)

		// Obtener geolocalización del sensor
		sgResp, _ := services.Metodo_get("API_CRUD_SENSOR",
			fmt.Sprintf("/v1/sensor_geolocalizacion?query=FkSensor:%d", sensorID), "")
		var sgd map[string]interface{}
		json.Unmarshal(sgResp, &sgd)

		var lat, lng float64
		if sArr, ok := sgd["Data"].([]interface{}); ok && len(sArr) > 0 {
			sg0 := sArr[0].(map[string]interface{})
			if geoS, ok := sg0["FkGeolocalizacionSensor"].(map[string]interface{}); ok {
				lat = geoS["Latitud"].(float64)
				lng = geoS["Longitud"].(float64)
			}
		}

		estado := "inactivo"
		if activo {
			estado = "activo"
		}

		sensors = append(sensors, map[string]interface{}{
			"idSensor":            sensorID,
			"identificadorSensor": identificador,
			"tipo_sensor":         tipoName,
			"estado":              estado,
			"ubicacionSensor": map[string]float64{
				"lat": lat,
				"lng": lng,
			},
		})
	}

	// 5. Responder con geolocalización de parcela y lista de sensores
	respuesta := map[string]interface{}{
		"mensaje":       "Geolocalización de la parcela y sensores",
		"nombreParcela": parObj["NombreParcela"],
		"geolocalizacionParcela": map[string]float64{
			"lat_inicial": li,
			"lng_inicial": lj,
			"lat_final":   lf,
			"lng_final":   ljf,
		},
		"sensors": sensors,
	}

	c.Data["json"] = respuesta
	c.ServeJSON()
}

// setError helper para manejar errores
func (c *SensorController) setError(code int, msg string, err error) {
	c.Ctx.Output.SetStatus(code)
	c.Data["json"] = map[string]interface{}{"error": msg, "detalle": errString(err)}
	c.ServeJSON()
}

func errString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
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
