package recorder

import (
	"fmt"
	"strings"
)

// TypeScriptGenerator compiles actions to TypeScript SDK code.
type TypeScriptGenerator struct{}

func (g *TypeScriptGenerator) Generate(actions []ActionIR) (string, error) {
	var sb strings.Builder
	sb.WriteString("import { mobile } from '@mobile/playwright';\n\n")
	sb.WriteString("(async () => {\n")
	sb.WriteString("  const device = await mobile.connect({ deviceId: 'device_id' });\n")
	sb.WriteString("  const session = await device.newSession({ appId: 'app_id' });\n\n")

	for _, a := range actions {
		switch a.Type {
		case ActionClick:
			sb.WriteString(fmt.Sprintf("  await session.locator('%s', '%s').click();\n", a.SelectorStrategy, a.Selector))
		case ActionFill:
			sb.WriteString(fmt.Sprintf("  await session.locator('%s', '%s').fill('%s');\n", a.SelectorStrategy, a.Selector, a.Value))
		case ActionSwipe:
			sb.WriteString(fmt.Sprintf("  await session.swipe({ startX: %d, startY: %d, endX: %d, endY: %d, durationMs: %d });\n", a.StartX, a.StartY, a.EndX, a.EndY, a.DurationMS))
		case ActionOpenApp:
			sb.WriteString(fmt.Sprintf("  await session.launchApp('%s');\n", a.Value))
		case ActionCloseApp:
			sb.WriteString(fmt.Sprintf("  await session.terminateApp('%s');\n", a.Value))
		}
	}

	sb.WriteString("\n  await session.close();\n")
	sb.WriteString("})();\n")
	return sb.String(), nil
}

// PythonGenerator compiles actions to Python SDK code.
type PythonGenerator struct{}

func (g *PythonGenerator) Generate(actions []ActionIR) (string, error) {
	var sb strings.Builder
	sb.WriteString("import asyncio\n")
	sb.WriteString("from mobile_playwright import mobile\n\n")
	sb.WriteString("async def main():\n")
	sb.WriteString("    device = await mobile.connect(device_id='device_id')\n")
	sb.WriteString("    session = await device.new_session(app_id='app_id')\n\n")

	for _, a := range actions {
		switch a.Type {
		case ActionClick:
			sb.WriteString(fmt.Sprintf("    await session.locator('%s', '%s').click()\n", a.SelectorStrategy, a.Selector))
		case ActionFill:
			sb.WriteString(fmt.Sprintf("    await session.locator('%s', '%s').fill('%s')\n", a.SelectorStrategy, a.Selector, a.Value))
		case ActionSwipe:
			sb.WriteString(fmt.Sprintf("    await session.swipe(start_x=%d, start_y=%d, end_x=%d, end_y=%d, duration_ms=%d)\n", a.StartX, a.StartY, a.EndX, a.EndY, a.DurationMS))
		case ActionOpenApp:
			sb.WriteString(fmt.Sprintf("    await session.launch_app('%s')\n", a.Value))
		case ActionCloseApp:
			sb.WriteString(fmt.Sprintf("    await session.terminate_app('%s')\n", a.Value))
		}
	}

	sb.WriteString("\n    await session.close()\n\n")
	sb.WriteString("asyncio.run(main())\n")
	return sb.String(), nil
}

// GoGenerator compiles actions to Go SDK code.
type GoGenerator struct{}

func (g *GoGenerator) Generate(actions []ActionIR) (string, error) {
	var sb strings.Builder
	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"context\"\n")
	sb.WriteString("\t\"github.com/ranjith035/morego/sdk\"\n")
	sb.WriteString(")\n\n")
	sb.WriteString("func main() {\n")
	sb.WriteString("\tctx := context.Background()\n")
	sb.WriteString("\tdevice, _ := sdk.Connect(ctx, \"device_id\")\n")
	sb.WriteString("\tsession, _ := device.NewSession(ctx, \"app_id\")\n")
	sb.WriteString("\tdefer session.Close(ctx)\n\n")

	for _, a := range actions {
		switch a.Type {
		case ActionClick:
			sb.WriteString(fmt.Sprintf("\t_ = session.Locator(\"%s\", \"%s\").Click(ctx)\n", a.SelectorStrategy, a.Selector))
		case ActionFill:
			sb.WriteString(fmt.Sprintf("\t_ = session.Locator(\"%s\", \"%s\").Fill(ctx, \"%s\")\n", a.SelectorStrategy, a.Selector, a.Value))
		case ActionSwipe:
			sb.WriteString(fmt.Sprintf("\t_ = session.Swipe(ctx, %d, %d, %d, %d, %d)\n", a.StartX, a.StartY, a.EndX, a.EndY, a.DurationMS))
		case ActionOpenApp:
			sb.WriteString(fmt.Sprintf("\t_ = session.LaunchApp(ctx, \"%s\")\n", a.Value))
		case ActionCloseApp:
			sb.WriteString(fmt.Sprintf("\t_ = session.TerminateApp(ctx, \"%s\")\n", a.Value))
		}
	}
	sb.WriteString("}\n")
	return sb.String(), nil
}

