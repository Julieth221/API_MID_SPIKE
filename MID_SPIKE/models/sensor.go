package models

import "time"

type Sensor struct {
	Id               int       `orm:"column(id_sensor);pk;auto"`
	NombreSensor     string    `orm:"column(nombre_sensor)"`
	FkTipoSensor     string    `orm:"column(fk_tipo_sensor);rel(fk)"`
	Activo           bool      `orm:"column(activo)"`
	FechaInstalacion time.Time `orm:"column(fecha_instalacion);type(timestamp with time zone);auto_now_add"`
}
