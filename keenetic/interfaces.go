package keenetic

import (
	"errors"
)

func (k *Keenetic) GetInterfaces() (map[string]bool, error) {
	r, err := k.c.R().SetResult(map[string]map[string]interface{}{}).Get("/rci/show/interface")
	if err != nil {
		return nil, err
	}
	if r.IsError() {
		return nil, errors.New(r.String())
	}

	interfaces := make(map[string]bool)

	for i, d := range *r.Result().(*map[string]map[string]interface{}) {
		up, _ := d["state"]
		u, _ := up.(string)
		interfaces[i] = u == "up"
	}

	return interfaces, nil
}

func (k *Keenetic) SetInterfaces(interfaces map[string]bool) error {
	var req []map[string]interface{}
	for i, u := range interfaces {
		req = append(req, map[string]interface{}{
			"interface": map[string]interface{}{
				i: map[string]interface{}{
					"up": map[string]interface{}{
						"no": !u,
					},
				},
			},
		})
	}
	req = append(req, map[string]interface{}{
		"system": map[string]interface{}{
			"configuration": map[string]interface{}{
				"save": true,
			},
		},
	})

	r, err := k.c.R().SetBody(req).Post("/rci/")
	if err != nil {
		return err
	}
	if r.IsError() {
		return errors.New(r.String())
	}

	return nil
}
