import React, { Component } from 'react';
import { connect } from 'react-redux';
import { func, number, shape } from 'prop-types';
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
import DateAndLocation from 'pages/MyMove/PPMBooking/DateAndLocation/DateAndLocation';

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
      newDutyLocationAddress,
      updateMTOShipment,
      serviceMember,
      orders,
    } = this.props;

    const { type } = qs.parse(location.search);

    // wait until MTO shipment has loaded to render form
    if (type || mtoShipment?.id) {
      if (type === SHIPMENT_OPTIONS.PPM || mtoShipment?.shipmentType === SHIPMENT_OPTIONS.PPM) {
        return (
          <DateAndLocation
            mtoShipment={mtoShipment}
            serviceMember={serviceMember}
            destinationDutyLocation={orders.new_duty_location}
          />
        );
      }

      return (
        <MtoShipmentForm
          match={match}
          history={history}
          mtoShipment={mtoShipment}
          selectedMoveType={type || mtoShipment.shipmentType}
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
  location: LocationShape.isRequired,
  match: MatchShape,
  history: HistoryShape,
  fetchCustomerData: func.isRequired,
  // technically this should be a [Generic]MtoShipmentShape
  // using hhg because it has all the props
  mtoShipment: HhgShipmentShape,
  currentResidence: AddressShape.isRequired,
  newDutyLocationAddress: SimpleAddressShape,
  updateMTOShipment: func.isRequired,
  serviceMember: shape({
    weight_allotment: shape({
      total_weight_self: number,
    }),
  }).isRequired,
  orders: OrdersShape,
};

CreateOrEditMtoShipment.defaultProps = {
  match: { isExact: false, params: { moveID: '' } },
  history: { goBack: () => {}, push: () => {} },
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
  orders: {},
};

function mapStateToProps(state, ownProps) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  const props = {
    serviceMember,
    orders: selectCurrentOrders(state) || {},
    mtoShipment: selectMTOShipmentById(state, ownProps.match.params.mtoShipmentId) || {},
    currentResidence: serviceMember?.residential_address || {},
    newDutyLocationAddress: selectCurrentOrders(state)?.new_duty_location?.address || {},
  };

  return props;
}

const mapDispatchToProps = {
  fetchCustomerData: fetchCustomerDataAction,
  updateMTOShipment: updateMTOShipmentAction,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(CreateOrEditMtoShipment));
