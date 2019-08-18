package lang

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
	"New deposit": {
		"en": "New deposit!\n You will receive at least ? BIP.",
		"ru": "Новый депозит!\n Вы получите минимум ? BIP.",
	},
	"Exchange is successful": {
		"en": "Exchange is successful, you received ? BIP.",
		"ru": "Обмен успешен!\n Вы получили ? BIP.",
	},
	"New deposit for sale": {
		"en": "New deposit for sale: ? BIP at ? $",
		"ru": "Новый депозит на продажу: ? BIP по ?",
	},
	"BIP exchanged": {
		"en": "? BIP exchanged for ? BTC",
		"ru": "? BIP обменяны на ? BTC",
	},
	"Price": {
		"en": "Currency",
		"ru": "Текущий курс",
	},
	"Sell": {
		"en": "Sell",
		"ru": "Продать",
	},
	"Buy": {
		"en": "Buy",
		"ru": "Купить",
	},
	"Sales": {
		"en": "My sales",
		"ru": "Мои продажи",
	},
}


func GetTranslate(key string, lang string) string {
	return Translates[key][lang]
}
