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

export async function beforeTransition(props, goto, func) {
  const { dirty, pageList, pageKey, handleSubmit } = props;
  const path = func(pageList, pageKey);
  if (dirty && handleSubmit) {
    const awaitSubmit = await handleSubmit();
    if (!awaitSubmit.error) {
      goto(path);
    }
  } else {
    debugger;
    goto(path);
  }
}
