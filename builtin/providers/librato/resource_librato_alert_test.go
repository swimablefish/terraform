package librato

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/henrikhodne/go-librato/librato"
)

func TestAccLibratoAlert_Basic(t *testing.T) {
	var alert librato.Alert

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLibratoAlertDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLibratoAlertConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLibratoAlertExists("librato_alert.foobar", &alert),
					testAccCheckLibratoAlertName(&alert, "FooBar"),
					resource.TestCheckResourceAttr(
						"librato_alert.foobar", "name", "FooBar"),
				),
			},
		},
	})
}

func TestAccLibratoAlert_Full(t *testing.T) {
	var alert librato.Alert

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLibratoAlertDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLibratoAlertConfig_full,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLibratoAlertExists("librato_alert.foobar", &alert),
					testAccCheckLibratoAlertName(&alert, "FooBar"),
					resource.TestCheckResourceAttr(
						"librato_alert.foobar", "name", "FooBar"),
				),
			},
		},
	})
}

func TestAccLibratoAlert_Updated(t *testing.T) {
	var alert librato.Alert

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLibratoAlertDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLibratoAlertConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLibratoAlertExists("librato_alert.foobar", &alert),
					testAccCheckLibratoAlertName(&alert, "FooBar"),
					resource.TestCheckResourceAttr(
						"librato_alert.foobar", "name", "FooBar"),
				),
			},
			resource.TestStep{
				Config: testAccCheckLibratoAlertConfig_new_value,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLibratoAlertExists("librato_alert.foobar", &alert),
					testAccCheckLibratoAlertName(&alert, "BarBaz"),
					resource.TestCheckResourceAttr(
						"librato_alert.foobar", "name", "BarBaz"),
				),
			},
		},
	})
}

func testAccCheckLibratoAlertDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*librato.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "librato_alert" {
			continue
		}

		id, err := strconv.ParseUint(rs.Primary.ID, 10, 0)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		_, _, err = client.Alerts.Get(uint(id))

		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
	}

	return nil
}

func testAccCheckLibratoAlertName(alert *librato.Alert, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if alert.Name == nil || *alert.Name != name {
			return fmt.Errorf("Bad name: %s", *alert.Name)
		}

		return nil
	}
}

func testAccCheckLibratoAlertExists(n string, alert *librato.Alert) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Alert ID is set")
		}

		client := testAccProvider.Meta().(*librato.Client)

		id, err := strconv.ParseUint(rs.Primary.ID, 10, 0)
		if err != nil {
			return fmt.Errorf("ID not a number")
		}

		foundAlert, _, err := client.Alerts.Get(uint(id))

		if err != nil {
			return err
		}

		if foundAlert.ID == nil || *foundAlert.ID != uint(id) {
			return fmt.Errorf("Alert not found")
		}

		*alert = *foundAlert

		return nil
	}
}

const testAccCheckLibratoAlertConfig_basic = `
resource "librato_alert" "foobar" {
    name = "FooBar"
    description = "A Test Alert"
}`

const testAccCheckLibratoAlertConfig_new_value = `
resource "librato_alert" "foobar" {
    name = "BarBaz"
    description = "A Test Alert"
}`

const testAccCheckLibratoAlertConfig_full = `
resource "librato_service" "foobar" {
    title = "Foo Bar"
    type = "mail"
    settings = <<EOF
{
  "addresses": "admin@example.com"
}
EOF
}

resource "librato_alert" "foobar" {
    name = "FooBar"
    description = "A Test Alert"
    services = [ "${librato_service.foobar.id}" ]
    condition {
      type = "above"
      threshold = 10
      metric_name = "librato.cpu.percent.idle"
    }
    attributes {
      runbook_url = "https://www.youtube.com/watch?v=oHg5SJYRHA0"
    }
    active = false
    rearm_seconds = 300
}`
