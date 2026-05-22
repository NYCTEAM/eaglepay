package model

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/utils"
	"gorm.io/gorm"
)

type Rate struct {
	Id
	Rate    string  `gorm:"column:rate;type:varchar(32);not null;comment:订单汇率" json:"rate"`
	Fiat    string  `gorm:"column:fiat;type:varchar(16);not null;comment:法币" json:"fiat"`
	Crypto  string  `gorm:"column:crypto;type:varchar(16);not null;comment:加密货币" json:"crypto"`
	RawRate float64 `gorm:"column:raw_rate;type:decimal(18,8);not null;comment:基准汇率" json:"raw_rate"`
	Syntax  string  `gorm:"column:syntax;type:varchar(32);not null;default:'';comment:浮动语法" json:"syntax"`
	AutoTimeAt
}

const efunDefaultRateSyntax = "~0.85"

const coinGeckoBscPlatform = "binance-smart-chain"

func (r *Rate) TableName() string {
	return "bep_rate"
}

func (r *Rate) BeforeCreate(*gorm.DB) error {
	syntax := GetK(ConfKey(fmt.Sprintf("rate_float_%s_%s", r.Crypto, r.Fiat)))
	syntax = fallbackRateSyntax(Crypto(r.Crypto), Fiat(r.Fiat), syntax)
	if syntax == "" {
		return nil
	}

	r.Syntax = syntax
	r.Rate = cast.ToString(ParseFloatRate(syntax, cast.ToFloat64(r.RawRate)))

	return nil
}

func CoingeckoRate() error {
	fiats := make([]string, 0)
	for k := range supportFiat {
		fiats = append(fiats, string(k))
	}

	ids := make([]string, 0)
	tokens := make(map[CoinId]Crypto)
	for token, id := range supportCrypto {
		if id == "" {
			continue
		}
		ids = append(ids, string(id))
		tokens[id] = token
	}

	if len(ids) == 0 {
		return errors.New("CoingeckoRate: no supported ids")
	}

	requestURL := fmt.Sprintf("%s/api/v3/simple/price?ids=%s&vs_currencies=%s", strings.TrimRight(GetC(RateSyncCoingeckoApiUrl), "/"), strings.Join(ids, ","), strings.Join(fiats, ","))
	client := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}

	if apiKey := strings.TrimSpace(GetC(RateSyncCoingeckoApiKey)); apiKey != "" {
		req.Header.Set("x-cg-demo-api-key", apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("CoingeckoRate: " + http.StatusText(resp.StatusCode))
	}

	data := gjson.ParseBytes(body)
	if data.Get("status.error_code").Exists() {
		return errors.New("CoingeckoRate: " + data.Get("status.error_message").String())
	}

	rows := make([]Rate, 0)
	for id, v := range data.Map() {
		token, ok := tokens[CoinId(id)]
		if !ok {
			continue
		}

		for fiat, val := range v.Map() {
			rows = append(rows, Rate{
				Rate:    val.String(),
				Fiat:    strings.ToUpper(fiat),
				Crypto:  string(token),
				RawRate: val.Float(),
			})
		}
	}

	if len(rows) == 0 {
		return errors.New("CoingeckoRate: no data")
	}

	Db.Create(&rows)
	return nil
}

