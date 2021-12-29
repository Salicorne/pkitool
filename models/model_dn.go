/*
 * Swagger Petstore - OpenAPI 3.0
 *
 * This is the REST API used by Salicorne/PKITool
 *
 * API version: 0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package models

type Dn struct {

	Country []string `json:"country,omitempty"`

	State []string `json:"state,omitempty"`

	Locality []string `json:"locality,omitempty"`

	OrgUnit []string `json:"orgUnit,omitempty"`

	Organization []string `json:"organization,omitempty"`

	CommonName string `json:"commonName,omitempty"`
}