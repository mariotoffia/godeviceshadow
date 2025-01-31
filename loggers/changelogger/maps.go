package changelogger

func (ml ManagedLogMap) All() []ManagedValue {
	size := 0

	for _, v := range ml {
		size += len(v)
	}

	all := make([]ManagedValue, 0, size)

	for _, v := range ml {
		all = append(all, v...)
	}

	return all
}

func (ml ManagedLogMap) Size() int {
	size := 0

	for _, v := range ml {
		size += len(v)
	}

	return size
}

func (pl PlainLogMap) All() []PlainValue {
	size := 0

	for _, v := range pl {
		size += len(v)
	}

	all := make([]PlainValue, 0, size)

	for _, v := range pl {
		all = append(all, v...)
	}

	return all
}

func (pl PlainLogMap) Size() int {
	size := 0

	for _, v := range pl {
		size += len(v)
	}

	return size
}
