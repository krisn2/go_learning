package main

import "fmt"

func PrintAllIntForm() {
	var num int = 98
	fmt.Print("\nprinting int in different bases:\n\n")
	fmt.Printf("Binary form: %b\n", num)
	fmt.Printf("Decimal form: %d\n", num)
	fmt.Printf("Octal form :%o\n", num)
	fmt.Printf("Octal with Prefix form: %O\n", num)
	fmt.Printf("single quoted character form: %q\n", num)
	fmt.Printf("Hexa Decimal 10 a-f lower-case form: %x\n", num)
	fmt.Printf("Hexa Decimal 10 A-F upper-case form: %X\n", num)
	fmt.Printf("Unicode format: %U\n", num)
}

func PrintAllFloatForm() {
	num := 123.456
	fmt.Println("\nPrinting Float in Different bases:")
	fmt.Printf("decimalless scinetific notation: %b\n", num)
	fmt.Printf("scinetific notation: %e\n", num)
	fmt.Printf("decimal point : %f\n", num)
	fmt.Printf("decimal point : %F\n", num)

}

func main() {
	// PrintAllIntForm()
	// PrintAllFloatForm()
	name := "ragnar"
	pointer := &name

	count := 0
	n := 5

	if n%2 == 0 {
		fmt.Print("even number \n")
	} else {
		fmt.Print("number is odd\n")
	}

	for i := 0; i <= n; i++ {
		count += i
	}
	fmt.Printf("%d\n", count)

	fmt.Printf("point to name: %s -> %p\n", name, pointer)
	name_change(&name)
	fmt.Printf("point to name: %s -> %p\n", name, pointer)

}

func name_change(name *string) {
	*name = "krisn"
}
