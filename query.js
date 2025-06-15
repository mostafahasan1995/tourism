




let sum = 0; // temp sum variable 
const computeSum = (arr, n) => {
    if (n === arr.length) { // stop condition ===> when n = length of the array
        console.log(sum);   // print the result
        return; 
    }
    sum += arr[n]; // cumulative sum 
    computeSum(arr, n + 1); // recursive function
};
computeSum([7, 8, 9, 99 ,10], 0); // Expected output: 123


const chunk = (arr, size) => {
    let result = []; //final array
    for (let i = 0; i < arr.length; i += size) {
        //iterate on the array
        const splitArr = arr.slice(i, i + size); // split every 2 items of the array
        result = [...result, splitArr]; //spread to result array (push)
    }
    console.log(result);
};
chunk([1, 2, 3, 4, 5, 6, 7, 8, 9], 3)