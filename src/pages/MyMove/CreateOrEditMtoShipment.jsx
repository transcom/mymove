import React, { Component } from 'react';
import { connect } from 'react-redux';
import { arrayOf, string, shape, bool, func } from 'prop-types';


import {
  loadMTOShipments as loadMTOShipmentsAction,
  selectMTOShipmentForMTO,
} from 'shared/Entities/modules/mtoShipments';
import EditShipment from 'components/Customer/EditShipment';
import HHGDetailsForm from 'components/Customer/MtoShipments/HHGDetailsForm';
import NTSDetailsForm from 'components/Customer/MtoShipments/NTSDetailsForm';
import NTSrDetailsForm from 'components/Customer/MtoShipments/NTSrDetailsForm';
import '../../ghc_index.scss';

class CreateOrEditMtoShipment extends Component {
  componentDidMount() {
    const { match, loadMTOShipments } = this.props;
    loadMTOShipments(match.params.moveId);
  }

  render() {
    const { pageList, pageKey, match, history, mtoShipment } = this.props;
    const isEditShipmentPage = match.path === '/moves/:moveId/mto-shipments/:mtoShipmentId/edit-shipment';
    const isHHGFormPage = match.path === '/moves/:moveId/hhg-start';
    return (
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
        {mtoShipment.ShipmentType === SHIPMENT_OPTIONS.NTS && (
          <NTSDetailsForm />
        )}
        {mtoShipment.ShipmentType === SHIPMENT_OPTIONS.NTSR && (
          <NTSrDetailsForm />
        )}
        {isEditShipmentPage && <EditShipment mtoShipment={mtoShipment} match={match} history={history} />}
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
  mtoShipment: ,
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
