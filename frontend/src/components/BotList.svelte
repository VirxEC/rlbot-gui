<script lang="ts">
    import { flip } from "svelte/animate";
    import {
        dndzone,
        TRIGGERS,
        SHADOW_ITEM_MARKER_PROPERTY_NAME,
    } from "svelte-dnd-action";
    import defaultIcon from "../assets/rlbot_mono.png";
    import type { DraggablePlayer } from "../index";
    let { items = [] }: { items: DraggablePlayer[] } = $props();
    const flipDurationMs = 100;

    let selectedTag = $state("All");
    const extraModeTags = ["hoops", "dropshot", "snow-day", "spike-rush", "heatseaker"];
    const categories = [
        ["All"],
        ["Standard", "Extra Modes", "Special bots/scripts"],
        ["Bots for 1v1", "Bots with teamplay", "Goalie bots"]
    ];

    let filteredItems: DraggablePlayer[] = $state([]);
    $effect(() => {
        filteredItems = filterBots(selectedTag);
    });

    function filterBots(filterTag: string) {
        return items.filter((bot) => {
            switch (filterTag) {
                case "Standard":
                    return !bot.tags.some((tag) =>
                        [...extraModeTags, "memebot", "human"].includes(tag),
                    );
                case "Extra Modes":
                    return bot.tags.some((tag) => extraModeTags.includes(tag));
                case "Special bots/scripts":
                    return bot.tags.some((tag) => tag === "memebot");
                case "Bots for 1v1":
                    return bot.tags.some((tag) => tag === "1v1");
                case "Bots with teamplay":
                    return bot.tags.some((tag) => tag === "teamplay");
                case "Goalie bots":
                    return bot.tags.some((tag) => tag === "goalie");
                default:
                    return items;
            }
        });
    }

    function handleTagClick(tag: string) {
        selectedTag = tag;
    }

    function handleDndConsider(e: any) {
        const { trigger, id } = e.detail.info;
        if (trigger === TRIGGERS.DRAG_STARTED) {
            const newId = `${id}-${Math.round(Math.random() * 100000)}`;
            const idx = filteredItems.findIndex((item) => item.id === id);
            console.log(e.detail.items.filter((item: any) => item[SHADOW_ITEM_MARKER_PROPERTY_NAME]).forEach((item: any) => item.id = newId));
            e.detail.items = e.detail.items.filter(
                (item: any) => !item[SHADOW_ITEM_MARKER_PROPERTY_NAME],
            );
            e.detail.items.splice(idx, 0, { ...filteredItems[idx], id: newId });
            filteredItems = e.detail.items;
        }
    }
    function handleDndFinalize(e: any) {
        filteredItems = e.detail.items;
    }
</script>

<div class="tag-buttons">
    {#each categories as tagGroup}
        <div class="tag-group">
            {#each tagGroup as tag}
                <button onclick={() => handleTagClick(tag)} class:selected={selectedTag === tag}>
                    {tag}
                </button>
            {/each}
        </div>
    {/each}
</div>

<div
    class="bots"
    use:dndzone={{
        items: filteredItems,
        flipDurationMs,
        centreDraggedOnCursor: true,
        dropTargetStyle: {},
        dropTargetClasses: ["dropTarget"],
    }}
    onconsider={handleDndConsider}
    onfinalize={handleDndFinalize}
>
    {#each filteredItems as bot (bot.id)}
        <div class="bot" animate:flip={{ duration: flipDurationMs }}>
            <img src={bot.icon || defaultIcon} alt="icon" />
            <p>{bot.displayName}</p>
        </div>
    {/each}
</div>

<style>
    .tag-buttons {
        display: flex;
        flex-direction: row;
        justify-content: center;

        gap: 0.5rem;
        margin-bottom: .6rem;
    }
    .tag-group {
        display: flex;
        gap: 2px;
        border: solid 1px gray;
        --border-radius: 0.25rem;
        border-radius: var(--border-radius);
    }
    .tag-buttons button {
        padding: 0.5rem 1rem;
        border-radius: 0;
        cursor: pointer;
        background-color: var(--background-alt);
    }
    .tag-buttons button:first-child {
        border-radius: var(--border-radius) 0 0 var(--border-radius);
    }
    .tag-buttons button:last-child {
        border-radius: 0 var(--border-radius) var(--border-radius) 0;
    }
    .tag-buttons button:only-child {
        border-radius: var(--border-radius);
    }
    .tag-buttons button.selected {
        background-color: var(--foreground);
        color: var(--background);
        /* font-size: 1.1rem; */
    }
    .bots {
        display: flex;
        gap: 0.5rem;
        flex-wrap: wrap;
    }
    .bot {
        display: flex;
        align-items: center;
        background-color: var(--background-alt);
        color: var(--foreground);
        height: 2.25rem;
        padding: 0.2rem;
        padding-right: 0.6rem;
        gap: 0.5rem;
        border-radius: 0.2rem;
    }
    .bot img {
        height: 2rem;
        width: auto;
    }
</style>
