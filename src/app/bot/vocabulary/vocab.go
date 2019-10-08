package vocab

var Translates = map[string]map[string]string{
	"Hello": {
		"en": "Hello, I'm an exchange bot for *BIP ⇆ BTC.* \n*Please choose your language.*",
		"ru": "Привет! Я бот для обмена *BIP ⇆ BTC.* \n*Пожалуйста, выберите язык.*",
	},
	"Language": {
		"en": "Language",
		"ru": "Язык",
	},
	"english": {
		"en": "English",
		"ru": "английский",
	},
	"russian": {
		"en": "Russian",
		"ru": "Русский",
	},
	"Installed": {
		"en": "Installed",
		"ru": "Установлен",
	},
	"Menu": {
		"en": "📕 Menu",
		"ru": "📕 Меню",
	},
	"Settings": {
		"en": "🔧 Settings",
		"ru": "🔧 Настройки",
	},
	"Yes": {
		"en": "Yes",
		"ru": "Да",
	},
	"No": {
		"en": "No",
		"ru": "Нет",
	},
	"Select": {
		"en": "The current rate is *$%.2f* (%s %%).\n\nHere you can buy or sell *BIP* and track the orders you've placed.",
		"ru": "Текущий курс: *$%.2f* (%s %%).\n\nЗдесь вы можете купить или продать *BIP*, а также отслеживать созданные заявки.",
	},
	// Buy
	// 1
	"New minter": {
		"en": "Enter your address on the Minter network.\n\n*Example:* Mx00000000000000000000000000000000000001",
		"ru": "Введите ваш адрес в сети Minter.\n\n*Пример:* Mx00000000000000000000000000000000000001",
	},
	"Select minter": {
		"en": "Choose your address on the Minter network or enter a new one.\n\n*Example:* Mx00000000000000000000000000000000000001",
		"ru": "Выберите ваш адрес в сети Minter или введите новый.\n\n*Пример:* Mx00000000000000000000000000000000000001",
	},
	"Wrong minter": {
		"en": "Re-check the Minter address you entered. It should contain *42 characters* and start with *Mx*.",
		"ru": "Проверьте правильность введённого адреса, он должен содержать *42 символа* и начинаться с *Mx*.",
	},
	// 2
	"New email": {
		"en": "Enter your e-mail address.\n\n*Example:* mail@example.com",
		"ru": "Введите ваш почтовый адрес.\n\n*Пример:* mail@example.com",
	},
	"Select email": {
		"en": "Choose your e-mail address or enter a new one.\n\n*Example:* mail@example.com",
		"ru": "Выберите ваш почтовый адрес или введите новый.\n\n*Пример:* mail@example.com",
	},
	"Wrong email": {
		"en": "Re-check the e-mail address you entered.",
		"ru": "Проверьте правильность введённого адреса.",
	},
	// 3 Send BTC ... 2 confirmations...
	"Send deposit": {
		"en": "Send BTC to the following address. After *2* confirmations, you will receive BIP to the Minter address you've specified before.\n\nThe *current rate* is $%.2f (%s %%)\n\n" +
			"💡 1 BTC will now buy you *%s* BIP. That's a *%.2f %% bonus* to the indicative price.",
		"ru": "Отправьте BTC на следующий адрес, после *2* подтверждений сети, вы получите BIP на указанный вами адрес в сети Minter.\n\n*Текущий курс:* $%.2f (%s %%)\n\n" +
			"💡 Сейчас за 1 BTC вы можете купить *%s* BIP, это на *%.2f %% больше* актуальной цены.",
	},
	// 4
	"Check": {
		"en": "Check",
		"ru": "Проверить",
	},
	"Wait deposit": {
		"en": "Waiting for the BTC transaction…",
		"ru": "Ожидание транзакции BTC…",
	},
	"New deposit": {
		"en": "BTC is already on the way. You will get at least %.2f BIP.",
		"ru": "BTC уже в пути, вы получите как минимум %.2f BIP.",
	},
	"No buy": {
		"en": "You've got no buy orders.",
		"ru": "У вас нет заявок на покупку.",
	},
	// 5
	"Exchange is successful": {
		"en": "🎉 *%.2f* BIP has been sent to your address.",
		"ru": "🎉 *%.2f* BIP были отправлены на ваш адрес.",
	},
	// 1
	"Coin": {
		"en": "Enter the ticker symbol of a coin you want to sell.\n\n*Example*: BIP",
		"ru": "Введите название монеты, которую хотите продать.\n\n*Пример*: BIP",
	},
	// 1
	"Wrong coin name": {
		"en": "⚠️ *Error*\n\nSuch a coin does not exist.",
		"ru": "⚠️ *Ошибка*\n\nТакой монеты не существует.",
	},
	// 2
	"Select price": {
		"en": "Specify the *USD* price at which you are willing to sell your coins.\n\n*Example*: 0.32",
		"ru": "Укажите цену в *USD*, которую хотите установить для монет.\n\n*Пример*: 0.32",
	},
	// 2
	"Wrong price": {
		"en": "⚠️ *Error*\n\nThe possible price range is *$0.01*–*$0.32*. The value should be strictly numerical.",
		"ru": "⚠️ *Ошибка*\n\nВозможный диапазон цены: от *$0.01* до *$0.32*, вводить цену нужно без символов обозначающих валюту и букв.",
	},
	// 3
	"New bitcoin": {
		"en": "Enter your *Bitcoin* address.",
		"ru": "Введите ваш *Bitcoin* адрес.",
	},
	"Send bitcoin": {
		"en": "Choose your *Bitcoin* address or enter a new one.",
		"ru": "Выберите ваш *Bitcoin* адрес или введите новый.",
	},
	"Select bitcoin": {
		"en": "Choose *Bitcoin* address or enter a new one.",
		"ru": "Выберите *Bitcoin* адрес или введите новый.",
	},
	// 3
	"Wrong bitcoin": {
		"en": "⚠️ *Error*\n\nRe-check the BTC address you entered.",
		"ru": "⚠️ *Ошибка*\n\nПроверьте правильность введённого BTC адреса.",
	},
	// 4
	"Save": {
		"en": "Do you want to save this address for future sales?",
		"ru": "Сохранить введённый адрес для следующих продаж?",
	},
	// 5
	"Send your coins": {
		"en": "Send *%s* to the address below.\n\n⚠️ Do not send less than *1 000 %s* in one transaction.\n\nYou can track your order at\n%s",
		"ru": "Отправьте *%s* на указанный ниже адрес.\n\n⚠️ Не отправляйте меньше *1 000 %s* в одной транзакции.\n\nВы можете отслеживать заявку по этой ссылке:\n%s",
	},
	//6
	"Share": {
		"en": "Share",
		"ru": "Поделиться",
	},
	"Wait deposit coin": {
		"en": "Waiting for %s for sale...",
		"ru": "Ожидание %s на продажу...",
	},
	"No sell": {
		"en": "You've got no sell orders.",
		"ru": "У вас нет заявок на продажу.",
	},
	// 7
	"New deposit for sale": {
		"en": "A new sell order: *%s* %s at *%.2f* $.",
		"ru": "Новая заявка на продажу: *%s* %s по *%.2f* $.",
	},
	// Заявки
	"Your loots": {
		"en": "📔 *Orders*\n\nIn this section, you can see all of your open orders.",
		"ru": "📔 *Заявки*\n\nВ этом разделе вы можете найти все активные заявки.",
	},
	"Loot": {
		"en": "Selling %s %s at $%v",
		"ru": "Продажа %s %s по $%v",
	},
	"Empty loots": {
		"en": "You've got no sell orders.",
		"ru": "У вас нет лотов на продажу.",
	},
	// Кнопки
	"Sell": {
		"en": "Sell",
		"ru": "Продать",
	},
	"Buy": {
		"en": "Buy",
		"ru": "Купить",
	},
	"Loots": {
		"en": "My orders",
		"ru": "Мои заявки",
	},
	"Cancel": {
		"en": "« Cancel",
		"ru": "« Назад",
	},
	// ------------------------------------------------------
	"Now": {
		"en": "The current rate is $%.2f (+%.2f %)",
		"ru": "Текущий курс: $%.2f (+%.2f %)",
	},
	"Coin exchanged": {
		"en": "%.4f %s has been exchanged for %.4f BTC.",
		"ru": "%.4f %s обменяны на %.4f BTC.",
	},
	"Minter deposit and tag": {
		"en": "Your Minter deposit address and tag:",
		"ru": "Ваш Minter адрес для депозита и tag:",
	},
	"Error": {
		"en": "⚠️*Error*",
		"ru": "⚠️*Ошибка*",
	},
	"timeout": {
		"en": "Deposit timed out.",
		"ru": "Время ожидания депозита истекло.",
	},
}

func GetTranslate(key string, lang string) string {
	return Translates[key][lang]
}
