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
    const tags = [
        ["All"],
        ["Standard", "Extra Modes", "Special bots/scripts"],
        ["Bots for 1v1", "Bots with teamplay", "Goalie bots"]
    ];

    function filterBots() {
        switch (selectedTag) {
            case "Standard":
                return items.filter((bot) => {
                    return !bot.tags.some((tag) =>
                        [...extraModeTags, "memebot", "human", "psyonix"].includes(tag),
                    );
                });
            case "Extra Modes":
                return items.filter((bot) => {
                    return bot.tags.some((tag) => extraModeTags.includes(tag));
                });
            case "Special bots/scripts":
                return items.filter((bot) => {
                    return bot.tags.some((tag) => tag === "memebot");
                });
            case "Bots for 1v1":
                return items.filter((bot) => {
                    return bot.tags.some((tag) => tag === "1v1");
                });
            case "Bots with teamplay":
                return items.filter((bot) => {
                    return bot.tags.some((tag) => tag === "teamplay");
                });
            case "Goalie bots":
                return items.filter((bot) => {
                    return bot.tags.some((tag) => tag === "goalie");
                });
            default:
                return items;
        }
    }

    function handleTagClick(tag: string) {
        selectedTag = tag;
    }

    let shouldIgnoreDndEvents = false;
    function handleDndConsider(e: any) {
        const { trigger, id } = e.detail.info;
        if (trigger === TRIGGERS.DRAG_STARTED) {
            const idx = items.findIndex((item) => item.id === id);
            const newId = `${id}_copy_${Math.round(Math.random() * 100000)}`;
            e.detail.items = e.detail.items.filter(
                (item: any) => !item[SHADOW_ITEM_MARKER_PROPERTY_NAME],
            );
            e.detail.items.splice(idx, 0, { ...items[idx], id: newId });
            items = e.detail.items;
            shouldIgnoreDndEvents = true;
        } else if (!shouldIgnoreDndEvents) {
            // items = e.detail.items;
        } else {
            items = [...items];
        }
    }
    function handleDndFinalize(e: any) {
        if (!shouldIgnoreDndEvents) {
            items = e.detail.items;
        } else {
            items = [...items];
            shouldIgnoreDndEvents = false;
        }
    }
</script>

<div class="tag-buttons">
    {#each tags as tagGroup}
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
        items,
        flipDurationMs,
        centreDraggedOnCursor: true,
        dropTargetStyle: {},
        dropTargetClasses: ["dropTarget"],
    }}
    onconsider={handleDndConsider}
    onfinalize={handleDndFinalize}
>
    {#each filterBots() as bot (bot.id)}
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
        gap: 0.5rem;
        margin-bottom: 1rem;
    }
    .tag-group {
        display: flex;
        gap: 0;
    }
    .tag-buttons button {
        padding: 0.5rem 1rem;
        border: 2px solid white;
        border-radius: 0;
        cursor: pointer;
        background-color: var(--background-alt);
    }
    .tag-buttons button.selected {
        background-color: white;
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
