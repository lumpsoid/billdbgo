package qrcode

import (
  "testing"

)

func TestParseImage(t *testing.T) {
  qrString, err := ParseImage("/home/qq/qr-test.png")
  if err != nil {
    t.Error(err)
    return
  }
  t.Errorf("qrString %s", qrString)
}
