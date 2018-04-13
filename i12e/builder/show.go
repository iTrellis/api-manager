package builder

import (
	"fmt"
	"strings"

	"github.com/dimiro1/banner"
	colorable "github.com/mattn/go-colorable"
)

// Show 显示项目信息
func Show(programVersion, compilerVersion, buildTime, author string) {
	bannerLogo :=
		`************************************************************
************************************************************
***   _______  ______  ______  _      _      _  ______   ***
***  / _____ \(  __  )|  __  \/ \    / \    / \/  __  \  ***
***  \/ | | \/| |  | || /  \_/| |    | |    | || (  \_/  ***
***     | |   | |__| || \___  | |    | |    | || (____   ***
***     | |   | |   _/|  ___) | |    | |    | |(____  )  ***
***     | |   | |\ \  | /   _ | |    | |    | |     ) |  ***
***     | |   | | \ \ | \__/ \| (___ | (___ | |/ \__) |  ***
***     \_/   \_/  \_)|______/(_____)(_____)\_/\______/  ***
***                                                      ***
************************************************************
****************** Compile Environment *********************
*** Program version : %s
*** Compiler version : %s
*** Build time : %s
*** Author : %s
************************************************************
****************** Running Environment *********************
*** Go running version : {{ .GoVersion }}
*** Go running OS : {{ .GOOS }} {{ .GOARCH }}
*** Startup time : {{ .Now "2006-01-02 15:04:05" }}
************************************************************
************************************************************
`
	newBanner := fmt.Sprintf(bannerLogo, programVersion, compilerVersion, buildTime, author)

	banner.Init(colorable.NewColorableStdout(), true, true, strings.NewReader(newBanner))
}
