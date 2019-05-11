package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"github.com/iron-io/iron_go3/mq"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/iron-io/iron_go3/config"
)

// dataSourceQueues() retrieves information about queues.
func dataSourceQueues() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"filter_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The name filter",
				ForceNew:    true,
			},
			"names": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"project_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The project id",
				ForceNew:    true,
			},
		},

		Read: dataSourceQueuesRead,
	}
}

// dataSourceQueuesRead reads information about available queues.
func dataSourceQueuesRead(d *schema.ResourceData, m interface{}) error {
	clientSettings := m.(ClientSettings)
	clientSettingsMQ := config.Settings{}
	clientSettingsMQ.UseSettings(&clientSettings.MQ)
	clientSettingsMQ.ProjectId = d.Get("project_id").(string)

	// Prepare the filters.
	filterName := d.Get("filter_name").(string)
	filterNameMode := 0

	if filterName != "" {
		if len(filterName) >= 2 && strings.HasPrefix(filterName, "*") && strings.HasSuffix(filterName, "*") {
			filterName = filterName[1 : len(filterName)-1]
			filterNameMode = 1
		} else if strings.HasPrefix(filterName, "*") {
			filterName = filterName[1:len(filterName)]
			filterNameMode = 2
		} else if strings.HasSuffix(filterName, "*") {
			filterName = filterName[0 : len(filterName)-1]
			filterNameMode = 3
		} else {
			filterNameMode = 4
		}

		if filterNameMode > 0 && filterName == "" {
			return errors.New("The name filter cannot be an empty wildcard filter")
		}
	}

	// Retrieve the list of projects.
	queues, err := mq.ListQueues(clientSettingsMQ, "", "", 1000)

	if err != nil {
		return err
	}

	// Parse and filter the results.
	names := make([]string, 0)

	for _, v := range queues {
		if filterNameMode == 1 && !strings.Contains(v.Name, filterName) {
			continue
		} else if filterNameMode == 2 && !strings.HasSuffix(v.Name, filterName) {
			continue
		} else if filterNameMode == 3 && !strings.HasPrefix(v.Name, filterName) {
			continue
		} else if filterNameMode == 4 && strings.Compare(v.Name, filterName) != 0 {
			continue
		}

		names = append(names, v.Name)
	}

	h := sha256.New()
	h.Write([]byte(strings.Join(names, ",")))

	d.SetId(fmt.Sprintf("%x", h.Sum(nil)))
	d.Set("names", names)

	return nil
}