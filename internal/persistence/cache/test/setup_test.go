package test

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"go-security/internal/model"
	"go-security/internal/persistence/cache"
	"go-security/internal/persistence/cache/impl"
	"testing"
)

var (
	ctx         = context.Background()
	permissions = map[model.Permission]struct{}{
		model.Permission{
			PermissionId:       1,
			PermissionTargetId: 100,
			TargetType:         "location",
		}: {},
	}
)

func setup(t *testing.T) (cache.UserPermissionCache, *miniredis.Miniredis) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}

	redisCache, err := impl.NewRedisCache(s.Addr(), "", 0)
	if err != nil {
		t.Fatalf("Failed to connect to miniredis: %v", err)
	}

	return redisCache, s
}
