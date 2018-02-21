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
