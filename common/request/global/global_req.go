package global

type ParamByID struct {
	ID string `params:"id" validate:"required,ulid"`
}