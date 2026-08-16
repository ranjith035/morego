// Active state variables
let sessionId = "";
let physicalWidth = 1080;
let physicalHeight = 2400;
let selectedLanguage = "typescript";
let activeElement = null;

// DOM Cache
const domElements = {
  treeContent: document.getElementById('hierarchy-tree-content'),
  deviceScreen: document.getElementById('device-screen'),
  nodeTag: document.getElementById('node-tag'),
  nodeClass: document.getElementById('node-class'),
  nodeText: document.getElementById('node-text'),
  nodeResourceId: document.getElementById('node-resource-id'),
  nodeDesc: document.getElementById('node-desc'),
  nodeBounds: document.getElementById('node-bounds'),
  actionTap: document.getElementById('action-tap'),
  locatorInput: document.getElementById('locator-input'),
  aiConfidence: document.getElementById('ai-confidence'),
  aiLocatorText: document.getElementById('ai-locator-text'),
  aiReason: document.getElementById('ai-reason'),
  generatedCodeBox: document.getElementById('generated-code-box'),
  languageSelect: document.getElementById('language-select'),
  refreshBtn: document.getElementById('refresh-btn'),
  logContent: document.getElementById('log-content'),
  deviceInfoText: document.getElementById('device-info-text')
};

// Initialize
addConsoleLog("Live Web Inspector initialized. Ready to sync with Go Core server.", "info");
fetchLiveState();

// Event Listeners
domElements.refreshBtn.addEventListener('click', fetchLiveState);

domElements.languageSelect.addEventListener('change', (e) => {
  selectedLanguage = e.target.value;
  if (activeElement) {
    updatePropertiesAndCode(activeElement);
  }
  addConsoleLog(`SDK compilation target updated to ${selectedLanguage.toUpperCase()}`, "info");
});

domElements.actionTap.addEventListener('click', async () => {
  if (!activeElement || !sessionId) return;
  const bounds = parseBounds(activeElement.bounds);
  if (!bounds) return;

  const clickX = Math.round(bounds.x1 + bounds.width / 2);
  const clickY = Math.round(bounds.y1 + bounds.height / 2);

  addConsoleLog(`Sending Tap action to device at (${clickX}, ${clickY})...`, "info");
  
  try {
    const resp = await fetch("/api/session/action/click", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        sessionId: sessionId,
        x: clickX,
        y: clickY
      })
    });
    
    if (resp.ok) {
      addConsoleLog(`Tap action injected successfully. Auto-refreshing in 1.5s...`, "success");
      setTimeout(fetchLiveState, 1500);
    } else {
      const errData = await resp.json();
      addConsoleLog(`Failed to inject tap: ${errData.error}`, "warn");
    }
  } catch (err) {
    addConsoleLog(`Error sending tap action: ${err.message}`, "warn");
  }
});

// Fetch layout XML and Base64 screenshot from server
async function fetchLiveState() {
  addConsoleLog("Syncing layout hierarchy and screen graphics from target device...", "info");
  domElements.refreshBtn.disabled = true;
  domElements.refreshBtn.style.opacity = "0.6";

  try {
    const response = await fetch("/api/session/state");
    if (!response.ok) {
      const errData = await response.json();
      throw new Error(errData.error || "Failed to fetch session state");
    }

    const data = await response.json();
    sessionId = data.sessionId;
    
    // 1. Update Device Info
    domElements.deviceInfoText.textContent = `Xiaomi a5fb8bad (Active Session: ${sessionId})`;
    
    // 2. Load Base64 screenshot
    domElements.deviceScreen.style.backgroundImage = `url(data:image/png;base64,${data.screenshot})`;
    
    // 3. Parse XML Tree
    const parser = new DOMParser();
    const xmlDoc = parser.parseFromString(data.xml, "text/xml");
    
    // 4. Extract physical dimensions from root hierarchy bounds if available
    const rootNode = xmlDoc.getElementsByTagName("hierarchy")[0];
    if (rootNode) {
      // Find children with bounds to get physical width/height
      const childNodes = xmlDoc.getElementsByTagName("node");
      for (let i = 0; i < childNodes.length; i++) {
        const bStr = childNodes[i].getAttribute("bounds");
        const b = parseBounds(bStr);
        if (b) {
          physicalWidth = Math.max(physicalWidth, b.x2);
          physicalHeight = Math.max(physicalHeight, b.y2);
        }
      }
    }
    
    // Scale screen element container to match device aspect ratio
    const screenHeight = 552; // fixed frame height
    const screenWidth = Math.round(screenHeight * (physicalWidth / physicalHeight));
    domElements.deviceScreen.style.width = `${screenWidth}px`;
    domElements.deviceScreen.style.height = `${screenHeight}px`;

    // Clear highlights overlay and hierarchy panel
    domElements.deviceScreen.innerHTML = "";
    domElements.treeContent.innerHTML = "";
    activeElement = null;
    resetPropertiesCard();

    // 5. Render layout tree and overlays recursively
    renderTree(xmlDoc.documentElement, domElements.treeContent);
    addConsoleLog("UI hierarchy mapped successfully. Interactive overlays ready.", "success");

  } catch (error) {
    addConsoleLog(`Sync error: ${error.message}`, "warn");
    domElements.deviceInfoText.textContent = "Offline (No Active Session)";
    domElements.treeContent.innerHTML = `<div style="padding: 15px; color: var(--accent-red); font-size: 13px;">Error: ${error.message}. Ensure core server is running with active physical device session.</div>`;
  } finally {
    domElements.refreshBtn.disabled = false;
    domElements.refreshBtn.style.opacity = "1";
  }
}

