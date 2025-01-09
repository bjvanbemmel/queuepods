<template>
  <div
    class="z-50 absolute top-0 h-screen w-screen flex justify-center items-center bg-black/20 hover:cursor-pointer"
    @click.once="$emit('close')"
  >
    <div
      class="mb-64 absolute z-50 flex flex-col bg-white rounded-xl border border-gray-400 shadow shadow-gray-400 hover:cursor-default"
      @click.stop
    >
      <div class="rounded-t-xl bg-gray-100 p-2 flex gap-4 justify-between">
        <h1 class="font-bold">{{ props.title }}</h1>
        <FormButton
          @click="$emit('close')"
        >
          <XMarkIcon class="h-5" />
        </FormButton>
      </div>
      <div class="p-2 bg-white rounded-b-xl">
        <slot />
        <HorizontalSeperator class="my-2" />
        <FormButton
          @click="$emit('close')"
        >OK</FormButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { XMarkIcon } from '@heroicons/vue/16/solid';
import FormButton from './FormButton.vue';
import HorizontalSeperator from './HorizontalSeperator.vue';

const props = defineProps<{
  title: string,
}>();

</script>
