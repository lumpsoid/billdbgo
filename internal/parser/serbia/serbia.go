package parser

import (
	"billdb/internal/bill"
	"billdb/internal/bill/country"
	"billdb/internal/bill/currency"
	"billdb/internal/bill/item"
	"billdb/internal/bill/tag"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

const (
	tokenXpath = "/html/head/script[9]"
	// for items fetch
	invoiceXpath = "//*[@id='invoiceNumberLabel']"
	priceXpath   = "//*[@id='totalAmountLabel']"
	buyDateXpath = "//*[@id='sdcDateTimeLabel']"
	billXpath    = "//*[@id='collapse3']/div/pre"
	nameXpath    = "//*[@id='shopFullNameLabel']"
	tokenRegex   = `viewModel\.Token\('(.*)'\);`
	dateLayout   = "2.1.2006. 15:04:05"
)

// Define a struct to represent your JSON data
type PostResponseJson struct {
	Success bool       `json:"Success"`
	Items   []ItemJson `json:"Items"`
}

type ItemJson struct {
	GTIN          string  `json:"GTIN"`
	Name          string  `json:"Name"`
	Quantity      float64 `json:"Quantity"`
	Total         float64 `json:"Total"`
	UnitPrice     float64 `json:"UnitPrice"`
	Label         string  `json:"Label"`
	LabelRate     float64 `json:"LabelRate"`
	TaxBaseAmount float64 `json:"TaxBaseAmount"`
	VatAmount     float64 `json:"VatAmount"`
}

// Parser is a parser for variant 1 of the URL.
type Parser struct {
}

func (p *Parser) Type() string {
	return "rs"
}

func cleanPrice(s string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(s, ".", ""),
		",",
		".",
	)
}

func cleanWhiteSpace(s string) string {
	return strings.Trim(s, " \r\n\t")
}

func queryNode(doc *html.Node, xpath string) (*html.Node, error) {
	resultNode, err := htmlquery.Query(doc, xpath)
	if resultNode == nil {
		log.WithField("xpath", xpath).Error(
			"Didn't find xpath.")
		return nil, err
	}
	if err != nil {
		log.WithField("xpath", xpath).Error(
			"Error querying xpath")
		return nil, err
	}
	return resultNode, nil
}

