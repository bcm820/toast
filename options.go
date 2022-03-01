package toast

type Option func(*File)

func WithPackageName(packageName string) Option {
	return func(f *File) {
		f.pkgName = packageName
	}
}

func WithCUEPackageName(packageName string) Option {
	return func(f *File) {
		f.cuePkgName = packageName
	}
}

func WithTransform(t Transform) Option {
	return func(f *File) {
		switch tt := t.(type) {
		case *ExcludeImport:
			f.eximports = append(f.eximports, tt)
		case *ModifyImport:
			f.modimports = append(f.modimports, tt)
		case *CopyIntoStruct:
			f.copies = append(f.copies, tt)
		case *GenEnumTypeTransform:
			f.genEnumTrans = append(f.genEnumTrans, tt)
		default:
			f.trans = append(f.trans, tt)
		}
	}
}
