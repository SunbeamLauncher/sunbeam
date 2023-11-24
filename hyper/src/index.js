// @ts-check

// @ts-ignore
const { globalShortcut, dialog, app, Tray, Menu, shell } = require("electron");
const { toggleWindows, hideWindows, showWindows } = require("./window");
const { getCenterOnCurrentScreen } = require("./screen");
const path = require("path");
const os = require("os");

let unload = () => { }

function onApp(app) {
  // Hide the dock icon
  app.dock.hide();

  // Create tray icon
  const tray = new Tray(path.join(__dirname, "../assets/trayiconTemplate.png"));
  const contextMenu = Menu.buildFromTemplate([
    {
      label: 'Show Sunbeam',
      click: () => {
        showWindows(app);
      },
    },
    { type: "separator" },
    {
      label: 'Edit Sunbeam Config',
      click: () => {
        shell.openPath(path.join(os.homedir(), '.config', 'sunbeam', 'sunbeam.json'))
      },
    },
    {
      label: 'Edit Hyper Config',
      click: () => {
        shell.openPath(path.join(os.homedir(), '.hyper.js'))
      },
    },
    { type: 'separator' },
    {
      label: 'Browse Documentation',
      click: () => {
        shell.openExternal('https://sunbeam.deno.dev/docs');
      },
    },
    {
      label: 'Open Github Repository',
      click: () => {
        shell.openExternal('https://github.com/pomdtr/sunbeam');
      }
    },
    { type: 'separator' },
    {
      label: 'Quit',
      click: () => {
        app.quit();
      },
    },
  ]);
  tray.setToolTip('Sunbeam');
  tray.setContextMenu(contextMenu);

  // Hide windows when the app looses focus
  const onBlur = () => {
    hideWindows(app);
  }
  app.on("browser-window-blur", onBlur);

  unload = () => {
    tray.destroy();
    globalShortcut.unregisterAll();
    app.removeListener("browser-window-blur", onBlur);
  };
};

function onWindow(win) {
  win.on("close", () => {
    app.hide()
  });
}


function onUnload() {
  unload();
}

// Hide window controls on macOS
function decorateBrowserOptions(defaults) {
  const bounds = getCenterOnCurrentScreen(defaults.width, defaults.height);
  return Object.assign({}, defaults, {
    ...bounds,
    titleBarStyle: '',
    transparent: true,
    frame: false,
    alwaysOnTop: true,
    type: "panel",
    skipTaskbar: true,
    movable: false,
    fullscreenable: false,
    minimizable: false,
    maximizable: false,
    resizable: false
  });
};


function decorateConfig(config) {
  globalShortcut.unregisterAll();

  if (config.sunbeam && config.sunbeam.hotkey) {
    const hotkey = config.sunbeam.hotkey;
    if (!globalShortcut.register(hotkey, () => toggleWindows(app))) {
      dialog.showMessageBox({
        message: `Could not register hotkey (${hotkey})`,
        buttons: ["Ok"]
      });
    }
  }

  const css = `
    .header_header {
      top: 0;
      right: 0;
      left: 0;
    }
    .tabs_borderShim {
      display: none;
    }
    .tabs_title {
      display: none;
    }
    .tabs_nav {
      height: auto;
    }
    .tabs_list {
      margin-left: 0;
    }
    .tab_tab:first-of-type {
      border-left-width: 0;
      padding-left: 1px;
    }
  `
  return Object.assign({}, config, {
    css: `
      ${config.css || ''}
      ${css}
    `
  });
}

// Removes the redundant space on mac if there is only one tab
function getTabsProps(parentProps, props) {
  var classTermsList = document.getElementsByClassName('terms_terms')
  if (classTermsList.length > 0) {
    var classTerms = classTermsList[0]
    var header = document.getElementsByClassName('header_header')[0]
    if (props.tabs.length <= 1) {
      // @ts-ignore
      classTerms.style.marginTop = 0
      // @ts-ignore
      header.style.visibility = 'hidden'
    } else {
      // @ts-ignore
      classTerms.style.marginTop = ''
      // @ts-ignore
      header.style.visibility = ''
    }
  }
  return Object.assign({}, parentProps, props)
}

module.exports = {
  onApp,
  onWindow,
  onUnload,
  decorateBrowserOptions,
  getTabsProps,
  decorateConfig,
};
