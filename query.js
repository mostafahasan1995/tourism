




const recursion = (arr, index = 0) => {
    return index === arr.length ? 0 : arr[index] + recursion(arr, index + 1);
  };


  console.log(recursion([1, 2, 3, 4, 5, 6,7,8,9,10]));



  const chunk = (arr,index) => {
    const result = [];
    for (let i = 0; i < arr.length; i + index) {
       result.push(arr.slice(i,i+index))
    }
    
}


console.log(chunk([1, 2, 3, 4, 5, 6,7,8,9,10], 2));







