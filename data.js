

const data = {
    table:[
        {
            id:"",
            isCore: true,
            name:"id",
            label:"id",
            type:"text-email-number....etc"
        },
        {
            id:"",
            isCore: true,
            name:"name",
            label:"name",
            type:"text-email-number....etc"
        },
        {
            id:"",
            isCore: true,
            name:"age",
            label:"age",
            type:"text-email-number....etc"
        },
        {
            id:"",
            isCore: true,
            name:"status",
            label:"status",
            type:"text-email-number....etc"
        }
    ],
    data :[
        {
            id:"g1",
            color:"color",
            items:[
                {
                    id:"1",
                    name:"feras",
                    age:40,
                    status:"pending"
                }, //each object is a row
                {
                    id:"2",
                    name:"adham",
                    age:47,
                    status:"approved"
                },
                {
                    id:"3",
                    name:"anas",
                    age:25,
                    status:"pending"
                }

            ],
        },
        {
            id:"g2",
            color:"color",
            data:[
                {
                    id:"4",
                    name:"yard",
                    age:24,
                    status:"pending"
                },
                {
                    id:"5",
                    name:"mona",
                    age:30,
                    status:"pending"
                },
                {
                    id:"6",
                    name:"sara",
                    age:28,
                    status:"approved"
                }

            ],
            
        }
    ]
}



function test(){
    var columns = data.table 

    //group -1
    var group1 = data.data[0]
    var group1Data = group1.data

    for (var i=0;i<group1Data.length;i++){
        var o = group1Data[i]
        for (var j = 0; j< columns.length;j++){
            var col = columns[j]
            var name = col.name
            console.log("column name",name)
            console.log("column data",o[name])

        }
    }
    
}

test()


var board = {
    tables :[
        {
            columns:[
                {
                    
                },
                {},
                {},
                {},
                {},

            ],
            groups:[
                {
                    groupId:"123f",
                    groupName:"name",
                    Items :[
                        {},
                        {},
                        {},
                        {},
                    ]
                },
                {
                    groupId:"123f",
                    groupName:"name",
                    Items :[
                        {},
                        {},
                        {},
                        {},
                    ]
                },
                {
                    groupId:"123f",
                    groupName:"name",
                    Items :[
                        {},
                        {},
                        {},
                        {},
                    ]
                }
            ]
        }
    ],
    forms :[],
    calenders:[]
}

var kan = [
    {
        "_id":"inprogress",
        "items":[]
    },
    {
        "_id":"completed",
        "items":[]
    },
    {
        "_id":"blank",
        "items":[]
    },
    {
        "_id":"stuck",
        "items":[]
    }

]