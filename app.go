package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"zaprLauncher/backend/update"
	"zaprLauncher/backend/utils"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context

	projectDir  string
	exeDir      string
	ExeFilePath string

	versionFilePath string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	done := make(chan struct{})

	go func() {
		a.projectDir = utils.GetAppDataPath("ZaprUI")
		a.exeDir = a.projectDir + "/bin"

		// creating ProjectDir in User/AppData/Roaming/
		if err := ensureAppDir(a.projectDir); err != nil {
			panic(fmt.Errorf("❗error creating app directory in AppData: %w", err))
		}
		if err := ensureAppDir(a.exeDir); err != nil { // temp dir for sessions data files
			panic(fmt.Errorf("❗error creating gitrepo directory in project dir: %w", err))
		}

		// =============================================== fetching zaprUI

		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		release, err := update.ParceLatestRelease(client) // Asking GitHub Releases about latest
		if err != nil {
			panic(fmt.Errorf("❗error parce latest release: %v", err))
		}

		if err := update.EnsureVersionFileExist(a.projectDir, release); err != nil { //  Check VersionFile
			panic(fmt.Errorf("❗version file ensure error: %v", err))
		}
		a.versionFilePath = filepath.Join(a.projectDir, "zaprUI_version.txt")

		// CHECKING Latest and Ready
		latest, err := update.IsLatestVersion(a.versionFilePath, release) // Trying version
		if err != nil {
			panic(fmt.Errorf("❗failed to check version: %v", err))
		}
		ready, err := update.IsReleaseReady(a.exeDir)
		if err != nil {
			panic(fmt.Errorf("❗failed to check release files: %v", err))
		}

		if latest && ready {
			fmt.Println("You use actual version!")
		} else {
			if err := update.DownloadReleaseExe(client, release, a.exeDir); err != nil {
				panic(fmt.Errorf("❗Downloading failed because of: %v", err))
			}
			a.ExeFilePath = filepath.Join(a.exeDir, "zaprUi.exe")
		}
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("Update finished")
	case <-ctx.Done():
		fmt.Println("Timeout reached")
	}

	// тут происходит запрос UAC для приложухи и её вызов

}

// Getting sure that ProjectDir created
func ensureAppDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// ===================================== WAILS API ==========================

// OpenURL opens the specified URL in the default browser
func (a *App) OpenURL(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}
