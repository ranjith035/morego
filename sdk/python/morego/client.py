import grpc
from morego.proto import session_pb2_grpc, session_pb2
from morego.proto import driver_pb2_grpc
from morego.proto import ai_pb2_grpc, ai_pb2
from .session import Session

class Device:
    def __init__(self, channel, device_id):
        self.channel = channel
        self.session_client = session_pb2_grpc.SessionServiceStub(channel)
        self.driver_client = driver_pb2_grpc.DriverServiceStub(channel)
        self.ai_client = ai_pb2_grpc.AIServiceStub(channel)
        self.device_id = device_id

    @staticmethod
    def connect(address: str, device_id: str):
        channel = grpc.insecure_channel(address)
        return Device(channel, device_id)

    def close(self):
        if self.channel:
            self.channel.close()

    def new_session(self, app_id: str, capabilities: dict = None) -> Session:
        req = session_pb2.CreateSessionRequest(
            device_id=self.device_id,
            app_id=app_id,
            capabilities=capabilities or {}
        )
        resp = self.session_client.CreateSession(req)
        return Session(self, resp.session_id, app_id)

    def suggest_locators(self, xml_hierarchy: str) -> list:
        req = ai_pb2.SuggestLocatorsRequest(
            session_id=self.device_id,
            view_hierarchy=xml_hierarchy
        )
        resp = self.ai_client.SuggestLocators(req)
        suggestions = []
        for s in resp.suggestions:
            strategy_name = "UNSPECIFIED"
            if s.locator.strategy == 1:
                strategy_name = "ACCESSIBILITY_ID"
            elif s.locator.strategy == 2:
                strategy_name = "RESOURCE_ID"
            elif s.locator.strategy == 3:
                strategy_name = "TEXT"
            elif s.locator.strategy == 4:
                strategy_name = "CLASS_NAME"
            elif s.locator.strategy == 5:
                strategy_name = "XPATH"

            suggestions.append({
                "strategy": strategy_name,
                "selector": s.locator.selector,
                "stability_score": s.stability_score
            })
        return suggestions
