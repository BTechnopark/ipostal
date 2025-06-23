package model

import "strings"

type PostalCode struct {
	ID         int    `json:"id"`
	PostalCode string `json:"postal_code"`
	Province   string `json:"province"`
	Regency    string `json:"regency"`
	District   string `json:"district"`
	Village    string `json:"village"`
}

type ListPostalCode []*PostalCode

func (l ListPostalCode) Provinces(p string) []*Province {
	result := []*Province{}

	mapProvince := map[string]struct{}{}
	for _, v := range l {
		if !strings.Contains(strings.ToLower(v.Province), strings.ToLower(p)) {
			continue
		}

		mapProvince[v.Province] = struct{}{}
	}

	for k := range mapProvince {
		result = append(result, &Province{
			Province: k,
		})
	}

	return result
}

func (l ListPostalCode) Regencies(p, r string) []*Regency {
	result := []*Regency{}
	mapProvinceRegency := map[string]map[string]struct{}{}

	setRegency := func(p, regency string) {
		if mapProvinceRegency[p] == nil {
			mapProvinceRegency[p] = map[string]struct{}{}
		}
		mapProvinceRegency[p][regency] = struct{}{}
	}

	for _, data := range l {
		if !strings.Contains(strings.ToLower(data.Province), strings.ToLower(p)) {
			continue
		}
		if !strings.Contains(strings.ToLower(data.Regency), strings.ToLower(r)) {
			continue
		}

		setRegency(data.Province, data.Regency)
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

func (l ListPostalCode) District(p, r, d string) []*District {
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
		if !strings.Contains(strings.ToLower(data.Province), strings.ToLower(p)) {
			continue
		}
		if !strings.Contains(strings.ToLower(data.Regency), strings.ToLower(r)) {
			continue
		}
		if !strings.Contains(strings.ToLower(data.District), strings.ToLower(d)) {
			continue
		}
		setDistrict(data.Province, data.Regency, data.District)
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
