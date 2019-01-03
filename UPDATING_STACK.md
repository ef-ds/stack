## Updating Stack and Checking the Results

If you want to make changes to stack and run the tests to check the effect on performance and memory,
we suggest you run all the benchmark tests locally once using below command.

```
go test -benchmem -count 10 -timeout 60m -bench="stack*" -run=^$ > testdata/stack.txt
```

Then make the changes and re-run the tests using below command (notice the output file now is stack2.txt).

```
go test -benchmem -count 10 -timeout 60m -bench="stack*" -run=^$ > testdata/stack2.txt
```

Then run the [test-splitter](https://github.com/ef-ds/tools/tree/master/testsplitter) tool as follow:

```
go run *.go --file PATH_TO_TESTDATA/stack2.txt --suffix 2
```

Test-splitter should create each file with the "2" suffix, so now we have the test file for both, the old and this new
test run. Use below commands to test the effect of the changes for each test suite.

```
benchstat testdata/BenchmarkMicroservice.txt testdata/BenchmarkMicroservice2.txt
benchstat testdata/BenchmarkFill.txt testdata/BenchmarkFill2.txt
benchstat testdata/BenchmarkRefill.txt testdata/BenchmarkRefill2.txt
benchstat testdata/BenchmarkRefillFull.txt testdata/BenchmarkRefillFull2.txt
benchstat testdata/BenchmarkSlowIncrease.txt testdata/BenchmarkSlowIncrease2.txt
benchstat testdata/BenchmarkSlowIncrease.txt testdata/BenchmarkSlowIncrease2.txt
benchstat testdata/BenchmarkStable.txt testdata/BenchmarkStable2.txt
```
