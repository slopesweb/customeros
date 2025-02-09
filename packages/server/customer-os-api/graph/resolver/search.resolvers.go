package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go/log"
)

// GcliSearch is the resolver for the gcli_search field.
func (r *queryResolver) GcliSearch(ctx context.Context, keyword string, limit *int) ([]*model.GCliItem, error) {
	ctx, span := tracing.StartGraphQLTracerSpan(ctx, "QueryResolver.GcliSearch", graphql.GetOperationContext(ctx))
	defer span.Finish()
	tracing.SetDefaultResolverSpanTags(ctx, span)
	span.LogFields(log.String("request.keyword", keyword))

	if keyword == "" {
		return []*model.GCliItem{}, nil
	}

	searchResultEntities, err := r.Services.SearchService.GCliSearch(ctx, keyword, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		graphql.AddErrorf(ctx, "Failed basic search for keyword %s", keyword)
		return nil, err
	}
	result := make([]*model.GCliItem, 0)
	for _, v := range *searchResultEntities {
		resultItem := model.GCliItem{}

		switch v.EntityType {
		case entity.SearchResultEntityTypeContact:
			resultItem = mapper.MapContactToGCliItem(*v.Node.(*neo4jentity.ContactEntity))
		case entity.SearchResultEntityTypeOrganization:
			resultItem = mapper.MapOrganizationToGCliItem(*v.Node.(*neo4jentity.OrganizationEntity))
		case entity.SearchResultEntityTypeEmail:
			resultItem = mapper.MapEmailToGCliItem(*v.Node.(*neo4jentity.EmailEntity))
		case entity.SearchResultEntityTypeState:
			resultItem = mapper.MapStateToGCliItem(*v.Node.(*neo4jentity.StateEntity))
		}
		result = append(result, &resultItem)
	}
	return result, nil
}
