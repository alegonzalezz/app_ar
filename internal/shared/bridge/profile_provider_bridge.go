package bridge

import (
	"context"

	adminDomain "gcp-serverless-app/internal/administrative/domain"
	authDomain "gcp-serverless-app/internal/auth/domain"
	userDomain "gcp-serverless-app/internal/user/domain"
	workerDomain "gcp-serverless-app/internal/worker/domain"
)

type profileProviderBridge struct {
	userRepo   userDomain.UserRepository
	workerRepo workerDomain.WorkerRepository
	adminRepo  adminDomain.AdministrativeRepository
}

func NewProfileProviderBridge(
	userRepo userDomain.UserRepository,
	workerRepo workerDomain.WorkerRepository,
	adminRepo adminDomain.AdministrativeRepository,
) authDomain.ProfileProvider {
	return &profileProviderBridge{
		userRepo:   userRepo,
		workerRepo: workerRepo,
		adminRepo:  adminRepo,
	}
}

func (b *profileProviderBridge) GetProfile(ctx context.Context, profileID string, profileType string) (string, interface{}, error) {
	switch profileType {
	case "user":
		user, err := b.userRepo.GetByID(ctx, profileID)
		if err != nil {
			return "", nil, err
		}
		return user.Name, user, nil

	case "worker":
		worker, err := b.workerRepo.GetByID(ctx, profileID)
		if err != nil {
			return "", nil, err
		}
		return worker.Name, worker, nil

	case "administrative":
		admin, err := b.adminRepo.GetByID(ctx, profileID)
		if err != nil {
			return "", nil, err
		}
		return admin.Name, admin, nil

	default:
		return "", nil, nil
	}
}
