package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/services"
)

// Monitoreo_cultivo_sensorController operations for Monitoreo_cultivo_sensor
type Monitoreo_cultivo_sensorController struct {
	beego.Controller
}

// URLMapping ...
func (c *Monitoreo_cultivo_sensorController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Monitoreo_cultivo_sensor
// @Param	body		body 	models.Monitoreo_cultivo_sensor	true		"body for Monitoreo_cultivo_sensor content"
// @Success 201 {object} models.Monitoreo_cultivo_sensor
// @Failure 403 body is empty
// @router / [post]
func (c *Monitoreo_cultivo_sensorController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Monitoreo_cultivo_sensor by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Monitoreo_cultivo_sensor
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Monitoreo_cultivo_sensorController) GetOne() {

}

// GetOne ...
// @Title GetOne
// @Description get Monitoreo_cultivo_sensor by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Monitoreo_cultivo_sensor
// @Failure 403 :id is empty
// @router /SensoresPorCultivo/:idCultivo [get]
func (c *Monitoreo_cultivo_sensorController) SensoresPorCultivo() {
	// 1) Leer parámetro de ruta:
	idCultivo := c.Ctx.Input.Param(":idCultivo")
	if idCultivo == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Falta idCultivo en la ruta"}
		c.ServeJSON()
		return
	}

	// 2) Llamar al API CRUD de Sensor para obtener sensores con FkCultivo = idCultivo
	//    Ejemplo: GET /v1/Sensor?query=FkCultivo:<idCultivo>
	query := fmt.Sprintf("?query=FkCultivo:%s", idCultivo)
	respBytes, err := services.Metodo_get("API_CRUD_SENSOR", "/v1/Sensor", query)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Error consultando API CRUD Sensor", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	// 3) Parsear respuesta del CRUD Sensor
	var crudResp map[string]interface{}
	if err := json.Unmarshal(respBytes, &crudResp); err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "No se pudo parsear respuesta de API CRUD Sensor", "detalle": err.Error()}
		c.ServeJSON()
		return
	}
	dataArray, ok := crudResp["Data"].([]interface{})
	if !ok || len(dataArray) == 0 {
		// Si no hay sensores, retornamos un arreglo vacío
		c.Data["json"] = map[string]interface{}{
			"IdCultivo": idCultivo,
			"Sensores":  []interface{}{},
		}
		c.ServeJSON()
		return
	}

	// 4) Extraer una lista de identificadores de sensor (y opcionalmente su Id interno)
	type SensorInfo struct {
		IdSensor            int    `json:"IdSensor"`
		IdentificadorSensor string `json:"IdentificadorSensor"`
	}
	sensores := make([]SensorInfo, 0, len(dataArray))
	identificadores := make([]string, 0, len(dataArray))
	for _, elem := range dataArray {
		if m, ok := elem.(map[string]interface{}); ok {
			// Suponemos que el CRUD devuelve algo así:
			// { "Id": 10, "FkCultivo": 5, "IdentificadorSensor": "sensor-001", … }
			idFloat, _ := m["Id"].(float64)
			ident, _ := m["IdentificadorSensor"].(string)
			sensores = append(sensores, SensorInfo{
				IdSensor:            int(idFloat),
				IdentificadorSensor: ident,
			})
			identificadores = append(identificadores, fmt.Sprintf(`"%s"`, ident))
		}
	}

	// 5) Llamar a Firebase RTDB para traer todas las lecturas cuyos identificadores estén en esta lista
	//    Usaremos la API REST de Realtime Database:
	//
	//    GET https://<DATABASE_URL>/LecturaSensor.json?
	//         orderBy="identificador_sensor"&
	//         startAt="<sensor-001>"&endAt="<sensor-001>\uf8ff"
	//
	//    Sin embargo, Firebase REST no permite un IN() directo. La forma más sencilla:
	//    - Hacer una consulta por cada identificador (por ejemplo: orderBy="identificador_sensor"&equalTo="sensor-001")
	//    - O si quieres bajar TODAS las lecturas y filtrar en el MID (discúlpalo si la cantidad es pequeña).
	//
	//    Aquí haremos la llamada por cada identificador de forma básica. Podrías optimizar haciendo un solo GET que baje todo (por ejemplo /LecturaSensor.json) y luego filtrar en Go, si el volumen lo permite.

	// ------------------------------------------------------------------------
	// Opción A) Hacer N llamadas a Firebase RTDB (una por cada identificador)
	// ------------------------------------------------------------------------
	tipoURL := os.Getenv("FIREBASE_RTDB_URL") // por ejemplo: "https://monitoreocultivoarroz-default-rtdb.firebaseio.com"
	// secret := os.Getenv("FIREBASE_SECRET")
	lecturaPorSensor := make(map[string][]map[string]interface{})
	client := &http.Client{}

	for _, ident := range identificadores {
		// ident ya viene con comillas dobles porque lo encerramos en fmt.Sprintf(`"%s"`, ident)
		url := fmt.Sprintf("%s/LecturaSensor.json?orderBy=\"identificador_sensor\"&equalTo=%s", tipoURL, ident)

		req, _ := http.NewRequest("GET", url, nil)
		// Si tu RTDB no requiere token extra, no hace falta Authorization. Si necesitas usar el token, agrégalo aquí:
		// req.Header.Add("Authorization", "Bearer "+<tu_token_de_firebase>)

		resp, err := client.Do(req)
		if err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Data["json"] = map[string]string{"error": "Error pidiendo datos a RTDB", "detalle": err.Error()}
			c.ServeJSON()
			return
		}
		defer resp.Body.Close()

		bodyBytes, errRead := ioutil.ReadAll(resp.Body)
		if errRead != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Data["json"] = map[string]string{"error": "Error leyendo body de RTDB", "detalle": errRead.Error()}
			c.ServeJSON()
			return
		}

		// Imprime en consola para depurar la respuesta tal cual llegó:
		fmt.Printf("RTDB Response (status %d): %s\n", resp.StatusCode, string(bodyBytes))

		if resp.StatusCode != http.StatusOK {
			// Si no es 200, devolvemos el contenido a Postman para ver qué está llegando
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Data["json"] = map[string]interface{}{
				"error":    "RTDB respondió con status != 200",
				"status":   resp.StatusCode,
				"response": string(bodyBytes),
				"url":      url,
			}
			c.ServeJSON()
			return
		}

		// Aquí sí hay Status 200: intentamos parsear el JSON
		var lecturaRes map[string]map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &lecturaRes); err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Data["json"] = map[string]string{"error": "No se pudo parsear respuesta RTDB", "detalle": err.Error()}
			c.ServeJSON()
			return
		}

		// lecturaRes viene como un map de IDs automáticos de Firebase => objeto con campos:
		//   {
		//     "-ORTA4sfeUZSJg8TeNUr": { "identificador_sensor": "sensor-001", "datos_sensor": { "temperatura": 26.1, "humedad": 48 }, "fecha_lectura": "..."},
		//     "-ORTA7MYSPRRp430nfMA": { … },
		//      …
		//   }
		// Queremos sólo la lista de valores sin la key interna:
		lista := make([]map[string]interface{}, 0, len(lecturaRes))
		for _, dato := range lecturaRes {
			lista = append(lista, dato)
		}
		// Guardamos la lista en el mapa con clave “identificador_sensor” sin comillas:
		claveLimpia := strings.Trim(ident, `"`)
		lecturaPorSensor[claveLimpia] = lista
	}

	// 6) Construir la respuesta final que se devolverá al front
	//    Podríamos devolver:
	//
	//    {
	//      "IdCultivo": 123,
	//      "Sensores": [
	//         {
	//           "IdSensor": 10,
	//           "IdentificadorSensor": "sensor-001",
	//           "Lecturas": [
	//              { "fecha_lectura": "...", "datos_sensor": { "temperatura": 25.5, "humedad": 50 } },
	//              …
	//           ]
	//         },
	//         …
	//      ]
	//    }
	//
	type RespuestaSensor struct {
		IdSensor            int                      `json:"IdSensor"`
		IdentificadorSensor string                   `json:"IdentificadorSensor"`
		Lecturas            []map[string]interface{} `json:"Lecturas"`
	}
	respuestaSensores := make([]RespuestaSensor, 0, len(sensores))
	for _, s := range sensores {
		respuestaSensores = append(respuestaSensores, RespuestaSensor{
			IdSensor:            s.IdSensor,
			IdentificadorSensor: s.IdentificadorSensor,
			Lecturas:            lecturaPorSensor[s.IdentificadorSensor],
		})
	}

	// 7) Devolver JSON final
	c.Data["json"] = map[string]interface{}{
		"IdCultivo": idCultivo,
		"Sensores":  respuestaSensores,
	}
	c.ServeJSON()
}

