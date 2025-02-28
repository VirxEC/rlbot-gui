<script lang="ts">
    import close from "../assets/close.svg";
    let {
        title = "Modal",
        visible = $bindable(true),
        children,
        minWidth = "20vw",
        minHeight = "20vh"
    } = $props();

    let background: EventTarget;
    let wrap: EventTarget;

    function handleOuter(e: MouseEvent) {
        if (e.target === background || e.target === wrap)
            visible = false;
    }
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class={visible ? "" : "hidden"} bind:this={background} onclick={handleOuter}>
    <div class="modalContainer" bind:this={wrap}>
        <div class="modal" style={`min-width: ${minWidth}; min-height: ${minHeight};`}>
            <header>
                <h2>{title}</h2>
                <button
                    onclick={() => {
                        visible = false;
                    }}><img src={close} alt="close" /></button
                >
            </header>
            <div class="modalBody">
                {@render children?.()}
            </div>
        </div>
    </div>
</div>

<style>
    .modalContainer {
        display: flex;
        height: 100vh;
        width: 100vw;
        position: fixed;
        top: 0px;
        left: 0px;
        justify-content: center;
        align-items: center;
        z-index: 100;
        background-color: #00000080;
        transition: opacity 0.2s;
    }

    .hidden * {
        opacity: 0;
        z-index: -100 !important;
    }

    .modal {
        background-color: var(--background);
        padding: 0.2rem;
        border-radius: 0.6rem;
    }
    header {
        padding: 0.2rem;
        padding-left: 0.8rem;
        border-bottom: 1px solid var(--background-alt);
        display: flex;
        justify-content: space-between;
        align-items: center;
    }
    header button {
        padding: 0px;
        background-color: transparent;
        filter: invert();
    }
    .modalBody {
        padding: 1rem;
        max-width: 80vw;
        max-height: 80vh;
    }
</style>