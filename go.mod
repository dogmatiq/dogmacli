module github.com/dogmatiq/dogmacli

go 1.16

replace github.com/dogmatiq/configkit v0.11.1-0.20210625234826-0eedd8ec2931 => ../configkit

require (
	github.com/dogmatiq/configkit v0.11.1-0.20210625234826-0eedd8ec2931
	github.com/spf13/cobra v1.1.3
	golang.org/x/tools v0.0.0-20201224043029-2b0845dc783e
)
