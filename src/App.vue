<template>
  <div>
    <input
      style="display: none;"
      ref="inputRef"
      type="file"
      multiple
      :disabled="files.length >= 2"
      @change="handleFileChange"
      placeholder="Choose Two Files"
    />
    <button
      type="button"
      :disabled="loading || files.length >= 2"
      @click="handleSelectFiles"
    >Select Two Files</button>
    <button
      type="button"
      :disabled="compareReady || files.length < 2"
      @click="handleParse"
    >Parse Files</button>
    <button type="button" :disabled="!compareReady" @click="handleGetTable">Get Table</button>
    <div>
      <p v-for="(item, index) in logList" :key="index">{{item}}</p>
    </div>
  </div>
</template>

<script>
import { defineComponent, onMounted, watchEffect, reactive, ref } from "vue";
import initWasm from "./assets/main.wasm";
import { useGoWasm } from "./hooks/useWasm";
import { FileStream } from "./utils/fileStream";

export default defineComponent({
  name: "App",

  setup() {
    const logList = reactive([]);
    const rawLog = console.log;
    console.log = function (...args) {
      rawLog(...args);
      logList.push(args.toString());
    };

    const compareReady = ref(false);
    const inputRef = ref(null);
    const files = ref([]);
    const { exports, loading } = useGoWasm();

    const table = ref(undefined);

    function handleFileChange(event) {
      const { target } = event;
      const { files: inputFiles } = target;
      const fileArray = Array.from(inputFiles);
      files.value = [...files.value, ...fileArray];
      if (files.value.length > 2) {
        files.value.length = 2;
      }
      console.log("selected files");
      files.value.forEach((file) => {
        console.log(`filename: ${file.name}; size: ${file.size};`);
      });
    }

    function handleSelectFiles() {
      if (inputRef.value) {
        inputRef.value.click();
      }
    }

    async function handleGetTable() {
      try {
        const resp = await exports.getJSON(
          "",
          5,
          true,
          true,
          true,
          0,
          true,
          1,
          Number.MAX_SAFE_INTEGER
        );
        console.log('success');
      } catch (error) {
        console.log(`Open Devtools to see the Error`);
        console.log(error);
      }
    }

    async function parseFile(file, i) {
      console.log(`Parsing ${file.name}`);
      const fileReader = new FileStream(file);
      for (;;) {
        const resp = await fileReader.readChunk();
        await exports.parse(i, resp.eof, resp.data, resp.data.length);

        if (resp.eof) {
          console.log(`${file.name} has been parsed`);
          break;
        }
      }
    }

    async function handleParse() {
      if (files.value.length === 2 && !loading.value) {
        await parseFile(files.value[0], 0);
        await parseFile(files.value[1], 1);
        compareReady.value = true;
      }
    }

    return {
      inputRef,

      compareReady,
      exports,
      loading,
      files,
      logList,
      table,

      handleFileChange,
      handleSelectFiles,
      handleParse,
      handleGetTable,
    };
  },
});
</script>
