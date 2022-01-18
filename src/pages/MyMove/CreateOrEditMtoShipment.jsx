import React, { Component } from 'react';
import { connect } from 'react-redux';
import { func } from 'prop-types';
import { withRouter } from 'react-router-dom';
import qs from 'query-string';

import MtoShipmentForm from 'components/Customer/MtoShipmentForm/MtoShipmentForm';
import { updateMTOShipment as updateMTOShipmentAction } from 'store/entities/actions';
import { fetchCustomerData as fetchCustomerDataAction } from 'store/onboarding/actions';
import { HhgShipmentShape, HistoryShape, MatchShape, OrdersShape } from 'types/customerShapes';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentOrders,
  selectMTOShipmentById,
} from 'store/entities/selectors';
import { AddressShape, SimpleAddressShape } from 'types/address';
import { LocationShape } from 'types/index';

export class CreateOrEditMtoShipment extends Component {
  componentDidMount() {
    const { fetchCustomerData } = this.props;
    fetchCustomerData();
  }

  render() {
    const {
      location,
      match,
      history,
      mtoShipment,
      currentResidence,
      newDutyStationAddress,
      updateMTOShipment,
      orders,
    } = this.props;

    const { type } = qs.parse(location.search);

    if (type === SHIPMENT_OPTIONS.PPM) {
      const { moveId } = match.params;

      history.replace(`/moves/${moveId}/ppm-start`);
      return <div />;
    }

    // wait until MTO shipment has loaded to render form
    if (type || mtoShipment?.id) {
      return (
        <MtoShipmentForm
          match={match}
          history={history}
          mtoShipment={mtoShipment}
          selectedMoveType={type || mtoShipment.shipmentType}
          isCreatePage={!!type}
          currentResidence={currentResidence}
          newDutyStationAddress={newDutyStationAddress}
          updateMTOShipment={updateMTOShipment}
          orders={orders}
        />
      );
    }

    return <LoadingPlaceholder />;
  }
}

CreateOrEditMtoShipment.propTypes = {
  location: LocationShape.isRequired,
  match: MatchShape,
  history: HistoryShape,
  fetchCustomerData: func.isRequired,
  // technically this should be a [Generic]MtoShipmentShape
  // using hhg because it has all the props
  mtoShipment: HhgShipmentShape,
  currentResidence: AddressShape.isRequired,
  newDutyStationAddress: SimpleAddressShape,
  updateMTOShipment: func.isRequired,
  orders: OrdersShape,
};

CreateOrEditMtoShipment.defaultProps = {
  match: { isExact: false, params: { moveID: '' } },
  history: { goBack: () => {}, push: () => {} },
  orders: {},
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
  newDutyStationAddress: {
    city: '',
    state: '',
    postalCode: '',
  },
};

function mapStateToProps(state, ownProps) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  const props = {
    orders: selectCurrentOrders(state) || {},
    mtoShipment: selectMTOShipmentById(state, ownProps.match.params.mtoShipmentId) || {},
    currentResidence: serviceMember?.residential_address || {},
    newDutyStationAddress: selectCurrentOrders(state)?.new_duty_station?.address || {},
  };

  return props;
}

const mapDispatchToProps = {
  fetchCustomerData: fetchCustomerDataAction,
  updateMTOShipment: updateMTOShipmentAction,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(CreateOrEditMtoShipment));
