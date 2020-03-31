package cmd

import(
	"k8s.io/klog"
	"github.com/spf13/cobra"
	"github.com/jwzl/beehive/pkg/core"
	"github.com/jwzl/edgecloud/pkg/web"
	"github.com/jwzl/edgecloud/pkg/eventhub"
	"github.com/jwzl/edgecloud/pkg/devicetwin"
)


/*
* new app command
*/
func NewAppCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use: "edgecloud",
		Long: `edgecloud is edgeOn server on cloud.  `,
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: To help debugging, immediately log version
			klog.Infof("###########  Start the edgeOn Server...! ###########")
			registerModules()
			// start all modules
			core.Run()
		},
	}

	return cmd
}

// register all module into beehive.
func registerModules(){
	web.Register()		//http server.
	devicetwin.Register()	
	eventhub.Register()	
}
