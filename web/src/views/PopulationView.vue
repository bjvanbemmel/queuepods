<template>
  <div class="h-screen flex flex-col justify-between">
  <ModalComponent
    title="Population trend"
    v-if="openAttraction !== ''"
      @close="() => openAttraction = ''"
  >
    <PopulationChart
      attraction="The Goliath"
    />
  </ModalComponent>
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
            class="group hover:cursor-pointer even:bg-gray-100 hover:bg-gray-200"
            :class="{'bg-yellow-400 hover:bg-yellow-500': population.population > population.capacity}"
            @click="() => openAttraction = population.attraction"
          >
            <td class="px-2">{{ key }}</td>
            <td class="px-2">{{ population.attraction }}</td>
            <td
              class="px-2 flex gap-2 items-center"
              :class="{'font-black underline': population.population > population.capacity}"
              title="Queue population exceeds capacity"
            >
              <ExclamationTriangleIcon
                class="h-6 fill-red-600"
                v-if="population.population > population.capacity"
              />
              {{ population.population }} / {{ population.capacity }}
            </td>
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
import ModalComponent from '@/components/ModalComponent.vue';
import PopulationChart from '@/components/PopulationChart.vue';
import { usePopulationStore } from '@/stores/populations';
import { Events } from '@/types/models';
import { ExclamationTriangleIcon } from '@heroicons/vue/16/solid';
import { computed, ref, type Ref } from 'vue';

const populations = usePopulationStore();
const state = computed(() => populations.populations.flatMap((x) => {
  switch (x.state) {
    case Events.QUEUE_FULL:
      return 'bg-red-400 group-hover:bg-red-500';
    case Events.QUEUE_EMPTY:
      return 'bg-green-400 group-hover:bg-green-500';
    case Events.QUEUE_ALMOST_EMPTY:
      return 'bg-green-300 group-hover:bg-green-400';
    case Events.QUEUE_ALMOST_FULL:
      return 'bg-orange-400 group-hover:bg-orange-500';
    default:
      return '';
  }
}))

const openAttraction: Ref<string> = ref('');

</script>
