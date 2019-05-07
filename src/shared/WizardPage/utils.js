export function getNextPagePath(pageList, currentPage) {
  const index = pageList.indexOf(currentPage);
  if (index === -1) return undefined;
  return pageList[index + 1];
}
export function getPreviousPagePath(pageList, currentPage) {
  const index = pageList.indexOf(currentPage);
  if (index === -1) return undefined;
  return pageList[index - 1];
}
export function isFirstPage(pageList, currentPage) {
  const index = pageList.indexOf(currentPage);
  return index === 0;
}

export function isLastPage(pageList, currentPage) {
  const index = pageList.indexOf(currentPage);
  return index === pageList.length - 1;
}

export async function beforeTransition(func, shouldHandleSubmit = true) {
  // Check whether there is work to do before transitioning to next page
  // If so, do it and make sure it succeeds before moving on
  const handleSubmit = this.props.handleSubmit;
  let gotoNext = true;
  if (this.props.dirty && handleSubmit && shouldHandleSubmit) {
    const awaitSubmit = await handleSubmit(); // may cause pagelist to change
    if (awaitSubmit && awaitSubmit.error) {
      console.error(awaitSubmit.error);
      gotoNext = false;
    }
  }

  // Good to go to the next page
  if (gotoNext) {
    // Fetch the pageList here in case handleSubmit causes the pageList to change
    const path = func(this.props.pageList, this.props.pageKey);
    if (path !== undefined) {
      this.goto(path);
    }
  }
}
