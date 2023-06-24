package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/J-khol-R/Email-Indexer/models"
)

func GenerateEmails() ([]models.Email, error) {
	var arrayEmails []models.Email
	var mutex sync.Mutex

	nombreArchivo := os.Args[1]

	var wg sync.WaitGroup

	err := filepath.Walk(nombreArchivo, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// email, err := ReadFile(path)
			// if err != nil {
			// 	return err
			// }
			// arrayEmails = append(arrayEmails, *email)
			wg.Add(1)
			go func() {
				defer wg.Done()
				email, err := ReadFile(path)
				if err != nil {
					return
				}
				mutex.Lock()
				arrayEmails = append(arrayEmails, *email)
				mutex.Unlock()
				fmt.Printf("\r%s%%", path)
			}()
		}

		return nil
	})

	if err != nil {
		return arrayEmails, fmt.Errorf("error en readfile: %v", err)
	}

	return arrayEmails, nil
}

func ReadFile(archivo string) (*models.Email, error) {
	file, err := os.Open(archivo)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	const maxTokenSize = 10 * 1024 * 1024
	buf := make([]byte, maxTokenSize)
	scanner.Buffer(buf, maxTokenSize)

	email := &models.Email{}

	for scanner.Scan() {
		linea := scanner.Text()
		if strings.TrimSpace(linea) == "" {
			continue
		}

		campoValor := strings.SplitN(linea, ":", 2)
		if len(campoValor) != 2 {
			continue
		}

		campo := strings.TrimSpace(campoValor[0])
		valor := strings.TrimSpace(campoValor[1])

		var mensaje strings.Builder

		switch campo {
		case "Message-ID":
			email.MessageID = valor
		case "Date":
			email.Date = valor
		case "From":
			email.From = valor
		case "To":
			email.To = valor
		case "Subject":
			email.Subject = valor
		case "Mime-Version":
			email.MimeVersion = valor
		case "Content-Type":
			email.ContentType = valor
		case "Content-Transfer-Encoding":
			email.ContentTransferEncoding = valor
		case "X-From":
			email.XFrom = valor
		case "X-To":
			email.XTo = valor
		case "X-cc":
			email.Xcc = valor
		case "X-bcc":
			email.Xbcc = valor
		case "X-Folder":
			email.XFolder = valor
		case "X-Origin":
			email.XOrigin = valor
		case "X-FileName":
			email.XFileName = valor

			for scanner.Scan() {
				linea := scanner.Text()
				if strings.TrimSpace(linea) == "" {
					continue
				}
				mensaje.WriteString(linea)
			}
			email.Content = mensaje.String()
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return email, nil
}

func GenerateNDJSON() (string, error) {
	arrayEmails, err := GenerateEmails()
	if err != nil {
		return "", fmt.Errorf("error al generar los emails: %v", err)
	}

	fmt.Print("todos los emails procesados")
	fmt.Print("creando ndjdon...")

	var r models.Request

	r.Index = "enron_mails"
	r.Records = arrayEmails

	fileName := r.Index + ".ndjson" // Nombre del archivo NDJSON
	file, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("error al crear el archivo ndjson: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	jsonData, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		return "", fmt.Errorf("error al convertir a json: %v", err)
	}

	_, err = writer.Write(jsonData)
	if err != nil {
		return "", fmt.Errorf("error al escribir en el archivo ndjson: %v", err)
	}

	err = writer.Flush()
	if err != nil {
		return "", fmt.Errorf("error al vaciar el búfer en el archivo ndjson: %v", err)
	}

	fmt.Print("archivo ndjson creado :)")

	return fileName, nil
}
