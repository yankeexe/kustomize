/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package loader

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-getter"

	"github.com/kubernetes-sigs/kustomize/pkg/fs"
)

// githubLoader loads files from a checkout github repo
type githubLoader struct {
	repo        string
	checkoutDir string
	loader      *fileLoader
}

// Root returns the root location for this Loader.
func (l *githubLoader) Root() string {
	return l.checkoutDir
}

// New delegates to fileLoader.New
func (l *githubLoader) New(newRoot string) (Loader, error) {
	return l.loader.New(newRoot)
}

// Load delegates to fileLoader.Load
func (l *githubLoader) Load(location string) ([]byte, error) {
	return l.loader.Load(location)
}

// Cleanup removes the checked out repo
func (l *githubLoader) Cleanup() error {
	return os.RemoveAll(l.checkoutDir)
}

// newGithubLoader returns a new fileLoader with given github Url.
func newGithubLoader(repoUrl string, fs fs.FileSystem) (*githubLoader, error) {
	dir, err := ioutil.TempDir("", "kustomize-")
	if err != nil {
		return nil, err
	}
	target := filepath.Join(dir, "repo")
	err = checkout(repoUrl, target)
	if err != nil {
		return nil, err
	}
	l := newFileLoaderAtRoot(target, fs)
	return &githubLoader{
		repo:        repoUrl,
		checkoutDir: dir,
		loader:      l,
	}, nil
}

// isRepoUrl checks if a string is a repo Url
func isRepoUrl(s string) bool {
	return strings.Contains(s, ".com") || strings.Contains(s, ".org") || strings.Contains(s, "https://")
}

// Checkout clones a github repo with specified commit/tag/branch
func checkout(url, dir string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	client := &getter.Client{
		Src:  url,
		Dst:  dir,
		Pwd:  pwd,
		Mode: getter.ClientModeDir,
	}
	return client.Get()
}
