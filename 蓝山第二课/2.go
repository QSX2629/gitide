package main

func CalculationFactory(operation string) func(int, int) int {
	switch operation {
	case "add":
		return func(a, b int) int {
			return a + b
		}
	case "subtract":
		return func(a, b int) int {
			return a - b
		}
	case "multiply":
		return func(a, b int) int {
			return a * b
		}
	case "divide":
		return func(a, b int) int {
			if b != 0 {
				return a / b
			}
			return 0
		}
	default:
		return nil
	}
}
func main() {

}
