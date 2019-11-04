package ezGmail

import (
        "encoding/json"
	"encoding/base64"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "net/url"
        "os"
        "os/user"
        "path/filepath"
	"strings"
	
        "golang.org/x/net/context"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/gmail/v1"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
        cacheFile, err := tokenCacheFile()
        if err != nil {
                log.Fatalf("Unable to get path to cached credential file. %v", err)
        }
        tok, err := tokenFromFile(cacheFile)
        if err != nil {
                tok = getTokenFromWeb(config)
                saveToken(cacheFile, tok)
        }
        return config.Client(ctx, tok)
}


// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
        authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
        fmt.Printf("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var code string
        if _, err := fmt.Scan(&code); err != nil {
                log.Fatalf("Unable to read authorization code %v", err)
        }

        tok, err := config.Exchange(oauth2.NoContext, code)
        if err != nil {
                log.Fatalf("Unable to retrieve token from web %v", err)
        }
        return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
        usr, err := user.Current()
        if err != nil {
                return "", err
        }
        tokenCacheDir := filepath.Join(usr.HomeDir, ".ezGmail")
        os.MkdirAll(tokenCacheDir, 0700)
        return filepath.Join(tokenCacheDir,
                url.QueryEscape("ezGmail.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
        f, err := os.Open(file)
        if err != nil {
                return nil, err
        }
        t := &oauth2.Token{}
        err = json.NewDecoder(f).Decode(t)
        defer f.Close()
        return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
        fmt.Printf("Saving credential file to: %s\n", file)
        f, err := os.Create(file)
        if err != nil {
                log.Fatalf("Unable to cache oauth token: %v", err)
        }
        defer f.Close()
        json.NewEncoder(f).Encode(token)
}

type GmailService struct {
        srv              *gmail.Service

	sUser            string
	sLabel           string
	iMaxResults      int64
	//search query data
        sFrom            string //
        sTo              string //
        sOlder           string //
        sNewer           string //
        sOlderRel        string //
        sNewerRel        string //
	sSubject         string //
	sInPlace         string //in:{inbox,sent,trash,spam,anywhere}
	sLarger          string //
	sSmaller         string //
	sFilename        string //
        sMatch           string
        sMatchExact      string
	sHasAttachment   bool
}

/*
	credentialFilePath is the loaction of 'credential-json' file
 */
func (gs *GmailService) InitSrvWithCrentialAt(credentialFilePath string) {
	// Connect and create gmail.Service object
	ctx := context.Background()

	b, err := ioutil.ReadFile(credentialFilePath)
	if err != nil {
		log.Fatalf("Unable to read client secret file at %s : %v", credentialFilePath,err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	gs.srv, err = gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	// Assign Defaults
	gs.iMaxResults = 50
	gs.sUser       = "me"
}

func (gs *GmailService) InitSrv() {
	gs.InitSrvWithCrentialAt("client_secret.json")
}

func (gs *GmailService) MaxResults	(maxres         int64 ) *GmailService	{ gs.iMaxResults	= maxres    ; return gs }
func (gs *GmailService) From		(from		string)	*GmailService   { gs.sFrom		= from	    ; return gs }
func (gs *GmailService) To		(to		string) *GmailService	{ gs.sTo		= to	    ; return gs }
func (gs *GmailService) OlderThan	(older		string) *GmailService	{ gs.sOlder		= older	    ; return gs }
func (gs *GmailService) NewerThan	(newer		string) *GmailService	{ gs.sNewer		= newer	    ; return gs }
func (gs *GmailService) OlderThanRel	(older		string) *GmailService	{ gs.sOlderRel		= older	    ; return gs }
func (gs *GmailService) NewerThanRel	(newer		string) *GmailService	{ gs.sNewerRel		= newer	    ; return gs }
func (gs *GmailService) Subject		(subject	string) *GmailService   { gs.sSubject		= subject   ; return gs }
func (gs *GmailService) InInbox		()		        *GmailService	{ gs.sInPlace		= "inbox"   ; return gs }
func (gs *GmailService) InSent		()		        *GmailService	{ gs.sInPlace		= "sent"    ; return gs }
func (gs *GmailService) InTrash		()		        *GmailService	{ gs.sInPlace		= "trash"   ; return gs }
func (gs *GmailService) InSpam		()		        *GmailService	{ gs.sInPlace		= "spam"    ; return gs }
func (gs *GmailService) InAnywhere	()		        *GmailService	{ gs.sInPlace		= "anywhere"; return gs }
func (gs *GmailService) LargerThan	(larger		string) *GmailService	{ gs.sLarger		= larger    ; return gs }
func (gs *GmailService) SmallerThan	(smaller	string) *GmailService	{ gs.sSmaller		= smaller   ; return gs }
func (gs *GmailService) Filename	(fname		string) *GmailService	{ gs.sFilename		= fname	    ; return gs }
func (gs *GmailService) HasAttachment   (hasatt		bool  ) *GmailService	{ gs.sHasAttachment     = hasatt    ; return gs }
func (gs *GmailService) Match		(match		string) *GmailService	{ gs.sMatch		= match	    ; return gs }
func (gs *GmailService) MatchExact	(matchex	string) *GmailService	{ gs.sMatchExact	= matchex   ; return gs }

func (gm *GmailMessage) parseMessagePart(origmsg *gmail.MessagePart, gs *GmailService) {
	var contentDisp string  // Content-Disposition
	gm.mimeFlow = append(gm.mimeFlow, origmsg.MimeType)
	
	for _, ii := range(origmsg.Headers) {
		if ii.Name == "Subject" { gm.sSubject = ii.Value }
		if ii.Name == "Content-Disposition" { contentDisp = ii.Value }
		if ii.Name == "Message-ID" { gm.sMessageId = ii.Value }
	}
	if strings.HasPrefix(origmsg.MimeType,"multipart") {
		for _, ii := range(origmsg.Parts) {
			gm.parseMessagePart(ii, gs)
		}
	}
	if strings.HasPrefix(contentDisp, "attachment") {
		var a = new(GmailAttachment)
		a.gmailService  = gs
		a.bDownloaded   = false
		a.sMessageId    = gm.sMessageId
		a.sAttachmentId = origmsg.Body.AttachmentId
		a.iSize         = origmsg.Body.Size
		a.sFilename     = strings.Trim(contentDisp[strings.Index(contentDisp, "filename=")+len("filename="):], "\"")
		a.sMimeType     = origmsg.MimeType
		gm.lAttachment  = append(gm.lAttachment, a)
	} else if origmsg.MimeType == "text/plain" {
		gm.sBodyText, _ = base64.URLEncoding.DecodeString(origmsg.Body.Data)
	} else if origmsg.MimeType == "text/html"  {
		gm.sBodyHtml, _ = base64.URLEncoding.DecodeString(origmsg.Body.Data)
	}
}

func (gs *GmailService) GetMessages() []*GmailMessage {
	messages  := gs.GetMessagesRaw()
	var gmessages []*GmailMessage
	for _, ii := range messages {
		var m = new(GmailMessage)
		m.parseMessagePart(ii.Payload, gs)
		gmessages = append(gmessages, m)
	}
	return gmessages
}

func (gs *GmailService) GetMessagesRaw() []*gmail.Message {
	var messages []*gmail.Message
	for _, ii := range (gs.GetListOnly().Messages) {
		m, err := gs.srv.Users.Messages.Get(gs.sUser, ii.Id).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve email messages %v", err)
		}
		messages = append(messages, m)
	}
	return messages
}

func (gs *GmailService) GetListOnly () *gmail.ListMessagesResponse {
	qStr   := ""
	gsCall := gs.srv.Users.Messages.List(gs.sUser).MaxResults(gs.iMaxResults)
	if len(gs.sLabel)           > 0 { gsCall = gsCall.LabelIds(gs.sLabel) }
	
	if len(gs.sFrom)            > 0 { qStr += "from:"		+ gs.sFrom		+ " "   }
	if len(gs.sTo)              > 0 { qStr += "to:"		        + gs.sTo		+ " "   }
	if len(gs.sOlder)           > 0 { qStr += "older:"		+ gs.sOlder		+ " "   }
	if len(gs.sNewer)           > 0 { qStr += "newer:"		+ gs.sNewer		+ " "   }
	if len(gs.sOlderRel)        > 0 { qStr += "older_than:"	        + gs.sOlderRel		+ " "   }
	if len(gs.sNewerRel)        > 0 { qStr += "newer_than:"	        + gs.sNewerRel		+ " "   }
	if len(gs.sSubject)         > 0 { qStr += "subject:\""		+ gs.sSubject		+ "\" " }
	if len(gs.sInPlace)         > 0 { qStr += "in:"			+ gs.sInPlace		+ " "	}
	if len(gs.sLarger)          > 0 { qStr += "larger:"		+ gs.sLarger		+ " "	}
	if len(gs.sSmaller)         > 0 { qStr += "smaller:"		+ gs.sSmaller		+ " "	}
	if len(gs.sFilename)        > 0 { qStr += "filename:\""		+ gs.sFilename		+ "\" " }
	if len(gs.sMatch)           > 0 { qStr += " "			+ gs.sMatch		+ " "	}
	if len(gs.sMatchExact)      > 0 { qStr += " "                   + gs.sMatchExact	+ " "	}
	if gs.sHasAttachment            { qStr += "has:attachment "					}
	if len(qStr) > 0 { gsCall = gsCall.Q(qStr) }
	do, err := gsCall.Do()
        if err != nil {
                log.Fatalf("Unable to retrieve email list %v", err)
        }
	return do
}

type GmailAttachment struct {
	gmailService  *GmailService
	sAttachmentId string
	sMessageId    string
	sFilename     string
	sMimeType     string
	sData         []byte
	iSize         int64
	bDownloaded   bool    
}

func (ga *GmailAttachment) GetFilename     ()  string  { return ga.sFilename     }
func (ga *GmailAttachment) GetMimeType     ()  string  { return ga.sMimeType     }
func (ga *GmailAttachment) GetSize         ()  int64   { return ga.iSize         }
func (ga *GmailAttachment) IsDownloaded    ()  bool    { return ga.bDownloaded   }
func (ga *GmailAttachment) GetAttachmentId ()  string  { return ga.sAttachmentId }
func (ga *GmailAttachment) GetMessageId    ()  string  { return ga.sMessageId    }

func (ga *GmailAttachment) GetData() []byte {
	if ga.bDownloaded { return ga.sData }
	att, err := ga.gmailService.srv.Users.Messages.Attachments.Get("me", ga.sMessageId, ga.sAttachmentId).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve email messages %v", err)
	}
	ga.sData, _ = base64.URLEncoding.DecodeString(att.Data)
	return ga.sData
}

type GmailMessage struct {
	sRaw             *gmail.Message
	sMessageId        string
	sSubject          string
	sBodyHtml         []byte
	sBodyText         []byte
	lAttachment       []*GmailAttachment
	mHeaders          map[string]string
	mimeFlow          []string  // multipart/mixed -> multipart/alternative etc
}


func (gm *GmailMessage) HasSubject	() bool			{ return len(gm.sSubject)	> 0	}
func (gm *GmailMessage) HasBodyHtml	() bool			{ return len(gm.sBodyHtml)	> 0	}
func (gm *GmailMessage) HasBodyText	() bool			{ return len(gm.sBodyText)	> 0	}
func (gm *GmailMessage) HasAttachments	() bool			{ return len(gm.lAttachment)	> 0	}
func (gm *GmailMessage) GetSubject	() string		{ return gm.sSubject			}
func (gm *GmailMessage) GetBodyText	() []byte		{ return gm.sBodyText			}
func (gm *GmailMessage) GetBodyHtml	() []byte		{ return gm.sBodyHtml			}
func (gm *GmailMessage) GetAttachments	() []*GmailAttachment	{ return gm.lAttachment			}
func (gm *GmailMessage) GetRawMessage	() *gmail.Message	{ return gm.sRaw			}
func (gm *GmailMessage) GetMessageId	() string               { return gm.sMessageId                  }
