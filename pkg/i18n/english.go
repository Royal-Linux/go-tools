package i18n

// TranslationSet is a set of localised strings for a given language
type TranslationSet struct {
	ProjectTitle                               string
	MainTitle                                  string
	GlobalTitle                                string
	Navigate                                   string
	Menu                                       string
	Execute                                    string
	Scroll                                     string
	Close                                      string
	ErrorTitle                                 string
	EditConfig                                 string
	AnonymousReportingTitle                    string
	AnonymousReportingPrompt                   string
	ConfirmQuit                                string
	ErrorOccurred                              string
	ConnectionFailed                           string

	Donate                     string
	Cancel                     string
	Remove                     string
	HideStopped                string
	ForceRemove                string
	Confirm                    string
	Return                     string
	FocusMain                  string
	RestartingStatus           string
	StoppingStatus             string
	RemovingStatus             string
	Stop                       string
	Restart                    string
	Rebuild                    string
	Recreate                   string
	PreviousContext            string
	NextContext                string
	Attach                     string
	ViewLogs                   string
	TopTitle                   string
	PressEnterToReturn         string
	ExecShell                  string

	LogsTitle                string
	ConfigTitle              string
	StatsTitle               string
	CreditsTitle             string

	No  string
	Yes string
}

func englishSet() TranslationSet {
	return TranslationSet{
		RemovingStatus:             "removing",
		RestartingStatus:           "restarting",
		StoppingStatus:             "stopping",

		ErrorOccurred:                     "An error occurred! Please create an issue at https://github.com/Royal-Linux/hornero/issues",
		ConnectionFailed:                  "connection failed. You may need to restart the client",

		Donate:  "Donate",
		Confirm: "Confirm",

		Return:              "return",
		FocusMain:           "focus main panel",
		Navigate:            "navigate",
		Execute:             "execute",
		Close:               "close",
		Menu:                "menu",
		Scroll:              "scroll",
		OpenConfig:          "open hornero config",
		EditConfig:          "edit hornero config",
		Cancel:              "cancel",
		Remove:              "remove",
		HideStopped:         "Hide/Show stopped containers",
		ForceRemove:         "force remove",
		Stop:                "stop",
		Restart:             "restart",
		Rebuild:             "rebuild",
		Recreate:            "recreate",
		PreviousContext:     "previous tab",
		NextContext:         "next tab",
		Attach:              "attach",
		ViewLogs:            "view logs",
		ExecShell:           "exec shell",

		AnonymousReportingTitle:  "Help make hornero better",
		AnonymousReportingPrompt: "Would you like to enable anonymous reporting data to help improve hornero?",

		GlobalTitle:               "Global",
		MainTitle:                 "Main",
		ProjectTitle:              "Project",
		ErrorTitle:                "Error",
		LogsTitle:                 "Logs",
		ConfigTitle:               "Config",
		TopTitle:                  "Top",
		StatsTitle:                "Stats",
		CreditsTitle:              "About",

		ConfirmQuit:                "Are you sure you want to quit?",
		PressEnterToReturn:         "Press enter to return to hornero (this prompt can be disabled in your config by setting `gui.returnImmediately: true`)",

		No:  "no",
		Yes: "yes",
	}
}
