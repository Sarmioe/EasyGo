package main

import (
	"archive/zip"
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	version = "0.0.2"
	hdv     = "0.0.4"
)
const hden = `
Using -h only English version ; If you using the others languages , Using ezgo -hes (es is for Espanol , mean is Sypanish)and more.
This is some basic commands.
| command                                                      | Function                                                     |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| ezgo -v                                                       | Display version                                              |
| ezgo -h                                                       | Print Help default English                                   |
| ezgo -hzc                                                     | Print Help of Chinese                                        |
| ezgo -hzt                                                     | Print Help of Chinese Traditional                            |
| ezgo -hes                                                     | Print Help of Spanish                                        |
| ezgo -update [version]                                        | Update EasyGo                                                |
| ezgo -clone [URL] [Localpath] -branch--[branchname] -depth--[number] | Clone repo from cloud                                        |
| ezgo -sync [localpath] [URL]                                  | Run sync                                                     |
| ezgo -sync auto [time defualt is second]                      | Auto sync                                                    |
| ezgo -sync incremental                                        | Synchronize only difference files                            |
| ezgo -config                                                  | Configure EasyGo                                             |
| ezgo -env | Automatic environment check |
| ezgo -logs [level] | Output ezgo logs |
| ezgo -logs git | Output git logs |
| ezgo -logs go | Output Go logs |
| ezgo -push [commit] | Commit to remote repository |
| ezgo -pull [branch] | Pull a branch |
| ezgo -checkout [branchname] | Switch branch name |
| ezgo -conflict [way] | Resolve cloud and local conflicts |
If you want see more , Please view this page :https://github.com/Sarmioe/EasyGo/blob/main/README.md
`
const hdes = `
Utilice -hzc para mostrar únicamente la versión en chino. Si utiliza otros idiomas, utilice -hes, por ejemplo, es es la abreviatura de español.
Aquí hay algunos comandos comunes
| comando                                                      | Función                                                      |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| ezgo -v       | Mostrar versión |
| ezgo -h | Imprimir ayuda predeterminada Inglés |
| ezgo -hzc | Imprimir ayuda de Chino |
| ezgo -hzt | Imprimir ayuda de Chino tradicional |
| ezgo -hes | Imprimir ayuda de Español |
| ezgo -update [versión] | Actualizar EasyGo |
| ezgo -clone [URL] [ruta local] -branch--[nombre de la rama] -depth--[número] | Clonar repositorio desde la nube |
| ezgo -sync [ruta local] [URL] | Ejecutar sincronización |
| ezgo -sync auto [el tiempo predeterminado es segundo] | Sincronización automática |
| ezgo -sync incremental | Sincronizar solo archivos de diferencia |
| ezgo -config | Configurar EasyGo |
| ezgo -env | Verificación automática del entorno |
| ezgo -logs [nivel] | Salida de registros de ezgo |
| ezgo -logs git | Salida de registros de git |
| ezgo -logs go | Salida de registros de Go |
| ezgo -push [confirmar] | Confirmar en repositorio remoto |
| ezgo -pull [rama] | Extraer una rama |
| ezgo -checkout [nombre de la rama] | Cambiar nombre de rama |
| ezgo -conflict [vía] | Resolver conflictos locales y en la nube |
¿Tu necesitas mas?Mirar el:https://github.com/Sarmioe/EasyGo/blob/main/Introducir.md
`
const hdzc = `
使用-hzc只能输出中文版,如果你在使用其他语言,请使用例如-hes,es是西班牙语的缩写.
这里列出了一些常用命令
| 命令                                                         | 功能                                                         |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| ezgo -v                                                       | 显示版本                                                     |
| ezgo -h                                                       | 输出帮助 默认英文                                            |
| ezgo -hzc                                                      | 输出简体中文的帮助                                           |
| ezgo -htw                                                      | 输出繁体中文的帮助                                           |
| ezgo -hes                                                      | 输出西班牙文的帮助                                           |
| ezgo -update [version]                                        | 更新EasyGo                                                   |
| ezgo -clone [URL] [Localpath] -branch--[branchname] -depth--[number] | 从远端直接克隆仓库                                           |
| ezgo -sync [localpath] [URL]                                  | 执行同步 URL是远程仓库的 如果将ezgo和项目置于一个文件夹 就把Localpath位置设置为./即可 |
| ezgo -sync auto [time defualt is second]                      | 指定时间后自动检测更改和同步                                 |
| ezgo -sync incremental                                        | 只同步差异文件 用git add . + git status实现                  |
| ezgo -config                                                  | 配置EasyGo 包括云端ssh密钥                                   |
| ezgo -env                                                     | 自动环境检查                                                 |
| ezgo -logs [level]                                            | 输出ezgo日志                                                 |
| ezgo -logs git                                                | 输出git日志                                                  |
| ezgo -logs go                                                 | 输出Go日志                                                   |
| ezgo -push [commit]                                           | 提交到远端存储库                                             |
| ezgo -pull [branch]                                           | 拉取一个分支                                                 |
| ezgo -checkout [branchname]                                   | 切换分支名称                                                 |
| ezgo -conflict [way]                                          | 解决云端和本地冲突                                           |
需要更多吗,访问:https://github.com/Sarmioe/EasyGo/blob/main/%E8%AF%BB%E6%88%91.md了解详情
`
const hdtw = `
使用-htw只能輸出中文版,如果你在使用其他語言,請使用例如-hes,es是西班牙文的縮寫.
這裡列出了一些常用指令
| 指令 | 功能 |
|------------------------------------------------- ----------- | -------------------------------------- ---------------------- |
| ezgo -v | 顯示版本 |
| ezgo -h | 輸出幫助 預設英文 |
| ezgo -hzc | 輸出簡體中文的幫助 |
| ezgo -htw | 輸出繁體中文的幫助 |
| ezgo -hes | 輸出西班牙文的幫助 |
| ezgo -update [version] | 更新EasyGo |
| ezgo -clone [URL] [Localpath] -branch--[branchname] -depth--[number] | 從遠端直接克隆倉庫 |
| ezgo -sync [localpath] [URL] | 執行同步 URL是遠端倉庫的 如果將ezgo和專案置於一個資料夾 就把Localpath位置設為./即可 |
| ezgo -sync auto [time defualt is second] | 指定時間後自動偵測變更與同步 |
| ezgo -sync incremental | 只同步差異檔 用git add . + git status實作 |
| ezgo -config | 設定EasyGo 包含雲端ssh金鑰 |
| ezgo -env | 自動環境檢查 |
| ezgo -logs [level] | 輸出ezgo日誌 |
| ezgo -logs git | 輸出git日誌 |
| ezgo -logs go | 輸出Go日誌 |
| ezgo -push [commit] | 提交到遠端儲存庫 |
| ezgo -pull [branch] | 拉取一個分支 |
| ezgo -checkout [branchname] | 切換分支名稱 |
| ezgo -conflict [way] | 解決雲端與本地衝突 |
需要更多嗎,瀏覽:https://github.com/Sarmioe/EasyGo/blob/main/%E8%AF%BB%E6%88%91.md了解詳情
`

