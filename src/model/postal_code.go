package model

type PostalCode struct {
	ID         int    `json:"id"`
	PostalCode string `json:"postal_code"`
	Province   string `json:"province"`
	Regency    string `json:"regency"`
	District   string `json:"district"`
	Village    string `json:"village"`
}

type ListPostalCode []*PostalCode

func (l ListPostalCode) Provinces() []string {
	result := []string{}

	mapProvince := map[string]struct{}{}
	for _, v := range l {
		mapProvince[v.Province] = struct{}{}
	}

	for k := range mapProvince {
		result = append(result, k)
	}

	return result
}

func (l ListPostalCode) Regencies(province string) []*Regency {
	result := []*Regency{}
	mapProvinceRegency := map[string]map[string]struct{}{}

	setRegency := func(province, regency string) {
		if mapProvinceRegency[province] == nil {
			mapProvinceRegency[province] = map[string]struct{}{}
		}
		mapProvinceRegency[province][regency] = struct{}{}
	}

	for _, data := range l {
		if province == "" {
			setRegency(data.Province, data.Regency)
		}
		if data.Province == province {
			setRegency(data.Province, data.Regency)
		}
	}

	for key, value := range mapProvinceRegency {
		r := Regency{
			Province: key,
		}

		for k := range value {
			r.Regency = k
			result = append(result, &r)
		}
	}

	return result
}

func (l ListPostalCode) District(province, regency string) []*District {
	result := []*District{}
	provinceRegencyDistrict := map[string]map[string]map[string]struct{}{}

	setDistrict := func(province, regency, district string) {
		if provinceRegencyDistrict[province] == nil {
			provinceRegencyDistrict[province] = map[string]map[string]struct{}{}
		}
		if provinceRegencyDistrict[province][regency] == nil {
			provinceRegencyDistrict[province][regency] = map[string]struct{}{}
		}
		provinceRegencyDistrict[province][regency][district] = struct{}{}
	}

	for _, data := range l {
		if province == "" {
			setDistrict(data.Province, data.Regency, data.District)
		}
		if data.Province == province {
			setDistrict(data.Province, data.Regency, data.District)
		}
	}

	for province, regencies := range provinceRegencyDistrict {
		d := District{
			Province: province,
		}

		for regency, districts := range regencies {
			d.Regency = regency

			for district := range districts {
				d.District = district
				result = append(result, &d)
			}
		}
	}

	return result
}
