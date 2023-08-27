package bundler

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/kanadesh/void/internal/logger"
	"github.com/kanadesh/void/internal/types"
)

func Bundle(option types.Option, audios []types.Audio, doLogging bool) error {
	if doLogging {
		fmt.Print("\n")
		logger.Rich(logger.ColorGreen, "ffmpeg", "bundling")
	}

	os.Remove(option.ResultFile)

	var tmp string
	if len(audios) == 0 {
		tmp = option.ResultFile
	} else {
		tmp = filepath.Join(option.CacheDir, "tmp.mp4")
	}

	if err := exec.Command(
		"ffmpeg",

		"-framerate",
		strconv.Itoa(option.Fps),

		"-start_number",
		"0",

		"-i",
		filepath.Join(option.CacheDir, "%010d.png"),

		"-vframes",
		strconv.Itoa(option.Frames),

		"-c:v",
		"libx264",
		"-pix_fmt",
		"yuv420p",

		"-vf",
		"scale=trunc(iw/2)*2:trunc(ih/2)*2",

		tmp,
	).Run(); err != nil {
		if doLogging {
			log.Fatal(err)
		} else {
			return err
		}
	}

	if len(audios) != 0 {
		cmd := []string{
			"-i",
			tmp,
		}
		var filter string

		for index, audio := range audios {
			cmd = append(cmd, "-i", audio.File)
			filter += "[" + strconv.Itoa(index+1) + ":a]adelay=" + strconv.FormatFloat(float64(audio.Start)/float64(option.Fps), 'f', 2, 64) + "s:all=1[a" + strconv.Itoa(index+1) + "];"
		}

		for index := range audios {
			filter += "[a" + strconv.Itoa(index+1) + "]"
		}

		filter += "amix=inputs=" + strconv.Itoa(len(audios)) + "[amixout]"

		cmd = append(cmd, "-filter_complex", filter, "-map", "0:v:0", "-map", "[amixout]", "-c:v", "copy", "-c:a", "aac", "-b:a", "192k", option.ResultFile)

		if err := exec.Command(
			"ffmpeg",
			cmd...,
		).Run(); err != nil {
			if doLogging {
				log.Fatal(err)
			} else {
				return err
			}
		}
	}

	os.RemoveAll(option.CacheDir)

	if doLogging {
		fmt.Print("\n")
		logger.Rich(logger.ColorGreen, "ffmpeg", "done: "+option.ResultFile)
	}

	return nil
}
