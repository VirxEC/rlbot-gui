<script lang="ts">
    import { App } from "../../bindings/gui";
    import Modal from "./Modal.svelte";
    import closeIcon from "../assets/close.svg";
    import toast from "svelte-5-french-toast";

    const GITHUB_URL = "https://github.com/";
    const OFFICIAL_BOTPACK_URL = GITHUB_URL + "VirxEC/botpack-test";

    let { visible = $bindable(false), paths = $bindable([]) } = $props();

    let botpacks: { url: string, installPath: string }[] = $state(
        JSON.parse(localStorage.getItem("BOTPACKS") || "[]"),
    );

    $effect(() => {
        localStorage.setItem("BOTPACKS", JSON.stringify(botpacks));
    });

    let addBotpackVisible = $state(false);
    let selectedBotpackType = $state("official");
    let customURL = $state("");
    let installPath = $state("");

    function setDefaultPath() {
        App.GetDefaultPath().then((result) => {
            installPath = result + "/RLBotPack";
        });
    }

    setDefaultPath();

    function removePath(index: number) {
        paths.splice(index, 1);
    }

    function removeBotpack(index: number) {
        botpacks.splice(index, 1);
    }

    function openAddBotpackModal() {
        visible = false;
        addBotpackVisible = true;
    }

    function closeAddBotpackModal() {
        addBotpackVisible = false;
        visible = true;
        selectedBotpackType = "official";
        customURL = "";
        setDefaultPath();
    }

    function addFolder() {
        App.PickFolder().then((result) => {
            if (result) {
                paths = [...paths, result];
            }
        });
    }

    function confirmAddBotpack() {
        if (!installPath) {
            toast.error("Install path cannot be blank");
            return;
        }
    
        let url = OFFICIAL_BOTPACK_URL;
        if (selectedBotpackType === "custom") {
            if (!customURL) {
                toast.error("URL cannot be blank");
                return;
            }

            if (!customURL.startsWith("https://github.com/")) {
                toast.error("Custom URL must start with 'https://github.com/'");
                return;
            }

            url = customURL;
        }

        if (botpacks.some((x) => x.url === url)) {
            toast.error("Botpack already added");
            return;
        }

        if (botpacks.some((x) => x.installPath === installPath)) {
            toast.error("Install path already in use");
            return;
        }

        App.DownloadBotpack(url, installPath)
            .then((result) => {
                if (!result.success) {
                    toast.error("Failed to download botpack: " + result.message);
                    return;
                }
                    
                botpacks = [...botpacks, { installPath, url }];
                closeAddBotpackModal();
            })
    }
</script>

<Modal title="Manage Paths" bind:visible={visible} minWidth="70vw" minHeight="50vh">
    <div class="paths">
        <div class="button-row">
            <button class="full-width" onclick={addFolder}>Add folder</button>
            <button class="full-width" onclick={openAddBotpackModal}>Add Botpack</button>
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

        {#each botpacks as botpack, i}
            <div class="path">
                <pre>{botpack.url} @ {botpack.installPath}</pre>
                <button class="close" onclick={() => removeBotpack(i)}>
                    <img src={closeIcon} alt="X" />
                </button>
            </div>
        {/each}
    </div>
</Modal>

<Modal title="Add Botpack" bind:visible={addBotpackVisible} minWidth="50vw">
    <div class="add-botpack">
        <label for="path">Install path</label>
        <input type="text" id="path" placeholder="Enter install path" bind:value={installPath} />
        <div class="radio-group">
            <label>
                <input type="radio" name="botpackType" value="official" bind:group={selectedBotpackType} />
                Official RLBotPack
            </label>
            <label>
                <input type="radio" name="botpackType" value="custom" bind:group={selectedBotpackType} />
                Custom
            </label>
        </div>
        {#if selectedBotpackType === "custom"}
            <input type="text" placeholder="https://github.com/owner/repo" bind:value={customURL} />
        {/if}
        <div class="button-row">
            <button onclick={confirmAddBotpack}>Confirm</button>
            <button onclick={closeAddBotpackModal}>Cancel</button>
        </div>
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
    .add-botpack {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }
    .radio-group {
        flex-direction: column;
        display: flex;
        gap: 1rem;
    }
    .radio-group label {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }
    .button-row {
        display: flex;
        gap: 1rem;
        justify-content: flex-end;
    }
</style>
