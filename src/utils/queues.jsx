import React from 'react';

import { DEFAULT_EMPTY_VALUE } from 'shared/constants';

const addAssignedOfficeUser = (users, assignedTo) => {
  const newAvailableOfficeUsers = users.slice();
  const { lastName, firstName, id } = assignedTo;
  newAvailableOfficeUsers.push({
    label: `${lastName}, ${firstName}`,
    value: id,
  });
  return newAvailableOfficeUsers;
};

export const formatOfficeUser = (user) => {
  const fullName = `${user?.lastName}, ${user?.firstName}`;
  return { label: fullName, value: user.officeUserId };
};

export const formatAvailableOfficeUsers = (users, isSupervisor, currentUserId) => {
  if (!users.length || isSupervisor === undefined || currentUserId === undefined) return [];

  // instantiate array with empty value for unassign purposes down the road
  const newAvailableOfficeUsers = [{ label: DEFAULT_EMPTY_VALUE, value: null }];

  // if they are a supervisor, push the whole list
  if (isSupervisor) {
    users.forEach((user) => {
      const newUser = formatOfficeUser(user);
      newAvailableOfficeUsers.push(newUser);
    });
  }

  // if they're not a supervisor, just populate with currentUserId
  if (!isSupervisor) {
    const currentUser = users?.filter((user) => user.officeUserId === currentUserId);
    newAvailableOfficeUsers.push(formatOfficeUser(currentUser[0]));
  }

  return newAvailableOfficeUsers;
};

export const formatAvailableOfficeUsersForRow = (row) => {
  // dupe the row to avoid issues with passing office user array by reference
  const updatedRow = { ...row };

  // if the move is assigned to a user not present in availableOfficeUsers
  // lets push them onto the end
  if (row.assignedTo !== undefined && !row.availableOfficeUsers?.some((user) => user.value === row.assignedTo.id)) {
    updatedRow.availableOfficeUsers = addAssignedOfficeUser(row.availableOfficeUsers, row.assignedTo);
  }
  const { assignedTo, availableOfficeUsers } = updatedRow;

  // if there is an assigned user, assign to a variable so we can set a default value below
  const assignedToUser = availableOfficeUsers.find((user) => user.value === assignedTo?.id);

  const formattedAvailableOfficeUsers = availableOfficeUsers.map(({ value, label }) => (
    <option value={value} key={`filterOption_${value}`}>
      {label}
    </option>
  ));

  return { formattedAvailableOfficeUsers, assignedToUser };
};
