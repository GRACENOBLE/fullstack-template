package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	ws "backend/internal/infrastructure/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// In debug mode (local dev) all origins are allowed so browsers can connect
	// from any port. In staging/production, restrict to BLUEPRINT_WS_ALLOWED_ORIGIN.
	CheckOrigin: func(r *http.Request) bool {
		if gin.IsDebugging() {
			return true
		}
		allowed := os.Getenv("BLUEPRINT_WS_ALLOWED_ORIGIN")
		if allowed == "" {
			return false
		}
		return r.Header.Get("Origin") == allowed
	},
}

// WsHandler godoc
//
//	@Summary		Open a WebSocket connection
//	@Description	Upgrades HTTP to WebSocket. Pass a Firebase ID token as `?token=<token>`. Returns 401 when the token is missing or invalid.
//	@Tags			websocket
//	@Produce		json
//	@Param			token	query		string	true	"Firebase ID token"
//	@Success		101		{string}	string	"Switching Protocols"
//	@Failure		401		{object}	map[string]string
//	@Router			/ws [get]
func (h *Handler) WsHandler(c *gin.Context) {
	if h.verifier != nil {
		token := c.Query("token")
		if token == "" {
			JSONError(c, http.StatusUnauthorized, "unauthorized", "missing token")
			return
		}
		if _, err := h.verifier.VerifyIDToken(c.Request.Context(), token); err != nil {
			JSONError(c, http.StatusUnauthorized, "unauthorized", "invalid or expired token")
			return
		}
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// Upgrader already wrote the error response.
		return
	}

	client := ws.NewClient(h.hub, conn)
	go client.WritePump()
	go client.ReadPump()
}
