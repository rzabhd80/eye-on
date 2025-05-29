package bitpin

import (
	"github.com/rzabhd80/eye-on/internal/database/models"
	"github.com/rzabhd80/eye-on/internal/helpers"
)

type BitpinSymbolRegistry struct{}

func (reg *BitpinSymbolRegistry) RegisterExchangeSymbols(bitpinExchange *models.Exchange) *[]models.TradingPair {
	pairs := []models.TradingPair{
		{ExchangeID: bitpinExchange.ID, Symbol: "BTC_IRT", BaseAsset: "BTC", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.00000001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "BTC_USDT", BaseAsset: "BTC", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(00000001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "ETH_USDT", BaseAsset: "ETH", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.00001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "ETH_IRT", BaseAsset: "ETH", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.00001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "XRP_USDT", BaseAsset: "XRP", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.00001),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "USDT_IRT", BaseAsset: "USDT", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.01), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "SOL_IRT", BaseAsset: "SOL", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "SOL_USDT", BaseAsset: "SOL", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.001),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "BNB_IRT", BaseAsset: "BNB", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "BNB_USDT", BaseAsset: "BNB", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.001),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "USDC_IRT", BaseAsset: "USDC", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.01), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "ADA_IRT", BaseAsset: "ADA", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.01), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "ADA_USDT", BaseAsset: "ADA", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.0001),
			StepSize: helpers.FloatPointer(0.01), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "DOGE_IRT", BaseAsset: "DOGE", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "DOGE_USDT", BaseAsset: "DOGE", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.00001),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "TRX_IRT", BaseAsset: "TRX", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.01), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "TRX_USDT", BaseAsset: "TRX", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.0001),
			StepSize: helpers.FloatPointer(0.01), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "LINK_USDT", BaseAsset: "LINK", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.0001),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "DOT_IRT", BaseAsset: "DOT", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "DOT_USDT", BaseAsset: "DOT", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.0001),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "LTC_IRT", BaseAsset: "LTC", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "LTC_USDT", BaseAsset: "LTC", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.001),
			StepSize: helpers.FloatPointer(0.001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "AVAX_IRT", BaseAsset: "AVAX", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "AVAX_USDT", BaseAsset: "AVAX", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.001),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "UNI_IRT", BaseAsset: "UNI", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "UNI_USDT", BaseAsset: "UNI", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.0001),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "TON_USDT", BaseAsset: "TON", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.0001),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "ATOM_USDT", BaseAsset: "ATOM", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.0001),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "XLM_USDT", BaseAsset: "XLM", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.0001),
			StepSize: helpers.FloatPointer(0.01), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "BCH_IRT", BaseAsset: "BCH", QuoteAsset: "IRT", TickSize: helpers.FloatPointer(1),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
		{ExchangeID: bitpinExchange.ID, Symbol: "BCH_USDT", BaseAsset: "BCH", QuoteAsset: "USDT", TickSize: helpers.FloatPointer(0.01),
			StepSize: helpers.FloatPointer(0.0001), IsActive: true},
	}

	return &pairs
}
