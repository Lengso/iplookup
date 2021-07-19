package runner

import (
	"flag"
	"io"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/projectdiscovery/gologger"
)

// Options contains the configuration options for tuning
// the subdomain enumeration process.
type Options struct {
	Verbose            bool   // Verbose flag indicates whether to show verbose output or not
	NoColor            bool   // No-Color disables the colored output
	JSON               bool   // JSON specifies whether to use json for output format or text file
	Silent             bool   // Silent suppresses any extra text and only writes subdomains to screen
	ListSources        bool   // ListSources specifies whether to list all available sources
	Stdin              bool   // Stdin specifies whether stdin input was given to the process
	Version            bool   // Version specifies if we should just show version and exit
	All                bool   // All specifies whether to use all (slow) sources.
	Threads            int    // Thread controls the number of threads to use for active enumerations
	Timeout            int    // Timeout is the seconds to wait for sources to respond
	MaxEnumerationTime int    // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
	Threshold          int    // Threshold is Number of domain name thresholds
	Ip                 string // Ip is the domain to find subdomains for
	IpsFile            string // IpsFile is the file containing list of domains to find subdomains for
	Output             io.Writer
	OutputFile         string // Output is the file to write found subdomains to.
	OutputDirectory    string // OutputDirectory is the directory to write results to in case list of domains is given
	Sources            string // Sources contains a comma-separated list of sources to use for enumeration
	ExcludeSources     string // ExcludeSources contains the comma-separated sources to not include in the enumeration process
	ConfigFile         string // ConfigFile contains the location of the config file

	YAMLConfig ConfigFile // YAMLConfig contains the unmarshalled yaml config file
}

// ParseOptions parses the command line flags provided by a user
func ParseOptions() *Options {
	options := &Options{}

	config, err := GetConfigDirectory()
	if err != nil {
		// This should never be reached
		gologger.Fatal().Msgf("Could not get user home: %s\n", err)
	}

	flag.BoolVar(&options.Verbose, "v", false, "Show Verbose output")
	flag.BoolVar(&options.NoColor, "nC", false, "Don't Use colors in output")
	flag.IntVar(&options.Threads, "t", 10, "Number of concurrent goroutines for resolving")
	flag.IntVar(&options.Timeout, "timeout", 30, "Seconds to wait before timing out")
	flag.BoolVar(&options.JSON, "json", false, "Write output in JSON lines Format")
	flag.IntVar(&options.MaxEnumerationTime, "max-time", 10, "Minutes to wait for enumeration results")
	flag.IntVar(&options.Threshold, "count", 50, "  Number of domain name thresholds")
	flag.StringVar(&options.Ip, "i", "", "ip to find domain for")
	flag.StringVar(&options.IpsFile, "iL", "", "File containing list of ips to enumerate")
	flag.StringVar(&options.OutputFile, "o", "", "File to write output to (optional)")
	flag.StringVar(&options.OutputDirectory, "oD", "", "Directory to write enumeration results to (optional)")
	flag.BoolVar(&options.Silent, "silent", false, "Show only subdomains in output")
	flag.BoolVar(&options.All, "all", false, "Use all sources (slow) for enumeration")
	flag.StringVar(&options.Sources, "sources", "", "Comma separated list of sources to use")
	//flag.BoolVar(&options.ListSources, "ls", false, "List all available sources")
	flag.StringVar(&options.ExcludeSources, "exclude-sources", "", "List of sources to exclude from enumeration")
	flag.StringVar(&options.ConfigFile, "config", path.Join(config, "config.yaml"), "Configuration file for API Keys, etc")
	flag.BoolVar(&options.Version, "version", false, "Show version of iplookup")
	flag.Parse()

	// Default output is stdout
	options.Output = os.Stdout

	// Check if stdin pipe was given
	options.Stdin = hasStdin()

	// Read the inputs and configure the logging
	options.configureOutput()

	// Show the user the banner
	ShowBanner()

	if options.Version {
		gologger.Info().Msgf("Current Version: %s\n", Version)
		os.Exit(0)
	}

	// Check if the config file exists. If not, it means this is the
	// first run of the program. Show the first run notices and initialize the config file.
	// Else show the normal banners and read the yaml fiile to the config
	//检查配置文件是否存在。 如果不是，则表示这是
	//程序的第一次运行。 显示首次运行通知并初始化配置文件。
	//否则显示正常横幅，并读取yaml文件，使其符合配置
	if !CheckConfigExists(options.ConfigFile) {
		options.firstRunTasks()
	} else {
		options.normalRunTasks()
	}

	if options.ListSources {
		listSources(options)
		os.Exit(0)
	}

	// Validate the options passed by the user and if any
	// invalid options have been used, exit.
	err = options.validateOptions()
	if err != nil {
		gologger.Fatal().Msgf("Program exiting: %s\n", err)
	}

	return options
}

func hasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	isPipedFromChrDev := (stat.Mode() & os.ModeCharDevice) == 0
	isPipedFromFIFO := (stat.Mode() & os.ModeNamedPipe) != 0

	return isPipedFromChrDev || isPipedFromFIFO
}

func listSources(options *Options) {
	gologger.Info().Msgf("Current list of available sources. [%d]\n", len(options.YAMLConfig.AllSources))
	gologger.Info().Msgf("Sources marked with an * needs key or token in order to work.\n")
	gologger.Info().Msgf("You can modify %s to configure your keys / tokens.\n\n", options.ConfigFile)

	keys := options.YAMLConfig.GetKeys()
	needsKey := make(map[string]interface{})
	keysElem := reflect.ValueOf(&keys).Elem()
	for i := 0; i < keysElem.NumField(); i++ {
		needsKey[strings.ToLower(keysElem.Type().Field(i).Name)] = keysElem.Field(i).Interface()
	}

	for _, source := range options.YAMLConfig.AllSources {
		message := "%s\n"
		if _, ok := needsKey[source]; ok {
			message = "%s *\n"
		}
		gologger.Silent().Msgf(message, source)
	}
}
