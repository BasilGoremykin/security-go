package impl

import (
	"go-security/internal/model"
	"go-security/internal/persistence"
	"go-security/internal/persistence/cache"
	"go-security/internal/persistence/repository"
	"go-security/internal/service"
	"log"
)

type CachedRepoPermissionService struct {
	repo  repository.UserPermissionRepository
	cache cache.UserPermissionCache
}

func NewCachedRepoPermissionService(repo repository.UserPermissionRepository, cache cache.UserPermissionCache) service.UserPermissionService {
	return &CachedRepoPermissionService{repo: repo, cache: cache}
}

func (s *CachedRepoPermissionService) FindByID(id int64) (*model.UserPermissions, error) {
	cachedPermissions := s.cache.GetPermissions(id)

	if cachedPermissions != nil && len(cachedPermissions) > 0 {
		return &model.UserPermissions{UserID: id, UserPermissions: cachedPermissions}, nil
	}

	userPermissionsDtos, err := s.repo.GetAllByID(id)
	if err != nil {
		return nil, err
	}
	permissions := persistence.MapDtosToModelSet(userPermissionsDtos)
	s.cache.SetPermissions(id, permissions) // no need to handle error

	return &model.UserPermissions{UserID: userPermissionsDtos[0].UserID, UserPermissions: permissions}, nil
}

func (s *CachedRepoPermissionService) AddPermissionToUser(userId int64, permission *model.Permission) error {
	err := s.cache.DeletePermissions(userId)
	if err != nil {
		log.Println("couldn't delete old permissions from cache")
		return err
	}
	return s.repo.AddPermissionToUser(userId, permission)
}

func (s *CachedRepoPermissionService) RemovePermissionFromUser(userId int64, permission *model.Permission) error {
	err := s.cache.DeletePermissions(userId)
	if err != nil {
		log.Println("couldn't delete old permissions from cache")
		return err
	}
	return s.repo.RemovePermissionFromUser(userId, permission)
}

func (s *CachedRepoPermissionService) DeletePermissions(userId int64) error {
	err := s.cache.DeletePermissions(userId)
	if err != nil {
		log.Println("couldn't delete old permissions from cache")
		return err
	}
	return s.repo.DeletePermissions(userId)
}
