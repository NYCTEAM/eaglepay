package model

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

const (
	OrderNotifyStateSucc = 1
	OrderNotifyStateFail = 0

	OrderStatusWaiting    = 1
	OrderStatusSuccess    = 2
	OrderStatusExpired    = 3
	OrderStatusCanceled   = 4
	OrderStatusConfirming = 5
	OrderStatusFailed     = 6

	BscBnb      TradeType = "bsc.bnb"
	EthereumEth TradeType = "ethereum.eth"
	TronTrx     TradeType = "tron.trx"

	UsdtTrc20    TradeType = "usdt.trc20"
	UsdcTrc20    TradeType = "usdc.trc20"
	UsdtPolygon  TradeType = "usdt.polygon"
	UsdcPolygon  TradeType = "usdc.polygon"
	UsdtArbitrum TradeType = "usdt.arbitrum"
	UsdcArbitrum TradeType = "usdc.arbitrum"
	UsdtErc20    TradeType = "usdt.erc20"
	UsdcErc20    TradeType = "usdc.erc20"
	UsdtBep20    TradeType = "usdt.bep20"
	UsdcBep20    TradeType = "usdc.bep20"
	EfunBep20    TradeType = "efun.bep20"
	UsdtXlayer   TradeType = "usdt.xlayer"
	UsdcXlayer   TradeType = "usdc.xlayer"
	UsdcBase     TradeType = "usdc.base"
	UsdtSolana   TradeType = "usdt.solana"
	UsdcSolana   TradeType = "usdc.solana"
	UsdtAptos    TradeType = "usdt.aptos"
	UsdcAptos    TradeType = "usdc.aptos"
	UsdtPlasma   TradeType = "usdt.plasma"
)

const (
	OrderApiTypeEpusdt = "epusdt"
	OrderApiTypeEpay   = "epay"
	OrderApiTypeAdmin  = "admin"
)

type Order struct {
	Id
	OrderId       string     `gorm:"column:order_id;type:varchar(128);not null;index;comment:merchant order id" json:"order_id"`
	TradeId       string     `gorm:"column:trade_id;type:varchar(128);not null;uniqueIndex;comment:local trade id" json:"trade_id"`
	TradeType     TradeType  `gorm:"column:trade_type;type:varchar(20);not null;index;comment:trade type" json:"trade_type"`
	Fiat          Fiat       `gorm:"column:fiat;type:varchar(16);not null;index;default:CNY;comment:fiat currency" json:"fiat"`
	Crypto        Crypto     `gorm:"column:crypto;type:varchar(16);not null;index;default:USDT;comment:crypto currency" json:"crypto"`
	CurrencyLimit string     `gorm:"column:currency_limit;type:varchar(255);not null;default:'';comment:currency limit" json:"currency_limit"`
	Rate          string     `gorm:"column:rate;type:varchar(10);not null;comment:exchange rate" json:"rate"`
	Amount        string     `gorm:"column:amount;type:varchar(32);not null;default:0.00;comment:crypto amount" json:"amount"`
	Money         string     `gorm:"column:money;type:varchar(32);not null;default:0.00;comment:fiat amount" json:"money"`
	Address       string     `gorm:"column:address;type:varchar(128);index;not null;comment:payin address" json:"address"`
	FromAddress   string     `gorm:"column:from_address;type:varchar(128);not null;default:'';comment:from address" json:"from_address"`
	AddressLocked bool       `gorm:"column:address_locked;not null;default:false;comment:address locked" json:"address_locked"`
	Status        int        `gorm:"column:status;not null;default:1;index;index:idx_order_notify_retry,priority:1;comment:status" json:"status"`
	Name          string     `gorm:"column:name;type:varchar(64);not null;default:'';comment:item name" json:"name"`
	ApiType       string     `gorm:"column:api_type;type:varchar(20);not null;default:'epusdt';comment:api type" json:"api_type"`
	ReturnUrl     string     `gorm:"column:return_url;type:varchar(255);not null;default:'';comment:return url" json:"return_url"`
	NotifyUrl     string     `gorm:"column:notify_url;type:varchar(255);not null;default:'';comment:notify url" json:"notify_url"`
	NotifyNum     int        `gorm:"column:notify_num;not null;default:0;index:idx_order_notify_retry,priority:3;comment:notify count" json:"notify_num"`
	NotifyState   int        `gorm:"column:notify_state;not null;default:0;index:idx_order_notify_retry,priority:2;comment:notify state" json:"notify_state"`
	RefHash       string     `gorm:"column:ref_hash;type:varchar(128);not null;default:'';index;comment:tx hash" json:"ref_hash"`
	RefBlockNum   int        `gorm:"column:ref_block_num;not null;default:0;comment:block number" json:"ref_block_num"`
	ExpiredAt     time.Time  `gorm:"column:expired_at;not null;comment:expired at" json:"expired_at"`
	ConfirmedAt   *time.Time `gorm:"column:confirmed_at;not null;comment:confirmed at" json:"confirmed_at"`
	AutoTimeAt
}

