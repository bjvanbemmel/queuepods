import type { Population } from "@/types/models";
import { defineStore } from "pinia";
import { ref, type Ref } from "vue";

export const usePopulationStore = defineStore('populations', () => {
  const populations: Ref<Array<Population>> = ref([]);

  function set(value: Array<Population>): void {
    populations.value = value;
  }

  return { populations, set };
});
