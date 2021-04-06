import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { push } from 'connected-react-router';

import { OrdersShape, ServiceMemberShape } from 'types/customerShapes';
import { selectServiceMemberFromLoggedInUser, selectCurrentOrders, selectCurrentMove } from 'store/entities/selectors';
import SectionWrapper from 'components/Customer/SectionWrapper';
import ServiceInfoTable from 'components/Customer/Review/ServiceInfoTable';

const Profile = ({ currentOrders, serviceMember }) => {
  const rank = currentOrders ? currentOrders.grade : serviceMember.rank;
  const currentStation = currentOrders ? currentOrders.origin_duty_station : serviceMember.current_station;
  const stationPhoneLines = currentStation?.transportation_office?.phone_lines;
  const stationPhone = stationPhoneLines ? stationPhoneLines[0] : '';

  const handleEditClick = (path) => {
    push(path);
  };

  return (
    <>
      <h1>Profile</h1>
      <SectionWrapper>
        <ServiceInfoTable
          firstName={serviceMember.first_name}
          lastName={serviceMember.last_name}
          currentDutyStationName={currentStation.name}
          currentDutyStationPhone={stationPhone}
          affiliation={serviceMember.affiliation}
          rank={rank}
          edipi={serviceMember.edipi}
          onEditClick={handleEditClick}
        />
      </SectionWrapper>
    </>
  );
};

Profile.propTypes = {
  serviceMember: ServiceMemberShape.isRequired,
  currentOrders: OrdersShape.isRequired,
};

function mapStateToProps(state) {
  return {
    serviceMember: selectServiceMemberFromLoggedInUser(state),
    move: selectCurrentMove(state) || {},
    currentOrders: selectCurrentOrders(state),
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
