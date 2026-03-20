<template>
    <p class="text-3xl font-semibold uppercase tracking-wider mb-6 text-center text-white">
        API keys
    </p>
    <div v-if="APIKeys.length === 0" class="text-center text-gray-400">
        No API keys yet
    </div>
    <div v-else class="space-y-4">
        <div v-for="(key, index) in APIKeys" :key="index" class="border-b border-gray-700">
            <div class="flex">
                <div class="text-white flex-4/5">
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
                <button class="text-gray-400 transition-colors rounded-md hover:text-red-500 hover:bg-red-50"
                    aria-label="Delete API Key">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                        stroke="currentColor" class="w-5 h-5">
                        <path stroke-linecap="round" stroke-linejoin="round"
                            d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                    </svg>
                </button>
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
