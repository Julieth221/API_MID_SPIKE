package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	auth_JWT "github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/auth_jwt"
	"github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/services"
)

// Gestion_fincaController operations for Gestion_finca
type Gestion_fincaController struct {
	beego.Controller
}

// URLMapping ...
func (c *Gestion_fincaController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetFinca)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Registrar Finca
// @Title Create
// @Description create Gestion_finca
// @Param	body		body 	models.Gestion_finca	true		"body for Gestion_finca content"
// @Success 201 {object} models.Gestion_finca
// @Failure 403 body is empty
// @router / [post]
func (c *Gestion_fincaController) Post() {
	fmt.Println("Registrar finca")
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

	// Obtener ID del Tipo de Suelo
	tipoSueloID, err := getTipoSueloID(body["tipo_suelo"].(string))
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al obtener el tipo de suelo"}
		c.ServeJSON()
		return
	}
	fmt.Println("Id de Tipo de suelo: ", tipoSueloID)

	// Crear Finca
	fincaID, err := crearFinca(body, userID, tipoSueloID)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al crear la finca"}
		c.ServeJSON()
		return
	}

	// Crear Parcelas y Geolocalización
	err = crearParcelas(body["parcelas"].([]interface{}), fincaID)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al crear parcelas"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"mensaje": "Finca creada con éxito", "finca_id": fincaID}
	c.ServeJSON()
}

// getTipoSueloID consulta el ID del tipo de suelo
func getTipoSueloID(nombreTipo string) (int, error) {
	response, err := services.Metodo_get("API_CRUD_FINCA", "/v1/tipo_suelo?query=nombre:", nombreTipo)
	if err != nil {
		return 0, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response, &data); err != nil {
		return 0, err
	}

	if resultados, ok := data["Data"].([]interface{}); ok && len(resultados) > 0 {
		tipoSuelo := resultados[0].(map[string]interface{})
		return int(tipoSuelo["Id"].(float64)), nil
	}
	return 0, fmt.Errorf("tipo de suelo no encontrado")
}

// crearFinca registra una finca en el API CRUD
func crearFinca(body map[string]interface{}, userID int, tipoSueloID int) (int, error) {
	jsonFinca := map[string]interface{}{
		"Nombre":         body["Nombre"],
		"AreaTotal":      body["AreaTotal"],
		"TotalParcelas":  body["TotalParcelas"],
		"TamañoParcelas": body["TamañoParcelas"],
		"FkFinca":        map[string]interface{}{"Id": tipoSueloID},
		"id_usuario":     userID,
	}

	jsonFincaByte, _ := json.Marshal(jsonFinca)
	response, err := services.Metodo_post("API_CRUD_FINCA", "/v1/Finca", jsonFincaByte)
	if err != nil {
		return 0, err
	}
	fmt.Println("Este es el reponse de crear finca: ", string(response))

	var fincaResponse map[string]interface{}
	if err := json.Unmarshal(response, &fincaResponse); err != nil {
		return 0, err
	}

	if data, ok := fincaResponse["Data"].(map[string]interface{}); ok {
		return int(data["Id"].(float64)), nil
	}
	return 0, fmt.Errorf("error al obtener ID de la finca")
}

// crearParcelas maneja la creación de parcelas y su geolocalización
func crearParcelas(parcelas []interface{}, fincaID int) error {
	for _, parcela := range parcelas {
		p := parcela.(map[string]interface{})

		// Crear geolocalización
		geoID, err := crearGeolocalizacion(p["geolocalizacion"].(map[string]interface{}))
		if err != nil {
			return err
		}

		// Crear parcela
		parcelaID, err := crearParcela(p, fincaID)
		if err != nil {
			return err
		}

		// Relacionar parcela con finca
		err = relacionarFincaParcela(fincaID, parcelaID, geoID)
		if err != nil {
			return err
		}
	}
	return nil
}

