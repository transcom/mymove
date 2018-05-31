import { get, union } from 'lodash';
import { showDutyStationTransportationOffice } from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

const getDutyStationTransportationOfficeType =
  'GET_DUTY_STATION_TRANSPORTATION_OFFICE';
export const GET_DUTY_STATION_TRANSPORTATION_OFFICE = ReduxHelpers.generateAsyncActionTypes(
  getDutyStationTransportationOfficeType,
);

export const loadDutyStationTransportationOffice = dutyStationId => dispatch => {
  const actions = ReduxHelpers.generateAsyncActions(
    getDutyStationTransportationOfficeType,
  );
  dispatch(actions.start);
  return showDutyStationTransportationOffice(dutyStationId)
    .then(transportationOffice =>
      dispatch(actions.success({ transportationOffice, dutyStationId })),
    )
    .catch(error => dispatch(actions.error(error)));
};

ReduxHelpers.generateAsyncActionCreator(
  getDutyStationTransportationOfficeType,
  showDutyStationTransportationOffice,
);

export const getDutyStationTransportationOffice = (state, dutyStationId) => {
  const officeId = get(
    state,
    `transportationOffices.byDutyStationId[${dutyStationId}]`,
  );
  return get(state, `transportationOffices.byId.${officeId}`);
};

function addTransportationOffice(state, action) {
  const { transportationOffice, dutyStationId } = action.payload;
  const { id } = transportationOffice;
  return {
    isLoading: false,
    hasErrored: false,
    hasLoaded: true,
    byId: {
      ...state.byId,
      [id]: transportationOffice,
    },
    allIds: union(state.allIds, [id]),
    byDutyStationId: {
      ...state.byDutyStationId,
      [dutyStationId]: id,
    },
  };
}

const initialState = {
  isLoading: false,
  hasErrored: false,
  hasLoaded: false,
  byId: {},
  allIds: [],
  byDutyStationId: {},
};

export const reducer = (state = initialState, action) => {
  switch (action.type) {
    case GET_DUTY_STATION_TRANSPORTATION_OFFICE.start:
      return Object.assign({}, state, {
        isLoading: true,
        hasErrored: false,
        hasLoaded: false,
      });
    case GET_DUTY_STATION_TRANSPORTATION_OFFICE.success: {
      return addTransportationOffice(state, action);
    }
    case GET_DUTY_STATION_TRANSPORTATION_OFFICE.failure: {
      const result = Object.assign({}, state, {
        isLoading: false,
        hasErrored: true,
        hasLoaded: false,
        error: action.error,
      });
      return result;
    }
    default: {
      return state;
    }
  }
};

export default reducer;
