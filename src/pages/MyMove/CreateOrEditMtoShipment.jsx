import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bool, string, func, shape, number } from 'prop-types';

import MtoShipmentForm from 'components/Customer/MtoShipmentForm/MtoShipmentForm';
import { selectMTOShipmentById } from 'shared/Entities/modules/mtoShipments';
import { updateMTOShipment as updateMTOShipmentAction } from 'store/entities/actions';
import { fetchCustomerData as fetchCustomerDataAction } from 'store/onboarding/actions';
import { HhgShipmentShape, HistoryShape, MatchShape, PageKeyShape, PageListShape } from 'types/customerShapes';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { selectServiceMemberFromLoggedInUser, selectCurrentOrders } from 'store/entities/selectors';
import { AddressShape, SimpleAddressShape } from 'types/address';

export class CreateOrEditMtoShipment extends Component {
  componentDidMount() {
    const { fetchCustomerData } = this.props;
    fetchCustomerData();
  }

  render() {
    const {
      match,
      history,
      pageList,
      pageKey,
      mtoShipment,
      selectedMoveType,
      currentResidence,
      newDutyStationAddress,
      updateMTOShipment,
      serviceMember,
      isCreate,
    } = this.props;

    // wait until MTO shipment has loaded to render form
    if (isCreate || mtoShipment?.id) {
      return (
        <MtoShipmentForm
          match={match}
          history={history}
          pageList={pageList}
          pageKey={pageKey}
          mtoShipment={mtoShipment}
          selectedMoveType={selectedMoveType}
          isCreatePage={isCreate}
          currentResidence={currentResidence}
          newDutyStationAddress={newDutyStationAddress}
          updateMTOShipment={updateMTOShipment}
          serviceMember={serviceMember}
        />
      );
    }

    return <LoadingPlaceholder />;
  }
}

CreateOrEditMtoShipment.propTypes = {
  match: MatchShape,
  history: HistoryShape,
  pageList: PageListShape,
  pageKey: PageKeyShape,
  fetchCustomerData: func.isRequired,
  selectedMoveType: string.isRequired,
  // technically this should be a [Generic]MtoShipmentShape
  // using hhg because it has all the props
  mtoShipment: HhgShipmentShape,
  currentResidence: AddressShape.isRequired,
  newDutyStationAddress: SimpleAddressShape,
  updateMTOShipment: func.isRequired,
  serviceMember: shape({
    weight_allotment: shape({
      total_weight_self: number,
    }),
  }).isRequired,
  isCreate: bool,
};

CreateOrEditMtoShipment.defaultProps = {
  pageList: [],
  pageKey: '',
  match: { isExact: false, params: { moveID: '' } },
  history: { goBack: () => {}, push: () => {} },
  mtoShipment: {
    customerRemarks: '',
    requestedPickupDate: '',
    requestedDeliveryDate: '',
    destinationAddress: {
      city: '',
      postal_code: '',
      state: '',
      street_address_1: '',
    },
  },
  newDutyStationAddress: {
    city: '',
    state: '',
    postal_code: '',
  },
  isCreate: false,
};

function mapStateToProps(state, ownProps) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  const props = {
    serviceMember,
    mtoShipment: selectMTOShipmentById(state, ownProps.match.params.mtoShipmentId),
    currentResidence: serviceMember?.residential_address || {},
    newDutyStationAddress: selectCurrentOrders(state)?.new_duty_station?.address || {},
  };

  return props;
}

const mapDispatchToProps = {
  fetchCustomerData: fetchCustomerDataAction,
  updateMTOShipment: updateMTOShipmentAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(CreateOrEditMtoShipment);
