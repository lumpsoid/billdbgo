package parser

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/antchfx/htmlquery"
)

const urlLink = `https://suf.purs.gov.rs:443/v/?vl=A1U2RVVRSDhUVTZFVVFIOFSmvAQAJ7oEAMRvQQMAAAAAAAABi8nEtiwAAACJEXBZdZJy/NmApRiEns0Sgulz4SpsZpL0dvJtAbJh7IOyoE6pEx+1qDfy59VX5fVpHsJwdGLNUg1a0R/y4+mVo85QwP7TNH4N/yzwrv6nrn1/m+rApP1xaGvy8K11wId0HqIuNIWi5XYQa3ah7fJ+LDi2Hyi/o5/SqDCYN58Hz2VnD4uTg+kmhnTSV6YjFtFRykSBoXx7mKh4SEj352l7r076EAtrrJmdqWFYpcY6qYCzxvwXicNpFnZOHrkuvxYqw86ktSB/nvTRvVGNDPkFmCEMe73K6NArhrajz0pPjsHECoT5FcX1ziqxwRPsv4k0ef1leofQ3djA+Wi3/dIrFixHLL7GbFV1l4r8giajLYOxBEdx0px1MIXuyperIu2OEJrjCiK5QpciFq1Payd1vggQnD7ccsbDXfNuG6r9JekuZvF6XGpgGqL+c9duSOpdW0Rrr+SX1RFmHLhOsFeu38HEVvSckjGaXUmC74bflQ0ggCl2fbic3tWUlfKT6gy3NATDpm7/hU/D2ljOJgu87bP6r7evdhLse9fnUn4DLwVioi32xKnOopaEVQZ508DgNEPCVOppgSXM93cHUOA2HGqzgFL+bR+cV4PmPdgeHWvPpyoHb9QPJZwUZcTHm3v17dR/5gbeKeLoMiSsfXsDrYfl9oYdF6Ml+p4pbyouh7T2pV3zexxL8OWcOlfoGJs=`

func TestCleanPrice(t *testing.T) {
	// Implement tests for cleanPrice function
	fmt.Println("Testing cleanPrice function")
	price := "1.000,00"
	result := cleanPrice(price)
	if result != "1000.00" {
		t.Errorf("Expected 1000.00, got %s", result)
	}
	price = "1.000.000,00"
	result = cleanPrice(price)
	if result != "1000000.00" {
		t.Errorf("Expected 1000000.00, got %s", result)
	}
}

func TestQueryNode(t *testing.T) {
	xpath := "//*[@id='invoiceNumberLabel']"
	u, err := url.Parse(urlLink)
	if err != nil {
		t.Errorf("Failed to parse URL: %v", err)
		return
	}
	doc, err := htmlquery.LoadURL(u.String())
	if err != nil {
		t.Errorf("Failed to load URL: %v", err)
		return
	}
	node, err := queryNode(doc, xpath)
	if err != nil {
		t.Errorf("Failed to query node: %v", err)
		return
	}

	if node.Attr[0].Val != "invoiceNumberLabel" {
		t.Errorf("Expected invoiceNumberLabel, got %s", node.Attr[0].Val)
	}
}

