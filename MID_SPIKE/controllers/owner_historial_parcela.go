package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/services"
)

// Owner_historial_parcelaController operations for Owner_historial_parcela
type Owner_historial_parcelaController struct {
	beego.Controller
}

// URLMapping ...
func (c *Owner_historial_parcelaController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Owner_historial_parcela
// @Param	body		body 	models.Owner_historial_parcela	true		"body for Owner_historial_parcela content"
// @Success 201 {object} models.Owner_historial_parcela
// @Failure 403 body is empty
// @router / [post]
func (c *Owner_historial_parcelaController) Post() {

}

// GetParcelasOriginalesPorFinca
// @Title GetParcelasOriginalesPorFinca
// @Description Obtener parcelas originales (sin padre) por ID de finca
// @Param   id    path    int  true  "ID de la finca"
// @Success 200 {object} []map[string]interface{}
// @Failure 403 :id inválido
// @router /parcelas/originales/:id [get]
func (c *Owner_historial_parcelaController) GetParcelasOriginalesPorFinca() {
	fmt.Println("Parcelas originales por finca")

	fincaID := c.Ctx.Input.Param(":id")
	if fincaID == "" {
		c.respondWithError("El ID de la finca es requerido")
		return
	}

	parcelas, err := obtenerParcelasOriginalesPorFinca(fincaID)
	if err != nil {
		c.respondWithError("Error al obtener parcelas originales", err.Error())
		return
	}

	var resultado []map[string]interface{}
	for _, p := range parcelas {
		resultado = append(resultado, map[string]interface{}{
			"Id":     p["Id"],
			"Nombre": p["NombreParcela"],
		})
	}

	c.Data["json"] = resultado
	c.ServeJSON()
}

func obtenerParcelasOriginalesPorFinca(fincaID string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("?query=FkFincaParcela.Id:%s", fincaID)
	resp, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Parcela", query)
	if err != nil {
		return nil, fmt.Errorf("error en la petición al API CRUD: %w", err)
	}

	var resultado map[string]interface{}
	if err := json.Unmarshal(resp, &resultado); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta del API CRUD: %w", err)
	}

	data, ok := resultado["Data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("respuesta inesperada del API CRUD: campo 'Data' inválido")
	}

	var parcelas []map[string]interface{}
	for _, item := range data {
		if p, ok := item.(map[string]interface{}); ok {
			// Verifica si no tiene FkParcelaPadre o está en nil
			if p["IdParcelaPadre"] == nil {
				parcelas = append(parcelas, p)
			}
		}
	}

	return parcelas, nil
}

// GetOne ...
// @Title GetOne
// @Description get Gestion_arrendamiento by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Gestion_arrendamiento
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Owner_historial_parcelaController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Owner_historial_parcela
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Owner_historial_parcela
// @Failure 403
// @router parcelaversiones/:id [get]
func (c *Owner_historial_parcelaController) GetAll() {
	idParcela := c.Ctx.Input.Param(":id")
	if idParcela == "" {
		c.respondWithError("El ID de la parcela es requerido")
		return
	}

	todasVersiones := buscarLineaVersionActiva(idParcela)
	if len(todasVersiones) == 0 {
		c.Data["json"] = []map[string]interface{}{}
		c.ServeJSON()
		return
	}

	var historial []map[string]interface{}

	for _, parcela := range todasVersiones {
		if len(parcela) == 0 {
			continue
		}

		idRaw, ok := parcela["Id"]
		if !ok {
			continue
		}
		idParcelaHija := int(idRaw.(float64))

		nombre := parcela["NombreParcela"].(string)
		version := int(parcela["Version"].(float64))
		tamano := parcela["TamanoParcela"].(float64)
		motivo := parcela["MotivoCambio"]
		fecha := parcela["FechaCreacion"].(string)

		arrRes, _ := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento_Parcela", fmt.Sprintf("?query=FkParcela.Id:%d", idParcelaHija))
		var arrendamientoParcelas []map[string]interface{}
		if err := json.Unmarshal(arrRes, &arrendamientoParcelas); err != nil {
			continue
		}

		arrendamientos := []map[string]interface{}{}

		for _, arrParcela := range arrendamientoParcelas {
			valor := arrParcela["Valor"].(string)
			arrObj := arrParcela["FkArrendamiento"].(map[string]interface{})

			activo := arrObj["Activo"].(bool)
			fechaInicioStr := arrObj["FechaInicio"].(string)
			fechaFinStr := arrObj["FechaFin"].(string)

			fechaInicio, err := time.Parse(time.RFC3339, fechaInicioStr)
			if err != nil {
				fmt.Println("Error parseando FechaInicio:", err)
			}
			fechaFin, err := time.Parse(time.RFC3339, fechaFinStr)
			if err != nil {
				fmt.Println("Error parseando FechaFin:", err)
			}

			estado := calcularEstadoArrendamiento(activo, fechaInicio, fechaFin)

			arrendatario := "-"
			if user, ok := arrObj["IdUserUserArrendatario"].(map[string]interface{}); ok {
				arrendatario = user["Nombre"].(string)
			}

			// 🔍 Buscar si el arrendamiento tiene otras parcelas asociadas
			idArrendamiento := int(arrObj["Id"].(float64))
			otrasParcelasRes, _ := services.Metodo_get("API_CRUD_FINCA", "/v1/Arrendamiento_Parcela", fmt.Sprintf("?query=FkArrendamiento.Id:%d", idArrendamiento))
			var otrasParcelas []map[string]interface{}
			if err := json.Unmarshal(otrasParcelasRes, &otrasParcelas); err != nil {
				fmt.Println("Error parseando parcelas adicionales:", err)
			}

			var parcelasAdicionales []string
			for _, p := range otrasParcelas {
				if fkParcela, ok := p["FkParcela"].(map[string]interface{}); ok {
					if int(fkParcela["Id"].(float64)) != idParcelaHija {
						nombreParcela := fkParcela["NombreParcela"].(string)
						parcelasAdicionales = append(parcelasAdicionales, nombreParcela)
					}
				}
			}

			arrendamientos = append(arrendamientos, map[string]interface{}{
				"id":                  idArrendamiento,
				"estado":              estado,
				"fechaInicio":         fechaInicioStr,
				"fechaFin":            fechaFinStr,
				"valor":               valor,
				"arrendatario":        arrendatario,
				"parcelasAdicionales": parcelasAdicionales,
			})
		}

		historial = append(historial, map[string]interface{}{
			"nombre":         nombre,
			"version":        version,
			"tamano":         tamano,
			"motivoCambio":   motivo,
			"fechaCreacion":  fecha,
			"arrendamientos": arrendamientos,
		})
	}

	c.Data["json"] = historial
	c.ServeJSON()
}

