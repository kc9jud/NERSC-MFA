package main

import (
	"bufio"
	"fmt"
	"os"

	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/howeyc/gopass"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableSorting: true,
	})
}

func main() {
	app := cli.NewApp()
	app.Name = "sshproxy"
	app.Usage = "sshproxy grabs keys from NERSC"
	app.Before = func(ctx *cli.Context) error {
		if ctx.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}
	app.Version = "1.0"
	app.Author = "Matt Dunford"
	app.Email = "mdunford@lbl.gov"
	app.Action = grabPpk
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Debug logging",
		},
		cli.StringFlag{
			Name:  "user, u",
			Usage: "Specify remote (NERSC) username",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Specify pathname for private key",
			Value: "nersckey.ppk",
		},
		cli.StringFlag{
			Name:  "scope, s",
			Usage: "key scope",
			Value: "default",
		},
		cli.StringFlag{
			Name:  "collab, c",
			Usage: "Specify a collaboration account",
		},
		cli.StringFlag{
			Name:  "server, U",
			Usage: "server to grab keys from",
			Value: "sshproxy.nersc.gov",
		},
	}
	app.Run(os.Args)
}

func promptPassword(prompt string) (password string, err error) {
	fmt.Print(prompt)
	bytes, err := gopass.GetPasswdMasked()
	password = string(bytes)
	return
}

func grabPpk(ctx *cli.Context) error {

	username := ctx.GlobalString("user")
	if len(username) < 1 {
		logrus.Fatal("-u username is required")
	}

	fname := ctx.GlobalString("output")
	if len(fname) < 1 {
		logrus.Fatal("no file name specified")
	}

	scope := ctx.GlobalString("scope")
	if len(scope) < 1 {
		logrus.Fatal("no scope specified")
	}

	server := ctx.GlobalString("server")
	if len(server) < 1 {
		logrus.Fatal("no scope specified")
	}

	prompt := fmt.Sprintf("Enter the password+OTP for %s: ", username)
	password, err := promptPassword(prompt)
	if err != nil {
		logrus.Fatal(err)
	}

	url := fmt.Sprintf("https://%s/create_pair/%s/?putty", server, scope)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	req.SetBasicAuth(username, password)

	logrus.Debugf("contacting %s", url)
	resp, err := client.Do(req)
	if err != nil {
		logrus.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logrus.Fatalf("%s returned status code %d: %s", url, resp.StatusCode, resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatal(err)
	}

	f, err := os.Create(fname)
	if err != nil {
		logrus.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	w.Write(body)
	w.Flush()

	fmt.Printf(`key was written to %s.  Run "pageant %s" to load the key. Then run putty instances like this: putty -agent %s@cori.nersc.gov`,
		fname, fname, username)
	return nil
}
