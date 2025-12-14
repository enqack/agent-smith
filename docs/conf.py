# Configuration file for the Sphinx documentation builder.

import os
import sys

# -- Project information -----------------------------------------------------
project = 'Agent Smith'
copyright = '2025, Enqack'
author = 'Enqack'

# Read version from VERSION file
with open(os.path.join(os.path.dirname(__file__), '..', 'VERSION')) as f:
    release = f.read().strip()

# -- General configuration ---------------------------------------------------
extensions = [
    'myst_parser',
    'sphinx_copybutton',
]

templates_path = ['_templates']
exclude_patterns = ['_build', 'Thumbs.db', '.DS_Store', 'README.md']

# -- Options for HTML output -------------------------------------------------
html_theme = 'furo'
html_static_path = []

# -- MyST Parser configuration -----------------------------------------------
myst_enable_extensions = [
    "colon_fence",
    "deflist",
    "fieldlist",
]
myst_heading_anchors = 3
