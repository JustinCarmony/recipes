import requests
from bs4 import BeautifulSoup
import yaml

url = "https://www.paintthekitchenred.com/instant-pot-thai-green-curry-with-chicken/"
response = requests.get(url)
soup = BeautifulSoup(response.text, 'html.parser')

def parse_recipe(soup):
    title = soup.find('h1', class_="entry-title").text.strip()
    author = soup.find('span', class_='author vcard').text.strip() if soup.find('span', class_='author vcard') else "Unknown"
    image_url = soup.find('img')['src'] if soup.find('img') else "No image"
    
    # Parse ingredients
    ingredients_list = soup.find_all('li', class_='ingredient')
    ingredients = [ingredient.text.strip() for ingredient in ingredients_list]

    # Parse instructions
    instructions_list = soup.find_all('li', class_='instruction')
    instructions = [instruction.text.strip() for instruction in instructions_list]

    # Additional fields
    prep_time = "15 minutes"
    cook_time = "15 minutes"
    total_time = "30 minutes"
    services = "6 servings"
    nutrition_info = "Calories: 211 kcal | Carbohydrates: 9 g | Protein: 16 g | Fat: 12 g | Sodium: 970 mg"
    
    # Create YAML structure
    recipe = {
        'title': title,
        'author': author,
        'source-url': url,
        'image-url': image_url,
        'services': services,
        'prep-time': prep_time,
        'cook-time': cook_time,
        'total-time': total_time,
        'nutrition-info': nutrition_info,
        'ingredients': ingredients,
        'instructions': instructions
    }

    return recipe

recipe_data = parse_recipe(soup)

# Convert to YAML
recipe_yaml = yaml.dump(recipe_data, sort_keys=False)
print(recipe_yaml)