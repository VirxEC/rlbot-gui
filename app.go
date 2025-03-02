package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ncruces/zenity"
	rlbot "github.com/swz-git/go-interface"
	"github.com/swz-git/go-interface/flat"
)

type RawReleaseInfo struct {
	repo    string
	content GhRelease
}

// App struct
type App struct {
	latest_release_json []RawReleaseInfo
	rlbot_address       string
}

func (a *App) IgnoreMe(
	_ BotInfo,
	_ PsyonixBotInfo,
	_ HumanInfo,
) {
}

func (a *App) GetDefaultPath() string {
	if runtime.GOOS == "windows" {
		localappdata := os.Getenv("LOCALAPPDATA")
		return filepath.Join(localappdata, "RLBotGUI")
	}

	// assume linux

	xdg_data_home := os.Getenv("XDG_DATA_HOME")
	if xdg_data_home == "" {
		home := os.Getenv("HOME")
		xdg_data_home = filepath.Join(home, ".local/share")
	}

	return filepath.Join(xdg_data_home, "rlbotgui")
}

func (a *App) GetLatestReleaseData(repo string) (*GhRelease, error) {
	latest_release_url := "https://api.github.com/repos/" + repo + "/releases/latest"

	resp, err := http.Get(latest_release_url)
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

	a.latest_release_json = append(a.latest_release_json, RawReleaseInfo{repo, content})

	return &a.latest_release_json[len(a.latest_release_json)-1].content, nil
}

func (a *App) DownloadBotpack(repo string, installPath string) (string, error) {
	var latest_release *GhRelease

	for _, release := range a.latest_release_json {
		if release.repo == repo {
			latest_release = &release.content
			break
		}
	}

	if latest_release == nil {
		content, err := a.GetLatestReleaseData(repo)
		if err != nil {
			return "", err
		}

		latest_release = content
	}

	var file_name string
	if runtime.GOOS == "windows" {
		file_name = "botpack_x86_64-windows.tar.xz"
	} else {
		file_name = "botpack_x86_64-linux.tar.xz"
	}

	var download_url string
	for _, asset := range latest_release.Assets {
		if asset.Name == file_name {
			download_url = asset.BrowserDownloadURL
			break
		}
	}

	err := DownloadExtractArchive(download_url, installPath)
	if err != nil {
		return "", err
	}

	return latest_release.TagName, nil
}

// NewApp creates a new App application struct
func NewApp() *App {
	ip := os.Getenv("RLBOT_SERVER_IP")
	if ip == "" {
		ip = "127.0.0.1"
	}

	port := os.Getenv("RLBOT_SERVER_PORT")
	if port == "" {
		port = "23234"
	}

	rlbot_address := ip + ":" + port

	var latest_release_json []RawReleaseInfo

	return &App{
		latest_release_json,
		rlbot_address,
	}
}

func recursiveTomlSearch(root, tomlType string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(info.Name()) != ".toml" {
			return nil
		}

		if info.Name() == tomlType+".toml" || strings.HasSuffix(info.Name(), "."+tomlType+".toml") {
			matches = append(matches, path)
		}

		return nil
	})
	return matches, err
}

type Result struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ExtraOptions struct {
	EnableRendering       bool `json:"enableRendering"`
	EnableStateSetting    bool `json:"enableStateSetting"`
	InstantStart          bool `json:"instantStart"`
	SkipReplays           bool `json:"skipReplays"`
	AutoSaveReplay        bool `json:"autoSaveReplay"`
	ExistingMatchBehavior byte `json:"existingMatchBehavior"`
}

type StartMatchOptions struct {
	Map             string                `json:"map"`
	GameMode        string                `json:"gameMode"`
	BluePlayers     []PlayerJs            `json:"bluePlayers"`
	OrangePlayers   []PlayerJs            `json:"orangePlayers"`
	MutatorSettings flat.MutatorSettingsT `json:"mutatorSettings"`
	ExtraOptions    ExtraOptions          `json:"extraOptions"`
	Launcher        string                `json:"launcher"`
	LauncherArg     string                `json:"launcherArg"`
}

