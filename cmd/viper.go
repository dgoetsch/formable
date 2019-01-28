package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgoetsch/formable/collections/stringarray"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initViper(cmd *cobra.Command) {
	viper.AutomaticEnv()
	viper.AddConfigPath(".")

	configDirectory, _ := cmd.Flags().GetString(configDirVar)
	if configDirectory == "" {
		configDirectory = viper.GetString(configDirVar)
	}
	if configDirectory != "" {
		viper.AddConfigPath(configDirectory)
	}

	viper.SetConfigName("formable")
	viper.ReadInConfig()
}

func copyViperFlagsToCmd(cmd *cobra.Command, args []string) error {
	initViper(cmd)
	for _, key := range viper.AllKeys() {
		viperValue := viper.GetString(key)
		if viperValue != "" && cmd.Flag(key).Value.String() == "" {
			cmd.Flags().Set(key, viperValue)
		}

	}
	return nil
}

const appName = "formable"

func configDirectory() string {
	return strings.TrimSuffix(viper.GetString(configDirVar), "/") + "/"
}

// func loadFile(filename string) {
// 	for configDirectory := range configDirectories() {
// 		for ext := range viper.SupportedExts {
// 			mergeConfig(configDirectory + filename + "." + ext)
// 		}
// 	}

// func configDirectories() []string {

// }

// func loadViperConfig() {

// 	configFileKeys := configFileKeys()
// 	for _, combination := range configFileKeys {
// 		configKey := strings.Join(combination, "-")
// 		loadFileType(configKey + ".json")
// 		loadFileType(configKey + ".yaml")
// 		loadFileType(configKey + ".yml")
// 	}
// }

// func loadFile(filename string) {
// 	for configDirectory := range configDirectories() {
// 		for ext := range viper.SupportedExts {
// 			mergeConfig(configDirectory + filename + "." + ext)
// 		}
// 	}

func configFileKeys() [][]string {
	project := viper.GetString(projectVar)
	serviceAccount := viper.GetString(serviceAccountVar)
	region := viper.GetString(regionVar)
	service := viper.GetString(serviceVar)
	keys := stringarray.Filter(
		[]string{project, serviceAccount, region, service},
		func(key string) bool {
			return key != ""
		})
	return allCombinations(keys, make([][]string, 0))

}
func mergeConfig(filename string) error {
	fmt.Printf("Loading %s config from %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	err = viper.MergeConfig(file)
	return err
}

func allCombinations(keys []string, sofar [][]string) [][]string {
	if len(keys) == 0 {
		return sofar
	}

	these := append(
		[][]string{[]string{keys[0]}},
		stringarray.Map(sofar, func(alreadyInSet []string) []string {
			return append(alreadyInSet, keys[0])
		})...)

	allOfTheseSoFar := append(sofar, these...)
	remainingKeys := keys[1:]

	return allCombinations(remainingKeys, allOfTheseSoFar)

}