func calcularEstadoArrendamiento(activo bool, inicio, fin time.Time) string {
	ahora := time.Now()
	if activo && ahora.After(inicio) && ahora.Before(fin) {
		return "Activo"
	}
	return "Inactivo"
}

func buscarLineaVersionActiva(id string) []map[string]interface{} {
	var historial []map[string]interface{}

	// 1. Consultar la parcela original (la raíz)
	query := fmt.Sprintf("?query=Id:%s", id)
	res, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Parcela", query)
	if err != nil {
		return historial
	}

	var respuesta map[string]interface{}
	if err := json.Unmarshal(res, &respuesta); err != nil {
		return historial
	}

	data, ok := respuesta["Data"].([]interface{})
	if !ok || len(data) == 0 {
		return historial
	}

	// Agregar la parcela original al historial
	parcelaOriginal, ok := data[0].(map[string]interface{})
	if !ok {
		return historial
	}
	historial = append(historial, parcelaOriginal)

	// 2. Buscar versiones hijas en cadena
	actual := id
	for {
		queryHijo := fmt.Sprintf("?query=IdParcelaPadre.Id:%s", actual)
		resHijo, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Parcela", queryHijo)
		if err != nil {
			break
		}

		var respuestaHijo map[string]interface{}
		if err := json.Unmarshal(resHijo, &respuestaHijo); err != nil {
			break
		}

		dataHijo, ok := respuestaHijo["Data"].([]interface{})
		if !ok || len(dataHijo) == 0 {
			break
		}

		// Solo consideramos la primera versión hija (como estaba en tu código original)
		parcelaHija, ok := dataHijo[0].(map[string]interface{})
		if !ok {
			break
		}

		historial = append(historial, parcelaHija)

		// Validar si es activa antes de romper el ciclo
		if activoRaw, ok := parcelaHija["Activo"]; ok {
			if activoBool, ok := activoRaw.(bool); ok && activoBool {
				break
			}
		}

		// Obtener ID de la hija para continuar
		if idRaw, ok := parcelaHija["Id"]; ok {
			if idFloat, ok := idRaw.(float64); ok {
				actual = fmt.Sprintf("%.0f", idFloat)
			} else {
				break
			}
		} else {
			break
		}
	}

	return historial
}

// Put ...
// @Title Put
// @Description update the Owner_historial_parcela
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Owner_historial_parcela	true		"body for Owner_historial_parcela content"
// @Success 200 {object} models.Owner_historial_parcela
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Owner_historial_parcelaController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Owner_historial_parcela
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Owner_historial_parcelaController) Delete() {

}

func (c *Owner_historial_parcelaController) respondWithError(msg string, details ...string) {
	errorResponse := map[string]string{"error": msg}
	if len(details) > 0 {
		errorResponse["detalle"] = details[0]
	}
	c.Data["json"] = errorResponse
	c.ServeJSON()
}
