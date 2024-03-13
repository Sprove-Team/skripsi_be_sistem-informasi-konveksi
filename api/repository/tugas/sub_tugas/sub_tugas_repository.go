package repo_sub_tugas

type SubTugas interface {
}

type subTugas struct {
}

func NewRepoSubTugas() SubTugas {
	return &subTugas{}
}
