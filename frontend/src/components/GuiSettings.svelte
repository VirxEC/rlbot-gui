<script lang="ts">
import Modal from "./Modal.svelte";

// TODO: Save settings, svelte store + localstorage?
//       Perhaps change all localstorage state to svelte stores?
let { visible = $bindable(false) } = $props();
</script>

<Modal bind:visible title="GUI Settings">
  <div class="inner">
    <h3>Not implemented yet</h3>
    <!-- TODO: Match start timeouts -->
    <!-- TODO: Dark/Light themes -->
    <!-- TODO: Refresh bots behavior (remove on refresh, remove not found agents, etc.) -->
    <!-- TODO: Telemetry settings if added -->
    <!-- TODO: Auto update botpack -->
  </div>
</Modal>

<style>
  .inner {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
  }
  h3 {
    padding: 5rem;
  }
</style>