// crearParcela registra una parcela en el API CRUD
func crearParcela(p map[string]interface{}, fincaID int) (int, error) {
	jsonParcela := map[string]interface{}{
		"NombreParcela": p["NombreParcela"],
		"TamañoParcela": p["TamañoParcela"],
		"FkFincaParcela": map[string]interface{}{
			"Id": fincaID,
		},
	}

	jsonParcelaByte, _ := json.Marshal(jsonParcela)
	response, err := services.Metodo_post("API_CRUD_FINCA", "/v1/Parcela", jsonParcelaByte)
	if err != nil {
		return 0, err
	}

	var parcelaResponse map[string]interface{}
	if err := json.Unmarshal(response, &parcelaResponse); err != nil {
		return 0, err
	}

	if data, ok := parcelaResponse["Data"].(map[string]interface{}); ok {
		return int(data["Id"].(float64)), nil
	}
	return 0, fmt.Errorf("error al obtener ID de la parcela")
}

// crearGeolocalizacion registra la ubicación en el API CRUD
func crearGeolocalizacion(geo map[string]interface{}) (int, error) {
	jsonGeo := map[string]interface{}{
		"LatitudInicial":  geo["LatitudInicial"],
		"LongitudInicial": geo["LongitudInicial"],
		"LatitudFinal":    geo["LatitudFinal"],
		"LongitudFinal":   geo["LongitudFinal"],
	}

	jsonGeoByte, _ := json.Marshal(jsonGeo)
	response, err := services.Metodo_post("API_CRUD_FINCA", "/v1/Geolocalizacion", jsonGeoByte)
	if err != nil {
		return 0, err
	}

	var geoResponse map[string]interface{}
	if err := json.Unmarshal(response, &geoResponse); err != nil {
		return 0, err
	}

	if data, ok := geoResponse["Data"].(map[string]interface{}); ok {
		return int(data["Id"].(float64)), nil
	}
	return 0, fmt.Errorf("error al obtener ID de la geolocalización")
}

// relacionarFincaParcela crea la relación entre finca, parcela y geolocalización
func relacionarFincaParcela(fincaID, parcelaID, geoID int) error {
	jsonRelacion := map[string]interface{}{
		"FkFincaParcela":    map[string]interface{}{"Id": fincaID},
		"FkParcelaFinca":    map[string]interface{}{"Id": parcelaID},
		"FkGeolocalizacion": map[string]interface{}{"Id": geoID},
	}

	jsonRelacionByte, _ := json.Marshal(jsonRelacion)
	_, err := services.Metodo_post("API_CRUD_FINCA", "/v1/FincaParcela", jsonRelacionByte)
	return err
}

// GetOne para mostrar finca
// @Title GetFinca
// @Description get Gestion_finca by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Gestion_finca
// @Failure 403 :id is empty
// @router /:buscarfinca/:porid/:id [get]
func (c *Gestion_fincaController) GetFinca() {
	fmt.Println("Función GET: Obtener finca por ID")

	// Obtener el ID de la finca desde la URL
	fincaID := c.Ctx.Input.Param(":id")
	if fincaID == "" {
		c.Data["json"] = map[string]interface{}{"error": "El ID de la finca es obligatorio"}
		c.ServeJSON()
		return
	}

	// Consultar la finca en el API CRUD FINCA
	response, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Finca/", fincaID)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al obtener la finca", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Parsear la respuesta JSON
	var fincaResponse map[string]interface{}
	if err := json.Unmarshal(response, &fincaResponse); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la respuesta", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Extraer los datos de la finca (se asume que la respuesta tiene la forma {"Data": { ... }})
	data, ok := fincaResponse["Data"].(map[string]interface{})
	if !ok {
		c.Data["json"] = map[string]interface{}{"error": "Finca no encontrada"}
		c.ServeJSON()
		return
	}

	// Construir la respuesta con los campos requeridos
	fincaInfo := map[string]interface{}{
		"Id":             data["Id"],
		"Nombre":         data["Nombre"],
		"AreaTotal":      data["AreaTotal"],
		"TotalParcelas":  data["TotalParcelas"],
		"TamañoParcelas": data["TamañoParcelas"],
	}

	// Si la finca incluye información del tipo de suelo, se extrae el nombre
	if fkFinca, exists := data["FkFinca"].(map[string]interface{}); exists {
		if nombre, ok := fkFinca["Nombre"].(string); ok {
			fincaInfo["TipoSuelo"] = nombre
		}
	}

	// Responder con la información de la finca
	c.Data["json"] = fincaInfo
	c.ServeJSON()
}

