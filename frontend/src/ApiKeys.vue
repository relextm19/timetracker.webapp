<template>
    <Modal v-if="showNewKeyModal" title="Add API key">
        <form @submit.prevent="addKey" class="flex justify-center items-center flex-col gap-4">
            <input v-model="newKeyName" type="text"
                class="border-2 border-gray-400 rounded-xs focus:outline-0 focus:ring-1 focus:ring-white"
                placeholder="Key name">
            <input
                class="bg-transparent p-1 w-full text-white border border-white rounded-xs hover:bg-white hover:text-black transition duration-200 cursor-pointer focus:outline-0 focus:ring-1 focus:ring-white"
                type="submit" value="Add">
        </form>
    </Modal>
    <Modal v-if="showKeyModal" title="Your API key">
        <div class="flex justify-center items-center">
            <span>{{ newKeyValue }}</span>
            <button @click="copyToClipboard(newKeyValue)"
                class="relative w-9 h-9 text-gray-400 transition-colors rounded-md hover:text-blue-500 hover:bg-blue-50/10 overflow-hidden"
                aria-label="Copy API Key">

                <span class="absolute inset-0 flex items-center justify-center transition-all duration-300 ease-in-out"
                    :class="copied ? 'opacity-0 scale-50 rotate-90' : 'opacity-100 scale-100 rotate-0'">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                        class="bi bi-clipboard-x" viewBox="0 0 16 16">
                        <path fill-rule="evenodd"
                            d="M6.146 7.146a.5.5 0 0 1 .708 0L8 8.293l1.146-1.147a.5.5 0 1 1 .708.708L8.707 9l1.147 1.146a.5.5 0 0 1-.708.708L8 9.707l-1.146 1.147a.5.5 0 0 1-.708-.708L7.293 9 6.146 7.854a.5.5 0 0 1 0-.708" />
                        <path
                            d="M4 1.5H3a2 2 0 0 0-2 2V14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V3.5a2 2 0 0 0-2-2h-1v1h1a1 1 0 0 1 1 1V14a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V3.5a1 1 0 0 1 1-1h1z" />
                        <path
                            d="M9.5 1a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-3a.5.5 0 0 1-.5-.5v-1a.5.5 0 0 1 .5-.5zm-3-1A1.5 1.5 0 0 0 5 1.5v1A1.5 1.5 0 0 0 6.5 4h3A1.5 1.5 0 0 0 11 2.5v-1A1.5 1.5 0 0 0 9.5 0z" />
                    </svg>
                </span>

                <span
                    class="absolute inset-0 flex items-center justify-center transition-all duration-300 ease-in-out text-green-500"
                    :class="copied ? 'opacity-100 scale-100 rotate-0' : 'opacity-0 scale-50 -rotate-90'">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2"
                        stroke="currentColor" class="w-5 h-5">
                        <path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
                    </svg>
                </span>
            </button>
        </div>
        <span class="text-red-600">This is your only chance to save the key!</span>
        <button @click="showKeyModal = !showKeyModal; copied = false;"
            class="bg-black p-1 w-full text-white border border-white rounded-xs hover:bg-white hover:text-black transition duration-200 cursor-pointer focus:outline-0 focus:ring-1 focus:ring-white">
            Close
        </button>
    </Modal>
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
import Modal from './components/Modal.vue';

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

const copied = ref(false)
const copyToClipboard = async (text: string) => {
    try {
        await navigator.clipboard.writeText(text);
        copied.value = true;
    } catch (err) {
        console.error("Failed to copy!", err);
    }
};
</script>
