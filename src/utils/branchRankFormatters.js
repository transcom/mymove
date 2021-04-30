import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders';

export default function friendlyBranchRank(branch, rank) {
  const friendlyBranch = ORDERS_BRANCH_OPTIONS[branch];
  const friendlyRank = ORDERS_RANK_OPTIONS[rank];
  if (friendlyBranch && friendlyRank) {
    return `${friendlyBranch}, ${friendlyRank}`;
  }
  return '';
}
