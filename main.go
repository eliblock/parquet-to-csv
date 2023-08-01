package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	goparquet "github.com/fraugster/parquet-go"
	"github.com/pkg/errors"
)

var (
	version = "dev"
)

func main() {
	versionFlag := flag.Bool("version", false, "")

	inputFlag := flag.String("in", "", "OPTIONAL: path to parquet file (if not set, first non-flag argument is used)")

	outputFlag := flag.String("out", "", "OPTIONAL: path to output csv file (defaults to standard out)")
	overwriteFlag := flag.Bool("overwrite", false, "OPTIONAL: allow overwriting an existing file (defaults to false)")

	rowLimitFlag := flag.Int("n", 0, "OPTIONAL: limit the number of rows to convert (defaults to all rows)")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("%s\n", version)
		return
	}

	if err := validateFlagsAndArgs(inputFlag, outputFlag, overwriteFlag); err != nil {
		log.Fatal(errors.Wrap(err, "invalid flags or arguments"))
	}

	// Open the output file for writing, erroring if we can't and defaulting to StdOut
	var outputFile *os.File
	var err error
	if *outputFlag != "" {
		outputFile, err = os.OpenFile(*outputFlag, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(errors.Wrapf(err, "failed to open %q as output file", *outputFlag))
		}
		defer outputFile.Close()
	} else {
		outputFile = os.Stdout
	}

	// Open the input file for reading
	var inputFile *os.File
	if *inputFlag != "" {
		inputFile, err = os.Open(*inputFlag)
		if err != nil {
			log.Fatal(errors.Wrapf(err, "failed to open %q as input file (using the value from the -in flag)", *inputFlag))
		}
	} else {
		inputFile, err = os.Open(flag.Args()[0])
		if err != nil {
			log.Fatal(errors.Wrapf(err, "failed to open %q as input file (using the first non-flag argument)", flag.Args()[0]))
		}
	}
	defer inputFile.Close()

	if err := convertAndOutput(inputFile, outputFile, *rowLimitFlag); err != nil {
		log.Fatal(errors.Wrap(err, "failed to convert parquet file to csv"))
	}
}

func convertAndOutput(inputFile, outputFile *os.File, rows int) error {
	// Construct a parquet reader on the input
	fr, err := goparquet.NewFileReader(inputFile)
	if err != nil {
		return err
	}

	// Construct a csv writer on the output
	fw := csv.NewWriter(outputFile)
	defer fw.Flush()

	// Note the file schema, and print the csv headers
	schema := fr.GetSchemaDefinition()
	log.Printf("Observed file schema: %s", schema)
	var rowHeaders []string
	for _, column := range schema.RootColumn.Children {
		rowHeaders = append(rowHeaders, column.SchemaElement.Name)
	}
	if err := fw.Write(rowHeaders); err != nil {
		return fmt.Errorf(errors.Wrap(err, "failed to write headers").Error())
	}

	count := 0
	for {
		// Read the next row
		count++
		row, err := fr.NextRow()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf(errors.Wrap(err, "reading record failed").Error())
		}

		// Output the row
		var rowStrings []string
		for _, v := range row {
			// Record the value, stringifying by casting or printing
			if vv, ok := v.([]byte); ok {
				rowStrings = append(rowStrings, string(vv))
			} else {
				rowStrings = append(rowStrings, fmt.Sprintf("%v", v))
			}
		}
		if err := fw.Write(rowStrings); err != nil {
			return fmt.Errorf(errors.Wrap(err, "failed to write row").Error())
		}

		// Status messages and row limits
		if count == rows {
			fw.Flush()
			log.Printf("Reached row limit (%d rows)", rows)
			return nil
		} else if count%1000 == 0 {
			fw.Flush()
			log.Printf("Processed %d records...", count)
		}
	}

	log.Printf("Reached end of file (%d records)", count)
	return nil
}

func validateFlagsAndArgs(inputFlag, outputFlag *string, overwriteFlag *bool) error {
	argCount := len(flag.Args())
	if !(argCount == 1 || (argCount == 0 && *inputFlag != "")) {
		return fmt.Errorf("must specify the -in flag or exactly one non-flag argument (you had %d arguments and a flag of %q)", argCount, *inputFlag)
	} else if argCount != 0 && *inputFlag != "" {
		return fmt.Errorf("must specify either the --in flag or a non-flag argument but not both")
	}

	if !*overwriteFlag && *outputFlag != "" {
		if _, err := os.Stat(*outputFlag); err == nil {
			return fmt.Errorf("output file %s already exists. Either specify --overwrite or pick a new path", *outputFlag)
		}
	}
	return nil
}
