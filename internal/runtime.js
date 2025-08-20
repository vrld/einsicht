function loadURL(_, url) {
  if (window.confirm(url)) {
    window.open(url, "_self")
  }
}

function loadImage(element, url) {
  if (window.confirm(url)) {
    console.log('loadImage', { element, url })
  }
}