func TestDateParse(t *testing.T) {
	fmt.Println("Testing date parsing functionality")
	links := []string{
		urlLink,
		"https://suf.purs.gov.rs/v/?vl=A0VSRVJEUzcyRVJFUkRTNzKHcgQAGm4EACRJjQAAAAAAAAABjzQuU1YAAAAZ719libGxJgZxbunqEioN0tqp7Aj7Y785qKhNyrIL+PjahS4Y4d3PsxoqORItqZA+Ar3WYyVhGZOe54alfec9pfVk3Ey6DaRHE+zcDYzbNCbm5CVQf6HvdcfppYjiwRaBll1GwLz79TZ4Aen72nDFf0nTnif1kqGT0vJdhzidWHGhpyV5bFMo0qN6vKaTJSfUwdV2YLuWXRs49oZ8sE+Lx+51Jn25HnEHt+Q/plHVstRbnWvnadLm4QRmosw4zJ2Ps+YoCbx5NKrIyPMdFClWE9XdB6ehX1V5myvFt19Uuo4IwPTMPPjP0CMB1Nn93dR/pvpXsYgUicHW3y0oBVX3dHA47CLWx+jMwmlIq+M9X7DkqQAZHp/sXdpdiFl13tgI9WHxB5Otjjd1yt4q/BUvo92guhvJtsc84HIPuvA2/Dmp/gYwa1Cs+YyJRRk3k1Iunudi463vjccH3YsHQAzjx1HMcrgadRF1J+0e99WglDlQd3EGS49G1cKGj+azBx4/Zlw6GcH3+uYOYpXmJRRjmMAyElGA1AVbbmgfTCuJgfMLwWaRW/icF3q+awkzdl+dy/imeha1WrT2BwcubcUtVoekkelyxrbWEC0ZVhcOp95vqtS4nBY18nDHE0LXfJtPoHjXiM9PIGzxWxccLPJhRPCvlWMiCf8Zn7mZlMchKrM9AFlqqpIRXd5Vo/bdHi0=",
	}
	expectedTimes := []time.Time{
		time.Date(2023, 11, 13, 18, 39, 54, 0, time.UTC),
		time.Date(2024, 5, 1, 14, 43, 22, 0, time.UTC),
	}

	for index, link := range links {
		u, err := url.Parse(link)
		if err != nil {
			t.Errorf("Failed to parse URL: %v", err)
			return
		}

		parser := &ParserSerbia{}
		billObject, err := parser.Parse(u)
		if err != nil {
			t.Errorf("Error parsing URL: %v", err)
			return
		}
		if billObject == nil {
			t.Errorf("Expected bill object, got nil")
			return
		}

		expectedTime := expectedTimes[index]
		if !billObject.Date.Equal(expectedTime) {
			t.Errorf("Wrong time. Expected %s, got %s", expectedTime, billObject.Date)
			return
		}
	}
}

func TestFetchItems(t *testing.T) {
	// Implement tests for fetchItems function
	fmt.Println("Testing fetchItems function")
	u, err := url.Parse(urlLink)
	if err != nil {
		t.Errorf("Failed to parse URL: %v", err)
		return
	}

	doc, err := htmlquery.LoadURL(u.String())
	if err != nil {
		t.Errorf("Failed to load URL: %v", err)
		return
	}
	timestamp := time.Now()

	items, err := fetchItems(doc, &timestamp)
	if err != nil {
		t.Errorf("Failed to fetch items: %v", err)
		return
	}
	if len(items) == 0 {
		t.Errorf("No items was fetched. Items in silce = %d", len(items))
		return
	}

	itemName := "Masl.ulje ekst.dev.G.Nature 1l/KOM"
	itemPrice := 1699.0
	if items[0].Name != itemName &&
		items[0].Price != itemPrice {
		t.Errorf(
			"Error in items parse. Expected item[0].Name: %s, got %s",
			itemName,
			items[0].Name,
		)
	}
}

