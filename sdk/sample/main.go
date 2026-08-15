package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ranjith035/morego/sdk"
)

func main() {
	ctx := context.Background()

	fmt.Println("Connecting to Mobile Automation Core Engine on port 50051...")
	// 1. Connect to gRPC core server
	device, err := sdk.Connect(ctx, "localhost:50051", "pixel_6_pro")
	if err != nil {
		fmt.Printf("Connection error: %v\n", err)
		return
	}
	defer device.Close()

	fmt.Println("Starting execution session on target app...")
	// 2. Start session (e.g. launching standard Google Calculator or Settings activity)
	session, err := device.NewSession(ctx, "com.android.settings", nil)
	if err != nil {
		fmt.Printf("Failed to create session: %v\n", err)
		return
	}
	defer session.Close(ctx)

	fmt.Printf("Session %s started successfully!\n", session.ID())

	// 3. Perform a simple swipe gesture
	fmt.Println("Executing swipe gesture...")
	err = session.Swipe(ctx, 500, 1500, 500, 500, 500*time.Millisecond)
	if err != nil {
		fmt.Printf("Swipe failed: %v\n", err)
	}

	// 4. Try locating and clicking the Search bar (resource id in settings app is usually search_action_bar)
	fmt.Println("Clicking Search Bar...")
	searchBar := session.Locator("RESOURCE_ID", "com.android.settings:id/search_action_bar")
	err = searchBar.Click(ctx)
	if err != nil {
		fmt.Printf("Could not click search bar (perhaps ID shifted): %v\n", err)
	} else {
		fmt.Println("Search bar clicked! Typing search query...")
		// 5. Fill input
		searchVal := session.Locator("CLASS_NAME", "android.widget.EditText")
		err = searchVal.Fill(ctx, "Wi-Fi")
		if err != nil {
			fmt.Printf("Failed to enter text: %v\n", err)
		} else {
			fmt.Println("Text 'Wi-Fi' typed successfully into Settings search field!")
		}
	}

	fmt.Println("Releasing session connections...")
}
