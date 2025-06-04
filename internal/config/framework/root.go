package framework

type Framework string

const (
	Bitrix  Framework = "bitrix"
	Laravel Framework = "laravel"
	Symfony Framework = "symfony"
	Vanilla Framework = "vanilla"
)

func GetAll() []Framework {
	return []Framework{
		Bitrix,
		Laravel,
		Symfony,
		Vanilla,
	}
}

func GetAllStrings() []string {
	all := GetAll()
	strs := make([]string, len(all))
	for i, f := range all {
		strs[i] = f.String()
	}
	return strs
}

func ParseFramework(s string) (Framework) {
	switch s {
	case Bitrix.String():
		return Bitrix
	case Laravel.String():
		return Laravel
	case Symfony.String():
		return Symfony
	case Vanilla.String():
		return Vanilla
	default:
		panic("Неизвестный фреймворк")
	}
}

func (f Framework) String() string {
	return string(f)
}
