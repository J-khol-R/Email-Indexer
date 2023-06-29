package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"time"
)

func main() {

	profileFile, err := os.Create("profile.pprof")
	if err != nil {
		log.Fatalf("Error al crear el archivo de perfil: %v", err)
	}
	defer profileFile.Close()

	// Inicia el perfil de CPU
	if err := pprof.StartCPUProfile(profileFile); err != nil {
		log.Fatalf("Error al iniciar el perfil de CPU: %v", err)
	}
	defer pprof.StopCPUProfile()

	// Preparar los datos para enviar
	filePath, err := GenerateNDJSON() // Ruta al archivo ndjson
	if err != nil {
		fmt.Println("Error al crear el archivo ndjson:", err)
		return
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return
	}

	fmt.Print("se leyo el archivo :)")

	// Crear una solicitud HTTP
	url := "http://localhost:4080/api/_bulkv2"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		fmt.Println("Error al crear la solicitud:", err)
		return
	}

	// Establecer la autenticación
	req.SetBasicAuth("admin", "Complexpass#123")

	// Establecer encabezados
	req.Header.Set("Content-Type", "application/octet-stream")

	// Realizar la petición
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error al hacer la petición:", err)
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error al leer la respuesta:", err)
		return
	}

	fmt.Print("emails indexados correctamente :)")

	// Imprimir la respuesta
	fmt.Println("Respuesta:", string(respBody))

	time.Sleep(30 * time.Second)

	memProfileFile, err := os.Create("memprofile.pprof")
	if err != nil {
		log.Fatalf("Error al crear el archivo de perfil de memoria: %v", err)
	}
	defer memProfileFile.Close()

	// Captura el perfil de memoria
	if err := pprof.WriteHeapProfile(memProfileFile); err != nil {
		log.Fatalf("Error al escribir el perfil de memoria: %v", err)
	}
}
