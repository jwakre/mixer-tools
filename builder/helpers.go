// Copyright © 2018 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package builder

import (
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/clearlinux/mixer-tools/helpers"
	"github.com/pkg/errors"
)

func parseUint32(s string) (uint32, error) {
	parsed, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, errors.Wrapf(err, "error parsing value %q", s)
	}
	return uint32(parsed), nil
}

// buildUpstreamURL builds the full upstream URL based on a b.UpstreamURL and a
// supplied subpath
func (b *Builder) buildUpstreamURL(subpath string) (string, error) {
	base, err := url.Parse(b.UpstreamURL)
	if err != nil {
		return "", err
	}

	for _, token := range strings.Split(subpath, "/") {
		base.Path = path.Join(base.Path, token)
	}

	return base.String(), nil
}

// cleanMix removes image and output (www) subdirectories that were generated
// for mixes greater than or equal to the current mixversion
func (b *Builder) cleanMix() error {
	mixVer, err := strconv.Atoi(b.MixVer)
	if err != nil {
		return err
	}

	outputDir := filepath.Join(b.Config.Builder.ServerStateDir, "www")
	if err = cleanMixDir(outputDir, mixVer); err != nil {
		return err
	}

	imageDir := filepath.Join(b.Config.Builder.ServerStateDir, "image")
	return cleanMixDir(imageDir, mixVer)
}

// cleanMixDir removes cleanDir's subdirectories that are numbered greater
// than or equal to cleanVer. This function is intended to be used to clean
// the output (www) or image directories.
func cleanMixDir(cleanDir string, cleanVer int) error {
	if _, err := os.Stat(cleanDir); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(cleanDir)
	if err != nil {
		return err
	}

	// Remove subdirectories numbered greater than or equal to clearVer
	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		dirVer, err := strconv.Atoi(f.Name())
		if err != nil {
			continue
		}

		if dirVer >= cleanVer {
			if err = os.RemoveAll(filepath.Join(cleanDir, f.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}

// DownloadFileFromUpstreamAsString will download a file from the Upstream URL
// joined with the passed subpath. It will trim leading and trailing whitespace
// from the result.
func (b *Builder) DownloadFileFromUpstreamAsString(subpath string) (string, error) {
	if b.UpstreamURL == "" {
		return b.State.Mix.Format, nil
	}
	url, err := b.buildUpstreamURL(subpath)
	if err != nil {
		return "", err
	}
	content, err := helpers.DownloadFileAsString(url)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

// DownloadFileFromUpstream will download a file from the Upstream URL
// joined with the passed subpath and write that file to the supplied file path.
// If the path is left empty, the file name will be inferred from the source
// and written to PWD.
func (b *Builder) DownloadFileFromUpstream(subpath string, filePath string) error {
	url, err := b.buildUpstreamURL(subpath)
	if err != nil {
		return err
	}
	return helpers.DownloadFile(url, filePath)
}

// TerminalWidth determines the screen width of the calling terminal.
func TerminalWidth() (int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	outStrs := strings.Fields(string(out))
	if len(outStrs) != 2 {
		return 0, errors.Errorf("Invalid stty output")
	}
	width, err := strconv.Atoi(outStrs[1])
	if err != nil {
		return 0, err
	}
	return width, nil
}
