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
	device, err := sdk.Connect(ctx, "localhost:50051", "a5fb8bad")
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

	// 3. Locate and click search bar using the new GetByText API first
	fmt.Println("Clicking Search Bar...")
	searchBar := session.GetByText("Search settings")
	err = searchBar.Click(ctx)
	if err != nil {
		fmt.Printf("Could not click search bar: %v\n", err)
	} else {
		fmt.Println("Search bar clicked! Typing search query...")
		// 4. Fill search input using its standard resource-id
		searchVal := session.Locator("RESOURCE_ID", "android:id/input")
		err = searchVal.Fill(ctx, "developer")
		if err != nil {
			fmt.Printf("Failed to enter text: %v\n", err)
		} else {
			fmt.Println("Text 'developer' typed successfully into Settings search field!")
			
			// Allow search results to load
			time.Sleep(2 * time.Second)

			// 5. Perform a swipe gesture on the search results screen
			fmt.Println("Executing swipe gesture on search results...")
			err = session.Swipe(ctx, 500, 1500, 500, 500, 500*time.Millisecond)
			if err != nil {
				fmt.Printf("Swipe failed: %v\n", err)
			}
		}
	}

	fmt.Println("Releasing session connections...")
}
