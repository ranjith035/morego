package com.morego.sdk.sample;

import com.morego.sdk.Device;
import com.morego.sdk.Session;
import com.morego.sdk.Locator;

public class Sample {
    public static void main(String[] args) {
        System.out.println("Connecting to morego Core Server...");
        try (Device device = Device.connect("localhost:50051", "pixel_6_pro")) {
            System.out.println("Creating automation session...");
            Session session = device.newSession("com.android.settings");
            System.out.println("Session established with ID: " + session.getSessionId());

            System.out.println("Executing swipe gesture...");
            session.swipe(500, 1500, 500, 500, 400);

            System.out.println("Locating search bar...");
            Locator searchBar = session.locator("RESOURCE_ID", "com.android.settings:id/search_action_bar");

            try {
                searchBar.click();
                System.out.println("Click sent!");
            } catch (Exception e) {
                System.out.println("Interaction details (expected if no physical adb hardware attached): " + e.getMessage());
            }
        } catch (Exception e) {
            System.out.println("Session error: " + e.getMessage());
        }
    }
}
