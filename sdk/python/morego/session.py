from morego.proto import session_pb2
from morego.proto import driver_pb2
from .locator import Locator

class Session:
    def __init__(self, device, session_id, app_id):
        self.device = device
        self.session_id = session_id
        self.app_id = app_id

    def close(self):
        req = session_pb2.CloseSessionRequest(session_id=self.session_id)
        self.device.session_client.CloseSession(req)

    def swipe(self, start_x: int, start_y: int, end_x: int, end_y: int, duration_ms: int):
        req = driver_pb2.SwipeRequest(
            driver_id=self.session_id,
            start=driver_pb2.Point(x=start_x, y=start_y),
            end=driver_pb2.Point(x=end_x, y=end_y),
            duration_ms=duration_ms
        )
        self.device.driver_client.Swipe(req)

    def get_source(self, format_type: str = "xml") -> str:
        req = driver_pb2.GetSourceRequest(
            driver_id=self.session_id,
            format=format_type
        )
        resp = self.device.driver_client.GetSource(req)
        return resp.source_data

    def locator(self, strategy: str, selector: str) -> Locator:
        return Locator(self, strategy, selector)
