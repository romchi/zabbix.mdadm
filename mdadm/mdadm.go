package mdadm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// MdadmDeviceStruct struct
type MdadmDeviceStruct struct {
	Name string
}

// RaidStats struct
type RaidStats struct {
	Capability       int64 // capability
	Dev              string
	DiscardAlignment int64 // discard_alignment
	ExtRange         int64 // ext_range
	Range            int64 // range
	Removable        int64 // removable
	RO               int64 // ro
	Size             int64 // size
	MD               *MDStats
}

// MDStats struct
type MDStats struct {
	Level                  string // md/level
	ArrayState             string // md/array_state
	Degraded               int64  // md/degraded
	MaxReadErrors          int64  // md/max_read_errors
	MetadataVersion        string // md/metadata_version
	MismatchCnt            int64  // md/mismatch_cnt
	PrereadBypassThreshold int64  // preread_bypass_threshold
	RaidDisks              int64  // md/raid_disks
	SyncAction             string // md/sync_action
}

var (
	systemPath = "/sys/block"
)

func valueString(file string, addr *string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	d, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	*addr = strings.TrimSpace(string(d))
	return nil
}

func valueInt64(file string, addr *int64) error {
	var s string
	if err := valueString(file, &s); err != nil {
		return err
	}

	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*addr = v

	return nil
}

func StatsRaidController(rName string) {
	raids, err := MdadmDevices()
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range raids {
		if rName == c.Name {
			i, err := c.RaidStats()
			if err != nil {
				log.Fatal(err)
			}
			m, err := c.MDStats()
			if err != nil {
				log.Fatal(err)
			}
			i.MD = m
			ji, _ := json.Marshal(i)
			fmt.Print(string(ji))
		}
	}
}

// MdadmDevicesList scan folder /sys/block and find md[] devices
func MdadmDevicesList() ([]string, error) {
	mdNames := "md[0-9]**"

	files, err := filepath.Glob(filepath.Join(systemPath, mdNames))

	if _, err = os.Stat(systemPath); os.IsNotExist(err) {
		return []string{}, err
	}

	if err != nil {
		return []string{}, err
	}

	result := make([]string, len(files))
	for ctr, name := range files {
		result[ctr] = filepath.Base(name)
	}

	return result, nil
}

// MdadmDevices s
func MdadmDevices() ([]MdadmDeviceStruct, error) {
	mdadmList, err := MdadmDevicesList()
	if err != nil {
		return nil, err
	}

	result := make([]MdadmDeviceStruct, len(mdadmList))
	for ctr, num := range mdadmList {
		result[ctr] = MdadmDeviceStruct{Name: filepath.Base(num)}
	}

	return result, nil
}

// RaidStats s
func (md MdadmDeviceStruct) RaidStats() (*RaidStats, error) {
	dir := filepath.Join(systemPath, md.Name)

	f, err := os.Open(dir)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	var files []string
	if files, err = f.Readdirnames(0); err != nil {
		return nil, err
	}

	raidStats := &RaidStats{}

	for _, f := range files {
		p := filepath.Join(dir, f)
		switch f {
		case "capability":
			if err = valueInt64(p, &raidStats.Capability); err != nil {
				return nil, err
			}
		case "dev":
			if err = valueString(p, &raidStats.Dev); err != nil {
				return nil, err
			}
		case "discard_alignment":
			if err = valueInt64(p, &raidStats.DiscardAlignment); err != nil {
				return nil, err
			}
		case "ext_range":
			if err = valueInt64(p, &raidStats.ExtRange); err != nil {
				return nil, err
			}
		case "range":
			if err = valueInt64(p, &raidStats.Range); err != nil {
				return nil, err
			}
		case "removable":
			if err = valueInt64(p, &raidStats.Removable); err != nil {
				return nil, err
			}
		case "ro":
			if err = valueInt64(p, &raidStats.RO); err != nil {
				return nil, err
			}
		case "size":
			if err = valueInt64(p, &raidStats.Size); err != nil {
				return nil, err
			}
		}
	}
	return raidStats, nil
}

// MDStats s
func (md MdadmDeviceStruct) MDStats() (*MDStats, error) {
	mdDir := md.Name + "/md"
	dir := filepath.Join(systemPath, mdDir)

	f, err := os.Open(dir)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	var files []string
	if files, err = f.Readdirnames(0); err != nil {
		return nil, err
	}

	mdStats := &MDStats{}

	for _, f := range files {
		p := filepath.Join(dir, f)
		switch f {
		case "level":
			if err = valueString(p, &mdStats.Level); err != nil {
				return nil, err
			}
		case "array_state":
			if err = valueString(p, &mdStats.ArrayState); err != nil {
				return nil, err
			}
		case "degraded":
			if err = valueInt64(p, &mdStats.Degraded); err != nil {
				return nil, err
			}
		case "max_read_errors":
			if err = valueInt64(p, &mdStats.MaxReadErrors); err != nil {
				return nil, err
			}
		case "mismatch_cnt":
			if err = valueInt64(p, &mdStats.MismatchCnt); err != nil {
				return nil, err
			}
		case "preread_bypass_threshold":
			if err = valueInt64(p, &mdStats.PrereadBypassThreshold); err != nil {
				return nil, err
			}
		case "raid_disks":
			if err = valueInt64(p, &mdStats.RaidDisks); err != nil {
				return nil, err
			}
		case "sync_action":
			if err = valueString(p, &mdStats.SyncAction); err != nil {
				return nil, err
			}
		}
	}
	return mdStats, nil
}
