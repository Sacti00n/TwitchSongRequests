package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/saxypandabear/twitchsongrequests/internal/util"
	"github.com/saxypandabear/twitchsongrequests/pkg/db"
)

const PrefFormExplicitKey = "explicit"

type PreferenceHandler struct {
	prefs       db.PreferenceStore
	redirectURL string
}

func NewPreferenceHandler(d db.PreferenceStore, redirectURL string) *PreferenceHandler {
	return &PreferenceHandler{
		prefs:       d,
		redirectURL: redirectURL,
	}
}

func (h *PreferenceHandler) SavePreferences(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetUserIDFromRequest(r)

	if err != nil {
		log.Println("failed to get Twitch ID from request", err)
		http.Redirect(w, r, h.redirectURL, http.StatusFound)
		return
	}

	p, err := h.prefs.GetPreference(userID)
	if err != nil {
		log.Println("failed to get user preferences for", userID, err)
		http.Redirect(w, r, h.redirectURL, http.StatusFound)
		return
	}

	e := r.Form.Get(PrefFormExplicitKey)
	if e != "" && e != fmt.Sprintf("%t", p.ExplicitSongs) {
		p.ExplicitSongs = e == "true"
	}

	err = h.prefs.UpdatePreference(p)
	if err != nil {
		log.Println("failed to update user preferences for", userID, err)
		http.Redirect(w, r, h.redirectURL, http.StatusFound)
		return
	}

	log.Println("successfully saved user preferences for", userID)
	http.Redirect(w, r, h.redirectURL, http.StatusFound)
}
