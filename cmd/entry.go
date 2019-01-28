package cmd

import (
	"bytes"
	"log"

	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const tfCmdVar string = "tf_cmd"
const tfDirVar string = "tf_dir"
const configDirVar string = "config_dir"
const inputVar string = "input"
const projectVar string = "project"
const serviceAccountVar string = "service_account"
const serviceVar string = "service"
const regionVar string = "region"
const zoneVar string = "zone"

var rootCmd = &cobra.Command{
	Use:   "formable",
	Short: "terraform execution infrastucture on google cloud",
	Long:  `terraform execution infrastucture in a docker container for auth and state and other ubiquitous things on google cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		tfCmd, _ := cmd.Flags().GetString(tfCmdVar)
		terraformDir, _ := cmd.Flags().GetString(tfDirVar)
		configDir, _ := cmd.Flags().GetString(configDirVar)
		promptForInput, _ := cmd.Flags().GetBool(inputVar)
		project, _ := cmd.Flags().GetString(projectVar)
		serviceAccount, _ := cmd.Flags().GetString(serviceAccountVar)
		service, _ := cmd.Flags().GetString(serviceVar)
		zone, _ := cmd.Flags().GetString(zoneVar)
		region, _ := cmd.Flags().GetString(regionVar)

		fmt.Printf("Initialilzed formable for project %s in account %s\n\tregion: %s\n\tzone: %s\n\tservice: %s\n", project, serviceAccount, region, zone, service)

		varFilesArgs, err := varFilesArgs(configDir)
		if err != nil {
			log.Fatalf("Could not load var files args from %s because of an error", configDir, err)
		}

		args = append(append(varFilesArgs, inputArg(promptForInput)), terraformDir)
		initArgs := append([]string{"init"}, append(backendConfig(service, project, region), args...)...)
		err = handleCmd(exec.Command("terraform", initArgs...))
		if err != nil {
			fmt.Printf("Error: %s", err)
			panic(err)
		}
		cmdArgs := append([]string{tfCmd}, args...)
		err = handleCmd(exec.Command("terraform", cmdArgs...))
		if err != nil {
			fmt.Printf("Error: %s", err)
			panic(err)
		}
	},
}

func backendConfig(service, project, region string) []string {
	return []string{
		fmt.Sprintf("-backend-config=bucket=%s-terraform-state", project),
		fmt.Sprintf("-backend-config=project=%s", project),
		fmt.Sprintf("-backend-config=region=%s", region),
		fmt.Sprintf("-backend-config=prefix=%s/", service),
	}
}

func handleCmd(cmd *exec.Cmd) error {
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	fmt.Printf("Executing %s, \n", cmd.Args)
	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()
	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err := cmd.Wait()
	return err
}

func inputArg(promptForInput bool) string {
	return fmt.Sprintf("-input=%t", promptForInput)
}

func varFilesArgs(configDir string) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(configDir)
	if err != nil {
		return make([]string, 0), err
	}

	var result = make([]string, 0)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			subDirVarFiles, err := varFilesArgs(configDir + "/" + fileInfo.Name())
			if err != nil {
				return result, err
			}
			result = append(result, subDirVarFiles...)
		} else {
			result = append(result, fmt.Sprintf("-var-file=%s/%s", configDir, fileInfo.Name()))
		}
	}

	return result, nil
}

func EntryPoint() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	bindRequiredStringFlag(rootCmd, tfCmdVar, "f", "", "the directory that contains terraform files")
	bindRequiredStringFlag(rootCmd, tfDirVar, "t", "", "the directory that contains terraform files")
	bindRequiredStringFlag(rootCmd, configDirVar, "c", "", "the directory that contains terraform configuration")
	bindRequiredStringFlag(rootCmd, serviceVar, "s", "", "the service name")
	bindRequiredStringFlag(rootCmd, projectVar, "p", "", "the name of the gc project")
	bindRequiredStringFlag(rootCmd, serviceAccountVar, "a", "", "the gc service account name")
	bindRequiredStringFlag(rootCmd, regionVar, "r", "", "the gc region to run against")
	bindStringFlag(rootCmd, zoneVar, "z", "", "the gc zone to run against")

	rootCmd.Flags().Bool(inputVar, false, "prompt for input on destructive commands")
	viper.BindPFlag(inputVar, rootCmd.Flags().Lookup(inputVar))
	rootCmd.PreRunE = copyViperFlagsToCmd
}

func bindRequiredStringFlag(cmd *cobra.Command, name string, shorthand string, defaultValue string, description string) {
	cmd.Flags().StringP(name, shorthand, defaultValue, description)
	viper.BindPFlag(name, cmd.Flags().Lookup(name))
	cmd.MarkFlagRequired(name)
}

func bindStringFlag(cmd *cobra.Command, name string, shorthand string, defaultValue string, description string) {
	cmd.Flags().StringP(name, shorthand, defaultValue, description)
	viper.BindPFlag(name, cmd.Flags().Lookup(name))
}
