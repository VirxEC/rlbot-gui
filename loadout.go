package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	rlbot "github.com/RLBot/go-interface"
	"github.com/RLBot/go-interface/flat"
)

func (a *App) SaveLoadoutToFile(basePath string, loadoutFile string, loadout LoadoutConfig) error {
	baseDir := filepath.Dir(basePath)
	fullPath := filepath.Join(baseDir, loadoutFile)

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileContents, err := toml.Marshal(loadout)
	if err != nil {
		return err
	}

	_, err = file.Write(fileContents)
	return err
}

type LoadoutPreviewOptions struct {
	Map         string            `json:"map"`
	Loadout     TeamLoadoutConfig `json:"loadout"`
	Team        uint32            `json:"team"`
	Launcher    string            `json:"launcher"`
	LauncherArg string            `json:"launcherArg"`
}

func (options LoadoutPreviewOptions) GetPreviewMatch(existingMatchBehavior flat.ExistingMatchBehavior) (*flat.MatchConfigurationT, error) {
	loadout := options.Loadout.ToPlayerLoadout()

	playerConfigs := []*flat.PlayerConfigurationT{
		{
			Variety: &flat.PlayerClassT{
				Type: flat.PlayerClassCustomBot,
				Value: &flat.CustomBotT{
					Name:       "Showcase",
					AgentId:    "gui/loadout-preview",
					RootDir:    "",
					RunCommand: "",
					Loadout:    loadout,
					Hivemind:   false,
				},
			},
			Team:     options.Team,
			PlayerId: 0,
		},
	}

	var launcher flat.Launcher
	switch options.Launcher {
	case "steam":
		launcher = flat.LauncherSteam
	case "epic":
		launcher = flat.LauncherEpic
	case "legendary", "heroic":
		launcher = flat.LauncherCustom
		options.LauncherArg = options.Launcher
	case "custom":
		launcher = flat.LauncherCustom
	case "nolaunch":
		launcher = flat.LauncherNoLaunch
	default:
		return nil, errors.New("no launcher specified")
	}

	return &flat.MatchConfigurationT{
		AutoStartAgents:      false,
		WaitForAgents:        false,
		GameMapUpk:           options.Map,
		PlayerConfigurations: playerConfigs,
		ScriptConfigurations: []*flat.ScriptConfigurationT{},
		GameMode:             flat.GameModeSoccar,
		Mutators: &flat.MutatorSettingsT{
			MatchLength: flat.MatchLengthMutatorUnlimited,
			BoostAmount: flat.BoostAmountMutatorUnlimitedBoost,
		},
		Freeplay:              false,
		EnableRendering:       flat.DebugRenderingAlwaysOff,
		EnableStateSetting:    true,
		InstantStart:          true,
		SkipReplays:           true,
		AutoSaveReplay:        false,
		Launcher:              launcher,
		LauncherArg:           options.LauncherArg,
		ExistingMatchBehavior: existingMatchBehavior,
	}, nil
}

func (a *App) LaunchPreviewLoadout(options LoadoutPreviewOptions, existingMatchBehavior flat.ExistingMatchBehavior) error {
	match, err := options.GetPreviewMatch(existingMatchBehavior)
	if err != nil {
		return err
	}

	return StartAndWaitForMatch(a.rlbotAddress, match)
}

func WaitForGamePacket(conn *rlbot.RLBotConnection) (*flat.GamePacketT, error) {
	var gamePacket *flat.GamePacketT
	for gamePacket == nil || (gamePacket.MatchInfo.MatchPhase != flat.MatchPhaseKickoff && gamePacket.MatchInfo.MatchPhase != flat.MatchPhaseActive) {
		packet, err := conn.RecvPacket()
		if err != nil {
			return nil, err
		}

		switch packet := packet.Value.(type) {
		case *flat.GamePacketT:
			gamePacket = packet
		case *flat.DisconnectSignalT:
			return nil, fmt.Errorf("received disconnect signal while waiting for game packet")
		}
	}

	return gamePacket, nil
}

func (a *App) SetLoadout(options LoadoutPreviewOptions) error {
	conn, err := rlbot.Connect(a.rlbotAddress)
	if err != nil {
		return err
	}

	conn.SendPacket(&flat.ConnectionSettingsT{
		AgentId:              "",
		WantsBallPredictions: false,
		WantsComms:           false,
		CloseBetweenMatches:  false,
	})

	conn.SendPacket(&flat.InitCompleteT{})

	gamePacket, err := WaitForGamePacket(&conn)
	if err != nil {
		return err
	}

	// if the match is over, launch a new match
	isMatchOver := gamePacket.MatchInfo.MatchPhase == flat.MatchPhaseEnded || gamePacket.MatchInfo.MatchPhase == flat.MatchPhaseInactive
	// ensure the player is a custom bot
	isPlayerBot := len(gamePacket.Players) == 1 && gamePacket.Players[0].IsBot
	// ensure unlimited time
	isUnlimitedTime := gamePacket.MatchInfo.IsUnlimitedTime

	if isMatchOver || !isPlayerBot || !isUnlimitedTime {
		match, err := options.GetPreviewMatch(flat.ExistingMatchBehaviorRestart)
		if err != nil {
			return err
		}

		conn.SendPacket(match)
		conn.SendPacket(&flat.DisconnectSignalT{})
		return nil
	}

	// ensure the player is on the correct team
	if gamePacket.Players[0].Team != options.Team {
		match, err := options.GetPreviewMatch(flat.ExistingMatchBehaviorContinueAndSpawn)
		if err != nil {
			return err
		}

		conn.SendPacket(match)
		conn.SendPacket(&flat.DisconnectSignalT{})
		return nil
	}

	conn.SendPacket(&flat.SetLoadoutT{
		Index:   0,
		Loadout: options.Loadout.ToPlayerLoadout(),
	})

	conn.SendPacket(&flat.DisconnectSignalT{})

	return nil
}

