package v1

import (
	"github.com/icinga/icingadb/pkg/contracts"
)

var Factories = []contracts.EntityFactoryFunc{
	NewActionUrl,
	NewCheckcommand,
	NewCheckcommandArgument,
	NewCheckcommandCustomvar,
	NewCheckcommandEnvvar,
	NewComment,
	NewCustomvar,
	NewDowntime,
	NewEndpoint,
	NewEventcommand,
	NewEventcommandArgument,
	NewEventcommandCustomvar,
	NewEventcommandEnvvar,
	NewHost,
	NewHostCustomvar,
	NewHostgroup,
	NewHostgroupCustomvar,
	NewHostgroupMember,
	NewIconImage,
	NewNotesUrl,
	NewNotification,
	NewNotificationcommand,
	NewNotificationcommandArgument,
	NewNotificationcommandCustomvar,
	NewNotificationcommandEnvvar,
	NewNotificationCustomvar,
	NewNotificationRecipient,
	NewNotificationUser,
	NewNotificationUsergroup,
	NewService,
	NewServiceCustomvar,
	NewServicegroup,
	NewServicegroupCustomvar,
	NewServicegroupMember,
	NewTimeperiod,
	NewTimeperiodCustomvar,
	NewTimeperiodOverrideExclude,
	NewTimeperiodOverrideInclude,
	NewTimeperiodRange,
	NewUser,
	NewUserCustomvar,
	NewUsergroup,
	NewUsergroupCustomvar,
	NewUsergroupMember,
	NewZone,
}
