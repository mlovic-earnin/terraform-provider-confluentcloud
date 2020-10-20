package ccloud

import (
	"fmt"
	"log"
  "strings"

	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform/helper/schema"
)

func apiKeyResource() *schema.Resource {
	return &schema.Resource{
		Create: apiKeyCreate,
		Read:   apiKeyRead,
		Delete: apiKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "",
			},
			"logical_clusters": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				ForceNew:    true,
				Description: "Logical Cluster ID List to create API KEY",
			},
			"user_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "User ID",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Environment ID",
			},
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
      "description": {
        Type:      schema.TypeString,
        // TODO define update and change
        ForceNew:    true,
        Optional: true,
      },
		},
	}
}

func apiKeyCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)

	clusterID := d.Get("cluster_id").(string)
	logicalClusters := d.Get("logical_clusters").([]interface{})
	accountID := d.Get("environment_id").(string)
	userID := d.Get("user_id").(int)

	logicalClustersReq := []ccloud.LogicalCluster{}

	if len(clusterID) > 0 {
		logicalClustersReq = append(logicalClustersReq, ccloud.LogicalCluster{ID: clusterID})
	}

	for i := range logicalClusters {
		if clusterID != logicalClusters[i].(string) {
			logicalClustersReq = append(logicalClustersReq, ccloud.LogicalCluster{
				ID: logicalClusters[i].(string),
			})
		}
	}

	req := ccloud.ApiKeyCreateRequest{
		AccountID:       accountID,
		UserID:          userID,
		LogicalClusters: logicalClustersReq,
	}

	key, err := c.CreateAPIKey(&req)
	if err == nil {
		d.SetId(fmt.Sprintf("%d", key.ID))

		err = d.Set("key", key.Key)
		if err != nil {
			return err
		}

		err = d.Set("secret", key.Secret)
		if err != nil {
			return err
		}
	} else {
		log.Printf("[ERROR] Could not create API key: %s", err)
	}

	return err
}

func apiKeyRead(d *schema.ResourceData, meta interface{}) error {
  fmt.Println("apiKeyRead() called")
  log.Printf("apiKeyRead")
	c := meta.(*ccloud.Client)

  s_parts := strings.Split(d.Id(), "/")
  env_id, cluster_id, key_v := s_parts[0], s_parts[1], s_parts[2]

  apiKeys, err := c.ListAPIKeys(cluster_id, env_id)
  log.Printf("first key: %s", apiKeys[0])

  key := apiKeys[0]
  log.Printf("first key: %s", key)


  for i := range apiKeys {
      if apiKeys[i].Key == key_v {
          log.Printf("found on index: %d", i)
          // Found!
          key = apiKeys[i]
          log.Printf("found key: %s", key)
          //var serviceAccount ServiceAccount = serviceAccounts[i]

      }
  }

  //key, err := c.GetAPIKey(d.Id())
	if err == nil {
		log.Printf("[WARN] hello %s", key.Key)
		//d.SetId(fmt.Sprintf("%d", key.ID))
		d.SetId(d.Id())

    // TODO
		//err = d.Set("cluster_id", key.)
		//if err != nil {
			//return err
		//}

    //TODO logical_clusters
		err = d.Set("user_id", key.UserID)
		if err != nil {
			return err
		}
    err = d.Set("environment_id", key.AccountID)
    if err != nil {
      return err
    }
    err = d.Set("description", key.Description)
    if err != nil {
      return err
    }
		err = d.Set("key", key.Key)
		if err != nil {
			return err
		}
		err = d.Set("secret", key.Secret)
		if err != nil {
			return err
		}
	} else {
		log.Printf("[ERROR] Could not read API Key: %s", err)
	}
	return nil
}

func apiKeyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] API key cannot be deleted: %s", d.Id())
	return nil
}
