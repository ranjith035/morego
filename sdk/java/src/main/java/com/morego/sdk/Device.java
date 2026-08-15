package com.morego.sdk;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import github.com.ranjith035.morego.proto.v1.SessionServiceGrpc;
import github.com.ranjith035.morego.proto.v1.DriverServiceGrpc;
import github.com.ranjith035.morego.proto.v1.AIServiceGrpc;
import github.com.ranjith035.morego.proto.v1.Session.CreateSessionRequest;
import github.com.ranjith035.morego.proto.v1.Session.CreateSessionResponse;

import java.util.HashMap;
import java.util.Map;

public class Device implements AutoCloseable {
    private final ManagedChannel channel;
    protected final SessionServiceGrpc.SessionServiceBlockingStub sessionStub;
    protected final DriverServiceGrpc.DriverServiceBlockingStub driverStub;
    protected final AIServiceGrpc.AIServiceBlockingStub aiStub;
    private final String deviceId;

    public Device(ManagedChannel channel, String deviceId) {
        this.channel = channel;
        this.sessionStub = SessionServiceGrpc.newBlockingStub(channel);
        this.driverStub = DriverServiceGrpc.newBlockingStub(channel);
        this.aiStub = AIServiceGrpc.newBlockingStub(channel);
        this.deviceId = deviceId;
    }

    public static Device connect(String address, String deviceId) {
        ManagedChannel channel = ManagedChannelBuilder.forTarget(address)
                .usePlaintext()
                .build();
        return new Device(channel, deviceId);
    }

    public Session newSession(String appId) {
        return newSession(appId, new HashMap<>());
    }

    public Session newSession(String appId, Map<String, String> capabilities) {
        CreateSessionRequest req = CreateSessionRequest.newBuilder()
                .setDeviceId(this.deviceId)
                .setAppId(appId)
                .putAllCapabilities(capabilities)
                .build();
        CreateSessionResponse resp = this.sessionStub.createSession(req);
        return new Session(this, resp.getSessionId(), appId);
    }

    @Override
    public void close() {
        if (channel != null) {
            channel.shutdown();
        }
    }
}
