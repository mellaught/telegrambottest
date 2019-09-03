package bot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	stct "telegrambottest/src/app/bipdev/structs"
	"telegrambottest/src/app/handler"

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

	chatid, _, err := b.DB.UpdateLoots(loot.Amount, loot.Tag)
	if err != nil {
		handler.ResponJSON(w, http.StatusBadGateway, err.Error())
		return
	}

	amount, err := strconv.ParseInt(loot.Amount, 10, 64)
	if err != nil {
		handler.ResponJSON(w, http.StatusBadGateway, err.Error())
		return
	}
	msg := tgbotapi.NewMessage(chatid, fmt.Sprintf("Now you have got %d", amount))
	b.Bot.Send(msg)

	handler.ResponJSON(w, http.StatusOK, "Notification has been sent")
	return
}
