package pkg

const LatestGoVersion = "1.20"

var SupportedGoVersions = []string{"1.20", "1.19", "1.18"}

type CreateCommandAnswer struct {
	HttpFramework string `survey:"http_framework"`
	Database      string `survey:"database"`
}
