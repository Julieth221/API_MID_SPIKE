package services

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/astaxie/beego"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
	"gopkg.in/gomail.v2"
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

// GenerarToken crea un token de 5 dígitos aleatorios y lo hashea
func GenerarToken() (string, string, error) {
	token := fmt.Sprintf("%05d", 10000+rand.Intn(90000)) // Token de 5 dígitos

	// Hashear el token antes de guardarlo
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return token, string(hashedToken), nil
}

// VerificarToken compara el token ingresado con el hash almacenado
func VerificarToken(tokenIngresado string, tokenGuardado string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(tokenGuardado), []byte(tokenIngresado))
	return err == nil
}

// EnviarCorreo envía un token de recuperación al usuario
func EnviarCorreo(destinatario string, token string) error {

	// Obtener credenciales del .env
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpSender := os.Getenv("SMTP_SENDER")

	// Convertir puerto a entero
	port := 587 // Valor por defecto
	fmt.Sscanf(smtpPort, "%d", &port)

	// Configurar mensaje
	mensaje := gomail.NewMessage()
	mensaje.SetHeader("From", smtpSender)
	mensaje.SetHeader("To", destinatario)
	mensaje.SetHeader("Subject", "Recuperación de contraseña")
	mensaje.SetBody("text/plain", fmt.Sprintf("Tu código de recuperación es: %s", token))

	// Configurar servidor SMTP
	dialer := gomail.NewDialer(smtpHost, port, smtpUser, smtpPass)

	// Enviar correo
	err := dialer.DialAndSend(mensaje)
	if err != nil {
		log.Printf("Error enviando el correo: %v", err)
		return err
	}

	fmt.Println("Correo enviado correctamente a", destinatario)
	return nil
}
