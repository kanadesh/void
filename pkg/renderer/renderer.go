package renderer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/kanadeishii/void/internal/bundler"
	"github.com/kanadeishii/void/internal/counter"
	"github.com/kanadeishii/void/internal/logger"
	"github.com/kanadeishii/void/internal/types"
	"github.com/kanadeishii/void/internal/utils"
)

func Render(link string, doLogging bool) error {
	pioneer, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()
	listener := `
	(() => {
		const channel = new BroadcastChannel('application')

		return new Promise((resolve) => {
			channel.addEventListener('message', ({ data }) => {
				if(data.action === 'response'){
					resolve(JSON.stringify(data.body))
				}
			})
			
			channel.postMessage({
				action: 'request',
				body: null
			}, '*')
		})
	})()
	`

	var raw string
	var option types.Option

	if err := chromedp.Run(pioneer,
		chromedp.Navigate(link),
		chromedp.Evaluate(
			listener,
			&raw,
			func(params *runtime.EvaluateParams) *runtime.EvaluateParams {
				return params.WithAwaitPromise(true)
			},
		),
	); err != nil {
		if doLogging {
			log.Fatal(err)
		} else {
			return err
		}
	}
	cancel()

	if err := json.Unmarshal([]byte(raw), &option); err != nil {
		if doLogging {
			log.Fatal(err)
		} else {
			return err
		}
	}

	cwd, _ := os.Getwd()

	option.CacheDir = filepath.Join(cwd, option.CacheDir)
	option.ResultFile = filepath.Join(cwd, option.ResultFile)

	os.RemoveAll(option.CacheDir)
	os.MkdirAll(option.CacheDir, 0777)

	var state int = 0
	var audios []types.Audio
	httpClient := new(http.Client)
	for index, audio := range option.Audios {
		request, _ := http.NewRequest("GET", utils.FixUrl(audio.Link, link), nil)
		response, err := httpClient.Do(request)
		if err != nil {
			if doLogging {
				log.Fatal(err)
			} else {
				return err
			}
		}
		defer response.Body.Close()

		contentType := response.Header.Get("Content-Type")
		if _, ok := utils.AudioMap[contentType]; ok {
			extension := utils.AudioMap[contentType]
			file := filepath.Join(option.CacheDir, strconv.Itoa(index)+"."+extension)
			binary, _ := io.ReadAll(response.Body)
			os.WriteFile(file, binary, 0644)

			audios = append(audios, types.Audio{
				File:  file,
				Start: audio.Start,
			})
		}
		response.Body.Close()

		if doLogging {
			counter.Count(index+1, len(option.Audios), "collecting audios")
		}
	}

	state = 0
	var waitGroup sync.WaitGroup

	for index := 0; index < option.Number; index++ {
		waitGroup.Add(1)

		go func(offset int, total int) {
			context, cancel := chromedp.NewContext(
				context.Background(),
			)
			defer cancel()

			chromedp.Run(
				context,
				emulation.SetDeviceMetricsOverride(int64(option.Width), int64(option.Height), 1, false),
				chromedp.Navigate(link),
				chromedp.Evaluate(listener,
					nil,
					func(params *runtime.EvaluateParams) *runtime.EvaluateParams {
						return params.WithAwaitPromise(true)
					},
				),
			)

			for frame := 0; frame < total; frame++ {
				var buffer []byte

				dispatcher := `
				(() => {
					const channel = new BroadcastChannel('application')
			
					return new Promise((resolve) => {
						channel.addEventListener('message', ({ data }) => {
							if(data.action === 'ok'){
								resolve(null)
							}
						})
						channel.postMessage({
							action: 'load',
							body: {
								frame: ` + strconv.Itoa(total*offset+frame) + `
							}
						}, '*')
					})
				})()
				`

				chromedp.Run(context,
					chromedp.Evaluate(
						dispatcher,
						nil,
						func(params *runtime.EvaluateParams) *runtime.EvaluateParams {
							return params.WithAwaitPromise(true)
						}),
					chromedp.Screenshot(
						"#application",
						&buffer,
						chromedp.NodeVisible,
						chromedp.ByQuery,
					),
				)
				os.WriteFile(option.CacheDir+"/"+fmt.Sprintf("%010d", total*offset+frame)+".png", buffer, 0644)

				state++
				if doLogging {
					counter.Count(state, option.Frames, "working with chromedp")
				}
			}

			waitGroup.Done()
		}(index, int(math.Ceil((float64(option.Frames) / float64(option.Number)))))
	}

	waitGroup.Wait()

	if doLogging {
		fmt.Println("\r⠿  completed!                     ")

		fmt.Print("\n")
		logger.Rich(logger.ColorBlue, "void", "done")
	}

	return bundler.Bundle(option, audios, doLogging)
}
