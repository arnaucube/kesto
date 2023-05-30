# babyjubjub-wasm-example

- To run it locally needs https://rustwasm.github.io/wasm-pack/installer/ and https://rustwasm.github.io/docs/book/game-of-life/setup.html
- Rust code lives in `./src/` directory
- Then, each time that the rust code is modified, just run `./re-build.sh`, which will call wasm-pack build and will serve the html files
- You can now go to http://127.0.0.1:8080 and the wasm will be running there
