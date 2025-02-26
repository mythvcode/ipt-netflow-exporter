package statparser

import (
	"bytes"
	"errors"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/mythvcode/ipt-netflow-exporter/internal/logger"
)

var readFile = os.ReadFile

var (
	isCPUStat    = regexp.MustCompile(`^cpu\d+$`).MatchString
	isSocketStat = regexp.MustCompile(`^sock\d+$`).MatchString
)

type StatCollector struct {
	filepath string
	log      *logger.Logger
}

func New(statPath string) *StatCollector {
	return &StatCollector{
		filepath: statPath,
		log:      logger.GetLogger().With(slog.String(logger.Component, "StatCollector")),
	}
}

func readStatFile(filePath string) ([]string, error) {
	result := make([]string, 0, 30)
	fileContent, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	for _, line := range bytes.Split(fileContent, []byte("\n")) {
		result = append(result, string(line))
	}

	return result, nil
}

func (s *StatCollector) CollectAndMarshal() (Statistics, error) {
	file, err := readStatFile(s.filepath)
	if err != nil {
		return Statistics{}, err
	}

	return s.parseFields(file)
}

func (s *StatCollector) parseFields(fileLines []string) (Statistics, error) {
	resultStruct := Statistics{}
	for _, line := range fileLines {
		if splitLine := strings.Fields(line); len(splitLine) > 0 {
			if err := s.processConfigLine(&resultStruct, splitLine); err != nil {
				return resultStruct, err
			}
		}
	}

	return resultStruct, nil
}

func (s *StatCollector) processConfigLine(statStruct *Statistics, splitLine []string) error {
	if isCPUStat(splitLine[0]) {
		if res := s.parseCPUFields(splitLine); res != nil {
			statStruct.CPUStatList = append(statStruct.CPUStatList, *res)
		}
		// do not return errors for specific metrics
		return nil
	}
	if isSocketStat(splitLine[0]) {
		if res := s.parseSocketFields(splitLine); res != nil {
			statStruct.SockStatList = append(statStruct.SockStatList, *res)
		}
		// do not return errors for specific metrics
		return nil
	}

	if err := setValueByName(statStruct, splitLine[0], splitLine[1]); err != nil {
		if errors.Is(errNotFoundField, err) {
			s.log.Debugf("found unsupported metrics in ipt_NETFLOW stat file: metric %s", splitLine[0])
		} else {
			return err
		}
	}

	return nil
}

func (s *StatCollector) parseCPUFields(cpuFields []string) *CPUStat {
	result := CPUStat{}
	if err := setValues(&result, cpuFields); err != nil {
		s.log.Errorf("%s", err.Error())

		return nil
	}

	return &result
}

func (s *StatCollector) parseSocketFields(cpuFields []string) *NFSockEntry {
	result := NFSockEntry{}
	if err := setValues(&result, cpuFields); err != nil {
		s.log.Errorf("%s", err.Error())

		return nil
	}

	return &result
}