func (a *App) StartMatch(options StartMatchOptions) Result {
	// TODO: Save this in App struct
	conn, err := rlbot.Connect(a.rlbot_address)
	if err != nil {
		return Result{false, "Failed to connect to RLBotServer at " + a.rlbot_address}
	}

	var gameMode flat.GameMode
	switch options.GameMode {
	case "Soccer":
		gameMode = flat.GameModeSoccer
	case "Hoops":
		gameMode = flat.GameModeHoops
	case "Dropshot":
		gameMode = flat.GameModeDropshot
	case "Hockey":
		gameMode = flat.GameModeHockey
	case "Rumble":
		gameMode = flat.GameModeRumble
	case "Heatseeker":
		gameMode = flat.GameModeHeatseeker
	case "Gridiron":
		gameMode = flat.GameModeGridiron
	case "Knockout":
		gameMode = flat.GameModeKnockout
	default:
		println("No mode chosen, defaulting to soccer")
		gameMode = flat.GameModeSoccer
	}

	var launcher flat.Launcher
	switch options.Launcher {
	case "steam":
		launcher = flat.LauncherSteam
	case "epic":
		launcher = flat.LauncherEpic
	case "custom":
		launcher = flat.LauncherCustom
	case "noLaunch":
		launcher = flat.LauncherNoLaunch
	default:
		println("No launcher chosen, defaulting to NoLaunch")
		launcher = flat.LauncherNoLaunch
	}

	var playerConfigs []*flat.PlayerConfigurationT

	for _, playerInfo := range options.BluePlayers {
		println(playerInfo.ToPlayer().ToPlayerConfig(0))
		playerConfigs = append(playerConfigs, playerInfo.ToPlayer().ToPlayerConfig(0))
	}
	for _, playerInfo := range options.OrangePlayers {
		playerConfigs = append(playerConfigs, playerInfo.ToPlayer().ToPlayerConfig(1))
	}

	conn.SendPacket(&flat.MatchConfigurationT{
		AutoStartBots:         true,
		GameMapUpk:            options.Map,
		PlayerConfigurations:  playerConfigs,
		GameMode:              gameMode,
		Mutators:              &options.MutatorSettings,
		EnableRendering:       options.ExtraOptions.EnableRendering,
		EnableStateSetting:    options.ExtraOptions.EnableStateSetting,
		InstantStart:          options.ExtraOptions.InstantStart,
		SkipReplays:           options.ExtraOptions.SkipReplays,
		AutoSaveReplay:        options.ExtraOptions.AutoSaveReplay,
		Launcher:              launcher,
		LauncherArg:           options.LauncherArg,
		ExistingMatchBehavior: flat.ExistingMatchBehavior(options.ExtraOptions.ExistingMatchBehavior),
	})

	return Result{true, ""}
}

func (a *App) StopMatch(shutdownServer bool) Result {
	// TODO: Save this in App struct
	// TODO: Make dynamic, pull from env var?
	conn, err := rlbot.Connect("127.0.0.1:23234")
	if err != nil {
		return Result{false, "Failed to connect to rlbot"}
	}

	conn.SendPacket(&flat.StopCommandT{
		ShutdownServer: shutdownServer,
	})

	return Result{true, ""}
}

func (a *App) PickFolder() string {
	path, err := zenity.SelectFile(zenity.Directory())
	if err != nil {
		println("ERR: File picker failed")
	}
	return path
}

func (a *App) ShowPathInExplorer(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	// if is dir
	var folder string
	if fileInfo.IsDir() {
		folder = path
	} else {
		folder = filepath.Dir(path)
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("explorer", folder)
		err := cmd.Run()
		if err != nil {
			return err
		}
	} else {
		cmd := exec.Command("xdg-open", folder)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
