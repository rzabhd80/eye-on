package nobitex

import (
	"github.com/rzabhd80/eye-on/internal/database/models"
	"github.com/rzabhd80/eye-on/internal/helpers"
)

type NobitexSymbolRegistry struct{}

func (reg *NobitexSymbolRegistry) RegisterExchangeSymbols(bitpinExchange *models.Exchange) *[]models.TradingPair {
	pairs := []models.TradingPair{
		{ExchangeID: bitpinExchange.ID, Symbol: "BTCIRT", BaseAsset: "BTC", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.00000001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "ETHIRT", BaseAsset: "ETH", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.00000001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "USDTIRT", BaseAsset: "USDT", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "BNBIRT", BaseAsset: "BNB", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "USDCIRT", BaseAsset: "USDC", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "BTCUSDT", BaseAsset: "BTC", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.00000001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "ETHUSDT", BaseAsset: "ETH", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.00000001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "LTCUSDT", BaseAsset: "LTC", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "XRPUSDT", BaseAsset: "XRP", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "BCHUSDT", BaseAsset: "BCH", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "BNBUSDT", BaseAsset: "BNB", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "EOSUSDT", BaseAsset: "EOS", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "XLMUSDT", BaseAsset: "XLM", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.01), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "ETCUSDT", BaseAsset: "ETC", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "TRXUSDT", BaseAsset: "TRX", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.01), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "DOGEUSDT", BaseAsset: "DOGE", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "UNIUSDT", BaseAsset: "UNI", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "DAIUSDT", BaseAsset: "DAI", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "LINKUSDT", BaseAsset: "LINK", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "DOTUSDT", BaseAsset: "DOT", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
	}
	return &pairs
}
