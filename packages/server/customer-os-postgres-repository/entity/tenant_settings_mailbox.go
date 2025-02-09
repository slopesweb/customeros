package entity

import "time"

type TenantSettingsMailbox struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt"`

	Username string `gorm:"column:user_name;type:varchar(255)" json:"userName"` // holds the email address of the user in the neo4j

	MailboxUsername string `gorm:"column:mailbox_username;type:varchar(255)" json:"mailboxUsername"`
	MailboxPassword string `gorm:"column:mailbox_password;type:varchar(255)" json:"mailboxPassword"`
}

func (TenantSettingsMailbox) TableName() string {
	return "tenant_settings_mailbox"
}
