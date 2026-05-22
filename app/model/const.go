package model

import "github.com/shopspring/decimal"

type ConfKey string
type Fiat string
type CoinId string
type Crypto string
type TradeType string
type MatchMode string
type Network string
type Range struct {
	MinAmount decimal.Decimal
	MaxAmount decimal.Decimal
}
type TradeTypeConf struct {
	Alias        string
	NetworkName  string
	Network      Network
	Crypto       Crypto
	Native       bool
	Contract     string
	Decimal      int32
	AmountRange  Range
	ExplorerFmt  string
	EndpointKey  ConfKey
	AddrCaseSens bool
}

const (
	AdminUsername ConfKey = "admin_username"
	AdminPassword ConfKey = "admin_password"
	AdminSecure   ConfKey = "admin_secure"
	AdminSecret   ConfKey = "admin_secret"
	AdminLoginIP  ConfKey = "admin_login_ip"
	AdminLoginAt  ConfKey = "admin_login_at"

	ApiAuthToken ConfKey = "api_auth_token"
	ApiAppUri    ConfKey = "api_app_uri"

	AtomUSDT ConfKey = "atom_usdt"
	AtomUSDC ConfKey = "atom_usdc"
	AtomEFUN ConfKey = "atom_efun"
	AtomTRX  ConfKey = "atom_trx"
	AtomBNB  ConfKey = "atom_bnb"
	AtomETH  ConfKey = "atom_eth"

	MonitorMinAmount  ConfKey = "monitor_min_amount"
	PaymentMinAmount  ConfKey = "payment_min_amount"
	PaymentMaxAmount  ConfKey = "payment_max_amount"
	PaymentTimeout    ConfKey = "payment_timeout"
	PaymentStaticPath ConfKey = "payment_static_path"
	PaymentMatchMode  ConfKey = "payment_match_mode"

	RpcEndpointPlasma         ConfKey = "rpc_endpoint_plasma"
	RpcEndpointBsc            ConfKey = "rpc_endpoint_bsc"
	RpcEndpointSolana         ConfKey = "rpc_endpoint_solana"
	RpcEndpointXlayer         ConfKey = "rpc_endpoint_xlayer"
	RpcEndpointPolygon        ConfKey = "rpc_endpoint_polygon"
	RpcEndpointArbitrum       ConfKey = "rpc_endpoint_arbitrum"
	RpcEndpointEthereum       ConfKey = "rpc_endpoint_ethereum"
	RpcEndpointBase           ConfKey = "rpc_endpoint_base"
	RpcEndpointAptos          ConfKey = "rpc_endpoint_aptos"
	RpcEndpointTron           ConfKey = "rpc_endpoint_tron"
	RpcEndpointTronGridApiKey ConfKey = "rpc_endpoint_tron_grid_api_key"

	RateSyncCoingeckoApiUrl ConfKey = "rate_sync_coingecko_api_url"
	RateSyncCoingeckoApiKey ConfKey = "rate_sync_coingecko_api_key"
	RateSyncCmcApiUrl       ConfKey = "rate_sync_cmc_api_url"
	RateSyncCmcApiKey       ConfKey = "rate_sync_cmc_api_key"
	RateSyncCmcSlugEfun     ConfKey = "rate_sync_cmc_slug_efun"
	RateSyncInterval        ConfKey = "rate_sync_interval"
	RateSyncHistoryDays     ConfKey = "rate_sync_history_days"

	NotifyMaxRetry     ConfKey = "notify_max_retry"
	BlockHeightMaxDiff ConfKey = "block_height_max_diff"
	BlockOffsetConfirm ConfKey = "block_offset_confirm"

	MqttHost        ConfKey = "mqtt_host"
	MqttPort        ConfKey = "mqtt_port"
	MqttUser        ConfKey = "mqtt_user"
	MqttPass        ConfKey = "mqtt_pass"
	MqttPublishQos  ConfKey = "mqtt_publish_qos"
	MqttTopicPrefix ConfKey = "mqtt_topic_prefix"
	MqttNetworks    ConfKey = "mqtt_networks"

	NotifierParams  ConfKey = "notifier_params"
	NotifierChannel ConfKey = "notifier_channel"

	SystemInstallLock ConfKey = "system_install_lock"
)

const (
	CNY Fiat = "CNY"
	USD Fiat = "USD"
	JPY Fiat = "JPY"
	EUR Fiat = "EUR"
	GBP Fiat = "GBP"
)

const (
	USDT Crypto = "USDT"
	USDC Crypto = "USDC"
	EFUN Crypto = "EFUN"
	TRX  Crypto = "TRX"
	BNB  Crypto = "BNB"
	ETH  Crypto = "ETH"
)

const (
	Classic   MatchMode = "classic"
	HasPrefix MatchMode = "has_prefix"
	RoundOff  MatchMode = "round_off"
)

var usdGeneralRange = Range{
	MinAmount: decimal.NewFromFloat(0.01),
	MaxAmount: decimal.NewFromFloat(1000000),
}

var networkTradesMap = make(map[Network][]TradeType)
var networkEndpointMap = make(map[Network]ConfKey)
var contractTradeMap = make(map[string]TradeType)
var contractDecimalMap = make(map[string]int32)
var tradeAmountRangeMap = make(map[TradeType]Range)
var explorerUrlMap = make(map[TradeType]string)
var cryptoAtomKeys = make(map[Crypto]ConfKey)
