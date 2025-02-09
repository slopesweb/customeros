package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"gorm.io/gorm"
)

type Repositories struct {
	AiLocationMappingRepository                 AiLocationMappingRepository
	AiPromptLogRepository                       AiPromptLogRepository
	AppKeyRepository                            AppKeyRepository
	FlowRepository                              FlowRepository
	FlowSequenceRepository                      FlowSequenceRepository
	FlowSequenceStepRepository                  FlowSequenceStepRepository
	FlowSequenceContactRepository               FlowSequenceContactRepository
	FlowSequenceSenderRepository                FlowSequenceSenderRepository
	PersonalIntegrationRepository               PersonalIntegrationRepository
	PersonalEmailProviderRepository             PersonalEmailProviderRepository
	TenantWebhookApiKeyRepository               TenantWebhookApiKeyRepository
	TenantWebhookRepository                     TenantWebhookRepository
	SlackChannelRepository                      SlackChannelRepository
	PostmarkApiKeyRepository                    PostmarkApiKeyRepository
	GoogleServiceAccountKeyRepository           GoogleServiceAccountKeyRepository
	CurrencyRateRepository                      CurrencyRateRepository
	EventBufferRepository                       EventBufferRepository
	TableViewDefinitionRepository               TableViewDefinitionRepository
	TrackingAllowedOriginRepository             TrackingAllowedOriginRepository
	TechLimitRepository                         TechLimitRepository
	ExternalAppKeysRepository                   ExternalAppKeysRepository
	EnrichDetailsBetterContactRepository        EnrichDetailsBetterContactRepository
	EnrichDetailsScrapInRepository              EnrichDetailsScrapInRepository
	EnrichDetailsPrefilterTrackingRepository    EnrichDetailsPrefilterTrackingRepository
	EnrichDetailsTrackingRepository             EnrichDetailsTrackingRepository
	UserEmailImportPageTokenRepository          UserEmailImportStateRepository
	RawEmailRepository                          RawEmailRepository
	OAuthTokenRepository                        OAuthTokenRepository
	SlackSettingsRepository                     SlackSettingsRepository
	SlackChannelNotificationRepository          SlackChannelNotificationRepository
	ApiCacheRepository                          ApiCacheRepository
	WorkflowRepository                          WorkflowRepository
	IndustryMappingRepository                   IndustryMappingRepository
	TrackingRepository                          TrackingRepository
	TenantSettingsRepository                    TenantSettingsRepository
	TenantSettingsOpportunityStageRepository    TenantSettingsOpportunityStageRepository
	TenantSettingsMailboxRepository             TenantSettingsMailboxRepository
	TenantSettingsEmailExclusionRepository      TenantSettingsEmailExclusionRepository
	EmailLookupRepository                       EmailLookupRepository
	EmailTrackingRepository                     EmailTrackingRepository
	TenantRepository                            TenantRepository
	CacheIpDataRepository                       CacheIpDataRepository
	CacheIpHunterRepository                     CacheIpHunterRepository
	CacheEmailValidationRepository              CacheEmailValidationRepository
	CacheEmailValidationDomainRepository        CacheEmailValidationDomainRepository
	StatsApiCallsRepository                     StatsApiCallsRepository
	CosApiEnrichPersonTempResultRepository      CosApiEnrichPersonTempResultRepository
	OranizationWebsiteHostingPlatformRepository OrganizationWebsiteHostingPlatformRepository
	CustomerOsIdsRepository                     CustomerOsIdsRepository
}