// Parse boundary string "[x1,y1][x2,y2]"
function parseBounds(boundsStr) {
  if (!boundsStr) return null;
  const match = boundsStr.match(/\[(\d+),(\d+)\]\[(\d+),(\d+)\]/);
  if (!match) return null;
  return {
    x1: parseInt(match[1]),
    y1: parseInt(match[2]),
    x2: parseInt(match[3]),
    y2: parseInt(match[4]),
    width: parseInt(match[3]) - parseInt(match[1]),
    height: parseInt(match[4]) - parseInt(match[2])
  };
}

// Recursively parse layout tree and append to DOM elements
function renderTree(xmlNode, parentDOM) {
  const tagName = xmlNode.tagName;
  if (!tagName || tagName === "hierarchy") {
    Array.from(xmlNode.childNodes).forEach(child => renderTree(child, parentDOM));
    return;
  }

  const className = xmlNode.getAttribute("class") || "";
  const text = xmlNode.getAttribute("text") || "";
  const resourceId = xmlNode.getAttribute("resource-id") || "";
  const contentDesc = xmlNode.getAttribute("content-desc") || "";
  const bounds = xmlNode.getAttribute("bounds") || "";

  const shortClassName = className.split(".").pop() || className;

  const treeNode = document.createElement("div");
  treeNode.className = "tree-node";

  const contentDiv = document.createElement("div");
  contentDiv.className = "tree-node-content";

  let labelText = `<span class="tree-tag">${shortClassName}</span>`;
  if (resourceId) {
    labelText += ` <span class="tree-attr">id</span>="${escapeHTML(resourceId.split("/").pop())}"`;
  }
  if (text) {
    labelText += ` <span class="tree-attr">text</span>="${escapeHTML(text.substring(0, 15))}"`;
  }
  if (contentDesc) {
    labelText += ` <span class="tree-attr">desc</span>="${escapeHTML(contentDesc.substring(0, 15))}"`;
  }
  contentDiv.innerHTML = `<span>&lt;</span>${labelText}<span>/&gt;</span>`;
  treeNode.appendChild(contentDiv);

  const elementData = {
    tagName: shortClassName,
    className: className,
    text: text,
    resourceId: resourceId,
    contentDesc: contentDesc,
    bounds: bounds
  };

  // Generate highlights box overlay
  const parsedBounds = parseBounds(bounds);
  let overlayBox = null;
  if (parsedBounds && parsedBounds.width > 0 && parsedBounds.height > 0) {
    overlayBox = createOverlayBox(elementData, parsedBounds, treeNode, contentDiv);
  }

  contentDiv.addEventListener("click", (e) => {
    e.stopPropagation();
    selectNode(elementData, treeNode, contentDiv, overlayBox);
  });

  // Recursively append children nodes
  const childrenContainer = document.createElement("div");
  childrenContainer.className = "tree-children";
  let hasChildren = false;
  
  Array.from(xmlNode.childNodes).forEach(child => {
    if (child.nodeType === 1) { // ELEMENT_NODE
      renderTree(child, childrenContainer);
      hasChildren = true;
    }
  });

  if (hasChildren) {
    treeNode.appendChild(childrenContainer);
  }
  parentDOM.appendChild(treeNode);
}

