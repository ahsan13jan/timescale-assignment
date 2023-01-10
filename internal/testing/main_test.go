//go:build integration
// +build integration

package testing

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/kelseyhightower/envconfig"

	"timescale/internal/config"
	"timescale/internal/db"
	"timescale/internal/logger"
	"timescale/internal/stats"
	util "timescale/internal/testing/test_util"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty",
}

var (
	conf        config.Config
	clientLocal *util.PGClientLocal
	client      *db.PGClient
	err         error
	log         = logger.GetLogger()
	metrics     = stats.StatsOutput{}
)

func init() {
	godog.BindCommandLineFlags("godog.", &opts)

	envconfig.MustProcess("", &conf)
	fmt.Println(conf.String())
	clientLocal, err = util.New(conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err = db.New(conf)
	if err != nil {
		log.Fatalln(err)
	}

	// aggregate = pg.NewAggregator(client, pg.NewStats())

}

func TestMain(m *testing.M) {
	flag.Parse()
	opts.Paths = flag.Args()

	status := godog.TestSuite{
		Name:                 "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^following data is in cpu usage table$`, theGivenCpuUsage)
	ctx.Step(`^aggregated per minute of host ([^\s]+) result should be from ([\d-T:\.Z]+) to ([\d-T:\.Z]+)$`, theAggregateShouldBe)
	ctx.Step(`^cli tool run with following input csv path (.*)$`, theCliToolExecutes)
	ctx.Step(`^num of quries in the metrics should be (.*)$`, theNumOfQuriesShouldBe)
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {

	metrics = stats.StatsOutput{}
	err := clientLocal.DeleteAll(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
}
