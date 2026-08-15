// State variables
let activeNode = 'submit';
let activePlatform = 'android';
let selectedLanguage = 'typescript';
let isRecording = true;

// Simulated element database
const elementDatabase = {
  title: {
    tagName: 'TextView',
    class: 'android.widget.TextView',
    text: 'Welcome Back',
    resourceId: 'com.example:id/title_label',
    contentDesc: '',
    accessibilityId: '',
    bestLocator: 'GetByText("Welcome Back")',
    aiSuggestion: 'GetByText("Welcome Back")',
    aiReason: 'Matched node text. Since this is static text, GetByText is the most robust strategy.',
    confidence: '92%',
    codeTemplate: {
      typescript: `await session.locator('TEXT', 'Welcome Back').isVisible();`,
      python: `await session.locator('TEXT', 'Welcome Back').is_visible()`,
      go: `_ = session.Locator("TEXT", "Welcome Back").AssertVisible(ctx)`,
      java: `session.locator("TEXT", "Welcome Back").assertVisible();`,
      csharp: `await session.Locator("TEXT", "Welcome Back").AssertVisibleAsync();`,
      kotlin: `session.locator("TEXT", "Welcome Back").assertVisible()`
    }
  },
  username: {
    tagName: 'EditText',
    class: 'android.widget.EditText',
    text: 'admin',
    resourceId: 'com.example:id/username_field',
    contentDesc: 'Username input',
    accessibilityId: 'Username input',
    bestLocator: 'GetByAccessibilityId("Username input")',
    aiSuggestion: 'GetByAccessibilityId("Username input")',
    aiReason: 'Highly stable accessibility identifier found. Accessibility IDs remain unchanged during layout moves.',
    confidence: '99%',
    codeTemplate: {
      typescript: `await session.locator('ACCESSIBILITY_ID', 'Username input').fill('admin');`,
      python: `await session.locator('ACCESSIBILITY_ID', 'Username input').fill('admin')`,
      go: `_ = session.Locator("ACCESSIBILITY_ID", "Username input").Fill(ctx, "admin")`,
      java: `session.locator("ACCESSIBILITY_ID", "Username input").fill("admin");`,
      csharp: `await session.Locator("ACCESSIBILITY_ID", "Username input").FillAsync("admin");`,
      kotlin: `session.locator("ACCESSIBILITY_ID", "Username input").fill("admin")`
    }
  },
  password: {
    tagName: 'EditText',
    class: 'android.widget.EditText',
    text: '••••••••',
    resourceId: 'com.example:id/password_field',
    contentDesc: 'Password input',
    accessibilityId: 'Password input',
    bestLocator: 'GetByAccessibilityId("Password input")',
    aiSuggestion: 'GetByAccessibilityId("Password input")',
    aiReason: 'Highly stable accessibility identifier found. Strongly preferred over index relative tags.',
    confidence: '99%',
    codeTemplate: {
      typescript: `await session.locator('ACCESSIBILITY_ID', 'Password input').fill('my_password');`,
      python: `await session.locator('ACCESSIBILITY_ID', 'Password input').fill('my_password')`,
      go: `_ = session.Locator("ACCESSIBILITY_ID", "Password input").Fill(ctx, "my_password")`,
      java: `session.locator("ACCESSIBILITY_ID", "Password input").fill("my_password");`,
      csharp: `await session.Locator("ACCESSIBILITY_ID", "Password input").FillAsync("my_password");`,
      kotlin: `session.locator("ACCESSIBILITY_ID", "Password input").fill("my_password")`
    }
  },
  submit: {
    tagName: 'Button',
    class: 'android.widget.Button',
    text: 'Login',
    resourceId: 'com.example:id/submit_btn',
    contentDesc: 'Login button',
    accessibilityId: 'Login button',
    bestLocator: 'GetByAccessibilityId("Login button")',
    aiSuggestion: 'GetByAccessibilityId("Login button")',
    aiReason: 'Found content-desc "Login button". Conforms to accessibility-first standards and is immune to UI location shifts.',
    confidence: '98%',
    codeTemplate: {
      typescript: `await session.locator('ACCESSIBILITY_ID', 'Login button').click();`,
      python: `await session.locator('ACCESSIBILITY_ID', 'Login button').click()`,
      go: `_ = session.Locator("ACCESSIBILITY_ID", "Login button").Click(ctx)`,
      java: `session.locator("ACCESSIBILITY_ID", "Login button").click();`,
      csharp: `await session.Locator("ACCESSIBILITY_ID", "Login button").ClickAsync();`,
      kotlin: `session.locator("ACCESSIBILITY_ID", "Login button").click()`
    }
  }
};

