package akuntansi

type GetAllJU struct {
	StartDate string `query:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `query:"end_date" validate:"required,datetime=2006-01-02"`
}

type GetAllBB struct {
	AkunID    string `query:"akun_id" validate:"omitempty"`
	StartDate string `query:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `query:"end_date" validate:"required,datetime=2006-01-02"`
}

type GetAllNC struct {
	Date string `query:"date" validate:"required,datetime=2006-01"`
}
