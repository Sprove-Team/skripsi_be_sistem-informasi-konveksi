package produk

type ParamByID struct {
	ID string `params:"id" validate:"required,uuidv4_no_hyphens"`
}