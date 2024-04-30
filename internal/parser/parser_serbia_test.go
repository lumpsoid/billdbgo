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
