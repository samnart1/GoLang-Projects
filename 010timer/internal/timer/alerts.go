package timer

import (
	"fmt"
	"os/exec"
	"runtime"
)

type AlertSystem struct {
	soundEnabled	bool
	desktopEnabled	bool
	soundFile		string
}

func NewAlertSystem(soundEnabled, desktopEnabled bool, soundFile string) *AlertSystem {
	return &AlertSystem{
		soundEnabled: soundEnabled,
		desktopEnabled: desktopEnabled,
		soundFile: soundFile,
	}
}

func (as *AlertSystem) PlaySound() error {
	if !as.soundEnabled {
		return nil
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		if as.soundFile != "" {
			cmd = exec.Command("afplay", as.soundFile)
		} else {
			cmd = exec.Command("say", "Time is up")
		}

	case "linux":
		if as.soundFile != "" {
			players := []string{"aplay", "paplay", "play"}
			for _, player := range players {
				if _, err := exec.LookPath(player); err == nil {
					cmd = exec.Command(player, as.soundFile)
					break
				}
			}
		}
		if cmd == nil {
			cmd = exec.Command("speaker-test", "-t", "sine", "-f", "1000", "-l", "1")
		}

	case "windows":
		if as.soundFile != "" {
			cmd = exec.Command("powershell", "-c", fmt.Sprintf(`(New-Object Media.SoundPlayer "%s").PlaySync()`, as.soundFile))
		} else {
			cmd = exec.Command("powershell", "-c", `[console]::beep(1000,500)`)
		}

	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if cmd != nil {
		return cmd.Run()
	}

	return fmt.Errorf("no suitable audia player found")
}

func (as *AlertSystem) ShowDesktopNotification(title, message string) error {
	if !as.desktopEnabled {
		return nil
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "%s"`, message, title))

	case "linux":
		cmd = exec.Command("notify-send", title, message)

	case "windows":
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf(`[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntim] > $null; $template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02); $toastXml = [xml] $template.GetXml(); $toastXml.GetElementsByTagName("text")[0].AppendChild($toastXml.CreateTextNode("%s")) > $null; $toastXml.GetElementsByTagName("text")[1].AppendChild($toastXml.CreateTextNode("%s")) > $null; $taost = [Windows.UI.Notifications..ToastNotification]::new($toastXml); [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("Timer").Show($toast);`, title, message))

	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}