package application

type CalculationMode string

const (
	ModeStrict        CalculationMode = "strict"
	ModeMissingAsZero CalculationMode = "missing_as_zero"
)

func (m CalculationMode) Valid() bool {
	return m == ModeStrict || m == ModeMissingAsZero
}
