








var q = encodeURIComponent(JSON.stringify({
  columns: [
    {
      name: "notes.en",
      exp: "in",
      value: "three,one"
    },
    {
      name: "name.en",
      exp: "like",
      value: "3"
    },
    {
      name: "_id",
      value: "6868dbbac3f6b2ba8e6aca89"
    },
    {
      name: "status",
      value: "active"
    }
  ]
}))

console.log(q)
