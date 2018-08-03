package common

import (
	"encoding/json"
	"fmt"
	"os"
	"github.com/fatih/color"
	"github.com/cruisechang/util/debug"
)

const (
	//IsPrint defermin if prints
	isPrint        = false
	isPrintInfo    = true
	isPrintWarning = true
	isPrintError   = true

	//ColorCyan text color
	ColorCyan = 0
	//ColorGreen text color
	ColorGreen = 2
	//ColorYellow text color
	ColorYellow = 3

	//LogInfo log info level
	LogInfo = 0
	//LogWarning log warring level
	LogWarning = 1
	//LogError log error level
	LogError = 2
)

//PrintClass turns class to json string
func PrintClass(a interface{}) {
	j, _ := json.Marshal(a)
	fmt.Println(string(j))
}

//Print print info by IsPrint
func Print(function interface{}, a ...interface{}) {
	if !isPrint {
		return
	}
	debug.LogInfo(os.Stdout, debug.GetShortPackageAndFunctionName(function)+":", a)
}

//PrintColor print info by IsPrint
func PrintColor(co int, function interface{}, v ...interface{}) {
	if !isPrint {
		return
	}
	var c *color.Color

	switch co {
	case ColorYellow:
		c = color.New(color.FgYellow)
	case ColorGreen:
		c = color.New(color.FgGreen)
	case ColorCyan:
		c = color.New(color.FgCyan)
	default:
		c = color.New(color.FgWhite)
	}

	c.EnableColor()
	c.Println(debug.GetShortPackageAndFunctionName(function)+":", v)
	c.DisableColor()
}

//PrintInfo print info by IsPrint
func PrintInfo(function interface{}, v ...interface{}) {
	if !isPrintInfo {
		return
	}
	c := color.New(color.FgBlue)
	c.EnableColor()
	c.Println("[Info]"+debug.GetShortPackageAndFunctionName(function)+":", v)
	c.DisableColor()
}

//PrintWarning print info by IsPrint
func PrintWarning(function interface{}, v ...interface{}) {
	if !isPrintWarning {
		return
	}
	c := color.New(color.FgYellow)
	c.EnableColor()
	c.Println("[Info]"+debug.GetShortPackageAndFunctionName(function)+":", v)
	c.DisableColor()
}

//PrintError prints err msg to stdout.
func PrintError(function interface{}, v ...interface{}) {
	if !isPrintError {
		return
	}

	c := color.New(color.FgRed)
	c.EnableColor()
	c.Println("[Error]"+debug.GetShortPackageAndFunctionName(function)+":", v)
	c.DisableColor()
	//debug.LogError(os.Stdout, debug.GetShortPackageAndFunctionName(function)+":", a)
}

//PrintFatal prints message and followed bt exit calling.
func PrintFatal(function interface{}, v ...interface{}) {
	c := color.New(color.FgHiMagenta)
	c.EnableColor()
	c.Println("[Fatal]"+debug.GetShortPackageAndFunctionName(function)+":", v)
	c.DisableColor()
	//debug.LogFatal(os.Stdout, debug.GetShortPackageAndFunctionName(function)+":", v)
}


/*
//LogGame log file info for game client.
func LogGame(targetUserID int, v ...interface{}) {
	filePath := nex.config.Server.Log.Path + "/" + config.Instance.Server.Log.GamePrefix + time.Now().Format("20060102") + config.Instance.Server.Log.Extension

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf(debug.GetShortPackageAndFunctionName(LogFileWarning)+" error opening file: %v \n", err)
	} else {
		debug.LogInfo(f, "["+strconv.Itoa(targetUserID)+"]", v)
	}
	defer f.Close()

}

//LogFile log file
func LogFile(logLevel int, filePath string, targetUserID int, v ...interface{}) {

	switch logLevel {
	case LogInfo:
		LogFileInfo(targetUserID, v)
		break
	case LogWarning:
		LogFileWarning(targetUserID, v)
		break
	case LogError:
		LogFileError(targetUserID, v)
		break
	default:
		LogFileInfo(targetUserID, v)
	}
}

//LogFileInfo log file info
func LogFileInfo(targetUserID int, v ...interface{}) {
	if config.Instance == nil || config.Instance.Server == nil || config.Instance.Server.Log == nil {
		return
	}
	filePath := config.Instance.Server.Log.Path + "/" + config.Instance.Server.Log.ServerPrefix + time.Now().Format("20060102") + config.Instance.Server.Log.Extension

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf(debug.GetShortPackageAndFunctionName(LogFileWarning)+" error opening file: %v \n", err)
	} else {
		debug.LogInfo(f, "["+strconv.Itoa(targetUserID)+"]", v)
	}
	defer f.Close()

}

//LogFileWarning log file warning
func LogFileWarning(targetUserID int, v ...interface{}) {
	if config.Instance == nil || config.Instance.Server == nil || config.Instance.Server.Log == nil {
		return
	}

	filePath := config.Instance.Server.Log.Path + "/" + config.Instance.Server.Log.ServerPrefix + time.Now().Format("20060102") + config.Instance.Server.Log.Extension

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf(debug.GetShortPackageAndFunctionName(LogFileWarning)+" error opening file: %v \n", err)
	} else {
		debug.LogWarning(f, "["+strconv.Itoa(targetUserID)+"]", v)
	}
	defer f.Close()
}

//LogFileError log file error
func LogFileError(targetUserID int, v ...interface{}) {
	if config.Instance == nil || config.Instance.Server == nil || config.Instance.Server.Log == nil {
		return
	}

	filePath := config.Instance.Server.Log.Path + "/" + config.Instance.Server.Log.ServerPrefix + time.Now().Format("20060102") + config.Instance.Server.Log.Extension

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf(debug.GetShortPackageAndFunctionName(LogFileError)+" error opening file: %v \n", err)
	} else {
		debug.LogError(f, "["+strconv.Itoa(targetUserID)+"]", v)
	}
	defer f.Close()
}

//MakeLogFile make a log file according to serve config file.
func MakeLogFile() {
	if config.Instance == nil || config.Instance.Server == nil || config.Instance.Server.Log == nil {
		return
	}
	ti := time.Now()
	path := config.Instance.Server.Log.Path + "/" + config.Instance.Server.Log.ServerPrefix + ti.Format("20060102") + config.Instance.Server.Log.Extension

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	defer f.Close()
	if err != nil {
		f, err = os.Create(path)
		if err != nil {

			fmt.Printf("error opening file: %v \n", err)
			return
		}
	}
}

func getFilePath() string {
	if config.Instance == nil || config.Instance.Server == nil || config.Instance.Server.Log == nil {
		return ""

	}
	return config.Instance.Server.Log.Path + "/" + config.Instance.Server.Log.ServerPrefix + time.Now().Format("20060102") + config.Instance.Server.Log.Extension
}
*/