func StaticSetter(rlbotAddress string, team uint32) error {
	conn, err := rlbot.Connect(rlbotAddress)
	if err != nil {
		return err
	}

	conn.SendPacket(&flat.ConnectionSettingsT{
		AgentId:              "",
		WantsBallPredictions: false,
		WantsComms:           false,
		CloseBetweenMatches:  true,
	})

	conn.SendPacket(&flat.InitCompleteT{})

	gameState := flat.DesiredGameStateT{
		CarStates: []*flat.DesiredCarStateT{
			{
				Physics: &flat.DesiredPhysicsT{
					Location:        Vector3P(0, 0, 20),
					Rotation:        RotatorP(0, 0, 0),
					Velocity:        Vector3P(0, 0, 0),
					AngularVelocity: Vector3P(0, 0, 0),
				},
			},
		},
	}

	for {
		packet, err := conn.RecvPacket()
		if err != nil {
			return err
		}

		switch packet.Value.(type) {
		case *flat.DisconnectSignalT:
			return nil
		case *flat.GamePacketT:
			conn.SendPacket(&gameState)
		}
	}
}

func (a *App) SetShowcaseType(showcaseType string, team uint32) error {
	conn, err := rlbot.Connect(a.rlbotAddress)
	if err != nil {
		return err
	}

	conn.SendPacket(&flat.ConnectionSettingsT{
		AgentId:              "",
		WantsBallPredictions: false,
		WantsComms:           false,
		CloseBetweenMatches:  false,
	})

	conn.SendPacket(&flat.InitCompleteT{})

	gamePacket, err := WaitForGamePacket(&conn)
	if err != nil {
		return err
	}

	ball := flat.DesiredPhysicsT{
		Location:        Vector3P(0, 0, -100),
		Velocity:        Vector3P(0, 0, 0),
		AngularVelocity: Vector3P(0, 0, 0),
	}

	car := flat.DesiredPhysicsT{
		Location:        Vector3P(0, 0, 20),
		Rotation:        RotatorP(0, 0, 0),
		Velocity:        Vector3P(0, 0, 0),
		AngularVelocity: Vector3P(0, 0, 0),
	}

	controller := flat.ControllerStateT{}

	var teamSign float32
	if gamePacket.Players[0].Team == 0 {
		teamSign = -1
	} else {
		teamSign = 1
	}

	// set initial game state
	switch showcaseType {
	case "static":
		controller.Boost = true

		go StaticSetter(a.rlbotAddress, team)
	case "boost":
		controller.Boost = true
		controller.Steer = 1

		car.Location.Y = Float(-1140)
		car.Velocity.X = Float(2300)
		car.AngularVelocity.Z = Float(3.5)
	case "throttle":
		controller.Throttle = 1
		controller.Steer = 0.56

		car.Location.Y = Float(-1140)
		car.Velocity.X = Float(1410)
		car.AngularVelocity.Z = Float(1.5)
	case "back-center-kickoff":
		car.Location.Y = Float(4608 * teamSign)
		car.Rotation.Yaw = Float(-0.5 * 3.14159 * teamSign)
	case "goal-explosion":
		car.Location.Y = Float(-2000 * teamSign)
		car.Rotation.Yaw = Float(-0.5 * 3.14159 * teamSign)
		car.Velocity.Y = Float(-2300 * teamSign)
		ball.Location = Vector3P(0, -3500*teamSign, 93)
	}

	conn.SendPacket(&flat.DesiredGameStateT{
		CarStates: []*flat.DesiredCarStateT{
			{
				Physics: &car,
			},
		},
		BallStates: []*flat.DesiredBallStateT{
			{
				Physics: &ball,
			},
		},
	})

	conn.SendPacket(&flat.PlayerInputT{
		PlayerIndex:     0,
		ControllerState: &controller,
	})

	conn.SendPacket(&flat.DisconnectSignalT{})

	return nil
}

func Float(x float32) *flat.FloatT {
	return &flat.FloatT{Val: x}
}

func Vector3P(x, y, z float32) *flat.Vector3PartialT {
	return &flat.Vector3PartialT{
		X: Float(x),
		Y: Float(y),
		Z: Float(z),
	}
}

func RotatorP(p, y, r float32) *flat.RotatorPartialT {
	return &flat.RotatorPartialT{
		Pitch: Float(p),
		Yaw:   Float(y),
		Roll:  Float(r),
	}
}
