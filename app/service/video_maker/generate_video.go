package videomaker

import (
	"context"
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/faiface/beep/mp3"

	"github.com/simpleAI/service-video-maker/app/structs/model"
)

func (s *serviceVideoMaker) GenerateVideo(ctx context.Context, request *model.GenerateVideoRequest) (string, error) {
	pwd, errPwd := os.Getwd()
	if errPwd != nil {
		return "", errPwd
	}

	tmpFolderPath := fmt.Sprintf("%s/tmp/%s", pwd, request.Id)

	if err := os.MkdirAll(tmpFolderPath, os.ModePerm); err != nil {
		return "", err
	}

	if err := s.downloadData(ctx, request, tmpFolderPath); err != nil {
		return "", err
	}

	outputPath := fmt.Sprintf("%s/video.mp4", tmpFolderPath)
	imageListFile := fmt.Sprintf("%s/image_list.txt", tmpFolderPath)
	noSound := fmt.Sprintf("%s/no_sound.mp4", tmpFolderPath)
	backgroundFile := fmt.Sprintf("%s/background.mp3", tmpFolderPath)
	voiceFile := fmt.Sprintf("%s/voice.mp3", tmpFolderPath)
	mixedAudio := fmt.Sprintf("%s/mixed_audio.mp3", tmpFolderPath)
	noSubtitle := fmt.Sprintf("%s/no_subtitle.mp4", tmpFolderPath)
	subtitlePath := fmt.Sprintf("%s/subtitles.srt", tmpFolderPath)

	errGenerateNoSound := s.repository.GetCommandRepository().Run(ctx,
		"ffmpeg", "-y", "-f", "concat", "-safe", "0", "-i", imageListFile, "-vf", "scale=1280:720", "-pix_fmt", "yuv420p", noSound,
	)
	if errGenerateNoSound != nil {
		return "", errGenerateNoSound
	}

	errGenerateMixedAudio := s.repository.GetCommandRepository().Run(ctx,
		"ffmpeg", "-y", "-i", backgroundFile, "-i", voiceFile,
		"-filter_complex", "[0:a]volume=0.3[a1];[a1][1:a]amix=inputs=2:duration=shortest",
		mixedAudio,
	)
	if errGenerateMixedAudio != nil {
		return "", errGenerateMixedAudio
	}

	errGenerateNoSubtitles := s.repository.GetCommandRepository().Run(ctx,
		"ffmpeg", "-y", "-i", noSound, "-i", mixedAudio, "-c:v", "copy",
		"-c:a", "aac", noSubtitle,
	)
	if errGenerateNoSubtitles != nil {
		return "", errGenerateNoSubtitles
	}

	errSubtitles := s.createSubtitles(subtitlePath, request)
	if errSubtitles != nil {
		return "", errSubtitles
	}

	fontSize := 40
	subtitleFontConfig := fmt.Sprintf(`subtitles=%s:force_style='Fontsize=%d`, subtitlePath, fontSize)

	errGenerateOutput := s.repository.GetCommandRepository().Run(ctx,
		"ffmpeg", "-y", "-i", noSubtitle, "-vf",
		subtitleFontConfig, outputPath,
	)
	if errGenerateOutput != nil {
		return "", errGenerateOutput
	}

	return outputPath, nil
}

func (s *serviceVideoMaker) createSubtitles(path string, request *model.GenerateVideoRequest) error {
	subtitlesFile, errCreate := os.Create(path)
	if errCreate != nil {
		return errCreate
	}
	defer subtitlesFile.Close()

	for i, transcript := range request.Data.Transcripts {
		startTime := s.formatTime(transcript.Start)
		endTime := s.formatTime(transcript.End)
		subtitlesFile.WriteString(fmt.Sprintf("%d\n", i+1))
		subtitlesFile.WriteString(fmt.Sprintf("%s --> %s\n", startTime, endTime))
		subtitlesFile.WriteString(fmt.Sprintf("%s\n\n", transcript.Words))
	}
	return nil
}

