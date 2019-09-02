package goshoplazza

import (
	"fmt"
	"time"
)

// ImageService is an interface for interacting with the image endpoints
// of the Shopify API.
// See https://help.shopify.com/api/reference/product_image
type ImageService interface {
	List(int64, interface{}) ([]Image, error)
	Count(int64, interface{}) (int, error)
	Get(int64, int64, interface{}) (*Image, error)
	Create(int64, Image) (*Image, error)
	Update(int64, Image) (*Image, error)
	Delete(int64, int64) error
}

// ImageServiceOp handles communication with the image related methods of
// the Shopify API.
type ImageServiceOp struct {
	client *Client
}

// Image represents a Shopify product's image.
type Image struct {
	ID         string     `json:"id,omitempty"`
	ProductID  string     `json:"product_id,omitempty"`
	Position   int        `json:"position,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	Width      interface{}      `json:"width,omitempty"` // 有两种类型的值 220 和 "220", 使用时注意asset类型并做转换
	Height     interface{}      `json:"height,omitempty"`// 有两种类型的值 220 和 "220", 使用时注意asset类型并做转换
	Src        string     `json:"src,omitempty"`
	Alt        string     `json:"alt,omitempty"`
}

// ImageResource represents the result form the products/X/images/Y.json endpoint
type ImageResource struct {
	Image *Image `json:"image"`
}

// ImagesResource represents the result from the products/X/images.json endpoint
type ImagesResource struct {
	Images []Image `json:"images"`
}

// List images
func (s *ImageServiceOp) List(productID int64, options interface{}) ([]Image, error) {
	path := fmt.Sprintf("%s/%s/%d/images.json", globalApiPathPrefix, productsBasePath, productID)
	resource := new(ImagesResource)
	err := s.client.Get(path, resource, options)
	return resource.Images, err
}

// Count images
func (s *ImageServiceOp) Count(productID int64, options interface{}) (int, error) {
	path := fmt.Sprintf("%s/%s/%d/images/count.json", globalApiPathPrefix, productsBasePath, productID)
	return s.client.Count(path, options)
}

// Get individual image
func (s *ImageServiceOp) Get(productID int64, imageID int64, options interface{}) (*Image, error) {
	path := fmt.Sprintf("%s/%s/%d/images/%d.json", globalApiPathPrefix, productsBasePath, productID, imageID)
	resource := new(ImageResource)
	err := s.client.Get(path, resource, options)
	return resource.Image, err
}

// Create a new image
//
// There are 2 methods of creating an image in Shopify:
// 1. Src
// 2. Filename and Attachment
//
// If both Image.Filename and Image.Attachment are supplied,
// then Image.Src is not needed.  And vice versa.
//
// If both Image.Attachment and Image.Src are provided,
// Shopify will take the attachment.
//
// Shopify will accept Image.Attachment without Image.Filename.
func (s *ImageServiceOp) Create(productID int64, image Image) (*Image, error) {
	path := fmt.Sprintf("%s/%s/%d/images.json", globalApiPathPrefix, productsBasePath, productID)
	wrappedData := ImageResource{Image: &image}
	resource := new(ImageResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Image, err
}

// Update an existing image
func (s *ImageServiceOp) Update(productID int64, image Image) (*Image, error) {
	path := fmt.Sprintf("%s/%s/%d/images/%d.json", globalApiPathPrefix, productsBasePath, productID, image.ID)
	wrappedData := ImageResource{Image: &image}
	resource := new(ImageResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.Image, err
}

// Delete an existing image
func (s *ImageServiceOp) Delete(productID int64, imageID int64) error {
	return s.client.Delete(fmt.Sprintf("%s/%s/%d/images/%d.json", globalApiPathPrefix, productsBasePath, productID, imageID))
}