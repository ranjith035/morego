package com.morego.sdk;

import github.com.ranjith035.morego.proto.v1.Session.CloseSessionRequest;
import github.com.ranjith035.morego.proto.v1.Driver.SwipeRequest;
import github.com.ranjith035.morego.proto.v1.Driver.GetSourceRequest;
import github.com.ranjith035.morego.proto.v1.Driver.GetSourceResponse;
import github.com.ranjith035.morego.proto.v1.Driver.Point;

public class Session {
    private final Device device;
    private final String sessionId;
    private final String appId;

    public Session(Device device, String sessionId, String appId) {
        this.device = device;
        this.sessionId = sessionId;
        this.appId = appId;
    }

    public String getSessionId() {
        return sessionId;
    }

    public void close() {
        CloseSessionRequest req = CloseSessionRequest.newBuilder()
                .setSessionId(this.sessionId)
                .build();
        this.device.sessionStub.closeSession(req);
    }

    public void swipe(int startX, int startY, int endX, int endY, int durationMs) {
        SwipeRequest req = SwipeRequest.newBuilder()
                .setDriverId(this.sessionId)
                .setStart(Point.newBuilder().setX(startX).setY(startY).build())
                .setEnd(Point.newBuilder().setX(endX).setY(endY).build())
                .setDurationMs(durationMs)
                .build();
        this.device.driverStub.swipe(req);
    }

    public String getSource() {
        return getSource("xml");
    }

    public String getSource(String format) {
        GetSourceRequest req = GetSourceRequest.newBuilder()
                .setDriverId(this.sessionId)
                .setFormat(format)
                .build();
        GetSourceResponse resp = this.device.driverStub.getSource(req);
        return resp.getSourceData();
    }

    public Locator locator(String strategy, String selector) {
        return new Locator(this, strategy, selector);
    }

    protected Device getDevice() {
        return device;
    }
}
