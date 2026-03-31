<template>
    <div class="flex justify-between items-center p-2">
        <TimeDisplay title="This Day" :time="aggregatedTimes.day" />
        <TimeDisplay title="This Week" :time="aggregatedTimes.week" />
        <TimeDisplay title="This Month" :time="aggregatedTimes.month" />
        <TimeDisplay title="This Year" :time="aggregatedTimes.year" />
    </div>
    <Calendar transparent borderless expanded :is-dark="true" :attributes="calendarAttributes"></Calendar>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import 'v-calendar/style.css';
import { Calendar } from 'v-calendar';
import type { timeData } from '@/Dashboard.vue';
import { formatToTimeOfDay } from '@/utils/formatTime';
import TimeDisplay from './TimeDisplay.vue';

const props = defineProps<{ timeData: timeData[] }>()

//maybe do this in the map so the loop runs once
const aggregatedTimes = computed(() => {
    let day = 0, week = 0, month = 0, year = 0;

    const now = new Date();
    const todayString = now.toDateString();
    const currentMonth = now.getMonth();
    const currentYear = now.getFullYear();

    const startOfWeek = new Date(now);
    const dayOfWeek = startOfWeek.getDay() || 7;
    startOfWeek.setDate(startOfWeek.getDate() - dayOfWeek + 1);
    startOfWeek.setHours(0, 0, 0, 0);

    const endOfWeek = new Date(startOfWeek);
    endOfWeek.setDate(startOfWeek.getDate() + 6);
    endOfWeek.setHours(23, 59, 59, 999);

    props.timeData.forEach(entry => {
        const entryDate = new Date(entry.name);
        const time = entry.totalTime;

        if (time <= 0) return;

        if (entryDate.toDateString() === todayString) {
            day += time;
        }

        if (entryDate >= startOfWeek && entryDate <= endOfWeek) {
            week += time;
        }

        if (entryDate.getMonth() === currentMonth && entryDate.getFullYear() === currentYear) {
            month += time;
        }

        if (entryDate.getFullYear() === currentYear) {
            year += time;
        }
    });

    return { day, week, month, year };
});

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
