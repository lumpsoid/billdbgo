package russia

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"unicode/utf8"
)

const testFolder = "../../../test/json/"

type TestData struct {
	Password        string `json:"password"`
	EncryptedString string `json:"encryptedString"`
	QrString        string `json:"qrData"`
}

var testData *TestData

func readTestData() error {
	if testData != nil {
		return nil
	}

	file, err := os.Open(testFolder + "parser_russia.json")
	if err != nil {
		return err
	}
	defer file.Close()

	data := &TestData{}
	err = json.NewDecoder(file).Decode(data)
	if err != nil {
		return err
	}
	testData = data

	return nil
}

func TestDecryption(t *testing.T) {
	err := readTestData()
	if err != nil {
		t.Error("Error reading test data:", err)
		return
	}

	passwordHash, err := getPasswordHash(testData.Password)
	if err != nil {
		t.Error("Error getting password hash:", err)
		return

	}

	encryptedData, err := base64.StdEncoding.DecodeString(testData.EncryptedString)
	if err != nil {
		t.Error("Error decoding base64:", err)
		return
	}

	decryptedData, err := decrypt(encryptedData, passwordHash)
	if err != nil {
		t.Error("Error decrypting:", err)
		return
	}

	if !utf8.ValidString(string(decryptedData)) {
		t.Error("Decrypted data is not utf-8:", string(decryptedData))
	}
}

func TestCreateBoundary(t *testing.T) {
	err := readTestData()
	if err != nil {
		t.Error("Error reading test data:", err)
		return
	}

	qrParams, err := parseQrString(testData.QrString)
	if err != nil {
		t.Error("Error parsing QR string:", err)
		return
	}
	token := computeToken(qrParams)
	tokenRight := "7"

	if token != tokenRight {
		t.Errorf("Expected %s, got %s", tokenRight, token)
	}
}

func TestSendFormatData(t *testing.T) {
	err := readTestData()
	if err != nil {
		t.Error("Error reading test data:", err)
		return
	}

	var requestBody bytes.Buffer

	qrParams, err := parseQrString(testData.QrString)
	if err != nil {
		t.Error("Error parsing QR string:", err)
		return
	}
	qr := "0" // 0 -> not our type

	token := computeToken(qrParams)

	contentType := createMultiPartForm(&requestBody, qrParams, qr, token)

	req, err := http.NewRequest("POST", "https://proverkacheka.com/api/v1/check/get", &requestBody)
	if err != nil {
		t.Errorf("Error while creating the POST request: %v\n", err)
		return
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
		t.Errorf("Error while sending the POST request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status: %v", resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading body: %v\n", err)
		return
	}

	passwordHash, err := getPasswordHash(testData.Password)
	if err != nil {
		t.Error("Error getting password hash:", err)
		return

	}

	decryptedData, err := decrypt(body, passwordHash)
	if err != nil {
		t.Error("Error decrypting:", err)
		return
	}

	var billJson BillJson
	err = json.Unmarshal(decryptedData, &billJson)
	if err != nil {
		t.Error("Error:", err)
		return
	}

	retailPlaceRight := "Магазин Аленка"
	if billJson.Data.Json.RetailPlace != retailPlaceRight {
		t.Errorf(
			"Error getting billJson RetailPlace got %s, expected %s\n",
			billJson.Data.Json.RetailPlace,
			retailPlaceRight,
		)
	}
	itemZeroNameRight := "КОНФ ВЕС Батончики Рот Фронт"
	if billJson.Data.Json.Items[0].Name != itemZeroNameRight {
		t.Errorf(
			"Error getting billJson Items0 Name got %s, expected %s\n",
			billJson.Data.Json.Items[0].Name,
			itemZeroNameRight,
		)
	}
}
