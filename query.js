

// function computeSum(arr) {
//   // Break condition 1: one element left
//   if (arr.length === 1) {
//       return arr[0];
//   }
//   // Break condition 2: No elements in the array
//   if (arr.length === 0) {
//       return 0;
//   }
//   // return first element and the result of the function with the rest of the array
//   return arr[0] + computeSum(arr.slice(1));
// }


// function computeSum (arr) {
//   // The base case is when there is no elements in the array
//   if(arr.length === 0) return 0;

//   // for each element of the array am just splicing the array and adding the last element to the whole sum
//   return computeSum(arr.splice(0, arr.length - 1)) + arr[arr.length - 1];
// }


// console.log(computeSum([1,3]));




// function chunk(arr, size) {
//   // the final result array
//   let result = [];
//   // a pointer to track array elements
//   let pointer = 0;
//   // array to put on the original array elements
//   let subArray = [];

//   // while loop to go through all array elements
//   while (pointer !== arr.length) {
//     // checking if the subArray length is less than the required size, if so the element get pushed to the subArray and increase the pointer by 1
//     if (subArray.length < size) {
//       subArray.push(arr[pointer]);
//       pointer++;
//     } else { // reaching the required size, push the subArray to the final result and clear the subArray 
//       result.push(subArray);
//       subArray = [];
//     }
//   }

//   return result;
// }


// console.log(chunk([1,2,3,4,5,6,7,8,9,10], 4));



// function chunk(arr, size) {
//   let results = [];
//   // iterate the array
//   let countOfChunks = Math.floor(arr.length / size)
//   let j = 0;
//   // collect each chunk from [chunkNum * size, chunkNum * size]
//   while (j < countOfChunks) { // 0 3
//       results.push(arr.slice(j * size, ((j + 1) * size) - 1))
//       j++;
//   }
//   // check if there will be remaining of the chunking by dividing the length on the size of the chunk
//   if (arr.length % size !== 0) {
//       let startIdx = countOfChunks * size;
//       results.push(arr.slice(startIdx))
//   }
//   return results;
// }


// console.log(chunk([1,2,3,4,5,6,7,8,9,10], 4));


function computeSum(arr)
{
    let result = 0; // initial result set to 0
    arr.map( number => result += number); // iterating through the array and adding the number to the sum
    return result; // returning the final sum
}
console.log(computeSum([7, 8, 9, 99]));


function chunk(array, size)
{
    let result = [];
    const numberOfArrays = Math.ceil(array.length / size); // get the number of arrays I need to create
    
    for(let i = 0; i < numberOfArrays; i++) // looping to create chunks
    {
        result.push([]); // creat a new chunk
        const start = i * size; // calculate the offset of the original array values to pe appended in the chunk
        let end =  start + size  // calculate the limit
        end = end > ( array.length ) ?  ( array.length ) : end; // check if the limit is bigger than the original array length to prevent errors
        for (let j = start; j < end; j++)
        {
            result[i].push(array[j]) // append the value in the designated chunk
        }
    }
    return result
}
console.log(chunk([1,2,3,4,5,6,7,8,9,10], 3));