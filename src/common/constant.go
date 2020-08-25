package common

const (
	OnlyChanged = 1 << 0
	OnlyAdded   = 1 << 1
	OnlyDeleted = 1 << 2
)

func GetMethodStatus(onlyChanged bool, onlyAdded bool, onlyDeleted bool) int {
	status := 0
	if onlyChanged {
		status |= OnlyChanged
	}
	if onlyAdded {
		status |= OnlyAdded
	}
	if onlyDeleted {
		status |= OnlyDeleted
	}
	return status
}

func ContainMethodStatus(status int, mask int) bool {
	return status&mask == mask
}

const (
	GroupControl   = 0
	GroupTreatment = 1
)

const (
	SortByDiffTotal = iota
	SortByDiffPartition
	SortByDiffSelf
	SortByDiffCount
	SortByControlTotal
	SortByControlSelf
	SortByControlCount
	SortByTreatmentTotal
	SortByTreatmentSelf
	SortByTreatmentCount
)
