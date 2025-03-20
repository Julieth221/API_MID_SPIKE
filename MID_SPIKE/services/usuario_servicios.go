package services

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
)

func Metodo_post(nombre_servicio string, endpoint string, data []byte) ([]byte, error) {
	baseURL := beego.AppConfig.String(nombre_servicio)
	if baseURL == "" {
		return nil, fmt.Errorf("no se encontró la configuración para %s", nombre_servicio)
	}

	// Asegurar que la URL tiene "http://"
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}

	url := baseURL + endpoint
	fmt.Println("URL construida:", url)

	// Enviar la solicitud HTTP POST
	response, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error en POST:", err)
		return nil, fmt.Errorf("error en POST a %s: %v", url, err)
	}
	defer response.Body.Close()

	// Verificar si la API respondió con un código de error
	// if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
	// 	body, _ := ioutil.ReadAll(response.Body)
	// 	fmt.Println("Error en la API:", response.StatusCode, "-", string(body))
	// 	return nil, fmt.Errorf("error en la API: %d - %s", response.StatusCode, string(body))
	// }

	// Leer la respuesta
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error al leer la respuesta:", err)
		return nil, fmt.Errorf("error al leer la respuesta: %v", err)
	}

	fmt.Println("Respuesta de la API:", string(body))

	return body, nil
}

func Metodo_get(nombre_servicio, endpoint, parametro string) ([]byte, error) {
	url := beego.AppConfig.String(nombre_servicio) + parametro
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body, nil
}

func Metodo_put(nombre_servicio string, endpoint string, id string, data []byte) ([]byte, error) {

	// Obtener la URL base desde la configuracion de Beego
	baseURL := beego.AppConfig.String(nombre_servicio)

	// Construir la URL final con el ID
	url := fmt.Sprintf("%s/%s", baseURL, id)

	// Crear la solicitud PUT
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Leer la respuesta
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body, nil

}
