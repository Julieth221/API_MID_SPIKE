package controllers

import (
	"encoding/json"
	"fmt"
	"log"
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

// GetTiposSueloUsuario
// @Title GetTiposSueloUsuario
// @Description Obtener todos los tipos de suelo registrados por el usuario autenticado
// @Success 200 {array} []map[string]interface{}
// @Failure 403 No se proporcionó token
// @router /tipos_suelo_usuario/usuario [get]
func (c *Gestion_fincaController) GetTiposSueloUsuario() {
	fmt.Println("Función GET: Tipos de suelo por usuario")

	// Obtener el token de autorización
	token := c.Ctx.Input.Header("Authorization")
	if token == "" {
		c.Data["json"] = map[string]interface{}{"error": "No se proporcionó token"}
		c.ServeJSON()
		return
	}

	// Validar token y obtener el ID del usuario
	claims, err := auth_JWT.ValidarJWT(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Token inválido o expirado"}
		c.ServeJSON()
		return
	}
	userID := claims.UserID

	// Construir la URL para consultar los tipos de suelo del usuario
	endpoint := fmt.Sprintf("?query=Id_Usuario:%d", userID)

	// Hacer la petición al API CRUD
	response, err := services.Metodo_get("API_CRUD_FINCA", "/v1/tipo_suelo", endpoint)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al consultar tipos de suelo", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Parsear la respuesta
	var parsedResp map[string]interface{}
	if err := json.Unmarshal(response, &parsedResp); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la respuesta", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Extraer y devolver los datos
	if tipos, ok := parsedResp["Data"].([]interface{}); ok {
		c.Data["json"] = tipos
	} else {
		c.Data["json"] = []interface{}{} // lista vacía si no hay tipos registrados
	}
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
		"TamañoParcelas": body["TamanoParcelas"],
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

// CrearParcelasParaFincaExistente
// @Title CrearParcelasParaFincaExistente
// @Description Crear parcelas adicionales para una finca ya registrada
// @Param   finca_id   query   int     true  "ID de la finca existente"
// @Param   body       body    []map[string]interface{} true "Lista de parcelas con geolocalización"
// @Success 201 {object} map[string]interface{}
// @Failure 400 finca_id inválido o faltan datos
// @router /crear_parcelas/finca_id [post]
func (c *Gestion_fincaController) CrearParcelasParaFincaExistente() {
	fmt.Println("Crear nuevas parcelas para finca existente")

	// Validar finca_id desde query param
	fincaIDStr := c.GetString("finca_id")
	if fincaIDStr == "" {
		c.Data["json"] = map[string]interface{}{"error": "Debe proporcionar el ID de la finca"}
		c.ServeJSON()
		return
	}

	fincaID, err := strconv.Atoi(fincaIDStr)
	if err != nil || fincaID <= 0 {
		c.Data["json"] = map[string]interface{}{"error": "ID de finca inválido"}
		c.ServeJSON()
		return
	}

	// Leer cuerpo con las nuevas parcelas
	var parcelas []interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &parcelas); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al leer las parcelas", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Reutilizar función existente
	if err := crearParcelas(parcelas, fincaID); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al crear parcelas", "details": err.Error()}
		c.ServeJSON()
		return
	}

	// Actualizar área total de la finca
	if err := actualizarAreaTotalFinca(fincaID); err != nil {
		fmt.Println("Advertencia: no se pudo actualizar el área total de la finca:", err)
	}

	c.Data["json"] = map[string]interface{}{"mensaje": "Parcelas creadas con éxito", "finca_id": fincaID}
	c.ServeJSON()
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
		"TamañoParcela": p["TamanoParcela"],
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
	query := fmt.Sprintf("?query=Id_Usuario:%d", userID)

	// Llamada al servicio para obtener todas las fincas
	response, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Finca", query)
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
				fmt.Println("fincas registradas del usuario: ", fincaInfo)

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
// @Param	body		body 	models.Gestion_finca	true		"body for Gestion_finca content"p
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

// Crear una nueva versión de la parcela
// @Title PostNuevaVersionParcela
// @Description Crea una nueva versión de la parcela referenciando a la anterior
// @Param	id		path 	int	true		"ID de la parcela original"
// @Param	body	body 	map[string]interface{}	true	"Datos de la nueva versión"
// @Success 200 {object} map[string]interface{}
// @Failure 400 Datos inválidos
// @router /parcela/versionar/:id [post]
func (c *Gestion_fincaController) PostNuevaVersionParcela() {
	fmt.Println("Creando nueva versión de la parcela")

	idParcelaAnterior := c.Ctx.Input.Param(":id") // Id de la parcela existente que se quiere versionar
	if idParcelaAnterior == "" {
		c.respondWithError("ID de parcela no proporcionado")
		return
	}

	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.respondWithError("Error al procesar el cuerpo de la solicitud")
		return
	}

	// 1. Obtener los datos de la parcela anterior
	parcelaAntData, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Parcela/", idParcelaAnterior)
	if err != nil {
		c.respondWithError("Error obteniendo la parcela anterior", err.Error())
		return
	}
	var parcelaAnterior map[string]interface{}
	_ = json.Unmarshal(parcelaAntData, &parcelaAnterior)
	idParcelaAnteriorInt, err := strconv.Atoi(idParcelaAnterior)
	if err != nil {
		c.respondWithError("ID de parcela anterior no es un número válido", err.Error())
		return
	}

	data, ok := parcelaAnterior["Data"].(map[string]interface{})
	if !ok {
		c.respondWithError("Error: no se encontró el campo Data en la respuesta de la parcela anterior")
		return
	}

	versionAnterior, ok := data["Version"].(float64)
	if !ok {
		c.respondWithError("Error: la versión de la parcela anterior no es válida o no existe")
		return
	}

	// 3. Desactivar la versión anterior
	data["Activo"] = false
	_ = updateParcelaActiva(idParcelaAnterior, data)

	// 4. Construir nueva parcela
	nuevaParcela := map[string]interface{}{
		"NombreParcela":  body["NombreParcela"],
		"TamanoParcela":  body["TamanoParcela"],
		"FkFincaParcela": map[string]interface{}{"Id": int(body["Finca"].(float64))},
		"MotivoCambio":   body["MotivoCambio"],
		"IdParcelaPadre": map[string]interface{}{"Id": idParcelaAnteriorInt},
		"Version":        int(versionAnterior) + 1, // Se incrementa la versión de forma segura
		"Activo":         true,
	}

	nuevaParcelaJSON, _ := json.Marshal(nuevaParcela)
	// fmt.Println("Este es el json que se va a envair a Parcela para crear la nueva version: ", string(nuevaParcelaJSON))
	resp, err := services.Metodo_post("API_CRUD_FINCA", "/v1/Parcela", nuevaParcelaJSON)
	if err != nil {
		c.respondWithError("Error al crear nueva versión de parcela", err.Error())
		return
	}

	// fmt.Println("Este es el reponse de crear una nueva version de parcela: ", string(resp))

	// 5. Crear geolocalización si está presente
	var geoID int
	if geoData, exists := body["Geolocalizacion"].(map[string]interface{}); exists {
		geoID, err = crearGeolocalizacion(geoData)
		if err != nil {
			c.respondWithError("Error al crear geolocalización", err.Error())
			return
		}
	}

	// 6. Obtener ID de la nueva parcela desde la respuesta
	var parcelaResponse map[string]interface{}
	_ = json.Unmarshal(resp, &parcelaResponse)
	nuevaParcelaID := int(parcelaResponse["Data"].(map[string]interface{})["Id"].(float64))

	// 7. Relacionar finca, parcela y geolocalización
	fincaID := int(body["Finca"].(float64))
	err = relacionarFincaParcela(fincaID, nuevaParcelaID, geoID)
	if err != nil {
		c.respondWithError("Error al relacionar FincaParcela", err.Error())
		return
	}

	// 8. Actualizar área total de finca
	if fincaID, ok := body["Finca"].(float64); ok {
		err := actualizarAreaTotalFinca(int(fincaID))
		if err != nil {
			fmt.Println("Advertencia al actualizar área total:", err)
		}
	}

	c.Data["json"] = map[string]interface{}{
		"mensaje":  "Nueva versión de parcela creada correctamente",
		"response": string(resp),
	}
	c.ServeJSON()
}

