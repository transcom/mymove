import React, { Component } from 'react';
import { connect } from 'react-redux';
import { arrayOf, string, shape, bool, func } from 'prop-types';

<<<<<<< HEAD
import {
  loadMTOShipments as loadMTOShipmentsAction,
  selectMTOShipmentForMTO,
} from 'shared/Entities/modules/mtoShipments';
=======
import styles from './HHGMove.module.scss';

>>>>>>> origin/MB-3484_pre-demo-cleanup
import HHGDetailsForm from 'components/Customer/HHGDetailsForm';
import '../../ghc_index.scss';

<<<<<<< HEAD
class HHGMoveSetup extends Component {
  componentDidMount() {
    const { match, loadMTOShipments } = this.props;
    loadMTOShipments(match.params.moveId);
  }

  render() {
    const { pageList, pageKey, match, push, mtoShipment } = this.props;

    return (
      <div>
        <h3>Now lets arrange details for the professional movers</h3>
        <HHGDetailsForm pageList={pageList} pageKey={pageKey} match={match} push={push} mtoShipment={mtoShipment} />
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
=======
const HHGMoveSetup = ({ pageList, pageKey, match, push }) => (
  <div className={styles.HHGMovePage}>
    <h3>Now let’s arrange details for the professional movers</h3>
    <HHGDetailsForm pageList={pageList} pageKey={pageKey} match={match} push={push} />
  </div>
);
>>>>>>> origin/MB-3484_pre-demo-cleanup

HHGMoveSetup.propTypes = {
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
  loadMTOShipments: func.isRequired,
  push: func.isRequired,
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

HHGMoveSetup.defaultProps = {
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

export default connect(mapStateToProps, mapDispatchToProps)(HHGMoveSetup);
