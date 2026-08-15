package com.morego.sdk;

import github.com.ranjith035.morego.proto.v1.Driver.FindElementRequest;
import github.com.ranjith035.morego.proto.v1.Driver.FindElementResponse;
import github.com.ranjith035.morego.proto.v1.Driver.ClickRequest;
import github.com.ranjith035.morego.proto.v1.Driver.FillRequest;
import github.com.ranjith035.morego.proto.v1.Driver.ScreenshotRequest;
import github.com.ranjith035.morego.proto.v1.Driver.ScreenshotResponse;
import github.com.ranjith035.morego.proto.v1.LocatorOuterClass.LocatorStrategy;

public class Locator {
    private final Session session;
    private final String strategy;
    private final String selector;

    public Locator(Session session, String strategy, String selector) {
        this.session = session;
        this.strategy = strategy;
        this.selector = selector;
    }

    private String resolve() {
        FindElementRequest req = FindElementRequest.newBuilder()
                .setSessionId(session.getSessionId())
                .setLocator(github.com.ranjith035.morego.proto.v1.LocatorOuterClass.Locator.newBuilder()
                        .setStrategy(mapStrategy(strategy))
                        .setSelector(selector)
                        .build())
                .build();
        FindElementResponse resp = session.getDevice().driverStub.findElement(req);
        if (!resp.getFound()) {
            throw new RuntimeException("Element not found: " + strategy + "=" + selector);
        }
        return resp.getElement().getElementId();
    }

    public void click() {
        String elementId = resolve();
        ClickRequest req = ClickRequest.newBuilder()
                .setDriverId(session.getSessionId())
                .setElementId(elementId)
                .build();
        session.getDevice().driverStub.click(req);
    }

    public void fill(String value) {
        String elementId = resolve();
        FillRequest req = FillRequest.newBuilder()
                .setDriverId(session.getSessionId())
                .setElementId(elementId)
                .setValue(value)
                .build();
        session.getDevice().driverStub.fill(req);
    }

    public byte[] screenshot() {
        String elementId = resolve();
        ScreenshotRequest req = ScreenshotRequest.newBuilder()
                .setDriverId(session.getSessionId())
                .setElementId(elementId)
                .build();
        ScreenshotResponse resp = session.getDevice().driverStub.screenshot(req);
        return resp.getImageData().toByteArray();
    }

    private LocatorStrategy mapStrategy(String strategyStr) {
        String s = strategyStr.toUpperCase();
        switch (s) {
            case "ACCESSIBILITY_ID":
            case "TEST_ID":
                return LocatorStrategy.LOCATOR_STRATEGY_ACCESSIBILITY_ID;
            case "RESOURCE_ID":
                return LocatorStrategy.LOCATOR_STRATEGY_RESOURCE_ID;
            case "TEXT":
                return LocatorStrategy.LOCATOR_STRATEGY_TEXT;
            case "CLASS_NAME":
            case "CLASS":
                return LocatorStrategy.LOCATOR_STRATEGY_CLASS_NAME;
            case "XPATH":
                return LocatorStrategy.LOCATOR_STRATEGY_XPATH;
            default:
                return LocatorStrategy.LOCATOR_STRATEGY_UNSPECIFIED;
        }
    }
}
