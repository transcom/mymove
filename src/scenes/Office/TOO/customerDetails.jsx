import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { getEntitlements } from 'shared/Entities/modules/moveTaskOrders';

class CustomerDetails extends React.Component {
  componentDidMount() {
    const fakeMoveTaskOrderID = 'fc633380-3c5d-428b-a701-15812e0b0ba7';
    this.props.getEntitlements(fakeMoveTaskOrderID);
  }

  render() {
    const { entitlements } = this.props;
    console.log(entitlements);
    const NTS = entitlements && entitlements.nonTemporaryStorage ? 'Y' : 'N';
    const POV = entitlements && entitlements.privatelyOwnedVehicle ? 'Y' : 'N';
    return (
      <>
        <h1>Customer Deets Page</h1>
        {entitlements && (
          <>
            <h2>Customer Entitlements</h2>
            <dl>
              <dt>Weight Entitlement</dt>
              <dd>{entitlements.totalWeightSelf}</dd>
              <dt>SIT Entitlement</dt>
              <dd>{entitlements.storageInTransit}</dd>
              <dt>NTS Entitlement</dt>
              <dd>{NTS}</dd>
              <dt>POV Entitlement</dt>
              <dd>{POV}</dd>
            </dl>
          </>
        )}
      </>
    );
  }
}
const mapStateToProps = state => {
  const entitlements = get(state, 'entities.entitlements');
  return {
    entitlements: entitlements && Object.values(entitlements).length > 0 ? Object.values(entitlements)[0] : null,
  };
};
const mapDispatchToProps = dispatch => bindActionCreators({ getEntitlements }, dispatch);

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(CustomerDetails);
