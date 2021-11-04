const YAML = require("yaml");
const glob = require("glob");
const path = require("path");
const fs = require("fs");

const serviceYamls = glob.sync(path.join("trackiam/services/*.yml"));

const services = {};
serviceYamls.forEach(sy => {
	const contents = fs.readFileSync(sy).toString();
	const service = path.basename(sy).replace(".yml", "");
	const parsedYaml = YAML.parse(contents);
	services[service] = parsedYaml.Actions.map(a => a.Name);
});
fs.writeFileSync("public/iamactions.json", JSON.stringify(services));
