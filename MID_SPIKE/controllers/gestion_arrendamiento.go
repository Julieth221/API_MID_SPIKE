package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/services"
)

// Gestion_arrendamientoController operations for Gestion_arrendamiento
type Gestion_arrendamientoController struct {
	beego.Controller
}

// URLMapping ...
func (c *Gestion_arrendamientoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post Registra arrendatario
// @Title Create
// @Description create Gestion_arrendamiento
// @Param	body		body 	models.Gestion_arrendamiento	true		"body for Gestion_arrendamiento content"
// @Success 201 {object} models.Gestion_arrendamiento
// @Failure 403 body is empty
// @router / [post]
func (c *Gestion_arrendamientoController) Post() {
	fmt.Println("Registrar arrendatario")

	// Leer el body de la petición
	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la solicitud"}
		c.ServeJSON()
		return
	}

	// Convertir a JSON para enviar a la API de arrendatarios
	jsonData, _ := json.Marshal(body)

	// Hacer la petición a la API CRUD de arrendatarios
	response, err := services.Metodo_post("API_CRUD_FINCA", "/v1/User_Arrendatario", jsonData)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al crear arrendatario", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	fmt.Println("Respuesta de la API:", string(response))

	// Respuesta exitosa, devolver mensaje e ID del arrendatario
	c.Data["json"] = map[string]interface{}{
		"menssage": "Arrendatario creado con éxito",
	}
	c.ServeJSON()
}

// Registrar Arredamiento
// @Title Create
// @Description create Gestion_finca
// @Param	body		body 	models.Gestion_finca	true		"body for Gestion_finca content"
// @Success 201 {object} models.Gestion_finca
// @Failure 403 body is empty
// @router /arrendamiento/ [post]
func (c *Gestion_arrendamientoController) Post_Arrendamiento() {
	fmt.Println("Registrar arrendamiento")

	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la solicitud"}
		c.ServeJSON()
		return
	}
	// Formatear las fechas antes de enviarla al API CRUD en formato timestamp
	if fechaInicio, ok := body["FechaInicio"].(string); ok {
		parsedFechaInicio, err := time.Parse("2006-01-02", fechaInicio)
		if err == nil {
			body["FechaInicio"] = parsedFechaInicio.Format(time.RFC3339Nano) // Formato con timestamp
		}
	}

	if fechaFin, ok := body["FechaFin"].(string); ok {
		parsedFechaFin, err := time.Parse("2006-01-02", fechaFin)
		if err == nil {
			body["FechaFin"] = parsedFechaFin.Format(time.RFC3339Nano) // Formato con timestamp
		}
	}

	// Convertir a JSON para enviar a la API de arrendatarios
	jsonData, _ := json.Marshal(body)
	fmt.Println("body: ", body)
	fmt.Println("json a enviar: ", string(jsonData))
	// Llamar a la API CRUD
	response, err := services.Metodo_post("API_CRUD_FINCA", "/v1/Arrendamiento", jsonData)
	if err != nil {
		c.respondWithError("Error al crear arrendamiento", err.Error())
		return
	}
	fmt.Println("Respuesta de la API:", string(response))

	// Respuesta exitosa
	c.Data["json"] = map[string]interface{}{
		"mensaje": "Arrendamiento creado con exito",
	}
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Gestion_arrendamiento by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Gestion_arrendamiento
// @Failure 403 :id is empty
// GetActivosPorFinca retorna los arrendamientos activos de una finca
// @router /activos/:id [get]
func (c *Gestion_arrendamientoController) GetActivosPorFinca() {
	fmt.Println("esta es la funcion de arrendamientos activos")
	fincaID := c.Ctx.Input.Param(":id")
	if fincaID == "" {
		c.respondWithError("El ID de la finca es requerido")
		return
	}

	response, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento", "")
	if err != nil {
		c.respondWithError("Error al obtener arrendamientos", err.Error())
		return
	}

	// Aquí se deserializa en una estructura tipo objeto que contiene un array en "Data"
	var responseWrapper struct {
		Data []map[string]interface{} `json:"Data"`
	}

	if err := json.Unmarshal(response, &responseWrapper); err != nil {
		c.respondWithError("Error al procesar la respuesta del API CRUD", err.Error())
		return
	}

	arrendamientos := responseWrapper.Data
	var activos []map[string]interface{}
	ahora := time.Now()

	for _, arr := range arrendamientos {
		fincaRef, ok := arr["FkArrendamientoFinca"].(map[string]interface{})
		if !ok {
			continue
		}
		fkFinca := fmt.Sprintf("%v", fincaRef["Id"])
		if fkFinca != fincaID {
			continue
		}

		fechaInicioStr := fmt.Sprintf("%v", arr["FechaInicio"])
		fechaFinStr := fmt.Sprintf("%v", arr["FechaFin"])

		fechaInicio, err1 := time.Parse(time.RFC3339, fechaInicioStr)
		fechaFin, err2 := time.Parse(time.RFC3339, fechaFinStr)

		activo, ok := arr["Activo"].(bool)
		if !ok {
			continue
		}

		if err1 == nil && err2 == nil && activo && ahora.After(fechaInicio) && ahora.Before(fechaFin) {
			activos = append(activos, arr)
		}
	}

	c.Data["json"] = activos
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Gestion_arrendamiento by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Gestion_arrendamiento
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Gestion_arrendamientoController) GetOne() {

}

