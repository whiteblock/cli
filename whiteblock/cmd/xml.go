package cmd

// import (
// 	"encoding/xml"
// 	"fmt"

// 	"github.com/spf13/cobra"
// )

// var xmlCmd = &cobra.Command{
// 	Use: "xml",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		data := `
// 			<get>
// 			<stats end="-1" start="-1" live="false" engine="1">
// 				<path id="PATH_1">
// 				<overall source="PORT_A">
// 					<tx_frames/>
// 					<tx_bytes/>
// 					<drops/>
// 				</overall>
// 				<overall source="PORT_B">
// 					<tx_frames/>
// 					<tx_bytes/>
// 					<drops/>
// 				</overall>
// 				<wan_access ep="PORT_A">
// 					<outbound>
// 					<tx_frames/>
// 					<tx_bytes/>
// 					<qdrops_frames/>
// 					<queue_frames/>
// 					<queue_bytes/>
// 					<bg_frames/>
// 					<bg_bytes/>
// 					<bg_qdrops/>
// 					</outbound>
// 					<inbound>
// 					<tx_frames/>
// 					<tx_bytes/>
// 					<qdrops_frames/>
// 					<queue_frames/>
// 					<queue_bytes/>
// 					<bg_frames/>
// 					<bg_bytes/>
// 					<bg_qdrops/>
// 					</inbound>
// 				</wan_access>
// 				<wan_access ep="PORT_B">
// 					<outbound>
// 					<tx_frames/>
// 					<tx_bytes/>
// 					<qdrops_frames/>
// 					<queue_frames/>
// 					<queue_bytes/>
// 					<bg_frames/>
// 					<bg_bytes/>
// 					<bg_qdrops/>
// 					</outbound>
// 					<inbound>
// 					<tx_frames/>
// 					<tx_bytes/>
// 					<qdrops_frames/>
// 					<queue_frames/>
// 					<queue_bytes/>
// 					<bg_frames/>
// 					<bg_bytes/>
// 					<bg_qdrops/>
// 					</inbound>
// 				</wan_access>
// 				<wan source="PORT_A">
// 					<drops_loss/>
// 					<dups/>
// 					<reorders/>
// 					<corrupts/>
// 				</wan>
// 				<wan source="PORT_B">
// 					<drops_loss/>
// 					<dups/>
// 					<reorders/>
// 					<corrupts/>
// 				</wan>
// 				</path>
// 				<bypass source_ep="PORT_A" dest_ep="PORT_B">
// 				<tx_frames/>
// 				<tx_bytes/>
// 				</bypass>
// 				<bypass source_ep="PORT_B" dest_ep="PORT_A">
// 				<tx_frames/>
// 				<tx_bytes/>
// 				</bypass>
// 			</stats>
// 			<status engine="1">
// 				<interface id="PORT_1"/>
// 				<interface id="PORT_2"/>
// 			</status>
// 			</get>
// 		`

// 		v := Result{Name: "none", Phone: "none"}
// 		err := xml.Unmarshal([]byte(data), &v)
// 		if err != nil {
// 			fmt.Printf("error: %v", err)
// 			return
// 		}
// 		fmt.Printf("XMLName: %#v\n", v.XMLName)
// 		fmt.Printf("Name: %q\n", v.Name)
// 		fmt.Printf("Phone: %q\n", v.Phone)
// 		fmt.Printf("Email: %v\n", v.Email)
// 		fmt.Printf("Groups: %v\n", v.Groups)
// 		fmt.Printf("Address: %v\n", v.Address)
// 	},
// }

// func init() {
// 	RootCmd.AddCommand(xmlCmd)
// }
