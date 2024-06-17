package parser

import (
	"billdb/internal/bill"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
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

// ParserSerbia is a parser for variant 1 of the URL.
type ParserSerbia struct {
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
		log.WithField("xpath", xpath).Errorf(
			"Didn't find xpath.")
		return nil, err
	}
	if err != nil {
		log.WithField("xpath", xpath).Errorf(
			"Error querying xpath: %e", err)
		return nil, err
	}
	return resultNode, nil
}

func fetchItems(
	doc *html.Node,
	billId *time.Time,
) ([]*bill.Item, error) {
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
		log.Error(
			"Error finding token string: ", err)
		return nil, err
	}
	// The first element is the full match
	// The second element (index 1) is the first capture group
	token := tokenSubmatches[1]

	postR, err := http.PostForm(
		"https://suf.purs.gov.rs/specifications",
		url.Values{
			"invoiceNumber": {invoceNumber},
			"token":         {token},
		},
	)
	if err != nil {
		log.Error(
			"Error making post request: ", err)
		return nil, err
	}
	defer postR.Body.Close()

	if postR.StatusCode != 200 {
		log.WithField("statusCode", postR.StatusCode).Error(
			"Error fetching items. Status code: ", err)
		return nil, err
	}

	var rJson PostResponseJson
	err = json.NewDecoder(postR.Body).Decode(&rJson)
	if err != nil {
		log.Error(
			"Error decoding json items: ", err)
		return nil, err
	}
	postR.Body.Close()

	if !rJson.Success {
		log.WithField("Success", rJson.Success).Error(
			"Error fetching items: ", err)
		return nil, err
	}

	items := make([]*bill.Item, 0)
	for _, item := range rJson.Items {
		items = append(items, bill.ItemNew(
			*billId,
			item.Name,
			item.Total,
			item.UnitPrice,
			item.Quantity,
		))
	}
	return items, nil
}

func (p *ParserSerbia) Parse(u *url.URL) (*bill.Bill, error) {
	doc, err := htmlquery.LoadURL(u.String())
	if err != nil {
		log.Error("Error loading URL: ", err)
		return nil, err
	}

	nodes := make(map[string]*html.Node)
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

	nodesStrings := make(map[string]string)
	for xpath, node := range nodes {
		nodesStrings[xpath] = htmlquery.InnerText(node)
	}

	cleanedDate := cleanWhiteSpace(nodesStrings[buyDateXpath])
	dateTime, err := DateParse(dateLayout, cleanedDate)
	if err != nil {
		log.WithField("date", nodesStrings[buyDateXpath]).Error(
			"Error parsing date: ", err)
		return nil, err
	}
	priceString := cleanWhiteSpace(cleanPrice(nodesStrings[priceXpath]))
	price, err := strconv.ParseFloat(priceString, 64)
	if err != nil {
		log.WithField("priceString", priceString).Error(
			"Error parsing price: ", err)
		return nil, err
	}

	country, err := bill.StringToCountry("serbia")
	if err != nil {
		log.Error("Error parsing country string: ", err)
		return nil, err
	}
	currency, err := bill.StringToCurrency("rsd")
	if err != nil {
		log.Error("Error parsing currency string: ", err)
		return nil, err
	}
	tag := bill.Tag("")

	timestamp := time.Now()

	items, err := fetchItems(doc, &timestamp)
	if err != nil {
		log.Error("Error fetching items: ", err)
		return nil, err
	}

	billObject := bill.BillNew(
		timestamp,
		nodesStrings[nameXpath],
		*dateTime,
		price,
		currency,
		1.0,
		country,
		items,
		tag,
		u.String(),
		nodesStrings[billXpath],
	)

	return billObject, nil
}
