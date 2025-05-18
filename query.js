

var travelReq = encodeURIComponent(JSON.stringify({
    Name: "1"
}))


//console.log(travelReq)


function computeSum(arr) {

    if (arr.length === 0) return 0;
    return arr[0] + computeSum(arr.slice(1));

}

console.log(computeSum([1,2,4,5,6]));




function chunk(arr,size)
{ 
  const result= [];
 // loop through array with increamting of size
  for (let i=0;i<arr.length;i+=size)
  {
   // extract chunk from current inedex to index+size
    const chunk = arr.slice(i,i+size);
  // add chank to results array 
    result.push(chunk);
  }
return result;
}

console.log(chunk([1,2,3,4,5,6], 3));
