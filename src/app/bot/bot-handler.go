package bot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	vocab "github.com/mrKitikat/telegrambottest/src/app/bot/vocabulary"
	"github.com/mrKitikat/telegrambottest/src/app/handler"
	stct "github.com/mrKitikat/telegrambottest/src/app/structs"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func (b *Bot) UpdateLoots(w http.ResponseWriter, r *http.Request) {

	loot := stct.UPDLoot{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loot); err != nil {
		handler.ResponError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	chatid, lang, err := b.DB.UpdateLoots(loot.Amount, loot.Tag)
	if err != nil {
		handler.ResponJSON(w, http.StatusBadGateway, err.Error())
		return
	}

	amount, err := strconv.ParseFloat(loot.SellAmount, 64)
	if err != nil {
		handler.ResponJSON(w, http.StatusBadGateway, err.Error())
		return
	}

	amountSell := float64((float64(loot.Price) / 100000.) * amount)
	ans := fmt.Sprintf(vocab.GetTranslate("Coin exchanged", lang), amount, loot.Coin, amountSell)
	msg := tgbotapi.NewMessage(chatid, ans)
	b.Bot.Send(msg)

	handler.ResponJSON(w, http.StatusOK, "Notification has been sent")

	return
}
