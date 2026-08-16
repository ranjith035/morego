package drivers

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// BaseStub provides shared placeholder structures for Driver stubs.
type BaseStub struct {
	Name string
}

func (b *BaseStub) Connect(ctx context.Context, capabilities map[string]string) error {
	return nil
}

func (b *BaseStub) Disconnect(ctx context.Context) error {
	return nil
}

func (b *BaseStub) Click(ctx context.Context, elementID string) error {
	return nil
}

func (b *BaseStub) ClickAt(ctx context.Context, x, y int) error {
	return nil
}

func (b *BaseStub) Fill(ctx context.Context, elementID string, value string) error {
	return nil
}

func (b *BaseStub) Swipe(ctx context.Context, startX, startY, endX, endY int, duration time.Duration) error {
	return nil
}

func (b *BaseStub) Screenshot(ctx context.Context, elementID string) ([]byte, error) {
	return []byte("screenshot_bytes"), nil
}

func (b *BaseStub) FindElement(ctx context.Context, strategy string, selector string) (string, error) {
	return "stub_element_id", nil
}

func (b *BaseStub) LaunchApp(ctx context.Context, appID string, args []string, env map[string]string) error {
	return nil
}

func (b *BaseStub) TerminateApp(ctx context.Context, appID string) error {
	return nil
}

func (b *BaseStub) InstallApp(ctx context.Context, appPath string) error {
	return nil
}

func (b *BaseStub) UninstallApp(ctx context.Context, appID string) error {
	return nil
}

func (b *BaseStub) GetSource(ctx context.Context, format string) (string, error) {
	return "<xml>source</xml>", nil
}

func (b *BaseStub) ExecuteScript(ctx context.Context, script string, arguments []string) (string, error) {
	return "result", nil
}

// Concrete Driver Stubs

type UiAutomator2Driver struct{ BaseStub }
type EspressoDriver struct{ BaseStub }

func init() {
	DefaultRegistry.Register("adb", func() Driver {
		return NewADBDriver()
	})
	DefaultRegistry.Register("uiautomator2", func() Driver {
		return &UiAutomator2Driver{BaseStub{Name: "UiAutomator2 Driver"}}
	})
	DefaultRegistry.Register("espresso", func() Driver {
		return &EspressoDriver{BaseStub{Name: "Espresso Driver"}}
	})
	DefaultRegistry.Register("xcuitest", func() Driver {
		return NewXCUITestDriver()
	})
	DefaultRegistry.Register("seetest", func() Driver {
		return NewSeeTestDriver()
	})
	DefaultRegistry.Register("browserstack", func() Driver {
		return NewCloudGridDriver("browserstack")
	})
	DefaultRegistry.Register("saucelabs", func() Driver {
		return NewCloudGridDriver("saucelabs")
	})
	DefaultRegistry.Register("lambdatest", func() Driver {
		return NewCloudGridDriver("lambdatest")
	})
	DefaultRegistry.Register("rich_stub", func() Driver {
		return &RichStubDriver{
			BaseStub: BaseStub{Name: "Rich Stub Driver"},
		}
	})
}

type RichStubDriver struct {
	BaseStub
	mu    sync.Mutex
	cache map[string]ADBXMLNode
	idGen int
}

func (r *RichStubDriver) GetSource(ctx context.Context, format string) (string, error) {
	return `<?xml version="1.0" encoding="utf-8"?>
	<hierarchy rotation="0">
		<node class="android.widget.FrameLayout" bounds="[0,0][1080,1920]" resource-id="root">
			<node class="android.widget.Button" text="Submit" resource-id="com.example:id/submit_btn" content-desc="Submit description" bounds="[100,200][300,300]" id="btn_1"/>
			<node class="android.widget.EditText" text="" placeholder="Enter Username" bounds="[100,400][300,500]" test-id="username_field" id="edit_1"/>
			<node class="android.widget.TextView" text="Username Label" bounds="[100,100][300,150]" id="label_1"/>
		</node>
	</hierarchy>`, nil
}

func (r *RichStubDriver) FindElement(ctx context.Context, strategy string, selector string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cache == nil {
		r.cache = make(map[string]ADBXMLNode)
	}
	r.idGen++
	elemID := fmt.Sprintf("rich_elem_%d", r.idGen)
	// Add dummy bounds representation
	r.cache[elemID] = ADBXMLNode{Bounds: selector}
	return elemID, nil
}


// Helper to check if a driver type is registered.
func IsDriverAvailable(driverType string) bool {
	_, err := DefaultRegistry.CreateDriver(driverType)
	return err == nil
}
