//+build ignore

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Type of a callback argument
type CType string

func (t CType) Ptr() bool {
	return strings.HasSuffix(string(t), "*")
}

func (t CType) Deref() CType {
	if !t.Ptr() {
		return t
	}
	if t.Short() == "voidptr" {
		return "int"
	}
	return CType(strings.Replace(string(t), "*", "", 1))
}

// "unsigned char" -> "unsigned char"
// "int*" -> "int*"
func (t CType) CNotation() string {
	return string(t)
}

// "unsigned char" -> "uchar"
// "int*" -> "intptr"
func (t CType) Short() (s string) {
	const unsigned_ = "unsigned "
	if strings.HasPrefix(string(t), unsigned_) {
		s += "u"
		t = CType([]byte(t)[len(unsigned_):])
	}
	s += strings.Replace(string(t), " ", "", -1)
	s = strings.Replace(s, "*", "ptr", 1)
	return s
}

// "unsigned char" -> "Uchar"
// "int*" -> "Intptr"
func (t CType) CapShort() string {
	bt := []byte(t.Short())
	if bt[0] >= 'a' && bt[0] <= 'z' {
		bt[0] -= 'a' - 'A'
	}
	return string(bt)
}

// "unsigned char" -> "c_uchar"
// "int*" -> "c_intptr"
func (t CType) GoNotation() string {
	return "c_" + t.Short()
}

// "unsigned char" -> "C.uchar"
// "int*" -> "*C.int"
func (t CType) CgoNotation() string {
	if t == "void*" {
		return "unsafe.Pointer"
	}
	s := "C." + t.Deref().Short()
	if t.Ptr() {
		s = "(*" + s + ")"
	}
	return s
}

func init() {
	rand.Seed(time.Now().Unix())
}

// Random value for a type in Go syntax
func (t CType) Random() interface{} {
	switch t.Short() {
	case "char":
		return int8(rand.Intn(math.MaxUint8) + math.MinInt8)
	case "uchar":
		return uint8(rand.Intn(math.MaxUint8))
	case "short":
		return int16(rand.Intn(math.MaxUint16) + math.MinInt16)
	case "ushort":
		return uint16(rand.Intn(math.MaxUint16))
	case "int", "long":
		return int32(rand.Intn(math.MaxInt32))
	case "uint", "ulong":
		return uint32(rand.Intn(math.MaxInt32) + math.MinInt32)
	case "longlong":
		return rand.Int63()
	case "ulonglong":
		return uint64(rand.Int63()) + uint64(rand.Int63())
	case "float":
		return float32(rand.NormFloat64())
	case "double":
		return rand.NormFloat64()
	}
	if t.Ptr() {
		return fmt.Sprintf("new%v(%v)", t.CapShort(), t.Deref().Random())
	}
	return ""
}

// Random value for a type in C syntax
func (t CType) CRandom() interface{} {
	if t.Ptr() {
		return CType("unsigned int").CRandom()
	}
	return t.Random()
}

// List of all supported C types
var Types = []CType{
	"char", "unsigned char",
	"short", "unsigned short",
	"int", "unsigned int",
	"long", "unsigned long",
	"long long", "unsigned long long",
	"float", "double",
	/*"int*",*/ "void*",
}

// Pointer to a type in the Types array
type TypePtr int

func RandomTypePtr() TypePtr {
	return TypePtr(rand.Intn(len(Types)))
}

func (t TypePtr) Void() bool {
	return t < 0 || t >= TypePtr(len(Types))
}

func (t *TypePtr) Next() bool {
	(*t)++
	return !t.Void()
}

func (t TypePtr) CNotation() string {
	if t.Void() {
		return "void"
	}
	return Types[t].CNotation()
}

func (t TypePtr) Short() string {
	if t.Void() {
		return "void"
	}
	return Types[t].Short()
}

func (t TypePtr) CapShort() string {
	if t.Void() {
		return "Void"
	}
	return Types[t].CapShort()
}

