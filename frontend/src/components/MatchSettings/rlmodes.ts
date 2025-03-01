import { mutators } from "./rlmutators"

export const gamemodes: { [x: string]: { match: { [x: string]: string }, mutators: { [x: string]: string } } } = {
  "Heatseeker Ricochet": {
    match: {
      game_mode: mutators.game_mode[5],
      game_map_upk: "Labs_PillarGlass_P",
    },
    mutators: {},
  },
  "Gotham City Rumble": {
    match: {
      game_mode: mutators.game_mode[0],
      game_map_upk: "Park_Bman_P",
    },
    mutators: {
      rumble_option: mutators.rumble_option[10],
    },
  },
  "Speed Demon": {
    match: {
      game_mode: mutators.game_mode[0],
    },
    mutators: {
      boost_amount_option: mutators.boost_amount_option[1],
      boost_strength_option: mutators.boost_strength_option[2],
      ball_max_speed_option: mutators.ball_max_speed_option[3],
      ball_bounciness_option: mutators.ball_bounciness_option[1],
      ball_size_option: mutators.ball_size_option[3],
      respawn_time_option: mutators.respawn_time_option[2],
      demolish_option: mutators.demolish_option[3],
    },
  },
  "Ghost Hunt": {
    match: {
      game_mode: mutators.game_mode[0],
      game_map_upk: "Haunted_TrainStation_P",
    },
    mutators: {
      ball_type_option: mutators.ball_type_option[6],
      rumble_option: mutators.rumble_option[8],
      game_event_option: mutators.game_event_option[1],
      audio_option: mutators.audio_option[1],
    },
  },
  "Dropshot Rumble": {
    match: {
      game_mode: mutators.game_mode[2],
      game_map_upk: "ShatterShot_P",
    },
    mutators: {
      rumble_option: mutators.rumble_option[1],
    },
  },
  "Nike Fc Showdown": {
    match: {
      game_mode: mutators.game_mode[0],
      game_map_upk: "swoosh_p",
    },
    mutators: {
      ball_max_speed_option: mutators.ball_max_speed_option[2],
      ball_weight_option: mutators.ball_weight_option[6],
      ball_bounciness_option: mutators.ball_bounciness_option[4],
      ball_type_option: mutators.ball_type_option[7],
    },
  },
  "Tactical Rumble": {
    match: {
      game_mode: mutators.game_mode[4],
    },
    mutators: {
      rumble_option: mutators.rumble_option[9],
    },
  },
  "Gforce Frenzy": {
    match: {
      game_mode: mutators.game_mode[0],
    },
    mutators: {
      boost_amount_option: mutators.boost_amount_option[1],
      boost_strength_option: mutators.boost_strength_option[3],
      gravity_option: mutators.gravity_option[1],
    },
  },
  "Spike Rush": {
    match: {
      game_mode: mutators.game_mode[0],
      game_map_upk: "ThrowbackStadium_P",
    },
    mutators: {
      rumble_option: mutators.rumble_option[7],
      respawn_time_option: mutators.respawn_time_option[2],
      game_event_option: mutators.game_event_option[2],
    },
  },
  "Spring Loaded": {
    match: {},
    mutators: {
      rumble_option: mutators.rumble_option[5],
    },
  },
  "Spooky Cube": {
    match: {
      game_mode: mutators.game_mode[0],
    },
    mutators: {
      ball_max_speed_option: mutators.ball_max_speed_option[3],
      ball_type_option: mutators.ball_type_option[8],
      ball_weight_option: mutators.ball_weight_option[1],
      ball_bounciness_option: mutators.ball_bounciness_option[2],
      boost_amount_option: mutators.boost_amount_option[1],
    },
  },
  "Beach Ball": {
    match: {
      game_mode: mutators.game_mode[0],
    },
    mutators: {
      ball_max_speed_option: mutators.ball_max_speed_option[2],
      ball_type_option: mutators.ball_type_option[4],
      ball_weight_option: mutators.ball_weight_option[5],
      ball_size_option: mutators.ball_size_option[2],
      ball_bounciness_option: mutators.ball_bounciness_option[2],
    },
  },
  "Super Cube": {
    match: {
      game_mode: mutators.game_mode[0],
    },
    mutators: {
      ball_max_speed_option: mutators.ball_max_speed_option[3],
      ball_type_option: mutators.ball_type_option[1],
      ball_weight_option: mutators.ball_weight_option[1],
      ball_bounciness_option: mutators.ball_bounciness_option[2],
      boost_amount_option: mutators.boost_amount_option[1],
    },
  },
  "Boomer Ball": {
    match: {
      game_mode: mutators.game_mode[0],
    },
    mutators: {
      boost_amount_option: mutators.boost_amount_option[1],
      boost_strength_option: mutators.boost_strength_option[1],
      ball_max_speed_option: mutators.ball_max_speed_option[3],
      ball_bounciness_option: mutators.ball_bounciness_option[2],
      ball_weight_option: mutators.ball_weight_option[3],
    },
  },
}
