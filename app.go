package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {

	if screens, err := runtime.ScreenGetAll(ctx); err == nil && len(screens) > 0 {
		w, h := runtime.WindowGetSize(ctx)
		runtime.WindowSetPosition(ctx, screens[0].Width-w-20, screens[0].Height-h-120)
	}
	a.ctx = ctx
}

//// Greet returns a greeting for the given name
//func (a *App) Greet(name string) string {
//	return fmt.Sprintf("Hello %s, It's show time!", name)
//}

func (a *App) onBeforeClose(ctx context.Context) (prevent bool) {
	if cmdScrcpy != nil {
		_ = cmdScrcpy.Process.Kill()
	}
	x, y := runtime.WindowGetPosition(ctx)
	fmt.Printf("Hello %d, %d It's show time!", x, y)
	return false
}

var cmdScrcpy *exec.Cmd

func (a *App) ConnPhone(isTop, isAwake bool) error {
	if cmdScrcpy != nil {
		return errors.New("请先退出app")
	}
	args := []string{}
	if isTop {
		args = append(args, "--always-on-top")
	}
	if isTop {
		args = append(args, "--stay-awake")
	}
	cmdScrcpy = exec.CommandContext(a.ctx, "./scrcpy/scrcpy.exe", args...)
	cmdScrcpy.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	//cmdScrcpy = exec.CommandContext(a.ctx, "wscript.exe", "CreateObject(\"Wscript.Shell\").Run strCommand, 0, false ./scrcpy/scrcpy.exe --always-on-top --stay-awake")
	//cmdScrcpy.SysProcAttr = &syscall.SysProcAttr{NoInheritHandles: false}
	log.Println("conn Phone")
	defer func() {
		cmdScrcpy = nil
		log.Println("exit Phone")
	}()
	cmdScrcpy.Stdout = log.New(os.Stdout, "【scrcpy】", log.LstdFlags).Writer()
	if err := cmdScrcpy.Run(); err != nil {
		return err
	}
	return nil
	//scrcpy --always-on-top --stay-awake
}

func (a *App) Msg(title, msg string) (string, error) {
	return runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   title,
		Message: msg,
	})
}

func (a *App) InstallApk() error {
	p, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{"apk", "*.apk"},
		},
	})
	if err != nil {
		log.Println("路径不对", err)
		return err
	}
	installCmd := exec.CommandContext(a.ctx, "./scrcpy/adb.exe", "install", p)
	installCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	installCmd.Stdout = log.New(os.Stdout, "【adb】", log.LstdFlags).Writer()
	if err := installCmd.Run(); err != nil {
		log.Println(err, p)
		return err
	}
	return nil
}

func (a *App) UnInstallApk() error {
	bs, err := ioutil.ReadFile(".package")
	if err != nil {
		log.Println(err)
		return err
	}
	installCmd := exec.CommandContext(a.ctx, "./scrcpy/adb.exe", "uninstall", string(bs))
	installCmd.Stdout = log.New(os.Stdout, "【adb】", log.LstdFlags).Writer()
	//installCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := installCmd.Run(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
