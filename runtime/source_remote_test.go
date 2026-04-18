package runtime

import (
	"testing"

	"github.com/cheryl-chun/confgen/internal/tree"
)

func TestRemoteConfigSource_KeyToPath(t *testing.T) {
	source := &RemoteConfigSource{Prefix: "config/app"}

	path := source.KeyToPath("/config/app/server/port")
	if path != "server.port" {
		t.Fatalf("KeyToPath() = %q, want %q", path, "server.port")
	}

	path = source.KeyToPath("config/app/database/host")
	if path != "database.host" {
		t.Fatalf("KeyToPath() = %q, want %q", path, "database.host")
	}
}

func TestRemoteConfigSource_PathToKey(t *testing.T) {
	source := &RemoteConfigSource{Prefix: "config/app"}

	key := source.PathToKey("server.host")
	if key != "config/app/server/host" {
		t.Fatalf("PathToKey() = %q, want %q", key, "config/app/server/host")
	}

	key = source.PathToKey("database.port")
	if key != "config/app/database/port" {
		t.Fatalf("PathToKey() = %q, want %q", key, "config/app/database/port")
	}
}

func TestEtcdSource_Priority(t *testing.T) {
	source := NewEtcdSource([]string{"127.0.0.1:2379"}, "config/app")
	if source.Priority() != tree.SourceRemote {
		t.Fatalf("Priority() = %v, want %v", source.Priority(), tree.SourceRemote)
	}
}
