package controllers

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

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
// @router /cultivos/:id [get]
func (c *Gestion_cultivoController) GetAll() {
	parcelaID := c.Ctx.Input.Param(":id")
	fmt.Printf("GetAll de cultivos para parcela %s\n", parcelaID)

	// 1. Validar token y obtener usuario y rol
	token := c.Ctx.Input.Header("Authorization")
	claims, err := auth_JWT.ValidarJWT(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		c.Data["json"] = map[string]string{"error": "Token inválido"}
		c.ServeJSON()
		return
	}
	userID := claims.UserID
	rol := claims.Role // "Propietario", "Arrendatario" o "Admin"

	// 2. Validar acceso a la parcela según rol
	idParcelaInt, err := strconv.Atoi(parcelaID)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Data["json"] = map[string]string{"error": "ID de parcela inválido"}
		c.ServeJSON()
		return
	}
	if rol == "Arrendatario" {
		okArr, err := validarParcelaEnArrendamiento(userID, idParcelaInt)
		if err != nil || !okArr {
			c.Ctx.ResponseWriter.WriteHeader(403)
			c.Data["json"] = map[string]string{"error": "Acceso denegado: no arrendatario de esta parcela"}
			c.ServeJSON()
			return
		}
	}
	if rol == "PROPIETARIO" {
		okProp, err := validarParcelaDelPropietario(userID, idParcelaInt)
		if err != nil || !okProp {
			c.Ctx.ResponseWriter.WriteHeader(403)
			c.Data["json"] = map[string]string{"error": "Acceso denegado: no propietario de esta parcela"}
			c.ServeJSON()
			return
		}
	}

	// 3. Obtener datos de finca/parcela para mostrar nombres
	var fincaName, parcelaName string
	respParc, _ := services.Metodo_get("API_CRUD_FINCA", "/v1/Parcela", fmt.Sprintf("?query=Id:%d,Activo:true", idParcelaInt))
	var parcData map[string]interface{}
	json.Unmarshal(respParc, &parcData)
	if list, ok := parcData["Data"].([]interface{}); ok && len(list) > 0 {
		p := list[0].(map[string]interface{})
		parcelaName = p["NombreParcela"].(string)
		if fincaObj, exists := p["FkFincaParcela"].(map[string]interface{}); exists {

			if name, ok2 := fincaObj["Nombre"].(string); ok2 {
				fincaName = name
			}

		}
	}

	// 4. Obtener cultivos asociados
	endpoint := fmt.Sprintf("?query=Id_Parcela:%s", parcelaID)
	resp, err := services.Metodo_get("API_CRUD_CULTIVO", "/v1/Registro_Cultivo", endpoint)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Error al obtener cultivos"}
		c.ServeJSON()
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(resp, &data); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.Data["json"] = map[string]string{"error": "Error al parsear cultivos"}
		c.ServeJSON()
		return
	}

	// Extraer “Data” con seguridad
	val, exists := data["Data"]
	if !exists || val == nil {
		// API devolvió null o no incluyó “Data”
		c.Data["json"] = map[string]interface{}{"data": []interface{}{}}
		c.ServeJSON()
		return
	}

	rawList, ok := data["Data"].([]interface{})
	if !ok || len(rawList) == 0 {
		// No hay registros → devolvemos array vacío directamente
		c.Data["json"] = map[string]interface{}{"data": []interface{}{}}
		c.ServeJSON()
		return
	}

	// 5. Procesar cada cultivo: etapa actual, estado e ingeniería de respuesta
	var results []map[string]interface{}
	hoy := time.Now()
	estadosFinales := map[string]bool{"Cosecha": true}

	for _, item := range rawList {

		if item == nil {
			continue
		}

		// 1) Asegurarnos de que el elemento es un map válido
		cultivo, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemMap, ok := item.(map[string]interface{})
		if !ok || len(itemMap) == 0 {
			continue // Saltar si está vacío o no es el tipo esperado
		}

		// 2) ID de cultivo (si no existe, saltamos)
		idRaw, ok := cultivo["Id"].(float64)
		if !ok {
			continue
		}
		cultID := int(idRaw)

		// 3) Campo opcional Id_Arrendamiento
		arrIDf, hasArr := cultivo["Id_Arrendamiento"].(float64)

		var tipoArroz string
		if fkTipo, ok := cultivo["FkTipoArroz"].(map[string]interface{}); ok {
			if nombre, ok2 := fkTipo["Nombre"].(string); ok2 {
				tipoArroz = nombre
			}
		}

		// 5a. Lógica de rol para incluir/excluir cultivos
		if rol == "Arrendatario" {
			if !hasArr {
				continue
			}
			// validar arrendatario y vigencia
			respArr, _ := services.Metodo_get("API_CRUD_CULTIVO", "/v1/Arrendamiento", fmt.Sprintf("?query=Id:%d,Activo:true", int(arrIDf)))
			var arrData map[string]interface{}
			json.Unmarshal(respArr, &arrData)
			arrList, _ := arrData["Data"].([]interface{})
			if len(arrList) == 0 {
				continue
			}
			details := arrList[0].(map[string]interface{})
			if int(details["IdUserUserArrendatario"].(float64)) != userID {
				continue
			}
		}
		if rol == "Propietario" {
			// propietario solo cultivos sin arrendamiento activo
			if hasArr {
				respArr, _ := services.Metodo_get("API_CRUD_CULTIVO", "/v1/Arrendamiento", fmt.Sprintf("?query=Id:%d,Activo:true", int(arrIDf)))
				var d map[string]interface{}
				json.Unmarshal(respArr, &d)
				if arrList, _ := d["Data"].([]interface{}); len(arrList) > 0 {
					continue
				}
			}
		}

		// 5b. Determinar etapa fenológica actual
		respFases, _ := services.Metodo_get("API_CRUD_CULTIVO", "/v1/Cultivo_Fase", fmt.Sprintf("?query=FkCultivoFase.Id:%d,Activo:true", int(cultID)))
		var fasesData map[string]interface{}
		json.Unmarshal(respFases, &fasesData)
		fasesList, _ := fasesData["Data"].([]interface{})
		sort.Slice(fasesList, func(i, j int) bool {
			fi := parseDate(fasesList[i].(map[string]interface{})["FechaInicio"].(string))
			fj := parseDate(fasesList[j].(map[string]interface{})["FechaInicio"].(string))
			return fi.Before(fj)
		})
		var actualFen string
		for idx, f := range fasesList {
			fase := f.(map[string]interface{})
			if fase["Completada"].(bool) && idx+1 < len(fasesList) {
				actualFen = fasesList[idx+1].(map[string]interface{})["FkFaseCultivo"].(map[string]interface{})["NombreFase"].(string)
			}
		}
		if actualFen == "" && len(fasesList) > 0 {
			actualFen = fasesList[0].(map[string]interface{})["FkFaseCultivo"].(map[string]interface{})["NombreFase"].(string)
		}

		// 5c. Estado activo/inactivo y sugerencia de nuevo arrendamiento

		// Activo
		activo, ok := cultivo["Activo"].(bool)
		if !ok {
			activo = false
		}
		sugerir := false
		var nuevoArr map[string]interface{}
		if hasArr {
			respArr2, _ := services.Metodo_get("API_CRUD_CULTIVO", "/v1/Arrendamiento", fmt.Sprintf("?query=Id:%d", int(arrIDf)))
			var arrData2 map[string]interface{}
			json.Unmarshal(respArr2, &arrData2)
			if arrList2, _ := arrData2["Data"].([]interface{}); len(arrList2) > 0 {
				arr := arrList2[0].(map[string]interface{})
				if hoy.After(parseDate(arr["FechaFin"].(string))) {
					estadoFen := cultivo["FkEstadoFenologicoCultivo"].(map[string]interface{})["Nombre"].(string)
					if estadosFinales[estadoFen] {
						activo = false
					} else {
						respSig, _ := services.Metodo_get("API_CRUD_CULTIVO", "/v1/Arrendamiento", fmt.Sprintf("?query=fk_arr_anterior.Id:%d", int(arrIDf)))
						var sigData map[string]interface{}
						json.Unmarshal(respSig, &sigData)
						if nextList, _ := sigData["Data"].([]interface{}); len(nextList) > 0 {
							sugerir = true
							nuevoArr = nextList[0].(map[string]interface{})
						} else {
							activo = false
						}
					}
				}
			}
		}

		// 5d. Armar objeto de respuesta
		obj := map[string]interface{}{
			"id":                           cultID,
			"nombre":                       cultivo["Nombre"],
			"finca":                        fincaName,
			"parcela":                      parcelaName,
			"areaSembrada":                 cultivo["AreaSembrada"],
			"tipoArroz":                    tipoArroz,
			"fechaSiembra":                 cultivo["FechaSiembra"],
			"etapaFenologica":              actualFen,
			"activo":                       activo,
			"sugerenciaNuevoArrendamiento": sugerir,
		}
		if sugerir {
			obj["nuevoArrendamiento"] = nuevoArr
		}
		results = append(results, obj)
	}

	c.Data["json"] = map[string]interface{}{"data": results}
	c.ServeJSON()
}

