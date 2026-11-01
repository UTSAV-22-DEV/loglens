
# Go Notes — LogLens

## 1. Variables and Syntax

- := Declare and initialize inside functions.
- var Declare variables.
- const Declare constants.
- func Declare functions.
- return Return values.
- nil Represents absence of a value for nil-capable types.

## 2. Built-in Functions

- len(x): Length.
- cap(x): Capacity.
- make(): Allocate slices, maps and channels.
- append(): Add elements to slices.
- copy(): Copy slice elements.
- new(): Allocate zero-valued object and return pointer.

## 3. Data Structures

- Array: Fixed-size collection.
- Slice: Flexible view over an underlying array.
- Map: Key-value collection.
- Struct: Group related fields.
- Interface: Defines behaviour through methods.

## 4. Pointers

- &x: Address of x.
- *p: Dereference pointer.
- *int: Pointer to integer.
- Pointer receiver: Can modify original object.

## 5. File Handling

- os.Open(): Open file.
- os.WriteFile(): Write file.
- file.Close(): Close file.
- defer: Schedule cleanup.
- bufio.NewScanner(): Create Scanner.
- scanner.Scan(): Read next token.
- scanner.Text(): Get current token.
- scanner.Err(): Check reading error.
- scanner.Buffer(): Configure buffer and token limit.

## 6. Memory Management

- byte: Alias for uint8.
- []byte: Slice of bytes.
- len(): Current slice length.
- cap(): Slice capacity.
- make([]byte, 0, n): Empty slice with capacity n.
- Buffer: Temporary memory for incoming data.

## 7. Error Handling

- error: Built-in error interface.
- fmt.Errorf(): Create formatted error.
- %w: Wrap an existing error.
- errors.Is(): Check error identity.
- errors.As(): Inspect error type.

## 8. Testing

- testing.T: Unit testing.
- t.Run(): Subtests.
- t.TempDir(): Temporary directory.
- t.Fatalf(): Fail immediately.
- t.Errorf(): Record failure.
- testing.B: Benchmarking.
- b.ResetTimer(): Reset benchmark timer.
- b.ReportAllocs(): Report allocations.
- b.SetBytes(): Set processed bytes per iteration.

## 9. Useful Commands

go run .
go build
go test ./...
go test -count=1 ./...
go test -v ./...
go test -bench=. -benchmem
go vet ./...
gofmt -w .
