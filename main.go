package main

import "flag"
import "github.com/go-ini/ini"
import "strings"
import (
	"os/user"
	"io/ioutil"
	"regexp"
	"os"
	"log"
	"fmt"
)

func main() {
	// Get the session name from the CLI
	var quiet = flag.Bool("q", false, "Quiet mode, avoid all output.")
	var verbose = flag.Bool("v", false, "Verbose mode, prints IDs of processes running in the given schroot session.")
	flag.Usage = func() {
		if !*quiet {
			fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... SCHROOT-SESSION-NAME\n", os.Args[0])
			fmt.Fprint(os.Stderr, "Options:\n")
			flag.PrintDefaults()
		}
	}
	flag.Parse()
	if *quiet {
		if *verbose {
			log.Fatalln("ERR: -q and -v are mutual exclusive.")
		}
		log.SetOutput(ioutil.Discard)
	}
	var sessionName = flag.Arg(0)
	if (sessionName == "") {
		flag.Usage()
		os.Exit(1)
	}

	matches, err := regexp.MatchString(`[/\\:,;~&()'">< ]`, sessionName)
	if err != nil {
		log.Println(err)
		log.Fatalln("ERR: Could not validate session name.")
	}
	if matches {
		log.Println(err)
		log.Fatalln("ERR: Session name does not pass validation.")
	}

	current_user, user_err := user.Current()
	if user_err != nil {
		log.Println(err)
		log.Fatalln("ERR: Could not retrieve the current user.")
	}
	// find original chroot
	schrootName, err := getSchrootName(sessionName)

	if err != nil {
		log.Println(err)
		log.Fatalln("ERR: Could not find the schroot or the session.")
	}
	allowedUsers, err := getAllowedUsers(schrootName)
	if err != nil {
		log.Println(err)
		log.Fatalln("ERR: Could not retrieve the list of users allowed to access the schroot.")
	}

	if (!isValueInCommaSeparatedList(current_user.Username, allowedUsers)) {
		log.Println(err)
		log.Fatalln("ERR: You are not an allowed user for the parent schroot of the given session.")
	}

	result, err := getAllProcessIdsInSchrootSessionDir("/var/lib/schroot/mount/" + sessionName, !*verbose)
	if err != nil {
		log.Println(err)
		log.Fatalln("ERR: Could not read all processes.")
	}

	if *verbose {
		message := "INFO: The following process IDs are running in the given session:"
		for _, id := range result {
			message += " " + id
		}
		log.Println(message)
	}

	if len(result) > 0 {
		log.Println("RESULT: There is at least one process active for the given session.")
		os.Exit(-1)
	} else {
		log.Println("RESULT: There is no process active for the given session.")
		os.Exit(0)
	}
}

func getAllProcessIdsInSchrootSessionDir(sessionMountDir string, earlyReturn bool) ([]string, error) {
	files, err := ioutil.ReadDir("/proc");
	processIdsInSession := make([]string, 0)
	if err != nil {
		return processIdsInSession, err
	}

	for _, f := range files {
		if f.IsDir() {
			isNumber, err := regexp.MatchString("[0-9]+", f.Name())
			if err != nil {
				return processIdsInSession, err
			}
			if isNumber {
				root, err := os.Readlink("/proc/" + f.Name() + "/root")
				if err != nil {
					return processIdsInSession, err
				}
				if root == sessionMountDir {
					processIdsInSession = append(processIdsInSession, f.Name())
					if earlyReturn {
						return processIdsInSession, nil
					}
				}
			}
		}
	}
	return processIdsInSession, nil
}

func getKeyFromIniFile(fileName string, sectionName string, keyName string) (string, error) {
	cfg, err := ini.Load(fileName)
	if err != nil {
		return "", err;
	}
	section, section_err := cfg.GetSection(sectionName)
	if section_err != nil {
		return "", err;
	}
	key, key_err := section.GetKey(keyName)
	if key_err != nil {
		return "", err;
	}
	return key.String(), nil
}

func getSchrootName(sessionName string) (string, error) {
	filename := "/var/lib/schroot/session/" + sessionName
	return getKeyFromIniFile(filename, sessionName, "original-name")
}

func getAllowedUsers(schrootName string) (string, error) {
	filename := "/etc/schroot/chroot.d/" + schrootName + ".conf"
	return getKeyFromIniFile(filename, schrootName, "users")
}

func isValueInCommaSeparatedList(value string, commaSeparatedList string) bool {
	allValues := strings.Split(commaSeparatedList, ",")
	for i := range allValues {
		if strings.Trim(allValues[i], " ") == value {
			return true
		}
	}
	return false
}
