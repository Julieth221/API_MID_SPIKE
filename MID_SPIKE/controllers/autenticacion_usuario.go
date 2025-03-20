package controllers

import (
	"github.com/astaxie/beego"
)

// Autenticacion_usuario_controllerController operations for Autenticacion_usuario_controller
type Autenticacion_usuarioController struct {
	beego.Controller
}

// URLMapping ...
func (c *Autenticacion_usuarioController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Autenticacion_usuario_controller
// @Param	body		body 	models.Autenticacion_usuario_controller	true		"body for Autenticacion_usuario_controller content"
// @Success 201 {object} models.Autenticacion_usuario_controller
// @Failure 403 body is empty
// @router / [post]
func (c *Autenticacion_usuarioController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Autenticacion_usuario_controller by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Autenticacion_usuario_controller
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Autenticacion_usuarioController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Autenticacion_usuario_controller
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Autenticacion_usuario_controller
// @Failure 403
// @router / [get]
func (c *Autenticacion_usuarioController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Autenticacion_usuario_controller
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Autenticacion_usuario_controller	true		"body for Autenticacion_usuario_controller content"
// @Success 200 {object} models.Autenticacion_usuario_controller
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Autenticacion_usuarioController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Autenticacion_usuario_controller
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Autenticacion_usuarioController) Delete() {

}
