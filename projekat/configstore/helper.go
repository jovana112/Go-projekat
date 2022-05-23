package configstore

import (
	"fmt"
	"github.com/google/uuid"
	"sort"
	"strings"
)

const (
	configs         = "configs"
	configId        = "configs/%s"
	configIdVersion = "configs/%s/%s"

	groups              = "groups"
	groupId             = "groups/%s"
	groupIdVersion      = "groups/%s/%s"
	groupIdVersionLabel = "groups/%s/%s/%s"
)

func createConfigsKey() string {
	return configs
}

func createNewConfigWithVersionKey(version string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(configIdVersion, id, version), id
}

func createConfigWithIdAndVersionKey(id string, version string) string {
	return fmt.Sprintf(configIdVersion, id, version)
}

func createConfigWithIdKey(id string) string {
	return fmt.Sprintf(configId, id)
}

func createGroupsKey() string {
	return groups
}

func createNewGroupWithVersionKey(version string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(groupIdVersion, id, version), id
}

func createGroupWithIdAndVersionKey(id string, version string) string {
	return fmt.Sprintf(groupIdVersion, id, version)
}
func createGroupWithIdAndVersionAndLabelKey(id string, version string, entries map[string]string) string {
	pairs := make([]string, 0, len(entries))
	for key := range entries {
		pairs = append(pairs, key+":"+entries[key])
	}
	sort.Strings(pairs)
	return fmt.Sprintf(groupIdVersionLabel, id, version, strings.Join(pairs, ";"))
}

func createGroupWithIdKey(id string) string {
	return fmt.Sprintf(groupId, id)
}
