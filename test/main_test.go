package main_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type CollectorResponse struct {
	Rank     int     `json:"Rank"`
	Name     string  `json:"Name"`
	FullName string  `json:"FullName"`
	Price    float32 `json:"Price"`
	Currency string  `json:"Currency"`
}

const (
	fakeAddr     string = "127.0.0.1:12222"
	collectorURL string = "http://localhost:12525"

	ccFakePath  string = "/data/top/totaltoptiervolfull"
	cmcFakePath string = "/v1/cryptocurrency/quotes/latest"

	expectedLimitCC string = "20"
	expectedTSYMCC  string = "USD"
	expectedPageCC  string = "0"

	ccFakeResonseFile string = "testdata/ccresponse.json"

	expectedConvertCMC string = "USD"
	expectedSymbolCMC  string = "BTC,TNCC,ETH,XRP,CRO,USDT,HYN,LINK,HMR,BCH,NYE,AION,CTAG,BSV,PLF,LTC,BNB,ADA,EOS,THX"

	cmcFakeResonseFile string = "testdata/cmcresponse.json"
)

var (
	ccFakeResonse  []byte
	cmcFakeResonse []byte

	expectedData = []CollectorResponse{
		{Rank: 1, Name: "BTC", FullName: "Bitcoin", Price: 9096.642, Currency: "USD"},
		{Rank: 2, Name: "TNCC", FullName: "TNC Coin", Price: 0, Currency: "USD"},
		{Rank: 3, Name: "ETH", FullName: "Ethereum", Price: 226.69994, Currency: "USD"},
		{Rank: 4, Name: "XRP", FullName: "XRP", Price: 0.17488791, Currency: "USD"},
		{Rank: 5, Name: "CRO", FullName: "Crypto.com Chain Token", Price: 0.1230979, Currency: "USD"},
		{Rank: 6, Name: "USDT", FullName: "Tether", Price: 1.0033295, Currency: "USD"},
		{Rank: 7, Name: "HYN", FullName: "Hyperion", Price: 0.51610065, Currency: "USD"},
		{Rank: 8, Name: "LINK", FullName: "Chainlink", Price: 4.7036905, Currency: "USD"},
		{Rank: 9, Name: "HMR", FullName: "Homeros", Price: 0.4304923, Currency: "USD"},
		{Rank: 10, Name: "BCH", FullName: "Bitcoin Cash", Price: 220.49957, Currency: "USD"},
		{Rank: 11, Name: "NYE", FullName: "NewYork Exchange", Price: 9.071555, Currency: "USD"},
		{Rank: 12, Name: "AION", FullName: "Aion", Price: 0.09388478, Currency: "USD"},
		{Rank: 13, Name: "CTAG", FullName: "CTAGtoken", Price: 0, Currency: "USD"},
		{Rank: 14, Name: "BSV", FullName: "Bitcoin SV", Price: 154.67036, Currency: "USD"},
		{Rank: 15, Name: "PLF", FullName: "PlayFuel", Price: 0.26513317, Currency: "USD"},
		{Rank: 16, Name: "LTC", FullName: "Litecoin", Price: 40.899677, Currency: "USD"},
		{Rank: 17, Name: "BNB", FullName: "Binance Coin", Price: 15.411868, Currency: "USD"},
		{Rank: 18, Name: "ADA", FullName: "Cardano", Price: 0.09066484, Currency: "USD"},
		{Rank: 19, Name: "EOS", FullName: "EOS", Price: 2.3326018, Currency: "USD"},
		{Rank: 20, Name: "THX", FullName: "Thorenext", Price: 0.08350019, Currency: "USD"},
	}
)

func ccFakeHandler(t *testing.T, w http.ResponseWriter, req *http.Request) {
	expected := require.New(t)

	q := req.URL.Query()

	expected.Equal(q.Get("limit"), expectedLimitCC)
	expected.Equal(q.Get("tsym"), expectedTSYMCC)
	expected.Equal(q.Get("page"), expectedPageCC)

	w.Header().Set("Content-Type", "application/json")
	bytes.NewBuffer(ccFakeResonse).WriteTo(w)
}

func cmcFakeHandler(t *testing.T, w http.ResponseWriter, req *http.Request) {
	expected := require.New(t)

	q := req.URL.Query()

	expected.Equal(q.Get("symbol"), expectedSymbolCMC)
	expected.Equal(q.Get("convert"), expectedConvertCMC)

	w.Header().Set("Content-Type", "application/json")
	bytes.NewBuffer(cmcFakeResonse).WriteTo(w)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	var err error

	ccFakeResonse, err = ioutil.ReadFile(ccFakeResonseFile)
	panicIfErr(err)

	cmcFakeResonse, err = ioutil.ReadFile(cmcFakeResonseFile)
	panicIfErr(err)

	go http.ListenAndServe(fakeAddr, nil)

	os.Exit(m.Run())
}

func TestAll(t *testing.T) {
	expected := require.New(t)

	http.HandleFunc(ccFakePath, func(w http.ResponseWriter, req *http.Request) {
		ccFakeHandler(t, w, req)
	})

	http.HandleFunc(cmcFakePath, func(w http.ResponseWriter, req *http.Request) {
		cmcFakeHandler(t, w, req)
	})

	request := collectorURL + "?limit=" + expectedLimitCC + "&format=json"

	resp, err := http.Get(request)
	expected.NoError(err)
	defer resp.Body.Close()

	actualData := []CollectorResponse{}

	err = json.NewDecoder(resp.Body).Decode(&actualData)
	expected.NoError(err)

	expected.Equal(expectedData, actualData)
}
