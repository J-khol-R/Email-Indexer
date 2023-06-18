package scripts

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/J-khol-R/Email-Indexer/models"
)

func GenerateNDJSON(arrayEmails []models.Email) (string, error) {
	var r models.Request
	r.Index = "enron_mails"
	r.Records = arrayEmails

	directory := "files"
	fileName := r.Index + ".ndjson"

	rutaAbsoluta, err := filepath.Abs(directory)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(filepath.Join(rutaAbsoluta, fileName))
	if err != nil {
		return "", fmt.Errorf("Error al crear el archivo NDJSON: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	jsonData, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		return "", fmt.Errorf("Error al convertir a JSON: %v", err)
	}

	_, err = writer.Write(jsonData)
	if err != nil {
		return "", fmt.Errorf("Error al escribir en el archivo NDJSON: %v", err)
	}

	err = writer.Flush()
	if err != nil {
		return "", fmt.Errorf("Error al vaciar el búfer en el archivo NDJSON: %v", err)
	}

	return fileName, nil
}
