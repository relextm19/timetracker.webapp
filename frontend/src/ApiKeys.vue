<template>
    <p class="text-3xl font-semibold uppercase tracking-wider mb-6 text-center text-white">
        API keys
    </p>
    <div v-if="APIKeys.length === 0" class="text-center text-gray-400">
        No API keys yet
    </div>
    <div v-else class="space-y-4">
        <div v-for="(key, index) in APIKeys" :key="index" class="border-b border-gray-700 pb-4  ">
            <div class="text-white">
                <div class="mb-2">
                    <span class="text-gray-400 text-sm">Name:</span>
                    <span class="ml-2">{{ key.name }}</span>
                </div>
                <div class="mb-2">
                    <span class="text-gray-400 text-sm">Key Hash:</span>
                    <span class="ml-2 font-mono text-sm">{{ key.keyHash }}</span>
                </div>
                <div>
                    <span class="text-gray-400 text-sm">Created At:</span>
                    <span class="ml-2">{{ formatToDate(key.createdAt) }}</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { formatToDate } from './utils/formatTime';

type APIKey = {
    name: string
    createdAt: number
    keyHash: string
}

const APIKeys = ref<APIKey[]>([]);

onMounted(async () => {
    const res = await fetch('/api/api_keys');
    const json = await res.json() as APIKey[];
    APIKeys.value.push(...json)
})
</script>
