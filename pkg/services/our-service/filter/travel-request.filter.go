package filter

import (
	"fmt"
	"regexp"
	"time"

	"larsa-tourism-microservices/pkg/services/our-service/enums"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DelegationFilter struct {
	DelegationType *string `json:"delegationType"`
}

type BusinessManFilter struct {
	Purpose *string `json:"purpose"`
}

type VipCarDestinationFilter struct {
	From *primitive.ObjectID `json:"from"`
	To   *primitive.ObjectID `json:"to"`
}

type VipCarFilter struct {
	Destination *VipCarDestinationFilter `json:"destination"`
}

type FlightTicketFilter struct {
	From          *primitive.ObjectID `json:"from"`
	To            *primitive.ObjectID `json:"to"`
	TripType      *string             `json:"tripType"`
	TravelClass   *string             `json:"travelClass"`
	DepartureDate *time.Time          `json:"departureDate"`
}

type TravelReqFilters struct {
	CustomerName *string             `json:"customerName"`
	Status       *string             `json:"status"`
	CustomerId   *primitive.ObjectID `json:"customerId"`
	ProgramId    *primitive.ObjectID `json:"programId"`
	ProgramName  *string             `json:"programName"`
	Date         *time.Time          `json:"date"`
	ServiceType  []enums.ServiceType `json:"serviceType"`
	Email        *string             `json:"email"`
	Phone        *string             `json:"phone"`
	Duration     []int               `json:"duration"` // Changed to array
	TripType     *string             `json:"tripType"` // for customPlan.tripType
	Delegation   *DelegationFilter   `json:"delegation"`
	BusinessMan  *BusinessManFilter  `json:"businessMan"`
	VipCar       *VipCarFilter       `json:"vipCar"`
	FlightTicket *FlightTicketFilter `json:"flightTicket"`
}

func (f TravelReqFilters) BuildPipeline(m bson.M) []bson.M {
	var ands bson.A
	if f.CustomerName != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.CustomerName)
		re, _ := regexp.Compile(pattern)
		nameFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$customerName",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}

		ands = append(ands, nameFilter)
	}

	if f.Status != nil {
		m["status"] = *f.Status
	}

	if f.CustomerId != nil {
		m["customerId"] = *f.CustomerId
	}

	if f.ProgramId != nil {
		if f.ProgramId.IsZero() {
			o := bson.M{
				"$or": bson.A{
					bson.M{"program": bson.M{"$exists": false}},
					bson.M{"program": primitive.NilObjectID},
				},
			}
			ands = append(ands, o)
		} else {
			m["program"] = *f.ProgramId
		}
	}

	// New filter fields
	if f.ProgramName != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.ProgramName)
		re, _ := regexp.Compile(pattern)
		programNameFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$programData.title",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}
		ands = append(ands, programNameFilter)
	}

	if f.Date != nil {
		m["date"] = *f.Date
	}

	if len(f.ServiceType) > 0 {
		serviceTypes := make([]string, len(f.ServiceType))
		for i, st := range f.ServiceType {
			serviceTypes[i] = string(st)
		}
		m["serviceType"] = bson.M{"$in": serviceTypes}
	}

	if f.Email != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.Email)
		re, _ := regexp.Compile(pattern)
		emailFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$clientEmail",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}
		ands = append(ands, emailFilter)
	}

	if f.Phone != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.Phone)
		re, _ := regexp.Compile(pattern)
		phoneFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$clientPhone.number",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}
		ands = append(ands, phoneFilter)
	}

	if len(f.Duration) > 0 {
		// Check if -1 is in the array (means 15 days and above)
		hasNegativeOne := false
		validDurations := []int{}
		for _, duration := range f.Duration {
			if duration == -1 {
				hasNegativeOne = true
			} else {
				validDurations = append(validDurations, duration)
			}
		}

		if hasNegativeOne && len(validDurations) > 0 {
			// Both -1 and specific durations: use $or
			durationFilter := bson.M{
				"$or": bson.A{
					bson.M{"tripDuration": bson.M{"$gte": 15}},
					bson.M{"tripDuration": bson.M{"$in": validDurations}},
				},
			}
			ands = append(ands, durationFilter)
		} else if hasNegativeOne {
			// Only -1: filter for 15 days and above
			m["tripDuration"] = bson.M{"$gte": 15}
		} else {
			// Only specific durations: use $in
			m["tripDuration"] = bson.M{"$in": validDurations}
		}
	}

	if f.TripType != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.TripType)
		re, _ := regexp.Compile(pattern)
		tripTypeFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$customPlan.tripType",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}
		ands = append(ands, tripTypeFilter)
	}

	// Delegation filter
	if f.Delegation != nil && f.Delegation.DelegationType != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.Delegation.DelegationType)
		re, _ := regexp.Compile(pattern)
		delegationFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$delegation.delegationType",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}
		ands = append(ands, delegationFilter)
	}

	// BusinessMan filter
	if f.BusinessMan != nil && f.BusinessMan.Purpose != nil {
		pattern := fmt.Sprintf(".*%s.*", *f.BusinessMan.Purpose)
		re, _ := regexp.Compile(pattern)
		businessManFilter := bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$businessMan.purpose",
					"regex":   re.String(),
					"options": "i",
				},
			},
		}
		ands = append(ands, businessManFilter)
	}

	// VipCar filter
	if f.VipCar != nil && f.VipCar.Destination != nil {
		vipCarConditions := bson.A{}
		if f.VipCar.Destination.From != nil {
			vipCarConditions = append(vipCarConditions, bson.M{
				"$eq": bson.A{"$$dest.destinationFrom", *f.VipCar.Destination.From},
			})
		}
		if f.VipCar.Destination.To != nil {
			vipCarConditions = append(vipCarConditions, bson.M{
				"$eq": bson.A{"$$dest.destinationTo", *f.VipCar.Destination.To},
			})
		}
		if len(vipCarConditions) > 0 {
			vipCarFilter := bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$ne": bson.A{"$vipCar.destinations", nil}},
						bson.M{"$isArray": "$vipCar.destinations"},
						bson.M{
							"$anyElementTrue": bson.M{
								"$map": bson.M{
									"input": "$vipCar.destinations",
									"as":    "dest",
									"in": bson.M{
										"$and": vipCarConditions,
									},
								},
							},
						},
					},
				},
			}
			ands = append(ands, vipCarFilter)
		}
	}

	// FlightTicket filter
	if f.FlightTicket != nil {
		flightConditions := bson.A{}
		if f.FlightTicket.From != nil {
			flightConditions = append(flightConditions, bson.M{
				"$eq": bson.A{"$$dest.destinationFrom", *f.FlightTicket.From},
			})
		}
		if f.FlightTicket.To != nil {
			flightConditions = append(flightConditions, bson.M{
				"$eq": bson.A{"$$dest.destinationTo", *f.FlightTicket.To},
			})
		}
		if f.FlightTicket.TripType != nil {
			flightConditions = append(flightConditions, bson.M{
				"$eq": bson.A{"$$dest.tripType", *f.FlightTicket.TripType},
			})
		}
		if f.FlightTicket.TravelClass != nil {
			flightConditions = append(flightConditions, bson.M{
				"$eq": bson.A{"$$dest.travelClass", *f.FlightTicket.TravelClass},
			})
		}
		if f.FlightTicket.DepartureDate != nil {
			flightConditions = append(flightConditions, bson.M{
				"$eq": bson.A{"$$dest.departureDate", *f.FlightTicket.DepartureDate},
			})
		}
		if len(flightConditions) > 0 {
			flightFilter := bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$ne": bson.A{"$flightTicketRequest.destinations", nil}},
						bson.M{"$isArray": "$flightTicketRequest.destinations"},
						bson.M{
							"$anyElementTrue": bson.M{
								"$map": bson.M{
									"input": "$flightTicketRequest.destinations",
									"as":    "dest",
									"in": bson.M{
										"$and": flightConditions,
									},
								},
							},
						},
					},
				},
			}
			ands = append(ands, flightFilter)
		}
	}

	if len(ands) > 0 {
		m["$and"] = ands
	}

	pipeline := []bson.M{
		{"$match": m},
	}

	return pipeline
}
