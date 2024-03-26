#!/usr/bin/env python
from lib.core.enums import PRIORITY
import re

__priority__ = PRIORITY.NORMAL

def dependencies():
    pass
def tamper(payload, **kwargs):
    retVal = ""    
    retVal = re.sub('\\bOR\\b', '||', payload)    
    retVal = re.sub('\\bAND\\b', '&&', retVal)    
    return retVal
