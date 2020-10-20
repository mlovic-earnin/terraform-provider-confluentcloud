package ccloud

import (
	"fmt"
	"log"
	"strconv"
  "time"

	ccloud "github.com/cgroschupp/go-client-confluent-cloud/confluentcloud"
	"github.com/hashicorp/terraform/helper/schema"
)

// TODO update
func serviceAccountResource() *schema.Resource {
	return &schema.Resource{
		Create: serviceAccountCreate,
		Read:   serviceAccountRead,
		Delete: serviceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Service Account Description",
			},
		},
	}
}

func serviceAccountCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)

	name := d.Get("name").(string)
	description := d.Get("description").(string)

	req := ccloud.ServiceAccountCreateRequest{
		Name:        name,
		Description: description,
	}

	serviceAccount, err := c.CreateServiceAccount(&req)
	if err == nil {
		d.SetId(fmt.Sprintf("%d", serviceAccount.ID))

		err = d.Set("name", serviceAccount.Name)
		if err != nil {
			return err
		}

		err = d.Set("description", serviceAccount.Description)
		if err != nil {
			return err
		}
	} else {
		log.Printf("[ERROR] Could not create Service Account: %s", err)
	}

	return nil
}

func serviceAccountRead(d *schema.ResourceData, meta interface{}) error {
  log.Printf("serviceAccountRead")
  log.Printf("waiting 3 seconds...")
  time.Sleep(3 * time.Second)
  //panic("mlovicc")
	c := meta.(*ccloud.Client)

  //c.SetTimeout(10 * time.Second)
	//accountID := d.Get("environment_id").(string)

  log.Printf("starting request")
	serviceAccounts, err := c.ListServiceAccounts()
  log.Printf("waiting 3 seconds...")
  log.Printf("starting request again")
	c = meta.(*ccloud.Client)
  time.Sleep(3 * time.Second)
	serviceAccounts, err = c.ListServiceAccounts()
  log.Printf("starting request again again")
	serviceAccounts, err = c.ListServiceAccounts()
  log.Printf("ending request")

	if err != nil {
		log.Printf("[ERROR] Error with ListServiceAccounts request")
    return err
  }

  //log.Printf("s accounts: %v", serviceAccounts)
  //log.Println(serviceAccounts)

  serviceAccount := serviceAccounts[0]
  log.Printf("first s account: %s", serviceAccount)
  log.Printf("first s account ID: %d", serviceAccount.ID)

  log.Printf("num SAs found: %d", len(serviceAccounts))

  ID, err := strconv.Atoi(d.Id())
  if err != nil {
    log.Printf("[ERROR] Could not parse Service Account ID %s to int", d.Id())
    return err
  }

  log.Printf("Looking for SA with ID %d", ID)

  found := false
  for i := range serviceAccounts {
      //if serviceAccounts[i].Name == "marko-provider-dev-test" {
      // TODO should the ID be the int instead? accept either?
      //   Yes, switching to using int/id


      if serviceAccounts[i].ID == ID {
          log.Printf("found on index: %d", i)
          // Found!
          found = true
          serviceAccount = serviceAccounts[i]
          log.Printf("found s account: %s", serviceAccount)
          //var serviceAccount ServiceAccount = serviceAccounts[i]

      }
  }
  if found == false {
    log.Printf("[ERROR] not found!")
    panic("[ERROR] not found!")
  }

	//serviceAccount, err := c.GetServiceAccount(d.Id(), accountID)
	//serviceAccount, err := c.GetServiceAccount(d.Id())
  //panic(fmt.Sprintf("%d", serviceAccount.ID))
	if err == nil {
		log.Printf("[WARN] hello %s", serviceAccount.Name)
    // why was I trying to set ID??
    //d.SetId(fmt.Sprintf("%d", serviceAccount.ID))
		//d.SetId(fmt.Sprintf("%d", "asdf"))

		err = d.Set("name", serviceAccount.Name)
		if err != nil {
      log.Printf("[ERROR] error setting name: %s", serviceAccount.Name)
			return err
		}

		err = d.Set("description", serviceAccount.Description)
		if err != nil {
      log.Printf("[ERROR] error setting desc: %s", serviceAccount.Description)
			return err
		}
	} else { log.Printf("[ERROR] Could not read Service Account: %s", err)
	}
  log.Printf("Finished reading SA with no errors")
  log.Printf("waiting 15 seconds...")
  time.Sleep(15 * time.Second)
	return nil
}

func serviceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ccloud.Client)

	ID, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Printf("[ERROR] Could not parse Service Account ID %s to int", d.Id())
		return err
	}

	err = c.DeleteServiceAccount(ID)
	if err != nil {
		log.Printf("[ERROR] Service Account can not be deleted: %d", ID)
		return err
	}

	log.Printf("[INFO] Service Account deleted: %d", ID)

	return nil
}