// GetOne para mostrar parcelas
// @Title GetParcelas
// @Description get Gestion_finca by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Gestion_finca
// @Failure 403 :id is empty
// @router /:parcelas/:id [get]
// La respuesta contendrá una lista de parcelas con su nombre, tamaño y geolocalización.
func (c *Gestion_fincaController) GetParcelas() {
	// Extraer el ID de la finca desde el parámetro de la URL.
	fincaID := c.Ctx.Input.Param(":id")
	if fincaID == "" {
		c.Data["json"] = map[string]interface{}{"error": "El ID de la finca es obligatorio"}
		c.ServeJSON()
		return
	}

	// Construir el query para obtener las relaciones en FincaParcela donde FkFincaParcela.Id = fincaID
	queryParam := "?query=FkFincaParcela.Id:" + fincaID

	// Realizar la consulta para obtener las parcelas
	response, err := services.Metodo_get("API_CRUD_FINCA", "/v1/FincaParcela", queryParam)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al obtener parcelas", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Parsear la respuesta JSON.
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la respuesta de parcelas", "details": err.Error()}
		c.ServeJSON()
		return
	}

	parcelasData, ok := result["Data"].([]interface{})
	if !ok {
		c.Data["json"] = map[string]interface{}{"mensaje": "No se encontraron parcelas para la finca"}
		c.ServeJSON()
		return
	}

	// Construir la respuesta filtrando los campos requeridos.
	var parcelas []map[string]interface{}
	for _, item := range parcelasData {
		parcelaMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// Se espera que en la relación FkParcelaFinca se encuentran los datos de la parcela.
		parcelaDatos, ok := parcelaMap["FkParcelaFinca"].(map[string]interface{})
		if !ok {
			continue
		}

		geoloc, _ := parcelaMap["FkGeolocalizacion"].(map[string]interface{})

		// Construir un mapa con los datos relevantes de la parcela.
		parcelaInfo := map[string]interface{}{
			"NombreParcela": parcelaDatos["NombreParcela"],
			"TamañoParcela": parcelaDatos["TamañoParcela"],
			"Geolocalizacion": map[string]interface{}{
				"LatitudInicial":  geoloc["LatitudInicial"],
				"LongitudInicial": geoloc["LongitudInicial"],
				"LatitudFinal":    geoloc["LatitudFinal"],
				"LongitudFinal":   geoloc["LongitudFinal"],
			},
		}
		parcelas = append(parcelas, parcelaInfo)
	}

	// Responder con la lista de parcelas filtradas.
	c.Data["json"] = parcelas
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Gestion_finca
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Gestion_finca
// @Failure 403
// @router / [get]
func (c *Gestion_fincaController) GetAll() {
	fmt.Println("Obteniendo todas las fincas")

	// Obtener el ID del usuario desde el token JWT
	token := c.Ctx.Input.Header("Authorization")
	if token == "" {
		c.Data["json"] = map[string]interface{}{"error": "No se proporcionó token"}
		c.ServeJSON()
		return
	}

	// Validar el token y extraer los claims
	claims, err := auth_JWT.ValidarJWT(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Token inválido o expirado"}
		c.ServeJSON()
		return
	}

	userID := claims.UserID // ID del usuario autenticado

	// Llamada al servicio para obtener todas las fincas
	response, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Finca", "")
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al obtener fincas"}
		c.ServeJSON()
		return
	}

	// Parsear la respuesta
	var data map[string]interface{}
	if err := json.Unmarshal(response, &data); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la respuesta"}
		c.ServeJSON()
		return
	}

	// Validar que hay datos disponibles
	if fincas, ok := data["Data"].([]interface{}); ok {
		var resultado []map[string]interface{}

		// Filtrar solo las fincas del usuario autenticado
		for _, finca := range fincas {
			fincaMap := finca.(map[string]interface{})
			if int(fincaMap["Id_Usuario"].(float64)) == userID { // Filtrado por usuario
				fincaInfo := map[string]interface{}{
					"Nombre":        fincaMap["Nombre"],
					"AreaTotal":     fincaMap["AreaTotal"],
					"TipoSuelo":     fincaMap["FkFinca"].(map[string]interface{})["Nombre"],
					"TotalParcelas": fincaMap["TotalParcelas"],
					"Id":            fincaMap["Id"],
				}
				resultado = append(resultado, fincaInfo)
			}
		}

		// Responder con la lista de fincas filtradas
		if len(resultado) > 0 {
			c.Data["json"] = resultado
		} else {
			c.Data["json"] = map[string]interface{}{"mensaje": "No tienes fincas registradas"}
		}
		c.ServeJSON()
	} else {
		c.Data["json"] = map[string]interface{}{"mensaje": "No hay fincas registradas"}
		c.ServeJSON()
	}
}

