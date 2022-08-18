package main

import (
	"github.com/loft-sh/vcluster-sdk/plugin"
	"github.com/loft-sh/vcluster-sync-all-configmaps/syncers"
)

func main() {
	ctx := plugin.MustInit()
	plugin.MustRegister(syncers.NewConfigMapSyncer(ctx))
	plugin.MustStart()
}
