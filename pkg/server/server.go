package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/snsvistunov/tg-bot-pocketer/pkg/repository"
	"github.com/zhashkevych/go-pocket-sdk"
)

type AuthorizationServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string
}

func NewAuthorizationServer(pocketClient *pocket.Client, tokenRepository repository.TokenRepository, redirectURL string) *AuthorizationServer {
	return &AuthorizationServer{pocketClient: pocketClient, tokenRepository: tokenRepository, redirectURL: redirectURL}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":8088",
		Handler: s,
	}
	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	chatIDParam := r.URL.Query().Get("chat_id")
	fmt.Println("------------------------------------------------------------------------------------------")
	fmt.Println(chatIDParam)
	fmt.Println("------------------------------------------------------------------------------------------")
	if chatIDParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatID, err := strconv.ParseInt(chatIDParam, 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestToken, err := s.tokenRepository.Get(chatID, repository.RequestTokens)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	authResp, err := s.pocketClient.Authorize(r.Context(), requestToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.tokenRepository.Save(chatID, authResp.AccessToken, repository.AccessTokens)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Location", s.redirectURL)
	w.WriteHeader(http.StatusMovedPermanently)
}
