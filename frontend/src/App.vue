<script setup lang="ts">

import { IHit } from './interfaces/IEmail';
import { EmailService } from './services/EmailService'
import { onBeforeMount, ref, Ref, computed } from 'vue'

const sort = ref("")
const totalEmails: Ref<number> = ref(0)
const offset: Ref<number> = ref(0)
const limit: Ref<number> = ref(25)
const max: Ref<number> = ref(limit.value + offset.value)
const emailActive = ref<IHit | null>(null)
const emails: Ref<Array<IHit>> = ref([])
const search: Ref<string> = ref("")
const loading: Ref<boolean> = ref(false)

const emailService = new EmailService()

const loadInitialData = async (): Promise<void> => {
  await emailService.fetchAllEmails(sort.value, offset.value, limit.value)
  emails.value = emailService.getEmails().value
  totalEmails.value = emailService.getTotalEmails().value
}

onBeforeMount(async () => {
  loadInitialData()
})

const handleChangeLimit = async (e: Event): Promise<void> => {
  const target = e.target as HTMLSelectElement
  limit.value = parseInt(target.value)
  max.value = limit.value + offset.value


  loading.value = true

  if (search.value) {
    await emailService.fetchSearchEmails(search.value, sort.value, offset.value, limit.value)
  } else {
    await emailService.fetchAllEmails(sort.value, offset.value, limit.value)
  }

  loading.value = false
  emails.value = emailService.getEmails().value
}

const handleSearchEmails = async (e: Event): Promise<void>  => {
  e.preventDefault()
  const target = e.target as HTMLInputElement
  search.value = target.value
  emailActive.value = null
  offset.value = 0
  
  loading.value = true

  if (search.value) {
    await emailService.fetchSearchEmails(search.value, sort.value, offset.value, limit.value)
  } else {
    await emailService.fetchAllEmails(sort.value, offset.value, limit.value)
  }

  loading.value = false
  emails.value = emailService.getEmails().value
  totalEmails.value = emailService.getTotalEmails().value
}

const nextPage = async (): Promise<void> => {
  if (max.value > totalEmails.value) return

  offset.value += limit.value
  max.value = limit.value + offset.value

  loading.value = true

  if (search.value) {
    await emailService.fetchSearchEmails(search.value, sort.value, offset.value, limit.value)
  } else {
    await emailService.fetchAllEmails(sort.value, offset.value, limit.value)
  }

  loading.value = false
  emails.value = emailService.getEmails().value

}

const prevPage = async (): Promise<void> => {
  const errOffset: number = offset.value - limit.value
  if (errOffset < 0) return

  offset.value -= limit.value
  max.value = limit.value + offset.value

  loading.value = true

  if (search.value) {
    await emailService.fetchSearchEmails(search.value, sort.value, offset.value, limit.value)
  } else {
    await emailService.fetchAllEmails(sort.value, offset.value, limit.value)
  }

  loading.value = false
  emails.value = emailService.getEmails().value

}

const highlightText = (text: string) => {
  if (!search.value) return text;

  // Escapar caracteres especiales en `search.value`
  const escapedSearchValue = search.value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

  // Crear la expresión regular con el valor escapado
  const regex = new RegExp(`(${escapedSearchValue})`, 'gi');

  // Reemplazar las coincidencias con la marca <mark>
  return text.replace(regex, '<mark class="bg-yellow-300">$1</mark>');
};


const highlightedEmails = computed(() => {
  return emails.value.map(email => {
    return {
      ...email,
      _source: {
        ...email._source,
        subject_mark: highlightText(email._source.subject),
        content_mark: highlightText(email._source.content)
      }
    };
  });
});

const toggleSort = async (): Promise<void> => {
  sort.value = sort.value === "" ? "-" : ""
  loading.value = true

  if (search.value) {
    await emailService.fetchSearchEmails(search.value, sort.value, offset.value, limit.value)
  } else {
    await emailService.fetchAllEmails(sort.value, offset.value, limit.value)
  }

  console.log(sort.value);
  

  loading.value = false
  emails.value = emailService.getEmails().value
  totalEmails.value = emailService.getTotalEmails().value
}
  

</script>

