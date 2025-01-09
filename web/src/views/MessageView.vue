<template>
  <main class="flex h-screen overflow-scroll">
    <div class="w-64"></div>
    <div class="fixed h-full w-64 bg-gray-100 border-r border-gray-400">
      <h1 class="w-full p-1 bg-gray-200 font-bold border-b border-gray-400">Filters</h1>
      <form @submit.prevent class="flex flex-col gap-8 p-1">
        <div>
          <div
            class="flex items-center gap-2"
            v-for="filter, key in filters"
            :key="key"
          >
            <input
              type="checkbox"
              :id="filter.name"
              :name="filter.name"
              v-model="filter.value"
              :checked="filter.value"
            />
            <label
              :for="filter.name"
            >{{filter.name}}</label>
          </div>
        </div>
        <div class="flex gap-2 items-center">
          <label>Limit</label>
          <input
            class="w-32"
            type="number"
            v-model="limit.value"
          />
        </div>
      </form>
    </div>
    <div>
      <table class="table-fixed">
        <thead>
          <tr class="bg-gray-200 border-b border-gray-400">
            <th class="px-2 py-1 text-left border-r border-gray-200">#</th>
            <th class="px-2 py-1 w-32 text-left border-r border-gray-200">Attraction</th>
            <th class="px-2 py-1 w-40 text-left border-r border-gray-200">Timestamp</th>
            <th class="px-2 py-1 w-32 text-left border-r border-gray-200">Event</th>
            <th class="px-2 py-1 text-left border-r border-gray-200">Value</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="message, key in messages.messages"
            :key="key"
            class="even:bg-gray-100"
          >
            <td class="px-2">{{ key }}</td>
            <td class="px-2" :title="message.attraction">{{ message.attraction }}</td>
            <td class="px-2" :title="new Date(message.timestamp).toLocaleString('nl-NL')">{{ new Date(message.timestamp).toLocaleString('nl-NL') }}</td>
            <td class="px-2">{{ message.event }}</td>
            <td class="px-2">{{ message.value }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </main>
</template>

<script setup lang="ts">
import { useMessageStore } from '@/stores/messages';
import { Events, Params, type Param } from '@/types/models';
import { ref, watch, type Ref } from 'vue';

interface Filter {
  name:  Events,
  value: boolean,
};

const messages = useMessageStore();
const limit: Ref<Param> = ref({
  name: Params.LIMIT,
  value: '200',
});
const filters: Ref<Array<Filter>> = ref([
  {
    name: Events.QUEUE_FULL,
    value: true,
  },
  {
    name: Events.QUEUE_EMPTY,
    value: true,
  },
  {
    name: Events.QUEUE_ALMOST_EMPTY,
    value: true,
  },
  {
    name: Events.QUEUE_ALMOST_FULL,
    value: true,
  },
  {
    name: Events.POPULATION_MONITORING,
    value: false,
  },
]);

watch(filters, (newFilters) => {
  const eventsParam: Param = {
    name: Params.EVENTS,
    value: newFilters.filter(x => x.value).flatMap(x => x.name).join(','),
  };

  messages.setParam(eventsParam);
}, { deep: true, immediate: true });

watch(limit, (newLimit: Param) => {
  messages.setParam(newLimit);
}, { immediate: true, });

</script>