// Create absolute positioned box overlay matching physical element dimensions
function createOverlayBox(elementData, parsedBounds, treeNode, contentDiv) {
  const scaleX = domElements.deviceScreen.clientWidth / physicalWidth;
  const scaleY = domElements.deviceScreen.clientHeight / physicalHeight;

  const box = document.createElement("div");
  box.className = "element-highlight-box";
  box.style.left = `${parsedBounds.x1 * scaleX}px`;
  box.style.top = `${parsedBounds.y1 * scaleY}px`;
  box.style.width = `${parsedBounds.width * scaleX}px`;
  box.style.height = `${parsedBounds.height * scaleY}px`;

  box.addEventListener("mouseenter", () => {
    box.classList.add("active");
    contentDiv.classList.add("active");
  });

  box.addEventListener("mouseleave", () => {
    box.classList.remove("active");
    contentDiv.classList.remove("active");
  });

  box.addEventListener("click", (e) => {
    e.stopPropagation();
    selectNode(elementData, treeNode, contentDiv, box);
  });

  domElements.deviceScreen.appendChild(box);
  return box;
}

// Update picked node details panel
function selectNode(elementData, treeNode, contentDiv, overlayBox) {
  activeElement = elementData;

  // Clear previous selections
  document.querySelectorAll(".tree-node-content").forEach(el => el.style.backgroundColor = "transparent");
  document.querySelectorAll(".element-highlight-box").forEach(el => el.classList.remove("selected"));

  // Apply visual indicators
  contentDiv.style.backgroundColor = "var(--bg-input)";
  if (overlayBox) {
    overlayBox.classList.add("selected");
  }

  updatePropertiesAndCode(elementData);
  
  // Expand tree node ancestors automatically to display selection path
  let parent = treeNode.parentElement;
  while (parent && parent.className === "tree-children") {
    parent.style.display = "block";
    parent = parent.parentElement.parentElement;
  }

  addConsoleLog(`Selected node: <${elementData.tagName} resource-id="${elementData.resourceId || 'N/A'}" />`, "info");
}

function updatePropertiesAndCode(elementData) {
  // Update properties fields
  domElements.nodeTag.textContent = elementData.tagName || "N/A";
  domElements.nodeClass.textContent = elementData.className || "N/A";
  domElements.nodeText.textContent = elementData.text || "N/A";
  domElements.nodeResourceId.textContent = elementData.resourceId || "N/A";
  domElements.nodeDesc.textContent = elementData.contentDesc || "N/A";
  domElements.nodeBounds.textContent = elementData.bounds || "N/A";

  // Calculate best locator strategy
  let strategy = "XPATH";
  let selector = "";
  let confidence = "50%";
  let reason = "Fallback to absolute layout path. No stable element identifiers found.";

  if (elementData.contentDesc) {
    strategy = "ACCESSIBILITY_ID";
    selector = elementData.contentDesc;
    confidence = "99%";
    reason = "Durable content-desc found. Immunizes selector against structural moves.";
  } else if (elementData.resourceId && (elementData.resourceId.includes("test_id") || elementData.resourceId.includes("testid"))) {
    strategy = "TEST_ID";
    selector = elementData.resourceId.split("/").pop() || elementData.resourceId;
    confidence = "98%";
    reason = "Durable developer test ID found. Ideal for UI integration tests.";
  } else if (elementData.text && elementData.text.length > 0 && elementData.text.length <= 30) {
    strategy = "TEXT";
    selector = elementData.text;
    confidence = "95%";
    reason = "Target matches unique page text. Preferred locator in Playwright paradigm.";
  } else if (elementData.resourceId) {
    strategy = "RESOURCE_ID";
    selector = elementData.resourceId;
    confidence = "90%";
    reason = "Platform layout resource ID found. Subject to vendor modification.";
  } else {
    // Generate simple XPath fallback
    strategy = "XPATH";
    selector = `//${elementData.tagName || "node"}[@class='${elementData.className}']`;
    confidence = "60%";
    reason = "Absolute class lookup fallback. Structural shifts may break this path.";
  }

  // Update Locator text
  domElements.locatorInput.value = formatLocatorDisplay(strategy, selector);

  // Update AI suggest card
  domElements.aiLocatorText.textContent = formatLocatorDisplay(strategy, selector);
  domElements.aiConfidence.textContent = `${confidence} Confidence`;
  domElements.aiReason.textContent = reason;

  // Generate code block
  domElements.generatedCodeBox.textContent = getSDKCodeSnippet(strategy, selector, selectedLanguage);
}

