import { writable } from 'svelte/store';
import { MAPS_STANDARD } from './arena-names';

export const mapStore = writable(localStorage.getItem("MS_MAP") || MAPS_STANDARD["DFH Stadium"]);
mapStore.subscribe(value => {
    localStorage.setItem("MS_MAP", value);
});