var targets = []struct {
	os   string
	arch string
}{
	{"windows", "amd64"},
	{"windows", "386"},
	{"windows", "arm"},
	{"windows", "arm64"},
	{"linux", "amd64"},
	{"linux", "386"},
	{"linux", "arm"},
	{"linux", "arm64"},
	{"linux", "ppc64"},
	{"linux", "ppc64le"},
	{"linux", "mips"},
	{"linux", "mipsle"},
	{"linux", "mips64"},
	{"linux", "mips64le"},
	{"linux", "riscv64"},
	{"darwin", "amd64"},
	{"darwin", "arm64"},
	{"freebsd", "amd64"},
	{"freebsd", "386"},
	{"freebsd", "arm"},
	{"freebsd", "arm64"},
	{"openbsd", "amd64"},
	{"openbsd", "386"},
	{"openbsd", "arm"},
	{"openbsd", "arm64"},
	{"netbsd", "amd64"},
	{"netbsd", "386"},
	{"netbsd", "arm"},
	{"netbsd", "arm64"},
	{"dragonfly", "amd64"},
	{"solaris", "amd64"},
	{"plan9", "amd64"},
	{"plan9", "386"},
	{"plan9", "arm"},
	{"aix", "ppc64"},
	{"illumos", "amd64"},
	{"hurd", "amd64"},
}