// JavaGenerator compiles actions to Java SDK code.
type JavaGenerator struct{}

func (g *JavaGenerator) Generate(actions []ActionIR) (string, error) {
	var sb strings.Builder
	sb.WriteString("import org.automation.Mobile;\n")
	sb.WriteString("import org.automation.Session;\n\n")
	sb.WriteString("public class RecordedTest {\n")
	sb.WriteString("    public static void main(String[] args) {\n")
	sb.WriteString("        var device = Mobile.connect(\"device_id\");\n")
	sb.WriteString("        var session = device.newSession(\"app_id\");\n\n")

	for _, a := range actions {
		switch a.Type {
		case ActionClick:
			sb.WriteString(fmt.Sprintf("        session.locator(\"%s\", \"%s\").click();\n", a.SelectorStrategy, a.Selector))
		case ActionFill:
			sb.WriteString(fmt.Sprintf("        session.locator(\"%s\", \"%s\").fill(\"%s\");\n", a.SelectorStrategy, a.Selector, a.Value))
		case ActionSwipe:
			sb.WriteString(fmt.Sprintf("        session.swipe(%d, %d, %d, %d, %d);\n", a.StartX, a.StartY, a.EndX, a.EndY, a.DurationMS))
		case ActionOpenApp:
			sb.WriteString(fmt.Sprintf("        session.launchApp(\"%s\");\n", a.Value))
		case ActionCloseApp:
			sb.WriteString(fmt.Sprintf("        session.terminateApp(\"%s\");\n", a.Value))
		}
	}

	sb.WriteString("\n        session.close();\n")
	sb.WriteString("    }\n")
	sb.WriteString("}\n")
	return sb.String(), nil
}

// CSharpGenerator compiles actions to C# SDK code.
type CSharpGenerator struct{}

func (g *CSharpGenerator) Generate(actions []ActionIR) (string, error) {
	var sb strings.Builder
	sb.WriteString("using System.Threading.Tasks;\n")
	sb.WriteString("using Automation;\n\n")
	sb.WriteString("class Program {\n")
	sb.WriteString("    static async Task Main() {\n")
	sb.WriteString("        var device = await Mobile.ConnectAsync(\"device_id\");\n")
	sb.WriteString("        var session = await device.NewSessionAsync(\"app_id\");\n\n")

	for _, a := range actions {
		switch a.Type {
		case ActionClick:
			sb.WriteString(fmt.Sprintf("        await session.Locator(\"%s\", \"%s\").ClickAsync();\n", a.SelectorStrategy, a.Selector))
		case ActionFill:
			sb.WriteString(fmt.Sprintf("        await session.Locator(\"%s\", \"%s\").FillAsync(\"%s\");\n", a.SelectorStrategy, a.Selector, a.Value))
		case ActionSwipe:
			sb.WriteString(fmt.Sprintf("        await session.SwipeAsync(%d, %d, %d, %d, %d);\n", a.StartX, a.StartY, a.EndX, a.EndY, a.DurationMS))
		case ActionOpenApp:
			sb.WriteString(fmt.Sprintf("        await session.LaunchAppAsync(\"%s\");\n", a.Value))
		case ActionCloseApp:
			sb.WriteString(fmt.Sprintf("        await session.TerminateAppAsync(\"%s\");\n", a.Value))
		}
	}

	sb.WriteString("\n        await session.CloseAsync();\n")
	sb.WriteString("    }\n")
	sb.WriteString("}\n")
	return sb.String(), nil
}

// KotlinGenerator compiles actions to Kotlin SDK code.
type KotlinGenerator struct{}

func (g *KotlinGenerator) Generate(actions []ActionIR) (string, error) {
	var sb strings.Builder
	sb.WriteString("import org.automation.mobile\n\n")
	sb.WriteString("suspend fun main() {\n")
	sb.WriteString("    val device = mobile.connect(\"device_id\")\n")
	sb.WriteString("    val session = device.newSession(\"app_id\")\n\n")

	for _, a := range actions {
		switch a.Type {
		case ActionClick:
			sb.WriteString(fmt.Sprintf("    session.locator(\"%s\", \"%s\").click()\n", a.SelectorStrategy, a.Selector))
		case ActionFill:
			sb.WriteString(fmt.Sprintf("    session.locator(\"%s\", \"%s\").fill(\"%s\")\n", a.SelectorStrategy, a.Selector, a.Value))
		case ActionSwipe:
			sb.WriteString(fmt.Sprintf("    session.swipe(%d, %d, %d, %d, %d)\n", a.StartX, a.StartY, a.EndX, a.EndY, a.DurationMS))
		case ActionOpenApp:
			sb.WriteString(fmt.Sprintf("    session.launchApp(\"%s\")\n", a.Value))
		case ActionCloseApp:
			sb.WriteString(fmt.Sprintf("    session.terminateApp(\"%s\")\n", a.Value))
		}
	}

	sb.WriteString("\n    session.close()\n")
	sb.WriteString("}\n")
	return sb.String(), nil
}
