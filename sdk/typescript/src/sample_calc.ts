import { Device } from './client';

async function main() {
  console.log("Connecting to Mobile Automation Core Engine on port 50051...");
  const device = await Device.connect("localhost:50051", "a5fb8bad");

  try {
    console.log("Starting calculator session...");
    const session = await device.newSession("com.google.android.calculator");
    console.log(`Session ${session.sessionId} started successfully!`);

    // Let's add a short sleep for screen to load
    await new Promise((resolve) => setTimeout(resolve, 2000));

    console.log("Clicking digit 7 using accessibility ID...");
    const digit7 = session.getByAccessibilityID("7");
    await digit7.click();

    console.log("Clicking plus (+) using accessibility ID...");
    const plus = session.getByAccessibilityID("plus");
    await plus.click();

    console.log("Clicking digit 8 using accessibility ID...");
    const digit8 = session.getByAccessibilityID("8");
    await digit8.click();

    console.log("Clicking equals (=) using accessibility ID...");
    const equals = session.getByAccessibilityID("equals");
    await equals.click();

    console.log("Calculation complete!");
    // Allow state to settle
    await new Promise((resolve) => setTimeout(resolve, 2000));

    console.log("Fetching screen source XML...");
    const src = await session.getSource();
    // On different calculators, the result container has text "15"
    if (src.includes("15")) {
      console.log("SUCCESS: Screen contains expected result '15'!");
    } else {
      console.log("Result '15' not found in layout XML. Current XML length: " + src.length);
    }

  } catch (error) {
    console.error("Test execution failed:", error);
  } finally {
    device.close();
  }
}

main().catch(console.error);
