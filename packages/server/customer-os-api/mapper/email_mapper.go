package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEmailInputToEntity(input *model.EmailInput) *neo4jentity.EmailEntity {
	if input == nil {
		return nil
	}
	emailEntity := neo4jentity.EmailEntity{
		RawEmail:      input.Email,
		Label:         utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
		Primary:       utils.IfNotNilBool(input.Primary),
		Source:        neo4jentity.DataSourceOpenline,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
		AppSource:     utils.IfNotNilString(input.AppSource),
	}
	if len(emailEntity.AppSource) == 0 {
		emailEntity.AppSource = constants.AppSourceCustomerOsApi
	}
	return &emailEntity
}

func MapEntityToEmail(entity *neo4jentity.EmailEntity) *model.Email {
	var label = model.EmailLabel(entity.Label)
	if !label.IsValid() {
		label = ""
	}
	return &model.Email{
		ID:            entity.Id,
		Email:         utils.StringPtrFirstNonEmptyNillable(entity.Email, entity.RawEmail),
		RawEmail:      utils.StringPtrNillable(entity.RawEmail),
		Label:         &label,
		Primary:       entity.Primary,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		EmailValidationDetails: &model.EmailValidationDetails{
			Validated:      entity.Validated,
			IsReachable:    entity.IsReachable,
			IsValidSyntax:  entity.IsValidSyntax,
			CanConnectSMTP: entity.CanConnectSMTP,
			AcceptsMail:    entity.AcceptsMail,
			HasFullInbox:   entity.HasFullInbox,
			IsCatchAll:     entity.IsCatchAll,
			IsDeliverable:  entity.IsDeliverable,
			IsDisabled:     entity.IsDisabled,
			Error:          entity.Error,
			IsDisposable:   entity.IsDisposable,
			IsRoleAccount:  entity.IsRoleAccount,
		},
	}
}

func MapEmailInputToLocalEntity(input *model.EmailInput) *neo4jentity.EmailEntity {
	if input == nil {
		return nil
	}
	emailEntity := neo4jentity.EmailEntity{
		RawEmail:      input.Email,
		Label:         utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
		Primary:       utils.IfNotNilBool(input.Primary),
		Source:        neo4jentity.DataSourceOpenline,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
		AppSource:     utils.IfNotNilString(input.AppSource),
	}
	if len(emailEntity.AppSource) == 0 {
		emailEntity.AppSource = constants.AppSourceCustomerOsApi
	}
	return &emailEntity
}

func MapEntitiesToEmails(entities *neo4jentity.EmailEntities) []*model.Email {
	if entities == nil {
		return nil
	}
	emails := make([]*model.Email, 0, len(*entities))
	for _, entity := range *entities {
		emails = append(emails, MapEntityToEmail(&entity))
	}
	return emails
}
