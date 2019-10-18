package main

import (
	"io"
	"strconv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"net/http"
	"net/smtp"
	"net"
    "net/mail"
	"crypto/tls"
	"github.com/magiconair/properties"
	"strings"
	googlesheet "github.com/webserg/alertEngine/readGoogleSheet"
)

//MarketData get from https://api.worldtradingdata.com/api/v1/stock?symbol=AAPL,MSFT,HSBA.L&api_token=BvKoM6i9WSlXECN46TyFnRK7VLSR5H9C19MKB0ZtSB4b4fxgfboeuOqqLCvD
type MarketData struct {
	SymbolsRequested int `json:"symbols_requested"`
	SymbolsReturned  int `json:"symbols_returned"`
	Data             []struct {
		Symbol             string `json:"symbol"`
		Name               string `json:"name"`
		Currency           string `json:"currency"`
		Price              string `json:"price"`
		PriceOpen          string `json:"price_open"`
		DayHigh            string `json:"day_high"`
		DayLow             string `json:"day_low"`
		Five2WeekHigh      string `json:"52_week_high"`
		Five2WeekLow       string `json:"52_week_low"`
		DayChange          string `json:"day_change"`
		ChangePct          string `json:"change_pct"`
		CloseYesterday     string `json:"close_yesterday"`
		MarketCap          string `json:"market_cap"`
		Volume             string `json:"volume"`
		VolumeAvg          string `json:"volume_avg"`
		Shares             string `json:"shares"`
		StockExchangeLong  string `json:"stock_exchange_long"`
		StockExchangeShort string `json:"stock_exchange_short"`
		Timezone           string `json:"timezone"`
		TimezoneName       string `json:"timezone_name"`
		GmtOffset          string `json:"gmt_offset"`
		LastTradeTime      string `json:"last_trade_time"`
	} `json:"data"`
}

func readFile(filaName string) ([]byte, error) {
	// Open our jsonFile
	// jsonFile, err := os.Open("/home/webserg/data/test_accounts_291218/data/data/accounts_1.json")
	jsonFile, err := os.Open(filaName)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// read our opened xmlFile as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)
	return byteValue, err
}

/*ReadData read json data to array*/
func readMarketData(byteValue []byte) (*MarketData, error) {

	// we initialize our Users array
	var marketData MarketData

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err := json.Unmarshal(byteValue, &marketData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(marketData.Data))

	return &marketData, nil
}

func sendEmail(mailAuth string, subj string, body string){
	from := mail.Address{"", "webserg@yahoo.com"}
    to   := mail.Address{"", "webserg@gmail.com"}

    // Setup headers
    headers := make(map[string]string)
    headers["From"] = from.String()
    headers["To"] = to.String()
    headers["Subject"] = subj
    // Setup message
    message := ""
    for k,v := range headers {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + body
    // Connect to the SMTP Server
    servername := "smtp.mail.yahoo.com:465"
    host, _, _ := net.SplitHostPort(servername)
    auth := smtp.PlainAuth("","webserg@yahoo.com", mailAuth, host)
    // TLS config
    tlsconfig := &tls.Config {
        InsecureSkipVerify: true,
        ServerName: host,
    }
    // Here is the key, you need to call tls.Dial instead of smtp.Dial
    // for smtp servers running on 465 that require an ssl connection
    // from the very beginning (no starttls)
    conn, err := tls.Dial("tcp", servername, tlsconfig)
	check(err)
    c, err := smtp.NewClient(conn, host)
    check(err)
    // Auth
    if err = c.Auth(auth); err != nil {
        log.Panic(err)
    }
    // To && From
    if err = c.Mail(from.Address); err != nil {
        log.Panic(err)
    }
    if err = c.Rcpt(to.Address); err != nil {
        log.Panic(err)
    }
    // Data
    w, err := c.Data()
	check(err)
    _, err = w.Write([]byte(message))
	check(err)
    err = w.Close()
	check(err)
	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}
}


func init() {
	fmt.Println("init")	
}

func check(e error) {
	if e != nil {
		log.Panic(e)
		panic(e)
	}
}

func keysString(m map[string]float64) string {
    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return strings.Join(keys, ",")
}

func main() {
	f, err := os.OpenFile(".\\alertEngine.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
    	log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(io.MultiWriter(f, os.Stdout))
	log.Println("This is a test log entry")
	log.Println("start")
	log.Println("read properties")
	props := properties.MustLoadFile(".\\config.properties", properties.UTF8)
	log.Println(props)

	log.Println("read target data")
	targetData,err := googlesheet.ReadTargetData()
	check(err)

	tickers:=keysString(targetData)
	log.Println("tickers=" + tickers)

	log.Println("read market data")
	marketURL := props.MustGetString("marketUrl")
	marketToken := props.MustGetString("marketToken")
	url := marketURL + tickers + marketToken
	resp, err := http.Get(url)
	check(err)
	byteValue, err := ioutil.ReadAll(resp.Body)
	// byteValue, err := readFile("C:/Users/webse/go/src/github.com/webserg/alertEngine/data.json")
	check(err)
	marketData, err := readMarketData(byteValue)
	check(err)
	log.Println(marketData)

	for _, s := range marketData.Data {
		log.Println(s.Symbol + " = " + s.Price)
		price, err := strconv.ParseFloat(s.Price,64)
		log.Println(price)
		check(err)
		targetPrice, exists := targetData[s.Symbol]
		log.Println(targetPrice)
		if(exists && targetPrice >= price){
			log.Println("alert!!! " + s.Symbol + " = " + s.Price)
			body := fmt.Sprintf("target = %.2f ; price = %.2f ", targetPrice, price)
			log.Println(body)
			mailAuth := props.MustGetString("mailAuth")
			sendEmail(mailAuth, "stock alert " + s.Symbol, body)
		}
	}
}