func updateParcelaActiva(id string, data map[string]interface{}) error {
	fmt.Println("Se va actualizar la parcela: ", id, data)
	parcelaJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = services.Metodo_put("API_CRUD_FINCA", "/v1/Parcela", id, parcelaJSON)
	return err
}

func actualizarAreaTotalFinca(fincaID int) error {
	log.Println("Este es un mensaje de depuración")

	fmt.Println("Antes del print de query params")
	query := fmt.Sprintf("?query=FkFincaParcela.Id:%d", fincaID)
	fmt.Printf("[DEBUG] Query params: %s\n", query)

	// Consultar la tabla FincaParcela
	resp, err := services.Metodo_get("API_CRUD_FINCA", "/v1/FincaParcela", query)
	if err != nil {
		return fmt.Errorf("error al obtener FincaParcela: %v", err)
	}
	fmt.Printf("[DEBUG] Respuesta cruda de FincaParcela: %s\n", string(resp))

	var fincaParcelaData map[string]interface{}
	if err := json.Unmarshal(resp, &fincaParcelaData); err != nil {
		return fmt.Errorf("error al deserializar respuesta: %v", err)
	}

	// Verifica si la respuesta contiene la propiedad que esperas
	data, exists := fincaParcelaData["Data"].([]interface{})
	if !exists {
		return fmt.Errorf("No se encontró el campo 'data' en la respuesta de FincaParcela")
	}

	// Inicializa variables para calcular el área total
	var areaTotal float64
	var totalParcelas int

	// Iterar sobre las parcelas para sumar el área total
	for _, parcela := range data {
		parcelaMap := parcela.(map[string]interface{})
		fkParcelaFinca, ok := parcelaMap["FkParcelaFinca"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("No se encontró FkParcelaFinca dentro de la parcela")
		}

		tamano, ok := fkParcelaFinca["TamanoParcela"].(float64)
		if !ok {
			return fmt.Errorf("No se encontró el tamaño de la parcela en FkParcelaFinca")
		}

		// Sumar el tamaño de la parcela al área total
		areaTotal += tamano
		totalParcelas++
	}

	// Construir el body con ambos campos
	body := map[string]interface{}{
		"AreaTotal":     areaTotal,
		"TotalParcelas": totalParcelas,
	}
	jsonBody, _ := json.Marshal(body)
	fmt.Printf("[DEBUG] Payload para PATCH: %s\n", string(jsonBody))

	// Actualizar la finca con los valores calculados
	idStr := strconv.Itoa(fincaID)
	_, err = services.Metodo_patch("API_CRUD_FINCA", "/v1/Finca", idStr, jsonBody)
	if err != nil {
		return fmt.Errorf("error al actualizar finca: %v", err)
	}

	fmt.Println("[DEBUG] Área total calculada:", areaTotal)
	fmt.Println("[DEBUG] Payload JSON a enviar:", string(jsonBody))

	return nil
}

