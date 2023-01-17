module.exports = {
  commonfate: {
    output: {
      clean: true,
      mode: "single",
      target: "./src/orval.ts",
      client: "swr",
      mock: false,
      // override: {
      //   mutator: {
      //     path: "./custom-instance.ts",
      //     name: "customInstanceRegistry",
      //   },
      // },
    },
    input: {
      target: "../openapi.yml",
    },
  },
};