func (t TypePtr) GoNotation() string {
	if t.Void() {
		return ""
	}
	return Types[t].GoNotation()
}

func (t TypePtr) CgoNotation() string {
	if t.Void() {
		return ""
	}
	return Types[t].CgoNotation()
}

func (t TypePtr) Random() interface{} {
	if t.Void() {
		return ""
	}
	return Types[t].Random()
}

func (t TypePtr) CRandom() interface{} {
	if t.Void() {
		return ""
	}
	return Types[t].CRandom()
}

type TypeCombo []TypePtr

func NewCombo(num int) TypeCombo {
	c := make(TypeCombo, num)
	for i := range c {
		c[i] = -1
	}
	if len(c) > 0 {
		c[0] = -2
	}
	return c
}

func (c TypeCombo) New() bool {
	return c[0] == -2
}

func (c TypeCombo) Next() bool {
	if len(c) == 0 {
		return false
	}
	// Just initialized
	if c.New() {
		c[0] = -1
		return true
	}
	if c[0].Next() {
		return true
	}
	c[0] = 0
	return c[1:].Next()
}

func (c TypeCombo) NonVoid() (s []TypePtr) {
	for _, t := range c {
		if t.Void() {
			break
		}
		s = append(s, t)
	}
	return
}

func (c TypeCombo) pack(fn func(TypePtr) string) (s []string) {
	for _, t := range c.NonVoid() {
		s = append(s, fn(t))
	}
	return
}

func (c TypeCombo) packi(fn func(int, TypePtr) string) (s []string) {
	for i, t := range c.NonVoid() {
		s = append(s, fn(i, t))
	}
	return
}

func (c TypeCombo) CNotations() []string   { return c.pack(TypePtr.CNotation) }
func (c TypeCombo) Shorts() []string       { return c.pack(TypePtr.Short) }
func (c TypeCombo) CapShorts() []string    { return c.pack(TypePtr.CapShort) }
func (c TypeCombo) GoNotations() []string  { return c.pack(TypePtr.GoNotation) }
func (c TypeCombo) CgoNotations() []string { return c.pack(TypePtr.CgoNotation) }

// [int int uint] -> "IntIntUint"
// [] -> Void
func (c TypeCombo) FuncName() string {
	caps := c.CapShorts()
	if len(caps) == 0 {
		return "Void"
	}
	return strings.Join(caps, "")
}

// [int int uint] -> "void *ptr, int arg1, int arg2, unsigned int arg3"
// [] -> "void *ptr" / "void"
func (c TypeCombo) CArgs(withPtr bool) string {
	var args []string
	if withPtr {
		args = append(args, "void *ptr")
	}
	args = append(args, c.packi(func(i int, t TypePtr) string {
		return fmt.Sprintf("%v arg%d", t.CNotation(), i+1)
	})...)
	if args == nil {
		args = append(args, "void")
	}
	return strings.Join(args, ", ")
}

// [int int uint] -> "f func(c_int, c_int, c_uint), arg1 c_int, arg2 c_int, arg3 c_uint"
func (c TypeCombo) GoArgs(fn *Func) string {
	var args []string
	if fn != nil {
		args = append(args, fmt.Sprintf("f %v", fn.GoAnon()))
	}
	args = append(args, c.packi(func(i int, t TypePtr) string {
		return fmt.Sprintf("arg%d %v", i+1, t.GoNotation())
	})...)
	return strings.Join(args, ", ")
}

// [int int uint] -> "int, int, unsigned int"
// [] -> "void"
func (c TypeCombo) CTypes() string {
	cnotns := c.CNotations()
	if len(cnotns) == 0 {
		return "void"
	}
	return strings.Join(c.CNotations(), ", ")
}

// [int int uint] -> "c_int, c_int, c_uint"
func (c TypeCombo) GoTypes() string {
	return strings.Join(c.GoNotations(), ", ")
}

// [int int uint] -> "arg1, arg2, arg3"
func (c TypeCombo) ListPattern(name string) string {
	p := c.packi(func(i int, _ TypePtr) string {
		return fmt.Sprintf("%v%d", name, i+1)
	})
	return strings.Join(p, ", ")
}

