from morego.proto import driver_pb2
from morego.proto import locator_pb2

def map_strategy(strategy_str: str) -> int:
    s = strategy_str.upper()
    if s in ("ACCESSIBILITY_ID", "TEST_ID"):
        return locator_pb2.LocatorStrategy.LOCATOR_STRATEGY_ACCESSIBILITY_ID
    if s == "RESOURCE_ID":
        return locator_pb2.LocatorStrategy.LOCATOR_STRATEGY_RESOURCE_ID
    if s == "TEXT":
        return locator_pb2.LocatorStrategy.LOCATOR_STRATEGY_TEXT
    if s in ("CLASS_NAME", "CLASS"):
        return locator_pb2.LocatorStrategy.LOCATOR_STRATEGY_CLASS_NAME
    if s == "XPATH":
        return locator_pb2.LocatorStrategy.LOCATOR_STRATEGY_XPATH
    return locator_pb2.LocatorStrategy.LOCATOR_STRATEGY_UNSPECIFIED

class Locator:
    def __init__(self, session, strategy: str, selector: str, index: int = 0):
        self.session = session
        self.strategy = strategy
        self.selector = selector
        self.index = index
        self.modifiers = []

    def _resolve(self) -> str:
        req = driver_pb2.FindElementRequest(
            session_id=self.session.session_id,
            locator=locator_pb2.Locator(
				strategy=map_strategy(self.strategy),
				selector=self.selector
			)
        )
        resp = self.session.device.driver_client.FindElement(req)
        if not resp.found:
            raise Exception(f"Element not found using locator: {self.strategy}={self.selector}")
        return resp.element.element_id

    def click(self):
        el_id = self._resolve()
        req = driver_pb2.ClickRequest(
            driver_id=self.session.session_id,
            element_id=el_id
        )
        self.session.device.driver_client.Click(req)

    def fill(self, value: str):
        el_id = self._resolve()
        req = driver_pb2.FillRequest(
            driver_id=self.session.session_id,
            element_id=el_id,
            value=value
        )
        self.session.device.driver_client.Fill(req)

    def screenshot(self) -> bytes:
        el_id = self._resolve()
        req = driver_pb2.ScreenshotRequest(
            driver_id=self.session.session_id,
            element_id=el_id
        )
        resp = self.session.device.driver_client.Screenshot(req)
        return resp.image_data

    def above(self, target):
        self.modifiers.append(("ABOVE", target))
        return self

    def below(self, target):
        self.modifiers.append(("BELOW", target))
        return self

    def left_of(self, target):
        self.modifiers.append(("LEFT_OF", target))
        return self

    def right_of(self, target):
        self.modifiers.append(("RIGHT_OF", target))
        return self

    def nth(self, index: int):
        return Locator(self.session, self.strategy, self.selector, index)
