package ios

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/solywsh/go-forensic/constant"
	"github.com/solywsh/go-forensic/utils/osx"
	"github.com/solywsh/go-forensic/utils/osx/ssh"
	"github.com/solywsh/go-forensic/utils/pathx"
	"github.com/solywsh/go-forensic/utils/printer"
	"howett.net/plist"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (t *SSHelper) init() error {
	t.spinner = printer.NewSpinner()
	t.spinner.SetSpinner(spinner.Line).Msg("loading...")
	t.spinner.Run()
	t.spinner.Msg("getting ssh connection...")
	var err error
	var sshClient *ssh.Client
	timeout := constant.GetTimeout()
	interval := constant.GetInterval()
	timeSpend := time.Duration(0)
	for i := int64(0); i < (int64(timeout) / int64(interval)); i++ {
		if timeSpend > timeout {
			return fmt.Errorf("ssh connection timeout")
		}
		sshClient, err = t._ssh().C()
		if err != nil {
			timeSpend += interval
			time.Sleep(interval)
			continue
		}
		break
	}
	if err != nil {
		return err
	}
	if sshClient == nil {
		return fmt.Errorf("ssh connection failed")
	}
	sftpClient := sshClient.SFTP()
	if sftpClient == nil {
		return nil
	}
	t.sshClient = sshClient
	t.sftpClient = sftpClient
	t.removeOutput()
	return nil
}

func (t *SSHelper) ExportAppDataByKeywords(keywords ...string) error {
	if len(keywords) == 0 {
		return fmt.Errorf("please input keywords")
	}
	err := t.init()
	if err != nil {
		return err
	}
	defer t.sshClient.Close()
	defer t.sftpClient.Close()
	defer t.spinner.Quit()
	defer t.spinner.Msg("export by keywords is done.")
	t.keywords = keywords
	for _, focusPath := range focusPathList {
		t.spinner.Msg(fmt.Sprintf("searching %s", focusPath))
		err := t.handleFocusPath(focusPath)
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func (t *SSHelper) handleFocusPath(tryPath string) error {
	//if t.cleanDirBefore {
	//	destinationDir := filepath.Join(t.output, tryPath)
	//	if pathx.PathExists(destinationDir) {
	//		// if you get stuck, then sudo is recommended
	//		t.spinner.Msg(fmt.Sprintf("cleaning %s", destinationDir))
	//		osx.RemoveAll(destinationDir)
	//	}
	//}
	tryFileList, err := t.sftpClient.ReadDir(tryPath)
	if err != nil {
		return err
	}
	for _, containerDir := range tryFileList {
		if !containerDir.IsDir() {
			continue
		}
		t.spinner.Msg(fmt.Sprintf("searching %s", containerDir.Name()))
		err := t.handleContainerDir(tryPath, containerDir)
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func (t *SSHelper) handleContainerDir(tryPath string, containerDir os.FileInfo) error {
	containerFileList, err := t.sftpClient.ReadDir(pathx.PathJoin(tryPath, containerDir.Name()))
	if err != nil {
		return err
	}
	var exist bool
	for _, file := range containerFileList {
		if file.IsDir() {
			continue
		}
		if file.Name() == ContainerMetaDataFile {
			exist = true
		}
	}
	if !exist {
		return nil
	}
	containerMetaDataFilePath := pathx.PathJoin(tryPath, containerDir.Name(), ContainerMetaDataFile)
	file, err := t.sftpClient.Open(containerMetaDataFilePath)
	defer file.Close()
	if err != nil {
		return err
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	packageInfo := &struct {
		PackageName string `plist:"MCMMetadataIdentifier"`
	}{}
	if _, err := plist.Unmarshal(content, packageInfo); err != nil {
		return err
	}
	if !t.ComparisonKeywords(packageInfo.PackageName) {
		return nil
	}
	remoteDir := pathx.PathJoin(tryPath, containerDir.Name())
	err = t.export(remoteDir)
	if err != nil {
		return err
	}
	return nil
}

func (t *SSHelper) ComparisonKeywords(str string) bool {
	for _, keyword := range t.keywords {
		if strings.Contains(strings.ToLower(str), strings.ToLower(keyword)) {
			//log.Info("Hit keywords", "packageName", str, "keywords", keyword)
			t.spinner.Msg(fmt.Sprintf("hit keywords: %s, package name: %s", keyword, str))
			return true
		}
	}
	return false
}

func (t *SSHelper) ExportBySpecify(pathList ...string) error {
	if len(pathList) == 0 {
		return fmt.Errorf("please input path")
	}
	err := t.init()
	if err != nil {
		return err
	}
	defer t.sshClient.Close()
	defer t.sftpClient.Close()
	defer t.spinner.Quit()
	defer t.spinner.Msg("export by specify path is done.")
	for _, _path := range pathList {
		//destinationDir := filepath.Join(t.output, filepath.Dir(_path))
		//t.spinner.Msg(fmt.Sprintf("decompressing %s", destinationDir))
		//if pathx.PathExists(destinationDir) {
		//	osx.RemoveAll(destinationDir)
		//}
		err := t.export(_path)
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func (t *SSHelper) export(remotePath string) error {
	var err error
	remoteTarPath := pathx.PathJoin(filepath.Dir(remotePath), ActiveFileName)
	t.spinner.Msg(fmt.Sprintf("packing %s", remotePath))
	_, err = t.sshClient.ExecuteCommand(fmt.Sprintf("tar --ignore-failed-read -cf %s %s", remoteTarPath, remotePath))
	if err != nil {
		return err
	}
	defer func() {
		t.sshClient.DeleteFiles(remoteTarPath)
	}()
	localTarPath := filepath.Join(t.output, ActiveFileName)
	t.spinner.Msg(fmt.Sprintf("downloading %s", remotePath))
	err = t.sftpClient.DownloadFile(remoteTarPath, localTarPath)
	if err != nil {
		return err
	}
	defer func() {
		osx.RemoveAll(localTarPath)
	}()
	t.spinner.Msg(fmt.Sprintf("decompressing %s", remotePath))
	_, err = t.sshClient.ExecuteCommand(fmt.Sprintf("rm -f %s", remoteTarPath))
	if err != nil {
		return err
	}
	destinationDir := filepath.Join(t.output, filepath.Dir(remotePath))
	err = os.MkdirAll(destinationDir, 0755)
	if err != nil {
		return err
	}
	err = osx.TarDecompression(localTarPath, t.output)
	if err != nil {
		return err
	}
	return nil
}