func (c TypeCombo) List() string {
	return c.ListPattern("arg")
}

// [int int uint] -> "unsafe.Pointer(ptr), C.int(arg1), C.int(arg2), C.uint(arg3)"
func (c TypeCombo) CgoCall(withPtr bool) string {
	var s []string
	if withPtr {
		s = append(s, "unsafe.Pointer(ptr)")
	}
	s = append(s, c.packi(func(i int, t TypePtr) string {
		return fmt.Sprintf("%v(arg%d)", t.CgoNotation(), i+1)
	})...)
	return strings.Join(s, ", ")
}

type Func struct {
	RetType TypePtr
	Args    TypeCombo

	GenRandom bool
	MaxGens   uint64
	Gens      map[uint64]struct{}
}

func NewFunc(argnum int) *Func {
	f := &Func{
		RetType: -2,
		Args:    NewCombo(argnum),
	}
	if argnum == 0 {
		// Hack to make zero-args funcs work
		f.RetType = -1
	}
	return f
}

// From void:short:uint -> void f(short, unsigned int)
func FuncFromShortNotation(s string) *Func {
	determineType := func(short string) TypePtr {
		if short == "void" {
			return -1
		}
		for i, t := range Types {
			if t.Short() == short {
				return TypePtr(i)
			}
		}
		panic("Unknown type: " + short)
	}

	s = strings.ToLower(s)
	if s == "" {
		s = "void"
	}
	shorts := strings.Split(s, ":")
	types := make(TypeCombo, len(shorts))
	for i := range shorts {
		types[i] = determineType(shorts[i])
	}
	if len(types) < 2 {
		types = append(types, -1)
	}

	return &Func{
		RetType: types[0],
		Args:    types[1:],
	}
}

func (f *Func) GenerateRandomly(maxnum uint64) {
	if maxnum > AllPossibleFuncs() {
		maxnum = AllPossibleFuncs()
	}

	f.GenRandom = true
	f.MaxGens = maxnum
	f.Gens = make(map[uint64]struct{})
}

// Generate unique number for this combination of return
// and arguments types
func (f *Func) Stamp() (n uint64) {
	n = uint64(f.RetType + 1)
	for _, t := range f.Args.NonVoid() {
		n = uint64(len(Types)+1)*n + uint64(t)
	}
	return
}

func (f *Func) Next() bool {
	if f.GenRandom {
		return f.NextRandom()
	} else {
		return f.NextSeq()
	}
}

func (f *Func) NextSeq() bool {
	if !f.RetType.Next() {
		f.RetType = -1
		return f.Args.Next()
	}
	return true
}

func (f *Func) NextRandom() bool {
	if len(f.Gens) == int(f.MaxGens) {
		return false
	}
	for {
		f.Randomize()
		if _, ok := f.Gens[f.Stamp()]; !ok {
			break
		}
	}
	f.Gens[f.Stamp()] = struct{}{}
	return true
}

func (f *Func) Randomize() {
	f.RetType = TypePtr(rand.Intn(len(Types)+1) - 1)
	argnum := rand.Intn(len(f.Args) + 1)
	for i := 0; i < argnum; i++ {
		f.Args[i] = RandomTypePtr()
	}
	if argnum < len(f.Args) {
		f.Args[argnum] = -1
	}
}

func (f *Func) Void() bool {
	return f.RetType < 0
}

func (f *Func) Name() string {
	return f.RetType.CapShort() + "_" + f.Args.FuncName()
}

func (f *Func) GoAnon() string {
	s := fmt.Sprintf("func(%v)", f.Args.GoArgs(nil))
	if !f.Void() {
		s += " " + f.RetType.GoNotation()
	}
	return s
}

