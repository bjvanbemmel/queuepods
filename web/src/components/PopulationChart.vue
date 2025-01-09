<template>
  <div class="flex justify-center items-center">
    <Line
      v-if="chartData.labels.length > 0"
      class="w-96"
      :options="chartOptions"
      :data="chartData"
      :key="key"
    />
    <div v-else class="p-8">
      <ArrowPathIcon class="h-5 animate-spin fill-gray-400" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Line } from 'vue-chartjs';
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip } from 'chart.js';
import { computed, onMounted, onUnmounted, ref, type Ref } from 'vue';
import { usePopulationStore } from '@/stores/populations';
import { Events, type Message } from '@/types/models';
import axios from 'axios';
import { ArrowPathIcon } from '@heroicons/vue/16/solid';

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip);

const props = defineProps<{
  attraction: string,
}>();

const populations = usePopulationStore();
const timestamps = computed<Array<string>>(() => populations.history.flatMap(x => new Date(x.timestamp).toLocaleTimeString('nl-NL')));
const values = computed<Array<number>>(() => populations.history.flatMap(x => +x.value));
const key: Ref<number> = ref(0);
const pollInterval: Ref<number> = ref(0);

onUnmounted(() => {
  clearInterval(pollInterval.value);
});

onMounted(async () => {
  pollInterval.value = setInterval(async () => {
    populations.setHistory(await fetchHistory())
    key.value += 1
  }, 1000);
});

const chartData = ref({
  labels: timestamps,
  datasets: [{
    data: values,
    borderColor: '#ff0505',
    pointRadius: 0,
    tension: 0.3,
    spanGaps: true,
  }],
});

const chartOptions = ref({
  responsive: true,
  animation: {
    duration: 0,
  },
  scales: {
    y: {
      suggestedMin: 0,
      suggestedMax: 150,
    }
  },
});

async function fetchHistory(): Promise<Array<Message>> {
  return await axios.get<Array<Message>>(
    `http://localhost:8888/messages?event=${Events.POPULATION_MONITORING}&attractions=${props.attraction}&from=${new Date(new Date().setSeconds(new Date().getSeconds() - 60)).toISOString()}`
  )
    .then(res => res.data)
    .catch((err) => {
      console.error(err);
      return [];
    })
}
</script>
