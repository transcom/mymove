import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { string, func } from 'prop-types';

import '../../ghc_index.scss';

import MtoShipmentForm from 'components/Customer/MtoShipmentForm/MtoShipmentForm';
import EditShipment from 'components/Customer/EditShipment';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { selectMTOShipmentForMTO } from 'shared/Entities/modules/mtoShipments';
import { fetchCustomerData as fetchCustomerDataAction } from 'store/onboarding/actions';
import { HhgShipmentShape, WizardPageShape } from 'types/customerShapes';

class CreateOrEditMtoShipment extends Component {
  componentDidMount() {
    const { fetchCustomerData } = this.props;
    fetchCustomerData();
  }

  // TODO: (in trailing PR) refactor edit component out of existence :)
  render() {
    const { wizardPage, mtoShipment, selectedMoveType } = this.props;
    const { match, history } = wizardPage;
    const isHHGFormPage = match.path === '/moves/:moveId/hhg-start';

    return (
      <div>
        {selectedMoveType === SHIPMENT_OPTIONS.HHG && !isHHGFormPage ? (
          <EditShipment mtoShipment={mtoShipment} match={match} history={history} />
        ) : (
          <MtoShipmentForm wizardPage={wizardPage} mtoShipment={mtoShipment} selectedMoveType={selectedMoveType} />
        )}
      </div>
    );
  }
}

function mapStateToProps(state, ownProps) {
  const props = {
    mtoShipment: selectMTOShipmentForMTO(state, ownProps.wizardPage.match.params.moveId),
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ fetchCustomerData: fetchCustomerDataAction }, dispatch);
}

CreateOrEditMtoShipment.propTypes = {
  wizardPage: WizardPageShape,
  fetchCustomerData: func.isRequired,
  selectedMoveType: string.isRequired,
  // technically this should be a [Generic]MtoShipmentShape
  // using hhg because it has all the props
  mtoShipment: HhgShipmentShape,
};

CreateOrEditMtoShipment.defaultProps = {
  wizardPage: {
    pageList: [],
    pageKey: '',
    match: { isExact: false, params: { moveID: '' } },
  },
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
};

export { CreateOrEditMtoShipment as CreateOrEditMtoShipmentComponent };
export default connect(mapStateToProps, mapDispatchToProps)(CreateOrEditMtoShipment);
