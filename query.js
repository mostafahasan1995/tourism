


function computeSum (arr) {
  // The base case is when there is no elements in the array
  if(arr.length === 0) return 0;

  // for each element of the array am just splicing the array and adding the last element to the whole sum
  return computeSum(arr.splice(0, arr.length - 1)) + arr[arr.length - 1];
}

console.log(computeSum([]));



function chunk(arr, size) {
  // the final result array
  let result = [];
  // a pointer to track array elements
  let pointer = 0;
  // array to put on the original array elements
  let subArray = [];

  // while loop to go through all array elements
  while (pointer !== arr.length) {
    // checking if the subArray length is less than the required size, if so the element get pushed to the subArray and increase the pointer by 1
    if (subArray.length < size) {
      subArray.push(arr[pointer]);
      pointer++;
    } else { // reaching the required size, push the subArray to the final result and clear the subArray 
      result.push(subArray);
      subArray = [];
    }
  }

  return result;
}


console.log(chunk([1, 2, 3, 4, 5, 6, 7, 8, 9, 10], 2));
