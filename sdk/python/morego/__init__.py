import sys
import os

# Seamlessly map internal proto references to morego namespace
current_dir = os.path.dirname(os.path.abspath(__file__))
parent_dir = os.path.dirname(current_dir)
if parent_dir not in sys.path:
    sys.path.insert(0, parent_dir)

import morego.proto
sys.modules['proto'] = sys.modules['morego.proto']

from .client import Device
from .session import Session
from .locator import Locator
