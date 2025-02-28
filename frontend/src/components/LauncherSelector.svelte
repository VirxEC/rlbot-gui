<script lang="ts">
    import Modal from "./Modal.svelte"
    // todo: remember these! global state etc, save locally through go?
    let {launcher = $bindable(), gamePath = $bindable()} = $props();
    let visible = $state(false)

    function uid(){
      return Date.now().toString(36) + Math.random().toString(36).substring(2);
    }
    let uids = [uid(), uid(), uid()];
</script>

<button onclick={()=>{visible = true}}>Launcher Options</button>

<Modal title="Select a launcher" bind:visible>
    <div class="container">
        <div class="launcherSelect">
            <input name="launcher" type="radio" id={`launcher-steam-${uids[0]}`} bind:group={launcher} value="steam">
            <label for={`launcher-steam-${uids[0]}`}>Steam</label><br>

            <input name="launcher" type="radio" id={`launcher-epic-${uids[1]}`} bind:group={launcher} value="epic">
            <label for={`launcher-epic-${uids[1]}`}>Epic</label><br>

            <input name="launcher" type="radio" id={`launcher-custom-${uids[2]}`} bind:group={launcher} value="custom">
            <label for={`launcher-custom-${uids[2]}`}>Custom</label>
        </div>
        <div class="gamePath">
            <label for="gamePath">Game path:</label>
            <input type="text" id="gamePath" bind:value={gamePath} placeholder="(Leave blank for default)" >
        </div>
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
    transform: scale(1.1);
}
.gamePath {
    display: flex;
    flex-direction: column;
}
</style>
