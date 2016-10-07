package aws

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

// See http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/enable-access-logs.html#attach-bucket-policy
var elbAccountIdPerRegionMap = map[string]string{
	"ap-northeast-1": "582318560864",
	"ap-northeast-2": "600734575887",
	"ap-south-1":     "718504428378",
	"ap-southeast-1": "114774131450",
	"ap-southeast-2": "783225319266",
	"cn-north-1":     "638102146993",
	"eu-central-1":   "054676820928",
	"eu-west-1":      "156460612806",
	"sa-east-1":      "507241528517",
	"us-east-1":      "127311923021",
	"us-gov-west":    "048591011584",
	"us-west-1":      "027434742980",
	"us-west-2":      "797873946194",
}

func dataSourceAwsElbServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsElbServiceAccountRead,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"arn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAwsElbServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	region := meta.(*AWSClient).region
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
	}

	if accid, ok := elbAccountIdPerRegionMap[region]; ok {
		d.SetId(accid)

		d.Set("arn", "arn:aws:iam::"+accid+":root")

		return nil
	}

	return fmt.Errorf("Unknown region (%q)", region)
}
