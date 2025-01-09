<template>
  <div class="h-screen flex flex-col justify-between">
    <main>
      <table class="table-fixed">
        <thead>
          <tr class="bg-gray-200 border-b border-gray-400">
            <th class="px-2 py-1 text-left border-r border-gray-200">#</th>
            <th class="px-2 py-1 w-32 text-left border-r border-gray-200">Attraction</th>
            <th class="px-2 py-1 w-32 text-left border-r border-gray-200">Population</th>
            <th class="px-2 py-1 w-32 text-left border-r border-gray-200">State</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="population, key in populations.populations"
            :key="key"
            class="even:bg-gray-100"
          >
            <td class="px-2">{{ key }}</td>
            <td class="px-2">{{ population.attraction }}</td>
            <td class="px-2">{{ population.population }} / {{ population.capacity }}</td>
            <td
              class="px-2"
              :class="state[key]"
            >{{ population.state }}</td>
          </tr>
        </tbody>
      </table>
    </main>
  </div>
</template>

<script setup lang="ts">
import { usePopulationStore } from '@/stores/populations';
import { Events } from '@/types/models';
import { computed } from 'vue';

const populations = usePopulationStore();
const state = computed(() => populations.populations.flatMap((x) => {
  switch (x.state) {
    case Events.QUEUE_FULL:
      return 'bg-red-400';
    case Events.QUEUE_EMPTY:
      return 'bg-green-400';
    case Events.QUEUE_ALMOST_EMPTY:
      return 'bg-green-300';
    case Events.QUEUE_ALMOST_FULL:
      return 'bg-orange-400';
    default:
      return '';
  }
}))
</script>
