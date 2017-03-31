package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"
	"torch/templates"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newCmd = &cobra.Command{
	Use:   "new app_path [-p prefix]",
	Short: "It create a torch game server application.",
	Long:  `It create the torch application in the $GOAPTH/src directory.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			if path.IsAbs(args[0]) {
				fmt.Println("app_path:", args[0], " should be a relative path.")
			} else {
				newApp(args[0])
			}

		} else {
			fmt.Println("No value provided for required arguments 'app_path'")
			fmt.Println(cmd.Use)
		}

	},
}

var prefix string

func init() {
	newCmd.Flags().StringVarP(&prefix, "Game directory's prefix", "p", "", "create the game's directories with prefix")
}

func newApp(appPath string) {

	gopath := viper.Get("gopath").(string)
	if len(gopath) == 0 {
		fmt.Println("Can not create the application without $GOPATH")
	}

	gopathSrc := path.Join(gopath, "src")

	createDirectories(appPath, gopathSrc)
	generateServerCmd(appPath, gopathSrc)
}

func createDirectories(app_path string, gopathSrc string) {
	dirs := []string{
		"bin",
		"cmd/server/cmd",
		"game",
		"vendor",
		"protobuf",
		"docs",
		"db",
		"includes",
	}
	subDirs := map[string][]string{
		"game": []string{"handlers", "pstructs", "scenes", "structs", "global", "middlewares"},
	}

	for _, dir := range dirs {

		appDir := path.Join(gopathSrc, app_path, dir)

		if err := os.MkdirAll(appDir, 0755); err != nil {
			fmt.Println("Failed to create directory: ", appDir, err.Error())
		} else {
			fmt.Println("Create directory: ", appDir)
		}
		list := subDirs[dir]
		if list != nil {
			for _, subDir := range list {
				pSubDir := path.Join(appDir, subDir)

				if len(prefix) > 0 {

					pSubDir = path.Join(appDir, fmt.Sprintf("%s_%s", prefix, subDir))
				}

				if err := os.MkdirAll(pSubDir, 0755); err != nil {
					fmt.Println("Failed to create directory: ", pSubDir, err.Error())

				} else {
					fmt.Println("Create directory: ", pSubDir)
				}
			}
		}
	}

}

func generateServerCmd(appPath string, gopathSrc string) {
	basename := path.Base(appPath)
	fpath := path.Join(gopathSrc, appPath, "cmd", "server", basename+".go")
	// bName := path.Base(appPath)
	writeToFile(fpath, func(f *os.File) { templates.WriteGenerateMain(f, appPath) }, 0644)

	fpath = path.Join(gopathSrc, appPath, "cmd", "server", "cmd", "server.go")

	writeToFile(fpath, func(f *os.File) { templates.WriteGenerateServerCmd(f, prefix, appPath, basename) }, 0644)

	fpath = path.Join(gopathSrc, appPath, "cmd", "server", "cmd", "root.go")

	short := "The game server of " + appPath
	long := `The command line tool for game server`

	writeToFile(fpath, func(f *os.File) { templates.WriteGenerateRootCmd(f, appPath, short, long) }, 0644)

	globalFilename := "global.go"
	globalDir := "global"
	bindingHandersFilename := "binding_handlers.go"
	bindingHandersDir := "handlers"
	structsDir := "structs"
	pstructsDir := "pstructs"
	scenesDir := "scenes"

	if len(prefix) > 0 {
		globalFilename = fmt.Sprintf("%s_%s", prefix, globalFilename)
		globalDir = fmt.Sprintf("%s_%s", prefix, globalDir)
		bindingHandersFilename = fmt.Sprintf("%s_%s", prefix, bindingHandersFilename)
		bindingHandersDir = fmt.Sprintf("%s_%s", prefix, bindingHandersDir)
		structsDir = fmt.Sprintf("%s_%s", prefix, structsDir)
		pstructsDir = fmt.Sprintf("%s_%s", prefix, pstructsDir)
		scenesDir = fmt.Sprintf("%s_%s", prefix, scenesDir)
	}

	fpath = path.Join(gopathSrc, appPath, "cmd", "server", "cmd", "generator.go")
	writeToFile(fpath, func(f *os.File) {
		templates.WriteGenerateGeneratorCmd(f, bindingHandersDir, appPath, scenesDir, prefix)
	}, 0644)

	fpath = path.Join(gopathSrc, appPath, "game", globalDir, globalFilename)
	writeToFile(fpath, func(f *os.File) { templates.WriteGenerateGlobal(f, prefix) }, 0644)

	fpath = path.Join(gopathSrc, appPath, "game", bindingHandersDir, bindingHandersFilename)

	writeToFile(fpath, func(f *os.File) { templates.WriteGenerateBindingHandlers(f, prefix, appPath) }, 0644)

	fpath = path.Join(gopathSrc, appPath, "game", bindingHandersDir, "example_handler.go")

	writeToFile(fpath, func(f *os.File) { templates.WriteGenerateExampleHandler(f, prefix, appPath, "Example") }, 0644)

	protobufDir := path.Join(gopathSrc, appPath, "protobuf")
	fpath = path.Join(protobufDir, "request.proto")

	writeToFile(fpath, func(f *os.File) { templates.WriteGenerateExampleRequestProtobuf(f, prefix, appPath) }, 0644)

	fpath = path.Join(protobufDir, "common.proto")

	writeToFile(fpath, func(f *os.File) { templates.WriteGenerateCommonProtobuf(f, prefix, appPath) }, 0644)

	fpath = path.Join(protobufDir, "response.proto")

	writeToFile(fpath, func(f *os.File) { templates.WriteGenerateResponseProtobuf(f, prefix, appPath) }, 0644)

	fpath = path.Join(gopathSrc, appPath, "game", structsDir, "example.go")

	writeToFile(fpath, func(f *os.File) { templates.WriteGenerateExampleStruct(f, prefix, appPath) }, 0644)

	fpath = path.Join(gopathSrc, appPath, "game", structsDir, "protobuf_id_map.go")

	writeToFile(fpath, func(f *os.File) { templates.WriteGenerateIdMapStruct(f, prefix, appPath) }, 0644)

	fpath = path.Join(gopathSrc, appPath, "bin", fmt.Sprintf("%s.sh", appPath))

	cmd := fmt.Sprintf("go run %s $@", path.Join("${GOPATH}", "src", appPath, "cmd", "server", appPath+".go"))

	writeToFile(fpath, func(f *os.File) { templates.WriteGenerateServerScript(f, cmd) }, 0755)

	fpath = path.Join(gopathSrc, appPath, "bin", "generate_protobuf_go_files.sh")

	writeToFile(fpath, func(f *os.File) {
		templates.WriteGenerateProtobufGofile(f,
			path.Join("${GOPATH}", path.Join("src", appPath, "game", pstructsDir)),
			protobufDir,
			path.Join(protobufDir, "*.proto"))
	}, 0755)

	generateProtobufGoFiles(path.Join(gopathSrc, appPath, "game", pstructsDir), protobufDir)

}

type templateWrite func(*os.File)

var overwriteFlag string

func writeToFile(fpath string, fun templateWrite, mode os.FileMode) {
Here:
	switch overwriteFlag {
	case "Y":
	case "N":
		if _, err := os.Stat(fpath); !os.IsNotExist(err) {
			return
		}
	default:
		if _, err := os.Stat(fpath); !os.IsNotExist(err) {
			fmt.Println("The file path: ", fpath, " is exist. Do you want to overwrite it? (n/N/y/Y)")
			if _, err := fmt.Scanf("%s", &overwriteFlag); err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			switch overwriteFlag {
			case "n", "N":
				return
			case "y", "Y":
			default:
				goto Here
			}
		}
	}

	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)

	defer f.Close()
	if err != nil {
		fmt.Println("Failed to generate the file:", fpath)
	} else {
		fmt.Println("Generate file:", fpath)
		fun(f)
	}

}

func generateProtobufGoFiles(pstructsPath string, protobufPath string) {
	var wg sync.WaitGroup
	fPath := path.Join(protobufPath, "*.proto")
	cmd := fmt.Sprintf("protoc --gofast_out=%s --proto_path=%s %s", pstructsPath, protobufPath, fPath)
	execCmd(cmd, &wg)
	wg.Wait()
}

func execCmd(cmd string, wg *sync.WaitGroup) {
	wg.Add(1)
	fmt.Println("command is ", cmd)

	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("%s\n", out)

	wg.Done() // Need to signal to waitgroup that this goroutine is done
}