function formatLocatorDisplay(strategy, selector) {
  switch (strategy) {
    case "ACCESSIBILITY_ID":
      return `GetByAccessibilityID("${selector}")`;
    case "TEST_ID":
      return `GetByTestID("${selector}")`;
    case "TEXT":
      return `GetByText("${selector}")`;
    case "RESOURCE_ID":
      return `Locator("RESOURCE_ID", "${selector}")`;
    default:
      return `Locator("XPATH", "${selector}")`;
  }
}

function getSDKCodeSnippet(strategy, selector, language) {
  const isClickable = true; // Default to click action generator
  switch (language) {
    case "typescript":
      if (strategy === "ACCESSIBILITY_ID") return `await session.getByAccessibilityID("${selector}").click();`;
      if (strategy === "TEXT") return `await session.getByText("${selector}").click();`;
      if (strategy === "TEST_ID") return `await session.getByTestID("${selector}").click();`;
      return `await session.locator("${strategy}", "${selector}").click();`;

    case "go":
      if (strategy === "ACCESSIBILITY_ID") return `err = session.GetByAccessibilityID("${selector}").Click(ctx)`;
      if (strategy === "TEXT") return `err = session.GetByText("${selector}").Click(ctx)`;
      if (strategy === "TEST_ID") return `err = session.GetByTestID("${selector}").Click(ctx)`;
      return `err = session.Locator("${strategy}", "${selector}").Click(ctx)`;

    case "python":
      if (strategy === "ACCESSIBILITY_ID") return `await session.get_by_accessibility_id("${selector}").click()`;
      if (strategy === "TEXT") return `await session.get_by_text("${selector}").click()`;
      if (strategy === "TEST_ID") return `await session.get_by_test_id("${selector}").click()`;
      return `await session.locator("${strategy}", "${selector}").click()`;

    case "java":
      if (strategy === "ACCESSIBILITY_ID") return `session.getByAccessibilityID("${selector}").click();`;
      if (strategy === "TEXT") return `session.getByText("${selector}").click();`;
      if (strategy === "TEST_ID") return `session.getByTestID("${selector}").click();`;
      return `session.locator("${strategy}", "${selector}").click();`;

    case "csharp":
      if (strategy === "ACCESSIBILITY_ID") return `await session.GetByAccessibilityID("${selector}").ClickAsync();`;
      if (strategy === "TEXT") return `await session.GetByText("${selector}").ClickAsync();`;
      if (strategy === "TEST_ID") return `await session.GetByTestID("${selector}").ClickAsync();`;
      return `await session.Locator("${strategy}", "${selector}").ClickAsync();`;

    default:
      return `session.locator("${strategy}", "${selector}").click()`;
  }
}

function resetPropertiesCard() {
  domElements.nodeTag.textContent = "N/A";
  domElements.nodeClass.textContent = "N/A";
  domElements.nodeText.textContent = "N/A";
  domElements.nodeResourceId.textContent = "N/A";
  domElements.nodeDesc.textContent = "N/A";
  domElements.nodeBounds.textContent = "N/A";
  domElements.locatorInput.value = "";
  domElements.aiLocatorText.textContent = "N/A";
  domElements.aiReason.textContent = "Select an element to calculate options.";
  domElements.generatedCodeBox.textContent = "Select an element to view code...";
}

function addConsoleLog(message, type = "") {
  const now = new Date();
  const timeStr = now.toTimeString().split(" ")[0];
  const logDiv = document.createElement("div");
  logDiv.className = `log-line ${type}`;
  logDiv.textContent = `[${timeStr}] ${message}`;
  domElements.logContent.appendChild(logDiv);
  domElements.logContent.scrollTop = domElements.logContent.scrollHeight;
}

function escapeHTML(str) {
  if (!str) return "";
  return str.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;");
}
