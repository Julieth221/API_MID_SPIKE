package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	auth_JWT "github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/auth_jwt"
	"github.com/sena_2824182/API_MID_SPIKE/MID_SPIKE/services"
)

// Gestion_insumo_cultivoController operations for Gestion_insumo_cultivo
type Gestion_insumo_cultivoController struct {
	beego.Controller
}

// URLMapping ...
func (c *Gestion_insumo_cultivoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Gestion_insumo_cultivo
// @Param	body		body 	models.Gestion_insumo_cultivo	true		"body for Gestion_insumo_cultivo content"
// @Success 201 {object} models.Gestion_insumo_cultivo
// @Failure 403 body is empty
// @router / [post]
func (c *Gestion_insumo_cultivoController) Post() {
	fmt.Println("Registrar insumo")
	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al procesar JSON"}
		c.ServeJSON()
		return
	}

	// Validar token
	token := c.Ctx.Input.Header("Authorization")
	if token == "" {
		c.Data["json"] = map[string]interface{}{"error": "Token faltante"}
		c.ServeJSON()
		return
	}
	claims, err := auth_JWT.ValidarJWT(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Token inválido"}
		c.ServeJSON()
		return
	}
	userID := claims.UserID

	// Validación obligatoria mínima
	if _, ok := body["FkRegistroCultivo"]; !ok {
		c.Data["json"] = map[string]interface{}{"error": "Falta FkRegistroCultivo"}
		c.ServeJSON()
		return
	}

	// Procesar campos relacionales condicionales
	fkTipoInsumoId, err := obtenerOFkCrearEntidad("tipo_insumo", "nombre", "", body, "FkTipoInsumo", userID)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error procesando TipoInsumo", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	fkCategoriaInsumoId, err := obtenerOFkCrearEntidad("Categoria_insumo", "nombre", "", body, "FkCategoriaInsumo", userID)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error procesando CategoriaInsumo", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	fkMetodoAplicacionId, err := obtenerOFkCrearEntidad("metodo_aplicacion_insumo", "nombre", "", body, "FkMetodoAplicacionInsumo", userID)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error procesando MetodoAplicacionInsumo", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	fechaRaw, ok := body["FechaAplicacion"].(string)
	if ok {
		t, err := time.Parse("2006-01-02", fechaRaw)
		if err == nil {
			body["FechaAplicacion"] = t.Format(time.RFC3339) // "2025-06-01T00:00:00Z"
		}
	}

	// Construcción del insumo
	insumo := map[string]interface{}{
		"FkRegistroCultivo":        map[string]interface{}{"Id": int(body["FkRegistroCultivo"].(float64))},
		"FkTipoInsumo":             map[string]interface{}{"Id": fkTipoInsumoId},
		"FkCategoriaInsumo":        map[string]interface{}{"Id": fkCategoriaInsumoId},
		"FkMetodoAplicacionInsumo": map[string]interface{}{"Id": fkMetodoAplicacionId},
		"NombreInsumo":             body["NombreInsumo"],
		"FechaAplicacion":          body["FechaAplicacion"],
		"CantidadAplicada":         body["CantidadAplicada"],
		"Observaciones":            body["Observaciones"],
		"Activo":                   true,
	}

	// POST al CRUD
	insumoJSON, _ := json.Marshal(insumo)
	fmt.Println("Este es el body que se envia a Insumo: ", string(insumoJSON))
	resp, err := services.Metodo_post("API_CRUD_CULTIVO", "/v1/Insumo", insumoJSON)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error al registrar insumo", "detalle": err.Error()}
		c.ServeJSON()
		return
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(resp, &parsed); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Error leyendo respuesta"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = parsed
	c.ServeJSON()

}

