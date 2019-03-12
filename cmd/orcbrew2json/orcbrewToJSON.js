require("node-go-require");
const Convert = require("./convert.go").Convert.Convert;
const { unlink, writeFile } = require("../../utils");

const dNow = new Date();
const sNow = `${dNow.getFullYear()}|${dNow.getMonth() +
  1}|${dNow.getDate()}|${dNow.getHours()}:${dNow.getMinutes()}:${dNow.getSeconds()}`;

/**
 * @returns {Promise<string>} returns a json string
 * @param {string} content expects the content of a valid .orcbrew file
 * @param {string?} filename
 */
function convert(content, filename = `${sNow}.temp.orcbrew`) {
  return writeFile(filename, content).then(() => {
    try {
      const json = Convert(filename);
      unlink(filename);
      return json;
    } catch (e) {
      console.debug(e);
      throw e;
    }
  });
}

module.exports = convert;
module.exports.Convert = Convert;
