package russia

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"mime/multipart"
	"strconv"
	"strings"
)

func createMultiPartForm(requestBody *bytes.Buffer, qrParams *QrRus, qr string, token string) string {
	multipartWriter := multipart.NewWriter(requestBody)

	multipartWriter.WriteField("fn", qrParams.Fn)
	multipartWriter.WriteField("fd", qrParams.Fd)
	multipartWriter.WriteField("fp", qrParams.Fp)
	multipartWriter.WriteField("n", qrParams.N)
	multipartWriter.WriteField("s", qrParams.Sum)
	multipartWriter.WriteField("t", qrParams.TimeString())
	multipartWriter.WriteField("qr", qr)
	multipartWriter.WriteField("token", "0."+token)

	multipartWriter.Close()

	contentType := multipartWriter.FormDataContentType()
	return contentType
}

// calculate token for form data part. Security technique for request, implemented by service.
//
// token - number
func computeToken(qrParams *QrRus) string {
	const p = "0" // a.append("qr", h), (b += h), (b = b.toString());
	data := (qrParams.Fn +
		qrParams.Fd +
		qrParams.Fp +
		qrParams.N +
		qrParams.Sum +
		qrParams.TimeString() +
		p)

	var hashString string
	var strD string

	for d := 0; d < 1000; d++ {
		strD = strconv.Itoa(d)
		dataHash := []byte(data + strD)

		hash := md5.Sum(dataHash)

		hashString = hex.EncodeToString(hash[:])

		split := strings.Split(hashString, "0")

		if len(split)-1 > 4 {
			break
		}
	}
	return strD
}
