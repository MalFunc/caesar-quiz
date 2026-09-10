package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"caesar-quiz/internal/api"
	"caesar-quiz/internal/models"
	"caesar-quiz/internal/services"
	"caesar-quiz/internal/ws"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(
		&models.Game{}, &models.Player{}, &models.Question{}, &models.Answer{}, &models.HallOfFame{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func postJSON(r http.Handler, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func newTestServer(t *testing.T) (*gorm.DB, http.Handler, *models.Game, *models.Question, *models.Player) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	h := ws.NewHub(db)
	go h.Run()

	game := &models.Game{Code: "TEST01", Shift: 3, HostToken: "secret-token", IsActive: true}
	if err := db.Create(game).Error; err != nil {
		t.Fatalf("create game: %v", err)
	}
	q := &models.Question{GameID: game.ID, Plaintext: "HELLO", Cipher: services.CaesarEncrypt("HELLO", 3), Shift: 3}
	if err := db.Create(q).Error; err != nil {
		t.Fatalf("create question: %v", err)
	}
	p := &models.Player{GameID: game.ID, Name: "Budi"}
	if err := db.Create(p).Error; err != nil {
		t.Fatalf("create player: %v", err)
	}

	r := gin.New()
	api.RegisterRoutes(r, db, h, []string{"*"})
	return db, r, game, q, p
}

func TestSubmitAnswerFlow(t *testing.T) {
	_, r, game, q, p := newTestServer(t)

	w := postJSON(r, "/game/"+game.ID.String()+"/start", map[string]any{}, map[string]string{"X-Host-Token": "secret-token"})
	if w.Code != http.StatusOK {
		t.Fatalf("start: expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	w = postJSON(r, "/game/"+game.ID.String()+"/answer", map[string]any{
		"player_id": p.ID.String(), "question_id": q.ID.String(), "answer": "HELLO",
	}, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("answer: expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if ok, _ := resp["correct"].(bool); !ok {
		t.Fatalf("expected correct=true, got %v", resp)
	}
	if ms, ok := resp["time_ms"].(float64); !ok || ms < 0 {
		t.Fatalf("expected numeric time_ms, got %v", resp["time_ms"])
	}

	// Duplicate submit harus idempotent (tidak dobel).
	w = postJSON(r, "/game/"+game.ID.String()+"/answer", map[string]any{
		"player_id": p.ID.String(), "question_id": q.ID.String(), "answer": "HELLO",
	}, nil)
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if dup, _ := resp["duplicate"].(bool); !dup {
		t.Fatalf("expected duplicate=true, got %v", resp)
	}
}

func TestStartRequiresHostToken(t *testing.T) {
	_, r, game, _, _ := newTestServer(t)

	w := postJSON(r, "/game/"+game.ID.String()+"/start", map[string]any{}, nil)
	if w.Code != http.StatusForbidden {
		t.Fatalf("start without token: expected 403 got %d", w.Code)
	}
}

func TestAnswerRejectsCrossGame(t *testing.T) {
	db, r, _, q, _ := newTestServer(t)

	otherGame := &models.Game{Code: "OTHER1", Shift: 5, IsActive: true}
	if err := db.Create(otherGame).Error; err != nil {
		t.Fatalf("create other game: %v", err)
	}
	intruder := &models.Player{GameID: otherGame.ID, Name: "Eve"}
	if err := db.Create(intruder).Error; err != nil {
		t.Fatalf("create intruder: %v", err)
	}

	w := postJSON(r, "/game/"+q.GameID.String()+"/answer", map[string]any{
		"player_id": intruder.ID.String(), "question_id": q.ID.String(), "answer": "HELLO",
	}, nil)
	if w.Code != http.StatusForbidden {
		t.Fatalf("cross-game answer: expected 403 got %d body=%s", w.Code, w.Body.String())
	}
}
