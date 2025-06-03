package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	auth_JWT "github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/auth_jwt"
	"github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/services"
)

// Gestion_cultivoController operations for Gestion_cultivo
type Gestion_cultivoController struct {
	beego.Controller
}

// URLMapping ...
func (c *Gestion_cultivoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Registrar nuevo cultivo
// @Title Create
// @Description create Gestion_cultivo
// @Param	body		body 	models.Gestion_cultivo	true		"body for Gestion_cultivo content"
// @Success 201 {object} models.Gestion_cultivo
// @Failure 403 body is empty
// @router / [post]
func (c *Gestion_cultivoController) Post() {
	fmt.Println("Registrar cultivo")

	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar la solicitud"}
		c.ServeJSON()
		return
	}

	// Obtener token y validar
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
	fmt.Println("Usuario autenticado:", userID)

	// Validar campos obligatorios
	idParcela, ok := body["Id_Parcela"].(float64)
	if !ok {
		c.Data["json"] = map[string]interface{}{"error": "Falta el campo id_parcela"}
		c.ServeJSON()
		return
	}

	// Validar si el usuario tiene derecho sobre la parcela
	// Si viene con fk_arrendamiento, validar que esa parcela esté en ese contrato
	// Si no viene, validar que la parcela pertenezca a una finca del propietario
	if fkArrRaw, exists := body["Id_Arrendamiento"]; exists && fkArrRaw != nil {
		// Validar que la parcela esté en el arrendamiento
		fkArrendamiento := int(fkArrRaw.(float64))
		isValida, err := validarParcelaEnArrendamiento(fkArrendamiento, int(idParcela))
		if err != nil || !isValida {
			c.Data["json"] = map[string]interface{}{"error": "La parcela no está asociada a ese arrendamiento o hay un error"}
			c.ServeJSON()
			return
		}
	} else {
		// Validar que la parcela pertenece a una finca del propietario
		esDelUsuario, err := validarParcelaDelPropietario(userID, int(idParcela))
		if err != nil || !esDelUsuario {
			c.Data["json"] = map[string]interface{}{"error": "La parcela no pertenece a ninguna finca del usuario"}
			c.ServeJSON()
			return
		}
	}

	fkTipoArrozRaw, ok := body["FkTipoArroz"].(float64)
	if !ok {
		c.Data["json"] = map[string]interface{}{"error": "Falta o es inválido el campo fk_tipo_arroz"}
		c.ServeJSON()
		return
	}

	fkMetodoSiembraRaw, ok := body["FkMetodoSiembra"].(float64)
	if !ok {
		c.Data["json"] = map[string]interface{}{"error": "Falta o es inválido el campo fk_metodo_siembra"}
		c.ServeJSON()
		return
	}

	fkEstadoRaw, ok := body["FkEstadoFenologicoCultivo"].(float64)
	if !ok {
		c.Data["json"] = map[string]interface{}{"error": "Falta o es inválido el campo fk_estado_fenologico_cultivo"}
		c.ServeJSON()
		return
	}

	// Construir json para enviar al API CRUD
	cultivo := map[string]interface{}{
		"FkTipoArroz":               map[string]interface{}{"Id": int(fkTipoArrozRaw)},
		"FkMetodoSiembra":           map[string]interface{}{"Id": int(fkMetodoSiembraRaw)},
		"FkEstadoFenologicoCultivo": map[string]interface{}{"Id": int(fkEstadoRaw)},
		"Id_Parcela":                int(idParcela),
		"FechaSiembra":              body["FechaSiembra"],
		"CicloDias":                 body["CicloDias"],
		"AreaSembrada":              body["AreaSembrada"],
		"Nombre":                    body["NombreCultivo"],
		"Id_Usuario":                userID,
	}

	// Id_Arrendamiento si aplica, mandar entero no objeto
	if fkArrRaw, exists := body["Id_Arrendamiento"]; exists && fkArrRaw != nil {
		cultivo["Id_Arrendamiento"] = int(fkArrRaw.(float64))
	}

	// Hacer POST al API CRUD CULTIVO
	cultivoJSON, _ := json.Marshal(cultivo)
	fmt.Println("Este es el json para regitrar cultivo: ", string(cultivoJSON))
	response, err := services.Metodo_post("API_CRUD_CULTIVO", "/v1/Registro_Cultivo", cultivoJSON)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al registrar cultivo", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	var parsedResp map[string]interface{}
	if err := json.Unmarshal(response, &parsedResp); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al leer respuesta del servidor", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	dataRaw := parsedResp["Data"]
	dataSlice, ok := dataRaw.([]interface{})
	if !ok || len(dataSlice) == 0 {
		c.Data["json"] = map[string]interface{}{"error": "Formato inesperado o vacío en el campo Data"}
		c.ServeJSON()
		return
	}

	// Acceder al primer objeto del arreglo
	dataMap, ok := dataSlice[0].(map[string]interface{})
	if !ok {
		c.Data["json"] = map[string]interface{}{"error": "Elemento en Data no es un objeto esperado"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"data": dataMap}
	c.ServeJSON()

}

// validarParcelaEnArrendamiento consulta si la parcela está en ese contrato
func validarParcelaEnArrendamiento(fkArrendamiento int, idParcela int) (bool, error) {
	endpoint := fmt.Sprintf("?query=FkArrendamiento.Id:%d,Id_Parcela:%d", fkArrendamiento, idParcela)
	resp, err := services.Metodo_get("API_CRUD_CULTIVO", "/v1/Arrendamiento_Parcela", endpoint)
	if err != nil {
		return false, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(resp, &data); err != nil {
		return false, err
	}

	if results, ok := data["Data"].([]interface{}); ok && len(results) > 0 {
		return true, nil
	}
	return false, nil
}

// validarParcelaDelPropietario consulta si la parcela está en finca del usuario
func validarParcelaDelPropietario(userID int, idParcela int) (bool, error) {
	// Paso 1: consultar parcela para obtener id finca
	endpointParcela := fmt.Sprintf("?query=Id:%d", idParcela)
	respParcela, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Parcela", endpointParcela)
	if err != nil {
		return false, err
	}
	var dataParcela map[string]interface{}
	if err := json.Unmarshal(respParcela, &dataParcela); err != nil {
		return false, err
	}
	if dataList, ok := dataParcela["Data"].([]interface{}); ok && len(dataList) > 0 {
		parcela := dataList[0].(map[string]interface{})
		finca := parcela["FkFincaParcela"].(map[string]interface{})
		idFinca := int(finca["Id"].(float64))

		// Paso 2: consultar finca y validar usuario
		endpointFinca := fmt.Sprintf("?query=Id:%d,Id_Usuario:%d", idFinca, userID)
		respFinca, err := services.Metodo_get("API_CRUD_FINCA", "/v1/Finca", endpointFinca)
		if err != nil {
			return false, err
		}
		var dataFinca map[string]interface{}
		if err := json.Unmarshal(respFinca, &dataFinca); err != nil {
			return false, err
		}
		if fincaList, ok := dataFinca["Data"].([]interface{}); ok && len(fincaList) > 0 {
			return true, nil
		}
	}
	return false, nil
}

// GetOne ...
// @Title GetOne
// @Description get Gestion_cultivo by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Gestion_cultivo
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Gestion_cultivoController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Gestion_cultivo
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Gestion_cultivo
// @Failure 403
// @router / [get]
func (c *Gestion_cultivoController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Gestion_cultivo
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Gestion_cultivo	true		"body for Gestion_cultivo content"
// @Success 200 {object} models.Gestion_cultivo
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Gestion_cultivoController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Gestion_cultivo
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Gestion_cultivoController) Delete() {

}

func (c *Gestion_cultivoController) respondWithError(msg string, details ...string) {
	errorResponse := map[string]string{"error": msg}
	if len(details) > 0 {
		errorResponse["detalle"] = details[0]
	}
	c.Data["json"] = errorResponse
	c.ServeJSON()
}
