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

  error = () => {
    const { error } = this.state;
    return error;
  };

  isCreatorActionable = () => {
    const { isCreatorActionable } = this.state;
    return isCreatorActionable;
  };

  sitEntitlement = () => {
    const { sitEntitlement } = this.props;
    return sitEntitlement;
  };

  render() {
    return (
      <div className="storage-in-transit-panel">
        <BasicPanel title="Storage in Transit (SIT)">
          {this.error() && (
            <Alert type="error" heading="Oops, something went wrong!" onRemove={this.closeError}>
              <span className="warning--header">Please refresh the page and try again.</span>
            </Alert>
          )}
          <div className="column-subhead">Entitlement: {this.sitEntitlement()} days</div>
          {this.isCreatorActionable() && <Creator />}
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
