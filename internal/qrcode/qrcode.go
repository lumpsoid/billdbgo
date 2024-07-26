package qrcode

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
)

func ParseImage(filePath string) (string, error) {
  _, err := os.Stat(filePath)
  if os.IsNotExist(err) {
    return "", errors.New("file is not exist")
  }

	cmd := exec.Command("zbarimg", filePath)

  var out bytes.Buffer

	cmd.Stdout = &out

  err = cmd.Run()
	if err != nil {
    if "exit status 4" == err.Error() {
      return "", errors.New("qr code was not detected")
    }
    return "", err
	}
  qrString, err := out.ReadString('\n')
  if !strings.HasPrefix(qrString, "QR-Code:") {
    return "", errors.New("qr code was not decoded")
  }

  qrString, found := strings.CutPrefix(qrString, "QR-Code:")
  // just to be sure
  if !found {
    return "", errors.New("string from qr code was not founded")
  }
  qrString = strings.Trim(qrString, "\n")

  return qrString, nil
}
