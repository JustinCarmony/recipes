I'm going to give you a website URL at the end of this prompt. I want you to write me a Golang program that will extract the following information from the website output the follow Yaml defined in this struct:

```golang
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
    Notes []string   `yaml:"note"`
	Ingredients   []string `yaml:"ingredients"`
	Instructions  []string `yaml:"instructions"`
}
```

Website URL: https://www.jocooks.com/wprm_print/instant-pot-mongolian-beef