/*
Copyright 2014 Workiva, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rest

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

// FooResource represents a domain model for which we want to perform CRUD operations with.
// Endpoints can operate on any type of entity -- primitive, struct, or composite -- so long
// as it is serializable (by default, this means JSON-serializable via either MarshalJSON
// or JSON struct tags).
type FooResource struct {
	ID     int    `json:"id"`
	Foobar string `json:"foobar"`
}

// FooHandler implements the ResourceHandler interface. It specifies the business
// logic for performing CRUD operations.
type FooHandler struct {
	BaseResourceHandler
}

// ResourceName is used to identify what resource a handler corresponds to and is used
// in the endpoint URLs, i.e. /api/:version/foo.
func (f FooHandler) ResourceName() string {
	return "foo"
}

// CreateResource is the logic that corresponds to creating a new resource at
// POST /api/:version/foo. Typically, this would insert a record into a database.
// It returns the newly created resource or an error if the create failed.
func (f FooHandler) CreateResource(ctx RequestContext, data Payload,
	version string) (Resource, error) {
	// Make a database call here.
	id := rand.Int()
	foobar, _ := data.GetString("foobar")
	created := &FooResource{ID: id, Foobar: foobar}
	return created, nil
}

// ReadResource is the logic that corresponds to reading a single resource by its ID at
// GET /api/:version/foo/{id}. Typically, this would make some sort of database query to
// load the resource. If the resource doesn't exist, nil should be returned along with
// an appropriate error.
func (f FooHandler) ReadResource(ctx RequestContext, id string,
	version string) (Resource, error) {
	// Make a database call here.
	if id == "42" {
		return &FooResource{ID: 42, Foobar: "hello world"}, nil
	}
	return nil, ResourceNotFound(fmt.Sprintf("No resource with id %s", id))
}

// ReadResourceList is the logic that corresponds to reading multiple resources, perhaps
// with specified query parameters accessed through the RequestContext. This is
// mapped to GET /api/:version/foo. Typically, this would make some sort of database query
// to fetch the resources. It returns the slice of results, a cursor (or empty) string,
// and error (or nil).
func (f FooHandler) ReadResourceList(ctx RequestContext, limit int,
	cursor string, version string) ([]Resource, string, error) {
	// Make a database call here.
	resources := make([]Resource, 0, limit)
	resources = append(resources, &FooResource{ID: 1, Foobar: "hello"})
	resources = append(resources, &FooResource{ID: 2, Foobar: "world"})
	return resources, "", nil
}

// UpdateResource is the logic that corresponds to updating an existing resource at
// PUT /api/:version/foo/{id}. Typically, this would make some sort of database update
// call. It returns the updated resource or an error if the update failed.
func (f FooHandler) UpdateResource(ctx RequestContext, id string, data Payload,
	version string) (Resource, error) {
	// Make a database call here.
	updateID, _ := strconv.Atoi(id)
	foobar, _ := data.GetString("foobar")
	foo := &FooResource{ID: updateID, Foobar: foobar}
	return foo, nil
}

// DeleteResource is the logic that corresponds to deleting an existing resource at
// DELETE /api/:version/foo/{id}. Typically, this would make some sort of database
// delete call. It returns the deleted resource or an error if the delete failed.
func (f FooHandler) DeleteResource(ctx RequestContext, id string,
	version string) (Resource, error) {
	// Make a database call here.
	deleteID, _ := strconv.Atoi(id)
	foo := &FooResource{ID: deleteID, Foobar: "Goodbye world"}
	return foo, nil
}

// Authenticate is logic that is used to authenticate requests. The default behavior
// of Authenticate, seen in BaseResourceHandler, always returns nil, meaning
// all requests are authenticated. Returning an error means that the request is
// unauthorized and any error message will be sent back with the response.
func (f FooHandler) Authenticate(r *http.Request) error {
	if secrets, ok := r.Header["Authorization"]; ok {
		if secrets[0] == "secret" {
			return nil
		}
	}

	return UnauthorizedRequest("You shall not pass")
}

// This example shows how to fully implement a basic ResourceHandler for performing
// CRUD operations.
func Example_simpleCrud() {
	api := NewAPI(NewConfiguration())

	// Call RegisterResourceHandler to wire up FooHandler.
	api.RegisterResourceHandler(FooHandler{})

	// We're ready to hit our CRUD endpoints.
	api.Start(":8080")
}
