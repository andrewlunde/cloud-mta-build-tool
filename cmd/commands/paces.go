package commands

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"cloud-mta-build-tool/cmd/constants"
	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/cmd/mta/metainfo"
	"cloud-mta-build-tool/cmd/proc"
)

// Prepare the process for execution
var prepare = &cobra.Command{
	Use:   "prepare",
	Short: "Prepare The environment For Build Process",
	Long:  "Prepare The environment For Build Process",
	Run: func(cmd *cobra.Command, args []string) {
		proc.Prepare()
	},
}

// Copy specific module for building purpose
var copyModule = &cobra.Command{
	Use:   "copy",
	Short: "copy module for build process",
	Long:  "copy module for build process",
	Run: func(cmd *cobra.Command, args []string) {
		logs.Logger.Infof("Executing Copy module %v", args[2])
		proc.CopyModule(args[0], args[1])
	},
}

// Zip specific module
var pack = &cobra.Command{
	Use:   "pack",
	Short: "Pack the module to zip format",
	Long:  "Pack the module to zip format",
	Run: func(cmd *cobra.Command, args []string) {
		// Define arguments variables
		if len(args) > 0 {
			tDir := args[0]
			mName := args[2]
			modRelPath := fs.ProjectPath() + "/" + args[1]
			modRelName := filepath.Join(tDir, mName)
			// Create empty folder with name as before the zip process
			// to put the file such as data.zip inside
			os.MkdirAll(modRelName, os.ModePerm)
			// zipping the build artifacts
			logs.Logger.Infof("Starting execute zipping module %v ", mName)
			if err := fs.Archive(modRelPath, tDir+"/"+args[1]+constants.DataZip, modRelPath); err != nil {
				logs.Logger.Error("Error occurred during ZIP module %v creation, error:   ", args[0], err)
			}
			logs.Logger.Infof("Execute zipping module %v finished successfully ", mName)
		} else {
			logs.Logger.Errorf("No path's provided to pack the module artifacts")
		}
	},
}

// Generate metadata info from deployment
var genMeta = &cobra.Command{
	Use:   "meta",
	Short: "Generate meta folder",
	Long:  "Generate meta folder",
	Run: func(cmd *cobra.Command, args []string) {
		logs.Logger.Info("Starting execute metadata creation")
		mtaStruct := proc.GetMta(fs.GetPath())
		mtarDir := args[0]
		// Generate meta info dir with required content
		metainfo.GenMetaInf(mtarDir, mtaStruct, args[1:])
		logs.Logger.Info("Metadata creation finish successfully")

	},
}

// Generate mtar from build artifacts
var genMtar = &cobra.Command{
	Use:   "mtar",
	Short: "Generate mtar",
	Long:  "Generate mtar",
	Run: func(cmd *cobra.Command, args []string) {
		logs.Logger.Info("Starting execute Build of mtar")
		mtaStruct := proc.GetMta(fs.GetPath())
		tDir := args[0]
		pDir := args[1]
		// Create MTAR from the building artifacts
		fs.Archive(tDir, pDir+constants.PathSep+mtaStruct.Id+constants.MtarSuffix, tDir)
		//logs.Logger.Infof("Build of mtar finished successfully, mtar location:  ", pDir+constants.PathSep+mtaStruct.Id+constants.MtarSuffix,tDir)
	},
}

// Cleanup temp artifacts
var cleanup = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove build temporary folder",
	Long:  "Remove build temporary folder",
	Run: func(cmd *cobra.Command, args []string) {
		logs.Logger.Info("Starting Cleanup process")
		// Remove temp folder
		os.RemoveAll(args[0])
		logs.Logger.Info("Done")
	},
}
