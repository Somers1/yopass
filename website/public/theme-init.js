(function () {
  try {
    var LIGHT = 'emerald';
    var DARK = 'dim';
    var storedMode = localStorage.getItem('themeMode'); // 'light' | 'dark'
    var mode =
      storedMode === 'light' || storedMode === 'dark'
        ? storedMode
        : 'light';
    var daisyTheme = mode === 'dark' ? DARK : LIGHT;
    document.documentElement.setAttribute('data-theme', daisyTheme);
  } catch (e) {}
})();
