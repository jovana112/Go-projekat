package configstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"os"
	"strings"
)

type ConfigStore struct {
	cli *api.Client
}

func CreateNewConfigStore() (*ConfigStore, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigStore{
		cli: client,
	}, nil
}

func (cs *ConfigStore) GetConfig(id string, version string) (*Config, error) {
	kv := cs.cli.KV()
	key := createConfigWithIdAndVersionKey(id, version)
	data, _, err := kv.Get(key, nil)

	if err != nil || data == nil {
		return nil, errors.New("That item does not exist!")
	}

	config := &Config{}
	err = json.Unmarshal(data.Value, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (cs *ConfigStore) DeleteConfig(id, version string) (map[string]string, error) {
	kv := cs.cli.KV()
	_, err := kv.Delete(createConfigWithIdAndVersionKey(id, version), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted config": id + version}, nil
}

func (cs *ConfigStore) GetAllConfig() ([]*Config, error) {
	kv := cs.cli.KV()
	data, _, err := kv.List(configs, nil)
	if err != nil {
		return nil, err
	}

	configs := []*Config{}
	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func (cs *ConfigStore) CreateConfig(config *Config) (*Config, error) {
	kv := cs.cli.KV()

	sid, rid := createNewConfigWithVersionKey(config.Version)
	config.Id = rid

	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	c := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (cs *ConfigStore) UpdateConfigWithNewVersion(config *Config, id string) (*Config, error) {
	kv := cs.cli.KV()

	_, er := cs.GetConfig(id, config.Version)

	if er == nil {
		return nil, errors.New("Config with that version already exists ")
	}

	sid := createConfigWithIdAndVersionKey(id, config.Version)
	config.Id = id

	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	c := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		return nil, err
	}
	return config, nil

}

func (cs *ConfigStore) CreateGroup(group *Group) (*Group, error) {
	kv := cs.cli.KV()
	sid, rid := createNewGroupWithVersionKey(group.Version)
	group.Id = rid

	data, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	g := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(g, nil)
	if err != nil {
		return nil, err
	}
	cs.CreateLabelsForGroupConfigs(group)

	return group, nil
}

func (cs *ConfigStore) GetGroup(id string, ver string) (*Group, error) {
	kv := cs.cli.KV()
	key := createGroupWithIdAndVersionKey(id, ver)
	data, _, err := kv.Get(key, nil)

	if err != nil || data == nil {
		return nil, errors.New("That item does not exist!")
	}

	group := &Group{}
	err = json.Unmarshal(data.Value, group)
	if err != nil {
		return nil, err
	}

	return group, nil
}
func (cs *ConfigStore) GetAllGroup() ([]*Group, error) {
	kv := cs.cli.KV()
	data, _, err := kv.List(groups, nil)
	if err != nil {
		return nil, err
	}

	groups := []*Group{}
	for _, pair := range data {
		group := &Group{}

		err = json.Unmarshal(pair.Value, group)
		if err == nil {

			groups = append(groups, group)
		}

	}

	return groups, nil
}

func (cs *ConfigStore) DeleteGroup(id, version string) (map[string]string, error) {
	kv := cs.cli.KV()
	_, err := kv.Delete(createGroupWithIdAndVersionKey(id, version), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted config": id + version}, nil
}

func (cs *ConfigStore) UpdateGroupWithNewVersion(group *Group, id string) (*Group, error) {
	kv := cs.cli.KV()

	_, er := cs.GetGroup(id, group.Version)

	if er == nil {
		return nil, errors.New("Group with that version already exists ")
	}

	sid := createGroupWithIdAndVersionKey(id, group.Version)
	group.Id = id

	data, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	c := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		return nil, err
	}

	cs.CreateLabelsForGroupConfigs(group)

	return group, nil

}

func (cs *ConfigStore) CreateLabelsForGroupConfigs(group *Group) {
	kv := cs.cli.KV()
	for _, configEntries := range group.Configs {
		labelKey := createGroupWithIdAndVersionAndLabelKey(group.Id, group.Version, configEntries)

		existingValue, _, err := kv.Get(labelKey, nil)
		entriesJson, err := json.Marshal([]map[string]string{configEntries})
		if err != nil {
			// nesto
		}

		if existingValue != nil {
			arr := []map[string]string{}
			err := json.Unmarshal(existingValue.Value, &arr)
			if err != nil {
				// nesto
			}
			arr = append(arr, configEntries)
			entriesJson, err = json.Marshal(arr)
			if err != nil {
				// nesto
			}
		}

		c := &api.KVPair{Key: labelKey, Value: entriesJson}
		_, err = kv.Put(c, nil)
		if err != nil {
			// nesto
		}
	}
}

func (cs *ConfigStore) UpdateGroup(group *Group, id string) (*Group, error) {
	kv := cs.cli.KV()

	sid := createGroupWithIdAndVersionKey(id, group.Version)
	group.Id = id

	data, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	c := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		return nil, err
	}

	cs.CreateLabelsForGroupConfigs(group)

	return group, nil

}

func (cs *ConfigStore) GetConfigsByLabels(id string, version string, params string) ([]map[string]string, error) {
	kv := cs.cli.KV()

	n := strings.Replace(params, "&", ";", -1)
	n1 := strings.Replace(n, "=", ":", -1)

	labelkey := fmt.Sprintf(groupIdVersionLabel, id, version, n1)

	keys, _, err := kv.List(labelkey, nil)
	if err != nil {
		return nil, err
	}

	arr := []map[string]string{}
	json.Unmarshal(keys[0].Value, &arr)

	return arr, nil
}