func fetchItems(
	doc *html.Node,
	billId *ksuid.KSUID,
	client *http.Client,
) ([]*item.Item, error) {
	invoceNode, err := queryNode(doc, invoiceXpath)
	if err != nil {
		return nil, err
	}
	tokenNode, err := queryNode(doc, tokenXpath)
	if err != nil {
		return nil, err
	}
	invoceNumber := strings.Trim(htmlquery.InnerText(invoceNode), " \r\n\t")
	pattern := regexp.MustCompile(tokenRegex)
	tokenSubmatches := pattern.FindStringSubmatch(htmlquery.InnerText(tokenNode))
	if len(tokenSubmatches) == 0 {
		log.Error("Error finding token string")
		return nil, fmt.Errorf("token not found in document")
	}
	// The first element is the full match
	// The second element (index 1) is the first capture group
	token := tokenSubmatches[1]

	// Prepare form data
	formData := url.Values{
		"invoiceNumber": {invoceNumber},
		"token":         {token},
	}

	// Create POST request
	req, err := http.NewRequest(
		http.MethodPost,
		"https://suf.purs.gov.rs/specifications",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		log.Error("Error creating post request: ", err)
		return nil, err
	}

	// Set headers for form data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.3")

	// Execute request
	postR, err := client.Do(req)
	if err != nil {
		log.Error("Error making post request: ", err)
		return nil, err
	}
	defer postR.Body.Close()

	if postR.StatusCode != 200 {
		log.WithField("statusCode", postR.StatusCode).Error("Error fetching items. Status code: ", postR.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", postR.StatusCode)
	}

	var rJson PostResponseJson
	err = json.NewDecoder(postR.Body).Decode(&rJson)
	if err != nil {
		log.Error("Error decoding json items: ", err)
		return nil, err
	}

	if !rJson.Success {
		log.WithField("Success", rJson.Success).
			WithField("Token", token).
			WithField("invoceNumber", invoceNumber).
			Error("Error fetching items")
		return nil, fmt.Errorf("Error fetching invoce items")
	}

	items := make([]*item.Item, 0)
	for _, itemCurrent := range rJson.Items {
		items = append(items, item.New(
			ksuid.New().String(),
			billId.String(),
			itemCurrent.Name,
			itemCurrent.Total,
			itemCurrent.UnitPrice,
			itemCurrent.Quantity,
		))
	}
	return items, nil
}

func dateParse(dateLayout string, dateString string) (*time.Time, error) {
	dateTime, err := time.Parse(dateLayout, dateString)
	if err != nil {
		log.Error("Error parsing date:", err)
		return nil, err
	}
	return &dateTime, nil
}

func (p *Parser) Parse(u string) (*bill.Bill, error) {
	maxAttempts := 3
	var doc *html.Node
	var billId ksuid.KSUID
	var items []*item.Item
	var nodes map[string]*html.Node
	var nodesStrings map[string]string
	var dateTime *time.Time
	var price float64
	var countryBill country.Country
	var currencyBill currency.Currency

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if attempt > 1 {
			log.WithField("attempt", attempt).Info("Refetching page for new token")
			time.Sleep(1 * time.Second)
		}

		client := &http.Client{
			Timeout: 15 * time.Second,
		}
		req, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			log.WithField("attempt", attempt).Error("creating request: ", err)
			if attempt == maxAttempts {
				return nil, err
			}
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.3")
		req.Header.Set("Referer", u)

		// fetch the page
		resp, err := client.Do(req)
		if err != nil {
			log.WithField("attempt", attempt).Error("request failed: ", err)
			if attempt == maxAttempts {
				return nil, err
			}
			continue
		}

		// check the HTTP status code
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			log.WithField("attempt", attempt).Errorf("unexpected status %d (%s) for %s", resp.StatusCode, resp.Status, u)
			if attempt == maxAttempts {
				return nil, fmt.Errorf("bad response: %d %s", resp.StatusCode, resp.Status)
			}
			continue
		}

		// parse the HTML document
		doc, err = htmlquery.Parse(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.WithField("attempt", attempt).Error("parsing HTML: ", err)
			if attempt == maxAttempts {
				return nil, err
			}
			continue
		}

		// Only parse these on first attempt
		if attempt == 1 {
			nodes = make(map[string]*html.Node)
			for _, nodeXpath := range []string{
				invoiceXpath,
				priceXpath,
				buyDateXpath,
				billXpath,
				nameXpath,
			} {
				node, err := queryNode(doc, nodeXpath)
				if err != nil {
					return nil, err
				}
				nodes[nodeXpath] = node
			}

			nodesStrings = make(map[string]string)
			for xpath, node := range nodes {
				nodesStrings[xpath] = htmlquery.InnerText(node)
			}

			cleanedDate := cleanWhiteSpace(nodesStrings[buyDateXpath])
			dateTime, err = dateParse(dateLayout, cleanedDate)
			if err != nil {
				log.WithField("date", nodesStrings[buyDateXpath]).Error(
					"Error parsing date: ", err)
				return nil, err
			}

			priceString := cleanWhiteSpace(cleanPrice(nodesStrings[priceXpath]))
			price, err = strconv.ParseFloat(priceString, 64)
			if err != nil {
				log.WithField("priceString", priceString).Error(
					"Error parsing price: ", err)
				return nil, err
			}

			countryBill, err = country.Parse("serbia")
			if err != nil {
				log.Error("Error parsing country string: ", err)
				return nil, err
			}

			currencyBill, err = currency.Parse("rsd")
			if err != nil {
				log.Error("Error parsing currency string: ", err)
				return nil, err
			}

			billId = ksuid.New()
		}

		// Try to fetch items
		items, err = fetchItems(doc, &billId, client)
		if err != nil {
			if err.Error() == "Error fetching invoce items" {
				log.WithField("attempt", attempt).Error("Failed to fetch items, will retry")
				if attempt == maxAttempts {
					return nil, fmt.Errorf("failed to fetch items after %d attempts", maxAttempts)
				}
				continue
			}
			// For other errors, return immediately
			log.Error("Error fetching items: ", err)
			return nil, err
		}

		// Success! Break out of retry loop
		break
	}

	billObject := bill.New(
		billId.String(),
		nodesStrings[nameXpath],
		*dateTime,
		price,
		currencyBill,
		//TODO exchange system migrate
		// 1.0,
		countryBill,
		items,
		tag.Empty(),
		u,
		nodesStrings[billXpath],
	)
	return billObject, nil
}
