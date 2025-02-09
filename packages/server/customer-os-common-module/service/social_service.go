package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type SocialService interface {
	Update(ctx context.Context, tenant string, entity neo4jentity.SocialEntity) (*neo4jentity.SocialEntity, error)
	GetAllForEntities(ctx context.Context, tenant string, linkedEntityType neo4jenum.EntityType, linkedEntityIds []string) (*neo4jentity.SocialEntities, error)
	Remove(ctx context.Context, tenant string, socialId string) error
}

type socialService struct {
	log      logger.Logger
	services *Services
}

func NewSocialService(log logger.Logger, services *Services) SocialService {
	return &socialService{
		log:      log,
		services: services,
	}
}

func (s *socialService) GetAllForEntities(ctx context.Context, tenant string, linkedEntityType neo4jenum.EntityType, linkedEntityIds []string) (*neo4jentity.SocialEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.GetAllForEntities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("linkedEntityType", string(linkedEntityType)), log.Object("linkedEntityIds", linkedEntityIds))

	socials, err := s.services.Neo4jRepositories.SocialReadRepository.GetAllForEntities(ctx, tenant, linkedEntityType, linkedEntityIds)
	if err != nil {
		return nil, err
	}
	socialEntities := make(neo4jentity.SocialEntities, 0, len(socials))
	for _, v := range socials {
		socialEntity := neo4jmapper.MapDbNodeToSocialEntity(v.Node)
		socialEntity.DataloaderKey = v.LinkedNodeId
		socialEntities = append(socialEntities, *socialEntity)
	}
	return &socialEntities, nil
}

func (s *socialService) Update(ctx context.Context, tenant string, socialEntity neo4jentity.SocialEntity) (*neo4jentity.SocialEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	updatedLocationNode, err := s.services.Neo4jRepositories.SocialWriteRepository.Update(ctx, tenant, socialEntity)
	if err != nil {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToSocialEntity(updatedLocationNode), nil
}

func (s *socialService) Remove(ctx context.Context, tenant string, socialId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.Remove")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, socialId)

	return s.services.Neo4jRepositories.SocialWriteRepository.PermanentlyDelete(ctx, tenant, socialId)
}
