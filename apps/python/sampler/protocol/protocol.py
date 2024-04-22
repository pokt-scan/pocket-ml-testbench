from typing import List, Literal, Optional
from pydantic import BaseModel, field_validator

class PocketNetworkRegisterTaskRequest(BaseModel):
    evaluation: Literal["lmeh", "helm"]
    tasks: str
    verbosity: Optional[Literal["CRITICAL", "ERROR", "WARNING", "INFO", "DEBUG"]] = "INFO"
    include_path: Optional[str] = None
    postgres_uri: Optional[str] = None
    mongodb_uri: Optional[str] = None

class RequesterArgs(BaseModel):
    address: str
    service: str
    method: str
    path: str
    
class PocketNetworkTaskRequest(PocketNetworkRegisterTaskRequest):
    requester_args: RequesterArgs
    blacklist: Optional[List[int]] = []
    qty: int
    doc_ids: Optional[List[int]] = []
    # assert that "qty" is greater than 0
    @field_validator("qty")
    def check_qty(cls, v):
        if v <= 0:
            raise ValueError("qty must be greater than 0")
        return v


class PocketNetworkMongoDBTask(BaseModel):
    evaluation: Literal["lmeh", "helm"]
    requester_args: RequesterArgs
    blacklist: Optional[List[int]] = []
    qty: int
    tasks: str
    total_instances: int
    request_type: str