func downloadZip(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download ZIP file: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save ZIP file: %w", err)
	}

	return nil
}
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer r.Close()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("failed to create extraction directory: %w", err)
	}

	for _, file := range r.File {
		filePath := filepath.Join(dest, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, file.Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		outFile, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer outFile.Close()

		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return fmt.Errorf("failed to extract file: %w", err)
		}
	}

	return nil
}
func buildSourceCode(srcDir string) error {
	cmd := exec.Command("go", "build", "-o", "EasyGo")
	cmd.Dir = srcDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	return nil
}
func getVersion(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
func atfs() {
	versionFlag := flag.Bool("v", false, "Display Version")
	helpFlag := flag.Bool("h", false, "Display Help")
	ayudaFlag := flag.Bool("hes", false, "Mostrar ayuda")
	bangzhuFlag := flag.Bool("hzc", false, "输出帮助")
	zhiyuanFlag := flag.Bool("htw", false, "輸出幫助")
	ezgoupdate := flag.Bool("update", false, "Update EasyGo")
	checkenv := flag.Bool("env", false, "Check environment")
	build := flag.Bool("gobuild", false, "Build the go project.")
	buildall := flag.Bool("gobuildall", false, "Build all the go project.")
	flag.Parse()
	if *versionFlag {
		fmt.Println("Version is:" + version)
		os.Exit(0)
	}
	if *helpFlag {
		fmt.Println(hden)
		os.Exit(0)
	}
	if *ayudaFlag {
		fmt.Println(hdes)
		os.Exit(0)
	}
	if *bangzhuFlag {
		fmt.Println(hdzc)
		os.Exit(0)
	}
	if *zhiyuanFlag {
		fmt.Println(hdtw)
		os.Exit(0)
	}
	if *ezgoupdate {
		fmt.Println("EasyGo Start run build to update , download zip from https://github.com/Sarmioe/EasyGo/archive/refs/heads/main.zip")
		zipURL := "https://github.com/Sarmioe/EasyGo/archive/refs/heads/main.zip"
		zipDest := "source.zip"
		extractDir := "source"

		fmt.Println("Downloading ZIP file...")
		if err := downloadZip(zipURL, zipDest); err != nil {
			fmt.Println("Error downloading ZIP file:", err)
			return
		}

		fmt.Println("Extracting ZIP file...")
		if err := unzip(zipDest, extractDir); err != nil {
			fmt.Println("Error extracting ZIP file:", err)
			return
		}

		fmt.Println("Building source code...")
		if err := buildSourceCode(extractDir); err != nil {
			fmt.Println("Error building source code:", err)
			return
		}

		fmt.Println("Build complete! The program is ready.")
		fmt.Println("After 5 seconds , the programm will be auto exit , you need restart it.")
		os.Exit(5)
	}
	if *checkenv {
		fmt.Println("Checking environment...")
		if _, err := exec.LookPath("git"); err != nil {
			fmt.Println("Git not found.")
			os.Exit(0)
		}
		if _, err := exec.LookPath("go"); err != nil {
			fmt.Println("Go not found.")
			os.Exit(0)
		}
		gitVersion, err := getVersion("git", "--version")
		if err != nil {
			fmt.Println("Error getting Git version:", err)
			os.Exit(0)
		} else {
			fmt.Println("Git version:", gitVersion)
		}
		goVersion, err := getVersion("go", "version")
		if err != nil {
			fmt.Println("Error getting Go version:", err)
			os.Exit(0)
		} else {
			fmt.Println("Go version:", goVersion)
		}
		fmt.Println("All the environment is ready.")
		os.Exit(0)
	}
	if *build {
		fmt.Println("Start build your go project.")
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter the absolute path to the Go project: ")
		projectPath, _ := reader.ReadString('\n')
		projectPath = strings.TrimSpace(projectPath)

		if _, err := os.Stat(projectPath); os.IsNotExist(err) {
			fmt.Println("The specified path does not exist. Please try again.")
			os.Exit(1)
		}

		fmt.Printf("Enter target OS (default: %s): ", runtime.GOOS)
		targetOS, _ := reader.ReadString('\n')
		targetOS = strings.TrimSpace(targetOS)
		if targetOS == "" {
			targetOS = runtime.GOOS
		}

		fmt.Printf("Enter target architecture (default: %s): ", runtime.GOARCH)
		targetArch, _ := reader.ReadString('\n')
		targetArch = strings.TrimSpace(targetArch)
		if targetArch == "" {
			targetArch = runtime.GOARCH
		}
		fmt.Println("INFO:If your opriting system is windows , please add ", "'.exe'", ".")
		fmt.Print("Enter output binary name (default: Go-project.exe): ")
		outputName, _ := reader.ReadString('\n')
		outputName = strings.TrimSpace(outputName)
		if outputName == "" {
			outputName = "Go-project.exe"
		}

		outputPath := filepath.Join(projectPath, outputName)

		fmt.Println("Setting up build environment...")
		env := os.Environ()
		env = append(env, "GOOS="+targetOS)
		env = append(env, "GOARCH="+targetArch)

		cmd := exec.Command("go", "build", "-o", outputPath, projectPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = env

		fmt.Println("Starting the build process...")
		if err := cmd.Run(); err != nil {
			fmt.Printf("Build failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Build succeeded! Output file: %s\n", outputPath)
		os.Exit(0)
	}
	if *buildall {
		fmt.Println("Start build your all the go project.")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the path to the Go project: ")
		projectPath, _ := reader.ReadString('\n')
		projectPath = strings.TrimSpace(projectPath)

		if _, err := os.Stat(projectPath); os.IsNotExist(err) {
			fmt.Println("The specified path does not exist. Please try again.")
			os.Exit(1)
		}

		fmt.Print("Enter output binary base name (default: Go-project): ")
		outputBaseName, _ := reader.ReadString('\n')
		outputBaseName = strings.TrimSpace(outputBaseName)
		if outputBaseName == "" {
			outputBaseName = "Go-project"
		}

		outputDir := filepath.Join(projectPath, "build")
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Printf("Failed to create output directory: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Starting the build process for all platforms...")
		fmt.Println("The ended output filename just like : Go-project-windows-amd64.exe")
		for _, target := range targets {
			targetOS := target.os
			targetArch := target.arch

			outputFile := fmt.Sprintf("%s-%s-%s", outputBaseName, targetOS, targetArch)
			if targetOS == "windows" {
				outputFile += ".exe"
			}
			outputPath := filepath.Join(outputDir, outputFile)
			env := os.Environ()
			env = append(env, "GOOS="+targetOS, "GOARCH="+targetArch)

			cmd := exec.Command("go", "build", "-o", outputPath, projectPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Env = env

			fmt.Printf("Building for %s/%s...\n", targetOS, targetArch)
			if err := cmd.Run(); err != nil {
				fmt.Printf("Build failed for %s/%s: %v\n", targetOS, targetArch, err)
				continue
			}
			fmt.Println("At the floader:", outputPath)
		}
		fmt.Println("Build succeeded: %s\n")
		fmt.Println("Created all the 12 Systems , 11 archs , and 41 files.")
		os.Exit(0)
	}
}
func main() {
	fmt.Println("Welcome to EasyGo!")
	fmt.Println("Powered by Sarmioe and Golang V1.23.4")
	atfs()
	fmt.Println("To get help document , view this page :https://github.com/Sarmioe/EasyGo/blob/main/README.md")
	fmt.Println("And if you not now how can using EasyGo , Using : 'ezgo -h' to get help")
	fmt.Println("Now , you no add any bool value , the programm will be exit...")
}
