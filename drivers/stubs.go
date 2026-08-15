package drivers

import (
	"context"
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
type SeeTestDriver struct{ BaseStub }

type CloudGridDriver struct {
	BaseStub
	Provider string
}

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
		return &SeeTestDriver{BaseStub{Name: "SeeTest Driver"}}
	})
	DefaultRegistry.Register("browserstack", func() Driver {
		return &CloudGridDriver{BaseStub: BaseStub{Name: "BrowserStack Driver"}, Provider: "browserstack"}
	})
	DefaultRegistry.Register("saucelabs", func() Driver {
		return &CloudGridDriver{BaseStub: BaseStub{Name: "SauceLabs Driver"}, Provider: "saucelabs"}
	})
	DefaultRegistry.Register("lambdatest", func() Driver {
		return &CloudGridDriver{BaseStub: BaseStub{Name: "LambdaTest Driver"}, Provider: "lambdatest"}
	})
}

// Helper to check if a driver type is registered.
func IsDriverAvailable(driverType string) bool {
	_, err := DefaultRegistry.CreateDriver(driverType)
	return err == nil
}
