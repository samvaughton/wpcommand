package types

type VersionInfo struct {
	Name          string `json:"name"`
	BuildDate     string `json:"buildDate"`
	Revision      string `json:"revision"`
	Version       string `json:"version"`
	APIVersion    string `json:"apiVersion"`
	GoVersion     string `json:"goVersion"`
	OS            string `json:"os"`
	Arch          string `json:"arch"`
	KernelVersion string `json:"kernelVersion"`
	Experimental  bool   `json:"experimental"`
}
