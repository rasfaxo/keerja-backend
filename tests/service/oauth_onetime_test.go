package service_test

import (
    "context"
    "testing"
    "time"

    "github.com/alicebob/miniredis/v2"
    "github.com/redis/go-redis/v9"
    "github.com/stretchr/testify/require"

    "keerja-backend/internal/service"
)

func TestCreateAndConsumeOneTimeCode_Redis(t *testing.T) {
    // Start miniredis
    mr, err := miniredis.Run()
    require.NoError(t, err)
    defer mr.Close()

    client := redis.NewClient(&redis.Options{Addr: mr.Addr(), DB: 0})
    stateStore := service.NewRedisOAuthStateStore(client)

    svc := service.NewOAuthService(nil, nil, service.OAuthConfig{}, "secret", time.Hour, stateStore, []string{"myapp://oauth-callback"})

    ctx := context.Background()
    jwt := "test.jwt.token"

    code, err := svc.CreateOneTimeCode(ctx, jwt, time.Minute)
    require.NoError(t, err)
    require.NotEmpty(t, code)

    got, err := svc.ConsumeOneTimeCode(ctx, code)
    require.NoError(t, err)
    require.Equal(t, jwt, got)

    // second consume should fail
    _, err = svc.ConsumeOneTimeCode(ctx, code)
    require.Error(t, err)
}

func TestCreateAndConsumeOneTimeCode_InMemory(t *testing.T) {
    stateStore := service.NewInMemoryOAuthStateStore()
    svc := service.NewOAuthService(nil, nil, service.OAuthConfig{}, "secret", time.Hour, stateStore, []string{"myapp://oauth-callback"})

    ctx := context.Background()
    jwt := "inmemory.jwt"

    code, err := svc.CreateOneTimeCode(ctx, jwt, time.Minute)
    require.NoError(t, err)
    require.NotEmpty(t, code)

    got, err := svc.ConsumeOneTimeCode(ctx, code)
    require.NoError(t, err)
    require.Equal(t, jwt, got)

    // second consume should fail
    _, err = svc.ConsumeOneTimeCode(ctx, code)
    require.Error(t, err)
}
