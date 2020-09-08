package fusionauth

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func handleStringSlice(key string, data *schema.ResourceData) []string {
	set := data.Get(key).(*schema.Set)
	l := set.List()
	s := make([]string, 0, len(l))

	for _, x := range l {
		s = append(s, x.(string))
	}

	return s
}
