




function chunk(arr, size){
  let newArr = []; // buffer for the chunk array

  for (let i = 0; i < arr.length / size ; i++)
  {
      let buf = []; // subarray buffer
      for (let j = 0 ; j < size ; j++)
      {
          if (!arr[j + (i * size)]) // check if there's an item with the corresponding index
              continue;

          buf.push(arr[j + (i * size)]); // push to the subarray
      }
      newArr.push(buf); // push the subarray to the chunk array
  }
  return newArr;
}

console.log(chunk([1, 2, 3, 4, 5, 5, 6, 7, 8, 9, 10], 3));
