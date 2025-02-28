<script lang="ts">
    import Modal from "./Modal.svelte";

    let visible = $state(false);
    let localLauncher = $state(localStorage.getItem('MS_LAUNCHER') || '');
    let localLauncherArg = $state(localStorage.getItem('MS_LAUNCHER_ARG') || '');

    function uid() {
        return Date.now().toString(36) + Math.random().toString(36).substring(2);
    }
    let uids = [uid(), uid(), uid(), uid(), uid()];

    $effect(() => {
        if (localLauncher === 'legendary') {
            localStorage.setItem('MS_LAUNCHER', 'custom');
            localStorage.setItem('MS_LAUNCHER_ARG', 'legendary');
        } else {
            localStorage.setItem('MS_LAUNCHER', localLauncher);
            localStorage.setItem('MS_LAUNCHER_ARG', localLauncherArg);
        }
    });
</script>

<button onclick={() => { visible = true }}>Launcher Options</button>

<Modal title="Select a launcher" bind:visible={visible}>
    <div class="container">
        <div class="launcherSelect">
            <input name="launcher" type="radio" id={`launcher-steam-${uids[0]}`} bind:group={localLauncher} value="steam">
            <label for={`launcher-steam-${uids[0]}`}>Steam</label><br>

            <input name="launcher" type="radio" id={`launcher-epic-${uids[1]}`} bind:group={localLauncher} value="epic">
            <label for={`launcher-epic-${uids[1]}`}>Epic</label><br>

            <input name="launcher" type="radio" id={`launcher-custom-${uids[2]}`} bind:group={localLauncher} value="custom">
            <label for={`launcher-custom-${uids[2]}`}>Custom</label><br>

            <input name="launcher" type="radio" id={`launcher-legendary-${uids[3]}`} bind:group={localLauncher} value="legendary">
            <label for={`launcher-legendary-${uids[3]}`}>Legendary</label><br>

            <input name="launcher" type="radio" id={`launcher-nolaunch-${uids[4]}`} bind:group={localLauncher} value="nolaunch">
            <label for={`launcher-nolaunch-${uids[4]}`}>Don't launch</label>
        </div>
        {#if localLauncher !== "legendary" && localLauncher !== "nolaunch"}
        <div class="launcherArg">
            <label for="launcherArg">Additional argument:</label>
            <input type="text" id="launcherArg" bind:value={localLauncherArg} placeholder="(Leave blank for default)">
        </div>
        {/if}
    </div>
</Modal>

<style>
.container {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}
.launcherSelect {
    font-size: 1.1rem;
}
input[type='radio'] {
    transform: scale(2);
}
.launcherArg {
    display: flex;
    flex-direction: column;
}
</style>
