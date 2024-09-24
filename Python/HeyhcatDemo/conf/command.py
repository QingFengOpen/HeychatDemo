from typing import List, Optional
from pydantic import BaseModel

# Constants
TYPE_SUB_COMMAND = 1
TYPE_SUB_COMMAND_GROUP = 2
TYPE_STRING = 3
TYPE_NUMBER= 4
TYPE_BOOLEAN = 5
TYPE_USER = 6
TYPE_CHANNEL = 7
TYPE_ROLE = 8
TYPE_SELECT = 9
TYPE_INTEGER  = 10


class CommandInfo(BaseModel):
    id: str
    name: str
    type: int
    options: Optional[List['Options']]


class Options(BaseModel):
    value: str
    name: str
    choices: Optional[List['Options']] = None
    type: int