func (f *Func) WriteCDecl(w io.Writer) {
	fmt.Fprintf(w, "%v %v(%v) {\n", f.RetType.CNotation(), f.Name(), f.Args.CArgs(true))

	if !f.Void() {
		fmt.Fprintf(w, "\t%v ret;\n", f.RetType.CNotation())
	}
	fmt.Fprintln(w, "\tvolatile int sp1, sp2;")
	fmt.Fprintln(w, "\t"+`asm("mov %%esp, %0\n\t":"=r"(sp1));`)

	fmt.Fprint(w, "\t")
	if !f.Void() {
		fmt.Fprint(w, "ret = ")
	}
	var attr string
	if Stdcall {
		attr = "__attribute__((stdcall)) "
	}
	fmt.Fprintf(w, "((%v (%v*)(%v))ptr)(%v);\n", f.RetType.CNotation(), attr, f.Args.CTypes(), f.Args.List())

	fmt.Fprintln(w, "\t"+`asm volatile ("mov %%esp, %0\n\t":"=r"(sp2));`)
	fmt.Fprintln(w, "\tif (sp1 != sp2) {")
	fmt.Fprintf(w, "\t\t"+`printf("%v: Stack pointers do not match! %%d bytes difference.\n", sp2-sp1);`+"\n", f.Name())
	fmt.Fprintln(w, "\t\texit(1);")
	fmt.Fprintln(w, "\t}")

	if !f.Void() {
		fmt.Fprintln(w, "\treturn ret;")
	}
	fmt.Fprint(w, "}\n")
}

func (f *Func) WriteGoDecl(w io.Writer) {
	fmt.Fprintf(w, "func %v(%v)%v {\n", f.Name(), f.Args.GoArgs(f), " "+f.RetType.GoNotation())
	if Stdcall {
		fmt.Fprint(w, "\tptr := callback.NewStdcall(f)\n")
	} else {
		fmt.Fprint(w, "\tptr := callback.New(f)\n")
	}
	fmt.Fprint(w, "\t")
	if !f.Void() {
		fmt.Fprint(w, "ret := ")
	}
	fmt.Fprintf(w, "C.%v(%v)\n", f.Name(), f.Args.CgoCall(true))
	fmt.Fprint(w, "\tcallback.Remove(ptr)\n")
	if !f.Void() {
		fmt.Fprintf(w, "\treturn %v(ret)\n", f.RetType.GoNotation())
	}
	fmt.Fprint(w, "}\n")
}

func (f *Func) WriteGoTestDecl(w io.Writer) {
	fmt.Fprintf(w, "func Test%v(t *testing.T) {\n", f.Name())
	fmt.Fprintf(w, "\tvar called bool\n")
	for i, t := range f.Args.NonVoid() {
		fmt.Fprintf(w, "\tvar set%d %v\n", i+1, t.GoNotation())
		fmt.Fprintf(w, "\tvar expect%d %v = %v\n", i+1, t.GoNotation(), t.Random())
	}
	if !f.Void() {
		fmt.Fprintf(w, "\tvar ret %v\n", f.RetType.GoNotation())
		fmt.Fprintf(w, "\tvar rete %v = %v\n", f.RetType.GoNotation(), f.RetType.Random())
	}
	fmt.Fprint(w, "\t")
	if !f.Void() {
		fmt.Fprint(w, "ret = ")
	}
	fmt.Fprintf(w, "%v(%v {\n", f.Name(), f.GoAnon())
	fmt.Fprintf(w, "\t\tcalled = true\n")
	for i := range f.Args.NonVoid() {
		fmt.Fprintf(w, "\t\tset%d = arg%d\n", i+1, i+1)
	}
	if !f.Void() {
		fmt.Fprintf(w, "\t\treturn rete\n")
	}
	fmt.Fprintf(w, "\t}, %v)\n", f.Args.ListPattern("expect"))
	fmt.Fprint(w, "\tif !called {\n")
	fmt.Fprint(w, "\t\tt.Fatal(\"Not called\")\n")
	fmt.Fprint(w, "\t}\n")
	for i := range f.Args.NonVoid() {
		fmt.Fprintf(w, "\tif set%d != expect%d {\n", i+1, i+1)
		fmt.Fprintf(w, "\t\tt.Errorf(\"Arg %d: expected %%v, got %%v\", expect%d, set%d)\n", i+1, i+1, i+1)
		fmt.Fprint(w, "\t}\n")
	}
	if !f.Void() {
		fmt.Fprint(w, "\tif ret != rete {\n")
		fmt.Fprintf(w, "\t\tt.Errorf(\"Ret: expected %%v, got %%v\", rete, ret)\n")
		fmt.Fprint(w, "\t}\n")
	}
	fmt.Fprint(w, "}\n")
}

