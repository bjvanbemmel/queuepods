import type { Message, Param } from "@/types/models";
import { defineStore } from "pinia";
import { ref, type Ref } from "vue";

export const useMessageStore = defineStore('events', () => {
  const messages: Ref<Array<Message>> = ref([]);
  const params: Ref<Array<Param>> = ref([]);

  function set(value: Array<Message>): void {
    messages.value = value;
  }
  
  function setParam(param: Param): void {
    if (!params.value.find(x => x.name == param.name)) {
      params.value.push(param);
      return;
    }

    const index: number = params.value.findIndex(x => x.name == param.name)
    console.log(`${index}`);
    console.log(params.value);
    if (index === -1) return;

    params.value[index] = param;

  }

  function serializeParams(): string {
    return params.value.flatMap(x => `${x.name}=${x.value}`).join("&");
  }

  return { messages, set, params, setParam, serializeParams };
});
