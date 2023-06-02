// Package pkg implements list function and variable that can be used by other packages
package pkg

import "github.com/AlecAivazis/survey/v2"

const LatestGoVersion = "1.20"

var SupportedGoVersions = []string{"1.20", "1.19", "1.18"}

// CreateCommandAnswer is the answer from the create command
type CreateCommandAnswer struct {
	HttpFramework string `survey:"http_framework"`
	Database      string `survey:"database"`
	UseOrm        bool   `survey:"use_orm"`
}

// CreateSurveyQuestion is the list of questions that will be asked to the user
var CreateSurveyQuestion = []*survey.Question{
	{
		Name: "http_framework",
		Prompt: &survey.Select{
			Message: "Choose a framework:",
			Options: []string{
				// "chi",
				// "echo",
				// "fiber",
				"gin",
				// "mux",
			},
			Default:  "gin",
			PageSize: 10,
		},
		Validate: survey.Required,
	},
	{
		Name: "database",
		Prompt: &survey.Select{
			Message: "Choose database",
			Options: []string{
				"mysql",
				"postgresql",
				// "sqlite",
				// "sqlserver",
			},
			Default:  "mysql",
			PageSize: 10,
		},
		Validate: survey.Required,
	},
	{
		Name: "use_orm",
		Prompt: &survey.Confirm{
			Message: "Do you want to use GORM?",
			Default: false,
		},
	},
}