func CmcRate() error {
	apiKey := strings.TrimSpace(GetC(RateSyncCmcApiKey))
	if apiKey == "" {
		return nil
	}

	slug := strings.TrimSpace(GetC(RateSyncCmcSlugEfun))
	if slug == "" {
		return errors.New("CmcRate: missing EFUN slug")
	}

	fiats := make([]string, 0)
	for k := range supportFiat {
		fiats = append(fiats, string(k))
	}

	endpoint := strings.TrimRight(GetC(RateSyncCmcApiUrl), "/")
	params := url.Values{}
	params.Set("slug", slug)
	params.Set("convert", strings.Join(fiats, ","))

	req, err := http.NewRequest(http.MethodGet, endpoint+"/v3/cryptocurrency/quotes/latest?"+params.Encode(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-CMC_PRO_API_KEY", apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("CmcRate: " + http.StatusText(resp.StatusCode))
	}

	data := gjson.ParseBytes(body)
	if errorCode := data.Get("status.error_code"); errorCode.Exists() && errorCode.Int() != 0 {
		return errors.New("CmcRate: " + data.Get("status.error_message").String())
	}

	rows := make([]Rate, 0)
	items := data.Get("data").Array()
	if len(items) == 0 {
		for _, item := range data.Get("data").Map() {
			items = append(items, item)
		}
	}

	for _, item := range items {
		quote := item.Get("quote")
		for fiat, val := range quote.Map() {
			price := val.Get("price")
			if !price.Exists() {
				continue
			}
			rows = append(rows, Rate{
				Rate:    price.String(),
				Fiat:    strings.ToUpper(fiat),
				Crypto:  string(EFUN),
				RawRate: price.Float(),
			})
		}
	}

	if len(rows) == 0 {
		return errors.New("CmcRate: no data")
	}

	Db.Create(&rows)
	return nil
}

func ParseFloatRate(syntax string, rawVal float64) float64 {
	if syntax == "" {
		return rawVal
	}

	if utils.IsNumber(syntax) {
		return cast.ToFloat64(syntax)
	}

	match, err := regexp.MatchString(`^[~+-]\d+(\.\d+)?$`, syntax)
	if !match || err != nil {
		log.Error("rate syntax parse error", err)
		return 0
	}

	act := syntax[0:1]
	raw := decimal.NewFromFloat(rawVal)
	base := decimal.NewFromFloat(cast.ToFloat64(syntax[1:]))
	result := 0.0

	switch act {
	case "~":
		result = raw.Mul(base).InexactFloat64()
	case "+":
		result = raw.Add(base).InexactFloat64()
	case "-":
		result = raw.Sub(base).InexactFloat64()
	}

	return round(result, 4)
}

func round(val float64, precision int) float64 {
	if precision == 0 {
		return math.Round(val)
	}

	p := math.Pow10(precision)
	if precision < 0 {
		return math.Floor(val*p+0.5) * math.Pow10(-precision)
	}

	return math.Floor(val*p+0.5) / p
}

func GetOrderRate(token Crypto, fiat Fiat, syntax string) (decimal.Decimal, error) {
	if token == EFUN {
		if liveRate, err := GetRealtimeTokenRate(token, fiat, syntax); err == nil {
			return liveRate, nil
		} else {
			log.Warn(fmt.Sprintf("GetOrderRate realtime fetch failed for %s %s: %s", token, fiat, err.Error()))
		}
	}

	var r Rate
	Db.Where("crypto = ? and fiat = ?", token, fiat).Order("created_at desc").Limit(1).Find(&r)
	if r.ID == 0 {
		return decimal.Decimal{}, fmt.Errorf("create order failed, latest rate not found: %s %s", token, fiat)
	}

	syntax = fallbackRateSyntax(token, fiat, syntax)
	if syntax == "" {
		return decimal.NewFromString(r.Rate)
	}

	return decimal.NewFromFloat(ParseFloatRate(syntax, r.RawRate)), nil
}

func fallbackRateSyntax(token Crypto, _ Fiat, syntax string) string {
	if syntax != "" {
		return syntax
	}
	if token == EFUN {
		return efunDefaultRateSyntax
	}
	return ""
}

func GetRealtimeTokenRate(token Crypto, fiat Fiat, syntax string) (decimal.Decimal, error) {
	switch token {
	case EFUN:
		return getCoinGeckoContractRate(coinGeckoBscPlatform, conf.EfunBep20, token, fiat, syntax)
	default:
		return decimal.Decimal{}, fmt.Errorf("realtime rate unsupported for token: %s", token)
	}
}

func getCoinGeckoContractRate(platformID, contract string, token Crypto, fiat Fiat, syntax string) (decimal.Decimal, error) {
	baseURL := strings.TrimRight(GetC(RateSyncCoingeckoApiUrl), "/")
	if baseURL == "" {
		baseURL = "https://api.coingecko.com"
	}

	params := url.Values{}
	params.Set("contract_addresses", strings.ToLower(strings.TrimSpace(contract)))
	params.Set("vs_currencies", strings.ToLower(string(fiat)))

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v3/simple/token_price/%s?%s", baseURL, platformID, params.Encode()), nil)
	if err != nil {
		return decimal.Decimal{}, err
	}

	if apiKey := strings.TrimSpace(GetC(RateSyncCoingeckoApiKey)); apiKey != "" {
		req.Header.Set("x-cg-demo-api-key", apiKey)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return decimal.Decimal{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return decimal.Decimal{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return decimal.Decimal{}, errors.New("CoinGecko token_price: " + http.StatusText(resp.StatusCode))
	}

	data := gjson.ParseBytes(body)
	priceNode := data.Get(strings.ToLower(contract) + "." + strings.ToLower(string(fiat)))
	if !priceNode.Exists() {
		return decimal.Decimal{}, errors.New("CoinGecko token_price: no data")
	}

	rawRate := priceNode.Float()
	effectiveRate := ParseFloatRate(fallbackRateSyntax(token, fiat, syntax), rawRate)
	if effectiveRate <= 0 {
		return decimal.Decimal{}, errors.New("CoinGecko token_price: invalid parsed rate")
	}

	return decimal.NewFromFloat(effectiveRate), nil
}
