package types

type PluginAction string

type PluginActionEnums struct {
	None      PluginAction
	Upgrade   PluginAction
	Downgrade PluginAction
	Install   PluginAction
	Uninstall PluginAction
}

var PluginActionEnum = &PluginActionEnums{
	None:      "none",
	Upgrade:   "upgrade",
	Downgrade: "downgrade",
	Install:   "install",
	Uninstall: "uninstall",
}

var PluginActionsList = []PluginAction{
	PluginActionEnum.None,
	PluginActionEnum.Upgrade,
	PluginActionEnum.Downgrade,
	PluginActionEnum.Install,
	PluginActionEnum.Uninstall,
}

type PluginActionSet struct {
	Actions []PluginActionItem
}

type PluginActionItem struct {
	Object ObjectBlueprint
	Action PluginAction
}
