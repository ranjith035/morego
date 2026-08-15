import time
from morego import Device

def main():
    print("Connecting to morego Core Server...")
    device = Device.connect("localhost:50051", "pixel_6_pro")
    
    try:
        print("Creating automation session...")
        session = device.new_session("com.android.settings")
        print(f"Session established with ID: {session.session_id}")
        
        print("Executing swipe gesture...")
        session.swipe(500, 1500, 500, 500, 400)
        
        print("Locating search bar...")
        # Since we use adb driver, we query search bar
        search_bar = session.locator("RESOURCE_ID", "com.android.settings:id/search_action_bar")
        
        # We can try to click (will require actual device connected to pass click resolution, 
        # or we catch connection errors gracefully if no physical adb device is plugged in!)
        try:
            search_bar.click()
            print("Click sent!")
        except Exception as e:
            print(f"Interaction details (expected if no physical adb hardware attached): {e}")
            
    finally:
        print("Closing device connection...")
        device.close()

if __name__ == "__main__":
    main()
