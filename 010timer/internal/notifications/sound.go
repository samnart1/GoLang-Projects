package notifications

import (
	"fmt"
	"os/exec"
	"runtime"
)

type SoundNotifier struct {
	enabled		bool
	soundFile	string
	volume		int
}

func NewSoundNotifier(enabled bool, soundFile string, volume int) *SoundNotifier {
	return &SoundNotifier{
		enabled: enabled,
		soundFile: soundFile,
		volume: volume,
	}
}

func (sn *SoundNotifier) Play() error {
	if !sn.enabled {
		return nil
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		if sn.soundFile == "" {
			cmd = exec.Command("afplay", sn.soundFile)
		} else {
			cmd = exec.Command("say", "Timer complete")
		}

	case "linux":
		if sn.soundFile != "" {
			players := []string{"aplay", "paplay", "play", "mpg123", "ogg123"}
			for _, player := range players {
				if _, err := exec.LookPath(player); err == nil {
					cmd = exec.Command(player, sn.soundFile)
					break
				}
			}
		}
		if cmd == nil {
			if _, err := exec.LookPath("beep"); err == nil {
				cmd = exec.Command("beep", "-f", "1000", "-l", "500")
			} else {
				fmt.Print("\a")
				return nil
			}
		}

	case "windows":
		if sn.soundFile != "" {
			cmd = exec.Command("powershell", "-c", fmt.Sprintf(`(New-Object Media.SoundPlayer "%s").PlaySync()`, sn.soundFile))
		} else {
			cmd = exec.Command("powershell", "-c", `[console]::beep(1000,500)`)
		}

	default:
		fmt.Print("\a")
		return nil
	}

	if cmd != nil {
		return cmd.Run()
	}

	return nil
}

func (sn *SoundNotifier) IsAvailable() bool {
	switch runtime.GOOS {
	case "darwin":
		_, err := exec.LookPath("afplay")
		return err == nil

	case "linux":
		players := []string{"aplay", "paplay", "play", "beep"}
		for _, player := range players {
			if _, err := exec.LookPath(player); err == nil {
				return true
			}
		}
		return false

	case "windows":
		_, err := exec.LookPath("powershell")
		return err == nil

	default:
		return false
	}
}