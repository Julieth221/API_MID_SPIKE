package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/beego/beego"
	"github.com/beego/beego/v2/server/web"
)

// SensorController operations for Sensor
type SensorController struct {
	beego.Controller
}
type Sensor struct {
	ID     int     `json:"id"`
	Nombre string  `json:"nombre"`
	Tipo   string  `json:"tipo"`
	Valor  float64 `json:"valor"`
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

}

// GetOne ...
// @Title GetOne
// @Description get Sensor by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Sensor
// @Failure 403 :id is empty
// @router /:id [get]
func (c *SensorController) GetOne() {
	apiCrudUrl, _ := web.AppConfig.String("API_CRUD_SPIKE")
	url := "http://" + apiCrudUrl + "/sensores" // Ruta del endpoint en API_CRUD_SPIKE

	// Hacer la petición HTTP
	resp, err := http.Get(url)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "No se pudo conectar con API_CRUD_SPIKE"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Decodificar la respuesta
	var sensores []Sensor
	if err := json.NewDecoder(resp.Body).Decode(&sensores); err != nil {
		c.Data["json"] = map[string]string{"error": "Error al leer la respuesta"}
		c.ServeJSON()
		return
	}

	// Enviar la respuesta de vuelta al cliente
	c.Data["json"] = sensores
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

}

// Delete ...
// @Title Delete
// @Description delete the Sensor
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *SensorController) Delete() {

}
