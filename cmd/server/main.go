package main

import (
	"fmt"
	"go-security/internal/persistence/cache/impl"
	impl2 "go-security/internal/persistence/repository/impl"
	"go-security/internal/properties"
	"go-security/internal/server"
	impl3 "go-security/internal/service/impl"
	away "go-security/pb/go_security"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
)

func main() {
	var props properties.SecurityProperties
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	if err = yaml.Unmarshal(data, &props); err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	postgresProps := props.Properties.Postgres
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable TimeZone=UTC",
		postgresProps.Username, postgresProps.Password, postgresProps.Host, postgresProps.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("can't establish connection to postgres")
	}

	cacheProps := props.Properties.Cache
	cache, err := impl.NewRedisCache(cacheProps.Addr, cacheProps.Password, cacheProps.Db)
	if err != nil {
		log.Fatal("can't establish connection to redis")
	}

	ssoProps := props.Properties.Sso
	oauthService := impl3.NewSsoService(ssoProps)

	permissionRepo := impl2.NewUserPermissionPostgresRepository(db)
	profileRepo := impl2.NewUserProfilePostgresRepository(db)

	permissionService := impl3.NewCachedRepoPermissionService(permissionRepo, cache)

	grpcProps := props.Properties.Grpc
	listener, err := net.Listen(grpcProps.Network, grpcProps.Port)
	if err != nil {
		log.Fatalf("couldn't establish listener on port %s", grpcProps.Port)
	}

	grpcServer := grpc.NewServer()
	away.RegisterSecurityServiceServer(grpcServer, server.NewSecurityService(permissionService, profileRepo, oauthService))

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