func TestBisareItemFetch(t *testing.T) {
	fmt.Println("Testing item fetch")

	links := []string{
		"https://suf.purs.gov.rs/v/?vl=Azg2S1JLM05TODZLUkszTlPvnwIAt58CAHh2GwAAAAAAAAABjys1m5AAAABdr7N9JtSFbKkkPX7j7AGIeLpbCPb4VyWoWwFARDxv7ujU8iQ3gEKOQQvsd33af1dFniwpKZ7KXj6W79d27qgmNEBRIUM0Ng%2Bzajx6BasiwuU8JygItQ%2Ba6Qd04P%2B1S2gpM9xy%2Bu6JRsT4KB7or%2FwmH7yi5%2Fk2E80jlOGqiXZBe%2B%2Bjl3Hj07y0kGGLEG4WaIzny6a9tH4HmTgLeZPCiR7A%2BSJlmYbU8H%2B%2B3JkQCY0tL0Te6dNuNoPXBivDBvUdQPdEWkAb%2BQce0vNH0JjNYO6p6mdMeWBRife3oqz%2Fic6OndEDVEO0B9gbqDLzgokqMLVRH0EW90MI7z5XaTs0%2BJnNrE6uMvZYwQKivzlDP06DVnnDvHqupqq9fOv2KZohhjfCpwyJORksbwgEuMZ%2F%2FQ7WGtGgFxdjKwbL20QjgfAFU0JUBKz3Nsc3rQxD5gJju%2BpRZMNoXef%2F2gV9r2HbGhVzNp%2F3qTUq2L7DsBfriQDEvJcKxRd%2FbhvO%2BG94m3Hr7XYG3x4sNU56soR%2BfOzz140jajKCj71oIDVqBgVGzLHbHThtg2QBWW7Un1bhHAkEkf69ali%2BY1VPVxOaNZTFSrWR1dKOn%2BrIJlvbSCEKZ8fbwgMsJRuwzbn5g5AojxRjTCIiqM3R2utOkEvSy1T7STQkb7Q3nGPuPLLR3t14ZdEJFes1U6UpCr%2Fc5Od8A8xr%2FoU%3D",
		"https://suf.purs.gov.rs/v/?vl=A0VSRVJEUzcyRVJFUkRTNzILbwQAnmoEAJSnMgAAAAAAAAABjxsAwK8AAACXvTQQQRgdPR4zBE%2FrPMUxeEVG%2B3vSzEfhjEhF4mnKVRz1M6uiXQ%2BP0d%2BzznDUprLIxnMVaCKVWkkAAa9q7qv%2BDb1peE9rsNC3s5QxczTT6g%2FanJ39Cgq4cKaCc4O5MFSpg%2Fdi7Xq48dthZKybtR9i%2B41DelnV7lyJVfVGhELSx6S3HMkfX8HCUh5gv5AjiG%2BWp8czS%2F1DNRS8cQtAzVae%2FGrH0INTt5%2BN4ABtmxBJESCW%2FOjMrcVOhR%2FagwzDYk%2FA6G9Eci0z6QTsuZQRUNs2mmN6vrbHEJ792bxd9jFc57Ef2h9ytT6QUL%2BJUHR23WhTmyPdOL8H7pAjeiS8I3y8Caks5n322CALYj7W2uLmyqV0ZcCJJ6JLUbrfS9fitt1bh2Z1R6ggMrTALW2hx%2F%2FkHIkzTli9dPOqp2grd5%2BoA6XXa4LPRZB1D8Xqr5bDMlnzNgWAd6u4roqxM2p2LLdVt50L6f4mC8gm96dG4h%2BWQ16382TFlesgf9GR0%2FT1piUmsiPmOyvZprgFVfB%2BitkIldPu54CpOYDatHesgxTnlQSg6YyUIBxs5UawivesN1%2FqwzoRvNGueclBJfM0UI2FBDe4k98P9lGn758o%2FZ6KPmIqpxJBaFhhSq2a1cGH3%2FwAcWHpsPFWoc%2FGC5gFrJmdyUabIOzX0l8uwZS9mATbcUZ0%2FUIRiJcQPZ0Y9KUL3Mg%3D",
		"https://suf.purs.gov.rs/v/?vl=A0VSRVJEUzcyRVJFUkRTNzKHcgQAGm4EACRJjQAAAAAAAAABjzQuU1YAAAAZ719libGxJgZxbunqEioN0tqp7Aj7Y785qKhNyrIL%2BPjahS4Y4d3PsxoqORItqZA%2BAr3WYyVhGZOe54alfec9pfVk3Ey6DaRHE%2BzcDYzbNCbm5CVQf6HvdcfppYjiwRaBll1GwLz79TZ4Aen72nDFf0nTnif1kqGT0vJdhzidWHGhpyV5bFMo0qN6vKaTJSfUwdV2YLuWXRs49oZ8sE%2BLx%2B51Jn25HnEHt%2BQ%2FplHVstRbnWvnadLm4QRmosw4zJ2Ps%2BYoCbx5NKrIyPMdFClWE9XdB6ehX1V5myvFt19Uuo4IwPTMPPjP0CMB1Nn93dR%2FpvpXsYgUicHW3y0oBVX3dHA47CLWx%2BjMwmlIq%2BM9X7DkqQAZHp%2FsXdpdiFl13tgI9WHxB5Otjjd1yt4q%2FBUvo92guhvJtsc84HIPuvA2%2FDmp%2FgYwa1Cs%2BYyJRRk3k1Iunudi463vjccH3YsHQAzjx1HMcrgadRF1J%2B0e99WglDlQd3EGS49G1cKGj%2BazBx4%2FZlw6GcH3%2BuYOYpXmJRRjmMAyElGA1AVbbmgfTCuJgfMLwWaRW%2FicF3q%2Bawkzdl%2Bdy%2Fimeha1WrT2BwcubcUtVoekkelyxrbWEC0ZVhcOp95vqtS4nBY18nDHE0LXfJtPoHjXiM9PIGzxWxccLPJhRPCvlWMiCf8Zn7mZlMchKrM9AFlqqpIRXd5Vo%2FbdHi0%3D",
	}
	itemsExpected := []string{
		"Cips Chio kecap 90g/KOM",
		"Mleko bez lakt.1,5mm 1l/KOM",
		"Grasak Gustona 400g/KOM",
	}

	for index, link := range links {
		u, err := url.Parse(link)
		if err != nil {
			t.Errorf("Failed to parse URL: %v", err)
			return
		}

		time := time.Now()
		doc, err := htmlquery.LoadURL(u.String())
		if err != nil {
			t.Error("Error loading URL: ", err)
			return
		}

		items, err := fetchItems(doc, &time)
		if err != nil {
			t.Errorf("Failed to fetch items: %v", err)
			return
		}
		if len(items) == 0 {
			t.Errorf("%d fetched items.", len(items))
			return
		}

		if items[0].Name != itemsExpected[index] {
			t.Errorf("Expected %s, got %s", itemsExpected[index], items[0].Name)
		}
	}
}

