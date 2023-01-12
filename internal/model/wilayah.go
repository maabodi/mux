package model

type WProvinsi struct {
	ProvinsiId int64   `gorm:"primaryKey" json:"-"`
	Provinsi   string  `json:"provinsi"`
	Kota       []WKota `gorm:"foreignKey:ProvinsiPrefer" json:"kota"`
}

type WKota struct {
	KotaId         int64        `gorm:"primaryKey" json:"-"`
	Kota           string       `json:"kota"`
	Kecamatan      []WKecamatan `gorm:"foreignKey:KotaPrefer" json:"kecamatan"`
	ProvinsiPrefer int64        `json:"-"`
}

type WKecamatan struct {
	KecamatanId int64        `gorm:"primaryKey" json:"-"`
	Kecamatan   string       `json:"kecamatan"`
	Kelurahan   []WKelurahan `gorm:"foreignKey:KecamatanPrefer" json:"kelurahan"`
	KotaPrefer  int64        `json:"-"`
}

type WKelurahan struct {
	KelurahanId     int64  `gorm:"primaryKey" json:"-"`
	Kelurahan       string `json:"kelurahan"`
	KecamatanPrefer int64  `json:"-"`
}
