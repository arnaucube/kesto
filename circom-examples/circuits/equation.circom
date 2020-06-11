include "../node_modules/circomlib/circuits/comparators.circom";

/*
	Circuit to check that prover knows a private inputs x & y,
	such as x**a - y**b = 0, where a & b are known parameters.
*/
template Equation(a, b) {
	signal private input x;
	signal private input y;

	
	signal aArray[a];
	signal bArray[b];
	for (var i=0; i<a; i++) {
		if (i==0) {
			aArray[0] <== x;
		} else {
			aArray[i] <== x * aArray[i-1];
		}
	}
	for(var i=0; i<b; i++) {
		if (i==0) {
			bArray[0] <== y;
		} else {
			bArray[i] <== y * bArray[i-1];
		}
	}
	
	component checkEq = IsEqual();
	checkEq.in[0] <== aArray[a-1];
	checkEq.in[1] <== bArray[b-1];
	checkEq.out === 1;
}

component main = Equation(4, 8);
