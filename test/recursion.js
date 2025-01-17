const fs = require("fs");

const depths = {
	"1K": 1_000,
	"10K": 10_000,
	"100K": 100_000,
	"1M": 1_000_000,
	"10M": 10_000_000,
};

for (let depthName in depths) {
	const depth = depths[depthName];
	console.log(`Generating ${depth} depth object`);

	let jsonString = "{";

	for (let i = 1; i < depth; i++) {
		jsonString += '"next":{';
	}

	jsonString += '"next":null';

	for (let i = 0; i < depth; i++) {
		jsonString += "}";
	}

	const fileName = `${depthName}_recursion.json`;

	fs.writeFileSync(__dirname + `/${fileName}`, jsonString);
	console.log(`File for depth ${depth} saved as ${fileName}`);
}
