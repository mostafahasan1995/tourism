




function computeSum(arr){
    //When the array is empty add 0 to end the recursion
    if(arr.length==0){
        return 0
    }
    //Add the first element of the array plus the sum of the reminder of the array
    return arr[0]+computeSum(arr.slice(1))
}

console.log(computeSum([1, 2, 3, 4, 5, 10]))


function chunk(arr,size){
    let chunks=[]
    let currentChunk=[]
    for(let i=0;i<arr.length;i++){
        n=arr[i]
        currentChunk.push(n)
        if(currentChunk.length==size){ //Add chunks as they reach size
            chunks.push(currentChunk)
            currentChunk=[]
        }else if(i==arr.length-1){ //If last chunk is lesser than size add it
            chunks.push(currentChunk)
        }
    }
    return chunks
}

console.log(chunk([1, 2, 3, 4, 5, 6, 7, 8, 9], 1))