// @Title CalcularRendimientoEstimado
// @Description Obtiene el rendimiento estimado del cultivo (en ton/m²) usando IA (Gemini) o fallback.
// @Param	idcultivo	path	string	true	"ID del cultivo"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @router prediccion/rendimientoestimado/:idcultivo [get]
func (c *Monitoreo_cultivo_sensorController) CalcularRendimientoEstimado() {
	// 1) Leer parámetro de ruta
	idCultivoStr := c.Ctx.Input.Param(":idcultivo")
	if idCultivoStr == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Falta idcultivo en la ruta"}
		c.ServeJSON()
		return
	}

	// 2) Convertir a entero
	idCultivo, err := strconv.Atoi(idCultivoStr)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "idcultivo inválido"}
		c.ServeJSON()
		return
	}

	// 3) Obtener datos del cultivo
	ruta := fmt.Sprintf("/v1/Registro_Cultivo/%d", idCultivo)
	respBytes, err := services.Metodo_get("API_CRUD_CULTIVO", ruta, "")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Error consultando API CRUD Cultivo", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	var crudResp map[string]interface{}
	if err := json.Unmarshal(respBytes, &crudResp); err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "No se pudo parsear respuesta de API CRUD Cultivo", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	// 4) Extraer "Data" en objCultivo
	var objCultivo map[string]interface{}
	switch raw := crudResp["Data"].(type) {
	case []interface{}:
		if len(raw) == 0 {
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.Data["json"] = map[string]string{"error": "RegistroCultivo no encontrado"}
			c.ServeJSON()
			return
		}
		elem, ok := raw[0].(map[string]interface{})
		if !ok {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Data["json"] = map[string]string{"error": "Formato inesperado de Data en respuesta CRUD Cultivo"}
			c.ServeJSON()
			return
		}
		objCultivo = elem

	case map[string]interface{}:
		objCultivo = raw

	default:
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Formato inesperado de Data en respuesta CRUD Cultivo"}
		c.ServeJSON()
		return
	}

	// 5) Obtener fase fenológica actual de Cultivo_Fase
	getFaseActual := func(idCultivo int) (string, error) {
		endpoint := fmt.Sprintf("?query=FkCultivoFase.Id:%d,Activo:true", idCultivo)
		respBytesFase, err := services.Metodo_get("API_CRUD_CULTIVO", "/v1/Cultivo_Fase", endpoint)
		if err != nil {
			return "", err
		}
		var respFase map[string]interface{}
		if err := json.Unmarshal(respBytesFase, &respFase); err != nil {
			return "", err
		}
		dataFases, ok := respFase["Data"].([]interface{})
		if !ok || len(dataFases) == 0 {
			return "", nil
		}

		// Buscar fase "en progreso" (Completada == false)
		var candidata map[string]interface{}
		for _, item := range dataFases {
			fm := item.(map[string]interface{})
			if completada, ok := fm["Completada"].(bool); ok && !completada {
				candidata = fm
				break
			}
		}

		// Si no hay en progreso, buscar última completada (mayor FechaFin)
		if candidata == nil {
			var ultima time.Time
			for _, item := range dataFases {
				fm := item.(map[string]interface{})
				if comp, _ := fm["Completada"].(bool); !comp {
					continue
				}
				if fs, ok := fm["FechaFin"].(string); ok {
					if f, err := time.Parse(time.RFC3339, fs); err == nil {
						if f.After(ultima) {
							ultima = f
							candidata = fm
						}
					}
				}
			}
		}

		if candidata == nil {
			return "", nil
		}
		if fkFase, ok := candidata["FkFaseCultivo"].(map[string]interface{}); ok {
			if nombre, ok2 := fkFase["NombreFase"].(string); ok2 {
				return nombre, nil
			}
		}
		return "", nil
	}

	// 6) Extraer campos del cultivo
	fkTipoMap, ok := objCultivo["FkTipoArroz"].(map[string]interface{})
	if !ok {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Falta o es inválido el campo FkTipoArroz"}
		c.ServeJSON()
		return
	}
	idTipoArrozF, ok := fkTipoMap["Id"].(float64)
	if !ok {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "FkTipoArroz.Id debe ser numérico"}
		c.ServeJSON()
		return
	}

	// Leer estado original (en caso de no haber Cultivo_Fase)
	var estadoOriginal string
	if fkEstMap, ok := objCultivo["FkEstadoFenologicoCultivo"].(map[string]interface{}); ok {
		if n, ok2 := fkEstMap["Nombre"].(string); ok2 {
			estadoOriginal = n
		}
	}

	// Determinar fase actual o fallback a estadoOriginal
	nombreEstado, err := getFaseActual(idCultivo)
	if err != nil {
		nombreEstado = estadoOriginal
	}
	if nombreEstado == "" {
		nombreEstado = estadoOriginal
	}

	densidadSiembra, ok := objCultivo["DensidadSiembra"].(float64)
	if !ok {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Falta o es inválido el campo DensidadSiembra"}
		c.ServeJSON()
		return
	}

	// 7) Fallback local: cálculo ton/ha
	productividadBase := map[int]float64{
		1: 5.0,
		2: 4.2,
		3: 4.8,
	}
	base, okBase := productividadBase[int(idTipoArrozF)]
	if !okBase {
		base = 4.5
	}

	var fd float64
	switch {
	case densidadSiembra < 100:
		fd = 0.8
	case densidadSiembra >= 100 && densidadSiembra <= 200:
		fd = 1.0
	default:
		fd = 0.9
	}

	fe := coeficienteFenologico(nombreEstado)
	rendTonHa := base * fd * fe
	if rendTonHa > base {
		rendTonHa = base
	}

	mensaje := fmt.Sprintf("%.4f ton/ha (Cálculo predeterminado)", rendTonHa)
	c.Data["json"] = map[string]interface{}{
		"rendimientoTonPorHa": rendTonHa,
		"mensaje":             mensaje,
	}
	c.ServeJSON()
}

