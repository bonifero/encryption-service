package cryptoutil

import (
	"path/filepath"
	"testing"
)

func TestLoadOrGenerateKeyPair_GeneratesThenReloadsSameKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test-keystore.p12")

	generated, err := LoadOrGenerateKeyPair(path, "test-password", "crypto-service", 2048)
	if err != nil {
		t.Fatalf("не удалось сгенерировать ключевую пару: %v", err)
	}

	reloaded, err := LoadOrGenerateKeyPair(path, "test-password", "crypto-service", 2048)
	if err != nil {
		t.Fatalf("не удалось перезагрузить ключевую пару: %v", err)
	}

	if !generated.Equal(reloaded) {
		t.Fatalf("второй вызов должен загрузить тот же ключ с диска, а не сгенерировать новый")
	}
}

func TestLoadOrGenerateKeyPair_RejectsWrongPasswordOnReload(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test-keystore.p12")

	if _, err := LoadOrGenerateKeyPair(path, "correct-password", "crypto-service", 2048); err != nil {
		t.Fatalf("не удалось сгенерировать ключевую пару: %v", err)
	}

	if _, err := LoadOrGenerateKeyPair(path, "wrong-password", "crypto-service", 2048); err == nil {
		t.Fatalf("ожидалась ошибка при открытии хранилища неверным паролем")
	}
}
