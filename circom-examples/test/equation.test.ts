const path = require("path");
const tester = require("circom").tester;
const chai = require("chai");
const assert = chai.assert;

export {};

describe("equation test", function () {
    this.timeout(200000);


    it("Test equation", async () => {
        const circuit = await tester(
            path.join(__dirname, "../circuits", "equation.circom"),
            {reduceConstraints: false}
        );
        
        const witness = await circuit.calculateWitness({
            "x": 4,
            "y": 2
        });
        await circuit.checkConstraints(witness);
    });
});
