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
			wg.Add(1)
			go func() {
				defer wg.Done()
				emails, err := readFile(path)
				if err != nil {
					fmt.Print(err)
				}
				for _, email := range emails {
					mutex.Lock()
					arrayEmails = append(arrayEmails, email)
					mutex.Unlock()
				}
				fmt.Printf("\r%s", "archivo leido correctamente:"+path)

			}()
		}

		return nil
	})

	if err != nil {
		return arrayEmails, err
	}

	return arrayEmails, nil
}

func readFile(archivo string) ([]models.Email, error) {
	file, err := os.Open(archivo)
	if err != nil {
		return []models.Email{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	const maxTokenSize = 10 * 1024 * 1024
	buf := make([]byte, maxTokenSize)
	scanner.Buffer(buf, maxTokenSize)

	var emails []models.Email
	var email models.Email

	//hace referencia al nombre del campo actual ej: from: <-- ese seria el campo
	var campoActual string

	//esta condicion hace referencia a cuando empezar a guardar el content del email
	//cuando es x-filename es porque es el cabezal del mensaje
	// y cuando es subject es porque son los submensajes del archivo
	condicion := "X-FileName"

	//sirve para saber cuando empiezan los submensajes
	contador := 0

	for scanner.Scan() {

		linea := scanner.Text()
		if strings.TrimSpace(linea) == "" {
			continue // Ignorar líneas vacías
		}

		//separamos las lineas con la conveccion campo : valor...
		campoValor := strings.SplitN(linea, ":", 2)

		//si no cumple la condicion pasada significa que el archivo todavia no ha llegado al
		// contenido del mensaje
		if campoValor[0] != condicion {

			//si campo valor no vale 2 significa que es una linea adicional del campo anterior
			if len(campoValor) != 2 {
				campoValor[0] = campoActual
				campoValor = append(campoValor, "\n"+linea)
			}
		}

		campo := strings.TrimSpace(campoValor[0])
		valor := strings.TrimSpace(campoValor[1])

		//aqui se guardara el contenido final de cada email
		var mensaje strings.Builder

		switch campo {
		case "Message-ID":
			email.MessageID += valor
			campoActual = campo
		case "Date":
			email.Date += valor
			campoActual = campo
		case "From":
			email.From += valor
			campoActual = campo
		case "Sent":
			email.Sent += valor
			campoActual = campo
		case "To":
			email.To += valor
			campoActual = campo
		case "Subject":
			email.Subject += valor
			campoActual = campo

			//si el contador es mayor a cero significa que ya estamos ubicados
			//en los submensajes del archivos
			if contador > 0 {

				//esta variable sirve para identificar si vamos a comenzar un nuevo email respecto
				// a unas directivas
				newEmail := false

				//aqui se empieza a añadir el "content" del mensaje
				for scanner.Scan() {
					linea := scanner.Text()

					if strings.TrimSpace(linea) == "" {
						continue //omitimos las lineas vacias
					}

					//si el mensaje contiene las siguientes cadenas significa que el mensaje
					//anterior termino y que el nuevo comienza
					if strings.Contains(linea, "-----Original Message-----") ||
						strings.Contains(linea, "---------------------- Forwarded by") ||
						strings.Contains(linea, "___________________") {

						//añado el mensaje a el email pasado antes de añadirlo al array
						email.Content = mensaje.String()

						emails = append(emails, email)
						email = models.Email{}

						//damos la indicacion de que el nuevo email comenzo
						newEmail = true

						//seteamos el campo actual por si el siguiente email no empieza con
						//campo:valor no los agregue al contenido del actual email
						campoActual = ""
						break
					}

					//si la linea no contiene las anteriores directivas añadimos al mensaje
					//el texto que encuentre
					mensaje.WriteString(linea)

				}

				//si nunca se dio la condicion de que habia un nuevo mensaje
				//guardamos el mensaje en el email actual
				if !newEmail {
					email.Content = mensaje.String()
				}
			}
		case "Cc":
			email.Cc += valor
			campoActual = campo
		case "cc":
			email.Cc += valor
			campoActual = campo
		case "Mime-Version":
			email.MimeVersion += valor
			campoActual = campo
		case "Content-Type":
			email.ContentType += valor
			campoActual = campo
		case "Content-Transfer-Encoding":
			email.ContentTransferEncoding += valor
			campoActual = campo
		case "Bcc":
			email.Bcc += valor
			campoActual = campo
		case "X-From":
			email.XFrom += valor
			campoActual = campo
		case "X-To":
			email.XTo += valor
			campoActual = campo
		case "X-cc":
			email.Xcc += valor
			campoActual = campo
		case "X-bcc":
			email.Xbcc += valor
			campoActual = campo
		case "X-Folder":
			email.XFolder += valor
			campoActual = campo
		case "X-Origin":
			email.XOrigin += valor
			campoActual = campo
		case "X-FileName":
			email.XFileName += valor
			campoActual = campo

			newEmail := false
			//empezar a añadir el mensaje
			for scanner.Scan() {
				linea := scanner.Text()

				if strings.TrimSpace(linea) == "" {
					continue
				}

				if strings.Contains(linea, " -----Original Message-----") ||
					strings.Contains(linea, "---------------------- Forwarded by") ||
					strings.Contains(linea, "___________________") {

					//añado el mensaje a el email pasado antes de añadirlo al array
					email.Content = mensaje.String()

					emails = append(emails, email)
					email = models.Email{}
					condicion = "Subject"
					newEmail = true
					contador++
					campoActual = ""
					break
				}

				mensaje.WriteString(linea)

			}
			//si nunca se dio la condicion de que habia un nuevo mensaje
			//guardamos el mensaje en el email actual
			if !newEmail {
				email.Content = mensaje.String()
			}

		}
	}

	if err := scanner.Err(); err != nil {
		return []models.Email{}, err
	}

	emails = append(emails, email)

	return emails, nil
}

// func ReadFile(archivo string) (*models.Email, error) {
// 	file, err := os.Open(archivo)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	const maxTokenSize = 10 * 1024 * 1024
// 	buf := make([]byte, maxTokenSize)
// 	scanner.Buffer(buf, maxTokenSize)

// 	email := &models.Email{}

// 	for scanner.Scan() {
// 		linea := scanner.Text()
// 		if strings.TrimSpace(linea) == "" {
// 			continue
// 		}

// 		campoValor := strings.SplitN(linea, ":", 2)
// 		if len(campoValor) != 2 {
// 			continue
// 		}

// 		campo := strings.TrimSpace(campoValor[0])
// 		valor := strings.TrimSpace(campoValor[1])

// 		var mensaje strings.Builder

// 		switch campo {
// 		case "Message-ID":
// 			email.MessageID = valor
// 		case "Date":
// 			email.Date = valor
// 		case "From":
// 			email.From = valor
// 		case "To":
// 			email.To = valor
// 		case "Subject":
// 			email.Subject = valor
// 		case "Mime-Version":
// 			email.MimeVersion = valor
// 		case "Content-Type":
// 			email.ContentType = valor
// 		case "Content-Transfer-Encoding":
// 			email.ContentTransferEncoding = valor
// 		case "X-From":
// 			email.XFrom = valor
// 		case "X-To":
// 			email.XTo = valor
// 		case "X-cc":
// 			email.Xcc = valor
// 		case "X-bcc":
// 			email.Xbcc = valor
// 		case "X-Folder":
// 			email.XFolder = valor
// 		case "X-Origin":
// 			email.XOrigin = valor
// 		case "X-FileName":
// 			email.XFileName = valor

// 			for scanner.Scan() {
// 				linea := scanner.Text()
// 				if strings.TrimSpace(linea) == "" {
// 					continue
// 				}
// 				mensaje.WriteString(linea)
// 			}
// 			email.Content = mensaje.String()
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return nil, err
// 	}

// 	return email, nil
// }

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
