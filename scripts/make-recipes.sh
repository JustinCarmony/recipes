#!/bin/bash

docker build -t recipes-latex .

for recipe_file in markdown/*.md; do
    recipe_slug=$(basename "$recipe_file" .md)
    recipe_title=$(yq --front-matter extract eval '.title' "$recipe_file")
    template_name=$(yq --front-matter extract eval '.template' "$recipe_file")

    # If no template is specified, default to basic template
    if [ "$template_name" == "null" ]; then
        template_name="basic"
    fi

    template_file="templates/$template_name.tex"
    tex_file="latex/$recipe_slug.tex"
    pdf_file="pdfs/$recipe_slug.pdf"

    echo "Recipe Slug: $recipe_slug"
    echo "Creating LaTeX file for $recipe_title using $template_name template"

    # Convert Markdown to LaTeX using Pandoc
    docker run --rm \
        -v "$(pwd)":/data \
        --user $(id -u):$(id -g) \
        pandoc/extra \
        "/data/$recipe_file" \
        --template="/data/$template_file" \
        -o "/data/$tex_file"

    echo "Creating PDF from LaTeX file using XeLaTeX"
    
    # Run XeLaTeX to convert LaTeX to PDF
    docker run --rm \
        -v "$(pwd)":/data \
        --user $(id -u):$(id -g) \
        recipes-latex \
        xelatex -output-directory="/data/pdfs" "/data/$tex_file"

    echo "   PDF created at $pdf_file"
done