// Put ...
// @Title Put
// @Description update the Gestion_finca
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Gestion_finca	true		"body for Gestion_finca content"
// @Success 200 {object} models.Gestion_finca
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Gestion_fincaController) Put() {
	fmt.Println("Actualizar completamente una finca")

	// Obtener ID de la finca desde la URL
	id := c.Ctx.Input.Param(":id")
	if id == "" {
		c.Data["json"] = map[string]interface{}{"error": "ID de finca requerido"}
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

	// Obtener datos del body
	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la solicitud"}
		c.ServeJSON()
		return
	}

	// Validar que se envíen todos los campos requeridos
	requiredFields := []string{"Nombre", "AreaTotal", "TotalParcelas", "TamañoParcelas", "tipo_suelo"}
	for _, field := range requiredFields {
		if _, exists := body[field]; !exists {
			c.Data["json"] = map[string]interface{}{"error": fmt.Sprintf("Falta el campo requerido: %s", field)}
			c.ServeJSON()
			return
		}
	}

	// Obtener ID del tipo de suelo si se está actualizando
	if tipoSuelo, exists := body["tipo_suelo"].(string); exists {
		tipoSueloID, err := getTipoSueloID(tipoSuelo)
		if err != nil {
			c.Data["json"] = map[string]interface{}{"error": "Error al obtener el tipo de suelo"}
			c.ServeJSON()
			return
		}
		body["FkFinca"] = map[string]interface{}{"Id": tipoSueloID}
		body["Id_Usuario"] = userID
		delete(body, "tipo_suelo") // Eliminar campo redundante
	}

	// Enviar actualización completa
	jsonBody, _ := json.Marshal(body)
	response, err := services.Metodo_put("API_CRUD_FINCA", "/v1/Finca/", id, jsonBody)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al actualizar finca", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"mensaje": "Finca actualizada completamente", "response": string(response)}
	c.ServeJSON()
}

// Patch para actualizar parcialmente una finca
// @Title PatchFinca
// @Description Actualiza parcialmente una finca
// @Param	id		path 	string	true		"ID de la finca"
// @Param	body	body 	map[string]interface{}	true	"Cuerpo con los datos a actualizar"
// @Success 200 {object} map[string]interface{}
// @Failure 400 body is empty
// @router /:id [patch]
func (c *Gestion_fincaController) Patch() {
	fmt.Println("Actualizar parcialmente una finca")

	// Obtener ID de la finca desde la URL
	id := c.Ctx.Input.Param(":id")
	if id == "" {
		c.Data["json"] = map[string]interface{}{"error": "ID de finca requerido"}
		c.ServeJSON()
		return
	}

	// Obtener datos del body
	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la solicitud"}
		c.ServeJSON()
		return
	}

	if len(body) == 0 {
		c.Data["json"] = map[string]interface{}{"error": "No se enviaron campos para actualizar"}
		c.ServeJSON()
		return
	}
	// Si el campo "tipo_suelo" está presente, obtener su ID
	if tipoSuelo, exists := body["tipo_suelo"].(string); exists {
		tipoSueloID, err := getTipoSueloID(tipoSuelo)
		if err != nil {
			c.Data["json"] = map[string]interface{}{"error": "Error al obtener el tipo de suelo"}
			c.ServeJSON()
			return
		}
		body["FkFinca"] = map[string]interface{}{"Id": tipoSueloID}
		delete(body, "tipo_suelo") // Eliminar el campo redundante
	}

	// Enviar actualización parcial
	jsonBody, _ := json.Marshal(body)
	response, err := services.Metodo_patch("API_CRUD_FINCA", "/v1/Finca/", id, jsonBody)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al actualizar finca", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"mensaje": "Finca actualizada parcialmente", "response": string(response)}
	c.ServeJSON()
}

