# Makefile

# Variables
SCRIPTS_DIR = scripts
SCRIPT = make-recipes.sh

# Default target
recipes:
	bash $(SCRIPTS_DIR)/$(SCRIPT)