func InitRepositories(db *gorm.DB) *Repositories {
	repositories := &Repositories{
		AiLocationMappingRepository:                 NewAiLocationMappingRepository(db),
		AiPromptLogRepository:                       NewAiPromptLogRepository(db),
		AppKeyRepository:                            NewAppKeyRepo(db),
		FlowRepository:                              NewFlowRepository(db),
		FlowSequenceRepository:                      NewFlowSequenceRepository(db),
		FlowSequenceStepRepository:                  NewFlowSequenceStepRepository(db),
		FlowSequenceContactRepository:               NewFlowSequenceContactRepository(db),
		FlowSequenceSenderRepository:                NewFlowSequenceSenderRepository(db),
		PersonalIntegrationRepository:               NewPersonalIntegrationsRepo(db),
		PersonalEmailProviderRepository:             NewPersonalEmailProviderRepository(db),
		TenantWebhookApiKeyRepository:               NewTenantWebhookApiKeyRepository(db),
		TenantWebhookRepository:                     NewTenantWebhookRepo(db),
		SlackChannelRepository:                      NewSlackChannelRepository(db),
		PostmarkApiKeyRepository:                    NewPostmarkApiKeyRepo(db),
		GoogleServiceAccountKeyRepository:           NewGoogleServiceAccountKeyRepository(db),
		CurrencyRateRepository:                      NewCurrencyRateRepository(db),
		EventBufferRepository:                       NewEventBufferRepository(db),
		TableViewDefinitionRepository:               NewTableViewDefinitionRepository(db),
		TrackingAllowedOriginRepository:             NewTrackingAllowedOriginRepository(db),
		TechLimitRepository:                         NewTechLimitRepository(db),
		ExternalAppKeysRepository:                   NewExternalAppKeysRepository(db),
		EnrichDetailsBetterContactRepository:        NewEnrichDetailsBetterContactRepository(db),
		EnrichDetailsScrapInRepository:              NewEnrichDetailsScrapInRepository(db),
		EnrichDetailsPrefilterTrackingRepository:    NewEnrichDetailsPrefilterTrackingRepository(db),
		EnrichDetailsTrackingRepository:             NewEnrichDetailsTrackingRepository(db),
		UserEmailImportPageTokenRepository:          NewUserEmailImportStateRepository(db),
		RawEmailRepository:                          NewRawEmailRepository(db),
		OAuthTokenRepository:                        NewOAuthTokenRepository(db),
		SlackSettingsRepository:                     NewSlackSettingsRepository(db),
		SlackChannelNotificationRepository:          NewSlackChannelNotificationRepository(db),
		ApiCacheRepository:                          NewApiCacheRepository(db),
		WorkflowRepository:                          NewWorkflowRepository(db),
		IndustryMappingRepository:                   NewIndustryMappingRepository(db),
		TrackingRepository:                          NewTrackingRepository(db),
		TenantSettingsRepository:                    NewTenantSettingsRepository(db),
		TenantSettingsOpportunityStageRepository:    NewTenantSettingsOpportunityStageRepository(db),
		TenantSettingsMailboxRepository:             NewTenantSettingsMailboxRepository(db),
		TenantSettingsEmailExclusionRepository:      NewEmailExclusionRepository(db),
		EmailLookupRepository:                       NewEmailLookupRepository(db),
		EmailTrackingRepository:                     NewEmailTrackingRepository(db),
		TenantRepository:                            NewTenantRepository(db),
		CacheIpDataRepository:                       NewCacheIpDataRepository(db),
		CacheIpHunterRepository:                     NewCacheIpHunterRepository(db),
		CacheEmailValidationRepository:              NewCacheEmailValidationRepository(db),
		CacheEmailValidationDomainRepository:        NewCacheEmailValidationDomainRepository(db),
		StatsApiCallsRepository:                     NewStatsApiCallsRepository(db),
		CosApiEnrichPersonTempResultRepository:      NewCosApiEnrichPersonTempResultRepository(db),
		OranizationWebsiteHostingPlatformRepository: NewOrganizationWebsiteHostingPlatformRepository(db),
		CustomerOsIdsRepository:                     NewCustomerOsIdsRepository(db),
	}

	return repositories
}

func (r *Repositories) Migration(db *gorm.DB) {

	//err = db.AutoMigrate(&entity.AppKey{})
	//if err != nil {
	//	panic(err)
	//}

	err := db.AutoMigrate(
		&entity.Tenant{},
		&entity.AiLocationMapping{},
		&entity.AiPromptLog{},
		&entity.PersonalIntegration{},
		&entity.PersonalEmailProvider{},
		&entity.TenantWebhookApiKey{},
		&entity.TenantWebhook{},
		&entity.SlackChannel{},
		&entity.PostmarkApiKey{},
		&entity.GoogleServiceAccountKey{},
		&entity.CurrencyRate{},
		&entity.EventBuffer{},
		&entity.TableViewDefinition{},
		&entity.TechLimit{},
		&entity.TrackingAllowedOrigin{},
		&entity.TenantSettingsEmailExclusion{},
		&entity.ExternalAppKeys{},
		&entity.EnrichDetailsBetterContact{},
		&entity.EnrichDetailsScrapIn{},
		&entity.UserEmailImportState{},
		&entity.UserEmailImportStateHistory{},
		&entity.RawEmail{},
		&entity.OAuthTokenEntity{},
		&entity.SlackSettingsEntity{},
		&entity.SlackChannelNotification{},
		&entity.ApiCache{},
		&entity.Workflow{},
		&entity.IndustryMapping{},
		&entity.Tracking{},
		&entity.EnrichDetailsPreFilterTracking{},
		&entity.EnrichDetailsTracking{},
		&entity.TenantSettingsOpportunityStage{},
		&entity.TenantSettingsMailbox{},
		&entity.FlowSequenceStepTemplateVariable{},
		&entity.Flow{},
		&entity.FlowSequence{},
		&entity.FlowSequenceStep{},
		&entity.FlowSequenceContact{},
		&entity.FlowSequenceSender{},
		&entity.EmailLookup{},
		&entity.EmailTracking{},
		&entity.CacheIpData{},
		&entity.CacheIpHunter{},
		&entity.CacheEmailValidation{},
		&entity.CacheEmailValidationDomain{},
		&entity.StatsApiCalls{},
		&entity.CosApiEnrichPersonTempResult{},
		&entity.OrganizationWebsiteHostingPlatform{},
		&entity.CustomerOsIds{})
	if err != nil {
		panic(err)
	}
}
