const { Device } = require('./dist');

async function main() {
  console.log("Connecting to morego Core Server...");
  const device = await Device.connect("localhost:50051", "pixel_6_pro");

  try {
    console.log("Creating automation session...");
    const session = await device.newSession("com.android.settings");
    console.log("Session established with ID:", session.sessionId);

    console.log("Executing swipe gesture...");
    await session.swipe(500, 1500, 500, 500, 400);

    console.log("Locating search bar...");
    const searchBar = session.locator("RESOURCE_ID", "com.android.settings:id/search_action_bar");
    
    try {
      await searchBar.click();
      console.log("Click sent!");
    } catch (e) {
      console.log("Interaction details (expected if no physical adb hardware attached):", e.message);
    }
  } finally {
    console.log("Closing device connection...");
    device.close();
  }
}

main().catch(console.error);
