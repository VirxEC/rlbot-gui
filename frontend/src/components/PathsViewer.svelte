<script lang="ts">
    import { App } from "../../bindings/gui";
    import Modal from "./Modal.svelte";
    import closeIcon from "../assets/close.svg";

    let { visible = $bindable(false), paths = $bindable([]) } = $props();

    function removePath(index: number) {
        paths.splice(index, 1);
    }
</script>

<Modal title="Manage Paths" bind:visible={visible} minWidth="70vw" minHeight="50vh">
    <div class="paths">
        <div class="button-row">
            <button
                class="full-width"
                onclick={async () => {
                    let result = await App.PickFolder();
                    console.log("PickFolder returned:", result);
                    if (result != "") {
                        paths = [...paths, result];
                    }
                }}>Add folder</button>
            <button class="full-width" onclick={alert.bind(null, "TODO: not implemented yet")}>Add Botpack</button>
            <button class="full-width" onclick={alert.bind(null, "TODO: not implemented yet")}>Add File</button>
        </div>

        {#each paths as path, i}
            <div class="path">
                <pre>{path}</pre>
                <button class="close" onclick={() => removePath(i)}>
                    <img src={closeIcon} alt="X" />
                </button>
            </div>
        {/each}
    </div>
</Modal>

<style>
    .paths {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }
    .button-row {
        display: flex;
        gap: 1rem;
    }
    .full-width {
        flex: 1;
    }
    .path {
        display: flex;
        align-items: center;
        gap: 1rem;
        justify-content: space-between;
    }
    .path pre {
        font-size: 1rem;
        margin: 0px;
    }
    .path button {
        padding: 0px;
    }
    .path button img {
        filter: invert();
    }
</style>
