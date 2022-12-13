/*
 */
package cmd

import (
	"fmt"
	"os"
	"github.com/megaproaktiv/orphan/groups"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "orphan",
    Short: "Delete orphaned log groups from AWS Lambda",
    Long:  `Any log group named /aws/lambda/$name wich has no Lambda function $name will be deleted`,
    // Uncomment the following line if your bare application
    // has an action associated with it:
    // Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		orphaned,err := groups.ListOrphans(groups.ClientLambda, groups.ClientLogs)
		if err != nil{
			panic(err)
		}
		killthem, _ := cmd.Flags().GetBool("no-dry-run")
		if err != nil{
			panic(err)
		}
		for _, group := range orphaned {
			if killthem {
				groups.DeleteLogGroup(*groups.ClientLogs, group)
				fmt.Printf("Deleted: %v \n", *group)
			}else {
				fmt.Printf("%v \n", *group)
			}
		}
		if len(orphaned) == 0 {
			fmt.Printf("All clearn, no orphans found.\n")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("no-dry-run", "n", false, "Really delete log groups")
}


