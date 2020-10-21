import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { string, func } from 'prop-types';

import '../../ghc_index.scss';

import MtoShipmentForm from 'components/Customer/MtoShipmentForm/MtoShipmentForm';
import { selectMTOShipmentForMTO } from 'shared/Entities/modules/mtoShipments';
import { fetchCustomerData as fetchCustomerDataAction } from 'store/onboarding/actions';
import { HhgShipmentShape, HistoryShape, MatchShape, PageKeyShape, PageListShape } from 'types/customerShapes';

class CreateOrEditMtoShipment extends Component {
  componentDidMount() {
    const { fetchCustomerData } = this.props;
    fetchCustomerData();
  }

  // TODO: (in trailing PR) refactor edit component out of existence :)
  render() {
    const { match, history, pageList, pageKey, mtoShipment, selectedMoveType } = this.props;
    const isCreatePage = match && match.path ? match.path.includes('start') : false;

    return (
      <MtoShipmentForm
        match={match}
        history={history}
        pageList={pageList}
        pageKey={pageKey}
        mtoShipment={mtoShipment}
        selectedMoveType={selectedMoveType}
        isCreatePage={isCreatePage}
      />
    );
  }
}

function mapStateToProps(state, ownProps) {
  const props = {
    mtoShipment: selectMTOShipmentForMTO(state, ownProps.match?.params.moveId),
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ fetchCustomerData: fetchCustomerDataAction }, dispatch);
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
};

export { CreateOrEditMtoShipment as CreateOrEditMtoShipmentComponent };
export default connect(mapStateToProps, mapDispatchToProps)(CreateOrEditMtoShipment);
