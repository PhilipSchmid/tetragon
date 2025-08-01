// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package sensors

import (
	"fmt"

	"github.com/cilium/tetragon/pkg/kernels"
	"github.com/cilium/tetragon/pkg/logger"
	"github.com/cilium/tetragon/pkg/option"
	"github.com/cilium/tetragon/pkg/sensors/program"
)

func observerLoadInstance(bpfDir string, load *program.Program, maps []*program.Map) error {
	version, _, err := kernels.GetKernelVersion(option.Config.KernelVersion, option.Config.ProcFS)
	if err != nil {
		return err
	}

	l := logger.GetLogger()
	l.Debug(fmt.Sprintf("observerLoadInstance %s %d", load.Name, version), "prog", load.Name, "kern_version", version)

	err = loadInstance(bpfDir, load, maps, version, option.Config.Verbosity)
	if err != nil && load.ErrorFatal {
		return fmt.Errorf("failed prog %s kern_version %d loadInstance: %w",
			load.Name, version, err)
	}
	return nil
}

func loadInstance(bpfDir string, load *program.Program, maps []*program.Map, _, verbose int) error {
	// Check if the load.type is a standard program type. If so, use the standard loader.
	loadFn, ok := standardTypes[load.Type]
	if ok {
		logger.GetLogger().Debug("Loading BPF program", "Program", load.Name, "Type", load.Type, "Attach", load.Attach)
		return loadFn(bpfDir, load, maps, verbose)
	}

	return fmt.Errorf("program %s has unregistered type '%s'", load.Label, load.Type)
}

func observerMinReqs() (bool, error) {
	return true, nil
}

func flushKernelSpec() {}

func (s *Sensor) preLoadMaps(_ string, _ []*program.Map) ([]*program.Map, error) {
	return nil, nil
}

func getCachedBTFFile() string {
	return ""
}
