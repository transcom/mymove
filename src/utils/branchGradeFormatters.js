import { ORDERS_BRANCH_OPTIONS, ORDERS_PAY_GRADE_OPTIONS } from 'constants/orders';

export default function friendlyBranchGrade(branch, grade) {
  const friendlyBranch = ORDERS_BRANCH_OPTIONS[branch];
  const friendlyGrade = ORDERS_PAY_GRADE_OPTIONS[grade];
  if (friendlyBranch && friendlyGrade) {
    return `${friendlyBranch}, ${friendlyGrade}`;
  }
  return '';
}
