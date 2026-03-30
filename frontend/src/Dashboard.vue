<template>
    <TotalTimeDisplay :totalTime="totalTime" />
    <DisplaySwitch @displayUpdated="handleDisplayUpdate" />
    <div v-if="hasData">
        <div v-if="displayComponent == MyCalendar">
            <displayComponent />
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
interface timeData {
    name: string
    totalTime: number
}

interface Data {
    byLanguage: timeData[]
    byProject: timeData[]
    byFile: timeData[]
    byTime: timeData[]
}

export enum GroupBy {
    Languages = "languages",
    Projects = "projects",
    Files = "files",
    TimeAggregated = "time aggregated"
}
</script>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import TotalTimeDisplay from './components/TotalTimeDisplay.vue'
import LanguageTimeDisplay from './components/LanguageTimeDisplay.vue'
import ProjectTimeDisplay from './components/ProjectTimeDisplay.vue'
import DisplaySwitch from './components/DisplaySwitch.vue'
import MyCalendar from './components/MyCalendar.vue'


const hasData = ref(true);
const data = ref<Data>();
const currentlyShown = ref<GroupBy>(GroupBy.TimeAggregated);

onMounted(async () => {
    const response = await fetch('/api/sessions')
    const json = await response.json()
    if (!(json.byLanguage && json.byProject && json.byTime)) {
        hasData.value = false;
        return;
    }
    data.value = json;
})

const currentData = computed(() => {
    console.log(currentlyShown.value)
    switch (currentlyShown.value) {
        case GroupBy.Languages:
            return data.value?.byLanguage;
        case GroupBy.Projects:
            return data.value?.byProject;
        case GroupBy.Files:
            console.log("halo", data.value?.byFile, data.value)
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
    console.log({ name: data.name, time: data.totalTime })
    return { name: data.name, time: data.totalTime }
}

const handleDisplayUpdate = (val: GroupBy): void => {
    currentlyShown.value = val;
}
</script>
