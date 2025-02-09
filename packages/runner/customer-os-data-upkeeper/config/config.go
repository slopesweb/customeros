package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	cronconf "github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/cron/config"
	commconf "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"log"
)

type Config struct {
	Postgres         commconf.PostgresConfig
	Neo4j            commconf.Neo4jConfig
	Logger           logger.Config
	Jaeger           tracing.JaegerConfig
	Cron             cronconf.Config
	GrpcClientConfig commconf.GrpcClientConfig
	CustomerOS       struct {
		CustomerOsAPI    string `env:"CUSTOMER_OS_API,required"`
		CustomerOsAPIKey string `env:"CUSTOMER_OS_API_KEY,required"`
	}
	PlatformAdminApi struct {
		Url    string `env:"PLATFORM_ADMIN_API_URL"`
		ApiKey string `env:"PLATFORM_ADMIN_API_KEY"`
	}
	EnrichmentApi struct {
		Url    string `env:"ENRICHMENT_API_URL" validate:"required"`
		ApiKey string `env:"ENRICHMENT_API_KEY" validate:"required"`
	}
	BetterContactApi struct {
		Url    string `env:"BETTER_CONTACT_API_URL" validate:"required"`
		ApiKey string `env:"BETTER_CONTACT_API_KEY" validate:"required"`
	}
	ProcessConfig      ProcessConfig
	EventNotifications EventNotifications
}

type ProcessConfig struct {
	CycleInvoicingEnabled                    bool `env:"CYCLE_INVOICING_ENABLED" envDefault:"true"`
	OffCycleInvoicingEnabled                 bool `env:"OFF_CYCLE_INVOICING_ENABLED" envDefault:"false"`
	DelaySendPayInvoiceNotificationInMinutes int  `env:"DELAY_SEND_PAY_INVOICE_NOTIFICATION_IN_MINUTES" envDefault:"60"`
	RetrySendPayInvoiceNotificationDays      int  `env:"RETRY_SEND_PAY_INVOICE_NOTIFICATION_DAYS" envDefault:"5"`
	DelayAutoPayInvoiceInMinutes             int  `env:"DELAY_AUTO_PAY_INVOICE_IN_MINUTES" envDefault:"5"`
	DelayRequestPaymentLinkInMinutes         int  `env:"DELAY_REQUEST_PAYMENT_LINK_IN_MINUTES" envDefault:"15"`
	RequestPaymentLinkLookBackWindowInDays   int  `env:"REQUEST_PAYMENT_LINK_LOOK_BACK_WINDOW_IN_DAYS" envDefault:"5"`
	DelayGenerateCycleInvoiceInMinutes       int  `env:"DELAY_GENERATE_CYCLE_INVOICE_IN_MINUTES" envDefault:"240"`
	DelayGenerateOffCycleInvoiceInMinutes    int  `env:"DELAY_GENERATE_OFF_CYCLE_INVOICE_IN_MINUTES" envDefault:"60"`
}

type EventNotifications struct {
	IntegrationAppEventWebhookUrls struct {
		GeneratePaymentLinkUrl string `env:"INVOICE_GENERATE_PAYMENT_LINK_URL" envDefault:"" required:"true"`
		InvoiceFinalizedUrl    string `env:"INVOICE_READY_URL" envDefault:"" required:"true"`
	}
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("Failed loading .env file")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v", err)
	}

	return &cfg
}
