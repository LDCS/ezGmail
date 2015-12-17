# ezGmail
ezGmail is a golang wrapper for the gmail API.  It is designed to be an easy way to access received email and process them.

The package is based on the Gmail API,  and so uses OAuth 2.0 authentication,  and the code is based on sample program provided in Gmail API's Go quickstart guide.

To download ezGmail and its dependencies,  use the command:
    ```go get github.com/LDCS/ezGmail```

To allow ezGmail to access a gmail account,  the Gmail API needs to be enabled for that account,  and the client_secret.json file needs to be in the same directory as the main program.

Follow step 1 in https://developers.google.com/gmail/api/quickstart/go to enable the Gmail API and download the client_secret.json file.

Here's a sample program that uses some of the functionality provided by the package:

```go
package main

import (
        "fmt"
	"github.com/LDCS/ezGmail"
)

func main() {
	var gs ezGmail.GmailService
	// InitSrv() uses client_secret.json to try to get a OAuth 2.0 token,  , if not present already.
	gs.InitSrv()

	// We compose a search statement with filter functions
	gs.InInbox().MaxResults(1).NewerThanRel("10d").Match("-address").HasAttachment(true)

	// GetMessages() tries to execute the search statement and get a list of messages
	for _, ii := range(gs.GetMessages()) {
		fmt.Println("\nTrying Subject")
		if ii.HasSubject() { fmt.Println(ii.GetSubject()) }
		fmt.Println("\nTrying BodyText")
		if ii.HasBodyText()    { fmt.Println(string(ii.GetBodyText())) }
		fmt.Println("\nTrying BodyHtml")
		if ii.HasBodyHtml()    { fmt.Println(string(ii.GetBodyHtml())) }
		fmt.Println("\nTrying Attachments")
		if ii.HasAttachments() {
			for _, jj := range(ii.GetAttachments()) {
				fmt.Println("\nMimeType")
				fmt.Println(jj.GetMimeType())
				if jj.GetFilename() == "readme.txt" {
					fmt.Println(string(jj.GetData()))
				}
			}
		}
	}
}
```
