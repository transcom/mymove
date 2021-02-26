const getRoleTypesFromRoles = (roles) => {
  let roleTypes = [];
  if (roles) {
    roleTypes = roles.map((role) => role.roleType);
  }

  return roleTypes;
};

export default getRoleTypesFromRoles;
