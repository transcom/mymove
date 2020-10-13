import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { string, func } from 'prop-types';

import '../../ghc_index.scss';

import EditShipment from 'components/Customer/EditShipment';
import HHGDetailsForm from 'components/Customer/MtoShipments/HHGDetailsForm';
import NTSDetailsForm from 'components/Customer/MtoShipments/NTSDetailsForm';
import NTSrDetailsForm from 'components/Customer/MtoShipments/NTSrDetailsForm';
import { HhgShipmentShape, WizardPageShape } from 'components/Customer/MtoShipments/propShapes';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import {
  loadMTOShipments as loadMTOShipmentsAction,
  selectMTOShipmentForMTO,
} from 'shared/Entities/modules/mtoShipments';

class CreateOrEditMtoShipment extends Component {
  componentDidMount() {
    const { wizardPage, loadMTOShipments } = this.props;
    loadMTOShipments(wizardPage.match.params.moveId);
  }

  // TODO: (in trailing PR) refactor edit component out of existence :)
  render() {
    const { wizardPage, mtoShipment, selectedMoveType } = this.props;
    const { match, history } = wizardPage;
    const isHHGFormPage = match.path === '/moves/:moveId/hhg-start';

    return (
      <div>
        {selectedMoveType === SHIPMENT_OPTIONS.HHG && (
          <div>
            {isHHGFormPage && <HHGDetailsForm wizardPage={wizardPage} mtoShipment={mtoShipment} />}
            {!isHHGFormPage && <EditShipment mtoShipment={mtoShipment} match={match} history={history} />}
          </div>
        )}
        {selectedMoveType === SHIPMENT_OPTIONS.NTS && (
          <NTSDetailsForm wizardPage={wizardPage} mtoShipment={mtoShipment} />
        )}
        {selectedMoveType === SHIPMENT_OPTIONS.NTSR && (
          <NTSrDetailsForm wizardPage={wizardPage} mtoShipment={mtoShipment} />
        )}
      </div>
    );
  }
}

function mapStateToProps(state, ownProps) {
  const props = {
    mtoShipment: selectMTOShipmentForMTO(state, ownProps.match.params.moveId),
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadMTOShipments: loadMTOShipmentsAction }, dispatch);
}

CreateOrEditMtoShipment.propTypes = {
  wizardPage: WizardPageShape,
  loadMTOShipments: func.isRequired,
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
