package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/model"
)

// OrganizationPlanCreate is the resolver for the organizationPlan_Create field.
func (r *mutationResolver) OrganizationPlanCreate(ctx context.Context, input model.OrganizationPlanInput) (*model.OrganizationPlan, error) {
	panic(fmt.Errorf("not implemented: OrganizationPlanCreate - organizationPlan_Create"))
}

// OrganizationPlanUpdate is the resolver for the organizationPlan_Update field.
func (r *mutationResolver) OrganizationPlanUpdate(ctx context.Context, input model.OrganizationPlanUpdateInput) (*model.OrganizationPlan, error) {
	panic(fmt.Errorf("not implemented: OrganizationPlanUpdate - organizationPlan_Update"))
}

// OrganizationPlanDuplicate is the resolver for the organizationPlan_Duplicate field.
func (r *mutationResolver) OrganizationPlanDuplicate(ctx context.Context, id string, organizationID string) (*model.OrganizationPlan, error) {
	panic(fmt.Errorf("not implemented: OrganizationPlanDuplicate - organizationPlan_Duplicate"))
}

// OrganizationPlanMilestoneCreate is the resolver for the organizationPlanMilestone_Create field.
func (r *mutationResolver) OrganizationPlanMilestoneCreate(ctx context.Context, input model.OrganizationPlanMilestoneInput) (*model.OrganizationPlanMilestone, error) {
	panic(fmt.Errorf("not implemented: OrganizationPlanMilestoneCreate - organizationPlanMilestone_Create"))
}

// OrganizationPlanMilestoneUpdate is the resolver for the organizationPlanMilestone_Update field.
func (r *mutationResolver) OrganizationPlanMilestoneUpdate(ctx context.Context, input model.OrganizationPlanMilestoneUpdateInput) (*model.OrganizationPlanMilestone, error) {
	panic(fmt.Errorf("not implemented: OrganizationPlanMilestoneUpdate - organizationPlanMilestone_Update"))
}

// OrganizationPlanMilestoneBulkUpdate is the resolver for the organizationPlanMilestone_BulkUpdate field.
func (r *mutationResolver) OrganizationPlanMilestoneBulkUpdate(ctx context.Context, input []*model.OrganizationPlanMilestoneUpdateInput) ([]*model.OrganizationPlanMilestone, error) {
	panic(fmt.Errorf("not implemented: OrganizationPlanMilestoneBulkUpdate - organizationPlanMilestone_BulkUpdate"))
}

// OrganizationPlanMilestoneReorder is the resolver for the organizationPlanMilestone_Reorder field.
func (r *mutationResolver) OrganizationPlanMilestoneReorder(ctx context.Context, input model.OrganizationPlanMilestoneReorderInput) (string, error) {
	panic(fmt.Errorf("not implemented: OrganizationPlanMilestoneReorder - organizationPlanMilestone_Reorder"))
}

// OrganizationPlanMilestoneDuplicate is the resolver for the organizationPlanMilestone_Duplicate field.
func (r *mutationResolver) OrganizationPlanMilestoneDuplicate(ctx context.Context, organizationID string, organizationPlanID string, id string) (*model.OrganizationPlanMilestone, error) {
	panic(fmt.Errorf("not implemented: OrganizationPlanMilestoneDuplicate - organizationPlanMilestone_Duplicate"))
}

// Milestones is the resolver for the milestones field.
func (r *organizationPlanResolver) Milestones(ctx context.Context, obj *model.OrganizationPlan) ([]*model.OrganizationPlanMilestone, error) {
	panic(fmt.Errorf("not implemented: Milestones - milestones"))
}

// RetiredMilestones is the resolver for the retiredMilestones field.
func (r *organizationPlanResolver) RetiredMilestones(ctx context.Context, obj *model.OrganizationPlan) ([]*model.OrganizationPlanMilestone, error) {
	panic(fmt.Errorf("not implemented: RetiredMilestones - retiredMilestones"))
}

// OrganizationPlan is the resolver for the organizationPlan field.
func (r *queryResolver) OrganizationPlan(ctx context.Context, id string) (*model.OrganizationPlan, error) {
	panic(fmt.Errorf("not implemented: OrganizationPlan - organizationPlan"))
}

// OrganizationPlansForOrganization is the resolver for the organizationPlansForOrganization field.
func (r *queryResolver) OrganizationPlansForOrganization(ctx context.Context, organizationID string) ([]*model.OrganizationPlan, error) {
	panic(fmt.Errorf("not implemented: OrganizationPlansForOrganization - organizationPlansForOrganization"))
}

// OrganizationPlans is the resolver for the organizationPlans field.
func (r *queryResolver) OrganizationPlans(ctx context.Context, retired *bool) ([]*model.OrganizationPlan, error) {
	panic(fmt.Errorf("not implemented: OrganizationPlans - organizationPlans"))
}

// OrganizationPlan returns generated.OrganizationPlanResolver implementation.
func (r *Resolver) OrganizationPlan() generated.OrganizationPlanResolver {
	return &organizationPlanResolver{r}
}

type organizationPlanResolver struct{ *Resolver }
