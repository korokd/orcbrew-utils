const homedir = require("os").homedir();
process.env.GOPATH = `${homedir}/go`;

const convert = require("./cmd/orcbrew2json/orcbrewToJSON");
module.exports = convert;

// require("fs").readFile(
//   "./samples/all-content.orcbrew",
//   "utf-8",
//   (err, data) => {
//     if (err) {
//       throw err;
//     }
//     convert(data).then(json => {
//       require("fs").writeFileSync("daledale", json, "utf-8");
//     });
//   }
// );
