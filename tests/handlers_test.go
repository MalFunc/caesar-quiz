package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"caesar-quiz/internal/api"
	"caesar-quiz/internal/ws"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&api.GameModel{}, &api.PlayerModel{}, &api.QuestionModel{}, &api.AnswerModel{}, &api.HallOfFameModel{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestSubmitAnswerFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	h := ws.NewHub(db)
	go h.Run()

	// create game
	game := &api.GameModel{Code: "TEST01", Shift: 3, ShiftRandom: false, Status: "waiting"}
	db.Create(game)

	// add question
	q := &api.QuestionModel{GameID: game.ID, Plaintext: "HELLO", Shift: 3, Ciphertext: api.CaesarEncrypt("HELLO", 3)}
	db.Create(q)

	// add player
	p := &api.PlayerModel{GameID: game.ID, Name: "Budi"}
	db.Create(p)

	// start
	now := time.Now()
	game.Status = "running"
	game.StartedAt = &now
	db.Save(game)

	r := gin.Default()
	api.RegisterRoutes(r, db, h)

	// submit correct answer
	body := map[string]interface{}{
		"player_id":   p.ID.String(),
		"question_id": q.ID.String(),
		"answer":      "HELLO",
		"time_ms":     1000,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/game/TEST01/answer", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if ok, _ := resp["is_correct"].(bool); !ok {
		t.Fatalf("expected correct")
	}
	if finished, _ := resp["finished"].(bool); !finished {
		t.Fatalf("expected finished true")
	}
}
	