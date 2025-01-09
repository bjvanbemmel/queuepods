<template>
  <div class="h-screen flex flex-col justify-between">
    <nav class="w-full bg-gray-300 flex gap-4 items-center p-4">
      <RouterLink to="/">Population</RouterLink>
      <RouterLink to="/messages">Messages</RouterLink>
    </nav>
    <RouterView />
    <footer class="z-10 p-3 flex justify-between items-center w-full bg-gray-100 border-t border-gray-400">
      <div id="footer-left" class="flex items-center">
        <p v-if="lastPoll">Last update: {{ lastPoll }}</p>
      </div>
      <div id="footer-right" class="h-full flex justify-end items-center gap-6 ">
        <div class="flex items-center gap-2">
          <span>Polling every</span>
          <select
            class="p-1 border border-gray-400 rounded focus:border-blue-400 bg-gray-50 hover:bg-gray-200"
            v-model="pollingRate"
          >
            <option
              v-for="val, key in range(1, 11)"
              :value="val * 1000"
              :key="key"
              :selected="val === 3"
            >
              {{val}}
            </option>
          </select>
          <span>seconds</span>
        </div>
        <VerticalSeperator />
        <FormButton
          @click="updateData"
        >
          <ArrowPathIcon class="h-6" />
        </FormButton>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, type Ref, watch } from 'vue';
import axios from 'axios';
import { range } from 'lodash';
import { ArrowPathIcon } from '@heroicons/vue/16/solid';
import VerticalSeperator from '@/components/VerticalSeperator.vue';
import type { Message, Param, Population } from '@/types/models';
import { usePopulationStore } from './stores/populations';
import FormButton from './components/FormButton.vue';
import { useMessageStore } from './stores/messages';

const populations = usePopulationStore();
const messages = useMessageStore();

const pollingRate: Ref<number> = ref(-1);
const pollingTimeout: Ref<number> = ref(0);
const lastPoll: Ref<string | null> = ref(null);

onMounted(async () => {
  pollingRate.value = 3000;
  messages.setParam({ name: 'limit', value: '500' } as Param);
  await updateData()
});

watch(pollingRate, async (rate: number) => {
  clearInterval(pollingTimeout.value);

  pollingTimeout.value = setInterval(async () => {
    await updateData();
  }, rate);
})

async function updateData(): Promise<void> {
  populations.set(await fetchPopulations());
  messages.set(await fetchMessages());
}

async function fetchPopulations(): Promise<Array<Population>> {
  return axios.get<Array<Population>>('http://localhost:8888/populations')
    .then(res => {
      lastPoll.value = new Date().toLocaleString('nl-NL');

      return res.data;
    })
    .catch((err) => {
      console.error(err);
      return [];
    });
}

async function fetchMessages(): Promise<Array<Message>> {
  return axios.get<Array<Message>>(`http://localhost:8888/messages?${messages.serializeParams()}`)
    .then((res) => {
      lastPoll.value = new Date().toLocaleString('nl-NL');

      return res.data;
    })
    .catch((err) => {
      console.error(err);
      return [];
    })
}
</script>
