const fs = require("fs");

module.exports = {
  writeFile(filename, content) {
    return new Promise((resolve, reject) => {
      fs.writeFile(filename, content, "utf-8", err => {
        if (err) {
          reject(err);
        }
        resolve();
      });
    });
  },
  unlink(filename) {
    return new Promise((resolve, reject) => {
      fs.unlink(filename, err => {
        if (err) {
          reject(err);
        }
        resolve();
      });
    });
  }
};