// Actualizar una parcela parcialmente (Patch)
// @Title PatchParcela
// @Description Actualiza solo los campos proporcionados de una parcela
// @Param	id		path 	int	true		"ID de la parcela a actualizar"
// @Param	body	body 	map[string]interface{}	true	"Datos de la parcela a actualizar"
// @Success 200 {object} map[string]interface{}
// @Failure 400 ID inválido o datos incorrectos
// @router /parcela/:id [patch]
func (c *Gestion_fincaController) PatchParcela() {
	fmt.Println("Actualizar parcialmente una parcela")

	idParcela := c.Ctx.Input.Param(":id")
	if idParcela == "" {
		c.respondWithError("ID de parcela no proporcionado")
		return
	}

	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.respondWithError("Error al procesar la solicitud")
		return
	}

	parcelaData := ParcelaData(body)

	response, err := services.Metodo_patch("API_CRUD_FINCA", "/v1/Parcela", idParcela, parcelaData)
	if err != nil {
		c.respondWithError("Error al actualizar la parcela", err.Error())
		return
	}

	// Actualizar geolocalización si se envía en la solicitud
	if geoData, exists := body["Geolocalizacion"].(map[string]interface{}); exists {
		if err := updateGeolocalizacion(geoData); err != nil {
			c.respondWithError("Error al actualizar la geolocalización", err.Error())
			return
		}
	}

	c.Data["json"] = map[string]interface{}{
		"mensaje":  "Parcela actualizada correctamente",
		"response": string(response),
	}
	c.ServeJSON()
}

// Construir los datos para la actualización de la parcela
func ParcelaData(body map[string]interface{}) []byte {
	data := make(map[string]interface{})
	if nombre, exists := body["NombreParcela"]; exists {
		data["NombreParcela"] = nombre
	}
	if tamaño, exists := body["TamañoParcela"]; exists {
		data["TamañoParcela"] = tamaño
	}
	jsonData, _ := json.Marshal(data)
	fmt.Println("Datos enviados:", string(jsonData))
	return jsonData
}

// Actualizar la geolocalización si existe en la solicitud
func updateGeolocalizacion(geoData map[string]interface{}) error {
	geoID, ok := geoData["id_geolocalizacion"].(float64)
	if !ok {
		return fmt.Errorf("ID de geolocalización inválido")
	}
	geoIDStr := strconv.Itoa(int(geoID))
	geoBody, _ := json.Marshal(geoData)
	fmt.Println("este es el body de geolocalizacion:", string(geoBody))
	_, err := services.Metodo_patch("API_CRUD_FINCA", "/v1/Geolocalizacion", geoIDStr, geoBody)
	return err
}

// **Sección del Arrendatario
// Registrar Finca
// @Title Create
// @Description create Gestion_finca
// @Param	body		body 	models.Gestion_finca	true		"body for Gestion_finca content"
// @Success 201 {object} models.Gestion_finca
// @Failure 403 body is empty
// @router /arrendatario/ [post]
func (c *Gestion_fincaController) Post_Arrendatario() {
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
// @router /arrendatario/arrendamiento/ [post]
func (c *Gestion_fincaController) Post_Arrendamiento() {
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

// Delete ...
// @Title Delete
// @Description delete the Gestion_finca
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Gestion_fincaController) Delete() {

}

// Manejo de errores reutilizable
func (c *Gestion_fincaController) respondWithError(msg string, details ...string) {
	errorResponse := map[string]string{"error": msg}
	if len(details) > 0 {
		errorResponse["detalle"] = details[0]
	}
	c.Data["json"] = errorResponse
	c.ServeJSON()
}
