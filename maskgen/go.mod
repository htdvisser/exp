module htdvisser.dev/exp/maskgen

go 1.14

replace htdvisser.dev/exp/clicontext => ../clicontext

replace htdvisser.dev/exp/stringslice => ../stringslice

require (
	github.com/fatih/structtag v1.2.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/tools v0.0.0-20200713160548-f739c553ea89
	htdvisser.dev/exp/clicontext v0.0.0-20200711072243-a90143cf50fc
	htdvisser.dev/exp/stringslice v0.0.0-20200711072243-a90143cf50fc
)