// parseDate: helper para convertir string a time.Time
func parseDate(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

// GetOne ...
// @Title GetOne
// @Description get Gestion_cultivo by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Gestion_cultivo
// @Failure 403 :id is empty
// @router /insumo/porcultivo/:id [get]
func (c *Gestion_cultivoController) InsumosPorCultivo() {
	fmt.Println("esta es una funcion para saber los insumos que se han aplicado en un cultivo")
	cultivoID := c.Ctx.Input.Param(":id")

	// 1) Llamar a la API CRUD para traer insumos
	endpoint := fmt.Sprintf("?query=FkRegistroCultivo.Id:%s,Activo:true", cultivoID)
	resp, err := services.Metodo_get("API_CRUD_CULTIVO", "/v1/Insumo", endpoint)
	if err != nil {
		c.respondWithError("Error al obtener insumos del cultivo", err.Error())
		return
	}

	// 2) Parsear JSON base
	var data map[string]interface{}
	if err := json.Unmarshal(resp, &data); err != nil {
		c.respondWithError("Error al parsear respuesta de insumos", err.Error())
		return
	}

	// 3) Extraer slice de forma segura
	raw, ok := data["Data"].([]interface{})
	if !ok || len(raw) == 0 {
		// No hay elementos o Data no es un slice → devolvemos vacío
		c.Data["json"] = map[string]interface{}{"data": []interface{}{}}
		c.ServeJSON()
		return
	}

	// 4) Procesar cada elemento con safe-casting
	results := make([]map[string]interface{}, 0, len(raw))
	for _, item := range raw {
		insumoMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		// Si es {} completamente vacío saltamos
		if len(insumoMap) == 0 {
			continue
		}

		// Extraer campo NombreInsumo (clave mínima para identificar un insumo)
		nombreInsumo, ok := insumoMap["NombreInsumo"].(string)
		if !ok {
			continue // sin nombre no es un registro válido
		}

		// Safe-get de cada nested map
		var tipoInsumo, categoriaInsumo, metodoAplicacion string
		if m, ok := insumoMap["FkTipoInsumo"].(map[string]interface{}); ok {
			if v, ok2 := m["Nombre"].(string); ok2 {
				tipoInsumo = v
			}
		}
		if m, ok := insumoMap["FkCategoriaInsumo"].(map[string]interface{}); ok {
			if v, ok2 := m["Nombre"].(string); ok2 {
				categoriaInsumo = v
			}
		}
		if m, ok := insumoMap["FkMetodoAplicacionInsumo"].(map[string]interface{}); ok {
			if v, ok2 := m["Nombre"].(string); ok2 {
				metodoAplicacion = v
			}
		}

		// Campos simples
		cantidadAplicada := insumoMap["CantidadAplicada"]
		fechaAplicacion := insumoMap["FechaAplicacion"]
		observaciones := insumoMap["Observaciones"]

		// Construir JSON de salida
		result := map[string]interface{}{
			"nombreInsumo":     nombreInsumo,
			"tipoInsumo":       tipoInsumo,
			"categoriaInsumo":  categoriaInsumo,
			"metodoAplicacion": metodoAplicacion,
			"cantidadAplicada": cantidadAplicada,
			"fechaAplicacion":  fechaAplicacion,
			"observaciones":    observaciones,
		}
		results = append(results, result)
	}

	// 5) Responder
	c.Data["json"] = map[string]interface{}{"data": results}
	c.ServeJSON()

}

// GetSeleccionables ...
// @Title GetSeleccionables
// @Description get Gestion_cultivo by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Gestion_cultivo
// @Failure 403 :id is empty
// @router /parcelas/seleccionables [get]
func (c *Gestion_cultivoController) GetSeleccionables() {
	fmt.Println("Funcion para llenar el select de parcelas")
	// 1. Validar JWT
	token := c.Ctx.Input.Header("Authorization")
	claims, err := auth_JWT.ValidarJWT(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		c.Data["json"] = map[string]string{"error": "Token inválido"}
		c.ServeJSON()
		return
	}
	userID := claims.UserID
	role := strings.ToLower(claims.Role) // "propietario", "arrendatario", "admin"

	var resultados []map[string]interface{}
	hoy := time.Now()

	switch role {
	case "propietario":
		// Traer parcelas propias activas
		resp, _ := services.Metodo_get("API_CRUD_FINCA", "/v1/Parcela",
			fmt.Sprintf("?query=FkFincaParcela.Id_Usuario:%d,Activo:true", userID))
		var data map[string]interface{}
		json.Unmarshal(resp, &data)
		if parcelas, ok := data["Data"].([]interface{}); ok {
			for _, pi := range parcelas {
				parcelaMap, ok := pi.(map[string]interface{})
				if !ok {
					continue
				}
				idP := int(parcelaMap["Id"].(float64))
				nombreParcela := parcelaMap["NombreParcela"].(string)
				// Obtener nombre de finca
				fincaName := ""
				if fkFP, exists := parcelaMap["FkFincaParcela"].(map[string]interface{}); exists {
					if name, ok2 := fkFP["Nombre"].(string); ok2 {
						fincaName = name
					}

				}
				// Verificar que no exista arrendamiento activo
				rspArr, _ := services.Metodo_get("API_CRUD_CULTIVO", "/v1/Arrendamiento_Parcela",
					fmt.Sprintf("?query=FkParcela.Id:%d,Activo:true", idP))
				var arrData map[string]interface{}
				json.Unmarshal(rspArr, &arrData)
				if arrList, ok := arrData["Data"].([]interface{}); ok && len(arrList) > 0 {
					continue // tiene arrendamiento activo
				}
				// Agregar al resultado con label combinada
				resultados = append(resultados, map[string]interface{}{
					"id":    idP,
					"label": fmt.Sprintf("%s (Finca: %s)", nombreParcela, fincaName),
				})
			}
		}

	case "arrendatario":
		// Traer arrendamientos parciales activos del arrendatario
		resp, _ := services.Metodo_get("API_CRUD_CULTIVO", "/v1/Arrendamiento_Parcela",
			fmt.Sprintf("?query=FkArrendamiento.IdUserUserArrendatario:%d,Activo:true", userID))
		var data map[string]interface{}
		json.Unmarshal(resp, &data)
		if arrPs, ok := data["Data"].([]interface{}); ok {
			for _, ai := range arrPs {
				ap, ok := ai.(map[string]interface{})
				if !ok {
					continue
				}
				// Validar rango de fechas
				if arrObj, exists := ap["FkArrendamiento"].(map[string]interface{}); exists {
					inicio, err1 := parseDateTime(arrObj["FechaInicio"])
					fin, err2 := parseDateTime(arrObj["FechaFin"])
					if err1 != nil || err2 != nil || hoy.Before(inicio) || hoy.After(fin) {
						continue
					}
				}
				// Parcela asociada
				if parcObj, exists := ap["FkParcela"].(map[string]interface{}); exists {
					idP := int(parcObj["Id"].(float64))
					nombreParcela := parcObj["NombreParcela"].(string)
					// Obtener nombre finca desde Parcela
					fincaName := ""
					if fkFP, exists2 := parcObj["FkFincaParcela"].(map[string]interface{}); exists2 {
						if fkFinca, ok2 := fkFP["FkFinca"].(map[string]interface{}); ok2 {
							fincaName, _ = fkFinca["Nombre"].(string)
						}
					}
					resultados = append(resultados, map[string]interface{}{
						"id":    idP,
						"label": fmt.Sprintf("%s (Finca: %s)", nombreParcela, fincaName),
					})
				}
			}
		}

	case "admin":
		// Traer todas las parcelas activas
		resp, _ := services.Metodo_get("API_CRUD_FINCA", "/v1/Parcela", "?query=Activo:true")
		var data map[string]interface{}
		json.Unmarshal(resp, &data)
		if parcelas, ok := data["Data"].([]interface{}); ok {
			for _, pi := range parcelas {
				parcelaMap, ok := pi.(map[string]interface{})
				if !ok {
					continue
				}
				idP := int(parcelaMap["Id"].(float64))
				nombreParcela := parcelaMap["NombreParcela"].(string)
				// Obtener nombre finca
				fincaName := ""
				if fkFP, exists := parcelaMap["FkFincaParcela"].(map[string]interface{}); exists {
					if name, ok2 := fkFP["Nombre"].(string); ok2 {
						fincaName = name
					}

				}
				resultados = append(resultados, map[string]interface{}{
					"id":    idP,
					"label": fmt.Sprintf("%s (Finca: %s)", nombreParcela, fincaName),
				})
			}
		}

	default:
		// Roles no contemplados devuelven vacío
	}

	// 3. Responder
	c.Data["json"] = map[string]interface{}{"data": resultados}
	c.ServeJSON()
}

// parseDateTime intenta parsear fecha de interface{} en formato ISO o time.Time
func parseDateTime(val interface{}) (time.Time, error) {
	switch v := val.(type) {
	case string:
		return time.Parse(time.RFC3339, v)
	case time.Time:
		return v, nil
	default:
		return time.Time{}, fmt.Errorf("tipo de fecha no soportado: %T", v)
	}
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
