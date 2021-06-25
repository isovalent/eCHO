# Generating new episode notes

* Download the episode planning spreadsheet and call it ./planning.csv
* `go run main.go`. You'll be prompted for 
  * the line number in the spreadsheet (this is one more than the episode number)
  * some info that might not be in the spreadsheet already 
* This should automatically generate a directory containing a README file for the episode
  * That README file also contains some updates for copying into the [main eCHO README](../README.md)