// coeficienteFenologico retorna el factor fe según la etapa fenológica.
func coeficienteFenologico(nombreEstado string) float64 {
	switch strings.ToLower(nombreEstado) {
	case "germinación", "germinacion":
		return 0.50
	case "plántula", "plantula":
		return 0.60
	case "macollamiento":
		return 0.70
	case "elongación del tallo", "elongacion del tallo":
		return 0.75
	case "iniciación de panícula", "iniciacion de panicula":
		return 0.80
	case "floración", "floracion":
		return 0.90
	case "grano lechoso":
		return 0.95
	case "grano pastoso":
		return 1.00
	case "grano maduro":
		return 1.00
	default:
		return 1.00
	}
}

// GetAll ...
// @Title GetAll
// @Description get Monitoreo_cultivo_sensor
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Monitoreo_cultivo_sensor
// @Failure 403
// @router / [get]
func (c *Monitoreo_cultivo_sensorController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Monitoreo_cultivo_sensor
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Monitoreo_cultivo_sensor	true		"body for Monitoreo_cultivo_sensor content"
// @Success 200 {object} models.Monitoreo_cultivo_sensor
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Monitoreo_cultivo_sensorController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Monitoreo_cultivo_sensor
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Monitoreo_cultivo_sensorController) Delete() {

}