// Full recorded script templates based on language selection
const fullScripts = {
  typescript: `import { mobile } from '@mobile/playwright';

(async () => {
  const device = await mobile.connect({ deviceId: 'pixel_6_pro' });
  const session = await device.newSession({ appId: 'com.example.loginapp' });

  await session.locator('ACCESSIBILITY_ID', 'Username input').fill('admin');
  await session.locator('ACCESSIBILITY_ID', 'Password input').fill('password');
  await session.locator('ACCESSIBILITY_ID', 'Login button').click();

  await session.close();
})();`,
  python: `import asyncio
from mobile_playwright import mobile

async def main():
    device = await mobile.connect(device_id='pixel_6_pro')
    session = await device.new_session(app_id='com.example.loginapp')

    await session.locator('ACCESSIBILITY_ID', 'Username input').fill('admin')
    await session.locator('ACCESSIBILITY_ID', 'Password input').fill('password')
    await session.locator('ACCESSIBILITY_ID', 'Login button').click()

    await session.close()

asyncio.run(main())`,
  go: `package main

import (
	"context"
	"github.com/ranjith035/morego/sdk"
)

func main() {
	ctx := context.Background()
	device, _ := sdk.Connect(ctx, "pixel_6_pro")
	session, _ := device.NewSession(ctx, "com.example.loginapp")
	defer session.Close(ctx)

	_ = session.Locator("ACCESSIBILITY_ID", "Username input").Fill(ctx, "admin")
	_ = session.Locator("ACCESSIBILITY_ID", "Password input").Fill(ctx, "password")
	_ = session.Locator("ACCESSIBILITY_ID", "Login button").Click(ctx)
}`,
  java: `import org.automation.Mobile;
import org.automation.Session;

public class RecordedTest {
    public static void main(String[] args) {
        var device = Mobile.connect("pixel_6_pro");
        var session = device.newSession("com.example.loginapp");

        session.locator("ACCESSIBILITY_ID", "Username input").fill("admin");
        session.locator("ACCESSIBILITY_ID", "Password input").fill("password");
        session.locator("ACCESSIBILITY_ID", "Login button").click();

        session.close();
    }
}`,
  csharp: `using System.Threading.Tasks;
using Automation;

class Program {
    static async Task Main() {
        var device = await Mobile.ConnectAsync("pixel_6_pro");
        var session = await device.NewSessionAsync("com.example.loginapp");

        await session.Locator("ACCESSIBILITY_ID", "Username input").FillAsync("admin");
        await session.Locator("ACCESSIBILITY_ID", "Password input").FillAsync("password");
        await session.Locator("ACCESSIBILITY_ID", "Login button").ClickAsync();

        await session.CloseAsync();
    }
}`,
  kotlin: `import org.automation.mobile

suspend fun main() {
    val device = mobile.connect("pixel_6_pro")
    val session = device.newSession("com.example.loginapp")

    session.locator("ACCESSIBILITY_ID", "Username input").fill("admin")
    session.locator("ACCESSIBILITY_ID", "Password input").fill("password")
    session.locator("ACCESSIBILITY_ID", "Login button").click()

    session.close()
}`
};

// Elements DOM Cache
const domElements = {
  boxes: {
    title: document.getElementById('box-title'),
    username: document.getElementById('box-username'),
    password: document.getElementById('box-password'),
    submit: document.getElementById('box-submit')
  },
  treeNodes: {
    title: document.getElementById('tree-title'),
    username: document.getElementById('tree-username'),
    password: document.getElementById('tree-password'),
    submit: document.getElementById('tree-submit')
  },
  locatorInput: document.getElementById('locator-input'),
  aiLocatorText: document.getElementById('ai-locator-text'),
  aiReasonText: document.querySelector('.ai-card-reason'),
  aiConfidenceBadge: document.querySelector('.ai-confidence-badge'),
  generatedCodeBox: document.getElementById('generated-code-box'),
  languageSelect: document.getElementById('language-select'),
  recordBtn: document.getElementById('record-btn'),
  recordText: document.getElementById('record-text'),
  pulseDot: document.querySelector('.pulse-dot'),
  logContent: document.getElementById('log-content'),
  platformAndroid: document.getElementById('opt-android'),
  platformIos: document.getElementById('opt-ios'),
  deviceInfoText: document.getElementById('device-info-text')
};

// Initial state updates
updateUI();

// Platform Toggles
domElements.platformAndroid.addEventListener('click', () => {
  activePlatform = 'android';
  domElements.platformAndroid.classList.add('active');
  domElements.platformIos.classList.remove('active');
  domElements.deviceInfoText.textContent = 'Pixel 6 Pro (ADB Connection: 127.0.0.1:5555)';
  addConsoleLog('Switched platform to Android. ADB connection established.', 'info');
});

domElements.platformIos.addEventListener('click', () => {
  activePlatform = 'ios';
  domElements.platformIos.classList.add('active');
  domElements.platformAndroid.classList.remove('active');
  domElements.deviceInfoText.textContent = 'iPhone 14 Pro Simulator (XCUITest Agent: Port 8100)';
  addConsoleLog('Switched platform to iOS. Connected to local XCUITest agent.', 'info');
});