func WriteWarning(w io.Writer) {
	fmt.Fprintln(w, "// Generated using run.go")
	fmt.Fprintln(w, "// Do not keep in the project")
	fmt.Fprintln(w)
}

func NewCallsH() *os.File {
	f, err := os.Create("calls.h")
	if err != nil {
		panic(err)
	}
	WriteWarning(f)
	fmt.Fprintln(f, "#include <stdlib.h>")
	fmt.Fprintln(f, "#include <stdio.h>")
	fmt.Fprintln(f)
	return f
}

func NewCallsGo() *os.File {
	f, err := os.Create("calls.go")
	if err != nil {
		panic(err)
	}
	WriteWarning(f)
	fmt.Fprintln(f, "package tests")
	fmt.Fprintln(f)
	fmt.Fprintln(f, `// #cgo CFLAGS: -w`)
	fmt.Fprintln(f, `// #include "calls.h"`)
	fmt.Fprintln(f, `import "C"`)
	fmt.Fprintln(f, `import (`)
	fmt.Fprintln(f, "\t", `"unsafe"`)
	fmt.Fprintln(f)
	fmt.Fprintln(f, "\t", `"github.com/yamnikov-oleg/cgo-callback"`)
	fmt.Fprintln(f, ")")
	fmt.Fprintln(f)
	for _, t := range Types {
		fmt.Fprintf(f, "type %v %v\n", t.GoNotation(), t.CgoNotation())
	}
	fmt.Fprintln(f)
	for _, t := range Types {
		if !t.Ptr() {
			continue
		}
		fmt.Fprintf(f, "func new%v(arg %v) %v {\n", t.CapShort(), t.Deref().GoNotation(), t.GoNotation())
		fmt.Fprintf(f, "\tv := %v(arg)\n", t.Deref().CgoNotation())
		fmt.Fprintf(f, "\treturn %v(&v)\n", t.GoNotation())
		fmt.Fprint(f, "}\n")
		fmt.Fprintln(f)
	}
	return f
}

func NewCallsTestGo() *os.File {
	f, err := os.Create("calls_test.go")
	if err != nil {
		panic(err)
	}
	WriteWarning(f)
	fmt.Fprintln(f, "package tests")
	fmt.Fprintln(f)
	fmt.Fprintln(f, `import "testing"`)
	fmt.Fprintln(f)
	return f
}

func CleanUp() {
	os.Remove("calls.h")
	os.Remove("calls.go")
	os.Remove("calls_test.go")
}

var (
	Verbose      bool
	StopOnFail   bool
	FuncsPerTest uint
	MaxArgs      uint
	OnlyGen      bool
	MaxRandom    uint

	Stdcall bool

	SpecTest string
	SpecFile string
)

func init() {
	flag.BoolVar(&Verbose, "v", false, "Pass -v flag to \"go test\"")
	flag.BoolVar(&StopOnFail, "e", false, "Stop at the first failed test run")
	flag.UintVar(&FuncsPerTest, "fn", 100, "Number of functions to generate per test run")
	flag.UintVar(&MaxArgs, "arg", 3, "Maximum number of arguments to test")
	flag.BoolVar(&OnlyGen, "gen", false, "Generate tests but not run")
	flag.UintVar(&MaxRandom, "rand", 0, "Do not generate functions sequentially, but randomly. Generate up to `n` functions.")

	flag.BoolVar(&Stdcall, "stdcall", false, "Generate tests using stdcall convention (useless on non-x86 machines)")

	flag.StringVar(&SpecTest, "t", "", "Run only specific test, e.g. void:float:ushort")
	flag.StringVar(&SpecFile, "f", "", "Run specific tests from file, e.g. void:float:ushort\\nvoid:double")
}

