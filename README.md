# vorto-vrp

Run on one training set:

```
go run main.go Training\ Problems/problem17.txt
```

Build solution:

```
go build main.go
```

Run evaluation:

```
python3 evaluateShared.py --cmd "./vorto-vrp" --problemDir Training\ Problems
```
