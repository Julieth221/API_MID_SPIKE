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

	// Formatear fechas
	if fechaInicio, ok := body["FechaInicio"].(string); ok {
		parsed, err := time.Parse("2006-01-02", fechaInicio)
		if err == nil {
			body["FechaInicio"] = parsed.Format(time.RFC3339Nano)
		}
	}
	if fechaFin, ok := body["FechaFin"].(string); ok {
		parsed, err := time.Parse("2006-01-02", fechaFin)
		if err == nil {
			body["FechaFin"] = parsed.Format(time.RFC3339Nano)
		}
	}

	// Extraer las parcelas antes de hacer el POST principal
	parcelas, ok := body["Parcelas"].([]interface{})
	if !ok {
		c.Data["json"] = map[string]interface{}{"error": "Parcelas no especificadas correctamente"}
		c.ServeJSON()
		return
	}
	delete(body, "Parcelas") // eliminamos del body principal

	// Enviar arrendamiento a API CRUD
	jsonData, _ := json.Marshal(body)
	fmt.Println("Enviando arrendamiento:", string(jsonData))
	response, err := services.Metodo_post("API_CRUD_FINCA", "/v1/Arrendamiento", jsonData)
	if err != nil {
		c.respondWithError("Error al crear arrendamiento", err.Error())
		return
	}

	// Parsear respuesta para obtener ID del arrendamiento creado
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		c.respondWithError("Error al leer respuesta del arrendamiento", err.Error())
		return
	}
	idArrendamiento := int(result["Id"].(float64)) // Asumiendo que la respuesta devuelve {"Id": 123}

	// Crear registros en ArrendamientoParcela
	for _, p := range parcelas {
		parcelaId := int(p.(float64))
		dataParcela := map[string]interface{}{
			"FkArrendamiento": map[string]interface{}{
				"Id": idArrendamiento,
			},
			"FkParcela": map[string]interface{}{
				"Id": parcelaId,
			},
		}
		parcelaJSON, _ := json.Marshal(dataParcela)
		fmt.Println("Enviando arrendamiento_parcela:", string(parcelaJSON))
		_, err := services.Metodo_post("API_CRUD_FINCA", "/v1/Arrendamiento_Parcela", parcelaJSON)
		if err != nil {
			c.respondWithError("Error al crear arrendamiento_parcela", err.Error())
			return
		}
	}

	c.Data["json"] = map[string]interface{}{
		"mensaje":          "Arrendamiento creado con éxito",
		"id_arrendamiento": idArrendamiento,
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
// @router /activos/:id/ [get]
func (c *Gestion_arrendamientoController) GetActivosPorFinca() {
	fmt.Println("Obteniendo arrendamientos activos con parcelas")

	fincaID := c.Ctx.Input.Param(":id")
	if fincaID == "" {
		c.respondWithError("El ID de la finca es requerido")
		return
	}

	// Obtener los arrendamientos
	response, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento", "")
	if err != nil {
		c.respondWithError("Error al obtener arrendamientos", err.Error())
		return
	}

	// Deserializar la respuesta
	var arrendamientos []map[string]interface{}
	if err := json.Unmarshal(response, &arrendamientos); err != nil {
		c.respondWithError("Error al procesar la respuesta del API CRUD", err.Error())
		return
	}

	// Filtrar arrendamientos activos
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
			// Si el arrendamiento es activo, buscamos las parcelas asociadas
			idArrendamiento := int(arr["Id"].(float64)) // Asumimos que "Id" es un número

			// Obtener las parcelas asociadas a este arrendamiento
			url := fmt.Sprintf("?query=FkArrendamiento.Id:%d", idArrendamiento)
			respParcela, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento_Parcela", url)
			if err != nil {
				c.respondWithError("Error al obtener parcelas del arrendamiento", err.Error())
				return
			}

			var arrendamientoParcelas []map[string]interface{}
			if err := json.Unmarshal(respParcela, &arrendamientoParcelas); err != nil {
				c.respondWithError("Error al procesar parcelas del arrendamiento", err.Error())
				return
			}

			parcelas := []string{}
			for _, ap := range arrendamientoParcelas {
				if parcelaData, ok := ap["FkParcela"].(map[string]interface{}); ok {
					if nombreParcela, ok := parcelaData["NombreParcela"].(string); ok {
						parcelas = append(parcelas, nombreParcela)
					}
				}
			}

			activos = append(activos, map[string]interface{}{
				"Id":          idArrendamiento,
				"Parcelas":    parcelas,
				"FechaInicio": arr["FechaInicio"],
				"FechaFin":    arr["FechaFin"],
				"Activo":      arr["Activo"],
			})
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
		idArr := int(arr["Id"].(float64)) // convertir a int desde interface{}

		// Obtener las parcelas asociadas al arrendamiento
		url := fmt.Sprintf("?query=FkArrendamiento.Id:%d", idArr)
		respParcela, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento_Parcela", url)
		if err != nil {
			c.respondWithError("Error al obtener parcelas del arrendamiento", err.Error())
			return
		}

		var arrendamientoParcelas []map[string]interface{}
		if err := json.Unmarshal(respParcela, &arrendamientoParcelas); err != nil {
			c.respondWithError("Error al procesar parcelas del arrendamiento", err.Error())
			return
		}

		parcelas := []string{}
		areas := []interface{}{}

		for _, ap := range arrendamientoParcelas {
			if parcelaData, ok := ap["FkParcela"].(map[string]interface{}); ok {
				parcelas = append(parcelas, parcelaData["NombreParcela"].(string))
				areas = append(areas, parcelaData["TamanoParcela"])
			}
		}

		// Información adicional del arrendamiento
		var finca string
		var fincaId float64
		if fincaMap, ok := arr["FkArrendamientoFinca"].(map[string]interface{}); ok {
			if fincaNombre, ok := fincaMap["Nombre"].(string); ok {
				finca = fincaNombre
			}
			if fincaID, ok := fincaMap["Id"].(float64); ok {
				fincaId = fincaID

			}

		}

		var arrendatario, contacto string
		if arrendatarioMap, ok := arr["IdUserUserArrendatario"].(map[string]interface{}); ok {
			if nombre, ok := arrendatarioMap["Nombre"].(string); ok {
				arrendatario = nombre
			}
			if cont, ok := arrendatarioMap["Contacto"].(string); ok {
				contacto = cont
			}
		}

		resultado = append(resultado, map[string]interface{}{
			"Finca":            finca,
			"FincaID":          fincaId,
			"Parcelas":         parcelas,
			"Areas":            areas,
			"Arrendatario":     arrendatario,
			"FechaInicio":      arr["FechaInicio"],
			"FechaFin":         arr["FechaFin"],
			"Contacto":         contacto,
			"Valor":            arr["Valor"],
			"Id":               idArr,
			"arrendamiento":    arr["Activo"],
			"CantidadParcelas": len(parcelas),
		})
	}

	// Responder con los datos formateados
	c.Data["json"] = resultado
	c.ServeJSON()

}

