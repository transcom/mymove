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
  const { dirty, pageList, pageKey, handleSubmit } = this.props;
  const path = func(pageList, pageKey);
  if (dirty && handleSubmit && shouldHandleSubmit) {
    const awaitSubmit = await handleSubmit();
    if (!awaitSubmit || !awaitSubmit.error) {
      this.goto(path);
    }
  } else {
    this.goto(path);
  }
}
