package odl

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOdlVirtualInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualInterfaceAdd,
		Read:   resourceVirtualInterfaceRead,
		Delete: resourceVirtualInterfaceDelete,
		Schema: map[string]*schema.Schema{
			"tenant_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"operation": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateOperation,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"terminal_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"bridge_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"interface_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}
func resourceVirtualInterfaceAdd(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	tenantName := d.Get("tenant_name").(string)
	bridgeName := d.Get("bridge_name").(string)
	interfaceName := d.Get("interface_name").(string)
	var body map[string]interface{}
	var input map[string]string
	input = make(map[string]string)

	log.Println("[DEBUG] Creating Interface with name " + interfaceName)
	input["tenant-name"] = tenantName
	input["update-mode"] = "UPDATE"
	input["bridge-name"] = bridgeName
	input["interface-name"] = interfaceName

	if operation, found := d.GetOk("operation"); found {
		input["operation"] = operation.(string)
	}
	if description, found := d.GetOk("description"); found {
		input["description"] = description.(string)
	}
	if terminalName, found := d.GetOk("terminal_name"); found {
		input["terminal-name"] = terminalName.(string)
	}
	if enabled, found := d.GetOk("enabled"); found {
		input["enabled"] = strconv.FormatBool(enabled.(bool))
	}

	log.Println("[DEBUG] All options collected for interface with name " + interfaceName)
	body = make(map[string]interface{})
	body["input"] = input
	response, err := config.PostRequest("restconf/operations/vtn-vinterface:update-vinterface", body)
	if err != nil {
		log.Printf("[ERROR] POST Request failed")
		return err
	}
	isCreated, output, errorOutput, err := Status(response)
	if isCreated {
		d.SetId(tenantName + bridgeName + interfaceName + output.Output.Status)
	} else {
		if errorOutput != nil {
			log.Printf("[ERROR] While creating interface %s", errorOutput.Errors.Error[0].Message)
			return fmt.Errorf("[ERROR] While creating interface %s", errorOutput.Errors.Error[0].Message)
		}
		if err != nil {
			return fmt.Errorf("[ERROR] Whlie creating interface %s", err.Error())
		}
	}

	return nil
}
func resourceVirtualInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	tenantName := d.Get("tenant_name").(string)
	bridgeName := d.Get("bridge_name").(string)
	interfaceName := d.Get("interface_name").(string)

	log.Println("[DEBUG] Read Interface with name " + interfaceName)
	response, err := config.GetRequest("restconf/operational/vtn:vtns")
	if err != nil {
		log.Printf("[ERROR] POST Request failed")
		return err
	}
	present, err := CheckResponseVirtualInterfaceExists(response, tenantName, bridgeName, interfaceName)
	if err != nil {
		log.Println("[ERROR] Interface Read failed")
		return fmt.Errorf("[ERROR] Interface could not be read %v", err)
	}
	if !present {
		log.Println("[DEBUG] Interface with name " + bridgeName + "was not found")
		d.SetId("")
	}
	return nil
}
func resourceVirtualInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	err := resourceVirtualInterfaceRead(d, meta)
	if d.Id() == "" {
		return fmt.Errorf("[ERROR] Interface does not exists")
	}
	config := meta.(*Config)
	tenantName := d.Get("tenant_name").(string)
	bridgeName := d.Get("bridge_name").(string)
	interfaceName := d.Get("interface_name").(string)

	var body map[string]interface{}
	var input map[string]string
	input = make(map[string]string)

	if terminalName, found := d.GetOk("terminal_name"); found {
		input["terminal-name"] = terminalName.(string)
	}
	input["tenant-name"] = tenantName
	input["bridge-name"] = bridgeName
	input["interface-name"] = interfaceName

	body = make(map[string]interface{})
	body["input"] = input

	log.Println("[DEBUG] All options collected for Interface with name " + interfaceName)
	log.Println("[DEBUG] Preparing to destroy Interface with name " + interfaceName)

	response, err := config.PostRequest("restconf/operations/vtn-vinterface:remove-vinterface", body)
	if err != nil {
		log.Printf("[ERROR] POST Request failed")
		return err
	}
	isDestroyed, _, errorOutput, err := Status(response)
	if isDestroyed {
		d.SetId("")
	} else {
		if errorOutput != nil {
			log.Printf("[ERROR] While destroying interface %s", errorOutput.Errors.Error[0].Message)
			return fmt.Errorf("[ERROR] While creating interface %s", errorOutput.Errors.Error[0].Message)
		}
		if err != nil {
			return fmt.Errorf("[ERROR] Whlie destroying interface %s", err.Error())
		}
	}

	return nil
}
