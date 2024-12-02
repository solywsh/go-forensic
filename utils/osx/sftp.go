package osx

import (
	"github.com/pkg/sftp"
	"io"
	"os"
	"path/filepath"
)

type SftpX struct {
	*SShX
	sftpClient *sftp.Client
}

func (s *SShX) SFTP() *SftpX {
	if s.client == nil {
		return nil
	}
	sftpClient, err := sftp.NewClient(s.client)
	if err != nil {
		return nil
	}
	return &SftpX{
		SShX:       s,
		sftpClient: sftpClient,
	}
}

func (s *SftpX) Close() error {
	if s.sftpClient == nil {
		return nil
	}
	return s.sftpClient.Close()
}

func (s *SftpX) ReadDir(path string) ([]os.FileInfo, error) {
	if s.sftpClient == nil {
		return nil, nil
	}
	return s.sftpClient.ReadDir(path)
}

func (s *SftpX) Open(path string) (*sftp.File, error) {
	if s.sftpClient == nil {
		return nil, nil
	}
	return s.sftpClient.Open(path)
}

func (s *SftpX) Client() *sftp.Client {
	if s.sftpClient == nil {
		return nil
	}
	return s.sftpClient
}

func (s *SftpX) DownloadFile(remoteFile, localFile string) error {
	remoteFileHandle, err := s.sftpClient.Open(remoteFile)
	if err != nil {
		return err
	}
	defer remoteFileHandle.Close()
	err = os.MkdirAll(filepath.Dir(localFile), 0755)
	if err != nil {
		return err
	}
	localFileHandle, err := os.Create(localFile)
	if err != nil {
		return err
	}
	defer localFileHandle.Close()

	_, err = io.Copy(localFileHandle, remoteFileHandle)
	if err != nil {
		return err
	}
	return nil

}
