import type { Message, Population } from "@/types/models";
import { defineStore } from "pinia";
import { ref, type Ref } from "vue";

export const usePopulationStore = defineStore('populations', () => {
  const populations: Ref<Array<Population>> = ref([]);
  const history: Ref<Array<Message>> = ref([]);

  function set(value: Array<Population>): void {
    populations.value = value;
  }

  function setHistory(value: Array<Message>): void {
    history.value = value;
  }

  return { populations, history, set, setHistory };
});
