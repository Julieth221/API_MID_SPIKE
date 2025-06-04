package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

	// Extraer y eliminar Parcelas del body principal
	parcelas, ok := body["Parcelas"].([]interface{})
	if !ok || len(parcelas) == 0 {
		c.Data["json"] = map[string]interface{}{"error": "Parcelas no especificadas correctamente"}
		c.ServeJSON()
		return
	}
	delete(body, "Parcelas")

	// BUSCAR ARRENDAMIENTO ANTERIOR (opcional)
	// parcelaPrincipal := parcelas[0].(map[string]interface{})
	// idParcela := int(parcelaPrincipal["IdParcela"].(float64))

	// query := fmt.Sprintf("FkParcela.Id:%d", idParcela)
	// parametros := fmt.Sprintf("?query=%s&limit=1&sortby=FechaFin&order=desc", url.QueryEscape(query))

	// resp, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento_Parcela", parametros)
	// if err == nil {
	// 	var historial []map[string]interface{}
	// 	if err := json.Unmarshal(resp, &historial); err == nil && len(historial) > 0 {
	// 		if arrendamiento, ok := historial[0]["FkArrendamiento"].(map[string]interface{}); ok {
	// 			idAnterior := int(arrendamiento["Id"].(float64))
	// 			body["FkArrendamientoAnterior"] = map[string]interface{}{"Id": idAnterior}
	// 			fmt.Println("Arrendamiento anterior detectado: ", idAnterior)
	// 		}
	// 	}
	// }

	// Enviar arrendamiento (solo los datos del modelo Arrendamiento)

	// Formatear fechas a formato RFC3339 (necesario para el modelo Arrendamiento)
	if fechaInicio, ok := body["FechaInicio"].(string); ok {
		if parsed, err := time.Parse("2006-01-02", fechaInicio); err == nil {
			body["FechaInicio"] = parsed.Format(time.RFC3339)
		}
	}
	if fechaFin, ok := body["FechaFin"].(string); ok {
		if parsed, err := time.Parse("2006-01-02", fechaFin); err == nil {
			body["FechaFin"] = parsed.Format(time.RFC3339)
		}
	}

	jsonData, _ := json.Marshal(body)
	fmt.Println("Enviando arrendamiento:", string(jsonData))
	response, err := services.Metodo_post("API_CRUD_FINCA", "/v1/Arrendamiento", jsonData)
	if err != nil {
		c.respondWithError("Error al crear arrendamiento", err.Error())
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		c.respondWithError("Error al leer respuesta del arrendamiento", err.Error())
		return
	}
	idArrendamiento := int(result["Id"].(float64))

	// Crear registros en ArrendamientoParcela
	for _, p := range parcelas {
		parcelaData, ok := p.(map[string]interface{})
		if !ok {
			continue
		}

		// Parsear fechas individualmente por parcela
		// if fechaInicio, ok := parcelaData["FechaInicio"].(string); ok {
		// 	parsed, err := time.Parse("2006-01-02", fechaInicio)
		// 	if err == nil {
		// 		parcelaData["FechaInicio"] = parsed.Format(time.RFC3339Nano)
		// 	}
		// }
		// if fechaFin, ok := parcelaData["FechaFin"].(string); ok {
		// 	parsed, err := time.Parse("2006-01-02", fechaFin)
		// 	if err == nil {
		// 		parcelaData["FechaFin"] = parsed.Format(time.RFC3339Nano)
		// 	}
		// }

		// Construir objeto para ArrendamientoParcela
		dataParcela := map[string]interface{}{
			"FkArrendamiento": map[string]interface{}{
				"Id": idArrendamiento,
			},
			"FkParcela": map[string]interface{}{
				"Id": int(parcelaData["IdParcela"].(float64)),
			},
			// "FechaInicio": parcelaData["FechaInicio"],
			// "FechaFin":    parcelaData["FechaFin"],
			"Valor":  parcelaData["Valor"],
			"Activo": true,
		}
		fmt.Println("Este es el body que se envia a Arrendamiento_Parcela: ", dataParcela)

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

// Crear un nuevo arrendamiento de un arrendamiento padre
// @Title PostNuevoArrendamiento
// @Description Desactiva arrendamiento anterior y crea uno nuevo con los cambios
// @Param	body		body 	map[string]interface{}	true	"body con datos del nuevo arrendamiento"
// @Success 200 {object} map[string]interface{}
// @Failure 403 body is empty
// @router /arrendamiento/versionar [post]
func (c *Gestion_arrendamientoController) PostNuevoArrendamiento() {
	fmt.Println("versionar arrendamiento")

	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la solicitud"}
		c.ServeJSON()
		return
	}

	// 1. Validar ID del arrendamiento anterior
	idAnterior, ok := body["IdArrendamientoAnterior"].(float64)
	if !ok {
		c.Data["json"] = map[string]interface{}{"error": "IdArrendamientoAnterior es requerido y debe ser numérico"}
		c.ServeJSON()
		return
	}

	// 2. Desactivar arrendamiento anterior
	update := map[string]interface{}{
		"Activo": false,
	}
	jsonUpdate, _ := json.Marshal(update)
	_, err := services.Metodo_patch("API_CRUD_FINCA", "/v1/Arrendamiento", fmt.Sprintf("%d", int(idAnterior)), jsonUpdate)
	if err != nil {
		c.respondWithError("Error al desactivar arrendamiento anterior", err.Error())
		return
	}

	// 3. Validar fechas en el nuevo arrendamiento
	parseFecha := func(fechaStr string) (time.Time, error) {
		return time.Parse("2006-01-02", fechaStr)
	}

	fechaInicioStr, okInicio := body["FechaInicio"].(string)
	fechaFinStr, okFin := body["FechaFin"].(string)

	if !okInicio || !okFin {
		c.Data["json"] = map[string]interface{}{"error": "FechaInicio y FechaFin deben estar presentes como strings"}
		c.ServeJSON()
		return
	}

	fechaInicio, errInicio := parseFecha(fechaInicioStr)
	fechaFin, errFin := parseFecha(fechaFinStr)

	if errInicio != nil || errFin != nil {
		c.Data["json"] = map[string]interface{}{"error": "Fechas en formato incorrecto (esperado: YYYY-MM-DD)"}
		c.ServeJSON()
		return
	}

	// Reemplazamos los strings por objetos time
	body["FechaInicio"] = fechaInicio.Format(time.RFC3339Nano)
	body["FechaFin"] = fechaFin.Format(time.RFC3339Nano)

	// 4. Preparar nuevo arrendamiento
	body["FkArrendamientoAnterior"] = map[string]interface{}{"Id": int(idAnterior)}
	delete(body, "IdArrendamientoAnterior")

	// Extraer parcelas
	parcelas, ok := body["Parcelas"].([]interface{})
	if !ok {
		c.Data["json"] = map[string]interface{}{"error": "Parcelas no especificadas correctamente"}
		c.ServeJSON()
		return
	}
	delete(body, "Parcelas") // Eliminamos para enviar solo lo del arrendamiento

	// 5. Crear nuevo arrendamiento
	jsonData, _ := json.Marshal(body)
	res, err := services.Metodo_post("API_CRUD_FINCA", "/v1/Arrendamiento", jsonData)
	if err != nil {
		c.respondWithError("Error al crear nuevo arrendamiento", err.Error())
		return
	}

	// Obtener ID del nuevo arrendamiento
	var result map[string]interface{}
	if err := json.Unmarshal(res, &result); err != nil {
		c.respondWithError("Error al procesar la respuesta", err.Error())
		return
	}
	nuevoID := int(result["Id"].(float64))

	// 6. Asociar parcelas al nuevo arrendamiento
	for _, p := range parcelas {
		parcelaMap, ok := p.(map[string]interface{})
		if !ok {
			c.respondWithError("Error en el formato de una parcela", "Se esperaba un objeto")
			return
		}

		idParcela, ok := parcelaMap["IdParcela"].(float64)
		if !ok {
			c.respondWithError("Parcela inválida", "IdParcela es requerido y debe ser numérico")
			return
		}

		valor, ok := parcelaMap["Valor"].(string)
		if !ok {
			c.respondWithError("Valor inválido", "Valor debe ser string")
			return
		}

		parcelaData := map[string]interface{}{
			"FkArrendamiento": map[string]interface{}{"Id": nuevoID},
			"FkParcela":       map[string]interface{}{"Id": int(idParcela)},
			"Valor":           valor,
			"Activo":          true,
		}

		parcelaJSON, _ := json.Marshal(parcelaData)
		fmt.Println("Este es el body que se envia a la tabla Arrendamiento_Parcela: ", string(parcelaJSON))
		_, err := services.Metodo_post("API_CRUD_FINCA", "/v1/Arrendamiento_Parcela", parcelaJSON)
		if err != nil {
			c.respondWithError("Error al crear arrendamiento_parcela", err.Error())
			return
		}
	}

	// 7. Respuesta exitosa
	c.Data["json"] = map[string]interface{}{
		"mensaje":             "Arrendamiento actualizado correctamente",
		"nuevo_arrendamiento": nuevoID,
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

// @Title GetDisponiblesPorFinca
// @Description Retorna las parcelas disponibles (no arrendadas activamente) de una finca
// @Param	id		path 	string	true		"ID de la finca"
// @Success 200 {object} []map[string]interface{}
// @Failure 403 :id is empty
// @router /disponibles/:id/ [get]
func (c *Gestion_arrendamientoController) GetDisponiblesPorFinca() {
	fmt.Println("arrendamiento parcelas disponible por finca")

	fincaID := c.Ctx.Input.Param(":id")
	if fincaID == "" {
		c.respondWithError("El ID de la finca es requerido")
		return
	}

	// 1. Obtener todas las parcelas de la finca
	urlParcelas := fmt.Sprintf("?query=FkFincaParcela.Id:%s", fincaID)
	respParcelas, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Parcela", urlParcelas)
	if err != nil {
		c.respondWithError("Error al obtener parcelas de la finca", err.Error())
		return
	}

	fmt.Println("Este es el response de la API:", string(respParcelas))

	// Deserializar correctamente el cuerpo de la respuesta
	var resParcelas map[string]interface{}
	if err := json.Unmarshal(respParcelas, &resParcelas); err != nil {
		c.respondWithError("Error al procesar parcelas", err.Error())
		return
	}

	// Extraer el arreglo de parcelas desde la clave "Data"
	dataRaw, ok := resParcelas["Data"].([]interface{})
	if !ok {
		c.respondWithError("Estructura inesperada en la respuesta de parcelas")
		return
	}

	var parcelas []map[string]interface{}
	for _, item := range dataRaw {
		if parcela, ok := item.(map[string]interface{}); ok {
			parcelas = append(parcelas, parcela)
		}
	}

	// 2. Obtener todos los arrendamientos por finca
	urlArrendamientos := fmt.Sprintf("?query=FkArrendamientoFinca.Id:%s", fincaID)
	respArrendamientos, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento", urlArrendamientos)
	if err != nil {
		c.respondWithError("Error al obtener arrendamientos", err.Error())
		return
	}
	fmt.Println("Esta es la respuesta de la api: ", string(respArrendamientos))

	var arrendamientos []map[string]interface{}
	if err := json.Unmarshal(respArrendamientos, &arrendamientos); err != nil {
		c.respondWithError("Error al procesar arrendamientos", err.Error())
		return
	}

	// dataArrRaw, ok := resArrendamientos["Data"].([]interface{})
	// if !ok {
	// 	c.respondWithError("Estructura inesperada en arrendamientos")
	// 	return
	// }

	var arrendamientoIDs []int
	for _, arrendamiento := range arrendamientos {
		if idFloat, ok := arrendamiento["Id"].(float64); ok {
			arrendamientoIDs = append(arrendamientoIDs, int(idFloat))
		}
	}

	// 3. Obtener todos los arrendamiento_parcela usando los IDs recolectados
	var arrParcela []map[string]interface{}
	fmt.Println("ArrendamientoIDs obtenidos:", arrendamientoIDs)

	for _, arrID := range arrendamientoIDs {
		url := fmt.Sprintf("?query=FkArrendamiento.Id:%d", arrID)
		fmt.Println("Consultando arrendamiento_parcela con URL:", url)

		respAP, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento_Parcela", url)
		if err != nil {
			fmt.Println("Error al obtener Arrendamiento_Parcela para arrID:", arrID, "Error:", err)
			continue // seguir aunque una falle
		}
		fmt.Println("Response de API Arrendamiento_Parcela:", string(respAP))

		if len(respAP) == 0 {
			fmt.Println("Respuesta vacía para arrID:", arrID)
			continue
		}

		if respAP[0] == '[' {
			var parcial []map[string]interface{}
			if err := json.Unmarshal(respAP, &parcial); err == nil {
				fmt.Println("Parsed como arreglo. Cantidad de elementos:", len(parcial))
				arrParcela = append(arrParcela, parcial...)
			} else {
				fmt.Println("Error al hacer Unmarshal del arreglo:", err)
			}
		} else {
			var temp map[string]interface{}
			if err := json.Unmarshal(respAP, &temp); err == nil {
				fmt.Println("Parsed como objeto único. ID:", temp["Id"])
				arrParcela = append(arrParcela, temp)
			} else {
				fmt.Println("Error al hacer Unmarshal del objeto:", err)
			}
		}
	}

	fmt.Println("Total de elementos en arrParcela:", len(arrParcela))

	arrendamientoMap := make(map[int]map[string]interface{})
	for _, a := range arrendamientos {
		if idFloat, ok := a["Id"].(float64); ok {
			arrendamientoMap[int(idFloat)] = a
		}
	}

	// 4. Identificar parcelas arrendadas actualmente
	parcelasArrendadas := make(map[int]bool)
	ahora := time.Now()

	for _, ap := range arrParcela {
		fmt.Println("Procesando arrParcela ID:", ap["Id"])

		fmt.Println("Procesando arrParcela ID:", ap["Id"])

		activo, _ := ap["Activo"].(bool)
		fmt.Println("Activo:", activo)

		// Obtener ID del arrendamiento
		var idArr int
		if fkArr, ok := ap["FkArrendamiento"].(map[string]interface{}); ok {
			if idFloat, ok := fkArr["Id"].(float64); ok {
				idArr = int(idFloat)
			}
		}

		// Buscar arrendamiento y obtener fechas
		arr, existe := arrendamientoMap[idArr]
		if !existe {
			fmt.Println("Arrendamiento no encontrado para ID:", idArr)
			continue
		}

		fechaInicioStr := fmt.Sprintf("%v", arr["FechaInicio"])
		fechaFinStr := fmt.Sprintf("%v", arr["FechaFin"])

		fechaInicio, err1 := time.Parse(time.RFC3339, fechaInicioStr)
		fechaFin, err2 := time.Parse(time.RFC3339, fechaFinStr)

		fmt.Println("FechaInicio (from Arrendamiento):", fechaInicioStr, "err1:", err1)
		fmt.Println("FechaFin (from Arrendamiento):", fechaFinStr, "err2:", err2)

		if activo && err1 == nil && err2 == nil && ahora.After(fechaInicio) && ahora.Before(fechaFin) {
			if parcelaData, ok := ap["FkParcela"].(map[string]interface{}); ok {
				if idFloat, ok := parcelaData["Id"].(float64); ok {
					parcelasArrendadas[int(idFloat)] = true
					fmt.Println("Parcela arrendada actualmente con ID:", int(idFloat))
				}
			}
		} else {
			fmt.Println("Arrendamiento inactivo o fuera de fechas.")
		}
	}

	// 5. Filtrar TODAS las parcelas de la finca y eliminar las que están arrendadas actualmente
	var disponibles []map[string]interface{}

	for _, p := range parcelas {
		if idFloat, ok := p["Id"].(float64); ok {
			id := int(idFloat)

			if !parcelasArrendadas[id] {
				disponibles = append(disponibles, map[string]interface{}{
					"Id":            id,
					"NombreParcela": p["NombreParcela"],
					"Estado":        "Disponible",
				})
			}
		}
	}
	c.Data["json"] = disponibles
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
// @router /arrendamiento/parcelas/porfinca/:id [get]
func (c *Gestion_arrendamientoController) GetAll() {
	fmt.Println("Parcelas arrendadas por finca (vista agrupada)")
	fincaID := c.Ctx.Input.Param(":id")
	if fincaID == "" {
		c.respondWithError("El ID de la finca es requerido")
		return
	}

	// Obtener arrendamientos activos de la finca
	arrendamientos, err := obtenerArrendamientosPorFinca(fincaID)
	if err != nil {
		c.respondWithError("Error al obtener arrendamientos", err.Error())
		return
	}

	// Obtener todas las relaciones arrendamiento-parcela
	arrParcela, err := obtenerArrendamientoParcelas(arrendamientos)
	if err != nil {
		c.respondWithError("Error al obtener arrendamiento_parcela", err.Error())
		return
	}

	// Obtener mapa de parcelas disponibles
	parcelasDisponibles, err := obtenerParcelasDisponiblesMap(fincaID)
	if err != nil {
		c.respondWithError("Error al verificar parcelas disponibles", err.Error())
		return
	}

	// Agrupar por arrendamiento
	agrupado := make(map[int]map[string]interface{})

	for _, ap := range arrParcela {
		parcela, ok := ap["FkParcela"].(map[string]interface{})
		if !ok {
			continue
		}
		arrendamiento, ok := ap["FkArrendamiento"].(map[string]interface{})
		if !ok {
			continue
		}

		idArr := int(arrendamiento["Id"].(float64))
		nombreParcela := fmt.Sprintf("%v", parcela["NombreParcela"])
		valorParcela, _ := strconv.ParseFloat(fmt.Sprintf("%v", ap["Valor"]), 64)

		estado := "Arrendada"
		if parcelasDisponibles[int(parcela["Id"].(float64))] {
			estado = "Inactivo"
		}

		if _, existe := agrupado[idArr]; !existe {
			// Nuevo arrendamiento: inicializar
			arrendatario := ""
			IdArrendatario := ""
			if user, ok := arrendamiento["IdUserUserArrendatario"].(map[string]interface{}); ok {
				arrendatario = fmt.Sprintf("%v", user["Nombre"])
				IdArrendatario = fmt.Sprintf("%v", user["Id"])
			}
			agrupado[idArr] = map[string]interface{}{
				"IdArrendamiento": idArr,
				"Arrendatario":    arrendatario,
				"IdArrendatario":  IdArrendatario,
				"FechaInicio":     formatearFecha(arrendamiento["FechaInicio"]),
				"FechaFin":        formatearFecha(arrendamiento["FechaFin"]),
				"Parcelas":        []string{nombreParcela},
				"ValorTotal":      valorParcela,
				"Estado":          estado,
			}
		} else {
			// Ya existe:
			entry := agrupado[idArr]
			entry["Parcelas"] = append(entry["Parcelas"].([]string), nombreParcela)
			entry["ValorTotal"] = entry["ValorTotal"].(float64) + valorParcela
		}
	}

	// Convertir a slice
	var resultado []map[string]interface{}
	for _, entry := range agrupado {
		// Opcional: concatenar parcelas como string
		entry["Parcelas"] = strings.Join(entry["Parcelas"].([]string), ", ")
		resultado = append(resultado, entry)
	}

	c.Data["json"] = resultado
	c.ServeJSON()
}

func formatearFecha(fecha interface{}) string {
	if str, ok := fecha.(string); ok {
		if t, err := time.Parse(time.RFC3339, str); err == nil {
			// Forzar a UTC sin conversión local
			return t.In(time.UTC).Format("2006-01-02")
		}
	}
	return ""
}

func obtenerArrendamientosPorFinca(fincaID string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("?query=FkArrendamientoFinca.Id:%s,Activo:true", fincaID)
	resp, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento", url)
	if err != nil {
		return nil, err
	}

	var arrendamientos []map[string]interface{}
	if err := json.Unmarshal(resp, &arrendamientos); err != nil {
		return nil, err
	}
	return arrendamientos, nil
}

func obtenerArrendamientoParcelas(arrendamientos []map[string]interface{}) ([]map[string]interface{}, error) {
	var resultado []map[string]interface{}
	for _, a := range arrendamientos {
		id, ok := a["Id"].(float64)
		if !ok {
			continue
		}
		url := fmt.Sprintf("?query=FkArrendamiento.Id:%d", int(id))
		resp, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento_Parcela", url)
		if err != nil || len(resp) == 0 {
			continue
		}

		if resp[0] == '[' {
			var arr []map[string]interface{}
			if err := json.Unmarshal(resp, &arr); err == nil {
				resultado = append(resultado, arr...)
			}
		} else {
			var obj map[string]interface{}
			if err := json.Unmarshal(resp, &obj); err == nil {
				resultado = append(resultado, obj)
			}
		}
	}
	return resultado, nil
}

func obtenerParcelasDisponiblesMap(fincaID string) (map[int]bool, error) {
	urlParcelas := fmt.Sprintf("?query=FkFincaParcela.Id:%s", fincaID)
	resp, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Parcela", urlParcelas)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, err
	}

	rawData, ok := res["Data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("estructura inesperada en parcelas")
	}

	parcelasDisponibles := make(map[int]bool)
	for _, item := range rawData {
		if p, ok := item.(map[string]interface{}); ok {
			id := int(p["Id"].(float64))
			parcelasDisponibles[id] = true
		}
	}

	// Obtener arrendamientos y arr_parcela para verificar activos
	arrs, _ := obtenerArrendamientosPorFinca(fincaID)
	apList, _ := obtenerArrendamientoParcelas(arrs)

	ahora := time.Now()
	for _, ap := range apList {
		activo, _ := ap["Activo"].(bool)

		arrendamiento, ok := ap["FkArrendamiento"].(map[string]interface{})
		if !ok {
			continue
		}

		fi, err1 := time.Parse(time.RFC3339, fmt.Sprintf("%v", arrendamiento["FechaInicio"]))
		ff, err2 := time.Parse(time.RFC3339, fmt.Sprintf("%v", arrendamiento["FechaFin"]))
		if err1 != nil || err2 != nil {
			continue
		}

		if activo && ahora.After(fi) && ahora.Before(ff) {
			if parcela, ok := ap["FkParcela"].(map[string]interface{}); ok {
				id := int(parcela["Id"].(float64))
				delete(parcelasDisponibles, id)
			}
		}
	}

	return parcelasDisponibles, nil

}

// GetParcelasPorArrendamiento obtiene las parcelas asociadas a un arrendamiento.
// @Title GetParcelaPorArrendamiento
// @Description Obtener parcelas asociadas a un arrendamiento específico.
// @Param	id		path 	string	true		"ID del arrendamiento"
// @Success 200 {object} []map[string]interface{}
// @Failure 403 :id is empty
// @router /parcelasarrendamiento/:id [get]
func (c *Gestion_arrendamientoController) GetParcelaPorArrendamiento() {
	arrendamientoID := c.Ctx.Input.Param(":id")
	if arrendamientoID == "" {
		c.Data["json"] = map[string]interface{}{"error": "El ID del arrendamiento es obligatorio"}
		c.ServeJSON()
		return
	}

	// 1. Obtener arrendamiento por ID
	resArr, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento/", arrendamientoID)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al obtener arrendamiento", "details": err.Error()}
		c.ServeJSON()
		return
	}

	var arrendamiento map[string]interface{}
	if err := json.Unmarshal(resArr, &arrendamiento); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al parsear arrendamiento", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// 2. Obtener datos del arrendatario
	userArr, ok := arrendamiento["IdUserUserArrendatario"].(map[string]interface{})
	if !ok || userArr["Id"] == nil {
		c.Data["json"] = map[string]interface{}{"error": "Datos de arrendatario no encontrados"}
		c.ServeJSON()
		return
	}

	// 3. Obtener parcelas del arrendamiento
	query := "?query=FkArrendamiento.Id:" + arrendamientoID
	respArrendParcela, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento_Parcela/", query)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al obtener arrendamiento_parcela", "details": err.Error()}
		c.ServeJSON()
		return
	}

	var arrendParcelas []map[string]interface{}
	if err := json.Unmarshal(respArrendParcela, &arrendParcelas); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al parsear arrendamiento_parcela", "details": err.Error()}
		c.ServeJSON()
		return
	}

	if len(arrendParcelas) == 0 {
		c.Data["json"] = map[string]interface{}{"error": "No se encontraron parcelas asociadas al arrendamiento"}
		c.ServeJSON()
		return
	}

	var parcelasResult []map[string]interface{}
	var valorTotal float64 = 0

	for _, ap := range arrendParcelas {
		fkParcela, ok := ap["FkParcela"].(map[string]interface{})
		if !ok || fkParcela["Id"] == nil {
			continue
		}

		parcelaID := fmt.Sprintf("%.0f", fkParcela["Id"].(float64))
		query := "?query=FkParcelaFinca.Id:" + parcelaID
		res, err := services.Metodo_get("API_CRUD_FINCA", "/v1/FincaParcela", query)
		if err != nil {
			c.Data["json"] = map[string]interface{}{"error": "Error al obtener la parcela", "details": err.Error()}
			c.ServeJSON()
			return
		}

		var result map[string]interface{}
		if err := json.Unmarshal(res, &result); err != nil {
			c.Data["json"] = map[string]interface{}{"error": "Error al parsear respuesta", "details": err.Error()}
			c.ServeJSON()
			return
		}

		parcelasData, ok := result["Data"].([]interface{})
		if !ok || len(parcelasData) == 0 {
			c.Data["json"] = map[string]interface{}{"error": "Parcela no encontrada"}
			c.ServeJSON()
			return
		}

		item, ok := parcelasData[0].(map[string]interface{})
		if !ok {
			c.Data["json"] = map[string]interface{}{"error": "Estructura de parcela inválida"}
			c.ServeJSON()
			return
		}

		parcelaDatos, ok := item["FkParcelaFinca"].(map[string]interface{})
		if !ok {
			c.Data["json"] = map[string]interface{}{"error": "No se encontraron datos de la parcela"}
			c.ServeJSON()
			return
		}

		geoloc, _ := item["FkGeolocalizacion"].(map[string]interface{})

		valorStr, _ := ap["Valor"].(string)
		valorFloat, _ := strconv.ParseFloat(valorStr, 64)
		valorTotal += valorFloat

		parcelaInfo := map[string]interface{}{
			"IdParcela":     parcelaDatos["Id"],
			"NombreParcela": parcelaDatos["NombreParcela"],
			"TamanoParcela": parcelaDatos["TamanoParcela"],
			"Valor":         ap["Valor"],
			"Geolocalizacion": map[string]interface{}{
				"IdGeolocalizacion": geoloc["Id"],
				"LatitudInicial":    geoloc["LatitudInicial"],
				"LongitudInicial":   geoloc["LongitudInicial"],
				"LatitudFinal":      geoloc["LatitudFinal"],
				"LongitudFinal":     geoloc["LongitudFinal"],
			},
		}
		parcelasResult = append(parcelasResult, parcelaInfo)
	}

	// Armar respuesta completa
	c.Data["json"] = map[string]interface{}{
		"IdArrendamiento":        arrendamiento["Id"],
		"NombreArrendatario":     userArr["Nombre"],
		"ContactoArrendatario":   userArr["Contacto"],
		"IdUserUserArrendatario": map[string]interface{}{"Id": userArr["Id"]},
		"FechaInicio":            formatearFecha(arrendamiento["FechaInicio"]),
		"FechaFin":               formatearFecha(arrendamiento["FechaFin"]),
		"ValorTotal":             valorTotal,
		"Parcelas":               parcelasResult,
	}
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
