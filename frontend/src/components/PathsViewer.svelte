<script lang="ts">
import toast from "svelte-5-french-toast";
import { App } from "../../bindings/gui";
import closeIcon from "../assets/close.svg";
import repairIcon from "../assets/repair.svg";
import Modal from "./Modal.svelte";
import Switch from "./Switch.svelte";
import ProgressBar from "./ProgressBar.svelte";
import { Events } from "@wailsio/runtime";

const OFFICIAL_BOTPACK_REPOS = [
  "VirxEC/botpack-test",
  "VirxEC/pytorch-archive",
];

let {
  visible = $bindable(false),
  paths = $bindable([]),
}: {
  visible?: boolean;
  paths?: {
    tagName: string | null;
    repo: string | null;
    installPath: string;
    visible: boolean;
    isDependency: boolean;
  }[];
} = $props();

let addBotpackVisible = $state(false);
let selectedBotpackType = $state("official");
let customRepo = $state("");
let installPath = $state("");

async function setDefaultPath() {
  const defaultPath = await App.GetDefaultPath();
  installPath = `${defaultPath}/RLBotPack`;
}

setDefaultPath();

function removePath(index: number) {
  const [removedItem] = paths.splice(index, 1);

  const dependency = paths.findIndex(
    (item) => item.installPath === removedItem.installPath,
  );
  if (dependency !== -1) removePath(dependency);
}

let downloadModalTitle = $state("Downloading & extracting repo/owner");
let downloadModalVisible = $state(false);
let downloadProgress = $state(0);
let downloadCurrentStep = $state(0);
let downloadTotalSteps = $state(0);

Events.On("monitor:download-progress", (event) => {
  const { status, done } = event.data.at(-1);

  if (done) {
    downloadProgress = 0;
    downloadCurrentStep += 1;
  } else {
    downloadProgress = status;
  }
});

async function repairBotpack(index: number) {
  const info = paths[index];
  if (!info.repo) {
    return; // can't update something that doesn't have a repo
  }

  downloadProgress = 0;
  downloadCurrentStep = 0;
  downloadTotalSteps = 0;

  for (const item of paths) {
    if (!item.repo || item.installPath !== info.installPath) continue;
    downloadTotalSteps += 2;
  }

  let clearInstallPath = true;
  for (const item of paths) {
    if (!item.repo || item.installPath !== info.installPath) continue;

    downloadModalTitle = `Redownloading & extracting ${item.repo}`;
    visible = false;
    downloadModalVisible = true;

    let tagName = await App.RepairBotpack(
      item.repo,
      item.installPath,
      clearInstallPath,
    ).catch((err) => {
      toast.error(`Failed to download botpack: ${err}`, {
        duration: 10000,
      });

      return null;
    });
    if (!tagName) break;

    clearInstallPath = false;
    downloadProgress = 0;
    downloadCurrentStep += 1;
  }

  downloadModalVisible = false;
  visible = true;
}

function openAddBotpackModal() {
  visible = false;
  addBotpackVisible = true;
}

function closeAddBotpackModal() {
  downloadModalVisible = false;
  addBotpackVisible = false;
  visible = true;
  selectedBotpackType = "official";
  customRepo = "";
  setDefaultPath();
}

function addInstallPath(installPath: string) {
  if (paths.some((x) => x.installPath === installPath)) {
    toast.error("Path already added");
    return;
  }

  paths.push({
    installPath,
    visible: true,
    tagName: null,
    repo: null,
    isDependency: false,
  });
}

async function addFolder() {
  const result = await App.PickFolder();

  if (!result) {
    return;
  }

  addInstallPath(result);
}

function addFile() {
  App.PickRLBotToml()
    .then((result) => {
      if (!result) {
        return;
      }

      addInstallPath(result);
    })
    .catch((_) => {
      toast.error(
        "Failed to add file: Only bot.toml or script.toml files are allowed",
      );
    });
}

