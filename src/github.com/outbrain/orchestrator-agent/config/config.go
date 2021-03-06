/*
   Copyright 2014 Outbrain Inc.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package config

import (
	"encoding/json"
	"os"

	"github.com/outbrain/golib/log"
)

// Configuration makes for orchestrator-agent configuration input, which can be provided by user via JSON formatted file.
type Configuration struct {
	SnapshotMountPoint                 string // The single, agreed-upon mountpoint for logical volume snapshots
	ContinuousPollSeconds              uint   // Poll interval for continuous operation
	ResubmitAgentIntervalMinutes       uint   // Poll interval for resubmitting this agent on orchestrator agents API
	CreateSnapshotCommand              string // Command which creates a snapshot logical volume. It's a "do it yourself" implementation
	AvailableLocalSnapshotHostsCommand string // Command which returns list of hosts (one host per line) with available snapshots in local datacenter
	AvailableSnapshotHostsCommand      string // Command which returns list of hosts (one host per line) with available snapshots in any datacenter
	SnapshotVolumesFilter              string // text pattern filtering agent logical volumes that are valid snapshots
	MySQLDatadirCommand                string // command expected to present with @@datadir
	MySQLPortCommand                   string // command expected to present with @@port
	MySQLDeleteDatadirContentCommand   string // command which deletes all content from MySQL datadir (does not remvoe directory itself)
	MySQLServiceStopCommand            string // Command to stop mysql, e.g. /etc/init.d/mysql stop
	MySQLServiceStartCommand           string // Command to start mysql, e.g. /etc/init.d/mysql start
	MySQLServiceStatusCommand          string // Command to check mysql status. Expects 0 return value when running, non-zero when not running, e.g. /etc/init.d/mysql status
	ReceiveSeedDataCommand             string // Accepts incoming data (e.g. tarball over netcat)
	SendSeedDataCommand                string // Sends date to remote host (e.g. tarball via netcat)
	PostCopyCommand                    string // command that is executed after seed is done and before MySQL starts
	AgentsServer                       string // HTTP address of the orchestrator agents server
	HTTPPort                           uint   // HTTP port on which this service listens
	HTTPAuthUser                       string // Username for HTTP Basic authentication (blank disables authentication)
	HTTPAuthPassword                   string // Password for HTTP Basic authentication
	UseSSL                             bool   // If true, service will serve HTTPS only
	SSLSkipVerify                      bool   // When using SSL, should we ignore SSL certification error
	SSLPrivateKeyFile                  string // Name of SSL private key file, applies only when UseSSL = true
	SSLCertFile                        string // Name of SSL certification file, applies only when UseSSL = true
	HttpTimeoutSeconds                 int    // Number of idle seconds before HTTP GET request times out (when accessing orchestrator)
}

var Config *Configuration = NewConfiguration()

func NewConfiguration() *Configuration {
	return &Configuration{
		SnapshotMountPoint:                 "",
		ContinuousPollSeconds:              60,
		ResubmitAgentIntervalMinutes:       60,
		CreateSnapshotCommand:              "",
		AvailableLocalSnapshotHostsCommand: "",
		AvailableSnapshotHostsCommand:      "",
		SnapshotVolumesFilter:              "",
		MySQLDatadirCommand:                "",
		MySQLPortCommand:                   "",
		MySQLDeleteDatadirContentCommand:   "",
		MySQLServiceStopCommand:            "",
		MySQLServiceStartCommand:           "",
		MySQLServiceStatusCommand:          "",
		ReceiveSeedDataCommand:             "",
		SendSeedDataCommand:                "",
		PostCopyCommand:                    "",
		AgentsServer:                       "",
		HTTPPort:                           3002,
		HTTPAuthUser:                       "",
		HTTPAuthPassword:                   "",
		UseSSL:                             false,
		SSLSkipVerify:                      false,
		SSLPrivateKeyFile:                  "",
		SSLCertFile:                        "",
		HttpTimeoutSeconds:                 10,
	}
}

// read reads configuration from given file, or silently skips if the file does not exist.
// If the file does exist, then it is expected to be in valid JSON format or the function bails out.
func read(file_name string) (*Configuration, error) {
	file, err := os.Open(file_name)
	if err == nil {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(Config)
		if err == nil {
			log.Infof("Read config: %s", file_name)
		} else {
			log.Fatal("Cannot read config file:", file_name, err)
		}
	}
	return Config, err
}

// Read reads configuration from zero, either, some or all given files, in order of input.
// A file can override configuration provided in previous file.
func Read(file_names ...string) *Configuration {
	for _, file_name := range file_names {
		read(file_name)
	}
	return Config
}

// ForceRead reads configuration from given file name or bails out if it fails
func ForceRead(file_name string) *Configuration {
	_, err := read(file_name)
	if err != nil {
		log.Fatal("Cannot read config file:", file_name, err)
	}
	return Config
}
