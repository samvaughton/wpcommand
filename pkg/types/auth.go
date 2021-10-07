package types

const RoleAdmin = "admin"
const RoleMember = "member"

const AuthObjectCasbin = "casbin"

const AuthObjectSite = "site"

const AuthObjectCommand = "command"
const AuthObjectCommandRunType = "command_run_type"

const AuthObjectCommandJob = "command_job"
const AuthObjectCommandJobEvent = "command_job_event"

const AuthObjectBlueprint = "blueprint"
const AuthObjectBlueprintObject = "blueprint_object"

const AuthObjectUser = "user"
const AuthObjectAccount = "account"
const AuthObjectConfig = "config"

const AuthObjectWordpressUser = "wp_user"

const AuthActionRead = "read"
const AuthActionReadSpecial = "read_special"
const AuthActionWrite = "write"
const AuthActionWriteSpecial = "write_special"
const AuthActionDelete = "delete"
const AuthActionRun = "run"                // things like deploy site etc
const AuthActionRunSpecial = "run_special" // things like setting up plugins/themes etc
const AuthActionConfigure = "configure"

type Authentication struct {
	Account  string
	Email    string
	Password string
}

type Token struct {
	Email       string
	TokenString string
}