async function confirmAddBotpack() {
  if (!installPath) {
    toast.error("Install path cannot be blank");
    return;
  }

  let repo: string, dep: string | null;

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
    dep = null;
  } else {
    [repo, dep] = OFFICIAL_BOTPACK_REPOS;
  }

  if (paths.some((x) => x.installPath === installPath)) {
    toast.error(`Install path "${installPath}" already in use for ${repo}`);
    return;
  }

  if (paths.some((x) => x.repo === repo)) {
    toast.error(`Botpack ${repo} already added`);
    return;
  }

  downloadModalTitle = `Downloading & extracting ${repo}`;
  downloadProgress = 0;
  addBotpackVisible = false;
  downloadModalVisible = true;
  downloadCurrentStep = 0;
  downloadTotalSteps = dep ? 4 : 2;

  const tagName = await App.DownloadBotpack(repo, installPath).catch((err) => {
    toast.error(`Failed to download botpack: ${err}`, {
      duration: 10000,
    });

    return null;
  });
  if (!tagName) {
    downloadModalVisible = false;
    addBotpackVisible = true;
    return;
  }

  if (dep) {
    downloadModalTitle = `Downloading & extracting ${dep}`;
    downloadProgress = 0;
    downloadCurrentStep += 1;

    // Download possible dependency of the given botpack
    const tagName = await App.DownloadBotpack(dep, installPath).catch((err) => {
      toast.error(`Failed to download botpack dependency: ${err}`, {
        duration: 10000,
      });

      return null;
    });
    if (!tagName) {
      downloadModalVisible = false;
      addBotpackVisible = true;
      return;
    }
  }

  paths.push({
    installPath,
    repo,
    tagName,
    visible: true,
    isDependency: false,
  });
  if (dep) {
    paths.push({
      installPath,
      repo: dep,
      tagName,
      visible: false,
      isDependency: true,
    });
  }

  toast.success("Botpack downloaded successfully!");
  closeAddBotpackModal();
}
</script>

<Modal title="Manage Paths" bind:visible={visible}>
  <div class="paths">
    <div class="button-row">
      <button onclick={addFolder}>Add folder</button>
      <button onclick={addFile}>Add File</button>
      <button onclick={openAddBotpackModal}>Add Botpack</button>
    </div>

    {#each paths as path, i}
      {#if path.isDependency}
        <div class="path" style="margin-left: 50px;">
          <pre>{path.repo}</pre>
        </div>
      {:else}
        <div class="path">
          <pre>{path.repo ? `${path.repo} @ ${path.installPath}` : path.installPath}</pre>
          {#if path.repo}
            <button class="no-left-margin repair" onclick={() => repairBotpack(i)}>
              <img src={repairIcon} alt="repair" />
            </button>
            <div><Switch bind:checked={path.visible} /></div>
          {:else}
            <div class="no-left-margin"><Switch bind:checked={path.visible} /></div>
          {/if}
          <button class="close" onclick={() => removePath(i)}>
            <img src={closeIcon} alt="X" />
          </button>
        </div>
      {/if}
    {/each}
  </div>
</Modal>

<Modal title="Add Botpack" bind:visible={addBotpackVisible}>
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

<Modal title={downloadModalTitle} bind:visible={downloadModalVisible} closeable={false}>
  <ProgressBar percentComplete={downloadProgress} currentStep={downloadCurrentStep} totalSteps={downloadTotalSteps} />
</Modal>

<style>
  .paths {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    min-width: 70vw;
    min-height: 50vh;
  }
  .button-row {
    display: flex;
    gap: 1rem;
  }
  .button-row {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
  }
  .button-row button {
    flex: 1;
  }
  .path {
    display: flex;
    align-items: center;
    gap: 1rem;
    justify-content: space-between;
  }
  .path .no-left-margin {
    margin-left: auto;
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
    min-width: 50vw;
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
  .repair {
    color: var(--foreground);
  }
</style>