// Sync Select Language
domElements.languageSelect.addEventListener('change', (e) => {
  selectedLanguage = e.target.value;
  updateUI();
  addConsoleLog(`SDK compilation target updated to ${selectedLanguage.toUpperCase()}`, 'info');
});

// Setup Device canvas highlight triggers
Object.keys(domElements.boxes).forEach((key) => {
  const box = domElements.boxes[key];
  
  box.addEventListener('mouseenter', () => {
    highlightElement(key);
  });

  box.addEventListener('mouseleave', () => {
    unhighlightElement(key);
  });

  box.addEventListener('click', () => {
    selectElement(key);
  });
});

// Setup Tree Node highlight triggers
Object.keys(domElements.treeNodes).forEach((key) => {
  const node = domElements.treeNodes[key];

  node.addEventListener('mouseenter', () => {
    highlightElement(key);
  });

  node.addEventListener('mouseleave', () => {
    unhighlightElement(key);
  });

  node.addEventListener('click', () => {
    selectElement(key);
  });
});

// Locator search matching preview filter
domElements.locatorInput.addEventListener('input', (e) => {
  const query = e.target.value.toLowerCase().trim();
  if (query === '') {
    clearHighlightBorders();
    return;
  }

  // Basic matcher
  Object.keys(elementDatabase).forEach((key) => {
    const data = elementDatabase[key];
    const box = domElements.boxes[key];
    
    const matchesQuery = 
      key.includes(query) || 
      data.text.toLowerCase().includes(query) || 
      data.resourceId.toLowerCase().includes(query) ||
      data.bestLocator.toLowerCase().includes(query);

    if (matchesQuery) {
      box.style.borderColor = '#06b6d4'; // cyan matching
      box.style.backgroundColor = 'rgba(6, 182, 212, 0.2)';
    } else {
      box.style.borderColor = 'transparent';
      box.style.backgroundColor = 'transparent';
    }
  });
});

// Setup Recording control toggles
domElements.recordBtn.addEventListener('click', () => {
  isRecording = !isRecording;
  if (isRecording) {
    domElements.recordText.textContent = 'Recording';
    domElements.pulseDot.style.display = 'block';
    domElements.recordBtn.style.background = 'linear-gradient(to right, #ef4444, #dc2626)';
    addConsoleLog('Recorder resumed. Listening for gestures...', 'success');
  } else {
    domElements.recordText.textContent = 'Paused';
    domElements.pulseDot.style.display = 'none';
    domElements.recordBtn.style.background = 'linear-gradient(to right, #4b5563, #374151)';
    addConsoleLog('Recorder paused. Execution stubs won\'t generate.', 'warn');
  }
  updateUI();
});

// Action functions
function highlightElement(key) {
  domElements.boxes[key].classList.add('active');
  domElements.treeNodes[key].classList.add('active');
}

function unhighlightElement(key) {
  domElements.boxes[key].classList.remove('active');
  domElements.treeNodes[key].classList.remove('active');
}

function selectElement(key) {
  activeNode = key;
  
  // Clear previous selections
  Object.keys(domElements.treeNodes).forEach((k) => {
    domElements.treeNodes[k].style.backgroundColor = 'transparent';
    domElements.boxes[k].style.borderStyle = 'solid';
  });

  // Apply visual style
  domElements.treeNodes[key].style.backgroundColor = 'var(--bg-input)';
  domElements.boxes[key].style.borderStyle = 'double';

  updateUI();
  
  const data = elementDatabase[key];
  addConsoleLog(`Inspector picked: <${data.tagName} resource-id="${data.resourceId}" />`, 'success');
}

function updateUI() {
  const data = elementDatabase[activeNode];
  if (!data) return;

  // Update Locator info
  domElements.locatorInput.value = data.bestLocator;
  
  // Update AI Box
  domElements.aiLocatorText.textContent = data.aiSuggestion;
  domElements.aiReasonText.textContent = data.aiReason;
  domElements.aiConfidenceBadge.textContent = `${data.confidence} Confidence`;

  // Update Generated Code Editor block
  if (isRecording) {
    domElements.generatedCodeBox.textContent = data.codeTemplate[selectedLanguage];
  } else {
    // Show full script if recording is paused/complete
    domElements.generatedCodeBox.textContent = fullScripts[selectedLanguage];
  }
}

function clearHighlightBorders() {
  Object.keys(domElements.boxes).forEach((k) => {
    domElements.boxes[k].style.borderColor = 'var(--accent-green)';
    domElements.boxes[k].style.backgroundColor = 'rgba(16, 185, 129, 0.15)';
  });
}

function addConsoleLog(message, type = '') {
  const now = new Date();
  const timeStr = now.toTimeString().split(' ')[0];
  
  const logDiv = document.createElement('div');
  logDiv.className = `log-line ${type}`;
  logDiv.textContent = `[${timeStr}] ${message}`;
  
  domElements.logContent.appendChild(logDiv);
  domElements.logContent.scrollTop = domElements.logContent.scrollHeight;
}
