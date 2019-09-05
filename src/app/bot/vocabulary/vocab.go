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
	"Select email": {
		"en": "Select email or enter a new one.",
		"ru": "Выберите email или введите новый.",
	},
	"Select bitcoin": {
		"en": "Select bitcoin address or enter a new one.",
		"ru": "Выберите bitcoin адрес или введите новый.",
	},
	"Select minter": {
		"en": "Select minter address or enter a new one.",
		"ru": "Выберите minter адрес или введите новый.",
	},
	"Select price": {
		"en": "Select price for your coins.",
		"ru": "Выберете цену за ваши монеты.",
	},
	"Coin": {
		"en": "Send name of the coin that you want to sell.",
		"ru": "Отправь название монеты которую ты хочешь продать.",
	},
	"Hello": {
		"en": "Hello, i'm an exchange bot BIP/BTC or BTC/BIP",
		"ru": "Привет, я бот для обмена BIP/BTC или BTC/BIP",
	},
	"Now": {
		"en": "📈 Now currency BIP/USD:",
		"ru": "📈 Сейчас курс BIP/USD:",
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
		"en": "Your minter deposit address and tag:",
		"ru": "Твой minter адрес для депозита и tag:",
	},
	"BTC deposit": {
		"en": "Your bitcoin deposit address:",
		"ru": "Твой bitcoin адрес для депозита:",
	},
	"Email": {
		"en": "Send me your email!\nExample: bip@thebest.com",
		"ru": "Отправь мне свой email!\nПример: bip@thebest.com",
	},
	"New email": {
		"en": "Enter new email",
		"ru": "Ввести новый email",
	},
	"New BTC": {
		"en": "Enter new bitcoin address",
		"ru": "Ввести новый bitcoin адрес",
	},
	"New minter": {
		"en": "Enter new minter address",
		"ru": "Ввести новый minter адрес",
	},
	"Coin price": {
		"en": "Send me a price for coins, format: 0.xxx.\nAllowable range: 0.1 - 0.32 $.",
		"ru": "Пришли мне цену за свои монеты, формат: 0.xxx.\nДопустимый диапазон: 0.1 - 0.32 $.",
	},
	"Coin name": {
		"en": "Wrong name of coin.\nExample: BIP, MNT ",
		"ru": "Неправильоне название монеты.\nПример: BIP, MNT.",
	},
	"New deposit": {
		"en": "New deposit!\nYou will receive at least  %.4f BIP.",
		"ru": "Новый депозит!\nВы получите минимум  %.4f BIP.",
	},
	"Exchange is successful": {
		"en": "Exchange is successful, you received  %.4f BIP.",
		"ru": "Обмен успешен!\nВы получили  %.4f BIP.",
	},
	"New deposit for sale": {
		"en": "New deposit for sale: %s %s at %.4f $.",
		"ru": "Новый депозит на продажу: %s %s по %.4f $.",
	},
	"Coin exchanged": {
		"en": "%.4f %s exchanged for %.4f BTC.",
		"ru": "%.4f %s обменяны на %.4f BTC.",
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
	"Wrong price": {
		"en": "Wrong price format!",
		"ru": "Неправильный формат цены!",
	},
	"timeout": {
		"en": "Deposit timed out.",
		"ru": "Время ожидания депозита истекло.",
	},
}

func GetTranslate(key string, lang string) string {
	return Translates[key][lang]
}
