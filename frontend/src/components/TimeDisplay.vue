<template>
    <div class="flex items-center gap-2 m-2">
        <span class="text-xs text-gray-500 uppercase tracking-wide">{{ props.title }}</span>
        <span class="font-mono  text-gray-300">{{formatToTimeOfDay(props.time).map(t => t.toString().padStart(2,
            "0")).join(":")}}</span>
    </div>
</template>

<script setup lang="ts">
import { formatToTimeOfDay } from '@/utils/formatTime';

const props = defineProps<{
    title: string,
    time: number
}>()
</script>
