package filters

import "go.mongodb.org/mongo-driver/bson"

type TransTestFilters struct {
	Note    *string  `json:"note"`
	Notes   []string `json:"notes"`
	Project *string  `json:"project"`
}

func (f TransTestFilters) BuildPipeline(m bson.M) []bson.M {

	if f.Note != nil {
		m["notes.en"] = bson.M{"$eq": *f.Note}
	}

	if len(f.Notes) > 0 {
		m["notes.en"] = bson.M{"$in": f.Notes}
	}

	if f.Project != nil {
		m["about.project.project_name.en"] = bson.M{"$eq": *f.Project}
	}

	return []bson.M{
		{"$match": m},
	}
}
