export const stringifyName = ({ first_name: firstName, last_name: lastName }) =>
  [lastName, firstName].filter(name => !!name).join(', ');
