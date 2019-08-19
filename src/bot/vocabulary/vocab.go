package vocab

var Translates = map[string]map[string]string{
	"Language": {
		"en": "Language",
		"ru": "Язык",
	},
	"English": {
		"en": "English",
		"ru": "Английский",
	},
	"Russian": {
		"en": "Russian",
		"ru": "Русский",
	},
	"Select": {
		"en": "Select, please, what you want:)",
		"ru": "Выберете, пожалуйста, что вы хотите:)",
	},
	"Hello!": {
		"en": "Hello, i'm an exchange bot: BIP/BTC or BTC/BIP",
		"ru": "Выберете, пожалуйста, что вы хотите:)",
	},
	"Now": {
		"en": "📈 Now currency BIP/USD %f $",
		"ru": "📈 Сейчас курс BIP/USD %f $",
	},
	"Send": {
		"en": "Send me your Minter Address:)",
		"ru": "Отправь мне свой адрес в Minter:)",
	},
	"New deposit": {
		"en": "New deposit!\n You will receive at least  %f BIP.",
		"ru": "Новый депозит!\n Вы получите минимум  %f BIP.",
	},
	"Exchange is successful": {
		"en": "Exchange is successful, you received  %f BIP.",
		"ru": "Обмен успешен!\n Вы получили  %f BIP.",
	},
	"New deposit for sale": {
		"en": "New deposit for sale: %f BIP at %f $",
		"ru": "Новый депозит на продажу: %f BIP по %f",
	},
	"BIP exchanged": {
		"en": "%f BIP exchanged for %f BTC",
		"ru": "%f BIP обменяны на %f BTC",
	},
	"Price": {
		"en": "💹Currency",
		"ru": "💹Текущий курс",
	},
	"Sell": {
		"en": "💰Sell",
		"ru": "💰Продать",
	},
	"Buy": {
		"en": "💰Buy",
		"ru": "💰Купить",
	},
	"Sales": {
		"en": "📃My sales",
		"ru": "📃Мои продажи",
	},
}

func GetTranslate(key string, lang string) string {
	return Translates[key][lang]
}
