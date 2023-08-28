package kb

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pkg/errors"
	kibana "github.com/samuel-deal-tisseo/go-kibana-rest/v7"
)

func TestAccKibanaObject(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckKibanaObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: getTestKibanaObject(),
				Check: resource.ComposeTestCheckFunc(
					testCheckKibanaObjectExists("kibana_object.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckKibanaObjectExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No object ID is set")
		}

		// Use static value that match the current test
		deepReference := true
		exportObject := map[string]string{}
		exportObject["id"] = "logstash-log-*"
		exportObject["type"] = "index-pattern"
		exportObjects := []map[string]string{exportObject}
		space := "default"

		meta := testAccProvider.Meta()

		client := meta.(*kibana.Client)
		data, err := client.API.KibanaSavedObject.Export(nil, exportObjects, deepReference, space)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			return errors.Errorf("Object %s not found", rs.Primary.ID)
		}

		return nil
	}
}

func testCheckKibanaObjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kibana_object" {
			continue
		}

	}

	return nil
}

func getTestKibanaObject() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(`
resource "kibana_object" "test" {
  name 				= "terraform-test"
  data				= file("%s/../fixtures/index-pattern.json")
  deep_reference	= "true"
  export_types    	= ["index-pattern"]
}
`, path)

}