func (o *Order) SetCanceled() error {
	o.Status = OrderStatusCanceled
	return Db.Save(o).Error
}

func (o *Order) SetExpired() {
	o.Status = OrderStatusExpired
	Db.Save(o)
}

func (o *Order) SetSuccess() {
	o.Status = OrderStatusSuccess
	Db.Save(o)
}

func (o *Order) SetFailed() {
	o.Status = OrderStatusFailed
	Db.Save(o)
}

func (o *Order) MarkConfirming(blockNum int, from, hash string, at time.Time, amount decimal.Decimal) {
	o.FromAddress = from
	o.ConfirmedAt = &at
	o.RefHash = hash
	o.RefBlockNum = blockNum
	o.Status = OrderStatusConfirming
	if o.AddressLocked {
		rate, _ := decimal.NewFromString(o.Rate)
		o.Amount = amount.String()
		o.Money = rate.Mul(amount).String()
	}
	Db.Save(o)
}

func (o *Order) SetNotifyState(state int) error {
	o.NotifyNum += 1
	o.NotifyState = state
	return Db.Save(o).Error
}

func (o *Order) GetStatusLabel() string {
	switch o.Status {
	case OrderStatusExpired:
		return "expired"
	case OrderStatusWaiting:
		return "waiting"
	case OrderStatusCanceled:
		return "canceled"
	default:
		return "success"
	}
}

func (o *Order) GetStatusEmoji() string {
	switch o.Status {
	case OrderStatusExpired:
		return "🔴"
	case OrderStatusWaiting:
		return "🟡"
	case OrderStatusCanceled:
		return "✖️"
	default:
		return "🟢"
	}
}

func (o *Order) GetTxUrl() string {
	return GetTxUrl(o.TradeType, o.RefHash)
}

func (o *Order) TableName() string {
	return "bep_order"
}

func GetTradeOrder(tradeId string) (Order, bool) {
	var order Order
	res := Db.Where("trade_id = ?", tradeId).Limit(1).Find(&order)
	return order, res.RowsAffected > 0
}

func GetOrderByStatus(status int) []Order {
	orders := make([]Order, 0)
	Db.Where("status = ?", status).Find(&orders)
	return orders
}

func GetNotifyFailedTradeOrders() ([]Order, error) {
	var orders []Order
	maxRetry := cast.ToInt(GetC(NotifyMaxRetry))
	if maxRetry <= 0 {
		maxRetry = cast.ToInt(defaultConf[NotifyMaxRetry])
	}

	res := Db.Where("status = ?", OrderStatusSuccess).
		Where("notify_state = ?", OrderNotifyStateFail).
		Where("notify_num <= ?", maxRetry).Find(&orders)

	return orders, res.Error
}

func CalcTradeAmount(addresses []string, rate decimal.Decimal, p OrderParams) (string, string, error) {
	if p.AddressLocked {
		return LockTradeAddress(addresses, p.TradeType)
	}

	var orders []Order
	lock := make(map[string]bool)
	status := []int{OrderStatusConfirming, OrderStatusWaiting}
	Db.Where("status in (?) and trade_type = ?", status, p.TradeType).Find(&orders)
	for _, order := range orders {
		lock[order.Address+order.Amount] = true
	}

	atom, precision := GetAtomicity(p.TradeType)
	if rate.LessThanOrEqual(decimal.Zero) || precision <= 0 {
		return "", "", fmt.Errorf("invalid atomicity config: %v - %v", atom, precision)
	}

	amount := p.Money.DivRound(rate, precision)
	if amount.LessThan(atom) {
		amount = atom
	}

	for i := 0; i <= 100; i++ {
		for _, addr := range addresses {
			key := addr + amount.String()
			if _, ok := lock[key]; ok {
				continue
			}
			return addr, amount.String(), nil
		}
		amount = amount.Add(atom)
	}

	return "", "", errors.New("failed to allocate amount")
}

func LockTradeAddress(addresses []string, t TradeType) (string, string, error) {
	zero := decimal.Zero.String()
	status := []int{OrderStatusConfirming, OrderStatusWaiting}
	for _, addr := range addresses {
		var o Order
		Db.Where("address = ? and status in (?) and trade_type = ? and address_locked = ?", addr, status, t, true).Order("id desc").Limit(1).Find(&o)
		if o.ID == 0 {
			return addr, zero, nil
		}
	}

	return "", zero, errors.New("no available wallet address")
}

func CalcTradeExpiredAt(sec int64) time.Time {
	if sec >= 180 && sec <= 3600 {
		return time.Now().Add(time.Duration(sec) * time.Second)
	}
	return time.Now().Add(time.Duration(cast.ToUint64(GetK(PaymentTimeout))) * time.Second)
}

func GetAtomicity(t TradeType) (decimal.Decimal, int32) {
	confKey, ok := GetTradeAtomKey(t)
	if !ok {
		confKey = AtomUSDT
	}

	atom, _ := decimal.NewFromString(GetK(confKey))
	return atom, cast.ToInt32(math.Abs(float64(atom.Exponent())))
}