<template>
  <main class="bg-cyan-700 min-h-screen px-[3em] pt-[1em] flex flex-col gap-[1em]">

    <!-- Web header -->
    <header class="flex justify-between items-center">

      <div class="flex-1 hidden md:block">
        <h1 class="text-2xl text-white">Emails received</h1>
      </div>

      <form class="flex items-center gap-[1em] bg-cyan-600 flex-1 max-w-full md:max-w-[40%] rounded-md">
        <input type="text" autocomplete="off" name="search" @input="handleSearchEmails"
          class="w-full bg-transparent text-white flex-1 outline-none border-none px-1">
        <button class="px-1">
          <v-icon name="la-search-solid" class="text-white"></v-icon>
        </button>
      </form>

    </header>

    <!-- web content -->
    <section class="w-full p-[0.5em] gap-2 bg-slate-100 rounded-md flex flex-col">

      <!-- Header in content -->
      <header class="flex justify-between" v-if="!emailActive">

        <div class="flex-1 max-w-[40%]  gap-2 items-center flex">
          <label for="countries"
            class="hidden text-sm font-medium text-gray-900 dark:text-white">Filas</label>
            <button @click="toggleSort">+</button>
          <select id="countries" @change="handleChangeLimit"
            class="bg-gray-50 border h-5 border-gray-300 text-gray-900 text-xs rounded-lg focus:ring-blue-500 focus:border-blue-500 block flex-1 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
            <option value="25">Ver 25 resultados</option>
            <option value="50">Ver 50 resultados</option>
            <option value="100">Ver 100 resultados</option>
            <option value="150">Ver 150 resultados</option>
          </select>
        </div>

        <div class="flex-1 flex gap-2 items-center justify-end">
          <p class="md:text-sm text-xs">{{ offset + 1 }} - {{ max }} de {{ totalEmails }}</p>
          <button @click="prevPage">
            <v-icon name="fa-arrow-left" class="text-cyan-700"></v-icon>
          </button>
          <button @click="nextPage">
            <v-icon name="fa-arrow-right" class="text-cyan-700"></v-icon>
          </button>
        </div>
      </header>

      <!-- Header in content | inside email -->
      <header v-else>
        <button @click="() => { emailActive = null }">
          <v-icon name="io-arrow-back-circle-sharp" class="text-cyan-700"></v-icon>
        </button>
      </header>

      <!-- Content -->
      <section class="w-full h-[85vh] overflow-auto relative" v-if="!emailActive">
        <div class="absolute w-full h-full flex justify-center items-center" v-if="loading">
          <v-icon name="io-reload-circle" scale="3" class="text-cyan-700 animate-spin"></v-icon>
        </div>

        <table class="w-full hidden md:table" v-if="!loading">
          <thead class="w-full">
            <tr>
              <!-- <th class="text-left w-1/3"></th>
              <th class="text-left w-1/3"></th>
              <th class="text-left w-1/3"></th> -->
            </tr>
          </thead>
          <tbody>
            <tr v-for="email in highlightedEmails" :key="email._id"
              class="border-y-2 hover:bg-gray-200 cursor-pointer w-full" @click="emailActive = email">
              <td class="w-[10em] pr-2 whitespace-nowrap overflow-hidden text-ellipsis">{{ email._source.x_from }}</td>
              <td class="whitespace-nowrap overflow-hidden text-ellipsis max-w-xs h-[1em]">
                <span> {{ email._source.subject }}</span> -
                <span class="opacity-50"> {{ email._source.content }}</span>
              </td>
              <td class="w-[10em] whitespace-nowrap overflow-hidden text-ellipsis text-right">{{ email._source.date }}</td>
            </tr>
          </tbody>
        </table>

        <div class="flex flex-col md:hidden" v-if="!loading">
          <div v-for="email in highlightedEmails" :key="email._id"
            class="flex gap-2 border-y-2 hover:bg-gray-200 cursor-pointer p-2" @click="emailActive = email">
            <div class="flex-2">
              <p class="text-sm">{{ email._source.x_from }}</p>
              <p class="max-w-[50vw] whitespace-nowrap overflow-hidden text-ellipsis text-sm">{{ email._source.subject }}</p>
              <p class="w-[50vw] whitespace-nowrap overflow-hidden text-ellipsis text-sm opacity-50"> {{ email._source.content }}</p>
            </div>
            <div class="flex-1 text-right">
              <p class="text-sm">{{ email._source.date }}</p>
            </div>
          </div>
        </div>

      </section>

      <!-- Content | inside email -->
      <section class="w-full h-[85vh] overflow-auto" v-if="emailActive">
        <div class="flex flex-col gap-2">
          <div class="flex gap-2 flex-col">
            <p><span class="font-semibold">From:</span> {{ emailActive._source.x_from }}</p>
            <p><span class="font-semibold">To:</span> {{ emailActive._source.x_to }}</p>
            <p><span class="font-semibold">Date:</span> {{ emailActive._source.date }}</p>
          </div>
          <h2 class="text-lg" v-html="emailActive._source.subject_mark"></h2>
          <p v-html="emailActive._source.content_mark"></p>
        </div>
      </section>

    </section>

  </main>
</template>
