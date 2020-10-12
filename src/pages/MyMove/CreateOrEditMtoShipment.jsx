import React, { Component } from 'react';
import { connect } from 'react-redux';
import { arrayOf, string, shape, bool, func } from 'prop-types';

import '../../ghc_index.scss';
import HHGShipmentSetup from './HHGShipmentSetup';

import NTSDetailsForm from 'components/Customer/MtoShipments/NTSDetailsForm';
import NTSrDetailsForm from 'components/Customer/MtoShipments/NTSrDetailsForm';
import { HhgShipmentShape } from 'components/Customer/MtoShipments/propShapes';
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

  render() {
    const { match, pageList, pageKey, history, mtoShipment, loadMTOShipments } = this.props;
    return (
      <div>
        {mtoShipment.ShipmentType === SHIPMENT_OPTIONS.HHG && (
          <HHGShipmentSetup
            match={match}
            pageKey={pageKey}
            pageList={pageList}
            history={history}
            loadMTOShipments={loadMTOShipments}
            mtoShipment={mtoShipment}
          />
        )}
        {mtoShipment.ShipmentType === SHIPMENT_OPTIONS.NTS && (
          <NTSDetailsForm
            match={match}
            pageKey={pageKey}
            pageList={pageList}
            history={history}
            mtoShipment={mtoShipment}
          />
        )}
        {mtoShipment.ShipmentType === SHIPMENT_OPTIONS.NTSR && (
          <NTSrDetailsForm
            match={match}
            pageKey={pageKey}
            pageList={pageList}
            history={history}
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
  pageList: arrayOf(string),
  pageKey: string,
  match: shape({
    isExact: bool.isRequired,
    params: shape({
      moveId: string.isRequired,
    }),
    path: string.isRequired,
    url: string.isRequired,
  }).isRequired,
  history: shape({
    goBack: func.isRequired,
    push: func.isRequired,
  }).isRequired,
  loadMTOShipments: func.isRequired,
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
  pageList: [],
  pageKey: '',
};

export default connect(mapStateToProps, mapDispatchToProps)(CreateOrEditMtoShipment);
