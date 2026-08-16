import { Device } from './client';

async function main() {
  console.log("Connecting to Mobile Automation Core Engine on port 50051...");
  const device = await Device.connect("localhost:50051", "a5fb8bad");

  try {
    console.log("Starting a persistent calculator session for live Web Inspector...");
    const session = await device.newSession("com.google.android.calculator");
    console.log(`Session ${session.sessionId} started successfully!`);
    console.log("==================================================================");
    console.log("The session is now ACTIVE for the next 90 seconds.");
    console.log("Please open the Live Web Inspector in your browser at:");
    console.log("👉 http://localhost:8082");
    console.log("==================================================================");

    // Keep session open for 90 seconds
    for (let i = 90; i > 0; i -= 10) {
      console.log(`Session will remain active for another ${i} seconds...`);
      await new Promise((resolve) => setTimeout(resolve, 10000));
    }

    console.log("Closing session...");
  } catch (error) {
    console.error("Session failed:", error);
  } finally {
    device.close();
  }
}

main().catch(console.error);
