# Golang Wasm Debug
[issue 40923](https://github.com/golang/go/issues/40923)

## Setup
```shell
# Before setup this demo, please make sure your environment has nodejs >= 10 and go >= 1.14

npm install
npm run dev

# or
yarn
yarn dev
```

## Test Step
- [download test files](https://drive.google.com/drive/folders/11outqa4qFCZoHeNIR0R9VUkhhgDj6pZU?usp=sharing)
- clone this repo
- run dev server
- select two files under testFile
- when both two file read, click parse
- when both two file parsed, click getTable
  - if log `success`, there is no error
  - but, the test files will cause error, open your devtools and see the error trace
