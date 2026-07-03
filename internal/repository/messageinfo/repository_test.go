package messageinfo_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/SelickSD/DemoBot.git/internal/db"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/SelickSD/DemoBot.git/internal/repository/messageinfo"
	"github.com/SelickSD/DemoBot.git/test/testdb"
)

var (
	ctx  = context.Background()
	repo *messageinfo.Repository
)

func TestMain(m *testing.M) {
	pg, err := testdb.StartPostgres(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := testdb.RunMigrations(pg.DSN); err != nil {
		log.Fatal(err)
	}

	dbPool, err := testdb.InitPool(ctx, pg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	db.Pool = dbPool

	repo = messageinfo.NewRepository()

	code := m.Run()

	_ = pg.Container.Terminate(ctx)
	os.Exit(code)
}

func TestSaveAndGetByChatID(t *testing.T) {
	err := repo.DeleteAll(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Save(ctx, messageinfo.MessageInfo{
		MessageID: 1,
		ChatID:    100,
		Message:   "hello test",
		UserID:    42,
	})
	if err != nil {
		t.Fatal(err)
	}

	msgs, err := repo.GetByChatID(ctx, 100, 10)
	if err != nil {
		t.Fatal(err)
	}

	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}

	if msgs[0].Message != "hello test" {
		t.Fatalf("unexpected message: %s", msgs[0].Message)
	}
}

func TestGetByChatID_OrderAndLimit(t *testing.T) {
	err := repo.DeleteAll(ctx)
	if err != nil {
		t.Fatal(err)
	}

	chatID := int64(200)

	for i := 1; i <= 5; i++ {
		err = repo.Save(ctx, messageinfo.MessageInfo{
			MessageID: int64(i),
			ChatID:    chatID,
			Message:   fmt.Sprintf("msg %d", i),
			UserID:    1,
		})
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(10 * time.Millisecond)
	}

	msgs, err := repo.GetByChatID(ctx, chatID, 3)
	if err != nil {
		t.Fatal(err)
	}

	if len(msgs) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(msgs))
	}

	// порядок: самые новые первые
	if msgs[0].MessageID <= msgs[1].MessageID {
		t.Fatal("messages are not ordered by created_at desc")
	}
}
