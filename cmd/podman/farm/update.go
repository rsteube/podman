package farm

import (
	"errors"
	"fmt"

	"github.com/containers/common/pkg/completion"
	"github.com/containers/common/pkg/config"
	"github.com/containers/podman/v4/cmd/podman/common"
	"github.com/containers/podman/v4/cmd/podman/registry"
	"github.com/containers/podman/v4/pkg/util"
	"github.com/spf13/cobra"
)

var (
	farmUpdateDescription = `Update an existing farm by adding a connection, removing a connection, or changing it to the default farm.`
	updateCommand         = &cobra.Command{
		Use:               "update [options] FARM",
		Short:             "Update an existing farm",
		Long:              farmUpdateDescription,
		RunE:              farmUpdate,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: common.AutoCompleteFarms,
		Example: `podman farm update --add con1 farm1
	podman update --remove con2 farm2
	podman update --default farm3`,
	}

	// Temporary struct to hold cli values.
	updateOpts = struct {
		Add     []string
		Remove  []string
		Default bool
	}{}
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Command: updateCommand,
		Parent:  farmCmd,
	})
	flags := updateCommand.Flags()

	addFlagName := "add"
	flags.StringSliceVarP(&updateOpts.Add, addFlagName, "a", nil, "add system connection(s) to farm")
	_ = updateCommand.RegisterFlagCompletionFunc(addFlagName, completion.AutocompleteNone)
	removeFlagName := "remove"
	flags.StringSliceVarP(&updateOpts.Remove, removeFlagName, "r", nil, "remove system connection(s) from farm")
	_ = updateCommand.RegisterFlagCompletionFunc(removeFlagName, completion.AutocompleteNone)
	defaultFlagName := "default"
	flags.BoolVarP(&updateOpts.Default, defaultFlagName, "d", false, "set the given farm as the default farm")
}

func farmUpdate(cmd *cobra.Command, args []string) error {
	farmName := args[0]

	defChanged := cmd.Flags().Changed("default")

	if len(updateOpts.Add) == 0 && len(updateOpts.Remove) == 0 && !defChanged {
		return fmt.Errorf("nothing to update for farm %q, please use the --add, --remove, or --default flags to update a farm", farmName)
	}

	cfg, err := config.ReadCustomConfig()
	if err != nil {
		return err
	}

	if len(cfg.Farms.List) == 0 {
		return errors.New("no farms are created at this time, there is nothing to update")
	}

	if _, ok := cfg.Farms.List[farmName]; !ok {
		return fmt.Errorf("cannot update farm, %q farm doesn't exist", farmName)
	}

	if defChanged {
		// Change the default to the given farm if --default=true
		if updateOpts.Default {
			cfg.Farms.Default = farmName
		} else {
			// if --default=false, user doesn't want any farms to be default so clear the DefaultFarm
			cfg.Farms.Default = ""
		}
	}

	if val, ok := cfg.Farms.List[farmName]; ok {
		cMap := make(map[string]int)
		for _, c := range val {
			cMap[c] = 0
		}

		for _, cRemove := range updateOpts.Remove {
			connections := cfg.Farms.List[farmName]
			if util.StringInSlice(cRemove, connections) {
				delete(cMap, cRemove)
			} else {
				return fmt.Errorf("cannot remove from farm, %q is not a connection in the farm", cRemove)
			}
		}

		for _, cAdd := range updateOpts.Add {
			if _, ok := cfg.Engine.ServiceDestinations[cAdd]; ok {
				if _, ok := cMap[cAdd]; !ok {
					cMap[cAdd] = 0
				}
			} else {
				return fmt.Errorf("cannot add to farm, %q is not a system connection", cAdd)
			}
		}

		updatedConnections := []string{}
		for k := range cMap {
			updatedConnections = append(updatedConnections, k)
		}
		cfg.Farms.List[farmName] = updatedConnections
	}

	if err := cfg.Write(); err != nil {
		return err
	}
	fmt.Printf("Farm %q updated\n", farmName)
	return nil
}
