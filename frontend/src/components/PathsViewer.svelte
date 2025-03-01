<script lang="ts">
import toast from "svelte-5-french-toast";
import { App } from "../../bindings/gui";
import closeIcon from "../assets/close.svg";
import Modal from "./Modal.svelte";

const OFFICIAL_BOTPACK_REPO = "VirxEC/botpack-test";

let {
  visible = $bindable(false),
  paths = $bindable([]),
}: {
  visible?: boolean;
  paths?: {
    tagName: string | null;
    repo: string | null;
    installPath: string;
  }[];
} = $props();

let addBotpackVisible = $state(false);
let selectedBotpackType = $state("official");
let customRepo = $state("");
let installPath = $state("");

function setDefaultPath() {
  App.GetDefaultPath().then((result) => {
    installPath = `${result}/RLBotPack`;
  });
}

setDefaultPath();

function removePath(index: number) {
  paths.splice(index, 1);
}

function openAddBotpackModal() {
  visible = false;
  addBotpackVisible = true;
}

function closeAddBotpackModal() {
  addBotpackVisible = false;
  visible = true;
  selectedBotpackType = "official";
  customRepo = "";
  setDefaultPath();
}

function addFolder() {
  App.PickFolder().then((result) => {
    if (result) {
      if (paths.some((x) => x.installPath === installPath)) {
        toast.error("Path already added");
        return;
      }

      paths.push({ installPath: result, repo: null, tagName: null });
    }
  });
}

function confirmAddBotpack() {
  if (!installPath) {
    toast.error("Install path cannot be blank");
    return;
  }

  let repo = OFFICIAL_BOTPACK_REPO;
  if (selectedBotpackType === "custom") {
    if (!customRepo) {
      toast.error("URL cannot be blank");
      return;
    }

    if (!/^[\w-]+\/[\w-]+$/.test(customRepo)) {
      toast.error("Custom repository must be in the format 'owner/repo'");
      return;
    }

    repo = customRepo;
  }

  if (paths.some((x) => x.repo === repo)) {
    toast.error("Botpack already added");
    return;
  }

  if (paths.some((x) => x.installPath === installPath)) {
    toast.error("Install path already in use");
    return;
  }

  const id = toast.loading("Downloading botpack...");
  App.DownloadBotpack(repo, installPath)
    .then((tagName) => {
      toast.success("Botpack downloaded successfully!", { id });

      paths.push({ installPath, repo, tagName });
      closeAddBotpackModal();
    })
    .catch((err) => {
      toast.error(`Failed to download botpack: ${err}`, {
        duration: 10000,
        id,
      });
    });
}
</script>

<Modal title="Manage Paths" bind:visible={visible} minWidth="70vw" minHeight="50vh">
    <div class="paths">
        <div class="button-row">
            <!-- TODO: this class is not needed -->
            <button class="full-width" onclick={addFolder}>Add folder</button>
            <button class="full-width" onclick={openAddBotpackModal}>Add Botpack</button>
            <button class="full-width" onclick={alert.bind(null, "TODO: not implemented yet")}>Add File</button>
        </div>

        {#each paths as path, i}
            <div class="path">
                <pre>{path.repo ? `${path.repo} @ ${path.installPath}` : path.installPath}</pre>
                <button class="close" onclick={() => removePath(i)}>
                    <img src={closeIcon} alt="X" />
                </button>
            </div>
        {/each}
    </div>
</Modal>

<Modal title="Add Botpack" bind:visible={addBotpackVisible} minWidth="50vw">
    <div class="add-botpack">
        <label for="path">Botpack install path</label>
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
            <input type="text" placeholder="owner/repo" bind:value={customRepo} />
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
