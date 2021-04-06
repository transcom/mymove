import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { push } from 'connected-react-router';

import { OrdersShape, ServiceMemberShape } from 'types/customerShapes';
import {
  selectServiceMemberFromLoggedInUser,
  selectMoveIsInDraft,
  selectCurrentOrders,
  selectCurrentMove,
  selectHasCurrentPPM,
} from 'store/entities/selectors';
import ServiceInfoTable from 'components/Customer/Review/ServiceInfoTable';

const Profile = ({ currentOrders, serviceMember }) => {
  const rank = currentOrders ? currentOrders.grade : serviceMember.rank;
  const currentStation = currentOrders ? currentOrders.origin_duty_station : serviceMember.current_station;
  const stationPhone = get(currentStation, 'transportation_office.phone_lines.0');

  return (
    <ServiceInfoTable
      firstName={serviceMember.first_name}
      lastName={serviceMember.last_name}
      currentDutyStationName={currentStation.name}
      currentDutyStationPhone={stationPhone}
      affiliation={serviceMember.affiliation}
      rank={rank}
      edipi={serviceMember.edipi}
    />
  );
};

Profile.propTypes = {
  serviceMember: ServiceMemberShape.isRequired,
  currentOrders: OrdersShape.isRequired,
};

function mapStateToProps(state) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  return {
    serviceMember,
    move: selectCurrentMove(state) || {},
    schema: get(state, 'swaggerInternal.spec.definitions.CreateServiceMemberPayload', {}),
    currentOrders: selectCurrentOrders(state),
    // The move still counts as in draft if there are no orders.
    moveIsInDraft: selectMoveIsInDraft(state) || !selectCurrentOrders(state),
    isPpm: selectHasCurrentPPM(state),
    schemaRank: get(state, 'swaggerInternal.spec.definitions.ServiceMemberRank', {}),
    schemaAffiliation: get(state, 'swaggerInternal.spec.definitions.Affiliation', {}),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      push,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(Profile);