func TestParserSerbia_Parse(t *testing.T) {
	// Implement tests for parsing variant 1 URL
	fmt.Println("Testing serbian parser")

	parser := &ParserSerbia{}

	u, err := url.Parse(urlLink)
	if err != nil {
		t.Errorf("Failed to parse URL: %v", err)
		return
	}

	// Call the Parse method
	billObject, err := parser.Parse(u)
	if err != nil {
		t.Errorf("Error parsing URL: %v", err)
		return
	}
	if billObject == nil {
		t.Errorf("Expected bill object, got nil")
		return
	}

	billName := "1002298-177 - Maxi"
	if billObject.Name != billName {
		t.Errorf("Expected %s, got %s", billName, billObject.Name)
	}
	if billObject.Id.IsZero() {
		t.Errorf("Expected non-zero bill ID")
	}
	if billObject.Date.IsZero() {
		t.Errorf("Expected non-zero bill date")
	}
	billPrice := 5462.01
	if billObject.Price != billPrice {
		t.Errorf("Expected bill price %f, got %f", billPrice, billObject.Price)
	}

	billDateBuy, err := time.Parse("02.01.2006. 15:04:05", "13.11.2023. 18:39:54")
	if err != nil {
		t.Errorf("Error parsing date: %v", err)
		return
	}
	if billObject.Date != billDateBuy {
		t.Errorf("Expected bill date %v, got %v", billDateBuy, billObject.Date)
	}
}
