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
  selectCurrentOrders,
  selectCurrentShipmentFromMove,
  selectAllMoves,
} from 'store/entities/selectors';
import { fetchCustomerData as fetchCustomerDataAction } from 'store/onboarding/actions';
import { AddressShape, SimpleAddressShape } from 'types/address';
import { MoveShape, OrdersShape, ServiceMemberShape } from 'types/customerShapes';
import { RouterShape } from 'types/index';
import { ShipmentShape } from 'types/shipment';
import { selectMove } from 'shared/Entities/modules/moves';
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
      mtoShipment,
      currentResidence,
      newDutyLocationAddress,
      updateMTOShipment,
      serviceMember,
      serviceMemberMoves,
      orders,
      move,
    } = this.props;

    const { type } = qs.parse(location.search);

    // loading placeholder while data loads - this handles any async issues
    if (!serviceMemberMoves) {
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
          newDutyLocationAddress={newDutyLocationAddress}
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
  mtoShipment: ShipmentShape,
  currentResidence: AddressShape.isRequired,
  newDutyLocationAddress: SimpleAddressShape,
  updateMTOShipment: func.isRequired,
  serviceMember: ServiceMemberShape,
  orders: OrdersShape,
  move: MoveShape,
};

CreateOrEditMtoShipment.defaultProps = {
  router: {},
  mtoShipment: {
    customerRemarks: '',
    requestedPickupDate: '',
    requestedDeliveryDate: '',
    destinationAddress: {
      city: '',
      postalCode: '',
      state: '',
      streetAddress1: '',
    },
  },
  newDutyLocationAddress: {
    city: '',
    state: '',
    postalCode: '',
  },
  serviceMember: {},
  orders: {},
  move: {},
};

function mapStateToProps(state, ownProps) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberMoves = selectAllMoves(state);
  const {
    router: {
      params: { mtoShipmentId, moveId },
    },
  } = ownProps;
  const mtoShipment = selectCurrentShipmentFromMove(state, moveId, mtoShipmentId) || {};

  const props = {
    serviceMember,
    serviceMemberMoves,
    orders: selectCurrentOrders(state) || {},
    mtoShipment,
    currentResidence: serviceMember?.residential_address || {},
    newDutyLocationAddress: selectCurrentOrders(state)?.new_duty_location?.address || {},
    move: selectMove(state, moveId),
  };

  return props;
}

const mapDispatchToProps = {
  fetchCustomerData: fetchCustomerDataAction,
  updateMTOShipment: updateMTOShipmentAction,
  updateAllMoves: updateAllMovesAction,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(CreateOrEditMtoShipment));