// Función reutilizable para obtener el ID de una entidad o crearla si es nueva
func obtenerOFkCrearEntidad(entidad string, campo1 string, campo2 string, body map[string]interface{}, clave string, userID int) (int, error) {
	valorRaw, ok := body[clave]
	if !ok || valorRaw == nil || valorRaw == "" {
		// Intentar crear con campo NuevoX si existe
		nuevoCampo := "Nuevo" + clave[2:] // ejemplo: FkMetodoAplicacionInsumo -> NuevoMetodoAplicacionInsumo
		if nuevoValor, ok := body[nuevoCampo]; ok && nuevoValor != "" {
			nuevo := map[string]interface{}{
				"Nombre":            nuevoValor.(string),
				"Activo":            true,
				"FechaCreacion":     time.Now().Format(time.RFC3339),
				"FechaModificacion": time.Now().Format(time.RFC3339),
			}

			// Agregar el campo adicional de usuario según entidad
			switch entidad {
			case "metodo_aplicacion_insumo":
				nuevo["FkUsuario"] = map[string]interface{}{"Id": userID}
			case "CategoriaInsumo", "TipoInsumo":
				nuevo["Id_Usuario"] = userID
			}

			payload, _ := json.Marshal(nuevo)
			resp, err := services.Metodo_post("API_CRUD_CULTIVO", "/v1/"+entidad, payload)
			if err != nil {
				return 0, fmt.Errorf("error al registrar %s: %v", entidad, err)
			}

			var parsed map[string]interface{}
			if err := json.Unmarshal(resp, &parsed); err != nil {
				return 0, fmt.Errorf("error al parsear respuesta de %s", entidad)
			}

			dataMap, ok := parsed["Data"].(map[string]interface{})
			if !ok {
				return 0, fmt.Errorf("la respuesta de %s no contiene un objeto válido en 'Data'", entidad)
			}
			return int(dataMap["Id"].(float64)), nil

		}

		return 0, fmt.Errorf("Campo %s faltante y no se proporcionó %s", clave, nuevoCampo)
	}

	switch v := valorRaw.(type) {
	case map[string]interface{}:
		// Si viene un objeto, es porque quiere crear uno nuevo
		valor := v["Nombre"].(string)
		nuevo := map[string]interface{}{
			"Activo":            true,
			"FechaCreacion":     time.Now().Format(time.RFC3339),
			"FechaModificacion": time.Now().Format(time.RFC3339),
		}

		switch entidad {
		case "TipoInsumo":
			nuevo["TipoInsumo"] = v["TipoInsumo"]
			nuevo["Nombre"] = valor
			nuevo["Id_Usuario"] = userID
		case "metodo_aplicacion_insumo":
			nuevo["Nombre"] = valor
			nuevo["FkUsuario"] = map[string]interface{}{"Id": userID}
		case "CategoriaInsumo":
			nuevo["Nombre"] = valor
			nuevo["Id_Usuario"] = userID
		}

		payload, _ := json.Marshal(nuevo)
		resp, err := services.Metodo_post("API_CRUD_CULTIVO", "/v1/"+entidad, payload)
		if err != nil {
			return 0, fmt.Errorf("error al registrar %s: %v", entidad, err)
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal(resp, &parsed); err != nil {
			return 0, fmt.Errorf("error al parsear respuesta de %s", entidad)
		}

		dataMap, ok := parsed["Data"].(map[string]interface{})
		if !ok {
			return 0, fmt.Errorf("la respuesta de %s no contiene un objeto válido en 'Data'", entidad)
		}
		return int(dataMap["Id"].(float64)), nil

	case float64:
		// Si es un número, usar el ID directamente
		return int(v), nil
	default:
		return 0, fmt.Errorf("formato no válido para %s", clave)
	}
}

// GetOne ...
// @Title GetOne
// @Description get Gestion_insumo_cultivo by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Gestion_insumo_cultivo
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Gestion_insumo_cultivoController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Gestion_insumo_cultivo
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Gestion_insumo_cultivo
// @Failure 403
// @router / [get]
func (c *Gestion_insumo_cultivoController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Gestion_insumo_cultivo
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Gestion_insumo_cultivo	true		"body for Gestion_insumo_cultivo content"
// @Success 200 {object} models.Gestion_insumo_cultivo
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Gestion_insumo_cultivoController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Gestion_insumo_cultivo
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Gestion_insumo_cultivoController) Delete() {

}
