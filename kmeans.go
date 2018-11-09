package noteshrink

type Value interface {
	Distance() float64
}

type Values interface {
	Get() []Value
	Average() Value
}