// GetAll para arrendamientos
// @Title GetAll
// @Description get Gestion_arrendamiento
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Gestion_arrendamiento
// @Failure 403
// @router / [get]
func (c *Gestion_arrendamientoController) GetAll() {
	fmt.Println("Obteniendo lista de arrendamientos")

	// Hacer la solicitud GET al API CRUD de arrendamientos
	response, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento", "")
	if err != nil {
		c.respondWithError("Error al obtener arrendamientos", err.Error())
		return
	}

	// Convertir la respuesta en un slice de mapas
	var arrendamientos []map[string]interface{}
	if err := json.Unmarshal(response, &arrendamientos); err != nil {
		c.respondWithError("Error al procesar la respuesta del API CRUD", err.Error())
		return
	}

	// Formatear la respuesta para incluir solo los campos necesarios
	var resultado []map[string]interface{}
	for _, arr := range arrendamientos {
		finca := arr["FkArrendamientoFinca"].(map[string]interface{})["Nombre"].(string)
		parcela := arr["FkArrendatamientoParcela"].(map[string]interface{})["Nombre"].(string)
		arrendatario := arr["IdUserUserArrendatario"].(map[string]interface{})["Nombre"].(string)
		contacto := arr["IdUserUserArrendatario"].(map[string]interface{})["Contacto"].(string)

		resultado = append(resultado, map[string]interface{}{
			"Finca":         finca,
			"Parcela":       parcela,
			"Arrendatario":  arrendatario,
			"AreaArrendada": arr["FkArrendatamientoParcela"].(map[string]interface{})["Area"],
			"FechaInicio":   arr["FechaInicio"],
			"FechaFin":      arr["FechaFin"],
			"Contacto":      contacto,
			"Id":            arr["Id"],
		})
	}

	// Responder con los datos formateados
	c.Data["json"] = resultado
	c.ServeJSON()

}

// Put ...
// @Title Put
// @Description update the Gestion_arrendamiento
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Gestion_arrendamiento	true		"body for Gestion_arrendamiento content"
// @Success 200 {object} models.Gestion_arrendamiento
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Gestion_arrendamientoController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Gestion_arrendamiento
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Gestion_arrendamientoController) Delete() {

}

func (c *Gestion_arrendamientoController) respondWithError(msg string, details ...string) {
	errorResponse := map[string]string{"error": msg}
	if len(details) > 0 {
		errorResponse["detalle"] = details[0]
	}
	c.Data["json"] = errorResponse
	c.ServeJSON()
}
