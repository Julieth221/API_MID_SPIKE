package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/beego/beego"
	"gorm.io/gorm"
)

// SensorController operations for Sensor
type SensorController struct {
	beego.Controller
	DB *gorm.DB
}
type Sensor struct {
	Id               int       `orm:"column(id_sensor);pk;auto"`
	NombreSensor     string    `orm:"column(nombre_sensor)"`
	FkTipoSensor     string    `orm:"column(fk_tipo_sensor);rel(fk)"`
	Activo           bool      `orm:"column(activo)"`
	FechaInstalacion time.Time `orm:"column(fecha_instalacion);type(timestamp with time zone);auto_now_add"`
}

// URLMapping ...
func (c *SensorController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Sensor
// @Param	body		body 	models.Sensor	true		"body for Sensor content"
// @Success 201 {object} models.Sensor
// @Failure 403 body is empty
// @router / [post]
func (c *SensorController) Post() {
	var sensor Sensor

	// Leer el cuerpo de la solicitud y convertirlo a una estructura Sensor
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &sensor); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Formato JSON inválido"}
		c.ServeJSON()
		return
	}

	// Validar que los campos requeridos no estén vacíos
	if sensor.NombreSensor == "" || sensor.FkTipoSensor == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Los campos 'nombre_sensor' y 'fk_tipo_sensor' son obligatorios"}
		c.ServeJSON()
		return
	}

	// Insertar en la base de datos
	if err := c.DB.Create(&sensor).Error; err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": fmt.Sprintf("No se pudo crear el sensor: %v", err)}
		c.ServeJSON()
		return
	}

	// Responder con el sensor creado
	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = sensor
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Sensor by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Sensor
// @Failure 403 :id is empty
// @router /:id [get]
func (c *SensorController) GetOne() {
	// Obtener el ID del sensor desde la URL
	id := c.Ctx.Input.Param(":id")

	var sensor Sensor
	// Buscar el sensor en la base de datos
	if err := c.DB.First(&sensor, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.Data["json"] = map[string]string{"error": "Sensor no encontrado"}
		} else {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Data["json"] = map[string]string{"error": "Error al obtener el sensor"}
		}
		c.ServeJSON()
		return
	}

	// Si el sensor existe, devolver los datos en formato JSON
	c.Data["json"] = sensor
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Sensor
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Sensor
// @Failure 403
// @router / [get]
func (c *SensorController) GetAll() {
	var sensores []Sensor

	// Obtener parámetros de la URL (filtros)
	nombreSensor := c.GetString("nombre_sensor")
	tipoSensor := c.GetString("fk_tipo_sensor")
	limit, _ := c.GetInt("limit", 10)  // Número de resultados por página (default: 10)
	offset, _ := c.GetInt("offset", 0) // Paginación (default: 0)

	// Construir consulta con filtros opcionales
	query := c.DB.Model(&Sensor{})
	if nombreSensor != "" {
		query = query.Where("nombre_sensor LIKE ?", "%"+nombreSensor+"%")
	}
	if tipoSensor != "" {
		query = query.Where("fk_tipo_sensor = ?", tipoSensor)
	}

	// Obtener los sensores con paginación
	if err := query.Limit(limit).Offset(offset).Find(&sensores).Error; err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Error al obtener los sensores"}
		c.ServeJSON()
		return
	}

	// Responder con la lista de sensores
	c.Data["json"] = sensores
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Sensor
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Sensor	true		"body for Sensor content"
// @Success 200 {object} models.Sensor
// @Failure 403 :id is not int
// @router /:id [put]
func (c *SensorController) Put() {
	// Obtener el ID desde la URL
	id := c.Ctx.Input.Param(":id")
	if id == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Se requiere un ID para actualizar el sensor"}
		c.ServeJSON()
		return
	}

	// Buscar el sensor en la base de datos
	var sensor Sensor
	if err := c.DB.First(&sensor, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.Data["json"] = map[string]string{"error": "Sensor no encontrado"}
		} else {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Data["json"] = map[string]string{"error": fmt.Sprintf("Error al buscar el sensor: %v", err)}
		}
		c.ServeJSON()
		return
	}

	// Leer los nuevos datos desde el cuerpo de la solicitud
	var updatedData Sensor
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &updatedData); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Formato JSON inválido"}
		c.ServeJSON()
		return
	}

	// Actualizar los campos si se proporcionan en la solicitud
	if updatedData.NombreSensor != "" {
		sensor.NombreSensor = updatedData.NombreSensor
	}
	if updatedData.FkTipoSensor != "" {
		sensor.FkTipoSensor = updatedData.FkTipoSensor
	}
	sensor.Activo = updatedData.Activo // Si no se envía, tomará el valor por defecto (false)

	// Guardar los cambios en la base de datos
	if err := c.DB.Save(&sensor).Error; err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": fmt.Sprintf("Error al actualizar el sensor: %v", err)}
		c.ServeJSON()
		return
	}

	// Responder con el sensor actualizado
	c.Data["json"] = sensor
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Sensor
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *SensorController) Delete() {
	// Obtener el ID desde la URL (para eliminar un solo sensor)
	id := c.Ctx.Input.Param(":id")

	// Si hay un ID en la URL, eliminamos solo ese sensor
	if id != "" {
		if err := c.DB.Delete(&Sensor{}, id).Error; err != nil {
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.Data["json"] = map[string]string{"error": fmt.Sprintf("No se pudo eliminar el sensor con ID %s", id)}
			c.ServeJSON()
			return
		}
		c.Data["json"] = map[string]string{"message": "Sensor eliminado exitosamente"}
		c.ServeJSON()
		return
	}

	// Si no hay ID en la URL, intentamos eliminar varios sensores desde el cuerpo de la solicitud
	var ids []int
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ids); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Se esperaba una lista de IDs en el cuerpo"}
		c.ServeJSON()
		return
	}

	// Verificamos si la lista está vacía
	if len(ids) == 0 {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "La lista de IDs está vacía"}
		c.ServeJSON()
		return
	}

	// Eliminamos múltiples sensores
	if err := c.DB.Delete(&Sensor{}, ids).Error; err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": fmt.Sprintf("No se pudieron eliminar los sensores con IDs %s", strings.Trim(fmt.Sprint(ids), "[]"))}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]string{"message": "Sensores eliminados exitosamente"}
	c.ServeJSON()
}
