package aws

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/helper/schema"
)

var awsSpotFleetRequestStateCancelled = map[string]struct{}{
	"cancelled":             {},
	"cancelled_running":     {},
	"cancelled_terminating": {},
}

func resourceAwsSpotFleetRequest() *schema.Resource {
	r := &schema.Resource{
		Create: resourceAwsSpotFleetRequestCreate,
		Read:   resourceAwsSpotFleetRequestRead,
		Delete: resourceAwsSpotFleetRequestDelete,
		Update: resourceAwsSpotFleetRequestUpdate,

		Schema: map[string]*schema.Schema{
			"ami": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"spot_price": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"target_capacity": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"security_groups": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"key_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"launch_specifications": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"weighted_capacity": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
							Default:  1.0,
						},
					},
				},
			},

			"wait_for_fulfillment": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"allocation_strategy": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},

			"client_token": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},

			"excess_capacity_termination_policy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"iam_fleet_role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"terminate_instances_with_expiration": &schema.Schema{
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},

			"valid_from": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},

			"valid_until": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},

			"tags": tagsSchema(),

			// Valid Values: submitted | active | cancelled | failed | cancelled_running | cancelled_terminating | modifying
			"spot_fleet_request_state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"create_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"active_instances": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"instance_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"spot_request_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}

	return r
}

func resourceAwsSpotFleetRequestCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	target_capacity := aws.Int64(int64(d.Get("target_capacity").(int)))
	if *target_capacity == int64(0) {
		log.Printf("target_capcity is 0, won't create")
		return nil
	}

	launchSpecs, err := buildAwsSpotFleetLaunchSpecifications(d, meta)
	if err != nil {
		return err
	}

	input := &ec2.RequestSpotFleetInput{
		SpotFleetRequestConfig: &ec2.SpotFleetRequestConfigData{
			IamFleetRole:         aws.String(d.Get("iam_fleet_role").(string)),
			LaunchSpecifications: launchSpecs,
			SpotPrice:            aws.String(d.Get("spot_price").(string)),
			TargetCapacity:       target_capacity,
		},
	}

	// optional params
	if v, ok := d.GetOk("allocation_strategy"); ok {
		input.SpotFleetRequestConfig.AllocationStrategy = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Requesting spot fleet: %s", input)
	resp, err := conn.RequestSpotFleet(input)
	if err != nil {
		return fmt.Errorf("Error requesting spot fleet: %s", err)
	}

	d.SetId(*resp.SpotFleetRequestId)

	// TODO: wait_for_fulfillment
	if d.Get("wait_for_fulfillment").(bool) {

	}

	return resourceAwsSpotFleetRequestRead(d, meta)
}

// Update spot state, etc
func resourceAwsSpotFleetRequestRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	// DescribeSpotFleetRequests
	{
		resp, err := conn.DescribeSpotFleetRequests(&ec2.DescribeSpotFleetRequestsInput{
			SpotFleetRequestIds: []*string{aws.String(d.Id())},
		})
		if err != nil {
			// If the spot fleet was not found, return nil so that we can show
			// that it is gone.
			if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidSpotFleetRequestID.NotFound" {
				d.SetId("")
				return nil
			}

			// Some other error, report it
			return err
		}

		// If nothing was found, then return no state
		if len(resp.SpotFleetRequestConfigs) == 0 {
			d.SetId("")
			return nil
		}

		request := resp.SpotFleetRequestConfigs[0]

		// if the request is cancelled, then it is gone
		if _, ok := awsSpotFleetRequestStateCancelled[*request.SpotFleetRequestState]; ok {
			d.SetId("")
			return nil
		}

		d.Set("spot_fleet_request_state", *request.SpotFleetRequestState)
		d.Set("create_time", request.CreateTime.Format(time.RFC3339))

		config := request.SpotFleetRequestConfig
		if config.AllocationStrategy != nil {
			d.Set("allocation_strategy", *config.AllocationStrategy)
		}
		if config.TargetCapacity != nil {
			d.Set("target_capacity", *config.TargetCapacity)
		}
		if config.ClientToken != nil {
			d.Set("client_token", *config.ClientToken)
		}
		if config.ExcessCapacityTerminationPolicy != nil {
			d.Set("excess_capacity_termination_policy", *config.ExcessCapacityTerminationPolicy)
		}
		if config.TerminateInstancesWithExpiration != nil {
			d.Set("terminate_instances_with_expiration", *config.TerminateInstancesWithExpiration)
		}
		if config.ValidFrom != nil {
			d.Set("valid_from", config.ValidFrom.Format(time.RFC3339))
		}
		if config.ValidUntil != nil {
			d.Set("valid_until", config.ValidUntil.Format(time.RFC3339))
		}
	}

	// DescribeSpotFleetInstances
	{
		resp, err := conn.DescribeSpotFleetInstances(&ec2.DescribeSpotFleetInstancesInput{
			SpotFleetRequestId: aws.String(d.Id()),
		})
		if err != nil {
			return err
		}

		activeInstances := make([]interface{}, 0)
		idNeedTag := make([]string, 0)
		for _, ins := range resp.ActiveInstances {
			instance := make(map[string]interface{})
			instance["instance_id"] = *ins.InstanceId
			instance["instance_type"] = *ins.InstanceType
			instance["spot_request_id"] = *ins.SpotInstanceRequestId

			activeInstances = append(activeInstances, instance)
			idNeedTag = append(idNeedTag, *ins.InstanceId, *ins.SpotInstanceRequestId)
		}
		d.Set("active_instances", activeInstances)

		// tag the spot request and instance
		conn.CreateTags(&ec2.CreateTagsInput{
			Resources: aws.StringSlice(idNeedTag),
			Tags:      tagsFromMap(d.Get("tags").(map[string]interface{})),
		})
	}

	return nil
}

func resourceAwsSpotFleetRequestUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	if d.HasChange("target_capacity") {
		target_capacity := aws.Int64(int64(d.Get("target_capacity").(int)))
		if *target_capacity == int64(0) {
			log.Printf("target_capcity is 0, delete the request")
			return resourceAwsSpotFleetRequestDelete(d, meta)
		} else {
			input := &ec2.ModifySpotFleetRequestInput{
				SpotFleetRequestId: aws.String(d.Id()),
				TargetCapacity:     target_capacity,
			}
			if v, ok := d.GetOk("excess_capacity_termination_policy"); ok {
				input.ExcessCapacityTerminationPolicy = aws.String(v.(string))
			}

			resp, err := conn.ModifySpotFleetRequest(input)
			if err != nil {
				return err
			}
			if !*resp.Return {
				return fmt.Errorf("Error modifying spot fleet (%s).", d.Id())
			}
		}
	}

	return resourceAwsSpotFleetRequestRead(d, meta)
}

func resourceAwsSpotFleetRequestDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).ec2conn

	_, err := conn.CancelSpotFleetRequests(&ec2.CancelSpotFleetRequestsInput{
		SpotFleetRequestIds: []*string{aws.String(d.Id())},
		TerminateInstances:  aws.Bool(true),
	})

	if err != nil {
		return fmt.Errorf("Error cancelling spot fleet (%s): %s", d.Id(), err)
	}

	return nil
}

func buildAwsSpotFleetLaunchSpecifications(
	d *schema.ResourceData, meta interface{}) ([]*ec2.SpotFleetLaunchSpecification, error) {
	specs := make([]*ec2.SpotFleetLaunchSpecification, 0)

	// subnet
	hasSubnet := false
	subnetID := ""
	if v, ok := d.GetOk("subnet_id"); ok {
		hasSubnet = true
		subnetID = v.(string)
	}

	// security group
	hasSecurityGroup := false
	groups := make([]*ec2.GroupIdentifier, 0)
	if v, ok := d.GetOk("security_groups"); ok {
		sgs := v.(*schema.Set).List()
		for _, v := range sgs {
			groups = append(groups,
				&ec2.GroupIdentifier{
					GroupId: aws.String(v.(string)),
				},
			)
		}
		if len(groups) > 0 {
			hasSecurityGroup = true
		}
	}

	for _, s := range d.Get("launch_specifications").([]interface{}) {
		specInput := s.(map[string]interface{})

		// required input
		spec := &ec2.SpotFleetLaunchSpecification{
			ImageId:      aws.String(d.Get("ami").(string)),
			InstanceType: aws.String(specInput["instance_type"].(string)),
			KeyName:      aws.String(d.Get("key_name").(string)),
		}

		if hasSubnet {
			spec.SubnetId = aws.String(subnetID)
		}

		if hasSecurityGroup {
			spec.SecurityGroups = groups
		}

		// weighted capacity
		if v, ok := specInput["weighted_capacity"]; ok {
			spec.WeightedCapacity = aws.Float64(v.(float64))
		}

		specs = append(specs, spec)
	}

	return specs, nil
}