func (s *serviceVideoMaker) formatTime(seconds float64) string {
	millis := int((seconds - float64(int(seconds))) * 1000)
	duration := time.Duration(int(seconds)) * time.Second
	timeStr := fmt.Sprintf("%02d:%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)
	return fmt.Sprintf("%s,%03d", timeStr, millis)
}

func (s *serviceVideoMaker) downloadData(ctx context.Context, request *model.GenerateVideoRequest, tmpFolderPath string) error {
	wg := new(sync.WaitGroup)
	wg.Add(2)
	var errs = make(chan error, 2)
	defer close(errs)

	go func() {
		defer wg.Done()

		errDownloadBackground := s.downloadFile(ctx, request.Data.BackgroundURL, tmpFolderPath+"/background.mp3")
		if errDownloadBackground != nil {
			errs <- errDownloadBackground
		}
	}()

	go func() {
		defer wg.Done()

		errDownloadVoice := s.downloadFile(ctx, request.Data.VoiceURL, tmpFolderPath+"/voice.mp3")
		if errDownloadVoice != nil {
			errs <- errDownloadVoice
		}
	}()

	wg.Wait()

	if len(errs) > 0 {
		return <-errs
	}

	wg.Add(2)

	var backgroundDuration time.Duration
	var voiceDuration time.Duration
	go func() {
		defer wg.Done()

		backgroundDurationTmp, err := s.getMP3Duration(tmpFolderPath + "/background.mp3")
		if err != nil {
			errs <- err
			return
		}
		backgroundDuration = backgroundDurationTmp
	}()
	go func() {
		defer wg.Done()

		voiceDurationTmp, err := s.getMP3Duration(tmpFolderPath + "/voice.mp3")
		if err != nil {
			errs <- err
			return
		}
		voiceDuration = voiceDurationTmp
	}()

	wg.Wait()

	if len(errs) > 0 {
		return <-errs
	}

	totalVideoDuration := math.Max(backgroundDuration.Seconds(), voiceDuration.Seconds())
	imageDuration := totalVideoDuration / float64(len(request.Data.ImageList))

	wg.Add(len(request.Data.ImageList))
	errs = make(chan error, len(request.Data.ImageList))
	defer close(errs)

	var imageList []string
	for i, url := range request.Data.ImageList {
		imagePath := fmt.Sprintf("%s/image-%04d.jpg", tmpFolderPath, i)

		go func() {
			defer wg.Done()

			err := s.downloadFile(ctx, url, imagePath)
			if err != nil {
				errs <- err
			}
		}()

		imageList = append(imageList, imagePath)
	}

	wg.Wait()

	if len(errs) > 0 {
		return <-errs
	}

	errSaveImageList := s.saveImageListToFile(imageList, imageDuration, tmpFolderPath+"/image_list.txt")
	if errSaveImageList != nil {
		return errSaveImageList
	}

	return nil
}

func (s *serviceVideoMaker) downloadFile(ctx context.Context, url, path string) error {
	_, err := s.restClient.R().SetContext(ctx).SetOutput(path).Get(url)

	if err != nil {
		return err
	}

	return nil
}

func (s *serviceVideoMaker) getMP3Duration(filePath string) (time.Duration, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return 0, err
	}
	defer streamer.Close()

	duration := time.Duration(streamer.Len()) * time.Second / time.Duration(format.SampleRate)

	return duration, nil
}

func (s *serviceVideoMaker) saveImageListToFile(imageList []string, duration float64, path string) error {
	file, errCreate := os.Create(path)
	if errCreate != nil {
		return errCreate
	}
	defer file.Close()

	for _, image := range imageList {
		if _, errWrite := file.WriteString(fmt.Sprintf("file '%s'\n", image)); errWrite != nil {
			return errWrite
		}
		if _, errWrite := file.WriteString(fmt.Sprintf("duration %.4f\n", duration)); errWrite != nil {
			return errWrite
		}
	}

	return nil
}
