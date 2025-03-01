<script lang="ts">
    import Select from "../NiceSelect.svelte";
    import Modal from "../Modal.svelte";
    import { MAPS_NON_STANDARD, MAPS_STANDARD } from "../../arena-names";
    import { mutators as mutatorOptions } from "./rlmutators";
    import LauncherSelector from "../LauncherSelector.svelte";
    import { gamemodes } from "./rlmodes";

    let {
        map = $bindable(),
        mode = $bindable(),
        extraOptions = $bindable(),
        mutators = $bindable(),
        onStart = (randomizeMap: boolean) => {},
        onStop = () => {},
    } = $props();
    let showExtraOptions = $state(false);
    let showMutators = $state(false);
    let randomizeMap = $state(localStorage.getItem("MS_RANDOMIZE_MAP") === "true");
    $effect(() => {
        localStorage.setItem("MS_RANDOMIZE_MAP", randomizeMap.toString());
    });

    const existingMatchBehaviors: { [n: string]: number } = {
        "Restart if different": 0,
        Restart: 1,
        "Continue and spawn": 2,
    };

    function cleanCase(toClean: string) {
        toClean = toClean.replaceAll("_", " ").replace(" option", "")
        return toClean.charAt(0).toUpperCase() + toClean.slice(1);
    }

    function resetMutators() {
        for (let key of Object.keys(mutators)) {
            mutators[key] = 0;
        }
    }

    function setPreset(preset: string) {
        if (preset === "") {
            mode = mutatorOptions.game_mode[0];
            map = MAPS_STANDARD["DFH Stadium"];
            randomizeMap = true;
            resetMutators();
            return;
        }

        const presetData = gamemodes[preset];
        if (!presetData) return;

        if (presetData.match["game_mode"] !== undefined) {
            mode = presetData.match["game_mode"];
        }

        if (presetData.match["game_map_upk"] !== undefined) {
            map = presetData.match["game_map_upk"];
            randomizeMap = false;
        } else {
            randomizeMap = true;
        }

        for (let key of filteredMutatorOptions) {
            if (presetData.mutators[key] !== undefined) {
                mutators[key] = mutatorOptions[key].indexOf(presetData.mutators[key]);
            } else {
                mutators[key] = 0;
            }
        }
    }

    const ALL_MAPS = {...MAPS_STANDARD, ...MAPS_NON_STANDARD};
    const filteredMutatorOptions = Object.keys(mutatorOptions).filter((key) => key !== 'game_mode');
</script>

<div class="matchSettings">
    <h1>Match Settings</h1>
    <div class="content">
        <div class="settings">
            <div class="left-controls">
                <Select
                    options={ALL_MAPS}
                    bind:value={map}
                    placeholder="Select map"
                />
                <Select
                    options={Object.fromEntries(mutatorOptions.game_mode.map((x) => [x, x]))}
                    bind:value={mode}
                    placeholder="Select mode"
                />
            </div>
            <div class="right-controls">
                <LauncherSelector />
            </div>
        </div>
        <div class="controls">
            <div class="left-controls">
                <button
                    onclick={() => {
                        showMutators = true;
                    }}>Mutators</button
                >
                <button
                    onclick={() => {
                        showExtraOptions = true;
                    }}>Extra</button
                >
                <input
                    type="checkbox"
                    id="randomizeMap"
                    bind:checked={randomizeMap}
                />
                <label for="randomizeMap">Randomize Map</label>
            </div>
            <div class="right-controls">
                <button class="start" onclick={()=>{onStart(randomizeMap)}}>Start Match</button>
                <button class="stop" onclick={()=>{onStop()}}>Stop</button>
            </div>
        </div>
    </div>
</div>

<Modal title="Rocket League Mutators" bind:visible={showMutators}>
    <div class="mutators">
        {#each filteredMutatorOptions as mutatorKey}
            <div class="mutator">
                <label
                    style={mutators[mutatorKey] == 0 ? "color:lightgrey" : ""}
                    for={mutatorKey}>{cleanCase(mutatorKey)}</label
                >

                <select
                    name={mutatorKey}
                    id={mutatorKey}
                    bind:value={mutators[mutatorKey]}
                >
                    {#each mutatorOptions[mutatorKey] as value, i}
                        <option value={i}>{value.replaceAll("_", " ")}</option>
                    {/each}
                </select>
            </div>
        {/each}
        <div></div>
        <div class="mutator">
            <label for="gamemode"><b>Presets</b></label>
            <select name="gamemode" id="gamemode" onchange={(e) => {
                // @ts-ignore
                setPreset(e.target.value)
            }}>
                {#each ["", ...Object.keys(gamemodes)] as gamemode}
                    <option value={gamemode}>{gamemode}</option>
                {/each}
            </select>
        </div>
    </div>
    <div class="bottomButtons">
        <p>Settings are saved automatically</p>
        <button
            class="mutatorResetButton"
            onclick={resetMutators}>Reset</button
        >
    </div>
</Modal>

<Modal title="RLBot Extra Options" bind:visible={showExtraOptions}>
    <div class="extraoptions">
        <input
            type="checkbox"
            id="enableRendering"
            bind:checked={extraOptions.enableRendering}
        />
        <label for="enableRendering">
            Enable Rendering (bots can draw on screen)
        </label>
        <br />
        <input
            type="checkbox"
            id="enableStateSetting"
            bind:checked={extraOptions.enableStateSetting}
        />
        <label for="enableStateSetting">
            Enable State Setting (bots can teleport)
        </label>
        <br />
        <input
            type="checkbox"
            id="autoSaveReplay"
            bind:checked={extraOptions.autoSaveReplay}
        />
        <label for="autoSaveReplay"> Auto Save Replay </label>
        <br />
        <input
            type="checkbox"
            id="skipReplays"
            bind:checked={extraOptions.skipReplays}
        />
        <label for="skipReplays"> Skip Replays </label>
        <br />
        <input
            type="checkbox"
            id="instantStart"
            bind:checked={extraOptions.instantStart}
        />
        <label for="instantStart"> Instant Start </label>
        <br />
        <select
            name="emb"
            id="emb"
            bind:value={extraOptions.existingMatchBehavior}
        >
            {#each Object.keys(existingMatchBehaviors) as key}
                <option value={existingMatchBehaviors[key]}>{key}</option>
            {/each}
        </select>
        <label for="emb"></label>
    </div>
    <div class="bottomButtons">
        <p>Settings are saved automatically</p>
    </div>
</Modal>

<style>
    h1 {
        margin-bottom: 0.6rem;
    }
    .settings,
    .controls {
        display: flex;
        justify-content: space-between;
        gap: 0.5rem;
    }
    #randomizeMap {
        transform: scale(1.2);
    }
    .left-controls {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }
    .right-controls {
        display: flex;
        gap: 0.5rem;
    }
    .content {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .mutators {
        display: grid;
        grid-template-columns: auto auto auto;
        gap: 1rem;
    }
    .mutator {
        display: flex;
        flex-direction: column;
        gap: 0.3rem;
    }
    .bottomButtons {
        display: flex;
        margin-top: 1rem;
        justify-content: space-between;
        align-items: center;
    }
    .mutatorResetButton {
        background-color: red;
    }

    .extraoptions > * {
        margin-bottom: 0.5rem;
    }

    button.start {
        background-color: #15680e;
    }
    button.stop {
        background-color: #cc1414;
    }
</style>
