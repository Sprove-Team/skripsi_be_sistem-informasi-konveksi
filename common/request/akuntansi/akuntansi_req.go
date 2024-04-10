package req_akuntansi

type GetAllJU struct {
	StartDate string `query:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `query:"end_date" validate:"required,datetime=2006-01-02"`
	Download  string `query:"download" validate:"omitempty"`
}

type GetAllBB struct {
	AkunID    []string `query:"akun_id" validate:"omitempty,dive,ulid"`
	StartDate string   `query:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string   `query:"end_date" validate:"required,datetime=2006-01-02"`
	Download  string   `query:"download" validate:"omitempty"`
}

type GetAllNC struct {
	Date     string `query:"date" validate:"required,datetime=2006-01"`
	Download string `query:"download" validate:"omitempty"`
}

type GetAllLBR struct {
	StartDate string `query:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `query:"end_date" validate:"required,datetime=2006-01-02"`
	Download  string `query:"download" validate:"omitempty"`
}
