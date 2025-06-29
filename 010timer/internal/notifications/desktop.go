package notifications

import (
	"fmt"
	"os/exec"
	"runtime"
)

type DesktopNotifier struct {
	enabled bool
}

func NewDesktopNotifier(enabled bool) *DesktopNotifier {
	return &DesktopNotifier{
		enabled: enabled,
	}
}

func (dn *DesktopNotifier) Show(title, message string) error {
	if !dn.enabled {
		return nil
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "%s"`, message, title))

	case "linux":
		if _, err := exec.LookPath("notify-send"); err == nil {
			cmd = exec.Command("notify-send", title, message, "-i", "clock")
		} else if _, err := exec.LookPath("zenity"); err == nil {
			cmd = exec.Command("zenity", "--info", "--title", title, "--text", message)
		}

	case "windows":
		script := fmt.Sprintf(`
			Add-Type -AssemblyName System.Windows.Forms
			$notification = New-Object System.Windows.Forms.NotifyIcon
			$notification.Icon = [System.Drawing.SystemIcons]::Information
			$notification.BalloonTipTitle = "%s"
			$notification.BalloonTipText = "%s"
			$notification.Visible = $true
			$notification.ShowBalloonTip(3000)
			Start-Sleep -Seconds 3
			$notification.Dispose()
		`, title, message)

		cmd = exec.Command("powershell", "-Command", script)

	default:
		return fmt.Errorf("desktop notifications not supported on %s", runtime.GOOS)
	}

	if cmd != nil {
		return cmd.Run()
	}

	return fmt.Errorf("no suitable notification system found")
}

func (dn *DesktopNotifier) IsAvailable() bool {
	switch runtime.GOOS {
	case "darwin":
		_, err := exec.LookPath("osascript")
		return err == nil

	case "linux":
		_, err1 := exec.LookPath("notify-send")
		_, err2 := exec.LookPath("zenity")
		return err1 == nil || err2 == nil

	case "windows":
		_, err := exec.LookPath("powershell")
		return err == nil

	default:
		return false
	}
}