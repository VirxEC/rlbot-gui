<script>
    let { onBack } = $props();
</script>

<div class="page">
    <button onclick={onBack}>Take be back to RLBot</button>
</div>

<style>
    .page {
        display: flex;
        align-items: center;
        justify-content: center;
    }
</style>
