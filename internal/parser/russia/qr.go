package russia

import (
	"fmt"
	"strings"
	"time"
)

type QrRus struct {
	Fn   string
	Fd   string
	Fp   string
	N    string
	Sum  string
	Time time.Time
}

func (qrRus *QrRus) String() string {
	return fmt.Sprintf("Fn: %s, Fd: %s, Fp: %s, N: %s, Sum: %s, Time: %s",
		qrRus.Fn, qrRus.Fd, qrRus.Fp, qrRus.N, qrRus.Sum, qrRus.Time)
}

func (qrRus *QrRus) TimeString() string {
	return qrRus.Time.Format("02.01.2006") + " " + qrRus.Time.Format("15:04")
}

func setParameter(qrRus *QrRus, key string, value string) error {
	switch key {
	case "fn":
		qrRus.Fn = value
	case "i":
		qrRus.Fd = value
	case "fp":
		qrRus.Fp = value
	case "n":
		qrRus.N = value
	case "s":
		sumValue := strings.Split(value, ".")
		qrRus.Sum = sumValue[0]
	case "t":
		layout := "20060102T1504"
		t, err := time.Parse(layout, value)
		if err != nil {
			return fmt.Errorf("error parsing time: %v", err)
		}
		qrRus.Time = t
	default:
		return fmt.Errorf("unknown key: %s", key)
	}
	return nil
}

func parseQrString(qrString string) (*QrRus, error) {
	parametersList := strings.Split(qrString, "&")
	if len(parametersList) < 6 {
		return nil, fmt.Errorf("invalid QR string")
	}

	qrRus := &QrRus{}
	for _, parameterKeyValue := range parametersList {
		key, value, haveSep := strings.Cut(parameterKeyValue, "=")
		if !haveSep {
			continue
		}
		err := setParameter(qrRus, key, value)
		if err != nil {
			return nil, err
		}
	}
	return qrRus, nil
}
