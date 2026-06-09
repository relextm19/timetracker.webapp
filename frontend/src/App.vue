<template>
    <header class="bg-black text-white shadow-md">
        <div class="flex items-center justify-between h-16 px-8 w-full">
            <div class="text-2xl font-bold">
                TimeTracker
            </div>

            <nav class="flex gap-6 items-center">
                <MenuElement to="/keys" text="Api keys" />
                <MenuElement to="/" text="Dashboard" />
                <button
                    v-if="canLogout"
                    class="relative px-1 py-2 text-lg font-medium transition-colors duration-200 hover:text-gray-300 cursor-pointer"
                    @click="handleLogout"
                >
                    Logout
                </button>
            </nav>
        </div>
    </header>

    <main class="bg-black text-white min-h-screen">
        <router-view></router-view>
    </main>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import router, { setLoggedIn } from './router';
import MenuElement from './components/MenuElement.vue';

const route = useRoute();
const canLogout = computed(() => !route.meta.public);

async function handleLogout() {
    try {
        const response = await fetch('/api/logout', {
            method: 'POST',
            credentials: 'include',
        });

        if (!response.ok && response.status !== 401) {
            return;
        }

        setLoggedIn(false);
        router.replace('/login');
    } catch (error) {
        console.error('Error during logout:', error);
    }
}
</script>
