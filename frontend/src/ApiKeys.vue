<template>
    <div v-for="key in APIKeys">
        {{ key }}
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

type APIKey = {
    name: String
    createdAt: Number
    keyHash: String
}

const APIKeys = ref<APIKey[]>([]);

onMounted(async () => {
    const res = await fetch('/api/api_keys');
    const json = await res.json() as APIKey[];
    for (const key of json) {
        APIKeys.value.push(key)
    }
})
</script>
