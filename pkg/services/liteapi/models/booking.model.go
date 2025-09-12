package models

type Booking map[string]any

// {
// 	"addons": null,
// 	"addonsRedemptions": null,
// 	"addonsTotalAmount": 0,
// 	"adults": 1,
// 	"agentId": null,
// 	"amountRefunded": 0,
// 	"apiCommission": 12.28,
// 	"bookedRooms": [
// 		{
// 			"adults": 1,
// 			"amount": 204.71,
// 			"board": "Room Only",
// 			"boardCode": "RO",
// 			"boardName": "Room Only",
// 			"boardType": "RO",
// 			"cancellationPolicies": {
// 				"cancelPolicyInfos": [
// 					{
// 						"amount": 204.71,
// 						"cancelTime": "2025-11-21 23:59:59",
// 						"currency": "USD",
// 						"timezone": "GMT",
// 						"type": "amount"
// 					}
// 				],
// 				"hotelRemarks": null,
// 				"refundableTag": "RFN"
// 			},
// 			"children": 0,
// 			"childrenAges": null,
// 			"children_count": 0,
// 			"currency": "USD",
// 			"firstName": "Sunny",
// 			"guests": [
// 				{
// 					"email": "s.mars@liteapi.travel",
// 					"firstName": "Sunny",
// 					"lastName": "Mars",
// 					"occupancyNumber": 1,
// 					"phone": "",
// 					"remarks": "quiet room please"
// 				}
// 			],
// 			"lastName": "Mars",
// 			"occupancy_number": 1,
// 			"rate": {
// 				"boardName": "RO",
// 				"boardType": "Room Only",
// 				"cancellationPolicies": {
// 					"cancelPolicyInfos": [
// 						{
// 							"amount": 204.71,
// 							"cancelTime": "2025-11-21 23:59:59",
// 							"currency": "USD",
// 							"timezone": "GMT",
// 							"type": "amount"
// 						}
// 					],
// 					"hotelRemarks": null,
// 					"refundableTag": "RFN"
// 				},
// 				"maxOccupancy": 1,
// 				"rateId": "",
// 				"remarks": "quiet room please",
// 				"retailRate": {
// 					"suggestedSellingPrice": {
// 						"source": "providerDirect"
// 					},
// 					"total": {
// 						"amount": 216.99,
// 						"currency": "USD"
// 					}
// 				}
// 			},
// 			"remarks": "quiet room please",
// 			"roomType": {
// 				"name": "Test rate 1",
// 				"roomTypeId": ""
// 			},
// 			"room_id": "GY2DMNJWMM3TKNZYGY2TEMBWGQ3GMNZVGYZDMYZWGUZDANZSGZTDMZRWMQZDANZXGY4TONBWHAZDANRTGY4TONBXHEZDANZWGY4TMNJXG4ZDAMRYGY3DONJWMM3GGMRQGY2DMZRXGU3DENTDGY2TEMBWGI3DKNRUGI4TEMBSHA3DENRVGY2DEMBXGQ3TSNZQGY2TEMBWHE3TGMRQG4ZTONJWGI3GCNRVGYZTONBSGA3TINTGGIYDMMJXGY3DCNRZGZRTMMJWGI3DSNTDGY4TONBXHEZDS7BREMZDAMRVGEYTEML4GIYDENJRGEZDE7DFNZPVKU34KVJXYVKTIR6DCQJQIN6DEMRYGN6FKMSWLBZHYMJXGU3TKNZXGI3DEOJQGVHWK22EEMYTANJYGA2CGMJTGU2C2MJQGU4DAND4KJDE47BSGAZDKMJRGA2DCNRQGB6DC7CSGEYTKMD4GIYDINZRENJE6I2SIZHCGMRQGI2S2MJRFUYDIIBRGY5DAMBD"
// 		}
// 	],
// 	"bookingId": "4oJRtrss3",
// 	"cancellationPolicies": {
// 		"cancelPolicyInfos": [
// 			{
// 				"amount": 204.71,
// 				"cancelTime": "2025-11-21 23:59:59",
// 				"currency": "USD",
// 				"timezone": "GMT",
// 				"type": "amount"
// 			}
// 		],
// 		"hotelRemarks": null,
// 		"refundableTag": "RFN"
// 	},
// 	"cancelledAt": null,
// 	"cancelledBy": null,
// 	"checkin": "2025-11-21",
// 	"checkout": "2025-11-22",
// 	"children": "",
// 	"childrenCount": 0,
// 	"clientCommission": 12.28,
// 	"clientReference": "",
// 	"commission": 12.28,
// 	"createdAt": "2025-09-11T08:06:19",
// 	"currency": "USD",
// 	"distributorCommission": 0,
// 	"distributorPrice": 0,
// 	"email": "fmnsha@gmail.com",
// 	"exchangeRate": 0.8538393725172431,
// 	"exchangeRateUsd": 1,
// 	"firstName": "Feras",
// 	"guestId": 0,
// 	"holder": {
// 		"email": "fmnsha@gmail.com",
// 		"firstName": "Feras",
// 		"lastName": "mnsha",
// 		"phone": ""
// 	},
// 	"holderTitle": "",
// 	"hotel": {
// 		"hotelId": "lp19d4c",
// 		"name": ""
// 	},
// 	"hotelConfirmationCode": "test",
// 	"hotelId": "lp19d4c",
// 	"hotelName": "",
// 	"knowBeforeYouGo": "Know before you go",
// 	"lastFreeCancellationDate": "2025-11-21T23:59:59Z",
// 	"lastName": "mnsha",
// 	"loyaltyGuestId": null,
// 	"mandatoryFees": "Mandatory fees",
// 	"nationality": "US",
// 	"optionalFees": "Optional fees",
// 	"paymentScheduledAt": null,
// 	"paymentStatus": "succeeded",
// 	"paymentTransactionId": "tr_cts_KP5DOYI80LcGWlJvIhW37",
// 	"prebookId": "4h8qPpHfR",
// 	"price": 216.99,
// 	"processingFee": 8.67,
// 	"rebookFrom": "",
// 	"refundType": "",
// 	"refundedAt": null,
// 	"remarks": "Remarks sandbox",
// 	"sandbox": 1,
// 	"sellingPrice": "216.99",
// 	"specialRemarks": "Example special remarks",
// 	"status": "CONFIRMED",
// 	"supplier": "nuitee",
// 	"supplierBookingId": "4oJRtrss3",
// 	"supplierBookingName": "nuitee",
// 	"supplierId": 2,
// 	"tag": "RFN",
// 	"trackingId": "",
// 	"updatedAt": "",
// 	"userId": 370021,
// 	"voucherCode": "",
// 	"voucherId": null,
// 	"voucherTotalAmount": 0,
// 	"voucherTransationId": null
// }
