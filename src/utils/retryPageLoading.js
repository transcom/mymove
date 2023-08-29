const MilmoveHasBeenForceRefreshed = 'milmove-has-been-force-refreshed';

export const retryPageLoading = (error) => {
  // if we see a chunk load error, try to reload the window to get
  // the latest version of the code
  if (!!error && error.name === 'ChunkLoadError' && !!window) {
    const pageHasAlreadyBeenForceRefreshed = window.localStorage.getItem(MilmoveHasBeenForceRefreshed) === 'true';

    if (!pageHasAlreadyBeenForceRefreshed) {
      window.localStorage.setItem(MilmoveHasBeenForceRefreshed, 'true');
      return window.location.reload();
    }
    window.localStorage.setItem(MilmoveHasBeenForceRefreshed, 'false');
  }
  return false;
};

export default retryPageLoading;
