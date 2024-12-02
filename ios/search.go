package ios

import (
	"fmt"
	"github.com/solywsh/go-forensic/utils/osx"
	"github.com/solywsh/go-forensic/utils/pathx"
	"howett.net/plist"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (t *IOS) SearchByKeywords(keywords ...string) error {
	sshClient, err := t.ssh().C()
	if err != nil {
		return err
	}
	defer sshClient.Close()
	sftpClient := sshClient.SFTP()
	if sftpClient == nil {
		return nil
	}
	defer sftpClient.Close()
	t.sshClient = sshClient
	t.sftpClient = sftpClient
	t.keywords = keywords
	for _, focusPath := range focusPathList {
		err = t.handleFiles(focusPath)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (t *IOS) handleFiles(tryPath string) error {
	tryFileList, err := t.sftpClient.ReadDir(tryPath)
	if err != nil {
		return err
	}
	for _, containerDir := range tryFileList {
		if !containerDir.IsDir() {
			continue
		}
		err := t.handle(tryPath, containerDir)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (t *IOS) handle(tryPath string, containerFile os.FileInfo) error {
	containerFileList, err := t.sftpClient.ReadDir(pathx.PathJoin(tryPath, containerFile.Name()))
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
	containerMetaDataFilePath := pathx.PathJoin(tryPath, containerFile.Name(), ContainerMetaDataFile)
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
	fmt.Println("hit: ", packageInfo.PackageName)
	remoteDir := pathx.PathJoin(tryPath, containerFile.Name())
	remoteTarPath := pathx.PathJoin(tryPath, ActiveFileName)
	_, err = t.sshClient.ExecuteCommand(fmt.Sprintf("tar --ignore-failed-read -cf %s %s", remoteTarPath, remoteDir))
	if err != nil {
		return err
	}
	defer func() {
		t.sshClient.DeleteFiles(remoteTarPath)
	}()
	localTarPath := filepath.Join(t.output, ActiveFileName)
	err = t.sftpClient.DownloadFile(remoteTarPath, localTarPath)
	if err != nil {
		return err
	}
	defer func() {
		os.RemoveAll(localTarPath)
	}()
	destinationDir := filepath.Join(t.output, tryPath, containerFile.Name())
	if pathx.PathExists(destinationDir) {
		os.RemoveAll(destinationDir)
	}
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

func (t *IOS) ComparisonKeywords(str string) bool {
	for _, keyword := range t.keywords {
		if strings.Contains(strings.ToLower(str), strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}
