<template>
    <div v-if="showNewKeyModal"
        class="absolute top-1/2 left-1/2 -translate-1/2 flex justify-center items-center flex-col border-2 border-white rounded-lg gap-4 p-4 z-50 bg-black">
        <h2 class="text-2xl tracking-wider">Add API key</h2>
        <form @submit.prevent="addKey" class="flex justify-center items-center flex-col gap-4">
            <input v-model="newKeyName" type="text"
                class="border-2 border-gray-400 rounded-xs focus:outline-0 focus:ring-1 focus:ring-white"
                placeholder="Key name">
            <input
                class="bg-transparent p-1 w-full text-white border border-white rounded-xs hover:bg-white hover:text-black transition duration-200 cursor-pointer focus:outline-0 focus:ring-1 focus:ring-white"
                type="submit" value="Add">
        </form>
    </div>
    <div v-if="showKeyModal"
        class="absolute top-1/2 left-1/2 -translate-1/2 flex justify-center items-center flex-col border-2 border-white rounded-lg gap-4 p-4 z-50 bg-black">
        <h2 class="text-2xl tracking-wider">Your API key</h2>
        <div>
            <span>{{ newKeyValue }}</span>
            <button @click="copyToClipboard(newKeyValue)"
                class="p-2 text-gray-400 transition-colors rounded-md hover:text-blue-500 hover:bg-blue-50"
                aria-label="Copy API Key">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                    stroke="currentColor" class="w-5 h-5">
                    <path stroke-linecap="round" stroke-linejoin="round"
                        d="M15.75 17.25v3.375c0 .621-.504 1.125-1.125 1.125h-9.75a1.125 1.125 0 0 1-1.125-1.125V7.875c0-.621.504-1.125 1.125-1.125H6.75a9.06 9.06 0 0 1 1.5-1.5H13.5a3.75 3.75 0 0 1 3.75 3.75v6.75a9.06 9.06 0 0 1-1.5 1.5Zm0 0c.342 0 .684.048 1.012.144.405.118.784.329 1.11.616l2.122 2.122c.287.287.498.666.616 1.071m-9.444-11.444L15 9.75M21 12v4.5A2.25 2.25 0 0 1 18.75 18.75h-4.5" />
                </svg>
            </button>
        </div>
        <span class="text-red-600">This is your only chance to save the key!</span>
        <button @click="showKeyModal = !showKeyModal"
            class="bg-black p-1 w-full text-white border border-white rounded-xs hover:bg-white hover:text-black transition duration-200 cursor-pointer focus:outline-0 focus:ring-1 focus:ring-white">
            Close
        </button>
    </div>
    <div :class="{ 'blur-xs': showNewKeyModal || showKeyModal }">
        <p class="text-3xl font-semibold uppercase tracking-wider mb-6 text-center text-white">
            API keys
        </p>
        <button @click="() => { showNewKeyModal = !showNewKeyModal; console.log(showNewKeyModal) }"
            class="p-2 text-gray-400 transition-colors rounded-md hover:text-green-600 flex justify-start items-center gap-2"
            aria-label="Add API Key">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                stroke="currentColor" class="w-5 h-5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
            </svg>
            <span>New key</span>
        </button>
        <div v-if="APIKeys.length === 0" class="text-center text-gray-400">
            No API keys yet
        </div>
        <div v-else class="space-y-4">
            <div v-for="(key, index) in APIKeys" :key="index" class="border-b border-gray-700">
                <div class="flex p-2">
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
                    <button @click="() => deleteKey(key.id)"
                        class="text-gray-400 transition-colors rounded-md hover:text-red-500 hover:bg-red-50"
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
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { formatToDate } from './utils/formatTime';

type APIKey = {
    id: number
    name: string
    createdAt: number
    keyHash: string
    key: string
}

const showNewKeyModal = ref(false)
const showKeyModal = ref(false)

const APIKeys = ref<APIKey[]>([]);

onMounted(async () => {
    const res = await fetch('/api/keys');
    const json = await res.json() as APIKey[];
    console.log(json)
    APIKeys.value.push(...json)
})

const deleteKey = async (id: number): Promise<void> => {
    await fetch('/api/keys/' + id, {
        method: "DELETE"
    });
    APIKeys.value = APIKeys.value.filter(k => k.id !== id)
}

const newKeyName = ref("")
const newKeyValue = ref("")

const addKey = async () => {
    if (!newKeyName.value) return;

    try {
        const res = await fetch('/api/keys', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: newKeyName.value
            })
        });

        if (res.ok) {
            const json = await res.json() as APIKey
            APIKeys.value.push(json)
            newKeyName.value = '';
            showNewKeyModal.value = false;
            showKeyModal.value = true;
            newKeyValue.value = json.key;
        }
    } catch (error) {
        console.error("Failed to add key:", error);
    }
}

const copyToClipboard = async (text: string) => {
    try {
        await navigator.clipboard.writeText(text);
    } catch (err) {
        console.error("Failed to copy!", err);
    }
};
</script>
