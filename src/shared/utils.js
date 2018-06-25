import { get, includes, findIndex } from 'lodash';

export const no_op = () => undefined;
export const no_op_action = () => {
  return function(dispatch) {
    dispatch({
      type: 'NO_OP_TYPE',
      item: null,
    });
  };
};

// Turn an array into an object with .reduce()
export const objFromArray = array =>
  array.reduce((accumulator, current) => {
    accumulator[current.id] = current;
    return accumulator;
  }, {});

export const upsert = (arr, newValue) => {
  const index = arr.findIndex(obj => obj.id === newValue.id);
  if (index !== -1) {
    arr.splice(index, 1, newValue);
  } else {
    arr.push(newValue);
  }
};

// Get index of active object in array of objects
export function fetchActive(foos) {
  if (!foos) {
    console.log('NO FOOS???');
    return null;
  }
  // Active moves, orders, or ppms cannot have a status of completed or canceled
  const activeIndex = findIndex(foos, function(o) {
    return includes(
      ['DRAFT', 'SUBMITTED', 'APPROVED', 'IN_PROGRESS'],
      get(o, 'status'),
    );
  }); // -1 is returned if no index is found
  if (activeIndex !== -1) {
    return activeIndex;
  } else {
    return null;
  }
}
