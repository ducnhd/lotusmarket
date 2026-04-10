"""lotusmarket — Vietnamese Stock Market Toolkit."""

import os

__version__ = "0.1.1"

if not os.environ.get("LOTUSMARKET_QUIET"):
    print(f"lotusmarket v{__version__} — Vietnamese Market Toolkit")
