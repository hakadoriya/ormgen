package iface

type Iface interface {
	FilepathAbs(path string) (string, error)
}

type Pkg struct {
	FilepathAbsFunc func(path string) (string, error)
}

func (p *Pkg) FilepathAbs(path string) (string, error) {
	return p.FilepathAbsFunc(path)
}
