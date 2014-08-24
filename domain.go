package mailgun

import (
	"fmt"
	"strconv"
)

type DomainResponse struct {
	Domain Domain `json:"domain"`
}

// Domain is used to represent a retrieved Domain. All properties
// are set as strings.
type Domain struct {
	Name         string `json:"name"`
	SmtpLogin    string `json:"smtp_login"`
	SmtpPassword string `json:"smtp_password"`
	SpamAction   string `json:"spam_action"`
	Wildcard     bool   `json:"wildcard"`
}

func (d *Domain) StringWildcard() string {
	return strconv.FormatBool(d.Wildcard)
}

// CreateDomain contains the request parameters to create a new
// domain.
type CreateDomain struct {
	Name         string // Name of the domain
	SmtpPassword string
	SpamAction   string // one of disabled or tag
	Wildcard     bool
}

// CreateDomain creates a domain from the parameters specified and
// returns an error if it fails. If no error and the name is returned,
// the Domain was succesfully created.
func (c *Client) CreateDomain(opts *CreateDomain) (string, error) {
	// Make the request parameters
	params := make(map[string]string)

	params["name"] = opts.Name
	params["smtp_password"] = opts.SmtpPassword
	params["spam_action"] = opts.SpamAction
	params["wildcard"] = strconv.FormatBool(opts.Wildcard)

	req, err := c.NewRequest(params, "POST", "/domains")
	if err != nil {
		return "", err
	}

	resp, err := checkResp(c.Http.Do(req))

	if err != nil {
		return "", fmt.Errorf("Error creating domain: %s", err)
	}

	domain := new(DomainResponse)

	err = decodeBody(resp, &domain)

	if err != nil {
		return "", fmt.Errorf("Error parsing domain response: %s", err)
	}

	// The request was successful
	return domain.Domain.Name, nil
}

// DestroyDomain destroys a domain by the ID specified and
// returns an error if it fails. If no error is returned,
// the Domain was succesfully destroyed.
func (c *Client) DestroyDomain(name string) error {
	req, err := c.NewRequest(map[string]string{}, "DELETE", fmt.Sprintf("/domains/%s", name))

	if err != nil {
		return err
	}

	_, err = checkResp(c.Http.Do(req))

	if err != nil {
		return fmt.Errorf("Error destroying domain: %s", err)
	}

	// The request was successful
	return nil
}

// RetrieveDomain gets  a domain by the ID specified and
// returns a Domain and an error. An error will be returned for failed
// requests with a nil Domain.
func (c *Client) RetrieveDomain(name string) (Domain, error) {
	req, err := c.NewRequest(map[string]string{}, "GET", fmt.Sprintf("/domains/%s", name))

	if err != nil {
		return Domain{}, err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return Domain{}, fmt.Errorf("Error destroying domain: %s", err)
	}

	domain := new(DomainResponse)

	err = decodeBody(resp, domain)

	if err != nil {
		return Domain{}, fmt.Errorf("Error decoding domain response: %s", err)
	}

	// The request was successful
	return domain.Domain, nil
}
