package cli

import (
	"fmt"
	"os"

	"net/http"

	"github.com/ory/hydra/config"
	"github.com/ory/hydra/pkg"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/ory/ladon"
	"github.com/spf13/cobra"
	"github.com/square/go-jose/json"
)

type PolicyHandler struct {
	Config *config.Config
}

func (h *PolicyHandler) newPolicyManager(cmd *cobra.Command) *hydra.PolicyApi {
	c := hydra.NewPolicyApiWithBasePath(h.Config.ClusterURL)
	c.Configuration.Transport = h.Config.OAuth2Client(cmd).Transport

	if term, _ := cmd.Flags().GetBool("fake-tls-termination"); term {
		c.Configuration.DefaultHeader["X-Forwarded-Proto"] = "https"
	}

	return c
}

func newPolicyHandler(c *config.Config) *PolicyHandler {
	return &PolicyHandler{
		Config: c,
	}
}

func (h *PolicyHandler) CreatePolicy(cmd *cobra.Command, args []string) {
	m := h.newPolicyManager(cmd)

	if files, _ := cmd.Flags().GetStringSlice("files"); len(files) > 0 {
		for _, path := range files {
			reader, err := os.Open(path)
			pkg.Must(err, "Could not open file %s: %s", path, err)

			var p hydra.Policy
			err = json.NewDecoder(reader).Decode(&p)
			pkg.Must(err, "Could not parse JSON: %s", err)

			_, response, err := m.CreatePolicy(p)
			checkResponse(response, err, http.StatusCreated)
			fmt.Printf("Imported policy %s from %s.\n", p.Id, path)
		}
		return
	}

	id, _ := cmd.Flags().GetString("id")
	description, _ := cmd.Flags().GetString("description")
	subjects, _ := cmd.Flags().GetStringSlice("subjects")
	resources, _ := cmd.Flags().GetStringSlice("resources")
	actions, _ := cmd.Flags().GetStringSlice("actions")
	isAllow, _ := cmd.Flags().GetBool("allow")
	if len(subjects) == 0 || len(resources) == 0 || len(actions) == 0 {
		fmt.Println(cmd.UsageString())
		fmt.Println("")
		fmt.Println("Got empty subject, resource or action list")
		return
	}

	effect := ladon.DenyAccess
	if isAllow {
		effect = ladon.AllowAccess
	}

	result, response, err := m.CreatePolicy(hydra.Policy{
		Id:          id,
		Description: description,
		Subjects:    subjects,
		Resources:   resources,
		Actions:     actions,
		Effect:      effect,
	})
	checkResponse(response, err, http.StatusCreated)
	fmt.Printf("Created policy %s.\n", result.Id)
}

func (h *PolicyHandler) AddResourceToPolicy(cmd *cobra.Command, args []string) {
	m := h.newPolicyManager(cmd)
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	p, response, err := m.GetPolicy(args[0])
	checkResponse(response, err, http.StatusOK)

	p.Resources = append(p.Resources, args[1:]...)

	_, response, err = m.UpdatePolicy(p.Id, *p)
	checkResponse(response, err, http.StatusOK)
	fmt.Printf("Added resources to policy %s", p.Id)
}

func (h *PolicyHandler) RemoveResourceFromPolicy(cmd *cobra.Command, args []string) {
	m := h.newPolicyManager(cmd)
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	p, response, err := m.GetPolicy(args[0])
	checkResponse(response, err, http.StatusOK)

	resources := []string{}
	for _, r := range p.Resources {
		var filter bool
		for _, a := range args[1:] {
			if r == a {
				filter = true
			}
		}
		if !filter {
			resources = append(resources, r)
		}
	}
	p.Resources = resources

	_, response, err = m.UpdatePolicy(p.Id, *p)
	checkResponse(response, err, http.StatusOK)
	fmt.Printf("Removed resources from policy %s", p.Id)
}

func (h *PolicyHandler) AddSubjectToPolicy(cmd *cobra.Command, args []string) {
	m := h.newPolicyManager(cmd)
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	p, response, err := m.GetPolicy(args[0])
	checkResponse(response, err, http.StatusOK)

	p.Subjects = append(p.Subjects, args[1:]...)

	_, response, err = m.UpdatePolicy(p.Id, *p)
	checkResponse(response, err, http.StatusOK)
	fmt.Printf("Added subjects to policy %s", p.Id)
}

func (h *PolicyHandler) RemoveSubjectFromPolicy(cmd *cobra.Command, args []string) {
	m := h.newPolicyManager(cmd)
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	p, response, err := m.GetPolicy(args[0])
	checkResponse(response, err, http.StatusOK)

	subjects := []string{}
	for _, r := range p.Subjects {
		var filter bool
		for _, a := range args[1:] {
			if r == a {
				filter = true
			}
		}
		if !filter {
			subjects = append(subjects, r)
		}
	}
	p.Subjects = subjects
	p.Subjects = append(p.Subjects, args[1:]...)

	_, response, err = m.UpdatePolicy(p.Id, *p)
	checkResponse(response, err, http.StatusOK)
	fmt.Printf("Removed subjects from policy %s", p.Id)
}

func (h *PolicyHandler) AddActionToPolicy(cmd *cobra.Command, args []string) {
	m := h.newPolicyManager(cmd)
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	p, response, err := m.GetPolicy(args[0])
	checkResponse(response, err, http.StatusOK)

	p.Actions = append(p.Actions, args[1:]...)

	_, response, err = m.UpdatePolicy(p.Id, *p)
	checkResponse(response, err, http.StatusOK)
	fmt.Printf("Added actions to policy %s", p.Id)
}

func (h *PolicyHandler) RemoveActionFromPolicy(cmd *cobra.Command, args []string) {
	m := h.newPolicyManager(cmd)
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	p, response, err := m.GetPolicy(args[0])
	checkResponse(response, err, http.StatusOK)

	actions := []string{}
	for _, r := range p.Actions {
		var filter bool
		for _, a := range args[1:] {
			if r == a {
				filter = true
			}
		}
		if !filter {
			actions = append(actions, r)
		}
	}
	p.Actions = actions

	_, response, err = m.UpdatePolicy(p.Id, *p)
	checkResponse(response, err, http.StatusOK)
	fmt.Printf("Removed actions from policy %s", p.Id)
}

func (h *PolicyHandler) GetPolicy(cmd *cobra.Command, args []string) {
	m := h.newPolicyManager(cmd)
	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	p, response, err := m.GetPolicy(args[0])
	checkResponse(response, err, http.StatusOK)

	fmt.Printf("%s\n", formatResponse(p))
}

func (h *PolicyHandler) DeletePolicy(cmd *cobra.Command, args []string) {
	m := h.newPolicyManager(cmd)
	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	for _, arg := range args {
		response, err := m.DeletePolicy(arg)
		checkResponse(response, err, http.StatusNoContent)
		fmt.Printf("Policy %s deleted.\n", arg)
	}
}
