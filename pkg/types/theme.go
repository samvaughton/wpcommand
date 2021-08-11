package types

type ThemeAction string

type ThemeActionEnums struct {
	None      ThemeAction
	Upgrade   ThemeAction
	Downgrade ThemeAction
	Install   ThemeAction
	Uninstall ThemeAction
}

var ThemeActionEnum = &ThemeActionEnums{
	None:      "none",
	Upgrade:   "upgrade",
	Downgrade: "downgrade",
	Install:   "install",
	Uninstall: "uninstall",
}

var ThemeActionsList = []ThemeAction{
	ThemeActionEnum.None,
	ThemeActionEnum.Upgrade,
	ThemeActionEnum.Downgrade,
	ThemeActionEnum.Install,
	ThemeActionEnum.Uninstall,
}

type ThemeActionSet struct {
	Items []ThemeActionItem
}

type ThemeActionItem struct {
	Name   string
	Object *ObjectBlueprint
	Action ThemeAction
}
