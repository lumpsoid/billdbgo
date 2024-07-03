package russia

import (
	B "billdb/internal/bill"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ParserRussia struct {
}

func (p *ParserRussia) Parse(qrString string, password string) (*B.Bill, error) {
	qrParams, err := parseQrString(qrString)
	if err != nil {
		return nil, err
	}
	qr := "0" // 0 -> not our type

	token := computeToken(qrParams)

	var requestBody bytes.Buffer
	contentType := createMultiPartForm(&requestBody, qrParams, qr, token)

	req, err := http.NewRequest("POST", "https://proverkacheka.com/api/v1/check/get", &requestBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	req.Header.Set("Host", "proverkacheka.com")
	req.Header.Set("Origin", "https://proverkacheka.com")
	req.Header.Set("Referer", "https://proverkacheka.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:126.0) Gecko/20100101 Firefox/126.0")
	req.Header.Set("DNT", "1")

	req.AddCookie(&http.Cookie{
		Name:  "ENGID",
		Value: "1.1",
	})

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	passwordHash, err := getPasswordHash(password)
	if err != nil {
		return nil, err

	}

	decryptedData, err := decrypt(body, passwordHash)
	if err != nil {
		return nil, err
	}

	var billJson BillJson
	err = json.Unmarshal(decryptedData, &billJson)
	if err != nil {
		return nil, err
	}
	bill, err := billJson.toBill()
	if err != nil {
		return nil, err
	}
	return bill, nil
}
