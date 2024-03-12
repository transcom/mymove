import React, { Component } from 'react';
import { connect } from 'react-redux';
import { func } from 'prop-types';
import qs from 'query-string';

import MtoShipmentForm from 'components/Customer/MtoShipmentForm/MtoShipmentForm';
import DateAndLocation from 'pages/MyMove/PPM/Booking/DateAndLocation/DateAndLocation';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import {
  updateMTOShipment as updateMTOShipmentAction,
  updateAllMoves as updateAllMovesAction,
} from 'store/entities/actions';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentShipmentFromMove,
  selectAllMoves,
  selectCurrentMoveFromAllMoves,
} from 'store/entities/selectors';
import { fetchCustomerData as fetchCustomerDataAction } from 'store/onboarding/actions';
import { AddressShape } from 'types/address';
import { ServiceMemberShape } from 'types/customerShapes';
import { RouterShape } from 'types/index';
import withRouter from 'utils/routing';
import { getAllMoves } from 'services/internalApi';

export class CreateOrEditMtoShipment extends Component {
  componentDidMount() {
    const { fetchCustomerData, serviceMember, updateAllMoves } = this.props;
    fetchCustomerData();
    getAllMoves(serviceMember.id).then((response) => {
      updateAllMoves(response);
    });
  }

  render() {
    const {
      router: { location },
      currentResidence,
      updateMTOShipment,
      serviceMember,
      serviceMemberMoves,
      moveId,
      mtoShipmentId,
    } = this.props;

    const { type } = qs.parse(location.search);

    const move = selectCurrentMoveFromAllMoves(serviceMemberMoves, moveId);
    const mtoShipment = selectCurrentShipmentFromMove(move, mtoShipmentId);
    const { orders } = move ?? {};

    // loading placeholder while data loads - this handles any async issues
    // loading placeholder while data loads - this handles any async issues
    if (!serviceMemberMoves || !serviceMemberMoves.currentMove || !serviceMemberMoves.previousMoves) {
      return <LoadingPlaceholder />;
    }

    // wait until MTO shipment has loaded to render form
    if (type || mtoShipment?.id) {
      if (type === SHIPMENT_OPTIONS.PPM || mtoShipment?.shipmentType === SHIPMENT_OPTIONS.PPM) {
        return (
          <DateAndLocation
            move={move}
            mtoShipment={mtoShipment}
            serviceMember={serviceMember}
            destinationDutyLocation={orders.new_duty_location}
          />
        );
      }

      return (
        <MtoShipmentForm
          mtoShipment={mtoShipment}
          shipmentType={type || mtoShipment.shipmentType}
          isCreatePage={!!type}
          currentResidence={currentResidence}
          newDutyLocationAddress={orders.new_duty_location?.address}
          updateMTOShipment={updateMTOShipment}
          serviceMember={serviceMember}
          orders={orders}
        />
      );
    }

    return <LoadingPlaceholder />;
  }
}

CreateOrEditMtoShipment.propTypes = {
  router: RouterShape,
  fetchCustomerData: func.isRequired,
  currentResidence: AddressShape.isRequired,
  updateMTOShipment: func.isRequired,
  serviceMember: ServiceMemberShape,
};

CreateOrEditMtoShipment.defaultProps = {
  router: {},
  serviceMember: {},
};

function mapStateToProps(state, ownProps) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberMoves = selectAllMoves(state);
  const {
    router: {
      params: { mtoShipmentId, moveId },
    },
  } = ownProps;
  const props = {
    serviceMember,
    serviceMemberMoves,
    moveId,
    mtoShipmentId,
    currentResidence: serviceMember?.residential_address || {},
  };

  return props;
}

const mapDispatchToProps = {
  fetchCustomerData: fetchCustomerDataAction,
  updateMTOShipment: updateMTOShipmentAction,
  updateAllMoves: updateAllMovesAction,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(CreateOrEditMtoShipment));
