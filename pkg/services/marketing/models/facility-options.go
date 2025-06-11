package models

// GetDefaultFacilityOptions returns the default set of facility options for a new exhibitor profile
func GetDefaultFacilityOptions() FacilitiesSection {
	return FacilitiesSection{
		Restaurants: []FacilityOption{
			{Id: "rest_1", Label: "Fine Dining Restaurant", Type: "checkbox", Category: "restaurants", IsSelected: false},
			{Id: "rest_2", Label: "Casual Dining", Type: "checkbox", Category: "restaurants", IsSelected: false},
			{Id: "rest_3", Label: "Buffet Restaurant", Type: "checkbox", Category: "restaurants", IsSelected: false},
			{Id: "rest_4", Label: "Specialty Restaurant", Type: "checkbox", Category: "restaurants", IsSelected: false},
			{Id: "rest_5", Label: "Breakfast Included", Type: "checkbox", Category: "restaurants", IsSelected: false},
		},
		PoolsBeaches: []FacilityOption{
			{Id: "pool_1", Label: "Outdoor Pool", Type: "checkbox", Category: "poolsBeaches", IsSelected: false},
			{Id: "pool_2", Label: "Indoor Pool", Type: "checkbox", Category: "poolsBeaches", IsSelected: false},
			{Id: "pool_3", Label: "Private Beach Access", Type: "checkbox", Category: "poolsBeaches", IsSelected: false},
			{Id: "pool_4", Label: "Kids Pool", Type: "checkbox", Category: "poolsBeaches", IsSelected: false},
			{Id: "pool_5", Label: "Adults-Only Pool", Type: "checkbox", Category: "poolsBeaches", IsSelected: false},
		},
		SpaGym: []FacilityOption{
			{Id: "spa_1", Label: "Full-Service Spa", Type: "checkbox", Category: "spaGym", IsSelected: false},
			{Id: "spa_2", Label: "Fitness Center", Type: "checkbox", Category: "spaGym", IsSelected: false},
			{Id: "spa_3", Label: "Yoga Classes", Type: "checkbox", Category: "spaGym", IsSelected: false},
			{Id: "spa_4", Label: "Massage Services", Type: "checkbox", Category: "spaGym", IsSelected: false},
			{Id: "spa_5", Label: "Sauna/Steam Room", Type: "checkbox", Category: "spaGym", IsSelected: false},
		},
		HotelServices: []FacilityOption{
			{Id: "serv_1", Label: "24-Hour Front Desk", Type: "checkbox", Category: "hotelServices", IsSelected: false},
			{Id: "serv_2", Label: "Concierge Service", Type: "checkbox", Category: "hotelServices", IsSelected: false},
			{Id: "serv_3", Label: "Room Service", Type: "checkbox", Category: "hotelServices", IsSelected: false},
			{Id: "serv_4", Label: "Airport Shuttle", Type: "checkbox", Category: "hotelServices", IsSelected: false},
			{Id: "serv_5", Label: "Laundry Service", Type: "checkbox", Category: "hotelServices", IsSelected: false},
			{Id: "serv_6", Label: "Free Parking", Type: "checkbox", Category: "hotelServices", IsSelected: false},
			{Id: "serv_7", Label: "Free WiFi", Type: "checkbox", Category: "hotelServices", IsSelected: false},
		},
		Business: []FacilityOption{
			{Id: "biz_1", Label: "Business Center", Type: "checkbox", Category: "business", IsSelected: false},
			{Id: "biz_2", Label: "Meeting Rooms", Type: "checkbox", Category: "business", IsSelected: false},
			{Id: "biz_3", Label: "Conference Facilities", Type: "checkbox", Category: "business", IsSelected: false},
			{Id: "biz_4", Label: "Event Space", Type: "checkbox", Category: "business", IsSelected: false},
		},
		KidsFacilities: []FacilityOption{
			{Id: "kids_1", Label: "Kids Club", Type: "checkbox", Category: "kidsFacilities", IsSelected: false},
			{Id: "kids_2", Label: "Babysitting/Childcare", Type: "checkbox", Category: "kidsFacilities", IsSelected: false},
			{Id: "kids_3", Label: "Children's Playground", Type: "checkbox", Category: "kidsFacilities", IsSelected: false},
			{Id: "kids_4", Label: "Family Rooms", Type: "checkbox", Category: "kidsFacilities", IsSelected: false},
		},
		Recreational: []FacilityOption{
			{Id: "rec_1", Label: "Tennis Courts", Type: "checkbox", Category: "recreational", IsSelected: false},
			{Id: "rec_2", Label: "Golf Course", Type: "checkbox", Category: "recreational", IsSelected: false},
			{Id: "rec_3", Label: "Water Sports", Type: "checkbox", Category: "recreational", IsSelected: false},
			{Id: "rec_4", Label: "Bicycle Rental", Type: "checkbox", Category: "recreational", IsSelected: false},
			{Id: "rec_5", Label: "Entertainment Programs", Type: "checkbox", Category: "recreational", IsSelected: false},
		},
	}
}
