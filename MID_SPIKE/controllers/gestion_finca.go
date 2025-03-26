package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

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
	c.Mapping("GetOne", c.GetOne)
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
// @Title GetOne
// @Description get Gestion_finca by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Gestion_finca
// @Failure 403 :id is empty
// @router /nombrefinca [get]
func (c *Gestion_fincaController) GetOne() {
	fmt.Println("Este es el GetOne para finca")

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
