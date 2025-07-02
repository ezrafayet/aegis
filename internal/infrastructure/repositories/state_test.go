package repositories

import (
	"aegis/internal/domain/entities"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestStateRepository_GetAndDeleteState(t *testing.T) {
	t.Run("should get and delete state", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&entities.State{})
		repo := NewStateRepository(db)
		state := entities.NewState("some-value")
		if err := repo.CreateState(state); err != nil {
			t.Fatal(err)
		}
		got, err := repo.GetAndDeleteState("some-value")
		if err != nil {
			t.Fatal(err)
		}
		if got.Value != state.Value {
			t.Fatal("expected state to be deleted", got.Value)
		}
		var count int64
		db.Model(&entities.State{}).Where("value = ?", state.Value).Count(&count)
		if count != 0 {
			t.Fatal("expected state to be deleted", count)
		}
	})
}
