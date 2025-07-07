package main

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ulikunitz/xz"
)

type GhRelease struct {
	URL             string    `json:"url"`
	AssetsURL       string    `json:"assets_url"`
	UploadURL       string    `json:"upload_url"`
	HTMLURL         string    `json:"html_url"`
	ID              int       `json:"id"`
	Author          GhAuthor  `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []GhAsset `json:"assets"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	Body            string    `json:"body"`
}

type GhAuthor struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type GhAsset struct {
	URL                string    `json:"url"`
	ID                 int       `json:"id"`
	NodeID             string    `json:"node_id"`
	Name               string    `json:"name"`
	Label              string    `json:"label"`
	Uploader           GhAuthor  `json:"uploader"`
	ContentType        string    `json:"content_type"`
	State              string    `json:"state"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	BrowserDownloadURL string    `json:"browser_download_url"`
}

func ParseReleaseData(data []byte) (GhRelease, error) {
	var release GhRelease
	err := json.Unmarshal(data, &release)
	if err != nil {
		return GhRelease{}, err
	}

	return release, nil
}

// taken from https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func extractTar(tr *tar.Reader, dst string) error {
	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

func DownloadExtractArchive(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// body is a .tar.xz file
	xzr, err := xz.NewReader(bytes.NewReader(body))
	if err != nil {
		return err
	}

	tr := tar.NewReader(xzr)
	err = extractTar(tr, dest)
	return err
}

func (a *App) GetLatestReleaseData(repo string) (*GhRelease, error) {
	for _, release := range a.latestReleaseJson {
		if release.repo == repo {
			return &release.content, nil
		}
	}

	latestReleaseUrl := "https://api.github.com/repos/" + repo + "/releases/latest"

	resp, err := http.Get(latestReleaseUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	content, err := ParseReleaseData(body)
	if err != nil {
		return nil, err
	}

	a.latestReleaseJson = append(a.latestReleaseJson, RawReleaseInfo{repo, content})

	return &a.latestReleaseJson[len(a.latestReleaseJson)-1].content, nil
}

func (a *App) UpdateBotpack(repo string, installPath string, currentTag string) (string, error) {
	latestRelease, err := a.GetLatestReleaseData(repo)
	if err != nil {
		return "", err
	}

	newestTag := latestRelease.TagName
	if newestTag == currentTag {
		return "", errors.New("already up to date")
	}

	currentVersion := strings.Split(currentTag, "-")[1]
	currentVersionNum, err := strconv.Atoi(currentVersion)
	if err != nil {
		return "", err
	}

	newestVersion := strings.Split(newestTag, "-")[1]
	newestVersionNum, err := strconv.Atoi(newestVersion)
	if err != nil {
		return "", err
	}

	var file_name string
	// todo: check platform (x86_64, etc)
	if runtime.GOOS == "windows" {
		file_name = "patch_x86_64-windows.bobdiff"
	} else {
		file_name = "patch_x86_64-linux.bobdiff"
	}

	var latestDownloadUrl string
	for _, asset := range latestRelease.Assets {
		if asset.Name == file_name {
			latestDownloadUrl = asset.BrowserDownloadURL
			break
		}
	}

	for i := currentVersionNum + 1; i <= newestVersionNum; i++ {
		tagI := strings.Replace(currentTag, currentVersion, strconv.Itoa(i), 1)
		downloadUrl := strings.Replace(latestDownloadUrl, newestTag, tagI, 1)

		resp, err := http.Get(downloadUrl)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		bytes, err := io.ReadAll(bytes.NewReader(body))
		if err != nil {
			return "", err
		}

		files, err := os.ReadDir(installPath)
		if err != nil {
			return "", err
		}

		var dir string
		for _, file := range files {
			if file.IsDir() {
				dir = file.Name()
				break
			}
		}
		if dir == "" {
			return "", errors.New("no directory found")
		}

		err = diffApply(filepath.Join(installPath, dir), bytes)
		if err != nil {
			return "", err
		}
	}

	return latestRelease.TagName, nil
}
