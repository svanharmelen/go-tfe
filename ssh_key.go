package tfe

import (
	"context"
	"fmt"
	"net/url"
)

// Compile-time proof of interface implementation.
var _ SSHKeys = (*sshKeys)(nil)

// SSHKeys describes all the SSH key related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs:
// https://www.terraform.io/docs/cloud/api/ssh-keys.html
type SSHKeys interface {
	// List all the SSH keys for a given organization
	List(ctx context.Context, organization string, options *SSHKeyListOptions) (*SSHKeyList, error)

	// Create an SSH key and associate it with an organization.
	Create(ctx context.Context, organization string, options SSHKeyCreateOptions) (*SSHKey, error)

	// Read an SSH key by its ID.
	Read(ctx context.Context, sshKeyID string) (*SSHKey, error)

	// Update an SSH key by its ID.
	Update(ctx context.Context, sshKeyID string, options SSHKeyUpdateOptions) (*SSHKey, error)

	// Delete an SSH key by its ID.
	Delete(ctx context.Context, sshKeyID string) error
}

// sshKeys implements SSHKeys.
type sshKeys struct {
	client *Client
}

// SSHKeyList represents a list of SSH keys.
type SSHKeyList struct {
	*Pagination
	Items []*SSHKey
}

// SSHKey represents a SSH key.
type SSHKey struct {
	ID   string `jsonapi:"primary,ssh-keys"`
	Name string `jsonapi:"attr,name"`
}

// SSHKeyListOptions represents the options for listing SSH keys.
type SSHKeyListOptions struct {
	ListOptions
}

// SSHKeyCreateOptions represents the options for creating an SSH key.
type SSHKeyCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,ssh-keys"`

	// A name to identify the SSH key.
	Name *string `jsonapi:"attr,name"`

	// The content of the SSH private key.
	Value *string `jsonapi:"attr,value"`
}

// SSHKeyUpdateOptions represents the options for updating an SSH key.
type SSHKeyUpdateOptions struct {
	// For internal use only!
	ID string `jsonapi:"primary,ssh-keys"`

	// Optional: A new name to identify the SSH key.
	Name *string `jsonapi:"attr,name,omitempty"`
}

// List all the SSH keys for a given organization
func (s *sshKeys) List(ctx context.Context, organization string, options *SSHKeyListOptions) (*SSHKeyList, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/ssh-keys", url.QueryEscape(organization))
	req, err := s.client.newRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	kl := &SSHKeyList{}
	err = s.client.do(ctx, req, kl)
	if err != nil {
		return nil, err
	}

	return kl, nil
}

// Create an SSH key and associate it with an organization.
func (s *sshKeys) Create(ctx context.Context, organization string, options SSHKeyCreateOptions) (*SSHKey, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("organizations/%s/ssh-keys", url.QueryEscape(organization))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	k := &SSHKey{}
	err = s.client.do(ctx, req, k)
	if err != nil {
		return nil, err
	}

	return k, nil
}

// Read an SSH key by its ID.
func (s *sshKeys) Read(ctx context.Context, sshKeyID string) (*SSHKey, error) {
	if !validStringID(&sshKeyID) {
		return nil, ErrInvalidSHHKeyID
	}

	u := fmt.Sprintf("ssh-keys/%s", url.QueryEscape(sshKeyID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	k := &SSHKey{}
	err = s.client.do(ctx, req, k)
	if err != nil {
		return nil, err
	}

	return k, nil
}

// Update an SSH key by its ID.
func (s *sshKeys) Update(ctx context.Context, sshKeyID string, options SSHKeyUpdateOptions) (*SSHKey, error) {
	if !validStringID(&sshKeyID) {
		return nil, ErrInvalidSHHKeyID
	}

	u := fmt.Sprintf("ssh-keys/%s", url.QueryEscape(sshKeyID))
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	k := &SSHKey{}
	err = s.client.do(ctx, req, k)
	if err != nil {
		return nil, err
	}

	return k, nil
}

// Delete an SSH key by its ID.
func (s *sshKeys) Delete(ctx context.Context, sshKeyID string) error {
	if !validStringID(&sshKeyID) {
		return ErrInvalidSHHKeyID
	}

	u := fmt.Sprintf("ssh-keys/%s", url.QueryEscape(sshKeyID))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

func (o SSHKeyCreateOptions) valid() error {
	if !validString(o.Name) {
		return ErrRequiredName
	}
	if !validString(o.Value) {
		return ErrRequiredValue
	}
	return nil
}
