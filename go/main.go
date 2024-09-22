package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"gopkg.in/yaml.v3"
)

// Recipe struct to store the extracted recipe information
type Recipe struct {
	Title         string   `yaml:"title"`
	Author        string   `yaml:"author"`
	SourceURL     string   `yaml:"source_url"`
	ImageURL      string   `yaml:"image_url"`
	Servings      string   `yaml:"servings"`
	PrepTime      string   `yaml:"prep_time"`
	CookTime      string   `yaml:"cook_time"`
	TotalTime     string   `yaml:"total_time"`
	NutritionInfo string   `yaml:"nutrition_info"`
	Notes         []string `yaml:"note"`
	Ingredients   []string `yaml:"ingredients"`
	Instructions  []string `yaml:"instructions"`
}

// ExtractTextFromNode retrieves the text content of an HTML node and its children
func ExtractTextFromNode(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(ExtractTextFromNode(c))
	}
	return sb.String()
}

// FindElementByTagAndClass recursively searches for elements with a specific tag name and class name
func FindElementByTagAndClass(n *html.Node, tagName, className string) []*html.Node {
	var result []*html.Node
	if n.Type == html.ElementNode && n.Data == tagName {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, className) {
				result = append(result, n)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result = append(result, FindElementByTagAndClass(c, tagName, className)...)
	}
	return result
}

// ParseRecipe parses a recipe from the HTML document and returns a Recipe struct
func ParseRecipe(doc *html.Node, sourceURL string) Recipe {
	var recipe Recipe
	recipe.SourceURL = sourceURL

	// Extract title
	titles := FindElementByTagAndClass(doc, "h2", "wprm-recipe-name")
	if len(titles) > 0 {
		recipe.Title = ExtractTextFromNode(titles[0])
	}

	// Extract author
	authors := FindElementByTagAndClass(doc, "span", "wprm-recipe-author")
	if len(authors) > 0 {
		recipe.Author = ExtractTextFromNode(authors[0])
	}

	// Extract image URL
	imageNodes := FindElementByTagAndClass(doc, "div", "wprm-recipe-image")
	if len(imageNodes) > 0 {
		img := imageNodes[0].FirstChild
		for img != nil && img.Data != "img" {
			img = img.NextSibling
		}
		if img != nil {
			for _, attr := range img.Attr {
				if attr.Key == "src" {
					recipe.ImageURL = attr.Val
					break
				}
			}
		}
	}

	// Extract servings, prep time, cook time, and total time
	recipe.Servings = strings.TrimSpace(ExtractTextFromNode(FindElementByTagAndClass(doc, "span", "wprm-recipe-servings")[0]))
	recipe.PrepTime = strings.TrimSpace(ExtractTextFromNode(FindElementByTagAndClass(doc, "span", "wprm-recipe-prep_time")[0]))
	recipe.CookTime = strings.TrimSpace(ExtractTextFromNode(FindElementByTagAndClass(doc, "span", "wprm-recipe-cook_time")[0]))
	recipe.TotalTime = strings.TrimSpace(ExtractTextFromNode(FindElementByTagAndClass(doc, "span", "wprm-recipe-total_time")[0]))

	// Extract nutrition info
	nutritionInfo := FindElementByTagAndClass(doc, "div", "wprm-nutrition-label")
	if len(nutritionInfo) > 0 {
		recipe.NutritionInfo = ExtractTextFromNode(nutritionInfo[0])
	}

	// Extract notes
	notes := FindElementByTagAndClass(doc, "div", "wprm-recipe-notes-container")
	for _, noteNode := range notes {
		recipe.Notes = append(recipe.Notes, strings.TrimSpace(ExtractTextFromNode(noteNode)))
	}

	// Extract ingredients
	ingredientNodes := FindElementByTagAndClass(doc, "li", "wprm-recipe-ingredient")
	for _, node := range ingredientNodes {
		ingredient := ExtractTextFromNode(node)
		recipe.Ingredients = append(recipe.Ingredients, strings.TrimSpace(ingredient))
	}

	// Extract instructions
	instructionNodes := FindElementByTagAndClass(doc, "div", "wprm-recipe-instruction-text")
	for _, node := range instructionNodes {
		instruction := ExtractTextFromNode(node)
		recipe.Instructions = append(recipe.Instructions, strings.TrimSpace(instruction))
	}

	return recipe
}

// FetchHTML retrieves the HTML content of a webpage
func FetchHTML(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %v", err)
	}
	return doc, nil
}

func main() {
	url := "https://www.paintthekitchenred.com/wprm_print/instant-pot-thai-green-curry-with-chicken"
	doc, err := FetchHTML(url)
	if err != nil {
		log.Fatalf("Failed to fetch and parse recipe: %v", err)
	}

	recipe := ParseRecipe(doc, url)

	// Convert the recipe struct to YAML
	yamlData, err := yaml.Marshal(&recipe)
	if err != nil {
		log.Fatalf("Failed to marshal recipe to YAML: %v", err)
	}

	// Print the YAML output
	fmt.Println(string(yamlData))
}