var allPossibleFuncs uint64

func AllPossibleFuncs() uint64 {
	if allPossibleFuncs > 0 {
		return allPossibleFuncs
	}
	allPossibleFuncs = 1
	for i := uint(0); i <= MaxArgs; i++ {
		allPossibleFuncs += uint64(math.Pow(float64(len(Types)), float64(i+1)))
	}
	return allPossibleFuncs
}

func CalcProgress(run int) (passed, all uint64) {
	if MaxRandom > 0 {
		all = uint64(MaxRandom)
		if all > AllPossibleFuncs() {
			all = AllPossibleFuncs()
		}
	} else {
		all = AllPossibleFuncs()
	}

	passed = uint64(run) * uint64(FuncsPerTest)
	if passed > all {
		passed = all
	}

	return
}

func RunTests(run int) bool {
	if OnlyGen {
		os.Exit(0)
	}

	if SpecTest == "" && SpecFile == "" {
		fnsPassed, fnsAll := CalcProgress(run)
		fmt.Printf("--- Party %d (%d/%d) ---\n", run, fnsPassed, fnsAll)
	}

	args := []string{"test"}
	if Verbose {
		args = append(args, "-v")
	}
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		return true
	}
	if _, ok := err.(*exec.ExitError); !ok {
		panic(err)
	}
	return false
}

func RunSpecificFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	ch := NewCallsH()
	cg := NewCallsGo()
	ctg := NewCallsTestGo()
	for scanner.Scan() {
		fn := FuncFromShortNotation(scanner.Text())
		fn.WriteCDecl(ch)
		fn.WriteGoDecl(cg)
		fn.WriteGoTestDecl(ctg)
	}
	ch.Close()
	cg.Close()
	ctg.Close()
	success := RunTests(1)
	if !success {
		os.Exit(1)
	}
	CleanUp()
}

func RunSpecificTest(test string) {
	fn := FuncFromShortNotation(test)

	ch := NewCallsH()
	fn.WriteCDecl(ch)
	ch.Close()

	cg := NewCallsGo()
	fn.WriteGoDecl(cg)
	cg.Close()

	ctg := NewCallsTestGo()
	fn.WriteGoTestDecl(ctg)
	ctg.Close()

	success := RunTests(1)
	if !success {
		os.Exit(1)
	}
	CleanUp()
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/yamnikov-oleg/cgo-callback/tests")

	flag.Parse()
	if SpecFile != "" {
		RunSpecificFile(fmt.Sprintf("%v/%v", wd, SpecFile))
		return
	}
	if SpecTest != "" {
		RunSpecificTest(SpecTest)
		return
	}

	var (
		funcs uint = 0
		run        = 1

		ch  *os.File
		cg  *os.File
		ctg *os.File

		success bool = true
	)

	prepare := func() {
		ch = NewCallsH()
		cg = NewCallsGo()
		ctg = NewCallsTestGo()
	}

	test := func() {
		ch.Close()
		cg.Close()
		ctg.Close()
		sc := RunTests(run)
		if success {
			success = sc
		}
		if !success && StopOnFail {
			fmt.Println("Test failed, exiting.")
			os.Exit(1)
		}
		run++
		funcs = 0
	}

	fn := NewFunc(int(MaxArgs))
	if MaxRandom > 0 {
		fn.GenerateRandomly(uint64(MaxRandom))
	}
	for fn.Next() {
		if funcs == 0 {
			prepare()
		}
		fn.WriteCDecl(ch)
		fn.WriteGoDecl(cg)
		fn.WriteGoTestDecl(ctg)
		funcs++
		if funcs >= FuncsPerTest {
			test()
		}
	}
	// Some functions left untested
	if funcs != 0 {
		test()
	}
	CleanUp()
}
