package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gluk256/crypto/cmd/common"
	"github.com/gluk256/crypto/crutils"
	"github.com/gluk256/crypto/terminal"
)

func help() {
	fmt.Printf("xcry v.1.%d \n", crutils.CipherVersion)
	fmt.Println("encrypt/decrypt a file")
	fmt.Println("USAGE: xcry flags [srcFile] [dstFile]")
	fmt.Println("\t -h help")
	fmt.Println("\t -s secure password/text input")
	fmt.Println("\t -x extra secure password/text input")
	fmt.Println("\t -f save to file")

	fmt.Println("\t -e encrypt (default mode)")
	fmt.Println("\t\t -r random password")

	fmt.Println("\t -d decrypt")
	fmt.Println("\t\t -p output decrypted content as text, don't save")
	fmt.Println("\t\t -g interactive grep (print specific text lines only)")
	fmt.Println("\t\t -G interactive grep with secure input")
	fmt.Println("\t\t -m mute - don't print the decrypted content")

	fmt.Println("\t -l load file")
	fmt.Println("\t\t -i insert file content into another file as steganographic content")

	fmt.Println("\t -t enter text (password mode)")
	fmt.Println("\t -T enter text (plain text mode)")
}

func processCommandArgs() (flags string, srcFile string, dstFile string) {
	if len(os.Args) == 1 {
		flags = "h"
	}
	if len(os.Args) > 1 {
		flags = os.Args[1]
	}
	if len(os.Args) > 2 {
		srcFile = os.Args[2]
	}
	if len(os.Args) > 3 {
		dstFile = os.Args[3]
	}

	if strings.Contains(flags, "h") || strings.Contains(flags, "?") {
		help()
		return "", srcFile, dstFile
	}

	if strings.Contains(flags, "d") && strings.Contains(flags, "r") {
		fmt.Println("Random password ('r') is incompatible with decryption ('d').")
		fmt.Println("ERROR: wrong flags.")
		return "", srcFile, dstFile
	}

	if strings.Contains(flags, "t") || strings.Contains(flags, "T") {
		if len(os.Args) > 2 {
			fmt.Println("ERROR: flag -t is incompatible with param srcFile")
			return "", srcFile, dstFile
		}
	}

	return flags, srcFile, dstFile
}

func enterText(flags string) (res []byte) {
	fmt.Println("Please enter your text:")
	if strings.Contains(flags, "s") {
		res = terminal.SecureInput(false)
	} else if strings.Contains(flags, "x") {
		res = terminal.SecureInput(true)
	} else if strings.Contains(flags, "T") {
		res = terminal.PlainTextInput()
	} else {
		res = terminal.PasswordModeInput()
	}

	if strings.Contains(flags, "d") && common.IsHexData(res) {
		h := make([]byte, len(res)/2)
		_, err := hex.Decode(h, res)
		if err != nil {
			fmt.Printf("Error decoding hex data: %s\n", err.Error())
			return nil
		} else {
			res = h
			fmt.Println("Loaded data in hex format")
		}
	}
	return res
}

func getData(flags string, srcFile string) (data []byte) {
	if strings.Contains(flags, "t") || strings.Contains(flags, "T") {
		data = enterText(flags)
	} else {
		data = loadDataFromFile(srcFile)
	}
	if len(data) == 0 {
		fmt.Println("Error: empty data")
	}
	return data
}

func loadDataFromFile(filename string) []byte {
	for i := 0; i < 8; i++ {
		if len(filename) == 0 || i > 0 {
			filename = common.GetFileName()
			if len(filename) == 0 || filename == string("q") {
				return nil
			}
		}
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("Failed to load data: %s\n", err.Error())
			fmt.Println("Please try again.")
			continue
		}
		return data
	}

	fmt.Println("The number of iterations exceeded allowed maximum. Exit.")
	return nil
}

func main() {
	flags, srcFile, dstFile := processCommandArgs()
	if len(flags) == 0 {
		return
	}

	defer crutils.ProveDataDestruction()
	data := getData(flags, srcFile)
	if len(data) == 0 {
		return
	}

	if strings.Contains(flags, "d") {
		processDecryption(flags, dstFile, data, false)
	} else {
		processEncryption(flags, dstFile, data, nil)
	}

	crutils.AnnihilateData(data)
}