// **Sección del Arrendatario
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

// GetDetalles obtiene una finca con todas sus parcelas y arrendamientos activos
// @Title GetDetalles
// @Description Obtener finca + parcelas + arrendamientos
// @Param    id      path    string  true        "ID de la finca"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404,500 {object} map[string]string
// @router /detalles/:id [get]
func (c *Gestion_fincaController) GetDetalles() {
	id := c.Ctx.Input.Param(":id")
	if id == "" {
		c.respondWithError("ID de finca obligatorio")
		return
	}

	// 1. Traer finca
	fincaResp, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Finca/", id)
	if err != nil {
		c.respondWithError("Error al obtener la finca", err.Error())
		return
	}
	var fincaData struct {
		Data struct {
			Id      int                    `json:"Id"`
			Nombre  string                 `json:"Nombre"`
			FkFinca map[string]interface{} `json:"FkFinca"`
		}
	}
	if err := json.Unmarshal(fincaResp, &fincaData); err != nil {
		c.respondWithError("Respuesta inválida de finca", err.Error())
		return
	}

	// 2. Traer parcelas de la finca
	qp := "?query=FkFincaParcela.Id:" + id
	parcelasResp, err := services.Metodo_get("API_CRUD_FINCA", "/v1/FincaParcela", qp)
	if err != nil {
		c.respondWithError("Error al obtener parcelas", err.Error())
		return
	}
	var fpData struct {
		Data []struct {
			FkParcelaFinca struct {
				Id            int     `json:"Id"`
				NombreParcela string  `json:"NombreParcela"`
				TamanoParcela float64 `json:"TamañoParcela"`
			} `json:"FkParcelaFinca"`
		}
	}
	if err := json.Unmarshal(parcelasResp, &fpData); err != nil {
		c.respondWithError("Respuesta inválida de parcelas", err.Error())
		return
	}

	// 3. Por cada parcela, revisar arrendamiento activo
	detalles := make([]map[string]interface{}, 0, len(fpData.Data))
	var AreaTotal float64
	for _, fp := range fpData.Data {
		dto := map[string]interface{}{
			"Id":            fp.FkParcelaFinca.Id,
			"NombreParcela": fp.FkParcelaFinca.NombreParcela,
			"TamanoParcela": fp.FkParcelaFinca.TamanoParcela,
		}
		AreaTotal += fp.FkParcelaFinca.TamanoParcela

		// arrendamiento activo
		query := fmt.Sprintf("?query=FkParcela:%d&Activo:true", fp.FkParcelaFinca.Id)
		arrResp, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento_Parcela", query)
		fmt.Println("esta es la respuesta de arrendamiento activo: ", string(arrResp))
		if err == nil {
			var arrData struct {
				Data []struct {
					FkArrendamiento struct {
						Id          int       `json:"Id"`
						FechaInicio time.Time `json:"FechaInicio"`
						FechaFin    time.Time `json:"FechaFin"`
						Activo      bool      `json:"Activo"`
					} `json:"FkArrendamiento"`
				}
			}
			if err2 := json.Unmarshal(arrResp, &arrData); err2 == nil && len(arrData.Data) > 0 {
				a := arrData.Data[0].FkArrendamiento
				dto["ArrendamientoId"] = a.Id
				dto["FechaInicio"] = a.FechaInicio
				dto["FechaFin"] = a.FechaFin
				dto["Activo"] = a.Activo
			}
		}
		detalles = append(detalles, dto)
	}

	// 4. Responder
	result := map[string]interface{}{
		"Id":            fincaData.Data.Id,
		"Nombre":        fincaData.Data.Nombre,
		"TipoSuelo":     fincaData.Data.FkFinca["Nombre"],
		"TotalParcelas": len(detalles),
		"AreaTotal":     AreaTotal,
		"Parcelas":      detalles,
	}
	c.Data["json"] = result
	c.ServeJSON()
}

// DesactivarArrendamiento marca como inactivo un arrendamiento existente
// @Title DesactivarArrendamiento
// @Description Pone Activo=false en un arrendamiento
// @Param    id      path    string  true        "ID del arrendamiento"
// @Success 200 {object} map[string]string
// @Failure 400,404,500 {object} map[string]string
// @router /arrendamiento/desactivar/:id [patch]
func (c *Gestion_fincaController) DesactivarArrendamiento() {
	id := c.Ctx.Input.Param(":id")
	if id == "" {
		c.respondWithError("ID de arrendamiento requerido")
		return
	}

	// Construir payload
	payload := map[string]interface{}{"Activo": false}
	data, _ := json.Marshal(payload)

	// PATCH al CRUD
	_, err := services.Metodo_patch("API_CRUD_FINCA", "/v1/Arrendamiento", id, data)
	if err != nil {
		c.respondWithError("Error al desactivar arrendamiento", err.Error())
		return
	}

	c.Data["json"] = map[string]string{"mensaje": "Arrendamiento desactivado"}
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
