package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type upOptions struct {
	Detach             bool
	Build              bool
	noBuild            bool
	Pull               string
	removeOrphans      bool
	scale              []string
	noColor            bool
	noPrefix           bool
	forceRecreate      bool
	noRecreate         bool
	noStart            bool
	cascadeStop        bool
	exitCodeFrom       string
	timeout            int
	timestamp          bool
	noDeps             bool
	attachDependencies bool
	wait               bool
	waitTimeout        int
	recreateDeps       bool
	noInherit          bool
	quietPull          bool
	attach             []string
	noAttach           []string
}

var upCommand = &cobra.Command{
	Use:   "up",
	Short: "Start go project using docker",
	Long:  "Start go project using docker.\n📔 Modify docker-compose.yml file for custom configuration.\n❗ Please make sure your device already installed docker to use this command",
	Run:   runUpCommand,
}

func init() {
	up := upOptions{}
	rootCommand.AddCommand(upCommand)
	flags := upCommand.Flags()
	flags.BoolVarP(&up.Detach, "detach", "d", false, "Detached mode: Run containers in the background")
	flags.BoolVar(&up.Build, "build", false, "Build images before starting containers.")
	flags.BoolVar(&up.noBuild, "no-build", false, "Don't build an image, even if it's missing.")
	flags.StringVar(&up.Pull, "pull", "missing", `Pull image before running ("always"|"missing"|"never")`)
	flags.BoolVar(&up.removeOrphans, "remove-orphans", false, "Remove containers for services not defined in the Compose file.")
	flags.StringArrayVar(&up.scale, "scale", []string{}, "Scale SERVICE to NUM instances. Overrides the `scale` setting in the Compose file if present.")
	flags.BoolVar(&up.noColor, "no-color", false, "Produce monochrome output.")
	flags.BoolVar(&up.noPrefix, "no-log-prefix", false, "Don't print prefix in logs.")
	flags.BoolVar(&up.forceRecreate, "force-recreate", false, "Recreate containers even if their configuration and image haven't changed.")
	flags.BoolVar(&up.noRecreate, "no-recreate", false, "If containers already exist, don't recreate them. Incompatible with --force-recreate.")
	flags.BoolVar(&up.noStart, "no-start", false, "Don't start the services after creating them.")
	flags.BoolVar(&up.cascadeStop, "abort-on-container-exit", false, "Stops all containers if any container was stopped. Incompatible with -d")
	flags.StringVar(&up.exitCodeFrom, "exit-code-from", "", "Return the exit code of the selected service container. Implies --abort-on-container-exit")
	flags.IntVarP(&up.timeout, "timeout", "t", 10, "Use this timeout in seconds for container shutdown when attached or when containers are already running.")
	flags.BoolVar(&up.timestamp, "timestamps", false, "Show timestamps.")
	flags.BoolVar(&up.noDeps, "no-deps", false, "Don't start linked services.")
	flags.BoolVar(&up.recreateDeps, "always-recreate-deps", false, "Recreate dependent containers. Incompatible with --no-recreate.")
	flags.BoolVarP(&up.noInherit, "renew-anon-volumes", "V", false, "Recreate anonymous volumes instead of retrieving data from the previous containers.")
	flags.BoolVar(&up.attachDependencies, "attach-dependencies", false, "Attach to dependent containers.")
	flags.BoolVar(&up.quietPull, "quiet-pull", false, "Pull without printing progress information.")
	flags.StringArrayVar(&up.attach, "attach", []string{}, "Attach to service output.")
	flags.StringArrayVar(&up.noAttach, "no-attach", []string{}, "Don't attach to specified service.")
	flags.BoolVar(&up.wait, "wait", false, "Wait for services to be running|healthy. Implies detached mode.")
	flags.IntVar(&up.waitTimeout, "wait-timeout", 0, "timeout waiting for application to be running|healthy.")
}

func runUpCommand(cmd *cobra.Command, args []string) {
	// Get raw flags from cmd
	var flags []string

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Changed {
			flags = append(flags, fmt.Sprintf("--%s", flag.Name))
		}
	})
	// Merge args and flags
	args = append(args, flags...)
	newArgs := []string{"up"}
	newArgs = append(newArgs, args...)
	// Run docker-compose up
	// Create command
	exec := exec.Command("docker-compose", newArgs...)

	// Set output to stdout
	exec.Stdout = os.Stdout
	exec.Stderr = os.Stderr

	// Run command
	err := exec.Run()
	if err != nil {
		log.Fatal(err)
	}
}
