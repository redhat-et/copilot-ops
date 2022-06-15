# Example Files for copilot-ops

This directory contains a collection of example files which can be used in tandem
with the copilot-ops tool.
Feel free to play around with these files and see how `copilot-ops` responds.

### stock-data.yaml 


```sh
# to change values within .data
go run main.go patch -f stock-data.yaml --request "@stock-data.yaml needs an additional field to hold 50000 units of AMC stock"

# to rename the configmap
go run main.go patch -f stock-data.yaml --request "rename the configmap to 'stock-holdings'"

# to erase existing fields
go run main.go patch -f examples/stock-data.yaml --request "remove the existing data fields in @stock-data.yaml"
```
<!--  -->