package ink

import (
	"testing"
)

func TestFs(t *testing.T) {
	router := NewRouter()
	router.RegisterFs("*", "/static/", "./", nil)

	_ = Run("127.0.0.1:80", map[string]*Router{"*": router})
}
