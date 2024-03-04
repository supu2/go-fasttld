package fasttld

import (
	"log"
	"encoding/json" // Import the encoding/json package
	"fmt"
	"github.com/elliotwutingfeng/go-fasttld"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var includePrivateSuffix, ignoreSubDomains, toPunyCode, outputJson  bool
var cacheFilePath string

var extractCmd = &cobra.Command{
	Use:     "extract",
	Aliases: []string{"ext"},
	Short:   "Extracts subcomponents from a URL.",
	Long: `Extracts subcomponents from a URL.

For Example
---
fasttld extract abc.example.com:5000/a/path
---
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		extractor, err := fasttld.New(fasttld.SuffixListParams{IncludePrivateSuffix: includePrivateSuffix, CacheFilePath: cacheFilePath})
		if err != nil {
			log.Fatal(err)
		}
		res, err := extractor.Extract(fasttld.URLParams{URL: args[0], IgnoreSubDomains: ignoreSubDomains, ConvertURLToPunyCode: toPunyCode})
		if err != nil {
			color.New(color.FgHiRed, color.Bold).Print("Error: ")
			color.New(color.FgHiWhite).Println(err)
		}
		if outputJson {
			// Convert the response to JSON
			jsonRes, err := json.Marshal(res)
			if err != nil {
				log.Fatal(err)
			}

			// Print the JSON response
			fmt.Println(string(jsonRes))	
		} else {
			fasttld.PrintRes(args[0], res)
		}
	},
}

func init() {
	extractCmd.Flags().BoolVarP(&includePrivateSuffix, "private-suffix", "p", false, "Include private suffix")
	extractCmd.Flags().BoolVarP(&ignoreSubDomains, "ignore-subdomains", "i", false, "Ignore subdomains")
	extractCmd.Flags().BoolVarP(&toPunyCode, "to-punycode", "t", false, "Convert to punycode")
	extractCmd.Flags().BoolVarP(&outputJson, "output-json", "j", false, "JSON output format")
	extractCmd.Flags().StringVarP(&cacheFilePath, "cache-file-path", "c", "/tmp", "Specify the cache path")
	rootCmd.AddCommand(extractCmd)
}
