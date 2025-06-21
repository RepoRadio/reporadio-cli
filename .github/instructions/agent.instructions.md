---
applyTo: '**'
---
We write Go code using the standard library and `go test` for testing. Avoid suggesting alternative tools unless asked.

Use Test Driven Development (TDD) practices. Write tests first to guide implementation, but expect that some tests may not be committed or used in automated testing.

Prefer writing as little code as possible. Avoid creating abstractions or helpers until their necessity is clearly discussed and justified.

Use a flat folder structure by default. Only suggest subfolders when there is a well-reasoned need.

Prefer integration tests and stand-in services over unit tests. Only write unit tests to verify specific programming logic.

Keep implementations simple and direct. Only introduce complexity after discussing and validating the need for it.

Do not overplan. Focus only on the current issue. Plan and write code for a single issue at a time.
