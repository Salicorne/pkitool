/*
 * Swagger Petstore - OpenAPI 3.0
 *
 * This is the REST API used by Salicorne/PKITool
 *
 * API version: 0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package models

type PkiRequest struct {

	Name string `json:"name,omitempty"`

	DN *Dn `json:"DN,omitempty"`

	ValidityDays int64 `json:"validityDays,omitempty"`
}