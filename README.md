# orcbrew2js (Node.js)

This repository exports a Node module that runs the [parser](https://github.com/jnwhiteh/orcbrew-utils) that @jnwhiteh wrote in Golang behind the scenes to transform `.orcbrew` files into JS objects.

## Install

`npm i orcbrew2js`

## In a nutshell

Basically, all you have to do is extract the content from the `.orcbrew` file and pass it to the exported function, which returns a Promise; then you can simply `.then(obj => { ... })`, where `obj` is a JS object.
