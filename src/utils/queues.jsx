import React from 'react';

import { deleteAssignedOfficeUserForMove, updateAssignedOfficeUserForMove } from 'services/ghcApi';
import { DEFAULT_EMPTY_VALUE } from 'shared/constants';

export const formatOfficeUser = (user) => {
  const fullName = `${user?.lastName}, ${user?.firstName}`;
  return { label: fullName, value: user.officeUserId };
};

const addAssignedOfficeUser = (users, assignedTo) => {
  const newAvailableOfficeUsers = users.slice();
  newAvailableOfficeUsers.push(formatOfficeUser(assignedTo));
  return newAvailableOfficeUsers;
};

export const formatAvailableOfficeUsers = (users, isSupervisor, currentUserId) => {
  if (!users?.length || isSupervisor === undefined || currentUserId === undefined) return [];

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

export const formatAvailableOfficeUsersForRow = (originalRow, supervisor, currentUserId) => {
  // dupe the row to avoid issues with passing office user array by reference
  const row = { ...originalRow };

  row.availableOfficeUsers = formatAvailableOfficeUsers(row.availableOfficeUsers, supervisor, currentUserId);
  // if the move is assigned to a user not present in availableOfficeUsers
  // lets push them onto the end
  if (
    row.assignedTo !== undefined &&
    !row.availableOfficeUsers?.some((user) => user.value === row.assignedTo.officeUserId)
  ) {
    row.availableOfficeUsers = addAssignedOfficeUser(row.availableOfficeUsers, row.assignedTo);
  }
  const { assignedTo, availableOfficeUsers } = row;

  // if there is an assigned user, assign to a variable so we can set a default value below
  const assignedToUser = availableOfficeUsers.find((user) => user.value === assignedTo?.officeUserId);

  const formattedAvailableOfficeUsers = availableOfficeUsers.map(({ value, label }) => (
    <option value={value} key={`filterOption_${value}`}>
      {label}
    </option>
  ));

  return { formattedAvailableOfficeUsers, assignedToUser };
};

export const handleQueueAssignment = (moveID, officeUserId, roleType) => {
  if (officeUserId === DEFAULT_EMPTY_VALUE) deleteAssignedOfficeUserForMove({ moveID, roleType });
  else updateAssignedOfficeUserForMove({ moveID, officeUserId, roleType });
};
