import React, { Component } from 'react';
import { connect } from 'react-redux';
import { string, func } from 'prop-types';

import '../../ghc_index.scss';

import EditShipment from 'components/Customer/EditShipment';
import HHGDetailsForm from 'components/Customer/MtoShipments/HHGDetailsForm';
import NTSDetailsForm from 'components/Customer/MtoShipments/NTSDetailsForm';
import NTSrDetailsForm from 'components/Customer/MtoShipments/NTSrDetailsForm';
import { HhgShipmentShape, wizardPageShape } from 'components/Customer/MtoShipments/propShapes';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import {
  loadMTOShipments as loadMTOShipmentsAction,
  selectMTOShipmentForMTO,
} from 'shared/Entities/modules/mtoShipments';

class CreateOrEditMtoShipment extends Component {
  componentDidMount() {
    const { match, loadMTOShipments } = this.props;
    loadMTOShipments(match.params.moveId);
  }

  // TODO: (in trailing PR) refactor edit component out of existence :)
  render() {
    const { match, pageList, pageKey, history, mtoShipment, selectedMoveType } = this.props;
    const isHHGFormPage = match.path === '/moves/:moveId/hhg-start';

    return (
      <div>
        {selectedMoveType === SHIPMENT_OPTIONS.HHG && (
          <div>
            {isHHGFormPage && (
              <HHGDetailsForm
                pageList={pageList}
                pageKey={pageKey}
                match={match}
                push={history.push}
                mtoShipment={mtoShipment}
              />
            )}
            {!isHHGFormPage && <EditShipment mtoShipment={mtoShipment} match={match} history={history} />}
          </div>
        )}
        {selectedMoveType === SHIPMENT_OPTIONS.NTS && (
          <NTSDetailsForm
            match={match}
            pageKey={pageKey}
            pageList={pageList}
            push={history.push}
            mtoShipment={mtoShipment}
          />
        )}
        {selectedMoveType === SHIPMENT_OPTIONS.NTSR && (
          <NTSrDetailsForm
            match={match}
            pageKey={pageKey}
            pageList={pageList}
            push={history.push}
            mtoShipment={mtoShipment}
          />
        )}
      </div>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  const props = {
    mtoShipment: selectMTOShipmentForMTO(state, ownProps.match.params.moveId),
  };
  return props;
};

const mapDispatchToProps = {
  loadMTOShipments: loadMTOShipmentsAction,
};

CreateOrEditMtoShipment.propTypes = {
  ...wizardPageShape,
  loadMTOShipments: func.isRequired,
  selectedMoveType: string.isRequired,
  // technically this should be a [Generic]MtoShipmentShape
  // using hhg because it has all the props
  mtoShipment: HhgShipmentShape,
};

CreateOrEditMtoShipment.defaultProps = {
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

export default connect(mapStateToProps, mapDispatchToProps)(CreateOrEditMtoShipment);
