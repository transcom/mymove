import React, { Component } from 'react';
import PropTypes from 'prop-types';

import { connect } from 'react-redux';

import BasicPanel from 'shared/BasicPanel';
import Alert from 'shared/Alert';
import Creator from 'shared/StorageInTransit/Creator';

export class StorageInTransitPanel extends Component {
  constructor() {
    super();
    this.state = {
      isCreatorActionable: true,
      error: null,
    };
  }

  closeError = () => {
    this.setState({ error: null });
  };

  render() {
    const { sitEntitlement } = this.props;
    const { error, isCreatorActionable } = this.state;
    const daysUsed = 0; // placeholder
    const daysRemaining = sitEntitlement - daysUsed;
    return (
      <div className="storage-in-transit-panel">
        <BasicPanel title="Storage in Transit (SIT)">
          {error && (
            <Alert type="error" heading="Oops, something went wrong!" onRemove={this.closeError}>
              <span className="warning--header">Please refresh the page and try again.</span>
            </Alert>
          )}
          <div className="column-subhead">
            Entitlement: {sitEntitlement} days <span className="unbold">({daysRemaining} remaining)</span>
          </div>
          {isCreatorActionable && <Creator />}
        </BasicPanel>
      </div>
    );
  }
}

StorageInTransitPanel.propTypes = {
  sitRequests: PropTypes.array,
  shipmentId: PropTypes.string,
  sitEntitlement: PropTypes.number,
};

/*
function mapStateToProps(state, ownProps) {
  return {
    sitRequests: selectSortedSitRequests(state, ownProps.shipmentId),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { createSitRequest, deleteSitRequest, approveSitRequest, updateSitRequest },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(StorageInTransitPanel);
*/
export default connect(null, null)(StorageInTransitPanel);
