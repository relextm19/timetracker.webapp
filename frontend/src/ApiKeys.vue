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
            <CopyToClipboardBtn :toCopy="newKeyValue" />
        </div>
        <span class="text-red-600">This is your only chance to save the key!</span>
        <button @click="showKeyModal = !showKeyModal"
            class="bg-black p-1 w-full text-white border border-white rounded-xs hover:bg-white hover:text-black transition duration-200 cursor-pointer focus:outline-0 focus:ring-1 focus:ring-white">
            Close
        </button>
    </Modal>
    <div :class="{ 'blur-xs': showNewKeyModal || showKeyModal }">
        <p class="text-3xl font-semibold uppercase tracking-wider mb-6 text-center text-white">
            API keys
        </p>
        <button @click="() => { showNewKeyModal = !showNewKeyModal }"
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
            <div v-for="key in APIKeys" :key="key.id" class="border-b border-gray-700">
                <ApiKeyDisplay :APIKey="key" @deleted="removeKey" />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import Modal from './components/Modal.vue';
import CopyToClipboardBtn from './components/CopyToClipboardBtn.vue';
import ApiKeyDisplay from './components/ApiKeyDisplay.vue';

export type APIKey = {
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

const removeKey = (id: number) => {
    APIKeys.value = APIKeys.value.filter(k => k.id !== id);
}

</script>
