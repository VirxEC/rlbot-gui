<script lang="ts">
import { Browser } from "@wailsio/runtime";
import { onMount } from "svelte";
import AlarmIcon from "../assets/alarm.svg";
import CalendarPlusIcon from "../assets/calendar-plus.svg";
import GeoIcon from "../assets/geo.svg";
import InfoIcon from "../assets/info_icon.svg";
import Modal from "./Modal.svelte";

let {
  visible = $bindable(false),
  eventsNow = $bindable(0),
  eventsFuture = $bindable(0),
} = $props();

let now = $state(new Date());

onMount(() => {
  const start = new Date();
  const timeToFirstUpdate = 60000 - (start.getTime() % 60000) + 1000;

  // schedule the interval to run every minute,
  // one second after the minute starts
  // this ensures that we aren't effected by inaccuracies
  setTimeout(() => {
    now = new Date();

    setInterval(() => {
      now = new Date();
    }, 60000); // Update every minute
  }, timeToFirstUpdate);
});

function dateTimeCheck(today: Date, event: any) {
  const name = event.summary;
  const start = event.start.dateTime;
  const newDate = new Date(start);

  if (event.recurrence) {
    try {
      const recurrence = event.recurrence[0].split(";");
      const recType = recurrence[0].split("=")[1];
      const interval = recurrence[2].split("=")[1];
      const endDateType = recurrence[2].split("=")[0];
      const endDateRaw = recurrence[2].split("=")[1];

      let endDate = new Date(newDate);
      if (endDateType === "COUNT") {
        if (recType === "WEEKLY") {
          endDate.setDate(newDate.getDate() + 7 * interval * endDateRaw);
        } else if (recType === "MONTHLY") {
          endDate.setDate(newDate.getDate() + 4 * interval * endDateRaw);
        }
      } else {
        endDate.setDate(endDateRaw);
      }

      if (recType === "WEEKLY") {
        while (newDate <= endDate) {
          if (newDate > today) {
            break;
          }

          newDate.setDate(newDate.getDate() + 7);
        }
      } else if (recType === "MONTHLY") {
        while (newDate <= endDate) {
          if (newDate > today) {
            break;
          }

          newDate.setDate(newDate.getDate() + 4);
        }
      }
    } catch (e) {
      console.error(`Error checking recurrence: ${e}`);
    }
  }

  return [name, newDate];
}

function formatFromNow(milliseconds: number) {
  const minuteMillis = 1000 * 60;
  const hourMillis = minuteMillis * 60;
  const dayMillis = hourMillis * 24;

  let format = "";

  const days = Math.floor(milliseconds / dayMillis);
  if (days > 0) {
    format += days;
    format += days > 1 ? " days " : " day ";
  }

  const hours = Math.floor((milliseconds % dayMillis) / hourMillis);
  if (hours > 0) {
    format += hours;
    format += hours > 1 ? " hours " : " hour ";
  }

  const minutes = Math.ceil((milliseconds % hourMillis) / minuteMillis);
  if (minutes > 0) {
    format += minutes;
    format += minutes > 1 ? " minutes " : " minute ";
  }

  return format;
}

async function fetchEvents() {
  const today = new Date();

  const apiKey = "AIzaSyBQ40UqlMPexzWxTNd7EYtTrkoFF_DqpqM";
  const timeMin = today.toISOString();
  const url = `https://www.googleapis.com/calendar/v3/calendars/rlbotofficial@gmail.com/events?maxResults=10&timeMin=${timeMin}&key=${apiKey}`;

  const response = await fetch(url);
  const data = await response.json();

  const events: {
    name: string;
    location: string;
    time: string;
    date: Date;
    moreInfo: string;
    logo: string;
  }[] = [];

  // compute dates and times
  for (const event of data.items) {
    const [name, newDate] = dateTimeCheck(today, event);

    const remainingTimeInMs = newDate.getTime() - today.getTime();
    if (remainingTimeInMs > 0) eventsFuture += 1;
    else eventsNow += 1;

    const logoSplit = event.description ? event.description.split("logo:") : [];
    const logo =
      logoSplit.length > 1
        ? logoSplit[1].replace("\n", "").split('href="')[1].split('"')[0]
        : null;

    const moreInfo =
      logoSplit.length > 0 && logoSplit[0].includes('href="')
        ? logoSplit[0].replace("\n", "").split('href="')[1].split('"')[0]
        : event.description;

    events.push({
      name,
      location: event.location,
      time: newDate.toLocaleString(),
      date: newDate,
      moreInfo,
      logo,
    });
  }

  // sort community events by start time
  events.sort((a, b) => {
    // @ts-ignore
    return a.date - b.date;
  });

  return events;
}
</script>

<Modal title="Community Events" bind:visible>
  <!-- svelte-ignore a11y_invalid_attribute -->
  <div id="community-events">
    {#await fetchEvents()}
      <p>Loading events...</p>
    {:then events }
      {#if events.length === 0}
        <p>There are no community events at this time.</p>
      {:else}
      {#each events as event}
        {@const remainingTimeInMs = event.date.getTime() - now.getTime()}
        {@const timeUntil = formatFromNow(Math.abs(remainingTimeInMs))}
        <div class="event">
          {#if event.logo}
            <img class="event-logo" src={event.logo} alt="event logo" />
          {/if}
          <div>
            <h2>{ event.name }</h2>
            {#if !timeUntil}
            <p>
              <img src={AlarmIcon} alt="alarm" /> Starting now!
            </p>
            {:else if remainingTimeInMs > 0}
            <p>
              <img src={CalendarPlusIcon} alt="calendar" /> Starts in <b>{ timeUntil }</b> ({ event.time })
            </p>
            {:else}
            <p>
              <img src={AlarmIcon} alt="alarm" /> Started <b>{ timeUntil }</b> ago, but you can still join!
            </p>
            {/if}
            <p>
              <img src={GeoIcon} alt="location"> <a href="#" onclick={() => {Browser.OpenURL(event.location)}} target="_blank">{ event.location }</a>
            </p>
            {#if event.moreInfo}
            <br>
            <p class="more-info">
              <img src={InfoIcon} alt="info"> <a href="#" onclick={() => {Browser.OpenURL(event.moreInfo)}} target="_blank">More info</a>
            </p>
            {/if}
          </div>
        </div>
      {/each}
      {/if}
    {/await}
  </div>
</Modal>

<style>
  h2 {
    font-size: 1.6em;
    margin-bottom: 0.25em;
  }
  p {
    align-items: center;
    vertical-align: middle;
    margin-top: 0.5em;
  }
  p img {
    filter: invert() brightness(90%);
    width: 24px;
    margin-left: 0.5em;
    margin-bottom: 3px;
    vertical-align: middle;
  }
  .event {
    display: flex;
    align-items: center;
    vertical-align: middle;
    margin-bottom: 1.2em;
  }
  .event > img {
    margin-left: 0.5em;
    margin-right: 1.2em;
  }
  .event > div {
    display: flex;
    flex-direction: column;
    justify-content: center;
  }
  .event-logo {
    max-height: auto;
    width: 18%;
  }
  #community-events {
    min-width: 500px;
    min-height: 150px;
    width: 100%;
    height: 100%;
    overflow: auto;
  }
</style>
