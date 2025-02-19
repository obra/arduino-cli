// This file is part of arduino-cli.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package resources

import (
	"fmt"
	"os"

	"github.com/arduino/go-paths-helper"
	"go.bug.st/downloader/v2"
)

// ArchivePath returns the path of the Archive of the specified DownloadResource relative
// to the specified downloadDir
func (r *DownloadResource) ArchivePath(downloadDir *paths.Path) (*paths.Path, error) {
	staging := downloadDir.Join(r.CachePath)
	if err := staging.MkdirAll(); err != nil {
		return nil, err
	}
	return staging.Join(r.ArchiveFileName), nil
}

// IsCached returns true if the specified DownloadResource has already been downloaded
func (r *DownloadResource) IsCached(downloadDir *paths.Path) (bool, error) {
	archivePath, err := r.ArchivePath(downloadDir)
	if err != nil {
		return false, fmt.Errorf(tr("getting archive path: %s"), err)
	}
	return archivePath.Exist(), nil
}

// Download a DownloadResource.
func (r *DownloadResource) Download(downloadDir *paths.Path, config *downloader.Config) (*downloader.Downloader, error) {
	path, err := r.ArchivePath(downloadDir)
	if err != nil {
		return nil, fmt.Errorf(tr("getting archive path: %s"), err)
	}

	if _, err := path.Stat(); os.IsNotExist(err) {
		// normal download
	} else if err == nil {
		// check local file integrity
		ok, err := r.TestLocalArchiveIntegrity(downloadDir)
		if err != nil || !ok {
			if err := path.Remove(); err != nil {
				return nil, fmt.Errorf(tr("removing corrupted archive file: %s"), err)
			}
		} else {
			// File is cached, nothing to do here
			return nil, nil
		}
	} else if err == nil {
		// resume download
	} else {
		return nil, fmt.Errorf(tr("getting archive file info: %s"), err)
	}

	return downloader.DownloadWithConfig(path.String(), r.URL, *config)
}