// GetParcelasPorArrendamiento obtiene las parcelas asociadas a un arrendamiento.
// @Title GetParcelasPorArrendamiento
// @Description Obtener parcelas asociadas a un arrendamiento específico.
// @Param	id		path 	string	true		"ID del arrendamiento"
// @Success 200 {object} []map[string]interface{}
// @Failure 403 :id is empty
// @router /parcelasarrendamiento/:id [get]
func (c *Gestion_arrendamientoController) GetParcelasPorArrendamiento() {
	arrendamientoID := c.Ctx.Input.Param(":id")
	if arrendamientoID == "" {
		c.Data["json"] = map[string]interface{}{"error": "El ID de arrendamiento es obligatorio"}
		c.ServeJSON()
		return
	}

	// Paso 1: Obtener las parcelas asociadas al arrendamiento
	queryArrendamiento := "?query=FkArrendamiento.Id:" + arrendamientoID
	resArrendamiento, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento_Parcela", queryArrendamiento)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al obtener arrendamientos", "details": err.Error()}
		c.ServeJSON()
		return
	}

	var resultArr []map[string]interface{}
	if err := json.Unmarshal(resArrendamiento, &resultArr); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al parsear arrendamientos", "details": err.Error()}
		c.ServeJSON()
		return
	}

	var arrendamientoData []map[string]interface{}
	if err := json.Unmarshal(resArrendamiento, &arrendamientoData); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al parsear arrendamientos", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Paso 2: Obtener los ID de las parcelas
	var parcelaIDs []string
	for _, item := range arrendamientoData {
		if parcela, ok := item["FkParcela"].(map[string]interface{}); ok {
			if id, ok := parcela["Id"].(float64); ok {
				parcelaIDs = append(parcelaIDs, fmt.Sprintf("%.0f", id))
			}
		}
	}

	if len(parcelaIDs) == 0 {
		c.Data["json"] = map[string]interface{}{"mensaje": "No se encontraron parcelas asociadas"}
		c.ServeJSON()
		return
	}

	// Paso 3: Construir la query para FincaParcela con todos los IDs

	// Paso 4: Armar la respuesta
	var parcelas []map[string]interface{}

	for _, parcelaID := range parcelaIDs {
		query := "?query=FkParcelaFinca.Id:" + parcelaID
		res, err := services.Metodo_get("API_CRUD_FINCA", "/v1/FincaParcela", query)
		if err != nil {
			continue // puedes registrar el error si quieres depurar
		}

		var result map[string]interface{}
		if err := json.Unmarshal(res, &result); err != nil {
			continue
		}

		parcelasData, ok := result["Data"].([]interface{})
		if !ok {
			continue
		}

		for _, item := range parcelasData {
			parcelaMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			parcelaDatos, ok := parcelaMap["FkParcelaFinca"].(map[string]interface{})
			if !ok {
				continue
			}

			geoloc, _ := parcelaMap["FkGeolocalizacion"].(map[string]interface{})

			parcelaInfo := map[string]interface{}{
				"IdParcela":     parcelaDatos["Id"],
				"NombreParcela": parcelaDatos["NombreParcela"],
				"TamanoParcela": parcelaDatos["TamanoParcela"],
				"Geolocalizacion": map[string]interface{}{
					"IdGeolocalizacion": geoloc["Id"],
					"LatitudInicial":    geoloc["LatitudInicial"],
					"LongitudInicial":   geoloc["LongitudInicial"],
					"LatitudFinal":      geoloc["LatitudFinal"],
					"LongitudFinal":     geoloc["LongitudFinal"],
				},
			}
			parcelas = append(parcelas, parcelaInfo)
		}
	}

	c.Data["json"] = parcelas
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
