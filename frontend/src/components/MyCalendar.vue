<template>
    <Calendar transparent borderless expanded :is-dark="true" :attributes="calendarAttributes"></Calendar>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import 'v-calendar/style.css';
import { Calendar } from 'v-calendar';
import type { timeData } from '@/Dashboard.vue';
import { formatToTimeOfDay } from '@/utils/formatTime';

const props = defineProps<{ timeData: timeData[] }>()

const calendarAttributes = computed(() => {
    return props.timeData
        .filter(entry => entry.totalTime > 0)
        .map(entry => {
            return {
                key: entry.name,
                highlight: 'blue',
                dates: [new Date(entry.name)],
                popover: {
                    label: `Time logged: ${formatToTimeOfDay(entry.totalTime).join(":")}`
                }
            };
        });
});
</script>
