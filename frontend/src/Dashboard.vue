<template>
    <TotalTimeDisplay :totalTime="totalTime" />
    <DisplaySwitch @displayUpdated="handleDisplayUpdate" />
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


const hasData = ref(false);
const groupedData = ref<GroupedData>({});
const selectedKeyHash = ref<string>('');
const currentlyShown = ref<GroupBy>(GroupBy.Languages);

onMounted(async () => {
    const response = await fetch('/api/sessions')
    const json = await response.json() as GroupedData
    console.log(json)
    const keyHashes = Object.keys(json);
    if (keyHashes.length > 0) {
        groupedData.value = json;
        selectedKeyHash.value = keyHashes[0];
        hasData.value = true;
    }
})

const data = computed(() => groupedData.value[selectedKeyHash.value])

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
</script>
