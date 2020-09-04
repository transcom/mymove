import React, { Component } from 'react';
import { connect } from 'react-redux';
import { arrayOf, string, shape, bool, func } from 'prop-types';

import EditShipment from '../../components/Customer/EditShipment';

import {
  loadMTOShipments as loadMTOShipmentsAction,
  selectMTOShipmentForMTO,
} from 'shared/Entities/modules/mtoShipments';
import HHGDetailsForm from 'components/Customer/HHGDetailsForm';
import '../../ghc_index.scss';

class HHGShipmentSetup extends Component {
  componentDidMount() {
    const { match, loadMTOShipments } = this.props;
    loadMTOShipments(match.params.moveId);
  }

  render() {
    const { pageList, pageKey, match, history, mtoShipment } = this.props;
    const isEditShipmentPage = match.path === '/moves/:moveId/edit-shipment';
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

HHGShipmentSetup.propTypes = {
  pageList: arrayOf(string).isRequired,
  pageKey: string.isRequired,
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
  mtoShipment: shape({
    agents: arrayOf(
      shape({
        firstName: string,
        lastName: string,
        phone: string,
        email: string,
        agentType: string,
      }),
    ),
    customerRemarks: string,
    requestedPickupDate: string,
    requestedDeliveryDate: string,
    pickupAddress: shape({
      city: string,
      postal_code: string,
      state: string,
      street_address_1: string,
    }),
    destinationAddress: shape({
      city: string,
      postal_code: string,
      state: string,
      street_address_1: string,
    }),
  }),
};

HHGShipmentSetup.defaultProps = {
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

export default connect(mapStateToProps, mapDispatchToProps)(HHGShipmentSetup);
