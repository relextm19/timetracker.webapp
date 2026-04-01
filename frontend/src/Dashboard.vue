<template>
    <TotalTimeDisplay :totalTime="totalTime" />
    <div class="flex justify-evenly items-center">
        <!-- FIXME: looks kinda shitty-->
        <DisplaySwitch @displayUpdated="handleDisplayUpdate" :options="Object.values(GroupBy)"
            :selected="GroupBy.Languages" />
        <DisplaySwitch v-if="hasData" @displayUpdated="handleAPIKeyFilterChange" :options="APIKeys!.map(k => k.name)"
            :selected="APIKeys![0].name" class="ml-auto" />
    </div>
    <div v-if="hasData">
        <div v-if="displayComponent == MyCalendar">
            <MyCalendar :timeData="data!.byTime" />
        </div>
        <div v-else v-for="entry in currentData" :key="entry.name">
            <component :is="displayComponent" v-bind="getProps(entry)" />
        </div>
    </div>
    <div v-else class="w-full text-center ">
        <span class="text-2xl">
            No data yet
        </span>
    </div>
</template>

<script lang="ts">
export interface timeData {
    name: string
    totalTime: number
}

interface Data {
    byLanguage: timeData[]
    byProject: timeData[]
    byFile: timeData[]
    byTime: timeData[]
}

type GroupedData = Record<string, Data>

export enum GroupBy {
    Languages = "languages",
    Projects = "projects",
    Files = "files",
    TimeAggregated = "calendar"
}
</script>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import TotalTimeDisplay from './components/TotalTimeDisplay.vue'
import LanguageTimeDisplay from './components/LanguageTimeDisplay.vue'
import ProjectTimeDisplay from './components/ProjectTimeDisplay.vue'
import DisplaySwitch from './components/DisplaySwitch.vue'
import MyCalendar from './components/MyCalendar.vue'
import type { APIKey } from './ApiKeys.vue'

const hasData = ref(false);
const groupedData = ref<GroupedData>({});
const selectedKeyHash = ref<APIKey>();
const APIKeys = ref<APIKey[]>();
const currentlyShown = ref<GroupBy>(GroupBy.Languages);

const fetchSessions = async (): Promise<GroupedData> => {
    const response = await fetch('/api/sessions')
    const json = await response.json() as GroupedData
    return json
}

const fetchAPIKeys = async (): Promise<APIKey[]> => {
    const res = await fetch('/api/keys');
    const json = await res.json() as APIKey[];
    return json
}

onMounted(async () => {
    const sessions = await fetchSessions();
    const keys = await fetchAPIKeys();
    if (keys.length > 0 && Object.values(sessions).length > 0) {
        groupedData.value = sessions;
        APIKeys.value = keys;
        selectedKeyHash.value = keys[0];
        hasData.value = true;
    }
})

const data = computed(() => {
    if (!selectedKeyHash.value) return undefined;
    return groupedData.value[selectedKeyHash.value.keyHash];
})

const currentData = computed(() => {
    switch (currentlyShown.value) {
        case GroupBy.Languages:
            return data.value?.byLanguage;
        case GroupBy.Projects:
            return data.value?.byProject;
        case GroupBy.Files:
            return data.value?.byFile;
        case GroupBy.TimeAggregated:
            return data.value?.byTime;
        default:
            return [];
    }
});
const displayComponent = computed(() => {
    switch (currentlyShown.value) {
        case GroupBy.Languages:
            return LanguageTimeDisplay;
        case GroupBy.Projects:
            return ProjectTimeDisplay;
        case GroupBy.Files:
            return ProjectTimeDisplay;
        case GroupBy.TimeAggregated:
            return MyCalendar;
        default:
            return [];
    }
});
const totalTime = computed(() => data.value?.byLanguage?.reduce((acc, e) => acc + e.totalTime, 0) ?? 0)

const getProps = (data: timeData) => {
    return { name: data.name, time: data.totalTime }
}

const handleDisplayUpdate = (val: GroupBy): void => {
    currentlyShown.value = val;
}

// TODO: change the val to APIKEY but in order to do ts i need to introduce a generic in the display switch
const handleAPIKeyFilterChange = (val: string): void => {
    selectedKeyHash.value = APIKeys.value!.filter(k => k.name === val)[0];
}
</script>
