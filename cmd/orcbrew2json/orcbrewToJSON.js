require("node-go-require");
const Convert = require("./convert.go").Convert.Convert;

const dNow = new Date();
const sNow = `${dNow.getFullYear()}|${dNow.getMonth() +
  1}|${dNow.getDate()}|${dNow.getHours()}:${dNow.getMinutes()}:${dNow.getSeconds()}`;

const sampleFile = "./samples/all-content.orcbrew";
const json = Convert(sampleFile);

require("fs").writeFileSync(
  `${sampleFile.split("/").pop()} - ${sNow}.json`,
  json,
  "utf-8"
);
console.debug("finished");
