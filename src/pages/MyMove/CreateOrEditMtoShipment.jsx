import React, { Component } from 'react';
import { connect } from 'react-redux';
import { string, func, shape, number } from 'prop-types';
import { get } from 'lodash';

import MtoShipmentForm from 'components/Customer/MtoShipmentForm/MtoShipmentForm';
import {
  selectMTOShipmentById,
  createMTOShipment as createMTOShipmentAction,
  updateMTOShipment as updateMTOShipmentAction,
} from 'shared/Entities/modules/mtoShipments';
import { fetchCustomerData as fetchCustomerDataAction } from 'store/onboarding/actions';
import { HhgShipmentShape, HistoryShape, MatchShape, PageKeyShape, PageListShape } from 'types/customerShapes';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { selectActiveOrLatestOrdersFromEntities } from 'shared/Entities/modules/orders';
import { selectServiceMemberFromLoggedInUser } from 'shared/Entities/modules/serviceMembers';
import { AddressShape, SimpleAddressShape } from 'types/address';

export class CreateOrEditMtoShipment extends Component {
  componentDidMount() {
    const { fetchCustomerData } = this.props;
    fetchCustomerData();
  }

  // TODO: (in trailing PR) refactor edit component out of existence :)
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
      createMTOShipment,
      updateMTOShipment,
      serviceMember,
    } = this.props;
    const isCreatePage = match && match.path ? match.path.includes('start') : false;

    // wait until MTO shipment has loaded to render form
    if (isCreatePage || mtoShipment?.id) {
      return (
        <MtoShipmentForm
          match={match}
          history={history}
          pageList={pageList}
          pageKey={pageKey}
          mtoShipment={mtoShipment}
          selectedMoveType={selectedMoveType}
          isCreatePage={isCreatePage}
          currentResidence={currentResidence}
          newDutyStationAddress={newDutyStationAddress}
          createMTOShipment={createMTOShipment}
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
  createMTOShipment: func.isRequired,
  updateMTOShipment: func.isRequired,
  serviceMember: shape({
    weight_allotment: shape({
      total_weight_self: number,
    }),
  }).isRequired,
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
};

function mapStateToProps(state, ownProps) {
  const orders = selectActiveOrLatestOrdersFromEntities(state);
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  const props = {
    serviceMember,
    mtoShipment: selectMTOShipmentById(state, ownProps.match.params.mtoShipmentId),
    currentResidence: get(selectServiceMemberFromLoggedInUser(state), 'residential_address', {}),
    newDutyStationAddress: get(orders, 'new_duty_station.address', {}),
  };

  return props;
}

const mapDispatchToProps = {
  fetchCustomerData: fetchCustomerDataAction,
  createMTOShipment: createMTOShipmentAction,
  updateMTOShipment: updateMTOShipmentAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(CreateOrEditMtoShipment);
