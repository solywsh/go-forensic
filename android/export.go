package android

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/electricbubble/gadb"
	"github.com/solywsh/go-forensic/utils/osx"
	"github.com/solywsh/go-forensic/utils/pathx"
	"github.com/solywsh/go-forensic/utils/printer"
	"path/filepath"
	"strings"
)

func (t *Android) ExportAppDataByKeywords(keywords ...string) error {
	if len(keywords) == 0 {
		return fmt.Errorf("please input keywords")
	}
	t.spinner = printer.NewSpinner()
	t.spinner.SetSpinner(spinner.Moon).Msg("loading...")
	t.spinner.Run()
	defer t.spinner.Msg("export by keywords is done.")
	defer t.spinner.Quit()
	adbClient, err := gadb.NewClient()
	if err != nil {
		return err
	}
	devices, err := adbClient.DeviceList()
	if err != nil {
		return err
	}
	if len(devices) == 0 {
		return fmt.Errorf("list of devices is empty")
	}
	t.device = devices[0]
	absOutput, err := filepath.Abs(t.output)
	if err != nil {
		return err
	}
	t.output = absOutput
	t.keywords = keywords
	t.removeOutput()
	for focusPath, layer := range androidFocusMap {
		t.spinner.Msg(fmt.Sprintf("searching %s", focusPath))
		resFileList, err := t.checkForMatchingSubDirs(focusPath, layer)
		if err != nil {
			log.Error(err)
			continue
		}
		resFileList = pathx.FilterSubDir(resFileList)
		for _, _file := range resFileList {
			err := t.exportWithAdb(_file)
			if err != nil {
				log.Error(err)
			}
		}
	}
	return nil
}

func (t *Android) exportWithAdb(remoteFilePath string) error {
	t.spinner.Msg(fmt.Sprintf("handling %s", remoteFilePath))
	tarFileName := filepath.Base(remoteFilePath) + ".tar"
	tarFileRemotePath := pathx.PathJoin("/sdcard/", tarFileName)
	localTarFilePath := filepath.Join(t.output, tarFileName)

	t.spinner.Msg(fmt.Sprintf("packing %s --> %s", remoteFilePath, tarFileRemotePath))
	if _, err := t.AdbRunShellCommand("tar", "-cf", tarFileRemotePath, remoteFilePath); err != nil {
		return fmt.Errorf("failed to create tar for %s: %s\n", remoteFilePath, err)
	}
	t.spinner.Msg(fmt.Sprintf("pulling %s --> %s", tarFileRemotePath, localTarFilePath))
	if err := t.AdbPullFile(tarFileRemotePath, localTarFilePath); err != nil {
		return fmt.Errorf("failed to pull tar file: %s\n", err)
	}
	t.spinner.Msg(fmt.Sprintf("removing %s", tarFileRemotePath))
	if _, err := t.AdbRunShellCommand("rm", tarFileRemotePath); err != nil {
		return fmt.Errorf("failed to remove %s: %s\n", tarFileRemotePath, err)
	}
	t.spinner.Msg(fmt.Sprintf("extracting %s", localTarFilePath))
	if err := osx.TarDecompression(localTarFilePath, t.output); err != nil {
		return err
	}
	return nil
}

// check for matching sub directories
func (t *Android) checkForMatchingSubDirs(dir string, layer int) ([]string, error) {
	var res []string
	if layer == 0 {
		return res, nil // max layer reached, no more recursion
	}
	output, err := t.AdbRunShellCommand("ls", dir)
	if err != nil {
		return res, fmt.Errorf("failed to list directory %s: %s\n", dir, err)
	}
	for _, item := range strings.Fields(output) {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		fullPath := pathx.PathJoin(dir, item)
		flag := false
		for _, keyword := range t.keywords {
			if strings.Contains(strings.ToLower(item), strings.ToLower(keyword)) {
				res = append(res, fullPath)
				flag = true
				break
			}
		}
		if flag {
			continue
		}
		// recursively check sub directories
		resChild, err := t.checkForMatchingSubDirs(fullPath, layer-1)
		if err != nil {
			log.Error(err)
			continue
		}
		res = append(res, resChild...)
	}
	return res, nil
}

func (t *Android) ExportBySpecify(pathList ...string) error {
	if len(pathList) == 0 {
		return fmt.Errorf("please input path")
	}
	t.spinner = printer.NewSpinner()
	t.spinner.SetSpinner(spinner.Moon).Msg("loading...")
	t.spinner.Run()
	defer t.spinner.Msg("export by specify is done.")
	defer t.spinner.Quit()
	adbClient, err := gadb.NewClient()
	if err != nil {
		return err
	}
	devices, err := adbClient.DeviceList()
	if err != nil {
		return err
	}
	if len(devices) == 0 {
		return fmt.Errorf("list of devices is empty")
	}
	t.device = devices[0]
	absOutput, err := filepath.Abs(t.output)
	if err != nil {
		return err
	}
	t.output = absOutput
	t.removeOutput()
	for _, _path := range pathList {
		err := t.exportWithAdb(_path)
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}
