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

type PluginActionItem struct {
	Name   string           // although object has the exact name, if the object does not exist, this needs to be set regardless
	Object *ObjectBlueprint `json:"-"`
	Action PluginAction
}

type PluginActionSet struct {
	Items []PluginActionItem
}
