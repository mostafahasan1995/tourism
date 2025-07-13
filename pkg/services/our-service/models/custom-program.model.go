package models

import (
	"larsa-tourism-microservices/pkg/transl"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomProgram struct {
	Delegation          Delegation                 `bson:"delegation,omitempty" json:"delegation,omitempty"`
	BusinessMan         BusinessMan                `bson:"businessMan,omitempty" json:"businessMan,omitempty"`
	CustomPlan          CustomPlan                 `bson:"customPlan,omitempty" json:"customPlan,omitempty"`
	HotelBooking        HotelBooking               `bson:"hotelBooking,omitempty" json:"hotelBooking,omitempty"`
	VipCar              ProgramVipCar              `bson:"vipCar,omitempty" json:"vipCar,omitempty"`
	FlightTicketRequest ProgramFlightTicketRequest `bson:"flightTicketRequest,omitempty" json:"flightTicketRequest,omitempty"`
	PartnerRequest      PartnerRequest             `bson:"partnerRequest,omitempty" json:"partnerRequest,omitempty"`
	Destinations        []ProgramDestination       `bson:"destinations,omitempty" json:"destinations,omitempty"`
}

type ProgramVipCar struct {
	Destinations []ProgramVipCarDestination `bson:"destinations" json:"destinations"`
}

type ProgramVipCarDestination struct {
	DestinationFrom       primitive.ObjectID `bson:"destinationFrom" json:"destinationFrom"`
	DestinationTo         primitive.ObjectID `bson:"destinationTo" json:"destinationTo"`
	ProgramTransportation `bson:",inline"`
	Services              ProgramServices `bson:"services" json:"services"`
}

type ProgramFlightTicketRequest struct {
	Destinations []ProgramFlightTicket `bson:"destinations" json:"destinations"`
}

type ProgramDestination struct {
	DestinationFrom primitive.ObjectID      `bson:"destinationFrom" json:"destinationFrom"`
	DestinationTo   primitive.ObjectID      `bson:"destinationTo" json:"destinationTo"`
	TripDetails     TripDetails             `bson:"tripDetails" json:"tripDetails"`
	Accommodation   []ProgramAccommodation  `bson:"accommodation" json:"accommodation"`
	FlightTickets   DestinationFlightTicket `bson:"flightTickets" json:"flightTickets"`
	Transportation  ProgramTransportation   `bson:"transportation" json:"transportation"`
	Activities      ProgramActivities       `bson:"activities" json:"activities"`
	Agenda          Agenda                  `bson:"agenda" json:"agenda"`
	Services        ProgramServices         `bson:"services" json:"services"`
}

func (pd *ProgramDestination) GetProgramDestServicePricing() ([]InvoiceService, error) {
	var services []InvoiceService

	for _, ac := range pd.Accommodation {
		service := InvoiceService{
			Item:  ac.Preferences,
			Price: ac.TotalStayCost,
			Qty:   1,
		}
		services = append(services, service)
	}

	if pd.FlightTickets.TotalCost > 0 {
		services = append(services, InvoiceService{
			Item:  pd.FlightTickets.TripType,
			Price: pd.FlightTickets.TotalCost,
			Qty:   1,
		})
	}

	if pd.Transportation.TotalCost > 0 {
		services = append(services, InvoiceService{
			Item:  pd.Transportation.TransType,
			Price: pd.Transportation.TotalCost,
			Qty:   1,
		})
	}

	if pd.Activities.TotalCost > 0 {
		services = append(services, InvoiceService{
			Item:  "Activities",
			Price: pd.Activities.TotalCost,
			Qty:   1,
		})
	}

	if pd.Services.TourGuide.Active {
		services = append(services, InvoiceService{
			Item:  "Tour Guide",
			Price: pd.Services.TourGuide.Cost,
			Qty:   1,
		})
	}

	if pd.Services.Translator.Active {
		services = append(services, InvoiceService{
			Item:  "Translator",
			Price: pd.Services.Translator.Cost,
			Qty:   1,
		})
	}

	if pd.Services.AirportPickup.Active {
		services = append(services, InvoiceService{
			Item:  "Airport Pickup",
			Price: pd.Services.AirportPickup.Cost,
			Qty:   1,
		})
	}

	if pd.Services.TourAfterMeeting.Active {
		services = append(services, InvoiceService{
			Item:  "Tour After Meeting",
			Price: pd.Services.TourAfterMeeting.Cost,
			Qty:   1,
		})
	}

	if pd.Services.Photography.Active {
		services = append(services, InvoiceService{
			Item:  "Photography",
			Price: pd.Services.Photography.Cost,
		})
	}

	if pd.Services.AirportMeetAndGreet.Active {
		services = append(services, InvoiceService{
			Item:  "Airport Meet and Greet",
			Price: pd.Services.AirportMeetAndGreet.Cost,
		})
	}

	if pd.Services.SimCardAndInternet.Active {
		services = append(services, InvoiceService{
			Item:  "Sim Card and Internet",
			Price: pd.Services.SimCardAndInternet.Cost,
		})
	}
	return services, nil
}

type ProgramAccommodation struct {
	Accommodation      `bson:",inline"`
	AccommodationHotel transl.Localizable[string] `bson:"accommodationHotel" json:"accommodationHotel"`
	StartDate          time.Time                  `bson:"startDate" json:"startDate"`
	EndDate            time.Time                  `bson:"endDate" json:"endDate"`
	PricePerNight      float64                    `bson:"pricePerNight" json:"pricePerNight"`
	TotalStayCost      float64                    `bson:"totalStayCost" json:"totalStayCost"`
}

type ProgramFlightTicket struct {
	FlightTicket `bson:",inline"`
	TotalCost    float64         `bson:"totalCost" json:"totalCost"`
	Services     ProgramServices `bson:"services" json:"services"`
}

type DestinationFlightTicket struct {
	FlightTicket `bson:",inline"`
	TotalCost    float64 `bson:"totalCost" json:"totalCost"`
}

type ProgramTransportation struct {
	Transportation `bson:",inline"`
	TotalCost      float64 `bson:"totalCost" json:"totalCost"`
}

type ProgramActivities struct {
	Activities []transl.Localizable[string] `bson:"activities" json:"activities"`
	TotalCost  float64                      `bson:"totalCost" json:"totalCost"`
}

type Service struct {
	Active bool    `bson:"active" json:"active"`
	Cost   float64 `bson:"cost" json:"cost"`
}

type ProgramServices struct {
	TourGuide           Service `bson:"tourGuide" json:"tourGuide"`
	Translator          Service `bson:"translator" json:"translator"`
	AirportPickup       Service `bson:"airportPickup" json:"airportPickup"`
	TourAfterMeeting    Service `bson:"tourAfterMeeting" json:"tourAfterMeeting"`
	Photography         Service `bson:"photography" json:"photography"`
	AirportMeetAndGreet Service `bson:"airportMeetAndGreet" json:"airportMeetAndGreet"`
	SimCardAndInternet  Service `bson:"simCardAndInternet" json:"simCardAndInternet"`
}

func (cp *CustomProgram) GetOtherServicePricing() ([]InvoiceService, error) {
	var services []InvoiceService

	for _, des := range cp.VipCar.Destinations {
		service := InvoiceService{
			Item:  des.TransType,
			Price: des.TotalCost,
			Qty:   1,
		}
		services = append(services, service)

		// Convert services to invoice services
		destServices, err := convertServicesToInvoiceServices(des.Services)
		if err != nil {
			return nil, err
		}
		services = append(services, destServices...)
	}

	for _, des := range cp.FlightTicketRequest.Destinations {
		service := InvoiceService{
			Item:  des.TripType,
			Price: des.TotalCost,
			Qty:   1,
		}
		services = append(services, service)
		// Convert services to invoice services
		destServices, err := convertServicesToInvoiceServices(des.Services)
		if err != nil {
			return nil, err
		}
		services = append(services, destServices...)
	}

	return services, nil
}

func (cp *CustomProgram) GetAllServicePricing() ([]InvoiceService, error) {

	var services []InvoiceService

	otherServices, err := cp.GetOtherServicePricing()
	if err != nil {
		return nil, err
	}

	services = append(services, otherServices...)

	for _, des := range cp.Destinations {
		destServices, err := des.GetProgramDestServicePricing()
		if err != nil {
			return nil, err
		}
		services = append(services, destServices...)
	}

	return services, nil
}

func (cp *CustomProgram) GetTotalPrice() (float64, error) {
	services, err := cp.GetAllServicePricing()
	if err != nil {
		return 0, err
	}

	var total float64
	for _, service := range services {
		total += service.Price * float64(service.Qty)
	}
	return total, nil
}

// convertServicesToInvoiceServices converts service configurations to invoice services
func convertServicesToInvoiceServices(servicesConfig ProgramServices) ([]InvoiceService, error) {
	serviceConfigs, err := util.StructToMap(servicesConfig)
	if err != nil {
		return nil, err
	}

	var services []InvoiceService
	for key, value := range serviceConfigs {
		if value.(Service).Active {
			services = append(services, InvoiceService{
				Item:  key,
				Price: value.(Service).Cost,
				Qty:   1,
			})
		}
	}

	return services, nil
}
