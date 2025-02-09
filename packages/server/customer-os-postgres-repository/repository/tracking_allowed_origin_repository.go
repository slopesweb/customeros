package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TrackingAllowedOriginRepository interface {
	GetTenantForOrigin(ctx context.Context, origin string) (*string, error)
}

type trackingAllowedOriginRepositoryImpl struct {
	gormDb *gorm.DB
}

func NewTrackingAllowedOriginRepository(gormDb *gorm.DB) TrackingAllowedOriginRepository {
	return &trackingAllowedOriginRepositoryImpl{gormDb: gormDb}
}

func (repo *trackingAllowedOriginRepositoryImpl) GetTenantForOrigin(ctx context.Context, origin string) (*string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TrackingAllowedOriginRepository.GetTenantForOrigin")
	defer span.Finish()
	span.LogFields(log.String("origin", origin))

	var result entity.TrackingAllowedOrigin
	err := repo.gormDb.Model(&entity.TrackingAllowedOrigin{}).Find(&result, "origin = ?", origin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logrus.Errorf("error while getting import allowed organization: %v", err)
		return nil, err
	}
	if result.Tenant == "" {
		return nil, nil
	}

	return &result.Tenant, nil
}