func decrypt(flags string, data []byte, unknownSize bool) (decrypted []byte, steg []byte, err error) {
	var key []byte
	defer crutils.AnnihilateData(key) // in case of panic

	for i := 0; i < 1024; i++ {
		key = common.GetPassword(flags)
		if unknownSize {
			decrypted, steg, err = crutils.DecryptStegContentOfUnknownSize(key, data)
		} else {
			decrypted, steg, err = crutils.Decrypt(key, data)
		}
		if err == nil {
			return decrypted, steg, err
		}
		crutils.AnnihilateData(key)
		fmt.Printf("Failed to decrypt data: %s\n", err.Error())
		fmt.Print("Do you want to retry? [y/n]: ")
		res := terminal.PlainTextInput()
		if len(res) > 0 && res[0] == 'n' {
			return nil, nil, err
		}
	}
	return decrypted, steg, err
}

func processDecryption(flags string, dstFile string, data []byte, unknownSize bool) {
	if len(data) == 0 {
		fmt.Println("Error: empty data")
		return
	}

	decrypted, steg, err := decrypt(flags, data, unknownSize)
	defer crutils.AnnihilateData(decrypted)
	defer crutils.AnnihilateData(steg)
	if err != nil {
		return
	}

	if !strings.Contains(flags, "m") {
		fmt.Printf("Decrypted:\n[%s]\n", string(decrypted))
	}

	if !strings.Contains(flags, "f") {
		for {
			fmt.Println("What do you want to do with decrypted content?")
			fmt.Println("options: save_File, Grep, Print, Decrypt, Secure_pass, eXtra_secure, Quit")
			fmt.Print("Please enter the command: ")
			flags = string(terminal.PlainTextInput())
			if strings.Contains(flags, "p") {
				fmt.Printf("Decrypted:\n[%s]\n", string(decrypted))
			} else if strings.Contains(flags, "q") {
				return
			} else {
				break
			}
		}
	}

	if strings.Contains(flags, "f") {
		common.SaveData(dstFile, decrypted)
	} else if strings.Contains(flags, "G") {
		runGrep(flags, decrypted)
	} else if strings.Contains(flags, "g") {
		runGrep(flags, decrypted)
	} else {
		processDecryption(flags, dstFile, steg, true) // recursively decrypt steg content
	}
}

func encrypt(key []byte, data []byte, steg []byte) (res []byte, err error) {
	if steg == nil {
		res, err = crutils.Encrypt(key, data)
	} else {
		res, err = crutils.EncryptSteg(key, data, steg)
	}
	return res, err
}

func processEncryption(flags string, dstFile string, data []byte, steg []byte) {
	key := common.GetPassword(flags)
	defer crutils.AnnihilateData(key) // in case of panic

	encrypted, err := encrypt(key, data, steg)
	defer crutils.AnnihilateData(encrypted)
	if err != nil {
		fmt.Printf("Failed to encrypt data: %s\nFATAL ERROR\n", err.Error())
		return
	}
	crutils.AnnihilateData(key)

	if !strings.Contains(flags, "f") {
		fmt.Println("What do you want to do with encrypted content?")
		fmt.Println("options: save_File, Encrypt, Rand_pass, Secure_pass, eXtra_secure, Quit")
		fmt.Print("Please enter the command: ")
		flags = string(terminal.PlainTextInput())
	}

	if strings.Contains(flags, "f") {
		err = common.SaveData(dstFile, encrypted)
	} else if strings.Contains(flags, "q") {
		return
	} else {
		for i := 0; i < 3; i++ {
			buf := getData(flags, "")
			if len(buf) == 0 {
				fmt.Println("no data, please try again")
			} else if len(buf) < len(encrypted)+4 {
				fmt.Printf("File size in insufficiant for steg encryption [%d vs. %d]. Please try again.\n", len(buf), len(encrypted)+4)
			} else {
				processEncryption(flags, dstFile, buf, encrypted) // recursively encrypt steg content
				break
			}
		}
	}
}

func runGrep(flags string, content []byte) {
	lines := strings.Split(string(content), "\n")
	var s []byte
	secure := strings.Contains(flags, "G")
	for {
		fmt.Print("please enter text: ")
		if secure {
			s = terminal.SecureInput(false)
		} else {
			s = terminal.PasswordModeInput()
		}
		if string(s) == "q" {
			break
		}
		found := false
		for _, ln := range lines {
			if strings.Contains(ln, string(s)) {
				fmt.Println(ln)
				found = true
			}
		}
		if !found {
			fmt.Println(">>> Requested text not found")
		}
	}
}
