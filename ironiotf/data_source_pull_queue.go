/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package ironiotf

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/iron-io/iron_go3/config"
	"github.com/iron-io/iron_go3/mq"
)

const dataSourcePullQueueMessageCountKey = "message_count"
const dataSourcePullQueueMessageCountTotalKey = "message_count_total"
const dataSourcePullQueueNameKey = "name"
const dataSourcePullQueueProjectIDKey = "project_id"

// dataSourcePullQueue reads information about IronMQ pull queues.
func dataSourcePullQueue() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			dataSourcePullQueueMessageCountKey: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of messages currently in the queue",
			},
			dataSourcePullQueueMessageCountTotalKey: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of messages which have been processed by the queue",
			},
			dataSourcePullQueueNameKey: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the queue",
				ForceNew:    true,
			},
			dataSourcePullQueueProjectIDKey: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The project id",
				ForceNew:    true,
			},
		},

		Read: dataSourcePullQueueRead,
	}
}

// dataSourcePullQueueRead reads information about an existing pull queue.
func dataSourcePullQueueRead(d *schema.ResourceData, m interface{}) error {
	clientSettings := m.(ClientSettings)
	clientSettingsMQ := config.Settings{}
	clientSettingsMQ.UseSettings(&clientSettings.MQ)

	projectID := d.Get(dataSourcePullQueueProjectIDKey).(string)
	queueName := d.Get(dataSourcePullQueueNameKey).(string)

	clientSettingsMQ.ProjectId = projectID

	queue := mq.ConfigNew(queueName, &clientSettingsMQ)
	queueInfo, err := queue.Info()

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")

			return nil
		}

		return err
	}

	if queueInfo.Type != "pull" {
		d.SetId("")

		return nil
	}

	d.Set(dataSourcePullQueueMessageCountKey, queueInfo.Size)
	d.Set(dataSourcePullQueueMessageCountTotalKey, queueInfo.TotalMessages)

	d.SetId(queueNameToID(clientSettingsMQ.ProjectId, queueName))

	return nil
}
