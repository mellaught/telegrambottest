package vocab

var Translates = map[string]map[string]string{
	"Language": {
		"en": "Language",
		"ru": "Язык",
	},
	"english": {
		"en": "english",
		"ru": "английский",
	},
	"russian": {
		"en": "russian",
		"ru": "русский",
	},
	"Installed": {
		"en": "Installed",
		"ru": "Установлен",
	},
	"Settings": {
		"en": "🔧 Settings",
		"ru": "🔧 Настройки",
	},
	"Select": {
		"en": "Select, please, what you want:)",
		"ru": "Выберете, пожалуйста, что вы хотите:)",
	},
	"Hello": {
		"en": "Hello, i'm an exchange bot BIP/BTC or BTC/BIP",
		"ru": "Привет, я бот для обмена BIP/BTC или BTC/BIP",
	},
	"Now": {
		"en": "📈 Now currency BIP/USD %f $",
		"ru": "📈 Сейчас курс BIP/USD %f $",
	},
	"Send minter": {
		"en": "Send me your Minter address.",
		"ru": "Отправь мне свой minter адрес.",
	},
	"Send BTC": {
		"en": "Send me your Bitcoin address.",
		"ru": "Отправь мне свой биткоин адрес.",
	},
	"Minter deposit and tag": {
		"en": "Your minter deposit address: %s\nTag: %s",
		"ru": "Твой minter адрес для депозита: %s\nTag: %s",
	},
	"BTC deposit": {
		"en": "Your bitcoin deposit address %s",
		"ru": "Твой bitcoin адрес для депозита: %s",
	},
	"Email": {
		"en": "Send me your email!\nExample: bip@thebest.com",
		"ru": "Отправь мне свой email!\nПример: bip@thebest.com",
	},
	"Coin price": {
		"en": "Send me a price for coins.",
		"ru": "Пришли мне цену за твои монеты.",
	},
	"New deposit": {
		"en": "New deposit!\nYou will receive at least  %f BIP.",
		"ru": "Новый депозит!\nВы получите минимум  %f BIP.",
	},
	"Exchange is successful": {
		"en": "Exchange is successful, you received  %f BIP.",
		"ru": "Обмен успешен!\nВы получили  %f BIP.",
	},
	"New deposit for sale": {
		"en": "New deposit for sale: %f %f at %f $.",
		"ru": "Новый депозит на продажу: %f %f по %f $.",
	},
	"BIP exchanged": {
		"en": "%f BIP exchanged for %f BTC.",
		"ru": "%f BIP обменяны на %f BTC.",
	},
	"Development": {
		"en": "In development stage",
		"ru": "В стадии разработки",
	},
	"Price": {
		"en": "💹 Currency",
		"ru": "💹 Текущий курс",
	},
	"Sell": {
		"en": "💰 Sell",
		"ru": "💰 Продать",
	},
	"Buy": {
		"en": "💰 Buy",
		"ru": "💰 Купить",
	},
	"Loots": {
		"en": "📃 My loots",
		"ru": "📃 Мои продажи",
	},
	"Empty loots": {
		"en": "You haven't got loots for sale.",
		"ru": "У вас нет лотов на продажу.",
	},
	"Error": {
		"en": "Something going wrong:(",
		"ru": "Что-то пошло не так:(",
	},
	"timeout": {
		"en": "Deposit timed out.",
		"ru": "Время ожидания депозита истекло.",
	},
}

func GetTranslate(key string, lang string) string {
	return Translates[key][lang]
}
