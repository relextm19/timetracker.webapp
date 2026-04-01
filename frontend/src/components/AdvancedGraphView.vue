<template>
    <div class="w-full mt-6 rounded-xl border border-zinc-800 bg-zinc-950/70 p-4">
        <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-semibold text-white">Activity patterns</h3>
            <DisplaySwitch :options="Object.values(BucketMode)" :selected="BucketMode.Hour"
                @displayUpdated="onModeChanged" />
        </div>

        <div class="w-full rounded-lg border border-zinc-800 bg-zinc-900/50 p-3">
            <svg :viewBox="`0 0 ${chartWidth} ${chartHeight}`" class="w-full h-64">
                <polyline fill="none" stroke="#FFFFFF" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"
                    :points="linePoints" />

                <circle v-for="point in points" :key="point.label" :cx="point.x" :cy="point.y" r="3.5"
                    class="fill-cyan-300" />

                <g v-for="label in labels" :key="`${label.label}-x`">
                    <text :x="label.x" :y="label.y" text-anchor="middle" class="fill-zinc-400 text-[10px]">
                        {{ label.label }}
                    </text>
                </g>
                <g v-for="point in points" :key="`${point.label}-x`">
                    <text :x="point.x" :y="chartHeight - 6" text-anchor="middle" class="fill-zinc-400 text-[10px]">
                        {{ point.label }}
                    </text>
                </g>
            </svg>

        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import type { timeData } from '@/Dashboard.vue';
import DisplaySwitch from './DisplaySwitch.vue';
import { formatToTimeOfDay } from '@/utils/formatTime';

enum BucketMode {
    Hour = 'hours of day',
    Weekday = 'days of week',
    Month = 'months'
}

interface ChartPoint {
    label: string;
    labelY: number;
    x: number;
    y: number;
}

const props = defineProps<{
    byHour: timeData[]
    byWeekday: timeData[]
    byMonth: timeData[]
}>()

const mode = ref<BucketMode>(BucketMode.Hour)

const activeData = computed(() => {
    switch (mode.value) {
        case BucketMode.Weekday:
            return props.byWeekday;
        case BucketMode.Month:
            return props.byMonth;
        case BucketMode.Hour:
        default:
            return props.byHour;
    }
})

const maxTime = computed(() => {
    const max = activeData.value.reduce((acc, e) => Math.max(acc, e.totalTime), 0)
    return max || 1
})

const chartWidth = 800
const chartHeight = 300
const paddingX = 40
const paddingTop = 16
const paddingBottom = 36

const points = computed(() => {
    const values = activeData.value
    if (values.length === 0) return [] as ChartPoint[];

    const usableWidth = chartWidth - paddingX * 2
    const usableHeight = chartHeight - paddingTop - paddingBottom
    const stepX = values.length > 1 ? usableWidth / (values.length - 1) : 0
    return values.map((entry, index) => {
        const x = paddingX + (index * stepX)
        const ratio = entry.totalTime / maxTime.value
        const y = paddingTop + (usableHeight - ratio * usableHeight)
        return {
            label: entry.name,
            x,
            y,
        }
    })
})

const labels = computed(() => {
    const values = getLabels(maxTime.value);
    const usableHeight = chartHeight - paddingTop - paddingBottom
    return values.map((entry, index) => {
        const x = 0
        const y = paddingTop + (usableHeight / values.length) * (index + 1)
        return {
            label: entry.toString(),
            x,
            y,
        }
    })
})

const linePoints = computed(() => points.value.map(point => `${point.x},${point.y}`).join(' '))

const onModeChanged = (value: BucketMode): void => {
    mode.value = value
}

const getLabels = (max: number): string[] => {
    const res = [formatToTimeOfDay(max).slice(0, -1).join(":")];
    for (let i = Math.floor(max / 3600); i >= 0; i--) {
        res.push(i.toString() + ":00");
    }
    return res;
}
</script>
